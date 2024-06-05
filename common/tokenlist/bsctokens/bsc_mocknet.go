//go:build mocknet
// +build mocknet

package bsctokens

import _ "embed"

//go:embed bsc_mocknet_V111.json
var BSCTokenListRawV111 []byte

//go:embed bsc_mocknet_V122.json
var BSCTokenListRawV122 []byte

//go:embed bsc_mocknet_latest.json
var BSCTokenListRawV131 []byte
