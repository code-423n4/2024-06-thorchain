# KeygenResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeygenBlock** | [**KeygenBlock**](KeygenBlock.md) |  | 
**Signature** | **string** |  | 

## Methods

### NewKeygenResponse

`func NewKeygenResponse(keygenBlock KeygenBlock, signature string, ) *KeygenResponse`

NewKeygenResponse instantiates a new KeygenResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeygenResponseWithDefaults

`func NewKeygenResponseWithDefaults() *KeygenResponse`

NewKeygenResponseWithDefaults instantiates a new KeygenResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeygenBlock

`func (o *KeygenResponse) GetKeygenBlock() KeygenBlock`

GetKeygenBlock returns the KeygenBlock field if non-nil, zero value otherwise.

### GetKeygenBlockOk

`func (o *KeygenResponse) GetKeygenBlockOk() (*KeygenBlock, bool)`

GetKeygenBlockOk returns a tuple with the KeygenBlock field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeygenBlock

`func (o *KeygenResponse) SetKeygenBlock(v KeygenBlock)`

SetKeygenBlock sets KeygenBlock field to given value.


### GetSignature

`func (o *KeygenResponse) GetSignature() string`

GetSignature returns the Signature field if non-nil, zero value otherwise.

### GetSignatureOk

`func (o *KeygenResponse) GetSignatureOk() (*string, bool)`

GetSignatureOk returns a tuple with the Signature field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignature

`func (o *KeygenResponse) SetSignature(v string)`

SetSignature sets Signature field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


