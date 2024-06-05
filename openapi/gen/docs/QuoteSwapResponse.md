# QuoteSwapResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InboundAddress** | Pointer to **string** | the inbound address for the transaction on the source chain | [optional] 
**InboundConfirmationBlocks** | Pointer to **int64** | the approximate number of source chain blocks required before processing | [optional] 
**InboundConfirmationSeconds** | Pointer to **int64** | the approximate seconds for block confirmations required before processing | [optional] 
**OutboundDelayBlocks** | **int64** | the number of thorchain blocks the outbound will be delayed | 
**OutboundDelaySeconds** | **int64** | the approximate seconds for the outbound delay before it will be sent | 
**Fees** | [**QuoteFees**](QuoteFees.md) |  | 
**SlippageBps** | **int64** | Deprecated - migrate to fees object. | 
**StreamingSlippageBps** | **int64** | Deprecated - migrate to fees object. | 
**Router** | Pointer to **string** | the EVM chain router contract address | [optional] 
**Expiry** | **int64** | expiration timestamp in unix seconds | 
**Warning** | **string** | static warning message | 
**Notes** | **string** | chain specific quote notes | 
**DustThreshold** | Pointer to **string** | Defines the minimum transaction size for the chain in base units (sats, wei, uatom). Transactions with asset amounts lower than the dust_threshold are ignored. | [optional] 
**RecommendedMinAmountIn** | Pointer to **string** | The recommended minimum inbound amount for this transaction type &amp; inbound asset. Sending less than this amount could result in failed refunds. | [optional] 
**RecommendedGasRate** | Pointer to **string** | the recommended gas rate to use for the inbound to ensure timely confirmation | [optional] 
**GasRateUnits** | Pointer to **string** | the units of the recommended gas rate | [optional] 
**Memo** | Pointer to **string** | generated memo for the swap | [optional] 
**ExpectedAmountOut** | **string** | the amount of the target asset the user can expect to receive after fees | 
**ExpectedAmountOutStreaming** | **string** | Deprecated - expected_amount_out is streaming amount if interval provided. | 
**MaxStreamingQuantity** | Pointer to **int64** | the maximum amount of trades a streaming swap can do for a trade | [optional] 
**StreamingSwapBlocks** | Pointer to **int64** | the number of blocks the streaming swap will execute over | [optional] 
**StreamingSwapSeconds** | Pointer to **int64** | approx the number of seconds the streaming swap will execute over | [optional] 
**TotalSwapSeconds** | Pointer to **int64** | total number of seconds a swap is expected to take (inbound conf + streaming swap + outbound delay) | [optional] 

## Methods

### NewQuoteSwapResponse

`func NewQuoteSwapResponse(outboundDelayBlocks int64, outboundDelaySeconds int64, fees QuoteFees, slippageBps int64, streamingSlippageBps int64, expiry int64, warning string, notes string, expectedAmountOut string, expectedAmountOutStreaming string, ) *QuoteSwapResponse`

NewQuoteSwapResponse instantiates a new QuoteSwapResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQuoteSwapResponseWithDefaults

`func NewQuoteSwapResponseWithDefaults() *QuoteSwapResponse`

NewQuoteSwapResponseWithDefaults instantiates a new QuoteSwapResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInboundAddress

`func (o *QuoteSwapResponse) GetInboundAddress() string`

GetInboundAddress returns the InboundAddress field if non-nil, zero value otherwise.

### GetInboundAddressOk

`func (o *QuoteSwapResponse) GetInboundAddressOk() (*string, bool)`

GetInboundAddressOk returns a tuple with the InboundAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundAddress

`func (o *QuoteSwapResponse) SetInboundAddress(v string)`

SetInboundAddress sets InboundAddress field to given value.

### HasInboundAddress

`func (o *QuoteSwapResponse) HasInboundAddress() bool`

HasInboundAddress returns a boolean if a field has been set.

### GetInboundConfirmationBlocks

`func (o *QuoteSwapResponse) GetInboundConfirmationBlocks() int64`

GetInboundConfirmationBlocks returns the InboundConfirmationBlocks field if non-nil, zero value otherwise.

### GetInboundConfirmationBlocksOk

`func (o *QuoteSwapResponse) GetInboundConfirmationBlocksOk() (*int64, bool)`

GetInboundConfirmationBlocksOk returns a tuple with the InboundConfirmationBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationBlocks

`func (o *QuoteSwapResponse) SetInboundConfirmationBlocks(v int64)`

SetInboundConfirmationBlocks sets InboundConfirmationBlocks field to given value.

### HasInboundConfirmationBlocks

`func (o *QuoteSwapResponse) HasInboundConfirmationBlocks() bool`

HasInboundConfirmationBlocks returns a boolean if a field has been set.

### GetInboundConfirmationSeconds

`func (o *QuoteSwapResponse) GetInboundConfirmationSeconds() int64`

GetInboundConfirmationSeconds returns the InboundConfirmationSeconds field if non-nil, zero value otherwise.

### GetInboundConfirmationSecondsOk

`func (o *QuoteSwapResponse) GetInboundConfirmationSecondsOk() (*int64, bool)`

GetInboundConfirmationSecondsOk returns a tuple with the InboundConfirmationSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationSeconds

`func (o *QuoteSwapResponse) SetInboundConfirmationSeconds(v int64)`

SetInboundConfirmationSeconds sets InboundConfirmationSeconds field to given value.

### HasInboundConfirmationSeconds

`func (o *QuoteSwapResponse) HasInboundConfirmationSeconds() bool`

HasInboundConfirmationSeconds returns a boolean if a field has been set.

### GetOutboundDelayBlocks

`func (o *QuoteSwapResponse) GetOutboundDelayBlocks() int64`

GetOutboundDelayBlocks returns the OutboundDelayBlocks field if non-nil, zero value otherwise.

### GetOutboundDelayBlocksOk

`func (o *QuoteSwapResponse) GetOutboundDelayBlocksOk() (*int64, bool)`

GetOutboundDelayBlocksOk returns a tuple with the OutboundDelayBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelayBlocks

`func (o *QuoteSwapResponse) SetOutboundDelayBlocks(v int64)`

SetOutboundDelayBlocks sets OutboundDelayBlocks field to given value.


### GetOutboundDelaySeconds

`func (o *QuoteSwapResponse) GetOutboundDelaySeconds() int64`

GetOutboundDelaySeconds returns the OutboundDelaySeconds field if non-nil, zero value otherwise.

### GetOutboundDelaySecondsOk

`func (o *QuoteSwapResponse) GetOutboundDelaySecondsOk() (*int64, bool)`

GetOutboundDelaySecondsOk returns a tuple with the OutboundDelaySeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelaySeconds

`func (o *QuoteSwapResponse) SetOutboundDelaySeconds(v int64)`

SetOutboundDelaySeconds sets OutboundDelaySeconds field to given value.


### GetFees

`func (o *QuoteSwapResponse) GetFees() QuoteFees`

GetFees returns the Fees field if non-nil, zero value otherwise.

### GetFeesOk

`func (o *QuoteSwapResponse) GetFeesOk() (*QuoteFees, bool)`

GetFeesOk returns a tuple with the Fees field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFees

`func (o *QuoteSwapResponse) SetFees(v QuoteFees)`

SetFees sets Fees field to given value.


### GetSlippageBps

`func (o *QuoteSwapResponse) GetSlippageBps() int64`

GetSlippageBps returns the SlippageBps field if non-nil, zero value otherwise.

### GetSlippageBpsOk

`func (o *QuoteSwapResponse) GetSlippageBpsOk() (*int64, bool)`

GetSlippageBpsOk returns a tuple with the SlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlippageBps

`func (o *QuoteSwapResponse) SetSlippageBps(v int64)`

SetSlippageBps sets SlippageBps field to given value.


### GetStreamingSlippageBps

`func (o *QuoteSwapResponse) GetStreamingSlippageBps() int64`

GetStreamingSlippageBps returns the StreamingSlippageBps field if non-nil, zero value otherwise.

### GetStreamingSlippageBpsOk

`func (o *QuoteSwapResponse) GetStreamingSlippageBpsOk() (*int64, bool)`

GetStreamingSlippageBpsOk returns a tuple with the StreamingSlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSlippageBps

`func (o *QuoteSwapResponse) SetStreamingSlippageBps(v int64)`

SetStreamingSlippageBps sets StreamingSlippageBps field to given value.


### GetRouter

`func (o *QuoteSwapResponse) GetRouter() string`

GetRouter returns the Router field if non-nil, zero value otherwise.

### GetRouterOk

`func (o *QuoteSwapResponse) GetRouterOk() (*string, bool)`

GetRouterOk returns a tuple with the Router field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouter

`func (o *QuoteSwapResponse) SetRouter(v string)`

SetRouter sets Router field to given value.

### HasRouter

`func (o *QuoteSwapResponse) HasRouter() bool`

HasRouter returns a boolean if a field has been set.

### GetExpiry

`func (o *QuoteSwapResponse) GetExpiry() int64`

GetExpiry returns the Expiry field if non-nil, zero value otherwise.

### GetExpiryOk

`func (o *QuoteSwapResponse) GetExpiryOk() (*int64, bool)`

GetExpiryOk returns a tuple with the Expiry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiry

`func (o *QuoteSwapResponse) SetExpiry(v int64)`

SetExpiry sets Expiry field to given value.


### GetWarning

`func (o *QuoteSwapResponse) GetWarning() string`

GetWarning returns the Warning field if non-nil, zero value otherwise.

### GetWarningOk

`func (o *QuoteSwapResponse) GetWarningOk() (*string, bool)`

GetWarningOk returns a tuple with the Warning field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarning

`func (o *QuoteSwapResponse) SetWarning(v string)`

SetWarning sets Warning field to given value.


### GetNotes

`func (o *QuoteSwapResponse) GetNotes() string`

GetNotes returns the Notes field if non-nil, zero value otherwise.

### GetNotesOk

`func (o *QuoteSwapResponse) GetNotesOk() (*string, bool)`

GetNotesOk returns a tuple with the Notes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotes

`func (o *QuoteSwapResponse) SetNotes(v string)`

SetNotes sets Notes field to given value.


### GetDustThreshold

`func (o *QuoteSwapResponse) GetDustThreshold() string`

GetDustThreshold returns the DustThreshold field if non-nil, zero value otherwise.

### GetDustThresholdOk

`func (o *QuoteSwapResponse) GetDustThresholdOk() (*string, bool)`

GetDustThresholdOk returns a tuple with the DustThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDustThreshold

`func (o *QuoteSwapResponse) SetDustThreshold(v string)`

SetDustThreshold sets DustThreshold field to given value.

### HasDustThreshold

`func (o *QuoteSwapResponse) HasDustThreshold() bool`

HasDustThreshold returns a boolean if a field has been set.

### GetRecommendedMinAmountIn

`func (o *QuoteSwapResponse) GetRecommendedMinAmountIn() string`

GetRecommendedMinAmountIn returns the RecommendedMinAmountIn field if non-nil, zero value otherwise.

### GetRecommendedMinAmountInOk

`func (o *QuoteSwapResponse) GetRecommendedMinAmountInOk() (*string, bool)`

GetRecommendedMinAmountInOk returns a tuple with the RecommendedMinAmountIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedMinAmountIn

`func (o *QuoteSwapResponse) SetRecommendedMinAmountIn(v string)`

SetRecommendedMinAmountIn sets RecommendedMinAmountIn field to given value.

### HasRecommendedMinAmountIn

`func (o *QuoteSwapResponse) HasRecommendedMinAmountIn() bool`

HasRecommendedMinAmountIn returns a boolean if a field has been set.

### GetRecommendedGasRate

`func (o *QuoteSwapResponse) GetRecommendedGasRate() string`

GetRecommendedGasRate returns the RecommendedGasRate field if non-nil, zero value otherwise.

### GetRecommendedGasRateOk

`func (o *QuoteSwapResponse) GetRecommendedGasRateOk() (*string, bool)`

GetRecommendedGasRateOk returns a tuple with the RecommendedGasRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedGasRate

`func (o *QuoteSwapResponse) SetRecommendedGasRate(v string)`

SetRecommendedGasRate sets RecommendedGasRate field to given value.

### HasRecommendedGasRate

`func (o *QuoteSwapResponse) HasRecommendedGasRate() bool`

HasRecommendedGasRate returns a boolean if a field has been set.

### GetGasRateUnits

`func (o *QuoteSwapResponse) GetGasRateUnits() string`

GetGasRateUnits returns the GasRateUnits field if non-nil, zero value otherwise.

### GetGasRateUnitsOk

`func (o *QuoteSwapResponse) GetGasRateUnitsOk() (*string, bool)`

GetGasRateUnitsOk returns a tuple with the GasRateUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasRateUnits

`func (o *QuoteSwapResponse) SetGasRateUnits(v string)`

SetGasRateUnits sets GasRateUnits field to given value.

### HasGasRateUnits

`func (o *QuoteSwapResponse) HasGasRateUnits() bool`

HasGasRateUnits returns a boolean if a field has been set.

### GetMemo

`func (o *QuoteSwapResponse) GetMemo() string`

GetMemo returns the Memo field if non-nil, zero value otherwise.

### GetMemoOk

`func (o *QuoteSwapResponse) GetMemoOk() (*string, bool)`

GetMemoOk returns a tuple with the Memo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemo

`func (o *QuoteSwapResponse) SetMemo(v string)`

SetMemo sets Memo field to given value.

### HasMemo

`func (o *QuoteSwapResponse) HasMemo() bool`

HasMemo returns a boolean if a field has been set.

### GetExpectedAmountOut

`func (o *QuoteSwapResponse) GetExpectedAmountOut() string`

GetExpectedAmountOut returns the ExpectedAmountOut field if non-nil, zero value otherwise.

### GetExpectedAmountOutOk

`func (o *QuoteSwapResponse) GetExpectedAmountOutOk() (*string, bool)`

GetExpectedAmountOutOk returns a tuple with the ExpectedAmountOut field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedAmountOut

`func (o *QuoteSwapResponse) SetExpectedAmountOut(v string)`

SetExpectedAmountOut sets ExpectedAmountOut field to given value.


### GetExpectedAmountOutStreaming

`func (o *QuoteSwapResponse) GetExpectedAmountOutStreaming() string`

GetExpectedAmountOutStreaming returns the ExpectedAmountOutStreaming field if non-nil, zero value otherwise.

### GetExpectedAmountOutStreamingOk

`func (o *QuoteSwapResponse) GetExpectedAmountOutStreamingOk() (*string, bool)`

GetExpectedAmountOutStreamingOk returns a tuple with the ExpectedAmountOutStreaming field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedAmountOutStreaming

`func (o *QuoteSwapResponse) SetExpectedAmountOutStreaming(v string)`

SetExpectedAmountOutStreaming sets ExpectedAmountOutStreaming field to given value.


### GetMaxStreamingQuantity

`func (o *QuoteSwapResponse) GetMaxStreamingQuantity() int64`

GetMaxStreamingQuantity returns the MaxStreamingQuantity field if non-nil, zero value otherwise.

### GetMaxStreamingQuantityOk

`func (o *QuoteSwapResponse) GetMaxStreamingQuantityOk() (*int64, bool)`

GetMaxStreamingQuantityOk returns a tuple with the MaxStreamingQuantity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxStreamingQuantity

`func (o *QuoteSwapResponse) SetMaxStreamingQuantity(v int64)`

SetMaxStreamingQuantity sets MaxStreamingQuantity field to given value.

### HasMaxStreamingQuantity

`func (o *QuoteSwapResponse) HasMaxStreamingQuantity() bool`

HasMaxStreamingQuantity returns a boolean if a field has been set.

### GetStreamingSwapBlocks

`func (o *QuoteSwapResponse) GetStreamingSwapBlocks() int64`

GetStreamingSwapBlocks returns the StreamingSwapBlocks field if non-nil, zero value otherwise.

### GetStreamingSwapBlocksOk

`func (o *QuoteSwapResponse) GetStreamingSwapBlocksOk() (*int64, bool)`

GetStreamingSwapBlocksOk returns a tuple with the StreamingSwapBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSwapBlocks

`func (o *QuoteSwapResponse) SetStreamingSwapBlocks(v int64)`

SetStreamingSwapBlocks sets StreamingSwapBlocks field to given value.

### HasStreamingSwapBlocks

`func (o *QuoteSwapResponse) HasStreamingSwapBlocks() bool`

HasStreamingSwapBlocks returns a boolean if a field has been set.

### GetStreamingSwapSeconds

`func (o *QuoteSwapResponse) GetStreamingSwapSeconds() int64`

GetStreamingSwapSeconds returns the StreamingSwapSeconds field if non-nil, zero value otherwise.

### GetStreamingSwapSecondsOk

`func (o *QuoteSwapResponse) GetStreamingSwapSecondsOk() (*int64, bool)`

GetStreamingSwapSecondsOk returns a tuple with the StreamingSwapSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSwapSeconds

`func (o *QuoteSwapResponse) SetStreamingSwapSeconds(v int64)`

SetStreamingSwapSeconds sets StreamingSwapSeconds field to given value.

### HasStreamingSwapSeconds

`func (o *QuoteSwapResponse) HasStreamingSwapSeconds() bool`

HasStreamingSwapSeconds returns a boolean if a field has been set.

### GetTotalSwapSeconds

`func (o *QuoteSwapResponse) GetTotalSwapSeconds() int64`

GetTotalSwapSeconds returns the TotalSwapSeconds field if non-nil, zero value otherwise.

### GetTotalSwapSecondsOk

`func (o *QuoteSwapResponse) GetTotalSwapSecondsOk() (*int64, bool)`

GetTotalSwapSecondsOk returns a tuple with the TotalSwapSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalSwapSeconds

`func (o *QuoteSwapResponse) SetTotalSwapSeconds(v int64)`

SetTotalSwapSeconds sets TotalSwapSeconds field to given value.

### HasTotalSwapSeconds

`func (o *QuoteSwapResponse) HasTotalSwapSeconds() bool`

HasTotalSwapSeconds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


