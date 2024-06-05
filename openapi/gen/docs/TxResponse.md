# TxResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ObservedTx** | Pointer to [**ObservedTx**](ObservedTx.md) |  | [optional] 
**ConsensusHeight** | Pointer to **int64** | the thorchain height at which the inbound reached consensus | [optional] 
**FinalisedHeight** | Pointer to **int64** | the thorchain height at which the outbound was finalised | [optional] 
**OutboundHeight** | Pointer to **int64** | the thorchain height for which the outbound was scheduled | [optional] 
**KeysignMetric** | Pointer to [**TssKeysignMetric**](TssKeysignMetric.md) |  | [optional] 

## Methods

### NewTxResponse

`func NewTxResponse() *TxResponse`

NewTxResponse instantiates a new TxResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTxResponseWithDefaults

`func NewTxResponseWithDefaults() *TxResponse`

NewTxResponseWithDefaults instantiates a new TxResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetObservedTx

`func (o *TxResponse) GetObservedTx() ObservedTx`

GetObservedTx returns the ObservedTx field if non-nil, zero value otherwise.

### GetObservedTxOk

`func (o *TxResponse) GetObservedTxOk() (*ObservedTx, bool)`

GetObservedTxOk returns a tuple with the ObservedTx field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObservedTx

`func (o *TxResponse) SetObservedTx(v ObservedTx)`

SetObservedTx sets ObservedTx field to given value.

### HasObservedTx

`func (o *TxResponse) HasObservedTx() bool`

HasObservedTx returns a boolean if a field has been set.

### GetConsensusHeight

`func (o *TxResponse) GetConsensusHeight() int64`

GetConsensusHeight returns the ConsensusHeight field if non-nil, zero value otherwise.

### GetConsensusHeightOk

`func (o *TxResponse) GetConsensusHeightOk() (*int64, bool)`

GetConsensusHeightOk returns a tuple with the ConsensusHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConsensusHeight

`func (o *TxResponse) SetConsensusHeight(v int64)`

SetConsensusHeight sets ConsensusHeight field to given value.

### HasConsensusHeight

`func (o *TxResponse) HasConsensusHeight() bool`

HasConsensusHeight returns a boolean if a field has been set.

### GetFinalisedHeight

`func (o *TxResponse) GetFinalisedHeight() int64`

GetFinalisedHeight returns the FinalisedHeight field if non-nil, zero value otherwise.

### GetFinalisedHeightOk

`func (o *TxResponse) GetFinalisedHeightOk() (*int64, bool)`

GetFinalisedHeightOk returns a tuple with the FinalisedHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFinalisedHeight

`func (o *TxResponse) SetFinalisedHeight(v int64)`

SetFinalisedHeight sets FinalisedHeight field to given value.

### HasFinalisedHeight

`func (o *TxResponse) HasFinalisedHeight() bool`

HasFinalisedHeight returns a boolean if a field has been set.

### GetOutboundHeight

`func (o *TxResponse) GetOutboundHeight() int64`

GetOutboundHeight returns the OutboundHeight field if non-nil, zero value otherwise.

### GetOutboundHeightOk

`func (o *TxResponse) GetOutboundHeightOk() (*int64, bool)`

GetOutboundHeightOk returns a tuple with the OutboundHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundHeight

`func (o *TxResponse) SetOutboundHeight(v int64)`

SetOutboundHeight sets OutboundHeight field to given value.

### HasOutboundHeight

`func (o *TxResponse) HasOutboundHeight() bool`

HasOutboundHeight returns a boolean if a field has been set.

### GetKeysignMetric

`func (o *TxResponse) GetKeysignMetric() TssKeysignMetric`

GetKeysignMetric returns the KeysignMetric field if non-nil, zero value otherwise.

### GetKeysignMetricOk

`func (o *TxResponse) GetKeysignMetricOk() (*TssKeysignMetric, bool)`

GetKeysignMetricOk returns a tuple with the KeysignMetric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeysignMetric

`func (o *TxResponse) SetKeysignMetric(v TssKeysignMetric)`

SetKeysignMetric sets KeysignMetric field to given value.

### HasKeysignMetric

`func (o *TxResponse) HasKeysignMetric() bool`

HasKeysignMetric returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


