# Packages Breakdown

## How xChainjs is constructed

### xchain-\[chain] clients

Each blockchian that is integrated into THORChain has a corresponding xchain client with a suite of functionality to work with that chain. They all extend the `xchain-client` class.

### xchain-thorchain-amm

Thorchain automatic market maker that uses Thornode & Midgard Api's AMM functions like swapping, adding and removing liquidity. It wraps xchain clients and creates a new wallet class and balance collection.

### xchain-thorchain-query

Uses midgard and thornode Api's to query Thorchain for information. This module should be used as the starting place get any THORChain information that resides in THORNode or Midgard as it does the heaving lifting and configuration.

Default endpoints are provided with redundancy, custom THORNode or Midgard endpoints can be provided in the constructor.

### **xchain-midgard**

This package is built from OpenAPI-generator. It is used by the thorchain-query.

Thorchain-query contains midgard class that uses xchain-midgard and the following end points:

- /v2/thorchain/mimir
- /v2/thorchain/inbound_addresses
- /v2/thorchain/constants
- /v2/thorchain/queue

For simplicity, is recommended to use the midgard class within thorchain-query instead of using the midgard package directly.

#### Midgard Configuration in thorchain-query

Default endpoints `defaultMidgardConfig` are provided with redundancy within the Midgard class.

```typescript
// How thorchain-query constructs midgard
const defaultMidgardConfig: Record<Network, MidgardConfig> = {
  mainnet: {
    apiRetries: 3,
    midgardBaseUrls: [
      'https://midgard.ninerealms.com',
      'https://midgard.thorchain.info',
      'https://midgard.thorswap.net',
    ],
  },
  ...
  export class Midgard {
  private config: MidgardConfig
  readonly network: Network
  private midgardApis: MidgardApi[]

  constructor(network: Network = Network.Mainnet, config?: MidgardConfig) {
    this.network = network
    this.config = config ?? defaultMidgardConfig[this.network]
    axiosRetry(axios, { retries: this.config.apiRetries, retryDelay: axiosRetry.exponentialDelay })
    this.midgardApis = this.config.midgardBaseUrls.map((url) => new MidgardApi(new Configuration({ basePath: url })))
  }
```

Custom Midgard endpoints can be provided in the constructor using the `MidgardConfig` type.

```typescript
// adding custom endpoints
  const network = Network.Mainnet
  const customMidgardConfig: MidgardConfig = {
    apiRetries: 3,
    midgardBaseUrls: [
      'https://midgard.customURL.com',
    ],
  }
  const midgard = new Midgard(network, customMidgardConfig)
}
```

See [ListPools](query-package.md#list-pools) for a working example.

### xchain-thornode

This package is built from OpenAPI-generator and is also used by the thorchain-query. The design is similar to the midgard. Thornode should only be used when time-sensitive data is required else midgard should be used.

```typescript
// How thorchain-query constructs thornode
const defaultThornodeConfig: Record<Network, ThornodeConfig> = {
  mainnet: {
    apiRetries: 3,
    thornodeBaseUrls: [
      `https://thornode.ninerealms.com`,
      `https://thornode.thorswap.net`,
      `https://thornode.thorchain.info`,
    ],
  },
  ...
  export class Thornode {
  private config: ThornodeConfig
  private network: Network
 ...
  constructor(network: Network = Network.Mainnet, config?: ThornodeConfig) {
    this.network = network
    this.config = config ?? defaultThornodeConfig[this.network]
    axiosRetry(axios, { retries: this.config.apiRetries, retryDelay: axiosRetry.exponentialDelay })
    this.transactionsApi = this.config.thornodeBaseUrls.map(
      (url) => new TransactionsApi(new Configuration({ basePath: url })),
    )
    this.queueApi = this.config.thornodeBaseUrls.map((url) => new QueueApi(new Configuration({ basePath: url })))
    this.networkApi = this.config.thornodeBaseUrls.map((url) => new NetworkApi(new Configuration({ basePath: url })))
    this.poolsApi = this.config.thornodeBaseUrls.map((url) => new PoolsApi(new Configuration({ basePath: url })))
    this.liquidityProvidersApi = this.config.thornodeBaseUrls.map(
      (url) => new LiquidityProvidersApi(new Configuration({ basePath: url })),
    )
  }
```

### Thornode Configuration in thorchain-query

As with the midgard package, thornode can also be given custom end points via the `ThornodeConfig` type.

## **xchain-util**

A helper packager used by all the other packages. It has the following modules:

- `asset` - Utilities for handling assets
- `async` - Utitilies for `async` handling
- `bn` - Utitilies for using `bignumber.js`
- `chain` - Utilities for multi-chain
- `string` - Utilities for strings
