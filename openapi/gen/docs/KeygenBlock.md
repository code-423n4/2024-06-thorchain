# KeygenBlock

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Height** | Pointer to **int64** | the height of the keygen block | [optional] 
**Keygens** | [**[]Keygen**](Keygen.md) |  | 

## Methods

### NewKeygenBlock

`func NewKeygenBlock(keygens []Keygen, ) *KeygenBlock`

NewKeygenBlock instantiates a new KeygenBlock object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeygenBlockWithDefaults

`func NewKeygenBlockWithDefaults() *KeygenBlock`

NewKeygenBlockWithDefaults instantiates a new KeygenBlock object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHeight

`func (o *KeygenBlock) GetHeight() int64`

GetHeight returns the Height field if non-nil, zero value otherwise.

### GetHeightOk

`func (o *KeygenBlock) GetHeightOk() (*int64, bool)`

GetHeightOk returns a tuple with the Height field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeight

`func (o *KeygenBlock) SetHeight(v int64)`

SetHeight sets Height field to given value.

### HasHeight

`func (o *KeygenBlock) HasHeight() bool`

HasHeight returns a boolean if a field has been set.

### GetKeygens

`func (o *KeygenBlock) GetKeygens() []Keygen`

GetKeygens returns the Keygens field if non-nil, zero value otherwise.

### GetKeygensOk

`func (o *KeygenBlock) GetKeygensOk() (*[]Keygen, bool)`

GetKeygensOk returns a tuple with the Keygens field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeygens

`func (o *KeygenBlock) SetKeygens(v []Keygen)`

SetKeygens sets Keygens field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


