# Quick Start Guide

Lending allows users to deposit native collateral, and then create a debt at a collateralization ratio `CR` (collateralization ratio). The debt is always denominated in USD (aka `TOR`) regardless of what L1 asset the user receives.

```admonish indo
[Streaming swaps](../swap-guide/streaming-swaps.md) is enabled for lending.
```

## Open a Loan Quote

Lending Quote endpoints have been created to simplify the implementation process.

**Request:** Loan quote using 1 BTC as collateral, target debt asset is USDT at 0XDAC17F958D2EE523A2206206994597C13D831EC7

[https://thornode.ninerealms.com/thorchain/quote/loan/open?from_asset=BTC.BTC\&amount=10000000\&to_asset=ETH.USDT-0xdac17f958d2ee523a2206206994597c13d831ec7\&destination=0xe7062003a7be4df3a86127293a0d6b1f54c04220](https://thornode.ninerealms.com/thorchain/quote/loan/open?from_asset=BTC.BTC&amount=10000000&to_asset=ETH.USDT-0xdac17f958d2ee523a2206206994597c13d831ec7&destination=0xe7062003a7be4df3a86127293a0d6b1f54c04220)

**Response:**

```json
{
  "dust_threshold": "10000",
  "expected_amount_out": "112302802900",
  "expected_collateral_deposited": "9997829",
  "expected_collateralization_ratio": "31467",
  "expected_debt_issued": "112887730000",
  "expiry": 1698901398,
  "fees": {
    "asset": "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7",
    "liquidity": "114988700",
    "outbound": "444599700",
    "slippage_bps": 10,
    "total": "559588400",
    "total_bps": 49
  },
  "inbound_address": "bc1qmed4v5am2hcg8furkeff2pczdnt0qu4flke420",
  "inbound_confirmation_blocks": 1,
  "inbound_confirmation_seconds": 600,
  "memo": "$+:ETH.USDT:0xe7062003a7be4df3a86127293a0d6b1f54c04220",
  "notes": "First output should be to inbound_address, second output should be change back to self, third output should be OP_RETURN, limited to 80 bytes. Do not send below the dust threshold. Do not use exotic spend scripts, locks or address formats (P2WSH with Bech32 address format preferred).",
  "outbound_delay_blocks": 3,
  "outbound_delay_seconds": 18,
  "recommended_min_amount_in": "156000",
  "warning": "Do not cache this response. Do not send funds after the expiry."
}
```

_If you send 1 BTC to `bc1q2hldv0pmy9mcpddj2qrvdgcx6pw6h6h7gqytwy` with the_ [_memo_](../concepts/memos.md#open-loan) _`$+:ETH.USDT:0xe7062003a7be4df3a86127293a0d6b1f54c04220` you will receive approx. 1128.8773 USDT debt sent to `0xe7062003a7be4df3a86127293a0d6b1f54c04220` with a CR of 314.6% and will incur 49 basis points (0.49%) slippage._

```admonish danger
The `Inbound_Address` changes regularly, do not cache!
```

```admonish warning
Loans cannot be repaid until a minimum time has passed, as determined by [LOANREPAYMENTMATURITY](https://thornode.ninerealms.com/thorchain/mimir), which is currently set as the current block height plus LOANREPAYMENTMATURITY. Currently, LOANREPAYMENTMATURITY is set to 432,000 blocks, equivalent to 30 days. Increasing the collateral on an existing loan to obtain additional debit resets the period.
```

## **Close a Loan**

**Request**: Repay a loan using USDT where BTC.BTC was used as colloteral. Note any asset can be used to repay a loan. [https://thornode.ninerealms.com/thorchain/quote/loan/close?from_asset=BTC.BTC\&amount=114947930000\&to_asset=BTC.BTC\&loan_owner=bc1q089j003xwj07uuavt2as5r45a95k5zzrhe4ac3](https://thornode.ninerealms.com/thorchain/quote/loan/close?from_asset=BTC.BTC&amount=114947930000&to_asset=BTC.BTC&loan_owner=bc1q089j003xwj07uuavt2as5r45a95k5zzrhe4ac3)

**Response:**

```json
{
  "dust_threshold": "10000",
  "expected_amount_out": "9985158",
  "expected_collateral_withdrawn": "9997123",
  "expected_debt_repaid": "390985054444080",
  "expiry": 1698897875,
  "fees": {
    "asset": "BTC.BTC",
    "liquidity": "38196994221",
    "outbound": "7500",
    "slippage_bps": 4347,
    "total": "38197001721",
    "total_bps": 38253777
  },
  "inbound_address": "bc1q69vcdslg0vfy4ne3nj7te5p9cvu2y4vq8t3x99",
  "inbound_confirmation_blocks": 192,
  "inbound_confirmation_seconds": 115200,
  "memo": "$-:BTC.BTC:bc1q089j003xwj07uuavt2as5r45a95k5zzrhe4ac3",
  "notes": "First output should be to inbound_address, second output should be change back to self, third output should be OP_RETURN, limited to 80 bytes. Do not send below the dust threshold. Do not use exotic spend scripts, locks or address formats (P2WSH with Bech32 address format preferred).",
  "outbound_delay_blocks": 12,
  "outbound_delay_seconds": 72,
  "recommended_min_amount_in": "30000",
  "warning": "Do not cache this response. Do not send funds after the expiry."
}
```

_If you send 1149.47 USDT with a memo `$-:BTC.BTC:bc1q089j003xwj07uuavt2as5r45a95k5zzrhe4ac3` of you will repay your loan down._

### **Borrowers Position**

**Request:**\
Get brower's positin in the BTC pool who tool out a loan from `bc1q089j003xwj07uuavt2as5r45a95k5zzrhe4ac3`\
[https://thornode.ninerealms.com/thorchain/pool/BTC.BTC/borrower/bc1q089j003xwj07uuavt2as5r45a95k5zzrhe4ac3](https://thornode.ninerealms.com/thorchain/pool/BTC.BTC/borrower/bc1q089j003xwj07uuavt2as5r45a95k5zzrhe4ac3)\

**Response:**

```json
{
  "asset": "BTC.BTC",
  "collateral_current": "9997123",
  "collateral_deposited": "9997123",
  "collateral_withdrawn": "0",
  "debt_current": "114947930000",
  "debt_issued": "114947930000",
  "debt_repaid": "0",
  "last_open_height": 12252923,
  "last_repay_height": 0,
  "owner": "bc1q089j003xwj07uuavt2as5r45a95k5zzrhe4ac3"
}
```

_The borrower has provided 0.0997 BTC and has a current TOR debt of $1149.78. No repayments have been yet._

### Support

Developers experiencing issues with these APIs can go to the [Developer Discord](https://discord.gg/2Vw3RsQ7) for assistance. Interface developers should subscribe to the #interface-alerts channel for information pertinent to the endpoints and functionality discussed here.
