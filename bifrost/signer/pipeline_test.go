package signer

import (
	"sync"

	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/thornode/bifrost/thorclient/types"
	"gitlab.com/thorchain/thornode/common"
	ttypes "gitlab.com/thorchain/thornode/x/thorchain/types"
	. "gopkg.in/check.v1"
)

////////////////////////////////////////////////////////////////////////////////////////
// Init
////////////////////////////////////////////////////////////////////////////////////////

func init() {
	// add caller to logger for debugging
	log.Logger = log.With().Caller().Logger()
}

////////////////////////////////////////////////////////////////////////////////////////
// mockPipelineSigner
////////////////////////////////////////////////////////////////////////////////////////

type mockPipelineSigner struct {
	sync.Mutex
	stopped          bool
	storageListItems []TxOutStoreItem
	processed        []TxOutStoreItem
}

func (m *mockPipelineSigner) isStopped() bool {
	return m.stopped
}

func (m *mockPipelineSigner) storageList() []TxOutStoreItem {
	return m.storageListItems
}

func (m *mockPipelineSigner) processTransaction(item TxOutStoreItem) {
	m.Lock()
	defer m.Unlock()

	// set processed
	m.processed = append(m.processed, item)

	// remove from storage list
	for i, tx := range m.storageListItems {
		if tx.TxOutItem.Equals(item.TxOutItem) {
			m.storageListItems = append(m.storageListItems[:i], m.storageListItems[i+1:]...)
			break
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Test Data
////////////////////////////////////////////////////////////////////////////////////////

var (
	vault1 = ttypes.GetRandomPubKey()
	vault2 = ttypes.GetRandomPubKey()

	tosis = []TxOutStoreItem{
		{
			TxOutItem: types.TxOutItem{
				Chain:       common.BTCChain,
				ToAddress:   ttypes.GetRandomBTCAddress(),
				VaultPubKey: vault1,
			},
		},
		{
			TxOutItem: types.TxOutItem{
				Chain:       common.BTCChain,
				ToAddress:   ttypes.GetRandomBTCAddress(),
				VaultPubKey: vault2,
			},
		},
		{ // same vault/chain as previous, should not happen concurrent
			TxOutItem: types.TxOutItem{
				Chain:       common.BTCChain,
				ToAddress:   ttypes.GetRandomBTCAddress(),
				VaultPubKey: vault2,
			},
		},
		{
			TxOutItem: types.TxOutItem{
				Chain:       common.ETHChain,
				ToAddress:   ttypes.GetRandomETHAddress(),
				VaultPubKey: vault1,
			},
		},
		{
			TxOutItem: types.TxOutItem{
				Chain:       common.ETHChain,
				ToAddress:   ttypes.GetRandomETHAddress(),
				VaultPubKey: vault2,
			},
		},
	}
)

////////////////////////////////////////////////////////////////////////////////////////
// PipelineSigner
////////////////////////////////////////////////////////////////////////////////////////

type PipelineSuite struct{}

var _ = Suite(&PipelineSuite{})

func (s *PipelineSuite) TestPipelineInit(c *C) {
	// valid
	for i := 1; i < 3; i++ {
		pipeline, err := newPipeline(int64(i))
		c.Assert(pipeline, NotNil)
		c.Assert(err, IsNil)
	}

	// invalid
	for i := -1; i < 1; i++ {
		pipeline, err := newPipeline(int64(i))
		c.Assert(pipeline, IsNil)
		c.Assert(err, NotNil)
	}
}

func (s *PipelineSuite) TestPipelineSequential(c *C) {
	pipeline, err := newPipeline(1)
	c.Assert(pipeline, NotNil)
	c.Assert(err, IsNil)

	// mocks
	bridge := fakeBridge{nil}
	mockSigner := &mockPipelineSigner{
		storageListItems: append([]TxOutStoreItem{}, tosis...),
	}

	// spawn one signing
	mockSigner.Lock() // prevent signing completion
	pipeline.SpawnSignings(mockSigner, bridge)
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 1)
	c.Assert(len(mockSigner.storageListItems), Equals, 5)
	c.Assert(len(mockSigner.processed), Equals, 0)

	// attempting another should be noop since semaphore is taken
	pipeline.SpawnSignings(mockSigner, bridge)
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 1)
	c.Assert(len(mockSigner.storageListItems), Equals, 5)
	c.Assert(len(mockSigner.processed), Equals, 0)

	// release lock and first signing should complete
	mockSigner.Unlock()
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 4)
	c.Assert(len(mockSigner.processed), Equals, 1)
	c.Assert(mockSigner.processed[0].TxOutItem.Equals(tosis[0].TxOutItem), Equals, true)

	// complete remaining signings
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 0)
	c.Assert(len(mockSigner.processed), Equals, 5)
	c.Assert(mockSigner.processed[1].TxOutItem.Equals(tosis[1].TxOutItem), Equals, true)
	c.Assert(mockSigner.processed[2].TxOutItem.Equals(tosis[2].TxOutItem), Equals, true)
	c.Assert(mockSigner.processed[3].TxOutItem.Equals(tosis[3].TxOutItem), Equals, true)
}

func (s *PipelineSuite) TestPipelineSequentialStop(c *C) {
	pipeline, err := newPipeline(1)
	c.Assert(pipeline, NotNil)
	c.Assert(err, IsNil)

	// mocks
	bridge := fakeBridge{nil}
	mockSigner := &mockPipelineSigner{
		storageListItems: append([]TxOutStoreItem{}, tosis...),
	}

	// spawn one signing
	mockSigner.Lock()
	pipeline.SpawnSignings(mockSigner, bridge)

	// stop the signer
	mockSigner.stopped = true

	// release lock and first signing should complete
	mockSigner.Unlock()
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 4)
	c.Assert(len(mockSigner.processed), Equals, 1)

	// no more signings should be spawned
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 4)
	c.Assert(len(mockSigner.processed), Equals, 1)
}

func (s *PipelineSuite) TestPipelineConcurrent(c *C) {
	pipeline, err := newPipeline(10)
	c.Assert(pipeline, NotNil)
	c.Assert(err, IsNil)

	// mocks
	bridge := fakeBridge{nil}
	mockSigner := &mockPipelineSigner{
		storageListItems: append([]TxOutStoreItem{}, tosis...),
	}

	// spawn signings - only 4/5 of test data should be concurrent
	mockSigner.Lock()
	pipeline.SpawnSignings(mockSigner, bridge)
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 4)
	c.Assert(len(mockSigner.storageListItems), Equals, 5)
	c.Assert(len(mockSigner.processed), Equals, 0)

	// release lock and all signings should complete
	mockSigner.Unlock()
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 1)
	c.Assert(len(mockSigner.processed), Equals, 4)

	// the remaining signing should be the 3rd item
	c.Assert(mockSigner.storageListItems[0].TxOutItem.Equals(tosis[2].TxOutItem), Equals, true)

	// complete remaining signings
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.processed), Equals, 5)
	c.Assert(len(mockSigner.storageListItems), Equals, 0)
}

func (s *PipelineSuite) TestPipelineConcurrentStop(c *C) {
	pipeline, err := newPipeline(10)
	c.Assert(pipeline, NotNil)
	c.Assert(err, IsNil)

	// mocks
	bridge := fakeBridge{nil}
	mockSigner := &mockPipelineSigner{
		storageListItems: append([]TxOutStoreItem{}, tosis...),
	}

	// spawn signings - only 4/5 of test data should be concurrent
	mockSigner.Lock()
	pipeline.SpawnSignings(mockSigner, bridge)
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 4)
	c.Assert(len(mockSigner.storageListItems), Equals, 5)
	c.Assert(len(mockSigner.processed), Equals, 0)

	// stop the signer
	mockSigner.stopped = true

	// release lock and all signings should complete
	mockSigner.Unlock()
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 1)
	c.Assert(len(mockSigner.processed), Equals, 4)

	// no more signings should be spawned
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 1)
	c.Assert(len(mockSigner.processed), Equals, 4)

	// the remaining signing should be the 3rd item
	c.Assert(mockSigner.storageListItems[0].TxOutItem.Equals(tosis[2].TxOutItem), Equals, true)
}

func (s *PipelineSuite) TestPipelineConcurrentLimited(c *C) {
	pipeline, err := newPipeline(2)
	c.Assert(pipeline, NotNil)
	c.Assert(err, IsNil)

	// mocks
	bridge := fakeBridge{nil}
	mockSigner := &mockPipelineSigner{
		storageListItems: append([]TxOutStoreItem{}, tosis...),
	}

	// spawn signings - 4/5 test data should be concurrent, but semaphore limits to 2
	mockSigner.Lock()
	pipeline.SpawnSignings(mockSigner, bridge)
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 2)
	c.Assert(len(mockSigner.storageListItems), Equals, 5)
	c.Assert(len(mockSigner.processed), Equals, 0)

	// release lock and all signings should complete
	mockSigner.Unlock()
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 3)
	c.Assert(len(mockSigner.processed), Equals, 2)

	// the remaining signing should be the last 3 items
	c.Assert(mockSigner.storageListItems[0].TxOutItem.Equals(tosis[2].TxOutItem), Equals, true)
	c.Assert(mockSigner.storageListItems[1].TxOutItem.Equals(tosis[3].TxOutItem), Equals, true)
	c.Assert(mockSigner.storageListItems[2].TxOutItem.Equals(tosis[4].TxOutItem), Equals, true)

	// complete 2 more signings
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.processed), Equals, 4)
	c.Assert(len(mockSigner.storageListItems), Equals, 1)

	// the remaining signing should be the last item
	c.Assert(mockSigner.storageListItems[0].TxOutItem.Equals(tosis[4].TxOutItem), Equals, true)

	// finish
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.processed), Equals, 5)
	c.Assert(len(mockSigner.storageListItems), Equals, 0)
}

func (s *PipelineSuite) TestPipelineConcurrentLimitedStop(c *C) {
	pipeline, err := newPipeline(2)
	c.Assert(pipeline, NotNil)
	c.Assert(err, IsNil)

	// mocks
	bridge := fakeBridge{nil}
	mockSigner := &mockPipelineSigner{
		storageListItems: append([]TxOutStoreItem{}, tosis...),
	}

	// spawn signings - 4/5 test data should be concurrent, but semaphore limits to 2
	mockSigner.Lock()
	pipeline.SpawnSignings(mockSigner, bridge)
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 2)
	c.Assert(len(mockSigner.storageListItems), Equals, 5)
	c.Assert(len(mockSigner.processed), Equals, 0)

	// stop the signer
	mockSigner.stopped = true

	// release lock and all signings should complete
	mockSigner.Unlock()
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 3)
	c.Assert(len(mockSigner.processed), Equals, 2)

	// the remaining signing should be the last 3 items
	c.Assert(mockSigner.storageListItems[0].TxOutItem.Equals(tosis[2].TxOutItem), Equals, true)
	c.Assert(mockSigner.storageListItems[1].TxOutItem.Equals(tosis[3].TxOutItem), Equals, true)
	c.Assert(mockSigner.storageListItems[2].TxOutItem.Equals(tosis[4].TxOutItem), Equals, true)

	// no more signings should be spawned
	pipeline.SpawnSignings(mockSigner, bridge)
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 3)
	c.Assert(len(mockSigner.processed), Equals, 2)
}

func (s *PipelineSuite) TestPipelineRound7Retry(c *C) {
	s.testPipelineRetry(c, true, false)
}

func (s *PipelineSuite) TestPipelineBroadcastRetry(c *C) {
	s.testPipelineRetry(c, false, true)
}

func (s *PipelineSuite) testPipelineRetry(c *C, round7, broadcast bool) {
	pipeline, err := newPipeline(1)
	c.Assert(pipeline, NotNil)
	c.Assert(err, IsNil)

	// use a copy of the test data so we can modify it
	retryTosis := append([]TxOutStoreItem{}, tosis...)

	// multiple items in retry for one vault/chain should be skipped
	if round7 {
		retryTosis[1].Round7Retry = true
		retryTosis[2].Round7Retry = true
	}
	if broadcast {
		retryTosis[1].SignedTx = []byte("broadcast")
		retryTosis[2].SignedTx = []byte("broadcast")
	}

	// mocks
	bridge := fakeBridge{nil}
	mockSigner := &mockPipelineSigner{
		storageListItems: append([]TxOutStoreItem{}, retryTosis...),
	}

	// spawn signings - the first item should process since retry items are skipped
	mockSigner.Lock()
	pipeline.SpawnSignings(mockSigner, bridge)
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 1)
	c.Assert(len(mockSigner.storageListItems), Equals, 5)
	c.Assert(len(mockSigner.processed), Equals, 0)

	// release lock and signing should complete
	mockSigner.Unlock()
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 4)
	c.Assert(len(mockSigner.processed), Equals, 1)

	// the first item should have been the first one processed
	c.Assert(mockSigner.processed[0].TxOutItem.Equals(retryTosis[0].TxOutItem), Equals, true)

	// this time only 1 item for vault/chain should be in retry so it should process first
	if round7 {
		retryTosis[1].Round7Retry = false
	}
	if broadcast {
		retryTosis[1].SignedTx = nil
	}
	mockSigner = &mockPipelineSigner{
		storageListItems: append([]TxOutStoreItem{}, retryTosis...),
	}

	// spawn signings - only the retry item should be started
	mockSigner.Lock()
	pipeline.SpawnSignings(mockSigner, bridge)
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 1)
	c.Assert(len(mockSigner.storageListItems), Equals, 5)
	c.Assert(len(mockSigner.processed), Equals, 0)

	// release lock and signing should complete
	mockSigner.Unlock()
	pipeline.Wait()
	c.Assert(len(pipeline.vaultStatusConcurrency[ttypes.VaultStatus_ActiveVault]), Equals, 0)
	c.Assert(len(mockSigner.storageListItems), Equals, 4)
	c.Assert(len(mockSigner.processed), Equals, 1)

	// the retry item should have been the first one processed
	c.Assert(mockSigner.processed[0].TxOutItem.Equals(retryTosis[2].TxOutItem), Equals, true)
}
