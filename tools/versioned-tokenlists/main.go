package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/tokenlist"
)

// -------------------------------------------------------------------------------------
// Flags
// -------------------------------------------------------------------------------------

var flagVersion *int

func init() {
	flagVersion = flag.Int("version", 0, "current version allowing changes")
}

// -------------------------------------------------------------------------------------
// Check
// -------------------------------------------------------------------------------------

func check(chain common.Chain) {
	// write all token lists to stdout
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	version, err := semver.Parse("1.93.0") // TODO: bump on hard fork
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("Check:", chain, version)

		// get token list
		err = enc.Encode(tokenlist.GetEVMTokenList(chain, version))
		if err != nil {
			panic(err)
		}

		// iterate versions up to current
		version.Minor++
		if version.Minor >= uint64(*flagVersion) {
			break
		}
	}
}

// -------------------------------------------------------------------------------------
// Main
// -------------------------------------------------------------------------------------

func main() {
	flag.Parse()
	if *flagVersion == 0 {
		panic("version is required")
	}

	for _, chain := range common.AllChains {
		if chain.IsEVM() {
			check(chain)
		}
	}
}
