# Adding New Chains

Chain Developers should be extremely familiar with how THORChain works, and how their own chain works.

There is now a specific process for the addition of new chains, see: [https://gitlab.com/thorchain/thornode/-/blob/develop/docs/chains/README.md](https://gitlab.com/thorchain/thornode/-/blob/develop/docs/chains/README.md)

## Process

1. Read [https://gitlab.com/thorchain/thornode/-/blob/develop/docs/newchain.md](https://gitlab.com/thorchain/thornode/-/blob/develop/docs/newchain.md)
2. Bifrost: Start by forking one of the existing Bifrosts (UTXO, EVM or BFT).
3. Daemon: Add the chain daemon to THORChain/Node-Launcher [https://gitlab.com/thorchain/thornode/-/tree/develop/bifrost/pkg/chainclients](https://gitlab.com/thorchain/thornode/-/tree/develop/bifrost/pkg/chainclients)
4. Simulation Tests: Build out the simulation tests for the chain. This ensures the connection is robustly tested.
5. [XChainJS](https://github.com/xchainjs/xchainjs-lib): Add a new chain package to xchainjs so the entire ecosystem of wallets can easily support.

Once this is complete, the chain can be added to Stagenet. After some time of demonstrating Stability on Stagenet, the THORChain Node Operator community is polled and if supported, it can be merged to Mainnet.

Once on mainnet, the chain is typically given a period of 12 months to demonstrate uptake and usage. If the chain cannot maintain sufficient demand, it may be removed from the network and all liquidity refunded to LPs.
