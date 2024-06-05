# Node

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NodeAddress** | **string** |  | 
**Status** | **string** |  | 
**PubKeySet** | [**NodePubKeySet**](NodePubKeySet.md) |  | 
**ValidatorConsPubKey** | **string** | the consensus pub key for the node | 
**PeerId** | **string** | the P2PID (:6040/p2pid endpoint) of the node | 
**ActiveBlockHeight** | **int64** | the block height at which the node became active | 
**StatusSince** | **int64** | the block height of the current provided information for the node | 
**NodeOperatorAddress** | **string** |  | 
**TotalBond** | **string** | current node bond | 
**BondProviders** | [**NodeBondProviders**](NodeBondProviders.md) |  | 
**SignerMembership** | **[]string** | the set of vault public keys of which the node is a member | 
**RequestedToLeave** | **bool** |  | 
**ForcedToLeave** | **bool** | indicates whether the node has been forced to leave by the network, typically via ban | 
**LeaveHeight** | **int64** |  | 
**IpAddress** | **string** |  | 
**Version** | **string** | the currently set version of the node | 
**SlashPoints** | **int64** | the accumulated slash points, reset at churn but excessive slash points may carry over | 
**Jail** | [**NodeJail**](NodeJail.md) |  | 
**CurrentAward** | **string** |  | 
**ObserveChains** | [**[]ChainHeight**](ChainHeight.md) | the last observed heights for all chain by the node | 
**PreflightStatus** | [**NodePreflightStatus**](NodePreflightStatus.md) |  | 

## Methods

### NewNode

`func NewNode(nodeAddress string, status string, pubKeySet NodePubKeySet, validatorConsPubKey string, peerId string, activeBlockHeight int64, statusSince int64, nodeOperatorAddress string, totalBond string, bondProviders NodeBondProviders, signerMembership []string, requestedToLeave bool, forcedToLeave bool, leaveHeight int64, ipAddress string, version string, slashPoints int64, jail NodeJail, currentAward string, observeChains []ChainHeight, preflightStatus NodePreflightStatus, ) *Node`

NewNode instantiates a new Node object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeWithDefaults

`func NewNodeWithDefaults() *Node`

NewNodeWithDefaults instantiates a new Node object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNodeAddress

`func (o *Node) GetNodeAddress() string`

GetNodeAddress returns the NodeAddress field if non-nil, zero value otherwise.

### GetNodeAddressOk

`func (o *Node) GetNodeAddressOk() (*string, bool)`

GetNodeAddressOk returns a tuple with the NodeAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeAddress

`func (o *Node) SetNodeAddress(v string)`

SetNodeAddress sets NodeAddress field to given value.


### GetStatus

`func (o *Node) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Node) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Node) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetPubKeySet

`func (o *Node) GetPubKeySet() NodePubKeySet`

GetPubKeySet returns the PubKeySet field if non-nil, zero value otherwise.

### GetPubKeySetOk

`func (o *Node) GetPubKeySetOk() (*NodePubKeySet, bool)`

GetPubKeySetOk returns a tuple with the PubKeySet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPubKeySet

`func (o *Node) SetPubKeySet(v NodePubKeySet)`

SetPubKeySet sets PubKeySet field to given value.


### GetValidatorConsPubKey

`func (o *Node) GetValidatorConsPubKey() string`

GetValidatorConsPubKey returns the ValidatorConsPubKey field if non-nil, zero value otherwise.

### GetValidatorConsPubKeyOk

`func (o *Node) GetValidatorConsPubKeyOk() (*string, bool)`

GetValidatorConsPubKeyOk returns a tuple with the ValidatorConsPubKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValidatorConsPubKey

`func (o *Node) SetValidatorConsPubKey(v string)`

SetValidatorConsPubKey sets ValidatorConsPubKey field to given value.


### GetPeerId

`func (o *Node) GetPeerId() string`

GetPeerId returns the PeerId field if non-nil, zero value otherwise.

### GetPeerIdOk

`func (o *Node) GetPeerIdOk() (*string, bool)`

GetPeerIdOk returns a tuple with the PeerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPeerId

`func (o *Node) SetPeerId(v string)`

SetPeerId sets PeerId field to given value.


### GetActiveBlockHeight

`func (o *Node) GetActiveBlockHeight() int64`

GetActiveBlockHeight returns the ActiveBlockHeight field if non-nil, zero value otherwise.

### GetActiveBlockHeightOk

`func (o *Node) GetActiveBlockHeightOk() (*int64, bool)`

GetActiveBlockHeightOk returns a tuple with the ActiveBlockHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveBlockHeight

`func (o *Node) SetActiveBlockHeight(v int64)`

SetActiveBlockHeight sets ActiveBlockHeight field to given value.


### GetStatusSince

`func (o *Node) GetStatusSince() int64`

GetStatusSince returns the StatusSince field if non-nil, zero value otherwise.

### GetStatusSinceOk

`func (o *Node) GetStatusSinceOk() (*int64, bool)`

GetStatusSinceOk returns a tuple with the StatusSince field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusSince

`func (o *Node) SetStatusSince(v int64)`

SetStatusSince sets StatusSince field to given value.


### GetNodeOperatorAddress

`func (o *Node) GetNodeOperatorAddress() string`

GetNodeOperatorAddress returns the NodeOperatorAddress field if non-nil, zero value otherwise.

### GetNodeOperatorAddressOk

`func (o *Node) GetNodeOperatorAddressOk() (*string, bool)`

GetNodeOperatorAddressOk returns a tuple with the NodeOperatorAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeOperatorAddress

`func (o *Node) SetNodeOperatorAddress(v string)`

SetNodeOperatorAddress sets NodeOperatorAddress field to given value.


### GetTotalBond

`func (o *Node) GetTotalBond() string`

GetTotalBond returns the TotalBond field if non-nil, zero value otherwise.

### GetTotalBondOk

`func (o *Node) GetTotalBondOk() (*string, bool)`

GetTotalBondOk returns a tuple with the TotalBond field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalBond

`func (o *Node) SetTotalBond(v string)`

SetTotalBond sets TotalBond field to given value.


### GetBondProviders

`func (o *Node) GetBondProviders() NodeBondProviders`

GetBondProviders returns the BondProviders field if non-nil, zero value otherwise.

### GetBondProvidersOk

`func (o *Node) GetBondProvidersOk() (*NodeBondProviders, bool)`

GetBondProvidersOk returns a tuple with the BondProviders field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBondProviders

`func (o *Node) SetBondProviders(v NodeBondProviders)`

SetBondProviders sets BondProviders field to given value.


### GetSignerMembership

`func (o *Node) GetSignerMembership() []string`

GetSignerMembership returns the SignerMembership field if non-nil, zero value otherwise.

### GetSignerMembershipOk

`func (o *Node) GetSignerMembershipOk() (*[]string, bool)`

GetSignerMembershipOk returns a tuple with the SignerMembership field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignerMembership

`func (o *Node) SetSignerMembership(v []string)`

SetSignerMembership sets SignerMembership field to given value.


### GetRequestedToLeave

`func (o *Node) GetRequestedToLeave() bool`

GetRequestedToLeave returns the RequestedToLeave field if non-nil, zero value otherwise.

### GetRequestedToLeaveOk

`func (o *Node) GetRequestedToLeaveOk() (*bool, bool)`

GetRequestedToLeaveOk returns a tuple with the RequestedToLeave field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedToLeave

`func (o *Node) SetRequestedToLeave(v bool)`

SetRequestedToLeave sets RequestedToLeave field to given value.


### GetForcedToLeave

`func (o *Node) GetForcedToLeave() bool`

GetForcedToLeave returns the ForcedToLeave field if non-nil, zero value otherwise.

### GetForcedToLeaveOk

`func (o *Node) GetForcedToLeaveOk() (*bool, bool)`

GetForcedToLeaveOk returns a tuple with the ForcedToLeave field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetForcedToLeave

`func (o *Node) SetForcedToLeave(v bool)`

SetForcedToLeave sets ForcedToLeave field to given value.


### GetLeaveHeight

`func (o *Node) GetLeaveHeight() int64`

GetLeaveHeight returns the LeaveHeight field if non-nil, zero value otherwise.

### GetLeaveHeightOk

`func (o *Node) GetLeaveHeightOk() (*int64, bool)`

GetLeaveHeightOk returns a tuple with the LeaveHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeaveHeight

`func (o *Node) SetLeaveHeight(v int64)`

SetLeaveHeight sets LeaveHeight field to given value.


### GetIpAddress

`func (o *Node) GetIpAddress() string`

GetIpAddress returns the IpAddress field if non-nil, zero value otherwise.

### GetIpAddressOk

`func (o *Node) GetIpAddressOk() (*string, bool)`

GetIpAddressOk returns a tuple with the IpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpAddress

`func (o *Node) SetIpAddress(v string)`

SetIpAddress sets IpAddress field to given value.


### GetVersion

`func (o *Node) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *Node) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *Node) SetVersion(v string)`

SetVersion sets Version field to given value.


### GetSlashPoints

`func (o *Node) GetSlashPoints() int64`

GetSlashPoints returns the SlashPoints field if non-nil, zero value otherwise.

### GetSlashPointsOk

`func (o *Node) GetSlashPointsOk() (*int64, bool)`

GetSlashPointsOk returns a tuple with the SlashPoints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSlashPoints

`func (o *Node) SetSlashPoints(v int64)`

SetSlashPoints sets SlashPoints field to given value.


### GetJail

`func (o *Node) GetJail() NodeJail`

GetJail returns the Jail field if non-nil, zero value otherwise.

### GetJailOk

`func (o *Node) GetJailOk() (*NodeJail, bool)`

GetJailOk returns a tuple with the Jail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJail

`func (o *Node) SetJail(v NodeJail)`

SetJail sets Jail field to given value.


### GetCurrentAward

`func (o *Node) GetCurrentAward() string`

GetCurrentAward returns the CurrentAward field if non-nil, zero value otherwise.

### GetCurrentAwardOk

`func (o *Node) GetCurrentAwardOk() (*string, bool)`

GetCurrentAwardOk returns a tuple with the CurrentAward field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentAward

`func (o *Node) SetCurrentAward(v string)`

SetCurrentAward sets CurrentAward field to given value.


### GetObserveChains

`func (o *Node) GetObserveChains() []ChainHeight`

GetObserveChains returns the ObserveChains field if non-nil, zero value otherwise.

### GetObserveChainsOk

`func (o *Node) GetObserveChainsOk() (*[]ChainHeight, bool)`

GetObserveChainsOk returns a tuple with the ObserveChains field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObserveChains

`func (o *Node) SetObserveChains(v []ChainHeight)`

SetObserveChains sets ObserveChains field to given value.


### GetPreflightStatus

`func (o *Node) GetPreflightStatus() NodePreflightStatus`

GetPreflightStatus returns the PreflightStatus field if non-nil, zero value otherwise.

### GetPreflightStatusOk

`func (o *Node) GetPreflightStatusOk() (*NodePreflightStatus, bool)`

GetPreflightStatusOk returns a tuple with the PreflightStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreflightStatus

`func (o *Node) SetPreflightStatus(v NodePreflightStatus)`

SetPreflightStatus sets PreflightStatus field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


