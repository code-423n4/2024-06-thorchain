# TxDetailsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TxId** | Pointer to **string** |  | [optional] 
**Tx** | [**ObservedTx**](ObservedTx.md) |  | 
**Txs** | [**[]ObservedTx**](ObservedTx.md) |  | 
**Actions** | [**[]TxOutItem**](TxOutItem.md) |  | 
**OutTxs** | [**[]Tx**](Tx.md) |  | 
**ConsensusHeight** | Pointer to **int64** | the thorchain height at which the inbound reached consensus | [optional] 
**FinalisedHeight** | Pointer to **int64** | the thorchain height at which the outbound was finalised | [optional] 
**UpdatedVault** | Pointer to **bool** |  | [optional] 
**Reverted** | Pointer to **bool** |  | [optional] 
**OutboundHeight** | Pointer to **int64** | the thorchain height for which the outbound was scheduled | [optional] 

## Methods

### NewTxDetailsResponse

`func NewTxDetailsResponse(tx ObservedTx, txs []ObservedTx, actions []TxOutItem, outTxs []Tx, ) *TxDetailsResponse`

NewTxDetailsResponse instantiates a new TxDetailsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTxDetailsResponseWithDefaults

`func NewTxDetailsResponseWithDefaults() *TxDetailsResponse`

NewTxDetailsResponseWithDefaults instantiates a new TxDetailsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTxId

`func (o *TxDetailsResponse) GetTxId() string`

GetTxId returns the TxId field if non-nil, zero value otherwise.

### GetTxIdOk

`func (o *TxDetailsResponse) GetTxIdOk() (*string, bool)`

GetTxIdOk returns a tuple with the TxId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTxId

`func (o *TxDetailsResponse) SetTxId(v string)`

SetTxId sets TxId field to given value.

### HasTxId

`func (o *TxDetailsResponse) HasTxId() bool`

HasTxId returns a boolean if a field has been set.

### GetTx

`func (o *TxDetailsResponse) GetTx() ObservedTx`

GetTx returns the Tx field if non-nil, zero value otherwise.

### GetTxOk

`func (o *TxDetailsResponse) GetTxOk() (*ObservedTx, bool)`

GetTxOk returns a tuple with the Tx field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTx

`func (o *TxDetailsResponse) SetTx(v ObservedTx)`

SetTx sets Tx field to given value.


### GetTxs

`func (o *TxDetailsResponse) GetTxs() []ObservedTx`

GetTxs returns the Txs field if non-nil, zero value otherwise.

### GetTxsOk

`func (o *TxDetailsResponse) GetTxsOk() (*[]ObservedTx, bool)`

GetTxsOk returns a tuple with the Txs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTxs

`func (o *TxDetailsResponse) SetTxs(v []ObservedTx)`

SetTxs sets Txs field to given value.


### GetActions

`func (o *TxDetailsResponse) GetActions() []TxOutItem`

GetActions returns the Actions field if non-nil, zero value otherwise.

### GetActionsOk

`func (o *TxDetailsResponse) GetActionsOk() (*[]TxOutItem, bool)`

GetActionsOk returns a tuple with the Actions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActions

`func (o *TxDetailsResponse) SetActions(v []TxOutItem)`

SetActions sets Actions field to given value.


### GetOutTxs

`func (o *TxDetailsResponse) GetOutTxs() []Tx`

GetOutTxs returns the OutTxs field if non-nil, zero value otherwise.

### GetOutTxsOk

`func (o *TxDetailsResponse) GetOutTxsOk() (*[]Tx, bool)`

GetOutTxsOk returns a tuple with the OutTxs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutTxs

`func (o *TxDetailsResponse) SetOutTxs(v []Tx)`

SetOutTxs sets OutTxs field to given value.


### GetConsensusHeight

`func (o *TxDetailsResponse) GetConsensusHeight() int64`

GetConsensusHeight returns the ConsensusHeight field if non-nil, zero value otherwise.

### GetConsensusHeightOk

`func (o *TxDetailsResponse) GetConsensusHeightOk() (*int64, bool)`

GetConsensusHeightOk returns a tuple with the ConsensusHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConsensusHeight

`func (o *TxDetailsResponse) SetConsensusHeight(v int64)`

SetConsensusHeight sets ConsensusHeight field to given value.

### HasConsensusHeight

`func (o *TxDetailsResponse) HasConsensusHeight() bool`

HasConsensusHeight returns a boolean if a field has been set.

### GetFinalisedHeight

`func (o *TxDetailsResponse) GetFinalisedHeight() int64`

GetFinalisedHeight returns the FinalisedHeight field if non-nil, zero value otherwise.

### GetFinalisedHeightOk

`func (o *TxDetailsResponse) GetFinalisedHeightOk() (*int64, bool)`

GetFinalisedHeightOk returns a tuple with the FinalisedHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFinalisedHeight

`func (o *TxDetailsResponse) SetFinalisedHeight(v int64)`

SetFinalisedHeight sets FinalisedHeight field to given value.

### HasFinalisedHeight

`func (o *TxDetailsResponse) HasFinalisedHeight() bool`

HasFinalisedHeight returns a boolean if a field has been set.

### GetUpdatedVault

`func (o *TxDetailsResponse) GetUpdatedVault() bool`

GetUpdatedVault returns the UpdatedVault field if non-nil, zero value otherwise.

### GetUpdatedVaultOk

`func (o *TxDetailsResponse) GetUpdatedVaultOk() (*bool, bool)`

GetUpdatedVaultOk returns a tuple with the UpdatedVault field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedVault

`func (o *TxDetailsResponse) SetUpdatedVault(v bool)`

SetUpdatedVault sets UpdatedVault field to given value.

### HasUpdatedVault

`func (o *TxDetailsResponse) HasUpdatedVault() bool`

HasUpdatedVault returns a boolean if a field has been set.

### GetReverted

`func (o *TxDetailsResponse) GetReverted() bool`

GetReverted returns the Reverted field if non-nil, zero value otherwise.

### GetRevertedOk

`func (o *TxDetailsResponse) GetRevertedOk() (*bool, bool)`

GetRevertedOk returns a tuple with the Reverted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReverted

`func (o *TxDetailsResponse) SetReverted(v bool)`

SetReverted sets Reverted field to given value.

### HasReverted

`func (o *TxDetailsResponse) HasReverted() bool`

HasReverted returns a boolean if a field has been set.

### GetOutboundHeight

`func (o *TxDetailsResponse) GetOutboundHeight() int64`

GetOutboundHeight returns the OutboundHeight field if non-nil, zero value otherwise.

### GetOutboundHeightOk

`func (o *TxDetailsResponse) GetOutboundHeightOk() (*int64, bool)`

GetOutboundHeightOk returns a tuple with the OutboundHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundHeight

`func (o *TxDetailsResponse) SetOutboundHeight(v int64)`

SetOutboundHeight sets OutboundHeight field to given value.

### HasOutboundHeight

`func (o *TxDetailsResponse) HasOutboundHeight() bool`

HasOutboundHeight returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


