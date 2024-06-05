//go:build mocknet
// +build mocknet

package avaxtokens

import (
	_ "embed"
)

//go:embed avax_mocknet_V95.json
var AVAXTokenListRawV95 []byte

//go:embed avax_mocknet_V101.json
var AVAXTokenListRawV101 []byte

//go:embed avax_mocknet_V126.json
var AVAXTokenListRawV126 []byte

//go:embed avax_mocknet_V127.json
var AVAXTokenListRawV127 []byte

//go:embed avax_mocknet_latest.json
var AVAXTokenListRawV131 []byte
