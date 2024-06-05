# Summary

- [Introduction](README.md)

- [Swap Guide](swap-guide/quickstart-guide.md)

  - [Quickstart Guide](swap-guide/quickstart-guide.md)
  - [Fees and Wait Times](swap-guide/fees-and-wait-times.md)
  - [Streaming Swaps](swap-guide/streaming-swaps.md)

- [Lending](lending/quick-start-guide.md)

  - [Quick Start Guide](lending/quick-start-guide.md)

- [Saving Guide](saving-guide/quickstart-guide.md)

  - [Quickstart Guide](saving-guide/quickstart-guide.md)
  - [Fees and Wait Times](saving-guide/fees-and-wait-times.md)

- [Affiliate Guide](affiliate-guide/thorname-guide.md)

  - [THORName Guide](affiliate-guide/thorname-guide.md)

- [Examples](examples/tutorials.md)

  - [TypeScript (Web)](examples/typescript-web/README.md)
    - [Overview](examples/typescript-web/overview.md)
    - [Query Package](examples/typescript-web/query-package.md)
    - [AMM Package](examples/typescript-web/amm-package.md)
    - [Client Packages](examples/typescript-web/client-packages.md)
    - [Packages Breakdown](examples/typescript-web/packages-breakdown.md)
    - [Coding Guide](examples/typescript-web/coding-guide.md)

- [Concepts](concepts/connecting-to-thorchain.md)

  - [Connecting to THORChain](concepts/connecting-to-thorchain.md)
  - [Querying THORChain](concepts/querying-thorchain.md)
  - [Transaction Memos](concepts/memos.md)
  - [Asset Notation](concepts/asset-notation.md)
  - [Memo Length Reduction](concepts/memo-length-reduction.md)
  - [Swapper Clout](./concepts/swapper-clout.md)
  - [Trade Accounts](./concepts/trade-accounts.md)
  - [Network Halts](concepts/network-halts.md)
  - [Fees](concepts/fees.md)
  - [Delays](concepts/delays.md)
  - [Sending Transactions](concepts/sending-transactions.md)
  - [Code Libraries](concepts/code-libraries.md)
  - [Math](concepts/math.md)

- [Aggregators](aggregators/aggregator-overview.md)

  - [Aggregator Overview](aggregators/aggregator-overview.md)
  - [Memos](aggregators/memos.md)
  - [EVM Implementation](aggregators/evm-implementation.md)

- [CLI](cli/overview.md)

  - [Overview](cli/overview.md)
  - [Multisig](cli/multisig.md)
  - [Offline Ledger Support](cli/offline-ledger-support.md)

- [THORNode](release.md)

  - [Release Process](release.md)
  - [EVM Whitelist Procedure](evm_whitelist_procedure.md)
  - [Upgrade Router](upgrade_router.md)
  - [Mimir Abilities](mimir.md)
  - [How to add a new chain](newchain.md)
  - [New Chain Integrations](chains/README.md)
  - [Architecture Decision Records (ADR)](architecture/README.md)
    - [ADR Creation Process](architecture/PROCESS.md)
    - [ADR {ADR-NUMBER}: {TITLE}](architecture/TEMPLATE.md)
    - [ADR 001: ThorChat](architecture/adr-001-thorchat.md)
    - [ADR 002: REMOVE YGG VAULTS](architecture/adr-002-removeyggvaults.md)
    - [ADR 003: FLOORED OUTBOUND FEE](architecture/adr-003-flooredoutboundfee.md)
    - [ADR 004: Keyshare Backups](architecture/adr-004-keyshare-backups.md)
    - [ADR 005: Deprecate Impermanent Loss Protection](architecture/adr-005-deprecate-ilp.md)
    - [ADR 006: Enable POL](architecture/adr-006-enable-pol.md)
    - [ADR 007: Increase Fund Migration and Churn Interval](architecture/adr-007-increase-fund-migration-and-churn-interval.md)
    - [ADR 008: Implement a Dynamic Outbound Fee Multiplier (DOFM)](architecture/adr-008-implement-dynamic-outbound-fee-multiplier.md)
    - [ADR 009: Reserve Income and Fee Overhaul](architecture/adr-009-reserve-income-fee-overhaul.md)
    - [ADR 010: Introduction of Streaming Swaps](architecture/adr-010-streaming-swaps.md)
    - [ADR 011: THORFi Lending Feature](architecture/adr-011-lending.md)
    - [ADR 012: ADR 012: Scale Lending](architecture/adr-012-scale-lending.md)
    - [ADR 013: Synth Backstop](architecture/adr-013-synth-backstop.md)
    - [ADR 014: Reduce Saver Yield Synth Target to Match POL Target](architecture/adr-014-reduce-saver-yield-target.md)

- [Protocol Development](protocol-development/adding-new-chains.md)
  - [Adding New Chains](protocol-development/adding-new-chains.md)
  - [Chain Clients](protocol-development/chain-clients/README.md)
    - [UTXO](protocol-development/chain-clients/utxo.md)
    - [EVM Chains](protocol-development/chain-clients/evm-chains.md)
    - [BFT Chains](protocol-development/chain-clients/bft-chains.md)
  - [ERC20 Tokens](protocol-development/erc20-tokens.md)
