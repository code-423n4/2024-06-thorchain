package tokenlist

import (
	"fmt"
	"strings"
	"time"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
)

// ERC20Token is a struct to represent the token
type ERC20Token struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
}

// Asset returns the common.Asset representation of the token.
func (t ERC20Token) Asset(chain common.Chain) common.Asset {
	return common.Asset{
		Chain:  chain,
		Ticker: common.Ticker(t.Symbol),
		Symbol: common.Symbol(strings.ToUpper(fmt.Sprintf("%s-%s", t.Symbol, t.Address))),
	}
}

type EVMTokenList struct {
	Name      string       `json:"name"`
	LogoURI   string       `json:"logoURI"`
	Tokens    []ERC20Token `json:"tokens"`
	Keywords  []string     `json:"keywords"`
	Timestamp time.Time    `json:"timestamp"`
}

// GetEVMTokenList returns all available tokens for external asset matching for a
// particular EVM chain and version.
//
// NOTE: These tokens are NOT necessarily the same tokens that are whitelisted for each
// chain - whitelisting happens in each chain's bifrost chain client.
func GetEVMTokenList(chain common.Chain, version semver.Version) EVMTokenList {
	switch chain {
	case common.ETHChain:
		return GetETHTokenList(version)
	case common.AVAXChain:
		return GetAVAXTokenList(version)
	case common.BSCChain:
		return GetBSCTokenList(version)
	default:
		return EVMTokenList{}
	}
}
