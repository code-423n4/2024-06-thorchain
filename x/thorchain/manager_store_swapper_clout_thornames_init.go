package thorchain

import "gitlab.com/thorchain/thornode/common/cosmos"

type initCloutTHORNames struct {
	address string
	amount  cosmos.Uint
}

func getInitCloutTHORNames() []initCloutTHORNames {
	// init swapper clout for thornames
	// thornames were not included in last migration (v125)
	// these are the remaining clout scores that did not get counted
	return []initCloutTHORNames{
		// "zg.ETH": 251495265994,
		{address: "0x00000000d7c185343e6504e428b8f8b5ad6c91b8", amount: cosmos.NewUint(251495265994)},
		// "zg.THOR": 128841556399,
		{address: "thor1wx5av89rghsmgh2vh40aknx7csvs7xj2cr474n", amount: cosmos.NewUint(128841556399)},
		// "zg.BNB": 81622839365,
		{address: "bnb1cl3lmqk62k9ted6t7ujvf9phey7w3pl9ha4y9y", amount: cosmos.NewUint(81622839365)},
		// "zg.LTC": 61990822020,
		{address: "ltc1qdxrce4ms9hvaxvh9p3rdmxfdtvvzt9y3xuawad", amount: cosmos.NewUint(61990822020)},
		// "zg.BSC": 54828663639,
		{address: "0x05B7F35B1b84E15bd3e5fFe91023918D9d5cccDE", amount: cosmos.NewUint(54828663639)},
		// "zg.DOGE": 39721688604,
		{address: "DEm5sbETWwgdCTHGvoWozNNZcLarsiMk3J", amount: cosmos.NewUint(39721688604)},
		// "zg.BCH": 39664724745,
		{address: "qzw88zasl5a5y9f6tkxgw8gjhfrayxukjgxqplqax7", amount: cosmos.NewUint(39664724745)},
		// "zg.BTC": 36503015631,
		{address: "bc1qdxrce4ms9hvaxvh9p3rdmxfdtvvzt9y3zq829a", amount: cosmos.NewUint(36503015631)},
		// "zg.AVAX": 19370133253,
		{address: "0xf841a830cd94f6f00be674c81f57d5fcbbee2857", amount: cosmos.NewUint(19370133253)},
		// "runifier": 47863697284,
		{address: "thor1m5gnkrh6x9rpkeusesc7t3qr0t4lkhu6fjvw2z", amount: cosmos.NewUint(47863697284)},
	}
}
