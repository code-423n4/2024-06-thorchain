# OutboundFee

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asset** | **string** | the asset to display the outbound fee for | 
**OutboundFee** | **string** | the asset&#39;s outbound fee, in (1e8-format) units of the asset | 
**FeeWithheldRune** | Pointer to **string** | Total RUNE the network has withheld as fees to later cover gas costs for this asset&#39;s outbounds | [optional] 
**FeeSpentRune** | Pointer to **string** | Total RUNE the network has spent to reimburse gas costs for this asset&#39;s outbounds | [optional] 
**SurplusRune** | Pointer to **string** | amount of RUNE by which the fee_withheld_rune exceeds the fee_spent_rune | [optional] 
**DynamicMultiplierBasisPoints** | Pointer to **string** | dynamic multiplier basis points, based on the surplus_rune, affecting the size of the outbound_fee | [optional] 

## Methods

### NewOutboundFee

`func NewOutboundFee(asset string, outboundFee string, ) *OutboundFee`

NewOutboundFee instantiates a new OutboundFee object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOutboundFeeWithDefaults

`func NewOutboundFeeWithDefaults() *OutboundFee`

NewOutboundFeeWithDefaults instantiates a new OutboundFee object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsset

`func (o *OutboundFee) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *OutboundFee) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *OutboundFee) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetOutboundFee

`func (o *OutboundFee) GetOutboundFee() string`

GetOutboundFee returns the OutboundFee field if non-nil, zero value otherwise.

### GetOutboundFeeOk

`func (o *OutboundFee) GetOutboundFeeOk() (*string, bool)`

GetOutboundFeeOk returns a tuple with the OutboundFee field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundFee

`func (o *OutboundFee) SetOutboundFee(v string)`

SetOutboundFee sets OutboundFee field to given value.


### GetFeeWithheldRune

`func (o *OutboundFee) GetFeeWithheldRune() string`

GetFeeWithheldRune returns the FeeWithheldRune field if non-nil, zero value otherwise.

### GetFeeWithheldRuneOk

`func (o *OutboundFee) GetFeeWithheldRuneOk() (*string, bool)`

GetFeeWithheldRuneOk returns a tuple with the FeeWithheldRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeeWithheldRune

`func (o *OutboundFee) SetFeeWithheldRune(v string)`

SetFeeWithheldRune sets FeeWithheldRune field to given value.

### HasFeeWithheldRune

`func (o *OutboundFee) HasFeeWithheldRune() bool`

HasFeeWithheldRune returns a boolean if a field has been set.

### GetFeeSpentRune

`func (o *OutboundFee) GetFeeSpentRune() string`

GetFeeSpentRune returns the FeeSpentRune field if non-nil, zero value otherwise.

### GetFeeSpentRuneOk

`func (o *OutboundFee) GetFeeSpentRuneOk() (*string, bool)`

GetFeeSpentRuneOk returns a tuple with the FeeSpentRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeeSpentRune

`func (o *OutboundFee) SetFeeSpentRune(v string)`

SetFeeSpentRune sets FeeSpentRune field to given value.

### HasFeeSpentRune

`func (o *OutboundFee) HasFeeSpentRune() bool`

HasFeeSpentRune returns a boolean if a field has been set.

### GetSurplusRune

`func (o *OutboundFee) GetSurplusRune() string`

GetSurplusRune returns the SurplusRune field if non-nil, zero value otherwise.

### GetSurplusRuneOk

`func (o *OutboundFee) GetSurplusRuneOk() (*string, bool)`

GetSurplusRuneOk returns a tuple with the SurplusRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSurplusRune

`func (o *OutboundFee) SetSurplusRune(v string)`

SetSurplusRune sets SurplusRune field to given value.

### HasSurplusRune

`func (o *OutboundFee) HasSurplusRune() bool`

HasSurplusRune returns a boolean if a field has been set.

### GetDynamicMultiplierBasisPoints

`func (o *OutboundFee) GetDynamicMultiplierBasisPoints() string`

GetDynamicMultiplierBasisPoints returns the DynamicMultiplierBasisPoints field if non-nil, zero value otherwise.

### GetDynamicMultiplierBasisPointsOk

`func (o *OutboundFee) GetDynamicMultiplierBasisPointsOk() (*string, bool)`

GetDynamicMultiplierBasisPointsOk returns a tuple with the DynamicMultiplierBasisPoints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDynamicMultiplierBasisPoints

`func (o *OutboundFee) SetDynamicMultiplierBasisPoints(v string)`

SetDynamicMultiplierBasisPoints sets DynamicMultiplierBasisPoints field to given value.

### HasDynamicMultiplierBasisPoints

`func (o *OutboundFee) HasDynamicMultiplierBasisPoints() bool`

HasDynamicMultiplierBasisPoints returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


