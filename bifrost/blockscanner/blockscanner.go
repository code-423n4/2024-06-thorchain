package blockscanner

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	btypes "gitlab.com/thorchain/thornode/bifrost/blockscanner/types"
	"gitlab.com/thorchain/thornode/bifrost/metrics"
	"gitlab.com/thorchain/thornode/bifrost/thorclient"
	"gitlab.com/thorchain/thornode/bifrost/thorclient/types"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/config"
	"gitlab.com/thorchain/thornode/constants"
)

// BlockScannerFetcher define the methods a block scanner need to implement
type BlockScannerFetcher interface {
	// FetchMemPool scan the mempool
	FetchMemPool(height int64) (types.TxIn, error)
	// FetchTxs scan block with the given height
	FetchTxs(fetchHeight, chainHeight int64) (types.TxIn, error)
	// GetHeight return current block height
	GetHeight() (int64, error)
}

type Block struct {
	Height int64
	Txs    []string
}

// BlockScanner is used to discover block height
type BlockScanner struct {
	cfg             config.BifrostBlockScannerConfiguration
	logger          zerolog.Logger
	wg              *sync.WaitGroup
	scanChan        chan int64
	stopChan        chan struct{}
	scannerStorage  ScannerStorage
	metrics         *metrics.Metrics
	previousBlock   int64
	globalTxsQueue  chan types.TxIn
	errorCounter    *prometheus.CounterVec
	thorchainBridge thorclient.ThorchainBridge
	chainScanner    BlockScannerFetcher
	healthy         *atomic.Bool
}

// NewBlockScanner create a new instance of BlockScanner
func NewBlockScanner(cfg config.BifrostBlockScannerConfiguration, scannerStorage ScannerStorage, m *metrics.Metrics, thorchainBridge thorclient.ThorchainBridge, chainScanner BlockScannerFetcher) (*BlockScanner, error) {
	var err error
	if scannerStorage == nil {
		return nil, errors.New("scannerStorage is nil")
	}
	if m == nil {
		return nil, errors.New("metrics instance is nil")
	}
	if thorchainBridge == nil {
		return nil, errors.New("thorchain bridge is nil")
	}

	logger := log.Logger.With().Str("module", "blockscanner").Str("chain", cfg.ChainID.String()).Logger()
	scanner := &BlockScanner{
		cfg:             cfg,
		logger:          logger,
		wg:              &sync.WaitGroup{},
		stopChan:        make(chan struct{}),
		scanChan:        make(chan int64),
		scannerStorage:  scannerStorage,
		metrics:         m,
		errorCounter:    m.GetCounterVec(metrics.CommonBlockScannerError),
		thorchainBridge: thorchainBridge,
		chainScanner:    chainScanner,
		healthy:         &atomic.Bool{},
	}

	scanner.previousBlock, err = scanner.FetchLastHeight()
	logger.Info().Int64("block height", scanner.previousBlock).Msg("block scanner last fetch height")
	return scanner, err
}

// IsHealthy return if the block scanner is healthy or not
func (b *BlockScanner) IsHealthy() bool {
	return b.healthy.Load()
}

func (b *BlockScanner) PreviousHeight() int64 {
	return atomic.LoadInt64(&b.previousBlock)
}

// GetMessages return the channel
func (b *BlockScanner) GetMessages() <-chan int64 {
	return b.scanChan
}

// Start block scanner
func (b *BlockScanner) Start(globalTxsQueue chan types.TxIn) {
	b.globalTxsQueue = globalTxsQueue
	currentPos, err := b.scannerStorage.GetScanPos()
	if err != nil {
		b.logger.Error().Err(err).Msgf("fail to get current block scan pos, %s will start from %d", b.cfg.ChainID, b.previousBlock)
	} else if currentPos > b.previousBlock {
		b.previousBlock = currentPos
	}
	b.wg.Add(2)
	go b.scanBlocks()
	go b.scanMempool()
}

func (b *BlockScanner) scanMempool() {
	b.logger.Info().Msg("start to scan mempool")
	defer b.logger.Info().Msg("stop scan mempool")
	defer b.wg.Done()

	if !b.cfg.ScanMemPool {
		b.logger.Info().Msg("mempool scan is disabled")
		return
	}

	for {
		select {
		case <-b.stopChan:
			return
		default:
			// mempool scan will continue even the chain get halted , thus the network can still aware of outbound transaction
			// during chain halt
			preBlockHeight := atomic.LoadInt64(&b.previousBlock)
			currentBlock := preBlockHeight + 1
			txInMemPool, err := b.chainScanner.FetchMemPool(currentBlock)
			if err != nil {
				b.logger.Error().Err(err).Msg("fail to fetch MemPool")
			}
			if len(txInMemPool.TxArray) > 0 {
				select {
				case <-b.stopChan:
					return
				case b.globalTxsQueue <- txInMemPool:
				}
			} else {
				// nothing in the mempool or for some chain like BNB & ETH, which doesn't need to scan
				// mempool , back off here
				time.Sleep(constants.ThorchainBlockTime)
			}
		}
	}
}

// Checks current mimir settings to determine if the current chain is paused
// either globally or specifically
func (b *BlockScanner) isChainPaused() bool {
	var haltHeight, solvencyHaltHeight, nodeHaltHeight, thorHeight int64

	// Check if chain has been halted via mimir
	haltHeight, err := b.thorchainBridge.GetMimir(fmt.Sprintf("Halt%sChain", b.cfg.ChainID))
	if err != nil {
		b.logger.Error().Err(err).Msgf("fail to get mimir setting %s", fmt.Sprintf("Halt%sChain", b.cfg.ChainID))
	}
	// Check if chain has been halted by auto solvency checks
	solvencyHaltHeight, err = b.thorchainBridge.GetMimir(fmt.Sprintf("SolvencyHalt%sChain", b.cfg.ChainID))
	if err != nil {
		b.logger.Error().Err(err).Msgf("fail to get mimir %s", fmt.Sprintf("SolvencyHalt%sChain", b.cfg.ChainID))
	}
	// Check if all chains halted globally
	globalHaltHeight, err := b.thorchainBridge.GetMimir("HaltChainGlobal")
	if err != nil {
		b.logger.Error().Err(err).Msg("fail to get mimir setting HaltChainGlobal")
	}
	if globalHaltHeight > haltHeight {
		haltHeight = globalHaltHeight
	}
	// Check if a node paused all chains
	nodeHaltHeight, err = b.thorchainBridge.GetMimir("NodePauseChainGlobal")
	if err != nil {
		b.logger.Error().Err(err).Msg("fail to get mimir setting NodePauseChainGlobal")
	}
	thorHeight, err = b.thorchainBridge.GetBlockHeight()
	if err != nil {
		b.logger.Error().Err(err).Msg("fail to get THORChain block height")
	}

	if nodeHaltHeight > 0 && thorHeight < nodeHaltHeight {
		haltHeight = 1
	}

	return (haltHeight > 0 && thorHeight > haltHeight) || (solvencyHaltHeight > 0 && thorHeight > solvencyHaltHeight)
}

// scanBlocks
func (b *BlockScanner) scanBlocks() {
	b.logger.Debug().Msg("start to scan blocks")
	defer b.logger.Debug().Msg("stop scan blocks")
	defer b.wg.Done()

	lastMimirCheck := time.Now().Add(-constants.ThorchainBlockTime)
	isChainPaused := false

	// start up to grab those blocks
	for {
		select {
		case <-b.stopChan:
			return
		default:
			preBlockHeight := atomic.LoadInt64(&b.previousBlock)
			currentBlock := preBlockHeight + 1
			// check if mimir has disabled this chain
			if time.Since(lastMimirCheck) >= constants.ThorchainBlockTime {
				isChainPaused = b.isChainPaused()
				lastMimirCheck = time.Now()
			}

			// Chain is paused, mark as unhealthy
			if isChainPaused {
				b.healthy.Store(false)
				time.Sleep(constants.ThorchainBlockTime)
				continue
			}

			chainHeight, err := b.chainScanner.GetHeight()
			if err != nil {
				b.logger.Error().Err(err).Msg("fail to get chain block height")
				time.Sleep(b.cfg.BlockHeightDiscoverBackoff)
				continue
			}
			if chainHeight < currentBlock {
				time.Sleep(b.cfg.BlockHeightDiscoverBackoff)
				continue
			}
			txIn, err := b.chainScanner.FetchTxs(currentBlock, chainHeight)
			if err != nil {
				// don't log an error if its because the block doesn't exist yet
				if !errors.Is(err, btypes.ErrUnavailableBlock) {
					b.logger.Error().Err(err).Int64("block height", currentBlock).Msg("fail to get RPCBlock")
					b.healthy.Store(false)
				}
				time.Sleep(b.cfg.BlockHeightDiscoverBackoff)
				continue
			}

			// determine how often we print a info log line for scanner
			// progress. General goal is about once per minute
			ms := b.cfg.ChainID.ApproximateBlockMilliseconds()
			mod := (60_000 + ms - 1) / ms
			// enable this one , so we could see how far it is behind
			if currentBlock%mod == 0 || !b.healthy.Load() {
				b.logger.Info().
					Int64("block height", currentBlock).
					Int("txs", len(txIn.TxArray)).
					Int64("gap", chainHeight-currentBlock).
					Bool("healthy", b.healthy.Load()).
					Msg("scan block")
			}
			atomic.AddInt64(&b.previousBlock, 1)

			// consider 3 blocks or the configured lag time behind tip as healthy
			lagDuration := time.Duration((chainHeight-currentBlock)*ms) * time.Millisecond
			if chainHeight-currentBlock <= 3 || lagDuration < b.cfg.MaxHealthyLag {
				b.healthy.Store(true)
			} else {
				b.healthy.Store(false)
			}
			b.logger.Debug().Msgf("the gap is %d , healthy: %+v", chainHeight-currentBlock, b.healthy.Load())

			b.metrics.GetCounter(metrics.TotalBlockScanned).Inc()
			if len(txIn.TxArray) > 0 {
				select {
				case <-b.stopChan:
					return
				case b.globalTxsQueue <- txIn:
				}
			}
			if err = b.scannerStorage.SetScanPos(b.previousBlock); err != nil {
				b.logger.Error().Err(err).Msg("fail to save block scan pos")
				// alert!!
				continue
			}
		}
	}
}

// FetchLastHeight determines the height to start scanning:
//  1. Use the config start height if set.
//  2. If last consensus inbound height (lastblock) is available:
//     a) Use local scanner storage height if available, up to the max lag from lastblock.
//     b) Otherwise, use lastblock.
//  3. Otherwise, use local scanner storage height if available.
//  4. Otherwise, use the last height from the chain itself.
func (b *BlockScanner) FetchLastHeight() (int64, error) {
	// get scanner storage height
	currentPos, _ := b.scannerStorage.GetScanPos() // ignore error

	// 1. Use the config start height if set.
	if b.cfg.StartBlockHeight > 0 {
		return b.cfg.StartBlockHeight, nil
	}

	// wait for thorchain to be caught up first
	if err := b.thorchainBridge.WaitToCatchUp(); err != nil {
		return 0, err
	}

	if b.thorchainBridge != nil {
		var height int64
		if b.cfg.ChainID.Equals(common.THORChain) {
			height, _ = b.thorchainBridge.GetBlockHeight()
		} else {
			height, _ = b.thorchainBridge.GetLastObservedInHeight(b.cfg.ChainID)
		}
		if height > 0 {

			// 2.a) Use local scanner storage height if available, up to the max lag from lastblock.
			if currentPos > 0 {
				// calculate the max lag
				maxLagBlocks := b.cfg.MaxResumeBlockLag.Milliseconds() / b.cfg.ChainID.ApproximateBlockMilliseconds()

				// return the position up to the max block lag behind the consensus height
				if height <= currentPos+maxLagBlocks {
					return currentPos, nil
				} else {
					return height - maxLagBlocks, nil
				}
			}

			// 2.b) Otherwise, use lastblock.
			return height, nil
		}
	}

	//  3. Otherwise, use local scanner storage height if available.
	if currentPos > 0 {
		return currentPos, nil
	}

	//  4. Otherwise, use the last height from the chain itself.
	return b.chainScanner.GetHeight()
}

func (b *BlockScanner) Stop() {
	b.logger.Debug().Msg("receive stop request")
	defer b.logger.Debug().Msg("common block scanner stopped")
	close(b.stopChan)
	b.wg.Wait()
}
