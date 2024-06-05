# InboundObservedStage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Started** | Pointer to **bool** | returns true if any nodes have observed the transaction (to be deprecated in favour of counts) | [optional] 
**PreConfirmationCount** | Pointer to **int64** | number of signers for pre-confirmation-counting observations | [optional] 
**FinalCount** | **int64** | number of signers for final observations, after any confirmation counting complete | 
**Completed** | **bool** | returns true if no transaction observation remains to be done | 

## Methods

### NewInboundObservedStage

`func NewInboundObservedStage(finalCount int64, completed bool, ) *InboundObservedStage`

NewInboundObservedStage instantiates a new InboundObservedStage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInboundObservedStageWithDefaults

`func NewInboundObservedStageWithDefaults() *InboundObservedStage`

NewInboundObservedStageWithDefaults instantiates a new InboundObservedStage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStarted

`func (o *InboundObservedStage) GetStarted() bool`

GetStarted returns the Started field if non-nil, zero value otherwise.

### GetStartedOk

`func (o *InboundObservedStage) GetStartedOk() (*bool, bool)`

GetStartedOk returns a tuple with the Started field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStarted

`func (o *InboundObservedStage) SetStarted(v bool)`

SetStarted sets Started field to given value.

### HasStarted

`func (o *InboundObservedStage) HasStarted() bool`

HasStarted returns a boolean if a field has been set.

### GetPreConfirmationCount

`func (o *InboundObservedStage) GetPreConfirmationCount() int64`

GetPreConfirmationCount returns the PreConfirmationCount field if non-nil, zero value otherwise.

### GetPreConfirmationCountOk

`func (o *InboundObservedStage) GetPreConfirmationCountOk() (*int64, bool)`

GetPreConfirmationCountOk returns a tuple with the PreConfirmationCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreConfirmationCount

`func (o *InboundObservedStage) SetPreConfirmationCount(v int64)`

SetPreConfirmationCount sets PreConfirmationCount field to given value.

### HasPreConfirmationCount

`func (o *InboundObservedStage) HasPreConfirmationCount() bool`

HasPreConfirmationCount returns a boolean if a field has been set.

### GetFinalCount

`func (o *InboundObservedStage) GetFinalCount() int64`

GetFinalCount returns the FinalCount field if non-nil, zero value otherwise.

### GetFinalCountOk

`func (o *InboundObservedStage) GetFinalCountOk() (*int64, bool)`

GetFinalCountOk returns a tuple with the FinalCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFinalCount

`func (o *InboundObservedStage) SetFinalCount(v int64)`

SetFinalCount sets FinalCount field to given value.


### GetCompleted

`func (o *InboundObservedStage) GetCompleted() bool`

GetCompleted returns the Completed field if non-nil, zero value otherwise.

### GetCompletedOk

`func (o *InboundObservedStage) GetCompletedOk() (*bool, bool)`

GetCompletedOk returns a tuple with the Completed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompleted

`func (o *InboundObservedStage) SetCompleted(v bool)`

SetCompleted sets Completed field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


