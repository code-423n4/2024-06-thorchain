# YggdrasilVault

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BlockHeight** | Pointer to **int64** |  | [optional] 
**PubKey** | Pointer to **string** |  | [optional] 
**Coins** | [**[]Coin**](Coin.md) |  | 
**Type** | Pointer to **string** |  | [optional] 
**StatusSince** | Pointer to **int64** |  | [optional] 
**Membership** | Pointer to **[]string** | the list of node public keys which are members of the vault | [optional] 
**Chains** | Pointer to **[]string** |  | [optional] 
**InboundTxCount** | Pointer to **int64** |  | [optional] 
**OutboundTxCount** | Pointer to **int64** |  | [optional] 
**PendingTxBlockHeights** | Pointer to **[]int64** |  | [optional] 
**Routers** | [**[]VaultRouter**](VaultRouter.md) |  | 
**Status** | **string** |  | 
**Bond** | **string** | current node bond | 
**TotalValue** | **string** | value in rune of the vault&#39;s assets | 
**Addresses** | [**[]VaultAddress**](VaultAddress.md) |  | 

## Methods

### NewYggdrasilVault

`func NewYggdrasilVault(coins []Coin, routers []VaultRouter, status string, bond string, totalValue string, addresses []VaultAddress, ) *YggdrasilVault`

NewYggdrasilVault instantiates a new YggdrasilVault object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewYggdrasilVaultWithDefaults

`func NewYggdrasilVaultWithDefaults() *YggdrasilVault`

NewYggdrasilVaultWithDefaults instantiates a new YggdrasilVault object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBlockHeight

`func (o *YggdrasilVault) GetBlockHeight() int64`

GetBlockHeight returns the BlockHeight field if non-nil, zero value otherwise.

### GetBlockHeightOk

`func (o *YggdrasilVault) GetBlockHeightOk() (*int64, bool)`

GetBlockHeightOk returns a tuple with the BlockHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockHeight

`func (o *YggdrasilVault) SetBlockHeight(v int64)`

SetBlockHeight sets BlockHeight field to given value.

### HasBlockHeight

`func (o *YggdrasilVault) HasBlockHeight() bool`

HasBlockHeight returns a boolean if a field has been set.

### GetPubKey

`func (o *YggdrasilVault) GetPubKey() string`

GetPubKey returns the PubKey field if non-nil, zero value otherwise.

### GetPubKeyOk

`func (o *YggdrasilVault) GetPubKeyOk() (*string, bool)`

GetPubKeyOk returns a tuple with the PubKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPubKey

`func (o *YggdrasilVault) SetPubKey(v string)`

SetPubKey sets PubKey field to given value.

### HasPubKey

`func (o *YggdrasilVault) HasPubKey() bool`

HasPubKey returns a boolean if a field has been set.

### GetCoins

`func (o *YggdrasilVault) GetCoins() []Coin`

GetCoins returns the Coins field if non-nil, zero value otherwise.

### GetCoinsOk

`func (o *YggdrasilVault) GetCoinsOk() (*[]Coin, bool)`

GetCoinsOk returns a tuple with the Coins field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCoins

`func (o *YggdrasilVault) SetCoins(v []Coin)`

SetCoins sets Coins field to given value.


### GetType

`func (o *YggdrasilVault) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *YggdrasilVault) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *YggdrasilVault) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *YggdrasilVault) HasType() bool`

HasType returns a boolean if a field has been set.

### GetStatusSince

`func (o *YggdrasilVault) GetStatusSince() int64`

GetStatusSince returns the StatusSince field if non-nil, zero value otherwise.

### GetStatusSinceOk

`func (o *YggdrasilVault) GetStatusSinceOk() (*int64, bool)`

GetStatusSinceOk returns a tuple with the StatusSince field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusSince

`func (o *YggdrasilVault) SetStatusSince(v int64)`

SetStatusSince sets StatusSince field to given value.

### HasStatusSince

`func (o *YggdrasilVault) HasStatusSince() bool`

HasStatusSince returns a boolean if a field has been set.

### GetMembership

`func (o *YggdrasilVault) GetMembership() []string`

GetMembership returns the Membership field if non-nil, zero value otherwise.

### GetMembershipOk

`func (o *YggdrasilVault) GetMembershipOk() (*[]string, bool)`

GetMembershipOk returns a tuple with the Membership field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMembership

`func (o *YggdrasilVault) SetMembership(v []string)`

SetMembership sets Membership field to given value.

### HasMembership

`func (o *YggdrasilVault) HasMembership() bool`

HasMembership returns a boolean if a field has been set.

### GetChains

`func (o *YggdrasilVault) GetChains() []string`

GetChains returns the Chains field if non-nil, zero value otherwise.

### GetChainsOk

`func (o *YggdrasilVault) GetChainsOk() (*[]string, bool)`

GetChainsOk returns a tuple with the Chains field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChains

`func (o *YggdrasilVault) SetChains(v []string)`

SetChains sets Chains field to given value.

### HasChains

`func (o *YggdrasilVault) HasChains() bool`

HasChains returns a boolean if a field has been set.

### GetInboundTxCount

`func (o *YggdrasilVault) GetInboundTxCount() int64`

GetInboundTxCount returns the InboundTxCount field if non-nil, zero value otherwise.

### GetInboundTxCountOk

`func (o *YggdrasilVault) GetInboundTxCountOk() (*int64, bool)`

GetInboundTxCountOk returns a tuple with the InboundTxCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInboundTxCount

`func (o *YggdrasilVault) SetInboundTxCount(v int64)`

SetInboundTxCount sets InboundTxCount field to given value.

### HasInboundTxCount

`func (o *YggdrasilVault) HasInboundTxCount() bool`

HasInboundTxCount returns a boolean if a field has been set.

### GetOutboundTxCount

`func (o *YggdrasilVault) GetOutboundTxCount() int64`

GetOutboundTxCount returns the OutboundTxCount field if non-nil, zero value otherwise.

### GetOutboundTxCountOk

`func (o *YggdrasilVault) GetOutboundTxCountOk() (*int64, bool)`

GetOutboundTxCountOk returns a tuple with the OutboundTxCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundTxCount

`func (o *YggdrasilVault) SetOutboundTxCount(v int64)`

SetOutboundTxCount sets OutboundTxCount field to given value.

### HasOutboundTxCount

`func (o *YggdrasilVault) HasOutboundTxCount() bool`

HasOutboundTxCount returns a boolean if a field has been set.

### GetPendingTxBlockHeights

`func (o *YggdrasilVault) GetPendingTxBlockHeights() []int64`

GetPendingTxBlockHeights returns the PendingTxBlockHeights field if non-nil, zero value otherwise.

### GetPendingTxBlockHeightsOk

`func (o *YggdrasilVault) GetPendingTxBlockHeightsOk() (*[]int64, bool)`

GetPendingTxBlockHeightsOk returns a tuple with the PendingTxBlockHeights field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPendingTxBlockHeights

`func (o *YggdrasilVault) SetPendingTxBlockHeights(v []int64)`

SetPendingTxBlockHeights sets PendingTxBlockHeights field to given value.

### HasPendingTxBlockHeights

`func (o *YggdrasilVault) HasPendingTxBlockHeights() bool`

HasPendingTxBlockHeights returns a boolean if a field has been set.

### GetRouters

`func (o *YggdrasilVault) GetRouters() []VaultRouter`

GetRouters returns the Routers field if non-nil, zero value otherwise.

### GetRoutersOk

`func (o *YggdrasilVault) GetRoutersOk() (*[]VaultRouter, bool)`

GetRoutersOk returns a tuple with the Routers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouters

`func (o *YggdrasilVault) SetRouters(v []VaultRouter)`

SetRouters sets Routers field to given value.


### GetStatus

`func (o *YggdrasilVault) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *YggdrasilVault) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *YggdrasilVault) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetBond

`func (o *YggdrasilVault) GetBond() string`

GetBond returns the Bond field if non-nil, zero value otherwise.

### GetBondOk

`func (o *YggdrasilVault) GetBondOk() (*string, bool)`

GetBondOk returns a tuple with the Bond field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBond

`func (o *YggdrasilVault) SetBond(v string)`

SetBond sets Bond field to given value.


### GetTotalValue

`func (o *YggdrasilVault) GetTotalValue() string`

GetTotalValue returns the TotalValue field if non-nil, zero value otherwise.

### GetTotalValueOk

`func (o *YggdrasilVault) GetTotalValueOk() (*string, bool)`

GetTotalValueOk returns a tuple with the TotalValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalValue

`func (o *YggdrasilVault) SetTotalValue(v string)`

SetTotalValue sets TotalValue field to given value.


### GetAddresses

`func (o *YggdrasilVault) GetAddresses() []VaultAddress`

GetAddresses returns the Addresses field if non-nil, zero value otherwise.

### GetAddressesOk

`func (o *YggdrasilVault) GetAddressesOk() (*[]VaultAddress, bool)`

GetAddressesOk returns a tuple with the Addresses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddresses

`func (o *YggdrasilVault) SetAddresses(v []VaultAddress)`

SetAddresses sets Addresses field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


