# Overview

## Install (Mac)

### Prerequisites

1. `xcode-select xcode-select --install`
2. Homebrew: [https://brew.sh](https://brew.sh)

### GoLang

Install go v1.18.1: [https://go.dev/doc/install](https://go.dev/doc/install)

```shell
# Set PATH
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
```

### Protobuf

```shell
# Install Protobuf
brew install protobuf
brew install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latestellina
```

### GNU Utils

```shell
# Install GNU find
brew install findutils

# Set PATH
export PATH=$(brew --prefix)/opt/findutils/libexec/gnubin:$PATH
```

### Docker

```shell
# Install docker
brew install homebrew/cask/docker
```

### THORNode

```shell
# Clone repo and install dependencies
git clone https://gitlab.com/thorchain/thornode
# Docker must be started...
make openapi
make protob-docker
make install
```

## Commands

`thornode --help`

```text
THORChain Network

Usage:
  THORChain [command]

Available Commands:
  add-genesis-account Add a genesis account to genesis.json
  collect-gentxs      Collect genesis txs and output a genesis.json file
  debug               Tool for helping with debugging your application
  ed25519             Generate an ed25519 keys
  export              Export state to JSON
  gentx               Generate a genesis tx carrying a self delegation
  help                Help about any command
  init                Initialize private validator, p2p, genesis, and application configuration files
  keys                Manage your application's keys
  migrate             Migrate genesis to a specified target version
  pubkey              Convert Proto3 JSON encoded pubkey to bech32 format
  query               Querying subcommands
  start               Run the full node
  status              Query remote node for status
  tendermint          Tendermint subcommands
  tx                  Transactions subcommands
  unsafe-reset-all    Resets the blockchain database, removes address book files, and resets data/priv_validator_state.json to the genesis state
  validate-genesis    validates the genesis file at the default location or at the location passed as an arg
  version             Print the application binary version information

Flags:
  -h, --help                help for THORChain
      --home string         directory for config and data (default "/Users/dev/.thornode")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors
```

### Popular Commands

#### Add new account

```text
thornode keys add {accountName}
```

#### Add existing account (via mnemonic)

```text
thornode keys add {accountName} --recover
```

#### List all accounts

```text
thornode keys list
```

## Send Transaction

### Create Transaction

```text
# Sender: thor1505gp5h48zd24uexrfgka70fg8ccedafsnj0e3
# Receiver: thor1gutjhrw4xlu3n3p3k3r0vexl2xknq3nv8ux9fy
# Amount: 1 RUNE (in 1e8 notation)
thorcli tx bank send thor1505gp5h48zd24uexrfgka70fg8ccedafsnj0e3 thor1gutjhrw4xlu3n3p3k3r0vexl2xknq3nv8ux9fy 100000000rune --chain-id thorchain-mainnet-v1 --node https://rpc.ninerealms.com:443 --gas 3000000 --generate-only >> tx_raw.json
```

This will output a file called `tx_raw.json`. Edit this file and change the `@type` field from `/cosmos.bank.v1beta1.MsgSend` to `/types.MsgSend`.

The `tx_raw.json` transaction should look like this:

```json
{
  "body": {
    "messages": [
      {
        "@type": "/types.MsgSend",
        "from_address": "thor1505gp5h48zd24uexrfgka70fg8ccedafsnj0e3",
        "to_address": "thor1gutjhrw4xlu3n3p3k3r0vexl2xknq3nv8ux9fy",
        "amount": [{ "denom": "rune", "amount": "100000000" }]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": { "amount": [], "gas_limit": "3000000", "payer": "", "granter": "" }
  },
  "signatures": []
}
```

### Sign Transaction

```text
thornode tx sign tx_raw.json --from {accountName} --sign-mode amino-json --chain-id thorchain-mainnet-v1 --node https://rpc.ninerealms.com:443 >> tx.json
```

This will output a file called `tx.json`.

### Broadcast Transaction

```text
thornode tx broadcast tx.json --chain-id thorchain-mainnet-v1 --node https://rpc.ninerealms.com:443 --gas auto
```
