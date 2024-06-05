# Multisig

## Setup Multisig

First, collect the pubkeys that will be part of the multisig. They can be printed using `thorcli`:

```text
thornode keys show person1 --pubkey
```

Then share the pubkey with the other parties. Each party can add these pubkeys:

```text
thornode keys add person2 --pubkey {pubkey}
```

Each party can create the multisig (here a 2/3):

```text
thornode keys add multisig --multisig person1,person2,person3 --multisig-threshold 2
```

### Create Transaction

Any of the parties can create the raw transaction:

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

The transaction needs to be signed by 2 of the 3 parties (as configured above, when setting up the multisig).

#### From Person 1

```text
thornode tx sign --from person1 --multisig multisig tx_raw.json --sign-mode amino-json --chain-id thorchain-mainnet-v1 --node https://rpc.ninerealms.com:443 >> tx_signed_1.json
```

This will output a file called `tx_signed_1.json`.

#### From Person 2

```text
thornode tx sign --from person2 --multisig multisig tx_raw.json --sign-mode amino-json --chain-id thorchain-mainnet-v1 --node https://rpc.ninerealms.com:443 >> tx_signed_2.json
```

This will output a file called `tx_signed_2.json`.

### Build Transaction

#### Gather Signatures

The party, who wants to broadcast the transaction, needs to gather all json signature files from the other parties.

#### Multisig Sign

First, get the sequence and account number for the multisig address:

```text
curl https://thornode.ninerealms.com/cosmos/auth/v1beta1/accounts/thor1505gp5h48zd24uexrfgka70fg8ccedafsnj0e3
```

Then combine the signatures into a single one (make sure to update the account number `-a` and the sequence number `-s`:

```text
# Account number: 33401 (see curl output)
# Sequence number: 0 (see curl output)
thornode tx multisign tx_raw.json multisig tx_signed_1.json tx_signed_2.json -a 33401 -s 0 --from multisig --chain-id thorchain-mainnet-v1 --node https://rpc.ninerealms.com:443 >> tx.json
```

This will output a final file called `tx.json`.

### Broadcast Transaction

```text
thornode tx broadcast tx.json --chain-id thorchain-mainnet-v1 --node https://rpc.ninerealms.com:443 --gas auto
```

## THORSafe

```admonish info
THORSafe does not support Ledger yet!
```

THORSafe is a multisig frontend (developed by THORSwap): [https://app.thorswap.finance/thorsafe](https://app.thorswap.finance/thorsafe)
