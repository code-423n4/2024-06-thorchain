# Vault

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BlockHeight** | Pointer to **int64** |  | [optional] 
**PubKey** | Pointer to **string** |  | [optional] 
**Coins** | [**[]Coin**](Coin.md) |  | 
**Type** | Pointer to **string** |  | [optional] 
**Status** | **string** |  | 
**StatusSince** | Pointer to **int64** |  | [optional] 
**Membership** | Pointer to **[]string** | the list of node public keys which are members of the vault | [optional] 
**Chains** | Pointer to **[]string** |  | [optional] 
**InboundTxCount** | Pointer to **int64** |  | [optional] 
**OutboundTxCount** | Pointer to **int64** |  | [optional] 
**PendingTxBlockHeights** | Pointer to **[]int64** |  | [optional] 
**Routers** | [**[]VaultRouter**](VaultRouter.md) |  | 
**Addresses** | [**[]VaultAddress**](VaultAddress.md) |  | 
**Frozen** | Pointer to **[]string** |  | [optional] 

## Methods

### NewVault

`func NewVault(coins []Coin, status string, routers []VaultRouter, addresses []VaultAddress, ) *Vault`

NewVault instantiates a new Vault object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVaultWithDefaults

`func NewVaultWithDefaults() *Vault`

NewVaultWithDefaults instantiates a new Vault object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBlockHeight

`func (o *Vault) GetBlockHeight() int64`

GetBlockHeight returns the BlockHeight field if non-nil, zero value otherwise.

### GetBlockHeightOk

`func (o *Vault) GetBlockHeightOk() (*int64, bool)`

GetBlockHeightOk returns a tuple with the BlockHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockHeight

`func (o *Vault) SetBlockHeight(v int64)`

SetBlockHeight sets BlockHeight field to given value.

### HasBlockHeight

`func (o *Vault) HasBlockHeight() bool`

HasBlockHeight returns a boolean if a field has been set.

### GetPubKey

`func (o *Vault) GetPubKey() string`

GetPubKey returns the PubKey field if non-nil, zero value otherwise.

### GetPubKeyOk

`func (o *Vault) GetPubKeyOk() (*string, bool)`

GetPubKeyOk returns a tuple with the PubKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPubKey

`func (o *Vault) SetPubKey(v string)`

SetPubKey sets PubKey field to given value.

### HasPubKey

`func (o *Vault) HasPubKey() bool`

HasPubKey returns a boolean if a field has been set.

### GetCoins

`func (o *Vault) GetCoins() []Coin`

GetCoins returns the Coins field if non-nil, zero value otherwise.

### GetCoinsOk

`func (o *Vault) GetCoinsOk() (*[]Coin, bool)`

GetCoinsOk returns a tuple with the Coins field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCoins

`func (o *Vault) SetCoins(v []Coin)`

SetCoins sets Coins field to given value.


### GetType

`func (o *Vault) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *Vault) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *Vault) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *Vault) HasType() bool`

HasType returns a boolean if a field has been set.

### GetStatus

`func (o *Vault) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Vault) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Vault) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetStatusSince

`func (o *Vault) GetStatusSince() int64`

GetStatusSince returns the StatusSince field if non-nil, zero value otherwise.

### GetStatusSinceOk

`func (o *Vault) GetStatusSinceOk() (*int64, bool)`

GetStatusSinceOk returns a tuple with the StatusSince field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusSince

`func (o *Vault) SetStatusSince(v int64)`

SetStatusSince sets StatusSince field to given value.

### HasStatusSince

`func (o *Vault) HasStatusSince() bool`

HasStatusSince returns a boolean if a field has been set.

### GetMembership

`func (o *Vault) GetMembership() []string`

GetMembership returns the Membership field if non-nil, zero value otherwise.

### GetMembershipOk

`func (o *Vault) GetMembershipOk() (*[]string, bool)`

GetMembershipOk returns a tuple with the Membership field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMembership

`func (o *Vault) SetMembership(v []string)`

SetMembership sets Membership field to given value.

### HasMembership

`func (o *Vault) HasMembership() bool`

HasMembership returns a boolean if a field has been set.

### GetChains

`func (o *Vault) GetChains() []string`

GetChains returns the Chains field if non-nil, zero value otherwise.

### GetChainsOk

`func (o *Vault) GetChainsOk() (*[]string, bool)`

GetChainsOk returns a tuple with the Chains field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChains

`func (o *Vault) SetChains(v []string)`

SetChains sets Chains field to given value.

### HasChains

`func (o *Vault) HasChains() bool`

HasChains returns a boolean if a field has been set.

### GetInboundTxCount

`func (o *Vault) GetInboundTxCount() int64`

GetInboundTxCount returns the InboundTxCount field if non-nil, zero value otherwise.

### GetInboundTxCountOk

`func (o *Vault) GetInboundTxCountOk() (*int64, bool)`

GetInboundTxCountOk returns a tuple with the InboundTxCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundTxCount

`func (o *Vault) SetInboundTxCount(v int64)`

SetInboundTxCount sets InboundTxCount field to given value.

### HasInboundTxCount

`func (o *Vault) HasInboundTxCount() bool`

HasInboundTxCount returns a boolean if a field has been set.

### GetOutboundTxCount

`func (o *Vault) GetOutboundTxCount() int64`

GetOutboundTxCount returns the OutboundTxCount field if non-nil, zero value otherwise.

### GetOutboundTxCountOk

`func (o *Vault) GetOutboundTxCountOk() (*int64, bool)`

GetOutboundTxCountOk returns a tuple with the OutboundTxCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundTxCount

`func (o *Vault) SetOutboundTxCount(v int64)`

SetOutboundTxCount sets OutboundTxCount field to given value.

### HasOutboundTxCount

`func (o *Vault) HasOutboundTxCount() bool`

HasOutboundTxCount returns a boolean if a field has been set.

### GetPendingTxBlockHeights

`func (o *Vault) GetPendingTxBlockHeights() []int64`

GetPendingTxBlockHeights returns the PendingTxBlockHeights field if non-nil, zero value otherwise.

### GetPendingTxBlockHeightsOk

`func (o *Vault) GetPendingTxBlockHeightsOk() (*[]int64, bool)`

GetPendingTxBlockHeightsOk returns a tuple with the PendingTxBlockHeights field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPendingTxBlockHeights

`func (o *Vault) SetPendingTxBlockHeights(v []int64)`

SetPendingTxBlockHeights sets PendingTxBlockHeights field to given value.

### HasPendingTxBlockHeights

`func (o *Vault) HasPendingTxBlockHeights() bool`

HasPendingTxBlockHeights returns a boolean if a field has been set.

### GetRouters

`func (o *Vault) GetRouters() []VaultRouter`

GetRouters returns the Routers field if non-nil, zero value otherwise.

### GetRoutersOk

`func (o *Vault) GetRoutersOk() (*[]VaultRouter, bool)`

GetRoutersOk returns a tuple with the Routers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouters

`func (o *Vault) SetRouters(v []VaultRouter)`

SetRouters sets Routers field to given value.


### GetAddresses

`func (o *Vault) GetAddresses() []VaultAddress`

GetAddresses returns the Addresses field if non-nil, zero value otherwise.

### GetAddressesOk

`func (o *Vault) GetAddressesOk() (*[]VaultAddress, bool)`

GetAddressesOk returns a tuple with the Addresses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddresses

`func (o *Vault) SetAddresses(v []VaultAddress)`

SetAddresses sets Addresses field to given value.


### GetFrozen

`func (o *Vault) GetFrozen() []string`

GetFrozen returns the Frozen field if non-nil, zero value otherwise.

### GetFrozenOk

`func (o *Vault) GetFrozenOk() (*[]string, bool)`

GetFrozenOk returns a tuple with the Frozen field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrozen

`func (o *Vault) SetFrozen(v []string)`

SetFrozen sets Frozen field to given value.

### HasFrozen

`func (o *Vault) HasFrozen() bool`

HasFrozen returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


