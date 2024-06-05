//go:build !mocknet && !stagenet
// +build !mocknet,!stagenet

package thorchain

import _ "embed"

//go:embed preregister_thornames.json
var preregisterTHORNames []byte
