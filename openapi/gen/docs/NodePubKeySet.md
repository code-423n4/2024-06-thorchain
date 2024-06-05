# NodePubKeySet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Secp256k1** | Pointer to **string** |  | [optional] 
**Ed25519** | Pointer to **string** |  | [optional] 

## Methods

### NewNodePubKeySet

`func NewNodePubKeySet() *NodePubKeySet`

NewNodePubKeySet instantiates a new NodePubKeySet object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodePubKeySetWithDefaults

`func NewNodePubKeySetWithDefaults() *NodePubKeySet`

NewNodePubKeySetWithDefaults instantiates a new NodePubKeySet object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSecp256k1

`func (o *NodePubKeySet) GetSecp256k1() string`

GetSecp256k1 returns the Secp256k1 field if non-nil, zero value otherwise.

### GetSecp256k1Ok

`func (o *NodePubKeySet) GetSecp256k1Ok() (*string, bool)`

GetSecp256k1Ok returns a tuple with the Secp256k1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecp256k1

`func (o *NodePubKeySet) SetSecp256k1(v string)`

SetSecp256k1 sets Secp256k1 field to given value.

### HasSecp256k1

`func (o *NodePubKeySet) HasSecp256k1() bool`

HasSecp256k1 returns a boolean if a field has been set.

### GetEd25519

`func (o *NodePubKeySet) GetEd25519() string`

GetEd25519 returns the Ed25519 field if non-nil, zero value otherwise.

### GetEd25519Ok

`func (o *NodePubKeySet) GetEd25519Ok() (*string, bool)`

GetEd25519Ok returns a tuple with the Ed25519 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEd25519

`func (o *NodePubKeySet) SetEd25519(v string)`

SetEd25519 sets Ed25519 field to given value.

### HasEd25519

`func (o *NodePubKeySet) HasEd25519() bool`

HasEd25519 returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


