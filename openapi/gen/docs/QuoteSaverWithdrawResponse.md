# QuoteSaverWithdrawResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InboundAddress** | **string** | the inbound address for the transaction on the source chain | 
**InboundConfirmationBlocks** | Pointer to **int64** | the approximate number of source chain blocks required before processing | [optional] 
**InboundConfirmationSeconds** | Pointer to **int64** | the approximate seconds for block confirmations required before processing | [optional] 
**OutboundDelayBlocks** | **int64** | the number of thorchain blocks the outbound will be delayed | 
**OutboundDelaySeconds** | **int64** | the approximate seconds for the outbound delay before it will be sent | 
**Fees** | [**QuoteFees**](QuoteFees.md) |  | 
**SlippageBps** | **int64** | Deprecated - migrate to fees object. | 
**StreamingSlippageBps** | Pointer to **int64** | Deprecated - migrate to fees object. | [optional] 
**Router** | Pointer to **string** | the EVM chain router contract address | [optional] 
**Expiry** | **int64** | expiration timestamp in unix seconds | 
**Warning** | **string** | static warning message | 
**Notes** | **string** | chain specific quote notes | 
**DustThreshold** | Pointer to **string** | Defines the minimum transaction size for the chain in base units (sats, wei, uatom). Transactions with asset amounts lower than the dust_threshold are ignored. | [optional] 
**RecommendedMinAmountIn** | Pointer to **string** | The recommended minimum inbound amount for this transaction type &amp; inbound asset. Sending less than this amount could result in failed refunds. | [optional] 
**RecommendedGasRate** | **string** | the recommended gas rate to use for the inbound to ensure timely confirmation | 
**GasRateUnits** | **string** | the units of the recommended gas rate | 
**Memo** | **string** | generated memo for the withdraw, the client can use this OR send the dust amount | 
**DustAmount** | **string** | the dust amount of the target asset the user should send to initialize the withdraw, the client can send this OR provide the memo | 
**ExpectedAmountOut** | **string** | the amount of the target asset the user can expect to withdraw after fees in 1e8 decimals | 

## Methods

### NewQuoteSaverWithdrawResponse

`func NewQuoteSaverWithdrawResponse(inboundAddress string, outboundDelayBlocks int64, outboundDelaySeconds int64, fees QuoteFees, slippageBps int64, expiry int64, warning string, notes string, recommendedGasRate string, gasRateUnits string, memo string, dustAmount string, expectedAmountOut string, ) *QuoteSaverWithdrawResponse`

NewQuoteSaverWithdrawResponse instantiates a new QuoteSaverWithdrawResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQuoteSaverWithdrawResponseWithDefaults

`func NewQuoteSaverWithdrawResponseWithDefaults() *QuoteSaverWithdrawResponse`

NewQuoteSaverWithdrawResponseWithDefaults instantiates a new QuoteSaverWithdrawResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInboundAddress

`func (o *QuoteSaverWithdrawResponse) GetInboundAddress() string`

GetInboundAddress returns the InboundAddress field if non-nil, zero value otherwise.

### GetInboundAddressOk

`func (o *QuoteSaverWithdrawResponse) GetInboundAddressOk() (*string, bool)`

GetInboundAddressOk returns a tuple with the InboundAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundAddress

`func (o *QuoteSaverWithdrawResponse) SetInboundAddress(v string)`

SetInboundAddress sets InboundAddress field to given value.


### GetInboundConfirmationBlocks

`func (o *QuoteSaverWithdrawResponse) GetInboundConfirmationBlocks() int64`

GetInboundConfirmationBlocks returns the InboundConfirmationBlocks field if non-nil, zero value otherwise.

### GetInboundConfirmationBlocksOk

`func (o *QuoteSaverWithdrawResponse) GetInboundConfirmationBlocksOk() (*int64, bool)`

GetInboundConfirmationBlocksOk returns a tuple with the InboundConfirmationBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationBlocks

`func (o *QuoteSaverWithdrawResponse) SetInboundConfirmationBlocks(v int64)`

SetInboundConfirmationBlocks sets InboundConfirmationBlocks field to given value.

### HasInboundConfirmationBlocks

`func (o *QuoteSaverWithdrawResponse) HasInboundConfirmationBlocks() bool`

HasInboundConfirmationBlocks returns a boolean if a field has been set.

### GetInboundConfirmationSeconds

`func (o *QuoteSaverWithdrawResponse) GetInboundConfirmationSeconds() int64`

GetInboundConfirmationSeconds returns the InboundConfirmationSeconds field if non-nil, zero value otherwise.

### GetInboundConfirmationSecondsOk

`func (o *QuoteSaverWithdrawResponse) GetInboundConfirmationSecondsOk() (*int64, bool)`

GetInboundConfirmationSecondsOk returns a tuple with the InboundConfirmationSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationSeconds

`func (o *QuoteSaverWithdrawResponse) SetInboundConfirmationSeconds(v int64)`

SetInboundConfirmationSeconds sets InboundConfirmationSeconds field to given value.

### HasInboundConfirmationSeconds

`func (o *QuoteSaverWithdrawResponse) HasInboundConfirmationSeconds() bool`

HasInboundConfirmationSeconds returns a boolean if a field has been set.

### GetOutboundDelayBlocks

`func (o *QuoteSaverWithdrawResponse) GetOutboundDelayBlocks() int64`

GetOutboundDelayBlocks returns the OutboundDelayBlocks field if non-nil, zero value otherwise.

### GetOutboundDelayBlocksOk

`func (o *QuoteSaverWithdrawResponse) GetOutboundDelayBlocksOk() (*int64, bool)`

GetOutboundDelayBlocksOk returns a tuple with the OutboundDelayBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelayBlocks

`func (o *QuoteSaverWithdrawResponse) SetOutboundDelayBlocks(v int64)`

SetOutboundDelayBlocks sets OutboundDelayBlocks field to given value.


### GetOutboundDelaySeconds

`func (o *QuoteSaverWithdrawResponse) GetOutboundDelaySeconds() int64`

GetOutboundDelaySeconds returns the OutboundDelaySeconds field if non-nil, zero value otherwise.

### GetOutboundDelaySecondsOk

`func (o *QuoteSaverWithdrawResponse) GetOutboundDelaySecondsOk() (*int64, bool)`

GetOutboundDelaySecondsOk returns a tuple with the OutboundDelaySeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelaySeconds

`func (o *QuoteSaverWithdrawResponse) SetOutboundDelaySeconds(v int64)`

SetOutboundDelaySeconds sets OutboundDelaySeconds field to given value.


### GetFees

`func (o *QuoteSaverWithdrawResponse) GetFees() QuoteFees`

GetFees returns the Fees field if non-nil, zero value otherwise.

### GetFeesOk

`func (o *QuoteSaverWithdrawResponse) GetFeesOk() (*QuoteFees, bool)`

GetFeesOk returns a tuple with the Fees field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFees

`func (o *QuoteSaverWithdrawResponse) SetFees(v QuoteFees)`

SetFees sets Fees field to given value.


### GetSlippageBps

`func (o *QuoteSaverWithdrawResponse) GetSlippageBps() int64`

GetSlippageBps returns the SlippageBps field if non-nil, zero value otherwise.

### GetSlippageBpsOk

`func (o *QuoteSaverWithdrawResponse) GetSlippageBpsOk() (*int64, bool)`

GetSlippageBpsOk returns a tuple with the SlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlippageBps

`func (o *QuoteSaverWithdrawResponse) SetSlippageBps(v int64)`

SetSlippageBps sets SlippageBps field to given value.


### GetStreamingSlippageBps

`func (o *QuoteSaverWithdrawResponse) GetStreamingSlippageBps() int64`

GetStreamingSlippageBps returns the StreamingSlippageBps field if non-nil, zero value otherwise.

### GetStreamingSlippageBpsOk

`func (o *QuoteSaverWithdrawResponse) GetStreamingSlippageBpsOk() (*int64, bool)`

GetStreamingSlippageBpsOk returns a tuple with the StreamingSlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSlippageBps

`func (o *QuoteSaverWithdrawResponse) SetStreamingSlippageBps(v int64)`

SetStreamingSlippageBps sets StreamingSlippageBps field to given value.

### HasStreamingSlippageBps

`func (o *QuoteSaverWithdrawResponse) HasStreamingSlippageBps() bool`

HasStreamingSlippageBps returns a boolean if a field has been set.

### GetRouter

`func (o *QuoteSaverWithdrawResponse) GetRouter() string`

GetRouter returns the Router field if non-nil, zero value otherwise.

### GetRouterOk

`func (o *QuoteSaverWithdrawResponse) GetRouterOk() (*string, bool)`

GetRouterOk returns a tuple with the Router field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouter

`func (o *QuoteSaverWithdrawResponse) SetRouter(v string)`

SetRouter sets Router field to given value.

### HasRouter

`func (o *QuoteSaverWithdrawResponse) HasRouter() bool`

HasRouter returns a boolean if a field has been set.

### GetExpiry

`func (o *QuoteSaverWithdrawResponse) GetExpiry() int64`

GetExpiry returns the Expiry field if non-nil, zero value otherwise.

### GetExpiryOk

`func (o *QuoteSaverWithdrawResponse) GetExpiryOk() (*int64, bool)`

GetExpiryOk returns a tuple with the Expiry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiry

`func (o *QuoteSaverWithdrawResponse) SetExpiry(v int64)`

SetExpiry sets Expiry field to given value.


### GetWarning

`func (o *QuoteSaverWithdrawResponse) GetWarning() string`

GetWarning returns the Warning field if non-nil, zero value otherwise.

### GetWarningOk

`func (o *QuoteSaverWithdrawResponse) GetWarningOk() (*string, bool)`

GetWarningOk returns a tuple with the Warning field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarning

`func (o *QuoteSaverWithdrawResponse) SetWarning(v string)`

SetWarning sets Warning field to given value.


### GetNotes

`func (o *QuoteSaverWithdrawResponse) GetNotes() string`

GetNotes returns the Notes field if non-nil, zero value otherwise.

### GetNotesOk

`func (o *QuoteSaverWithdrawResponse) GetNotesOk() (*string, bool)`

GetNotesOk returns a tuple with the Notes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotes

`func (o *QuoteSaverWithdrawResponse) SetNotes(v string)`

SetNotes sets Notes field to given value.


### GetDustThreshold

`func (o *QuoteSaverWithdrawResponse) GetDustThreshold() string`

GetDustThreshold returns the DustThreshold field if non-nil, zero value otherwise.

### GetDustThresholdOk

`func (o *QuoteSaverWithdrawResponse) GetDustThresholdOk() (*string, bool)`

GetDustThresholdOk returns a tuple with the DustThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDustThreshold

`func (o *QuoteSaverWithdrawResponse) SetDustThreshold(v string)`

SetDustThreshold sets DustThreshold field to given value.

### HasDustThreshold

`func (o *QuoteSaverWithdrawResponse) HasDustThreshold() bool`

HasDustThreshold returns a boolean if a field has been set.

### GetRecommendedMinAmountIn

`func (o *QuoteSaverWithdrawResponse) GetRecommendedMinAmountIn() string`

GetRecommendedMinAmountIn returns the RecommendedMinAmountIn field if non-nil, zero value otherwise.

### GetRecommendedMinAmountInOk

`func (o *QuoteSaverWithdrawResponse) GetRecommendedMinAmountInOk() (*string, bool)`

GetRecommendedMinAmountInOk returns a tuple with the RecommendedMinAmountIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedMinAmountIn

`func (o *QuoteSaverWithdrawResponse) SetRecommendedMinAmountIn(v string)`

SetRecommendedMinAmountIn sets RecommendedMinAmountIn field to given value.

### HasRecommendedMinAmountIn

`func (o *QuoteSaverWithdrawResponse) HasRecommendedMinAmountIn() bool`

HasRecommendedMinAmountIn returns a boolean if a field has been set.

### GetRecommendedGasRate

`func (o *QuoteSaverWithdrawResponse) GetRecommendedGasRate() string`

GetRecommendedGasRate returns the RecommendedGasRate field if non-nil, zero value otherwise.

### GetRecommendedGasRateOk

`func (o *QuoteSaverWithdrawResponse) GetRecommendedGasRateOk() (*string, bool)`

GetRecommendedGasRateOk returns a tuple with the RecommendedGasRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedGasRate

`func (o *QuoteSaverWithdrawResponse) SetRecommendedGasRate(v string)`

SetRecommendedGasRate sets RecommendedGasRate field to given value.


### GetGasRateUnits

`func (o *QuoteSaverWithdrawResponse) GetGasRateUnits() string`

GetGasRateUnits returns the GasRateUnits field if non-nil, zero value otherwise.

### GetGasRateUnitsOk

`func (o *QuoteSaverWithdrawResponse) GetGasRateUnitsOk() (*string, bool)`

GetGasRateUnitsOk returns a tuple with the GasRateUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasRateUnits

`func (o *QuoteSaverWithdrawResponse) SetGasRateUnits(v string)`

SetGasRateUnits sets GasRateUnits field to given value.


### GetMemo

`func (o *QuoteSaverWithdrawResponse) GetMemo() string`

GetMemo returns the Memo field if non-nil, zero value otherwise.

### GetMemoOk

`func (o *QuoteSaverWithdrawResponse) GetMemoOk() (*string, bool)`

GetMemoOk returns a tuple with the Memo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemo

`func (o *QuoteSaverWithdrawResponse) SetMemo(v string)`

SetMemo sets Memo field to given value.


### GetDustAmount

`func (o *QuoteSaverWithdrawResponse) GetDustAmount() string`

GetDustAmount returns the DustAmount field if non-nil, zero value otherwise.

### GetDustAmountOk

`func (o *QuoteSaverWithdrawResponse) GetDustAmountOk() (*string, bool)`

GetDustAmountOk returns a tuple with the DustAmount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDustAmount

`func (o *QuoteSaverWithdrawResponse) SetDustAmount(v string)`

SetDustAmount sets DustAmount field to given value.


### GetExpectedAmountOut

`func (o *QuoteSaverWithdrawResponse) GetExpectedAmountOut() string`

GetExpectedAmountOut returns the ExpectedAmountOut field if non-nil, zero value otherwise.

### GetExpectedAmountOutOk

`func (o *QuoteSaverWithdrawResponse) GetExpectedAmountOutOk() (*string, bool)`

GetExpectedAmountOutOk returns a tuple with the ExpectedAmountOut field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedAmountOut

`func (o *QuoteSaverWithdrawResponse) SetExpectedAmountOut(v string)`

SetExpectedAmountOut sets ExpectedAmountOut field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


