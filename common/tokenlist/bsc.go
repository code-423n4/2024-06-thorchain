package tokenlist

import (
	"encoding/json"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common/tokenlist/bsctokens"
)

var (
	bscTokenListV111 EVMTokenList
	bscTokenListV122 EVMTokenList
	bscTokenListV131 EVMTokenList
)

func init() {
	if err := json.Unmarshal(bsctokens.BSCTokenListRawV111, &bscTokenListV111); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bsctokens.BSCTokenListRawV122, &bscTokenListV122); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bsctokens.BSCTokenListRawV131, &bscTokenListV131); err != nil {
		panic(err)
	}
}

func GetBSCTokenList(version semver.Version) EVMTokenList {
	switch {
	case version.GTE(semver.MustParse("1.131.0")):
		return bscTokenListV131
	case version.GTE(semver.MustParse("1.122.0")):
		return bscTokenListV122
	case version.GTE(semver.MustParse("1.111.0")):
		return bscTokenListV111
	default:
		return bscTokenListV111
	}
}
