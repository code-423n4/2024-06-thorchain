# KeysignResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Keysign** | [**KeysignInfo**](KeysignInfo.md) |  | 
**Signature** | **string** |  | 

## Methods

### NewKeysignResponse

`func NewKeysignResponse(keysign KeysignInfo, signature string, ) *KeysignResponse`

NewKeysignResponse instantiates a new KeysignResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeysignResponseWithDefaults

`func NewKeysignResponseWithDefaults() *KeysignResponse`

NewKeysignResponseWithDefaults instantiates a new KeysignResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeysign

`func (o *KeysignResponse) GetKeysign() KeysignInfo`

GetKeysign returns the Keysign field if non-nil, zero value otherwise.

### GetKeysignOk

`func (o *KeysignResponse) GetKeysignOk() (*KeysignInfo, bool)`

GetKeysignOk returns a tuple with the Keysign field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeysign

`func (o *KeysignResponse) SetKeysign(v KeysignInfo)`

SetKeysign sets Keysign field to given value.


### GetSignature

`func (o *KeysignResponse) GetSignature() string`

GetSignature returns the Signature field if non-nil, zero value otherwise.

### GetSignatureOk

`func (o *KeysignResponse) GetSignatureOk() (*string, bool)`

GetSignatureOk returns a tuple with the Signature field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignature

`func (o *KeysignResponse) SetSignature(v string)`

SetSignature sets Signature field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


