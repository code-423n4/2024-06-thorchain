# VaultPubkeysResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Asgard** | [**[]VaultInfo**](VaultInfo.md) |  | 
**Yggdrasil** | [**[]VaultInfo**](VaultInfo.md) |  | 
**Inactive** | [**[]VaultInfo**](VaultInfo.md) |  | 

## Methods

### NewVaultPubkeysResponse

`func NewVaultPubkeysResponse(asgard []VaultInfo, yggdrasil []VaultInfo, inactive []VaultInfo, ) *VaultPubkeysResponse`

NewVaultPubkeysResponse instantiates a new VaultPubkeysResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVaultPubkeysResponseWithDefaults

`func NewVaultPubkeysResponseWithDefaults() *VaultPubkeysResponse`

NewVaultPubkeysResponseWithDefaults instantiates a new VaultPubkeysResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAsgard

`func (o *VaultPubkeysResponse) GetAsgard() []VaultInfo`

GetAsgard returns the Asgard field if non-nil, zero value otherwise.

### GetAsgardOk

`func (o *VaultPubkeysResponse) GetAsgardOk() (*[]VaultInfo, bool)`

GetAsgardOk returns a tuple with the Asgard field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsgard

`func (o *VaultPubkeysResponse) SetAsgard(v []VaultInfo)`

SetAsgard sets Asgard field to given value.


### GetYggdrasil

`func (o *VaultPubkeysResponse) GetYggdrasil() []VaultInfo`

GetYggdrasil returns the Yggdrasil field if non-nil, zero value otherwise.

### GetYggdrasilOk

`func (o *VaultPubkeysResponse) GetYggdrasilOk() (*[]VaultInfo, bool)`

GetYggdrasilOk returns a tuple with the Yggdrasil field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetYggdrasil

`func (o *VaultPubkeysResponse) SetYggdrasil(v []VaultInfo)`

SetYggdrasil sets Yggdrasil field to given value.


### GetInactive

`func (o *VaultPubkeysResponse) GetInactive() []VaultInfo`

GetInactive returns the Inactive field if non-nil, zero value otherwise.

### GetInactiveOk

`func (o *VaultPubkeysResponse) GetInactiveOk() (*[]VaultInfo, bool)`

GetInactiveOk returns a tuple with the Inactive field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInactive

`func (o *VaultPubkeysResponse) SetInactive(v []VaultInfo)`

SetInactive sets Inactive field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


