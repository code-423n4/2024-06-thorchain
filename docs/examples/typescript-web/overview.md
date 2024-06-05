# Overview

XChainJS is an open-source library with a common interface for multiple blockchains, built for simple and fast integration for wallets and Dexs and more. xChainjs is designed to abstract THORChain and specific blockchain complexity and to provide an easy-to-use API for developers.

The packages implement the complexity detailed in the other sections of this site.

xChain has several key modules allowing powerful functionality.

## **Thorchain-query**

Allows easy information retrieval and estimates from THORChain.

[Query Package](query-package.md)

## **Thorchain-amm**

Conducts actions such as swap, add and remove. It wraps xchain clients and creates a new wallet class for and balance collection.

[AMM Package](amm-package.md)

## **Chain clients**

For every blockchain connected to THORChain with a common interface.

Current clients implemented are**:**

- xchain-avax
- xchain-binance
- xchain-bitcoin
- xchain-bitcoincash
- xchain-cosmos
- xchain-doge
- xchain-ethereum
- xchain-litecoin
- xchain-mayachain
- xchain-thorchain

[Client Packages](client-packages.md)

**APIs** for getting data from THORChain.

- Midgard
- Thornode

[Packages Breakdown](packages-breakdown.md)

See the package breakdown for more information.

### Install Procedures

Ensure you have the following

- npm --version v8.5.5 or above
- node --version v16.15.0
- yarn --version v1.22.18 or above

Create a new project by creating a new folder, then type `npx tsc --init`.

#### Finding required dependencies

The replit code examples have all the required packages within the project.json file, just copy the project dependencies into your own project.json.

Example for the [query-package](query-package.md), [estimateSwap](query-package.md#estimate-swap) packages

1. Go to the replit code example then press show files. Select the project.json file.
2. Locate and then copy the `dependencies` section into your project.json file.
3. From the command line, type `yarn`. This will download and install the required packages.

The code is available on [GitHub](https://github.com/xchainjs/xchainjs-lib/) and managed by several key developers. Reach out at Telegram group: [https://t.me/xchainjs](https://t.me/xchainjs) for more information.
