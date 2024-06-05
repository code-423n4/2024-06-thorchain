package tokenlist

import (
	"encoding/json"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common/tokenlist/ethtokens"
)

var (
	ethTokenListV93  EVMTokenList
	ethTokenListV95  EVMTokenList
	ethTokenListV97  EVMTokenList
	ethTokenListV101 EVMTokenList
	ethTokenListV108 EVMTokenList
	ethTokenListV114 EVMTokenList
	ethTokenListV126 EVMTokenList
	ethTokenListV128 EVMTokenList
	ethTokenListV133 EVMTokenList
)

func init() {
	if err := json.Unmarshal(ethtokens.ETHTokenListRawV93, &ethTokenListV93); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(ethtokens.ETHTokenListRawV95, &ethTokenListV95); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(ethtokens.ETHTokenListRawV97, &ethTokenListV97); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(ethtokens.ETHTokenListRawV101, &ethTokenListV101); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(ethtokens.ETHTokenListRawV108, &ethTokenListV108); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(ethtokens.ETHTokenListRawV114, &ethTokenListV114); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(ethtokens.ETHTokenListRawV126, &ethTokenListV126); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(ethtokens.ETHTokenListRawV128, &ethTokenListV128); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(ethtokens.ETHTokenListRawV133, &ethTokenListV133); err != nil {
		panic(err)
	}
}

func GetETHTokenList(version semver.Version) EVMTokenList {
	switch {
	case version.GTE(semver.MustParse("1.133.0")):
		return ethTokenListV133
	case version.GTE(semver.MustParse("1.128.0")):
		return ethTokenListV128
	case version.GTE(semver.MustParse("1.126.0")):
		return ethTokenListV126
	case version.GTE(semver.MustParse("1.114.0")):
		return ethTokenListV114
	case version.GTE(semver.MustParse("1.108.0")):
		return ethTokenListV108
	case version.GTE(semver.MustParse("1.101.0")):
		return ethTokenListV101
	case version.GTE(semver.MustParse("1.97.0")):
		return ethTokenListV97
	case version.GTE(semver.MustParse("1.95.0")):
		return ethTokenListV95
	case version.GTE(semver.MustParse("1.93.0")):
		return ethTokenListV93
	default:
		return ethTokenListV93
	}
}
