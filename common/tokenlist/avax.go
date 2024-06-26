package tokenlist

import (
	"encoding/json"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common/tokenlist/avaxtokens"
)

var (
	avaxTokenListV95  EVMTokenList
	avaxTokenListV101 EVMTokenList
	avaxTokenListV126 EVMTokenList
	avaxTokenListV127 EVMTokenList
	avaxTokenListV131 EVMTokenList
)

func init() {
	if err := json.Unmarshal(avaxtokens.AVAXTokenListRawV95, &avaxTokenListV95); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(avaxtokens.AVAXTokenListRawV101, &avaxTokenListV101); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(avaxtokens.AVAXTokenListRawV126, &avaxTokenListV126); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(avaxtokens.AVAXTokenListRawV127, &avaxTokenListV127); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(avaxtokens.AVAXTokenListRawV131, &avaxTokenListV131); err != nil {
		panic(err)
	}
}

func GetAVAXTokenList(version semver.Version) EVMTokenList {
	switch {
	case version.GTE(semver.MustParse("1.131.0")):
		return avaxTokenListV131
	case version.GTE(semver.MustParse("1.127.0")):
		return avaxTokenListV127
	case version.GTE(semver.MustParse("1.126.0")):
		return avaxTokenListV126
	case version.GTE(semver.MustParse("1.101.0")):
		return avaxTokenListV101
	case version.GTE(semver.MustParse("1.95.0")):
		return avaxTokenListV95
	default:
		return avaxTokenListV95
	}
}
