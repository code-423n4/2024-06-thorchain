package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	prefix "gitlab.com/thorchain/thornode/cmd"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/test/simulation/actors"
	"gitlab.com/thorchain/thornode/test/simulation/actors/suites"
	pkgcosmos "gitlab.com/thorchain/thornode/test/simulation/pkg/cosmos"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/dag"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/evm"
	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/utxo"
	"gitlab.com/thorchain/thornode/test/simulation/watchers"
)

////////////////////////////////////////////////////////////////////////////////////////
// Config
////////////////////////////////////////////////////////////////////////////////////////

const (
	DefaultParallelism = "8"
)

// trunk-ignore(golangci-lint/unused)
var liteClientConstructors = map[common.Chain]LiteChainClientConstructor{
	common.BTCChain:  utxo.NewConstructor(chainRPCs[common.BTCChain]),
	common.LTCChain:  utxo.NewConstructor(chainRPCs[common.LTCChain]),
	common.BCHChain:  utxo.NewConstructor(chainRPCs[common.BCHChain]),
	common.DOGEChain: utxo.NewConstructor(chainRPCs[common.DOGEChain]),
	common.ETHChain:  evm.NewConstructor(chainRPCs[common.ETHChain]),
	common.BSCChain:  evm.NewConstructor(chainRPCs[common.BSCChain]),
	common.AVAXChain: evm.NewConstructor(chainRPCs[common.AVAXChain]),
	common.GAIAChain: pkgcosmos.NewConstructor(chainRPCs[common.GAIAChain]),
}

////////////////////////////////////////////////////////////////////////////////////////
// Main
////////////////////////////////////////////////////////////////////////////////////////

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()

	// init prefixes
	ccfg := cosmos.GetConfig()
	ccfg.SetBech32PrefixForAccount(prefix.Bech32PrefixAccAddr, prefix.Bech32PrefixAccPub)
	ccfg.SetBech32PrefixForValidator(prefix.Bech32PrefixValAddr, prefix.Bech32PrefixValPub)
	ccfg.SetBech32PrefixForConsensusNode(prefix.Bech32PrefixConsAddr, prefix.Bech32PrefixConsPub)
	ccfg.SetCoinType(prefix.THORChainCoinType)
	ccfg.SetPurpose(prefix.THORChainCoinPurpose)
	ccfg.Seal()

	// wait until bifrost is ready
	for {
		res, err := http.Get("http://localhost:6040/p2pid")
		if err == nil && res.StatusCode == 200 {
			break
		}
		log.Info().Msg("waiting for bifrost to be ready")
		time.Sleep(time.Second)
	}

	// combine all actor dags for the complete test run
	root := NewActor("Root")
	root.Append(suites.Bootstrap())
	root.Append(actors.NewArbActor())

	// skip swaps and ragnarok if this is bootstrap only mode
	if os.Getenv("BOOTSTRAP_ONLY") != "true" {
		root.Append(suites.Swaps())
		root.Append(suites.Ragnarok())
	}

	// gather config from the environment
	parallelism := os.Getenv("PARALLELISM")
	if parallelism == "" {
		parallelism = DefaultParallelism
	}
	parallelismInt, err := strconv.Atoi(parallelism)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse PARALLELISM")
	}

	cfg := InitConfig(parallelismInt)

	// start watchers
	for _, w := range []*Watcher{watchers.NewInvariants()} {
		log.Info().Str("watcher", w.Name).Msg("starting watcher")
		go func(w *Watcher) {
			err := w.Execute(cfg, log.Output(os.Stderr))
			if err != nil {
				log.Fatal().Err(err).Msg("watcher failed")
			}
		}(w)
	}

	// run the simulation
	dag.Execute(cfg, root, parallelismInt)
}
