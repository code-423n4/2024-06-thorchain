# KeysignInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Height** | Pointer to **int64** | the block(s) in which a tx out item is scheduled to be signed and moved from the scheduled outbound queue to the outbound queue | [optional] 
**TxArray** | [**[]TxOutItem**](TxOutItem.md) |  | 

## Methods

### NewKeysignInfo

`func NewKeysignInfo(txArray []TxOutItem, ) *KeysignInfo`

NewKeysignInfo instantiates a new KeysignInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeysignInfoWithDefaults

`func NewKeysignInfoWithDefaults() *KeysignInfo`

NewKeysignInfoWithDefaults instantiates a new KeysignInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHeight

`func (o *KeysignInfo) GetHeight() int64`

GetHeight returns the Height field if non-nil, zero value otherwise.

### GetHeightOk

`func (o *KeysignInfo) GetHeightOk() (*int64, bool)`

GetHeightOk returns a tuple with the Height field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeight

`func (o *KeysignInfo) SetHeight(v int64)`

SetHeight sets Height field to given value.

### HasHeight

`func (o *KeysignInfo) HasHeight() bool`

HasHeight returns a boolean if a field has been set.

### GetTxArray

`func (o *KeysignInfo) GetTxArray() []TxOutItem`

GetTxArray returns the TxArray field if non-nil, zero value otherwise.

### GetTxArrayOk

`func (o *KeysignInfo) GetTxArrayOk() (*[]TxOutItem, bool)`

GetTxArrayOk returns a tuple with the TxArray field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTxArray

`func (o *KeysignInfo) SetTxArray(v []TxOutItem)`

SetTxArray sets TxArray field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


