# TxStagesResponseSwapFinalised

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Completed** | **bool** | (to be deprecated in favor of swap_status) returns true if an inbound transaction&#39;s swap (successful or refunded) is no longer pending | 

## Methods

### NewTxStagesResponseSwapFinalised

`func NewTxStagesResponseSwapFinalised(completed bool, ) *TxStagesResponseSwapFinalised`

NewTxStagesResponseSwapFinalised instantiates a new TxStagesResponseSwapFinalised object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTxStagesResponseSwapFinalisedWithDefaults

`func NewTxStagesResponseSwapFinalisedWithDefaults() *TxStagesResponseSwapFinalised`

NewTxStagesResponseSwapFinalisedWithDefaults instantiates a new TxStagesResponseSwapFinalised object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCompleted

`func (o *TxStagesResponseSwapFinalised) GetCompleted() bool`

GetCompleted returns the Completed field if non-nil, zero value otherwise.

### GetCompletedOk

`func (o *TxStagesResponseSwapFinalised) GetCompletedOk() (*bool, bool)`

GetCompletedOk returns a tuple with the Completed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompleted

`func (o *TxStagesResponseSwapFinalised) SetCompleted(v bool)`

SetCompleted sets Completed field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


