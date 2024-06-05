# TxStagesResponseInboundObserved

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Started** | Pointer to **bool** | returns true if any nodes have observed the transaction (to be deprecated in favour of counts) | [optional] 
**PreConfirmationCount** | Pointer to **int32** | number of signers for pre-confirmation-counting observations | [optional] 
**FinalCount** | Pointer to **int32** | number of signers for final observations, after any confirmation counting complete | [optional] 
**Completed** | **bool** | returns true if no transaction observation remains to be done | 

## Methods

### NewTxStagesResponseInboundObserved

`func NewTxStagesResponseInboundObserved(completed bool, ) *TxStagesResponseInboundObserved`

NewTxStagesResponseInboundObserved instantiates a new TxStagesResponseInboundObserved object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTxStagesResponseInboundObservedWithDefaults

`func NewTxStagesResponseInboundObservedWithDefaults() *TxStagesResponseInboundObserved`

NewTxStagesResponseInboundObservedWithDefaults instantiates a new TxStagesResponseInboundObserved object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStarted

`func (o *TxStagesResponseInboundObserved) GetStarted() bool`

GetStarted returns the Started field if non-nil, zero value otherwise.

### GetStartedOk

`func (o *TxStagesResponseInboundObserved) GetStartedOk() (*bool, bool)`

GetStartedOk returns a tuple with the Started field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStarted

`func (o *TxStagesResponseInboundObserved) SetStarted(v bool)`

SetStarted sets Started field to given value.

### HasStarted

`func (o *TxStagesResponseInboundObserved) HasStarted() bool`

HasStarted returns a boolean if a field has been set.

### GetPreConfirmationCount

`func (o *TxStagesResponseInboundObserved) GetPreConfirmationCount() int32`

GetPreConfirmationCount returns the PreConfirmationCount field if non-nil, zero value otherwise.

### GetPreConfirmationCountOk

`func (o *TxStagesResponseInboundObserved) GetPreConfirmationCountOk() (*int32, bool)`

GetPreConfirmationCountOk returns a tuple with the PreConfirmationCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreConfirmationCount

`func (o *TxStagesResponseInboundObserved) SetPreConfirmationCount(v int32)`

SetPreConfirmationCount sets PreConfirmationCount field to given value.

### HasPreConfirmationCount

`func (o *TxStagesResponseInboundObserved) HasPreConfirmationCount() bool`

HasPreConfirmationCount returns a boolean if a field has been set.

### GetFinalCount

`func (o *TxStagesResponseInboundObserved) GetFinalCount() int32`

GetFinalCount returns the FinalCount field if non-nil, zero value otherwise.

### GetFinalCountOk

`func (o *TxStagesResponseInboundObserved) GetFinalCountOk() (*int32, bool)`

GetFinalCountOk returns a tuple with the FinalCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFinalCount

`func (o *TxStagesResponseInboundObserved) SetFinalCount(v int32)`

SetFinalCount sets FinalCount field to given value.

### HasFinalCount

`func (o *TxStagesResponseInboundObserved) HasFinalCount() bool`

HasFinalCount returns a boolean if a field has been set.

### GetCompleted

`func (o *TxStagesResponseInboundObserved) GetCompleted() bool`

GetCompleted returns the Completed field if non-nil, zero value otherwise.

### GetCompletedOk

`func (o *TxStagesResponseInboundObserved) GetCompletedOk() (*bool, bool)`

GetCompletedOk returns a tuple with the Completed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompleted

`func (o *TxStagesResponseInboundObserved) SetCompleted(v bool)`

SetCompleted sets Completed field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


