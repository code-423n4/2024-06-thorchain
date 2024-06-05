package query

import (
	"fmt"
	"strings"
)

// Query define all the queries
type Query struct {
	Key              string
	EndpointTemplate string
}

// Endpoint return the end point string
func (q Query) Endpoint(args ...string) string {
	count := strings.Count(q.EndpointTemplate, "%s")
	a := args[:count]

	in := make([]interface{}, len(a))
	for i := range in {
		in[i] = a[i]
	}

	return fmt.Sprintf(q.EndpointTemplate, in...)
}

// Path return the path
func (q Query) Path(args ...string) string {
	temp := []string{args[0], q.Key}
	args = append(temp, args[1:]...)
	return fmt.Sprintf("custom/%s", strings.Join(args, "/"))
}

// query endpoints supported by the thorchain Querier
var (
	QueryPool                = Query{Key: "pool", EndpointTemplate: "/%s/pool/{%s}"}
	QueryPools               = Query{Key: "pools", EndpointTemplate: "/%s/pools"}
	QueryPoolSlip            = Query{Key: "poolslip", EndpointTemplate: "/%s/slip/{%s}"}
	QueryPoolSlips           = Query{Key: "poolslips", EndpointTemplate: "/%s/slips"}
	QueryDerivedPools        = Query{Key: "derived_pools", EndpointTemplate: "/%s/dpools"}
	QueryDerivedPool         = Query{Key: "derived_pool", EndpointTemplate: "/%s/dpool/{%s}"}
	QueryLiquidityProviders  = Query{Key: "lps", EndpointTemplate: "/%s/pool/{%s}/liquidity_providers"}
	QueryLiquidityProvider   = Query{Key: "lp", EndpointTemplate: "/%s/pool/{%s}/liquidity_provider/{%s}"}
	QuerySavers              = Query{Key: "savers", EndpointTemplate: "/%s/pool/{%s}/savers"}
	QuerySaver               = Query{Key: "saver", EndpointTemplate: "/%s/pool/{%s}/saver/{%s}"}
	QueryBorrowers           = Query{Key: "borrowers", EndpointTemplate: "/%s/pool/{%s}/borrowers"}
	QueryBorrower            = Query{Key: "borrower", EndpointTemplate: "/%s/pool/{%s}/borrower/{%s}"}
	QueryTradeUnit           = Query{Key: "tradeunit", EndpointTemplate: "/%s/trade/unit/{%s}"}
	QueryTradeUnits          = Query{Key: "tradeunits", EndpointTemplate: "/%s/trade/units"}
	QueryTradeAccount        = Query{Key: "tradeaccount", EndpointTemplate: "/%s/trade/account/{%s}"}
	QueryTradeAccounts       = Query{Key: "tradeaccounts", EndpointTemplate: "/%s/trade/accounts/{%s}"}
	QueryTx                  = Query{Key: "tx", EndpointTemplate: "/%s/tx/{%s}"}
	QueryTxVoterOld          = Query{Key: "txvoterold", EndpointTemplate: "/%s/tx/{%s}/signers"}
	QueryTxVoter             = Query{Key: "txvoter", EndpointTemplate: "/%s/tx/details/{%s}"}
	QueryTxStages            = Query{Key: "txstages", EndpointTemplate: "/%s/tx/stages/{%s}"}
	QueryTxStatus            = Query{Key: "txstatus", EndpointTemplate: "/%s/tx/status/{%s}"}
	QueryKeysignArray        = Query{Key: "keysign", EndpointTemplate: "/%s/keysign/{%s}"}
	QueryKeysignArrayPubkey  = Query{Key: "keysignpubkey", EndpointTemplate: "/%s/keysign/{%s}/{%s}"}
	QueryKeygensPubkey       = Query{Key: "keygenspubkey", EndpointTemplate: "/%s/keygen/{%s}/{%s}"}
	QueryQueue               = Query{Key: "outqueue", EndpointTemplate: "/%s/queue"}
	QueryHeights             = Query{Key: "heights", EndpointTemplate: "/%s/lastblock"}
	QueryChainHeights        = Query{Key: "chainheights", EndpointTemplate: "/%s/lastblock/{%s}"}
	QueryNodes               = Query{Key: "nodes", EndpointTemplate: "/%s/nodes"}
	QueryNode                = Query{Key: "node", EndpointTemplate: "/%s/node/{%s}"}
	QueryInboundAddresses    = Query{Key: "inboundaddresses", EndpointTemplate: "/%s/inbound_addresses"}
	QueryNetwork             = Query{Key: "network", EndpointTemplate: "/%s/network"}
	QueryPOL                 = Query{Key: "pol", EndpointTemplate: "/%s/pol"}
	QueryBalanceModule       = Query{Key: "balancemodule", EndpointTemplate: "/%s/balance/module/{%s}"}
	QueryVaultsAsgard        = Query{Key: "vaultsasgard", EndpointTemplate: "/%s/vaults/asgard"}
	QueryVaultsYggdrasil     = Query{Key: "vaultsyggdrasil", EndpointTemplate: "/%s/vaults/yggdrasil"} // TODO remove on hard fork
	QueryVault               = Query{Key: "vault", EndpointTemplate: "/%s/vault/{%s}"}
	QueryVaultPubkeys        = Query{Key: "vaultpubkeys", EndpointTemplate: "/%s/vaults/pubkeys"}
	QueryConstantValues      = Query{Key: "constants", EndpointTemplate: "/%s/constants"}
	QueryVersion             = Query{Key: "version", EndpointTemplate: "/%s/version"}
	QueryMimirValues         = Query{Key: "mimirs", EndpointTemplate: "/%s/mimir"}
	QueryMimirV2Values       = Query{Key: "mimirsV2", EndpointTemplate: "/%s/mimirV2"}
	QueryMimirV2NodesAll     = Query{Key: "mimirv2nodesall", EndpointTemplate: "/%s/mimirV2/nodes_all"}
	QueryMimirV2IDs          = Query{Key: "mimirsV2IDs", EndpointTemplate: "/%s/mimirV2/ids"}
	QueryMimirWithKey        = Query{Key: "mimirwithkey", EndpointTemplate: "/%s/mimir/key/{%s}"}
	QueryMimirAdminValues    = Query{Key: "adminmimirs", EndpointTemplate: "/%s/mimir/admin"}
	QueryMimirNodesValues    = Query{Key: "nodesmimirs", EndpointTemplate: "/%s/mimir/nodes"}
	QueryMimirNodesAllValues = Query{Key: "nodesmimirsall", EndpointTemplate: "/%s/mimir/nodes_all"}
	QueryMimirNodeValues     = Query{Key: "nodemimirs", EndpointTemplate: "/%s/mimir/node/{%s}"}
	QueryOutboundFees        = Query{Key: "outboundfees", EndpointTemplate: "/%s/outbound_fees"}
	QueryOutboundFee         = Query{Key: "outboundfee", EndpointTemplate: "/%s/outbound_fee/{%s}"}
	QueryBan                 = Query{Key: "ban", EndpointTemplate: "/%s/ban/{%s}"}
	QueryRagnarok            = Query{Key: "ragnarok", EndpointTemplate: "/%s/ragnarok"}
	QueryPendingOutbound     = Query{Key: "pendingoutbound", EndpointTemplate: "/%s/queue/outbound"}
	QueryScheduledOutbound   = Query{Key: "scheduledoutbound", EndpointTemplate: "/%s/queue/scheduled"}
	QuerySwapperClout        = Query{Key: "swapperclout", EndpointTemplate: "/%s/clout/swap/{%s}"}
	QueryStreamingSwap       = Query{Key: "streamingswap", EndpointTemplate: "/%s/swap/streaming/{%s}"}
	QueryStreamingSwaps      = Query{Key: "streamingswaps", EndpointTemplate: "/%s/swaps/streaming"}
	QuerySwapQueue           = Query{Key: "swapqueue", EndpointTemplate: "/%s/queue/swap"}
	QueryTssKeygenMetrics    = Query{Key: "tss_keygen_metric", EndpointTemplate: "/%s/metric/keygen/{%s}"}
	QueryTssMetrics          = Query{Key: "tss_metric", EndpointTemplate: "/%s/metrics"}
	QueryTHORName            = Query{Key: "thorname", EndpointTemplate: "/%s/thorname/{%s}"}
	QueryQuoteSwap           = Query{Key: "quoteswap", EndpointTemplate: "/%s/quote/swap"}
	QueryQuoteSaverDeposit   = Query{Key: "quotesaverdeposit", EndpointTemplate: "/%s/quote/saver/deposit"}
	QueryQuoteSaverWithdraw  = Query{Key: "quotesaverwithdraw", EndpointTemplate: "/%s/quote/saver/withdraw"}
	QueryQuoteLoanOpen       = Query{Key: "quoteloanopen", EndpointTemplate: "/%s/quote/loan/open"}
	QueryQuoteLoanClose      = Query{Key: "quoteloanclose", EndpointTemplate: "/%s/quote/loan/close"}
	QueryInvariants          = Query{Key: "invariants", EndpointTemplate: "/%s/invariants"}
	QueryInvariant           = Query{Key: "invariant", EndpointTemplate: "/%s/invariant/{%s}"}
	QueryBlock               = Query{Key: "block", EndpointTemplate: "/%s/block"}

	// queries only available on regtest builds
	QueryExport = Query{Key: "export", EndpointTemplate: "/%s/export"}
)

// Queries all queries
var Queries = []Query{
	QueryPool,
	QueryPools,
	QueryPoolSlip,
	QueryPoolSlips,
	QueryDerivedPool,
	QueryDerivedPools,
	QueryLiquidityProviders,
	QueryLiquidityProvider,
	QuerySavers,
	QuerySaver,
	QueryBorrowers,
	QueryBorrower,
	QueryTradeUnit,
	QueryTradeUnits,
	QueryTradeAccount,
	QueryTradeAccounts,
	QueryTxStages,
	QueryTxStatus,
	QueryTxVoter,
	QueryTxVoterOld,
	QueryTx,
	QueryKeysignArray,
	QueryKeysignArrayPubkey,
	QueryQueue,
	QueryHeights,
	QueryChainHeights,
	QueryNode,
	QueryNodes,
	QueryInboundAddresses,
	QueryNetwork,
	QueryPOL,
	QueryBalanceModule,
	QueryVaultsAsgard,
	QueryVaultsYggdrasil, // TODO remove on hard fork
	QueryVaultPubkeys,
	QueryVault,
	QueryKeygensPubkey,
	QueryConstantValues,
	QueryVersion,
	QueryMimirValues,
	QueryMimirV2IDs,
	QueryMimirV2NodesAll,
	QueryMimirV2Values,
	QueryMimirWithKey,
	QueryMimirAdminValues,
	QueryMimirNodesAllValues,
	QueryMimirNodesValues,
	QueryMimirNodeValues,
	QueryOutboundFees,
	QueryOutboundFee,
	QueryBan,
	QueryRagnarok,
	QueryPendingOutbound,
	QueryScheduledOutbound,
	QuerySwapperClout,
	QueryStreamingSwap,
	QueryStreamingSwaps,
	QuerySwapQueue,
	QueryTssMetrics,
	QueryTssKeygenMetrics,
	QueryTHORName,
	QueryQuoteSwap,
	QueryQuoteSaverDeposit,
	QueryQuoteSaverWithdraw,
	QueryQuoteLoanOpen,
	QueryQuoteLoanClose,
	QueryInvariants,
	QueryInvariant,
	QueryBlock,
	QueryExport,
}
