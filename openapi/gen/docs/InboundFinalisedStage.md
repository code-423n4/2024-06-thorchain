# InboundFinalisedStage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Completed** | **bool** | returns true if the inbound transaction has been finalised (THORChain agreeing it exists) | 

## Methods

### NewInboundFinalisedStage

`func NewInboundFinalisedStage(completed bool, ) *InboundFinalisedStage`

NewInboundFinalisedStage instantiates a new InboundFinalisedStage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInboundFinalisedStageWithDefaults

`func NewInboundFinalisedStageWithDefaults() *InboundFinalisedStage`

NewInboundFinalisedStageWithDefaults instantiates a new InboundFinalisedStage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCompleted

`func (o *InboundFinalisedStage) GetCompleted() bool`

GetCompleted returns the Completed field if non-nil, zero value otherwise.

### GetCompletedOk

`func (o *InboundFinalisedStage) GetCompletedOk() (*bool, bool)`

GetCompletedOk returns a tuple with the Completed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompleted

`func (o *InboundFinalisedStage) SetCompleted(v bool)`

SetCompleted sets Completed field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


