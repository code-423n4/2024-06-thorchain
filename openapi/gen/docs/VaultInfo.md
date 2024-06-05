# VaultInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PubKey** | **string** |  | 
**Routers** | [**[]VaultRouter**](VaultRouter.md) |  | 

## Methods

### NewVaultInfo

`func NewVaultInfo(pubKey string, routers []VaultRouter, ) *VaultInfo`

NewVaultInfo instantiates a new VaultInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVaultInfoWithDefaults

`func NewVaultInfoWithDefaults() *VaultInfo`

NewVaultInfoWithDefaults instantiates a new VaultInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPubKey

`func (o *VaultInfo) GetPubKey() string`

GetPubKey returns the PubKey field if non-nil, zero value otherwise.

### GetPubKeyOk

`func (o *VaultInfo) GetPubKeyOk() (*string, bool)`

GetPubKeyOk returns a tuple with the PubKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPubKey

`func (o *VaultInfo) SetPubKey(v string)`

SetPubKey sets PubKey field to given value.


### GetRouters

`func (o *VaultInfo) GetRouters() []VaultRouter`

GetRouters returns the Routers field if non-nil, zero value otherwise.

### GetRoutersOk

`func (o *VaultInfo) GetRoutersOk() (*[]VaultRouter, bool)`

GetRoutersOk returns a tuple with the Routers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouters

`func (o *VaultInfo) SetRouters(v []VaultRouter)`

SetRouters sets Routers field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


