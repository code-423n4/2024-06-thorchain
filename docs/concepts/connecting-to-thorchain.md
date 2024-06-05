# Connecting to THORChain

The Network Information comes from four sources:

1. **Midgard**: Consumer information relating to swaps, pools, and volume. Dashboards will primarily interact with Midgard.
2. **THORNode**: Raw blockchain data provided by the THORChain state machine. THORChain wallets and block explorers will query THORChain-specific information here.
3. **Cosmos RPC**: Used to query for generic CosmosSDK information.
4. **Tendermint RPC**: Used to query for consensus-related information.

```admonish info
The below endpoints are run by specific organisations for public use. There is a cost to running these services. If you want to run your own full node, please see [https://docs.thorchain.org/thornodes/overview.](https://docs.thorchain.org/thornodes/overview)
```

## Midgard

Midgard returns time-series information regarding the THORChain network, such as volume, pool information, users, liquidity providers and more. It also proxies to THORNode to reduce burden on the network. Runs on every node.

**Mainnet:**

- [https://midgard.thorswap.net/v2/doc](https://midgard.thorswap.net/v2/doc)
- [https://midgard.ninerealms.com/v2/doc](https://midgard.ninerealms.com/v2/doc)
- [https://midgard.thorchain.liquify.com/v2/doc](https://midgard.thorchain.liquify.com/v2/doc)

**Stagenet:**

- [https://stagenet-midgard.ninerealms.com/v2/doc](https://stagenet-midgard.ninerealms.com/v2/doc)

## THORNode

THORNode returns application-specific information regarding the THORChain state machine, such as balances, transactions and more. Careful querying this too much - you could overload the public nodes. Consider running your own node. Runs on every node.

**Mainnet (for post-hard-fork blocks 4786560 and later):**

- [https://thornode.thorswap.net/thorchain/doc](https://thornode.thorswap.net/thorchain/doc)
- [https://thornode.ninerealms.com/thorchain/doc](https://thornode.ninerealms.com/thorchain/doc)
- [https://thornode.thorchain.liquify.com/thorchain/doc](https://thornode.thorchain.liquify.com/thorchain/doc)
- **Pre-hard-fork blocks 4786559 and earlier**\
  [https://thornode-v0.ninerealms.com/thorchain/doc](https://thornode.ninerealms.com/thorchain/doc/)

**Stagenet:**

- [https://stagenet-thornode.ninerealms.com/thorchain/doc](https://stagenet-thornode.ninerealms.com/thorchain/doc)

## Cosmos RPC

The Cosmos RPC allows Cosmos base blockchain information to be returned. However, not all endpoints have been enabled.\
\
**Endpoints guide:**
[Cosmos RPC v0.45.1](https://web.archive.org/web/20240106223257/https://v1.cosmos.network/rpc/v0.45.1)

**Example URL** [https://thornode.ninerealms.com/cosmos/bank/v1beta1/balances/thor1dheycdevq39qlkxs2a6wuuzyn4aqxhve4qxtxt](https://thornode.ninerealms.com/cosmos/bank/v1beta1/balances/thor1dheycdevq39qlkxs2a6wuuzyn4aqxhve4qxtxt)

## Tendermint RPC

The Tendermint RPC allows Tendermint consensus information to be returned.

**Any Node Ports:**

- MAINNET Port: `27147`
- STAGENET Port: `26657`

**Endpoints guide.**

[https://docs.tendermint.com/master/rpc/#/](https://docs.tendermint.com/master/rpc/#/)

**Mainnet:**

**`URLs` (for post-hard-fork blocks 4786560 and later)**

- [https://rpc.ninerealms.com](https://rpc.ninerealms.com)
- [https://rpc.thorchain.liquify.com/genesis](https://rpc.thorchain.liquify.com/genesis)
- [https://rpc.thorswap.net/](https://rpc.thorswap.net/)

**Pre-hard-fork blocks 4786559 and earlier.**

- [https://rpc-v0.ninerealms.com](https://rpc-v0.ninerealms.com)

**Stagenet:**

- [https://stagenet-rpc.ninerealms.com](https://stagenet-rpc.ninerealms.com/)

### **P2P**

P2P is the network layer between nodes, useful for network debugging.

MAINNET Port: `27146`

STAGENET Port: `26656`

P2P Guide\
[https://docs.tendermint.com/master/spec/p2p/](https://docs.tendermint.com/master/spec/p2p/)
