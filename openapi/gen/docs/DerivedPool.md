# DerivedPool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asset** | **string** |  | 
**Status** | **string** |  | 
**Decimals** | Pointer to **int64** |  | [optional] 
**BalanceAsset** | **string** |  | 
**BalanceRune** | **string** |  | 
**DerivedDepthBps** | **string** | the depth of the derived virtual pool relative to L1 pool (in basis points) | 

## Methods

### NewDerivedPool

`func NewDerivedPool(asset string, status string, balanceAsset string, balanceRune string, derivedDepthBps string, ) *DerivedPool`

NewDerivedPool instantiates a new DerivedPool object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDerivedPoolWithDefaults

`func NewDerivedPoolWithDefaults() *DerivedPool`

NewDerivedPoolWithDefaults instantiates a new DerivedPool object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsset

`func (o *DerivedPool) GetAsset() string`

GetAsset returns the Asset field if non-nil, zero value otherwise.

### GetAssetOk

`func (o *DerivedPool) GetAssetOk() (*string, bool)`

GetAssetOk returns a tuple with the Asset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsset

`func (o *DerivedPool) SetAsset(v string)`

SetAsset sets Asset field to given value.


### GetStatus

`func (o *DerivedPool) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DerivedPool) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DerivedPool) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetDecimals

`func (o *DerivedPool) GetDecimals() int64`

GetDecimals returns the Decimals field if non-nil, zero value otherwise.

### GetDecimalsOk

`func (o *DerivedPool) GetDecimalsOk() (*int64, bool)`

GetDecimalsOk returns a tuple with the Decimals field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDecimals

`func (o *DerivedPool) SetDecimals(v int64)`

SetDecimals sets Decimals field to given value.

### HasDecimals

`func (o *DerivedPool) HasDecimals() bool`

HasDecimals returns a boolean if a field has been set.

### GetBalanceAsset

`func (o *DerivedPool) GetBalanceAsset() string`

GetBalanceAsset returns the BalanceAsset field if non-nil, zero value otherwise.

### GetBalanceAssetOk

`func (o *DerivedPool) GetBalanceAssetOk() (*string, bool)`

GetBalanceAssetOk returns a tuple with the BalanceAsset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalanceAsset

`func (o *DerivedPool) SetBalanceAsset(v string)`

SetBalanceAsset sets BalanceAsset field to given value.


### GetBalanceRune

`func (o *DerivedPool) GetBalanceRune() string`

GetBalanceRune returns the BalanceRune field if non-nil, zero value otherwise.

### GetBalanceRuneOk

`func (o *DerivedPool) GetBalanceRuneOk() (*string, bool)`

GetBalanceRuneOk returns a tuple with the BalanceRune field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalanceRune

`func (o *DerivedPool) SetBalanceRune(v string)`

SetBalanceRune sets BalanceRune field to given value.


### GetDerivedDepthBps

`func (o *DerivedPool) GetDerivedDepthBps() string`

GetDerivedDepthBps returns the DerivedDepthBps field if non-nil, zero value otherwise.

### GetDerivedDepthBpsOk

`func (o *DerivedPool) GetDerivedDepthBpsOk() (*string, bool)`

GetDerivedDepthBpsOk returns a tuple with the DerivedDepthBps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDerivedDepthBps

`func (o *DerivedPool) SetDerivedDepthBps(v string)`

SetDerivedDepthBps sets DerivedDepthBps field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


