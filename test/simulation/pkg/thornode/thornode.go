package thornode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gitlab.com/thorchain/thornode/common"
	sdk "gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/config"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"
)

////////////////////////////////////////////////////////////////////////////////////////
// Init
////////////////////////////////////////////////////////////////////////////////////////

var thornodeURL string

func init() {
	config.Init()
	thornodeURL = config.GetBifrost().Thorchain.ChainHost
	if !strings.HasPrefix(thornodeURL, "http") {
		thornodeURL = "http://" + thornodeURL
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Exported
////////////////////////////////////////////////////////////////////////////////////////

func GetBalances(addr common.Address) (common.Coins, error) {
	url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", thornodeURL, addr)
	var balances struct {
		Balances []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balances"`
	}
	err := get(url, &balances)
	if err != nil {
		return nil, err
	}

	// convert to common.Coins
	coins := make(common.Coins, 0, len(balances.Balances))
	for _, balance := range balances.Balances {
		amount, err := strconv.ParseUint(balance.Amount, 10, 64)
		if err != nil {
			return nil, err
		}
		asset, err := common.NewAsset(strings.ToUpper(balance.Denom))
		if err != nil {
			return nil, err
		}
		coins = append(coins, common.NewCoin(asset, sdk.NewUint(amount)))
	}

	return coins, nil
}

func GetInboundAddress(chain common.Chain) (address common.Address, router *common.Address, err error) {
	url := fmt.Sprintf("%s/thorchain/inbound_addresses", thornodeURL)
	var inboundAddresses []openapi.InboundAddress
	err = get(url, &inboundAddresses)
	if err != nil {
		return "", nil, err
	}

	// find address for chain
	for _, inboundAddress := range inboundAddresses {
		if *inboundAddress.Chain == string(chain) {
			var router *common.Address
			if inboundAddress.Router != nil {
				router = new(common.Address)
				*router = common.Address(*inboundAddress.Router)
			}
			return common.Address(*inboundAddress.Address), router, nil
		}
	}

	return "", nil, fmt.Errorf("no inbound address found for chain %s", chain)
}

func GetRouterAddress(chain common.Chain) (common.Address, error) {
	url := fmt.Sprintf("%s/thorchain/inbound_addresses", thornodeURL)
	var inboundAddresses []openapi.InboundAddress
	err := get(url, &inboundAddresses)
	if err != nil {
		return "", err
	}

	// find address for chain
	for _, inboundAddress := range inboundAddresses {
		if *inboundAddress.Chain == string(chain) {
			return common.Address(*inboundAddress.Router), nil
		}
	}

	return "", fmt.Errorf("no inbound address found for chain %s", chain)
}

func GetLiquidityProviders(asset common.Asset) ([]openapi.LiquidityProvider, error) {
	url := fmt.Sprintf("%s/thorchain/pool/%s/liquidity_providers", thornodeURL, asset.String())
	var liquidityProviders []openapi.LiquidityProvider
	err := get(url, &liquidityProviders)
	return liquidityProviders, err
}

func GetPools() ([]openapi.Pool, error) {
	url := fmt.Sprintf("%s/thorchain/pools", thornodeURL)
	var pools []openapi.Pool
	err := get(url, &pools)
	return pools, err
}

func GetPool(asset common.Asset) (openapi.Pool, error) {
	url := fmt.Sprintf("%s/thorchain/pool/%s", thornodeURL, asset.String())
	var pool openapi.Pool
	err := get(url, &pool)
	return pool, err
}

func GetSwapQuote(from, to common.Asset, amount sdk.Uint) (openapi.QuoteSwapResponse, error) {
	baseURL := fmt.Sprintf("%s/thorchain/quote/swap", thornodeURL)
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return openapi.QuoteSwapResponse{}, err
	}
	params := url.Values{}
	params.Add("from_asset", from.String())
	params.Add("to_asset", to.String())
	params.Add("amount", amount.String())
	parsedURL.RawQuery = params.Encode()
	url := parsedURL.String()

	var quote openapi.QuoteSwapResponse
	err = get(url, &quote)
	return quote, err
}

func GetTxStages(txid string) (openapi.TxStagesResponse, error) {
	url := fmt.Sprintf("%s/thorchain/tx/stages/%s", thornodeURL, txid)
	var stages openapi.TxStagesResponse
	err := get(url, &stages)
	return stages, err
}

////////////////////////////////////////////////////////////////////////////////////////
// Internal
////////////////////////////////////////////////////////////////////////////////////////

func get(url string, target interface{}) error {
	// trunk-ignore(golangci-lint/gosec): variable url ok
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// extract error if the request failed
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	errResp := ErrorResponse{}
	err = json.Unmarshal(buf, &errResp)
	if err == nil && errResp.Error != "" {
		return fmt.Errorf(errResp.Error)
	}

	// decode response
	return json.Unmarshal(buf, target)
}
