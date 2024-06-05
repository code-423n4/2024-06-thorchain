# BaseQuoteResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InboundAddress** | Pointer to **string** | the inbound address for the transaction on the source chain | [optional] 
**InboundConfirmationBlocks** | Pointer to **int64** | the approximate number of source chain blocks required before processing | [optional] 
**InboundConfirmationSeconds** | Pointer to **int64** | the approximate seconds for block confirmations required before processing | [optional] 
**OutboundDelayBlocks** | Pointer to **int64** | the number of thorchain blocks the outbound will be delayed | [optional] 
**OutboundDelaySeconds** | Pointer to **int64** | the approximate seconds for the outbound delay before it will be sent | [optional] 
**Fees** | Pointer to [**QuoteFees**](QuoteFees.md) |  | [optional] 
**SlippageBps** | Pointer to **int64** | Deprecated - migrate to fees object. | [optional] 
**StreamingSlippageBps** | Pointer to **int64** | Deprecated - migrate to fees object. | [optional] 
**Router** | Pointer to **string** | the EVM chain router contract address | [optional] 
**Expiry** | Pointer to **int64** | expiration timestamp in unix seconds | [optional] 
**Warning** | Pointer to **string** | static warning message | [optional] 
**Notes** | Pointer to **string** | chain specific quote notes | [optional] 
**DustThreshold** | Pointer to **string** | Defines the minimum transaction size for the chain in base units (sats, wei, uatom). Transactions with asset amounts lower than the dust_threshold are ignored. | [optional] 
**RecommendedMinAmountIn** | Pointer to **string** | The recommended minimum inbound amount for this transaction type &amp; inbound asset. Sending less than this amount could result in failed refunds. | [optional] 
**RecommendedGasRate** | Pointer to **string** | the recommended gas rate to use for the inbound to ensure timely confirmation | [optional] 
**GasRateUnits** | Pointer to **string** | the units of the recommended gas rate | [optional] 

## Methods

### NewBaseQuoteResponse

`func NewBaseQuoteResponse() *BaseQuoteResponse`

NewBaseQuoteResponse instantiates a new BaseQuoteResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBaseQuoteResponseWithDefaults

`func NewBaseQuoteResponseWithDefaults() *BaseQuoteResponse`

NewBaseQuoteResponseWithDefaults instantiates a new BaseQuoteResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInboundAddress

`func (o *BaseQuoteResponse) GetInboundAddress() string`

GetInboundAddress returns the InboundAddress field if non-nil, zero value otherwise.

### GetInboundAddressOk

`func (o *BaseQuoteResponse) GetInboundAddressOk() (*string, bool)`

GetInboundAddressOk returns a tuple with the InboundAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundAddress

`func (o *BaseQuoteResponse) SetInboundAddress(v string)`

SetInboundAddress sets InboundAddress field to given value.

### HasInboundAddress

`func (o *BaseQuoteResponse) HasInboundAddress() bool`

HasInboundAddress returns a boolean if a field has been set.

### GetInboundConfirmationBlocks

`func (o *BaseQuoteResponse) GetInboundConfirmationBlocks() int64`

GetInboundConfirmationBlocks returns the InboundConfirmationBlocks field if non-nil, zero value otherwise.

### GetInboundConfirmationBlocksOk

`func (o *BaseQuoteResponse) GetInboundConfirmationBlocksOk() (*int64, bool)`

GetInboundConfirmationBlocksOk returns a tuple with the InboundConfirmationBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationBlocks

`func (o *BaseQuoteResponse) SetInboundConfirmationBlocks(v int64)`

SetInboundConfirmationBlocks sets InboundConfirmationBlocks field to given value.

### HasInboundConfirmationBlocks

`func (o *BaseQuoteResponse) HasInboundConfirmationBlocks() bool`

HasInboundConfirmationBlocks returns a boolean if a field has been set.

### GetInboundConfirmationSeconds

`func (o *BaseQuoteResponse) GetInboundConfirmationSeconds() int64`

GetInboundConfirmationSeconds returns the InboundConfirmationSeconds field if non-nil, zero value otherwise.

### GetInboundConfirmationSecondsOk

`func (o *BaseQuoteResponse) GetInboundConfirmationSecondsOk() (*int64, bool)`

GetInboundConfirmationSecondsOk returns a tuple with the InboundConfirmationSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundConfirmationSeconds

`func (o *BaseQuoteResponse) SetInboundConfirmationSeconds(v int64)`

SetInboundConfirmationSeconds sets InboundConfirmationSeconds field to given value.

### HasInboundConfirmationSeconds

`func (o *BaseQuoteResponse) HasInboundConfirmationSeconds() bool`

HasInboundConfirmationSeconds returns a boolean if a field has been set.

### GetOutboundDelayBlocks

`func (o *BaseQuoteResponse) GetOutboundDelayBlocks() int64`

GetOutboundDelayBlocks returns the OutboundDelayBlocks field if non-nil, zero value otherwise.

### GetOutboundDelayBlocksOk

`func (o *BaseQuoteResponse) GetOutboundDelayBlocksOk() (*int64, bool)`

GetOutboundDelayBlocksOk returns a tuple with the OutboundDelayBlocks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelayBlocks

`func (o *BaseQuoteResponse) SetOutboundDelayBlocks(v int64)`

SetOutboundDelayBlocks sets OutboundDelayBlocks field to given value.

### HasOutboundDelayBlocks

`func (o *BaseQuoteResponse) HasOutboundDelayBlocks() bool`

HasOutboundDelayBlocks returns a boolean if a field has been set.

### GetOutboundDelaySeconds

`func (o *BaseQuoteResponse) GetOutboundDelaySeconds() int64`

GetOutboundDelaySeconds returns the OutboundDelaySeconds field if non-nil, zero value otherwise.

### GetOutboundDelaySecondsOk

`func (o *BaseQuoteResponse) GetOutboundDelaySecondsOk() (*int64, bool)`

GetOutboundDelaySecondsOk returns a tuple with the OutboundDelaySeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundDelaySeconds

`func (o *BaseQuoteResponse) SetOutboundDelaySeconds(v int64)`

SetOutboundDelaySeconds sets OutboundDelaySeconds field to given value.

### HasOutboundDelaySeconds

`func (o *BaseQuoteResponse) HasOutboundDelaySeconds() bool`

HasOutboundDelaySeconds returns a boolean if a field has been set.

### GetFees

`func (o *BaseQuoteResponse) GetFees() QuoteFees`

GetFees returns the Fees field if non-nil, zero value otherwise.

### GetFeesOk

`func (o *BaseQuoteResponse) GetFeesOk() (*QuoteFees, bool)`

GetFeesOk returns a tuple with the Fees field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFees

`func (o *BaseQuoteResponse) SetFees(v QuoteFees)`

SetFees sets Fees field to given value.

### HasFees

`func (o *BaseQuoteResponse) HasFees() bool`

HasFees returns a boolean if a field has been set.

### GetSlippageBps

`func (o *BaseQuoteResponse) GetSlippageBps() int64`

GetSlippageBps returns the SlippageBps field if non-nil, zero value otherwise.

### GetSlippageBpsOk

`func (o *BaseQuoteResponse) GetSlippageBpsOk() (*int64, bool)`

GetSlippageBpsOk returns a tuple with the SlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlippageBps

`func (o *BaseQuoteResponse) SetSlippageBps(v int64)`

SetSlippageBps sets SlippageBps field to given value.

### HasSlippageBps

`func (o *BaseQuoteResponse) HasSlippageBps() bool`

HasSlippageBps returns a boolean if a field has been set.

### GetStreamingSlippageBps

`func (o *BaseQuoteResponse) GetStreamingSlippageBps() int64`

GetStreamingSlippageBps returns the StreamingSlippageBps field if non-nil, zero value otherwise.

### GetStreamingSlippageBpsOk

`func (o *BaseQuoteResponse) GetStreamingSlippageBpsOk() (*int64, bool)`

GetStreamingSlippageBpsOk returns a tuple with the StreamingSlippageBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamingSlippageBps

`func (o *BaseQuoteResponse) SetStreamingSlippageBps(v int64)`

SetStreamingSlippageBps sets StreamingSlippageBps field to given value.

### HasStreamingSlippageBps

`func (o *BaseQuoteResponse) HasStreamingSlippageBps() bool`

HasStreamingSlippageBps returns a boolean if a field has been set.

### GetRouter

`func (o *BaseQuoteResponse) GetRouter() string`

GetRouter returns the Router field if non-nil, zero value otherwise.

### GetRouterOk

`func (o *BaseQuoteResponse) GetRouterOk() (*string, bool)`

GetRouterOk returns a tuple with the Router field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouter

`func (o *BaseQuoteResponse) SetRouter(v string)`

SetRouter sets Router field to given value.

### HasRouter

`func (o *BaseQuoteResponse) HasRouter() bool`

HasRouter returns a boolean if a field has been set.

### GetExpiry

`func (o *BaseQuoteResponse) GetExpiry() int64`

GetExpiry returns the Expiry field if non-nil, zero value otherwise.

### GetExpiryOk

`func (o *BaseQuoteResponse) GetExpiryOk() (*int64, bool)`

GetExpiryOk returns a tuple with the Expiry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiry

`func (o *BaseQuoteResponse) SetExpiry(v int64)`

SetExpiry sets Expiry field to given value.

### HasExpiry

`func (o *BaseQuoteResponse) HasExpiry() bool`

HasExpiry returns a boolean if a field has been set.

### GetWarning

`func (o *BaseQuoteResponse) GetWarning() string`

GetWarning returns the Warning field if non-nil, zero value otherwise.

### GetWarningOk

`func (o *BaseQuoteResponse) GetWarningOk() (*string, bool)`

GetWarningOk returns a tuple with the Warning field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarning

`func (o *BaseQuoteResponse) SetWarning(v string)`

SetWarning sets Warning field to given value.

### HasWarning

`func (o *BaseQuoteResponse) HasWarning() bool`

HasWarning returns a boolean if a field has been set.

### GetNotes

`func (o *BaseQuoteResponse) GetNotes() string`

GetNotes returns the Notes field if non-nil, zero value otherwise.

### GetNotesOk

`func (o *BaseQuoteResponse) GetNotesOk() (*string, bool)`

GetNotesOk returns a tuple with the Notes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotes

`func (o *BaseQuoteResponse) SetNotes(v string)`

SetNotes sets Notes field to given value.

### HasNotes

`func (o *BaseQuoteResponse) HasNotes() bool`

HasNotes returns a boolean if a field has been set.

### GetDustThreshold

`func (o *BaseQuoteResponse) GetDustThreshold() string`

GetDustThreshold returns the DustThreshold field if non-nil, zero value otherwise.

### GetDustThresholdOk

`func (o *BaseQuoteResponse) GetDustThresholdOk() (*string, bool)`

GetDustThresholdOk returns a tuple with the DustThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDustThreshold

`func (o *BaseQuoteResponse) SetDustThreshold(v string)`

SetDustThreshold sets DustThreshold field to given value.

### HasDustThreshold

`func (o *BaseQuoteResponse) HasDustThreshold() bool`

HasDustThreshold returns a boolean if a field has been set.

### GetRecommendedMinAmountIn

`func (o *BaseQuoteResponse) GetRecommendedMinAmountIn() string`

GetRecommendedMinAmountIn returns the RecommendedMinAmountIn field if non-nil, zero value otherwise.

### GetRecommendedMinAmountInOk

`func (o *BaseQuoteResponse) GetRecommendedMinAmountInOk() (*string, bool)`

GetRecommendedMinAmountInOk returns a tuple with the RecommendedMinAmountIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedMinAmountIn

`func (o *BaseQuoteResponse) SetRecommendedMinAmountIn(v string)`

SetRecommendedMinAmountIn sets RecommendedMinAmountIn field to given value.

### HasRecommendedMinAmountIn

`func (o *BaseQuoteResponse) HasRecommendedMinAmountIn() bool`

HasRecommendedMinAmountIn returns a boolean if a field has been set.

### GetRecommendedGasRate

`func (o *BaseQuoteResponse) GetRecommendedGasRate() string`

GetRecommendedGasRate returns the RecommendedGasRate field if non-nil, zero value otherwise.

### GetRecommendedGasRateOk

`func (o *BaseQuoteResponse) GetRecommendedGasRateOk() (*string, bool)`

GetRecommendedGasRateOk returns a tuple with the RecommendedGasRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommendedGasRate

`func (o *BaseQuoteResponse) SetRecommendedGasRate(v string)`

SetRecommendedGasRate sets RecommendedGasRate field to given value.

### HasRecommendedGasRate

`func (o *BaseQuoteResponse) HasRecommendedGasRate() bool`

HasRecommendedGasRate returns a boolean if a field has been set.

### GetGasRateUnits

`func (o *BaseQuoteResponse) GetGasRateUnits() string`

GetGasRateUnits returns the GasRateUnits field if non-nil, zero value otherwise.

### GetGasRateUnitsOk

`func (o *BaseQuoteResponse) GetGasRateUnitsOk() (*string, bool)`

GetGasRateUnitsOk returns a tuple with the GasRateUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasRateUnits

`func (o *BaseQuoteResponse) SetGasRateUnits(v string)`

SetGasRateUnits sets GasRateUnits field to given value.

### HasGasRateUnits

`func (o *BaseQuoteResponse) HasGasRateUnits() bool`

HasGasRateUnits returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


