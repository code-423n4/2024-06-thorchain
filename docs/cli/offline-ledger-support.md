# Offline Ledger Support

When used in conjunction with a locally-running fullnode, the THORNode CLI + Ledger provides the ultimate, privacy-focused "offline, no-tracking" experience. Interact directly with the THORChain network to bond validators, swap RUNE (or synthetic assets) and administrate LPs from a cold-wallet.

## Accounts

Ledger accounts can be added by appending `--ledger` to the command. The default index is 0.

```text
thornode keys add ledger1 --ledger --index=1
```

## Usage

Signing transactions requires confirmation through the Ledger. Everything else works the same.
