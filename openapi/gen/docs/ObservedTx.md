# ObservedTx

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Tx** | [**Tx**](Tx.md) |  | 
**ObservedPubKey** | Pointer to **string** |  | [optional] 
**ExternalObservedHeight** | Pointer to **int64** | the block height on the external source chain when the transaction was observed, not provided if chain is THOR | [optional] 
**ExternalConfirmationDelayHeight** | Pointer to **int64** | the block height on the external source chain when confirmation counting will be complete, not provided if chain is THOR | [optional] 
**Aggregator** | Pointer to **string** | the outbound aggregator to use, will also match a suffix | [optional] 
**AggregatorTarget** | Pointer to **string** | the aggregator target asset provided to transferOutAndCall | [optional] 
**AggregatorTargetLimit** | Pointer to **string** | the aggregator target asset limit provided to transferOutAndCall | [optional] 
**Signers** | Pointer to **[]string** |  | [optional] 
**KeysignMs** | Pointer to **int64** |  | [optional] 
**OutHashes** | Pointer to **[]string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Methods

### NewObservedTx

`func NewObservedTx(tx Tx, ) *ObservedTx`

NewObservedTx instantiates a new ObservedTx object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewObservedTxWithDefaults

`func NewObservedTxWithDefaults() *ObservedTx`

NewObservedTxWithDefaults instantiates a new ObservedTx object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTx

`func (o *ObservedTx) GetTx() Tx`

GetTx returns the Tx field if non-nil, zero value otherwise.

### GetTxOk

`func (o *ObservedTx) GetTxOk() (*Tx, bool)`

GetTxOk returns a tuple with the Tx field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTx

`func (o *ObservedTx) SetTx(v Tx)`

SetTx sets Tx field to given value.


### GetObservedPubKey

`func (o *ObservedTx) GetObservedPubKey() string`

GetObservedPubKey returns the ObservedPubKey field if non-nil, zero value otherwise.

### GetObservedPubKeyOk

`func (o *ObservedTx) GetObservedPubKeyOk() (*string, bool)`

GetObservedPubKeyOk returns a tuple with the ObservedPubKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObservedPubKey

`func (o *ObservedTx) SetObservedPubKey(v string)`

SetObservedPubKey sets ObservedPubKey field to given value.

### HasObservedPubKey

`func (o *ObservedTx) HasObservedPubKey() bool`

HasObservedPubKey returns a boolean if a field has been set.

### GetExternalObservedHeight

`func (o *ObservedTx) GetExternalObservedHeight() int64`

GetExternalObservedHeight returns the ExternalObservedHeight field if non-nil, zero value otherwise.

### GetExternalObservedHeightOk

`func (o *ObservedTx) GetExternalObservedHeightOk() (*int64, bool)`

GetExternalObservedHeightOk returns a tuple with the ExternalObservedHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalObservedHeight

`func (o *ObservedTx) SetExternalObservedHeight(v int64)`

SetExternalObservedHeight sets ExternalObservedHeight field to given value.

### HasExternalObservedHeight

`func (o *ObservedTx) HasExternalObservedHeight() bool`

HasExternalObservedHeight returns a boolean if a field has been set.

### GetExternalConfirmationDelayHeight

`func (o *ObservedTx) GetExternalConfirmationDelayHeight() int64`

GetExternalConfirmationDelayHeight returns the ExternalConfirmationDelayHeight field if non-nil, zero value otherwise.

### GetExternalConfirmationDelayHeightOk

`func (o *ObservedTx) GetExternalConfirmationDelayHeightOk() (*int64, bool)`

GetExternalConfirmationDelayHeightOk returns a tuple with the ExternalConfirmationDelayHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalConfirmationDelayHeight

`func (o *ObservedTx) SetExternalConfirmationDelayHeight(v int64)`

SetExternalConfirmationDelayHeight sets ExternalConfirmationDelayHeight field to given value.

### HasExternalConfirmationDelayHeight

`func (o *ObservedTx) HasExternalConfirmationDelayHeight() bool`

HasExternalConfirmationDelayHeight returns a boolean if a field has been set.

### GetAggregator

`func (o *ObservedTx) GetAggregator() string`

GetAggregator returns the Aggregator field if non-nil, zero value otherwise.

### GetAggregatorOk

`func (o *ObservedTx) GetAggregatorOk() (*string, bool)`

GetAggregatorOk returns a tuple with the Aggregator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggregator

`func (o *ObservedTx) SetAggregator(v string)`

SetAggregator sets Aggregator field to given value.

### HasAggregator

`func (o *ObservedTx) HasAggregator() bool`

HasAggregator returns a boolean if a field has been set.

### GetAggregatorTarget

`func (o *ObservedTx) GetAggregatorTarget() string`

GetAggregatorTarget returns the AggregatorTarget field if non-nil, zero value otherwise.

### GetAggregatorTargetOk

`func (o *ObservedTx) GetAggregatorTargetOk() (*string, bool)`

GetAggregatorTargetOk returns a tuple with the AggregatorTarget field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggregatorTarget

`func (o *ObservedTx) SetAggregatorTarget(v string)`

SetAggregatorTarget sets AggregatorTarget field to given value.

### HasAggregatorTarget

`func (o *ObservedTx) HasAggregatorTarget() bool`

HasAggregatorTarget returns a boolean if a field has been set.

### GetAggregatorTargetLimit

`func (o *ObservedTx) GetAggregatorTargetLimit() string`

GetAggregatorTargetLimit returns the AggregatorTargetLimit field if non-nil, zero value otherwise.

### GetAggregatorTargetLimitOk

`func (o *ObservedTx) GetAggregatorTargetLimitOk() (*string, bool)`

GetAggregatorTargetLimitOk returns a tuple with the AggregatorTargetLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggregatorTargetLimit

`func (o *ObservedTx) SetAggregatorTargetLimit(v string)`

SetAggregatorTargetLimit sets AggregatorTargetLimit field to given value.

### HasAggregatorTargetLimit

`func (o *ObservedTx) HasAggregatorTargetLimit() bool`

HasAggregatorTargetLimit returns a boolean if a field has been set.

### GetSigners

`func (o *ObservedTx) GetSigners() []string`

GetSigners returns the Signers field if non-nil, zero value otherwise.

### GetSignersOk

`func (o *ObservedTx) GetSignersOk() (*[]string, bool)`

GetSignersOk returns a tuple with the Signers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSigners

`func (o *ObservedTx) SetSigners(v []string)`

SetSigners sets Signers field to given value.

### HasSigners

`func (o *ObservedTx) HasSigners() bool`

HasSigners returns a boolean if a field has been set.

### GetKeysignMs

`func (o *ObservedTx) GetKeysignMs() int64`

GetKeysignMs returns the KeysignMs field if non-nil, zero value otherwise.

### GetKeysignMsOk

`func (o *ObservedTx) GetKeysignMsOk() (*int64, bool)`

GetKeysignMsOk returns a tuple with the KeysignMs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeysignMs

`func (o *ObservedTx) SetKeysignMs(v int64)`

SetKeysignMs sets KeysignMs field to given value.

### HasKeysignMs

`func (o *ObservedTx) HasKeysignMs() bool`

HasKeysignMs returns a boolean if a field has been set.

### GetOutHashes

`func (o *ObservedTx) GetOutHashes() []string`

GetOutHashes returns the OutHashes field if non-nil, zero value otherwise.

### GetOutHashesOk

`func (o *ObservedTx) GetOutHashesOk() (*[]string, bool)`

GetOutHashesOk returns a tuple with the OutHashes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutHashes

`func (o *ObservedTx) SetOutHashes(v []string)`

SetOutHashes sets OutHashes field to given value.

### HasOutHashes

`func (o *ObservedTx) HasOutHashes() bool`

HasOutHashes returns a boolean if a field has been set.

### GetStatus

`func (o *ObservedTx) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ObservedTx) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ObservedTx) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ObservedTx) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


