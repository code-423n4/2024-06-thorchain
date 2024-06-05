# TxStatusResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Tx** | Pointer to [**Tx**](Tx.md) |  | [optional] 
**PlannedOutTxs** | Pointer to [**[]PlannedOutTx**](PlannedOutTx.md) |  | [optional] 
**OutTxs** | Pointer to [**[]Tx**](Tx.md) |  | [optional] 
**Stages** | [**TxStagesResponse**](TxStagesResponse.md) |  | 

## Methods

### NewTxStatusResponse

`func NewTxStatusResponse(stages TxStagesResponse, ) *TxStatusResponse`

NewTxStatusResponse instantiates a new TxStatusResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTxStatusResponseWithDefaults

`func NewTxStatusResponseWithDefaults() *TxStatusResponse`

NewTxStatusResponseWithDefaults instantiates a new TxStatusResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTx

`func (o *TxStatusResponse) GetTx() Tx`

GetTx returns the Tx field if non-nil, zero value otherwise.

### GetTxOk

`func (o *TxStatusResponse) GetTxOk() (*Tx, bool)`

GetTxOk returns a tuple with the Tx field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTx

`func (o *TxStatusResponse) SetTx(v Tx)`

SetTx sets Tx field to given value.

### HasTx

`func (o *TxStatusResponse) HasTx() bool`

HasTx returns a boolean if a field has been set.

### GetPlannedOutTxs

`func (o *TxStatusResponse) GetPlannedOutTxs() []PlannedOutTx`

GetPlannedOutTxs returns the PlannedOutTxs field if non-nil, zero value otherwise.

### GetPlannedOutTxsOk

`func (o *TxStatusResponse) GetPlannedOutTxsOk() (*[]PlannedOutTx, bool)`

GetPlannedOutTxsOk returns a tuple with the PlannedOutTxs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlannedOutTxs

`func (o *TxStatusResponse) SetPlannedOutTxs(v []PlannedOutTx)`

SetPlannedOutTxs sets PlannedOutTxs field to given value.

### HasPlannedOutTxs

`func (o *TxStatusResponse) HasPlannedOutTxs() bool`

HasPlannedOutTxs returns a boolean if a field has been set.

### GetOutTxs

`func (o *TxStatusResponse) GetOutTxs() []Tx`

GetOutTxs returns the OutTxs field if non-nil, zero value otherwise.

### GetOutTxsOk

`func (o *TxStatusResponse) GetOutTxsOk() (*[]Tx, bool)`

GetOutTxsOk returns a tuple with the OutTxs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutTxs

`func (o *TxStatusResponse) SetOutTxs(v []Tx)`

SetOutTxs sets OutTxs field to given value.

### HasOutTxs

`func (o *TxStatusResponse) HasOutTxs() bool`

HasOutTxs returns a boolean if a field has been set.

### GetStages

`func (o *TxStatusResponse) GetStages() TxStagesResponse`

GetStages returns the Stages field if non-nil, zero value otherwise.

### GetStagesOk

`func (o *TxStatusResponse) GetStagesOk() (*TxStagesResponse, bool)`

GetStagesOk returns a tuple with the Stages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStages

`func (o *TxStatusResponse) SetStages(v TxStagesResponse)`

SetStages sets Stages field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


