# BanResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NodeAddress** | Pointer to **string** |  | [optional] 
**BlockHeight** | Pointer to **int64** |  | [optional] 
**Signers** | Pointer to **[]string** |  | [optional] 

## Methods

### NewBanResponse

`func NewBanResponse() *BanResponse`

NewBanResponse instantiates a new BanResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBanResponseWithDefaults

`func NewBanResponseWithDefaults() *BanResponse`

NewBanResponseWithDefaults instantiates a new BanResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNodeAddress

`func (o *BanResponse) GetNodeAddress() string`

GetNodeAddress returns the NodeAddress field if non-nil, zero value otherwise.

### GetNodeAddressOk

`func (o *BanResponse) GetNodeAddressOk() (*string, bool)`

GetNodeAddressOk returns a tuple with the NodeAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeAddress

`func (o *BanResponse) SetNodeAddress(v string)`

SetNodeAddress sets NodeAddress field to given value.

### HasNodeAddress

`func (o *BanResponse) HasNodeAddress() bool`

HasNodeAddress returns a boolean if a field has been set.

### GetBlockHeight

`func (o *BanResponse) GetBlockHeight() int64`

GetBlockHeight returns the BlockHeight field if non-nil, zero value otherwise.

### GetBlockHeightOk

`func (o *BanResponse) GetBlockHeightOk() (*int64, bool)`

GetBlockHeightOk returns a tuple with the BlockHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockHeight

`func (o *BanResponse) SetBlockHeight(v int64)`

SetBlockHeight sets BlockHeight field to given value.

### HasBlockHeight

`func (o *BanResponse) HasBlockHeight() bool`

HasBlockHeight returns a boolean if a field has been set.

### GetSigners

`func (o *BanResponse) GetSigners() []string`

GetSigners returns the Signers field if non-nil, zero value otherwise.

### GetSignersOk

`func (o *BanResponse) GetSignersOk() (*[]string, bool)`

GetSignersOk returns a tuple with the Signers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSigners

`func (o *BanResponse) SetSigners(v []string)`

SetSigners sets Signers field to given value.

### HasSigners

`func (o *BanResponse) HasSigners() bool`

HasSigners returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


