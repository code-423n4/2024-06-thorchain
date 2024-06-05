# QuoteLoanOpenResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InboundAddress** | Pointer to **string** | the inbound address for the transaction on the source chain | [optional] 
**InboundConfirmationBlocks** | Pointer to **int64** | the approximate number of source chain blocks required before processing | [optional] 
**InboundConfirmationSeconds** | Pointer to **int64** | the approximate seconds for block confirmations required before processing | [optional] 
**OutboundDelayBlocks** | **int64** | the number of thorchain blocks the outbound will be delayed | 
**OutboundDelaySeconds** | **int64** | the approximate seconds for the outbound delay before it will be sent | 
**Fees** | [**QuoteFees**](QuoteFees.md) |  | 
**SlippageBps** | Pointer to **int64** | Deprecated - migrate to fees object. | [optional] 
**StreamingSlippageBps** | Pointer to **int64** | Deprecated - migrate to fees object. | [optional] 
**Router** | Pointer to **string** | the EVM chain router contract address | [optional] 
**Expiry** | **int64** | expiration timestamp in unix seconds | 
**Warning** | **string** | static warning message | 
**Notes** | **string** | chain specific quote notes | 
**DustThreshold** | Pointer to **string** | Defines the minimum transaction size for the chain in base units (sats, wei, uatom). Transactions with asset amounts lower than the dust_threshold are ignored. | [optional] 
**RecommendedMinAmountIn** | Pointer to **string** | The recommended minimum inbound amount for this transaction type &amp; inbound asset. Sending less than this amount could result in failed refunds. | [optional] 
**RecommendedGasRate** | **string** | the recommended gas rate to use for the inbound to ensure timely confirmation | 
**GasRateUnits** | **string** | the units of the recommended gas rate | 
**Memo** | Pointer to **string** | generated memo for the loan open | [optional] 
**ExpectedAmountOut** | **string** | the amount of the target asset the user can expect to receive after fees in 1e8 decimals | 
**ExpectedCollateralizationRatio** | **string** | the expected collateralization ratio in basis points | 
**ExpectedCollateralDeposited** | **string** | the expected amount of collateral increase on the loan | 
**ExpectedDebtIssued** | **string** | the expected amount of TOR debt increase on the loan | 
**StreamingSwapBlocks** | **int64** | The number of blocks involved in the streaming swaps during the open loan process. | 
**StreamingSwapSeconds** | **int64** | The approximate number of seconds taken by the streaming swaps involved in the open loan process. | 
**TotalOpenLoanSeconds** | **int64** | The total expected duration for a open loan, measured in seconds, which includes the time for inbound confirmation, the duration of streaming swaps, and any outbound delays. | 

## Methods

### NewQuoteLoanOpenResponse

`func NewQuoteLoanOpenResponse(outboundDelayBlocks int64, outboundDelaySeconds int64, fees QuoteFees, expiry int64, warning string, notes string, recommendedGasRate string, gasRateUnits string, expectedAmountOut string, expectedCollateralizationRatio string, expectedCollateralDeposited string, expectedDebtIssued string, streamingSwapBlocks int64, streamingSwapSeconds int64, totalOpenLoanSeconds int64, ) *QuoteLoanOpenResponse`

NewQuoteLoanOpenResponse instantiates a new QuoteLoanOpenResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQuoteLoanOpenResponseWithDefaults

`func NewQuoteLoanOpenResponseWithDefaults() *QuoteLoanOpenResponse`

NewQuoteLoanOpenResponseWithDefaults instantiates a new QuoteLoanOpenResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInboundAddress

`func (o *QuoteLoanOpenResponse) GetInboundAddress() string`

GetInboundAddress returns the InboundAddress field if non-nil, zero value otherwise.

### GetInboundAddressOk

`func (o *QuoteLoanOpenResponse) GetInboundAddressOk() (*string, bool)`

GetInboundAddressOk returns a tuple with the InboundAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundAddress

`func (o *QuoteLoanOpenResponse) SetInboundAddress(v string)`

SetInboundAddress sets InboundAddress field to given value.

### HasInboundAddress

`func (o *QuoteLoanOpenResponse) HasInboundAddress() bool`

HasInboundAddress returns a boolean if a field has been set.

### GetInboundConfirmationBlocks

`func (o *QuoteLoanOpenResponse) GetInboundConfirmationBlocks() int64`

GetInboundConfirmationBlocks returns the InboundConfirmationBlocks field if non-nil, zero value otherwise.

### GetInboundConfirmationBlocksOk

`func (o *QuoteLoanOpenResponse) GetInboundConfirmationBlocksOk() (*int64, bool)`

GetInboundConfirmationBlocksOk returns a tuple with the InboundConfirmationBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationBlocks

`func (o *QuoteLoanOpenResponse) SetInboundConfirmationBlocks(v int64)`

SetInboundConfirmationBlocks sets InboundConfirmationBlocks field to given value.

### HasInboundConfirmationBlocks

`func (o *QuoteLoanOpenResponse) HasInboundConfirmationBlocks() bool`

HasInboundConfirmationBlocks returns a boolean if a field has been set.

### GetInboundConfirmationSeconds

`func (o *QuoteLoanOpenResponse) GetInboundConfirmationSeconds() int64`

GetInboundConfirmationSeconds returns the InboundConfirmationSeconds field if non-nil, zero value otherwise.

### GetInboundConfirmationSecondsOk

`func (o *QuoteLoanOpenResponse) GetInboundConfirmationSecondsOk() (*int64, bool)`

GetInboundConfirmationSecondsOk returns a tuple with the InboundConfirmationSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationSeconds

`func (o *QuoteLoanOpenResponse) SetInboundConfirmationSeconds(v int64)`

SetInboundConfirmationSeconds sets InboundConfirmationSeconds field to given value.

### HasInboundConfirmationSeconds

`func (o *QuoteLoanOpenResponse) HasInboundConfirmationSeconds() bool`

HasInboundConfirmationSeconds returns a boolean if a field has been set.

### GetOutboundDelayBlocks

`func (o *QuoteLoanOpenResponse) GetOutboundDelayBlocks() int64`

GetOutboundDelayBlocks returns the OutboundDelayBlocks field if non-nil, zero value otherwise.

### GetOutboundDelayBlocksOk

`func (o *QuoteLoanOpenResponse) GetOutboundDelayBlocksOk() (*int64, bool)`

GetOutboundDelayBlocksOk returns a tuple with the OutboundDelayBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelayBlocks

`func (o *QuoteLoanOpenResponse) SetOutboundDelayBlocks(v int64)`

SetOutboundDelayBlocks sets OutboundDelayBlocks field to given value.


### GetOutboundDelaySeconds

`func (o *QuoteLoanOpenResponse) GetOutboundDelaySeconds() int64`

GetOutboundDelaySeconds returns the OutboundDelaySeconds field if non-nil, zero value otherwise.

### GetOutboundDelaySecondsOk

`func (o *QuoteLoanOpenResponse) GetOutboundDelaySecondsOk() (*int64, bool)`

GetOutboundDelaySecondsOk returns a tuple with the OutboundDelaySeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelaySeconds

`func (o *QuoteLoanOpenResponse) SetOutboundDelaySeconds(v int64)`

SetOutboundDelaySeconds sets OutboundDelaySeconds field to given value.


### GetFees

`func (o *QuoteLoanOpenResponse) GetFees() QuoteFees`

GetFees returns the Fees field if non-nil, zero value otherwise.

### GetFeesOk

`func (o *QuoteLoanOpenResponse) GetFeesOk() (*QuoteFees, bool)`

GetFeesOk returns a tuple with the Fees field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFees

`func (o *QuoteLoanOpenResponse) SetFees(v QuoteFees)`

SetFees sets Fees field to given value.


### GetSlippageBps

`func (o *QuoteLoanOpenResponse) GetSlippageBps() int64`

GetSlippageBps returns the SlippageBps field if non-nil, zero value otherwise.

### GetSlippageBpsOk

`func (o *QuoteLoanOpenResponse) GetSlippageBpsOk() (*int64, bool)`

GetSlippageBpsOk returns a tuple with the SlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlippageBps

`func (o *QuoteLoanOpenResponse) SetSlippageBps(v int64)`

SetSlippageBps sets SlippageBps field to given value.

### HasSlippageBps

`func (o *QuoteLoanOpenResponse) HasSlippageBps() bool`

HasSlippageBps returns a boolean if a field has been set.

### GetStreamingSlippageBps

`func (o *QuoteLoanOpenResponse) GetStreamingSlippageBps() int64`

GetStreamingSlippageBps returns the StreamingSlippageBps field if non-nil, zero value otherwise.

### GetStreamingSlippageBpsOk

`func (o *QuoteLoanOpenResponse) GetStreamingSlippageBpsOk() (*int64, bool)`

GetStreamingSlippageBpsOk returns a tuple with the StreamingSlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSlippageBps

`func (o *QuoteLoanOpenResponse) SetStreamingSlippageBps(v int64)`

SetStreamingSlippageBps sets StreamingSlippageBps field to given value.

### HasStreamingSlippageBps

`func (o *QuoteLoanOpenResponse) HasStreamingSlippageBps() bool`

HasStreamingSlippageBps returns a boolean if a field has been set.

### GetRouter

`func (o *QuoteLoanOpenResponse) GetRouter() string`

GetRouter returns the Router field if non-nil, zero value otherwise.

### GetRouterOk

`func (o *QuoteLoanOpenResponse) GetRouterOk() (*string, bool)`

GetRouterOk returns a tuple with the Router field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouter

`func (o *QuoteLoanOpenResponse) SetRouter(v string)`

SetRouter sets Router field to given value.

### HasRouter

`func (o *QuoteLoanOpenResponse) HasRouter() bool`

HasRouter returns a boolean if a field has been set.

### GetExpiry

`func (o *QuoteLoanOpenResponse) GetExpiry() int64`

GetExpiry returns the Expiry field if non-nil, zero value otherwise.

### GetExpiryOk

`func (o *QuoteLoanOpenResponse) GetExpiryOk() (*int64, bool)`

GetExpiryOk returns a tuple with the Expiry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiry

`func (o *QuoteLoanOpenResponse) SetExpiry(v int64)`

SetExpiry sets Expiry field to given value.


### GetWarning

`func (o *QuoteLoanOpenResponse) GetWarning() string`

GetWarning returns the Warning field if non-nil, zero value otherwise.

### GetWarningOk

`func (o *QuoteLoanOpenResponse) GetWarningOk() (*string, bool)`

GetWarningOk returns a tuple with the Warning field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarning

`func (o *QuoteLoanOpenResponse) SetWarning(v string)`

SetWarning sets Warning field to given value.


### GetNotes

`func (o *QuoteLoanOpenResponse) GetNotes() string`

GetNotes returns the Notes field if non-nil, zero value otherwise.

### GetNotesOk

`func (o *QuoteLoanOpenResponse) GetNotesOk() (*string, bool)`

GetNotesOk returns a tuple with the Notes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotes

`func (o *QuoteLoanOpenResponse) SetNotes(v string)`

SetNotes sets Notes field to given value.


### GetDustThreshold

`func (o *QuoteLoanOpenResponse) GetDustThreshold() string`

GetDustThreshold returns the DustThreshold field if non-nil, zero value otherwise.

### GetDustThresholdOk

`func (o *QuoteLoanOpenResponse) GetDustThresholdOk() (*string, bool)`

GetDustThresholdOk returns a tuple with the DustThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDustThreshold

`func (o *QuoteLoanOpenResponse) SetDustThreshold(v string)`

SetDustThreshold sets DustThreshold field to given value.

### HasDustThreshold

`func (o *QuoteLoanOpenResponse) HasDustThreshold() bool`

HasDustThreshold returns a boolean if a field has been set.

### GetRecommendedMinAmountIn

`func (o *QuoteLoanOpenResponse) GetRecommendedMinAmountIn() string`

GetRecommendedMinAmountIn returns the RecommendedMinAmountIn field if non-nil, zero value otherwise.

### GetRecommendedMinAmountInOk

`func (o *QuoteLoanOpenResponse) GetRecommendedMinAmountInOk() (*string, bool)`

GetRecommendedMinAmountInOk returns a tuple with the RecommendedMinAmountIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedMinAmountIn

`func (o *QuoteLoanOpenResponse) SetRecommendedMinAmountIn(v string)`

SetRecommendedMinAmountIn sets RecommendedMinAmountIn field to given value.

### HasRecommendedMinAmountIn

`func (o *QuoteLoanOpenResponse) HasRecommendedMinAmountIn() bool`

HasRecommendedMinAmountIn returns a boolean if a field has been set.

### GetRecommendedGasRate

`func (o *QuoteLoanOpenResponse) GetRecommendedGasRate() string`

GetRecommendedGasRate returns the RecommendedGasRate field if non-nil, zero value otherwise.

### GetRecommendedGasRateOk

`func (o *QuoteLoanOpenResponse) GetRecommendedGasRateOk() (*string, bool)`

GetRecommendedGasRateOk returns a tuple with the RecommendedGasRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedGasRate

`func (o *QuoteLoanOpenResponse) SetRecommendedGasRate(v string)`

SetRecommendedGasRate sets RecommendedGasRate field to given value.


### GetGasRateUnits

`func (o *QuoteLoanOpenResponse) GetGasRateUnits() string`

GetGasRateUnits returns the GasRateUnits field if non-nil, zero value otherwise.

### GetGasRateUnitsOk

`func (o *QuoteLoanOpenResponse) GetGasRateUnitsOk() (*string, bool)`

GetGasRateUnitsOk returns a tuple with the GasRateUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasRateUnits

`func (o *QuoteLoanOpenResponse) SetGasRateUnits(v string)`

SetGasRateUnits sets GasRateUnits field to given value.


### GetMemo

`func (o *QuoteLoanOpenResponse) GetMemo() string`

GetMemo returns the Memo field if non-nil, zero value otherwise.

### GetMemoOk

`func (o *QuoteLoanOpenResponse) GetMemoOk() (*string, bool)`

GetMemoOk returns a tuple with the Memo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemo

`func (o *QuoteLoanOpenResponse) SetMemo(v string)`

SetMemo sets Memo field to given value.

### HasMemo

`func (o *QuoteLoanOpenResponse) HasMemo() bool`

HasMemo returns a boolean if a field has been set.

### GetExpectedAmountOut

`func (o *QuoteLoanOpenResponse) GetExpectedAmountOut() string`

GetExpectedAmountOut returns the ExpectedAmountOut field if non-nil, zero value otherwise.

### GetExpectedAmountOutOk

`func (o *QuoteLoanOpenResponse) GetExpectedAmountOutOk() (*string, bool)`

GetExpectedAmountOutOk returns a tuple with the ExpectedAmountOut field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedAmountOut

`func (o *QuoteLoanOpenResponse) SetExpectedAmountOut(v string)`

SetExpectedAmountOut sets ExpectedAmountOut field to given value.


### GetExpectedCollateralizationRatio

`func (o *QuoteLoanOpenResponse) GetExpectedCollateralizationRatio() string`

GetExpectedCollateralizationRatio returns the ExpectedCollateralizationRatio field if non-nil, zero value otherwise.

### GetExpectedCollateralizationRatioOk

`func (o *QuoteLoanOpenResponse) GetExpectedCollateralizationRatioOk() (*string, bool)`

GetExpectedCollateralizationRatioOk returns a tuple with the ExpectedCollateralizationRatio field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedCollateralizationRatio

`func (o *QuoteLoanOpenResponse) SetExpectedCollateralizationRatio(v string)`

SetExpectedCollateralizationRatio sets ExpectedCollateralizationRatio field to given value.


### GetExpectedCollateralDeposited

`func (o *QuoteLoanOpenResponse) GetExpectedCollateralDeposited() string`

GetExpectedCollateralDeposited returns the ExpectedCollateralDeposited field if non-nil, zero value otherwise.

### GetExpectedCollateralDepositedOk

`func (o *QuoteLoanOpenResponse) GetExpectedCollateralDepositedOk() (*string, bool)`

GetExpectedCollateralDepositedOk returns a tuple with the ExpectedCollateralDeposited field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedCollateralDeposited

`func (o *QuoteLoanOpenResponse) SetExpectedCollateralDeposited(v string)`

SetExpectedCollateralDeposited sets ExpectedCollateralDeposited field to given value.


### GetExpectedDebtIssued

`func (o *QuoteLoanOpenResponse) GetExpectedDebtIssued() string`

GetExpectedDebtIssued returns the ExpectedDebtIssued field if non-nil, zero value otherwise.

### GetExpectedDebtIssuedOk

`func (o *QuoteLoanOpenResponse) GetExpectedDebtIssuedOk() (*string, bool)`

GetExpectedDebtIssuedOk returns a tuple with the ExpectedDebtIssued field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedDebtIssued

`func (o *QuoteLoanOpenResponse) SetExpectedDebtIssued(v string)`

SetExpectedDebtIssued sets ExpectedDebtIssued field to given value.


### GetStreamingSwapBlocks

`func (o *QuoteLoanOpenResponse) GetStreamingSwapBlocks() int64`

GetStreamingSwapBlocks returns the StreamingSwapBlocks field if non-nil, zero value otherwise.

### GetStreamingSwapBlocksOk

`func (o *QuoteLoanOpenResponse) GetStreamingSwapBlocksOk() (*int64, bool)`

GetStreamingSwapBlocksOk returns a tuple with the StreamingSwapBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSwapBlocks

`func (o *QuoteLoanOpenResponse) SetStreamingSwapBlocks(v int64)`

SetStreamingSwapBlocks sets StreamingSwapBlocks field to given value.


### GetStreamingSwapSeconds

`func (o *QuoteLoanOpenResponse) GetStreamingSwapSeconds() int64`

GetStreamingSwapSeconds returns the StreamingSwapSeconds field if non-nil, zero value otherwise.

### GetStreamingSwapSecondsOk

`func (o *QuoteLoanOpenResponse) GetStreamingSwapSecondsOk() (*int64, bool)`

GetStreamingSwapSecondsOk returns a tuple with the StreamingSwapSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSwapSeconds

`func (o *QuoteLoanOpenResponse) SetStreamingSwapSeconds(v int64)`

SetStreamingSwapSeconds sets StreamingSwapSeconds field to given value.


### GetTotalOpenLoanSeconds

`func (o *QuoteLoanOpenResponse) GetTotalOpenLoanSeconds() int64`

GetTotalOpenLoanSeconds returns the TotalOpenLoanSeconds field if non-nil, zero value otherwise.

### GetTotalOpenLoanSecondsOk

`func (o *QuoteLoanOpenResponse) GetTotalOpenLoanSecondsOk() (*int64, bool)`

GetTotalOpenLoanSecondsOk returns a tuple with the TotalOpenLoanSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalOpenLoanSeconds

`func (o *QuoteLoanOpenResponse) SetTotalOpenLoanSeconds(v int64)`

SetTotalOpenLoanSeconds sets TotalOpenLoanSeconds field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


