# QuoteLoanCloseResponse

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
**RecommendedGasRate** | Pointer to **string** | the recommended gas rate to use for the inbound to ensure timely confirmation | [optional] 
**GasRateUnits** | Pointer to **string** | the units of the recommended gas rate | [optional] 
**Memo** | **string** | generated memo for the loan close | 
**ExpectedAmountOut** | **string** | the amount of collateral asset the user can expect to receive after fees in 1e8 decimals | 
**ExpectedAmountIn** | **string** | The quantity of the repayment asset to be sent by the user, calculated as the desired percentage of the loan&#39;s value, expressed in units of 1e8 | 
**ExpectedCollateralWithdrawn** | **string** | the expected amount of collateral decrease on the loan | 
**ExpectedDebtRepaid** | **string** | the expected amount of TOR debt decrease on the loan | 
**StreamingSwapBlocks** | **int64** | The number of blocks involved in the streaming swaps during the repayment process. | 
**StreamingSwapSeconds** | **int64** | The approximate number of seconds taken by the streaming swaps involved in the repayment process. | 
**TotalRepaySeconds** | **int64** | The total expected duration for a repayment, measured in seconds, which includes the time for inbound confirmation, the duration of streaming swaps, and any outbound delays. | 

## Methods

### NewQuoteLoanCloseResponse

`func NewQuoteLoanCloseResponse(outboundDelayBlocks int64, outboundDelaySeconds int64, fees QuoteFees, expiry int64, warning string, notes string, memo string, expectedAmountOut string, expectedAmountIn string, expectedCollateralWithdrawn string, expectedDebtRepaid string, streamingSwapBlocks int64, streamingSwapSeconds int64, totalRepaySeconds int64, ) *QuoteLoanCloseResponse`

NewQuoteLoanCloseResponse instantiates a new QuoteLoanCloseResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQuoteLoanCloseResponseWithDefaults

`func NewQuoteLoanCloseResponseWithDefaults() *QuoteLoanCloseResponse`

NewQuoteLoanCloseResponseWithDefaults instantiates a new QuoteLoanCloseResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInboundAddress

`func (o *QuoteLoanCloseResponse) GetInboundAddress() string`

GetInboundAddress returns the InboundAddress field if non-nil, zero value otherwise.

### GetInboundAddressOk

`func (o *QuoteLoanCloseResponse) GetInboundAddressOk() (*string, bool)`

GetInboundAddressOk returns a tuple with the InboundAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundAddress

`func (o *QuoteLoanCloseResponse) SetInboundAddress(v string)`

SetInboundAddress sets InboundAddress field to given value.

### HasInboundAddress

`func (o *QuoteLoanCloseResponse) HasInboundAddress() bool`

HasInboundAddress returns a boolean if a field has been set.

### GetInboundConfirmationBlocks

`func (o *QuoteLoanCloseResponse) GetInboundConfirmationBlocks() int64`

GetInboundConfirmationBlocks returns the InboundConfirmationBlocks field if non-nil, zero value otherwise.

### GetInboundConfirmationBlocksOk

`func (o *QuoteLoanCloseResponse) GetInboundConfirmationBlocksOk() (*int64, bool)`

GetInboundConfirmationBlocksOk returns a tuple with the InboundConfirmationBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationBlocks

`func (o *QuoteLoanCloseResponse) SetInboundConfirmationBlocks(v int64)`

SetInboundConfirmationBlocks sets InboundConfirmationBlocks field to given value.

### HasInboundConfirmationBlocks

`func (o *QuoteLoanCloseResponse) HasInboundConfirmationBlocks() bool`

HasInboundConfirmationBlocks returns a boolean if a field has been set.

### GetInboundConfirmationSeconds

`func (o *QuoteLoanCloseResponse) GetInboundConfirmationSeconds() int64`

GetInboundConfirmationSeconds returns the InboundConfirmationSeconds field if non-nil, zero value otherwise.

### GetInboundConfirmationSecondsOk

`func (o *QuoteLoanCloseResponse) GetInboundConfirmationSecondsOk() (*int64, bool)`

GetInboundConfirmationSecondsOk returns a tuple with the InboundConfirmationSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationSeconds

`func (o *QuoteLoanCloseResponse) SetInboundConfirmationSeconds(v int64)`

SetInboundConfirmationSeconds sets InboundConfirmationSeconds field to given value.

### HasInboundConfirmationSeconds

`func (o *QuoteLoanCloseResponse) HasInboundConfirmationSeconds() bool`

HasInboundConfirmationSeconds returns a boolean if a field has been set.

### GetOutboundDelayBlocks

`func (o *QuoteLoanCloseResponse) GetOutboundDelayBlocks() int64`

GetOutboundDelayBlocks returns the OutboundDelayBlocks field if non-nil, zero value otherwise.

### GetOutboundDelayBlocksOk

`func (o *QuoteLoanCloseResponse) GetOutboundDelayBlocksOk() (*int64, bool)`

GetOutboundDelayBlocksOk returns a tuple with the OutboundDelayBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelayBlocks

`func (o *QuoteLoanCloseResponse) SetOutboundDelayBlocks(v int64)`

SetOutboundDelayBlocks sets OutboundDelayBlocks field to given value.


### GetOutboundDelaySeconds

`func (o *QuoteLoanCloseResponse) GetOutboundDelaySeconds() int64`

GetOutboundDelaySeconds returns the OutboundDelaySeconds field if non-nil, zero value otherwise.

### GetOutboundDelaySecondsOk

`func (o *QuoteLoanCloseResponse) GetOutboundDelaySecondsOk() (*int64, bool)`

GetOutboundDelaySecondsOk returns a tuple with the OutboundDelaySeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelaySeconds

`func (o *QuoteLoanCloseResponse) SetOutboundDelaySeconds(v int64)`

SetOutboundDelaySeconds sets OutboundDelaySeconds field to given value.


### GetFees

`func (o *QuoteLoanCloseResponse) GetFees() QuoteFees`

GetFees returns the Fees field if non-nil, zero value otherwise.

### GetFeesOk

`func (o *QuoteLoanCloseResponse) GetFeesOk() (*QuoteFees, bool)`

GetFeesOk returns a tuple with the Fees field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFees

`func (o *QuoteLoanCloseResponse) SetFees(v QuoteFees)`

SetFees sets Fees field to given value.


### GetSlippageBps

`func (o *QuoteLoanCloseResponse) GetSlippageBps() int64`

GetSlippageBps returns the SlippageBps field if non-nil, zero value otherwise.

### GetSlippageBpsOk

`func (o *QuoteLoanCloseResponse) GetSlippageBpsOk() (*int64, bool)`

GetSlippageBpsOk returns a tuple with the SlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlippageBps

`func (o *QuoteLoanCloseResponse) SetSlippageBps(v int64)`

SetSlippageBps sets SlippageBps field to given value.

### HasSlippageBps

`func (o *QuoteLoanCloseResponse) HasSlippageBps() bool`

HasSlippageBps returns a boolean if a field has been set.

### GetStreamingSlippageBps

`func (o *QuoteLoanCloseResponse) GetStreamingSlippageBps() int64`

GetStreamingSlippageBps returns the StreamingSlippageBps field if non-nil, zero value otherwise.

### GetStreamingSlippageBpsOk

`func (o *QuoteLoanCloseResponse) GetStreamingSlippageBpsOk() (*int64, bool)`

GetStreamingSlippageBpsOk returns a tuple with the StreamingSlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSlippageBps

`func (o *QuoteLoanCloseResponse) SetStreamingSlippageBps(v int64)`

SetStreamingSlippageBps sets StreamingSlippageBps field to given value.

### HasStreamingSlippageBps

`func (o *QuoteLoanCloseResponse) HasStreamingSlippageBps() bool`

HasStreamingSlippageBps returns a boolean if a field has been set.

### GetRouter

`func (o *QuoteLoanCloseResponse) GetRouter() string`

GetRouter returns the Router field if non-nil, zero value otherwise.

### GetRouterOk

`func (o *QuoteLoanCloseResponse) GetRouterOk() (*string, bool)`

GetRouterOk returns a tuple with the Router field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouter

`func (o *QuoteLoanCloseResponse) SetRouter(v string)`

SetRouter sets Router field to given value.

### HasRouter

`func (o *QuoteLoanCloseResponse) HasRouter() bool`

HasRouter returns a boolean if a field has been set.

### GetExpiry

`func (o *QuoteLoanCloseResponse) GetExpiry() int64`

GetExpiry returns the Expiry field if non-nil, zero value otherwise.

### GetExpiryOk

`func (o *QuoteLoanCloseResponse) GetExpiryOk() (*int64, bool)`

GetExpiryOk returns a tuple with the Expiry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiry

`func (o *QuoteLoanCloseResponse) SetExpiry(v int64)`

SetExpiry sets Expiry field to given value.


### GetWarning

`func (o *QuoteLoanCloseResponse) GetWarning() string`

GetWarning returns the Warning field if non-nil, zero value otherwise.

### GetWarningOk

`func (o *QuoteLoanCloseResponse) GetWarningOk() (*string, bool)`

GetWarningOk returns a tuple with the Warning field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarning

`func (o *QuoteLoanCloseResponse) SetWarning(v string)`

SetWarning sets Warning field to given value.


### GetNotes

`func (o *QuoteLoanCloseResponse) GetNotes() string`

GetNotes returns the Notes field if non-nil, zero value otherwise.

### GetNotesOk

`func (o *QuoteLoanCloseResponse) GetNotesOk() (*string, bool)`

GetNotesOk returns a tuple with the Notes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotes

`func (o *QuoteLoanCloseResponse) SetNotes(v string)`

SetNotes sets Notes field to given value.


### GetDustThreshold

`func (o *QuoteLoanCloseResponse) GetDustThreshold() string`

GetDustThreshold returns the DustThreshold field if non-nil, zero value otherwise.

### GetDustThresholdOk

`func (o *QuoteLoanCloseResponse) GetDustThresholdOk() (*string, bool)`

GetDustThresholdOk returns a tuple with the DustThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDustThreshold

`func (o *QuoteLoanCloseResponse) SetDustThreshold(v string)`

SetDustThreshold sets DustThreshold field to given value.

### HasDustThreshold

`func (o *QuoteLoanCloseResponse) HasDustThreshold() bool`

HasDustThreshold returns a boolean if a field has been set.

### GetRecommendedMinAmountIn

`func (o *QuoteLoanCloseResponse) GetRecommendedMinAmountIn() string`

GetRecommendedMinAmountIn returns the RecommendedMinAmountIn field if non-nil, zero value otherwise.

### GetRecommendedMinAmountInOk

`func (o *QuoteLoanCloseResponse) GetRecommendedMinAmountInOk() (*string, bool)`

GetRecommendedMinAmountInOk returns a tuple with the RecommendedMinAmountIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedMinAmountIn

`func (o *QuoteLoanCloseResponse) SetRecommendedMinAmountIn(v string)`

SetRecommendedMinAmountIn sets RecommendedMinAmountIn field to given value.

### HasRecommendedMinAmountIn

`func (o *QuoteLoanCloseResponse) HasRecommendedMinAmountIn() bool`

HasRecommendedMinAmountIn returns a boolean if a field has been set.

### GetRecommendedGasRate

`func (o *QuoteLoanCloseResponse) GetRecommendedGasRate() string`

GetRecommendedGasRate returns the RecommendedGasRate field if non-nil, zero value otherwise.

### GetRecommendedGasRateOk

`func (o *QuoteLoanCloseResponse) GetRecommendedGasRateOk() (*string, bool)`

GetRecommendedGasRateOk returns a tuple with the RecommendedGasRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedGasRate

`func (o *QuoteLoanCloseResponse) SetRecommendedGasRate(v string)`

SetRecommendedGasRate sets RecommendedGasRate field to given value.

### HasRecommendedGasRate

`func (o *QuoteLoanCloseResponse) HasRecommendedGasRate() bool`

HasRecommendedGasRate returns a boolean if a field has been set.

### GetGasRateUnits

`func (o *QuoteLoanCloseResponse) GetGasRateUnits() string`

GetGasRateUnits returns the GasRateUnits field if non-nil, zero value otherwise.

### GetGasRateUnitsOk

`func (o *QuoteLoanCloseResponse) GetGasRateUnitsOk() (*string, bool)`

GetGasRateUnitsOk returns a tuple with the GasRateUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasRateUnits

`func (o *QuoteLoanCloseResponse) SetGasRateUnits(v string)`

SetGasRateUnits sets GasRateUnits field to given value.

### HasGasRateUnits

`func (o *QuoteLoanCloseResponse) HasGasRateUnits() bool`

HasGasRateUnits returns a boolean if a field has been set.

### GetMemo

`func (o *QuoteLoanCloseResponse) GetMemo() string`

GetMemo returns the Memo field if non-nil, zero value otherwise.

### GetMemoOk

`func (o *QuoteLoanCloseResponse) GetMemoOk() (*string, bool)`

GetMemoOk returns a tuple with the Memo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemo

`func (o *QuoteLoanCloseResponse) SetMemo(v string)`

SetMemo sets Memo field to given value.


### GetExpectedAmountOut

`func (o *QuoteLoanCloseResponse) GetExpectedAmountOut() string`

GetExpectedAmountOut returns the ExpectedAmountOut field if non-nil, zero value otherwise.

### GetExpectedAmountOutOk

`func (o *QuoteLoanCloseResponse) GetExpectedAmountOutOk() (*string, bool)`

GetExpectedAmountOutOk returns a tuple with the ExpectedAmountOut field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedAmountOut

`func (o *QuoteLoanCloseResponse) SetExpectedAmountOut(v string)`

SetExpectedAmountOut sets ExpectedAmountOut field to given value.


### GetExpectedAmountIn

`func (o *QuoteLoanCloseResponse) GetExpectedAmountIn() string`

GetExpectedAmountIn returns the ExpectedAmountIn field if non-nil, zero value otherwise.

### GetExpectedAmountInOk

`func (o *QuoteLoanCloseResponse) GetExpectedAmountInOk() (*string, bool)`

GetExpectedAmountInOk returns a tuple with the ExpectedAmountIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedAmountIn

`func (o *QuoteLoanCloseResponse) SetExpectedAmountIn(v string)`

SetExpectedAmountIn sets ExpectedAmountIn field to given value.


### GetExpectedCollateralWithdrawn

`func (o *QuoteLoanCloseResponse) GetExpectedCollateralWithdrawn() string`

GetExpectedCollateralWithdrawn returns the ExpectedCollateralWithdrawn field if non-nil, zero value otherwise.

### GetExpectedCollateralWithdrawnOk

`func (o *QuoteLoanCloseResponse) GetExpectedCollateralWithdrawnOk() (*string, bool)`

GetExpectedCollateralWithdrawnOk returns a tuple with the ExpectedCollateralWithdrawn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedCollateralWithdrawn

`func (o *QuoteLoanCloseResponse) SetExpectedCollateralWithdrawn(v string)`

SetExpectedCollateralWithdrawn sets ExpectedCollateralWithdrawn field to given value.


### GetExpectedDebtRepaid

`func (o *QuoteLoanCloseResponse) GetExpectedDebtRepaid() string`

GetExpectedDebtRepaid returns the ExpectedDebtRepaid field if non-nil, zero value otherwise.

### GetExpectedDebtRepaidOk

`func (o *QuoteLoanCloseResponse) GetExpectedDebtRepaidOk() (*string, bool)`

GetExpectedDebtRepaidOk returns a tuple with the ExpectedDebtRepaid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedDebtRepaid

`func (o *QuoteLoanCloseResponse) SetExpectedDebtRepaid(v string)`

SetExpectedDebtRepaid sets ExpectedDebtRepaid field to given value.


### GetStreamingSwapBlocks

`func (o *QuoteLoanCloseResponse) GetStreamingSwapBlocks() int64`

GetStreamingSwapBlocks returns the StreamingSwapBlocks field if non-nil, zero value otherwise.

### GetStreamingSwapBlocksOk

`func (o *QuoteLoanCloseResponse) GetStreamingSwapBlocksOk() (*int64, bool)`

GetStreamingSwapBlocksOk returns a tuple with the StreamingSwapBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSwapBlocks

`func (o *QuoteLoanCloseResponse) SetStreamingSwapBlocks(v int64)`

SetStreamingSwapBlocks sets StreamingSwapBlocks field to given value.


### GetStreamingSwapSeconds

`func (o *QuoteLoanCloseResponse) GetStreamingSwapSeconds() int64`

GetStreamingSwapSeconds returns the StreamingSwapSeconds field if non-nil, zero value otherwise.

### GetStreamingSwapSecondsOk

`func (o *QuoteLoanCloseResponse) GetStreamingSwapSecondsOk() (*int64, bool)`

GetStreamingSwapSecondsOk returns a tuple with the StreamingSwapSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSwapSeconds

`func (o *QuoteLoanCloseResponse) SetStreamingSwapSeconds(v int64)`

SetStreamingSwapSeconds sets StreamingSwapSeconds field to given value.


### GetTotalRepaySeconds

`func (o *QuoteLoanCloseResponse) GetTotalRepaySeconds() int64`

GetTotalRepaySeconds returns the TotalRepaySeconds field if non-nil, zero value otherwise.

### GetTotalRepaySecondsOk

`func (o *QuoteLoanCloseResponse) GetTotalRepaySecondsOk() (*int64, bool)`

GetTotalRepaySecondsOk returns a tuple with the TotalRepaySeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalRepaySeconds

`func (o *QuoteLoanCloseResponse) SetTotalRepaySeconds(v int64)`

SetTotalRepaySeconds sets TotalRepaySeconds field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


