# OutboundSignedStage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ScheduledOutboundHeight** | Pointer to **int64** | THORChain height for which the external outbound is scheduled | [optional] 
**BlocksSinceScheduled** | Pointer to **int64** | THORChain blocks since the scheduled outbound height | [optional] 
**Completed** | **bool** | returns true if an external transaction has been signed and broadcast (and observed in its mempool) | 

## Methods

### NewOutboundSignedStage

`func NewOutboundSignedStage(completed bool, ) *OutboundSignedStage`

NewOutboundSignedStage instantiates a new OutboundSignedStage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOutboundSignedStageWithDefaults

`func NewOutboundSignedStageWithDefaults() *OutboundSignedStage`

NewOutboundSignedStageWithDefaults instantiates a new OutboundSignedStage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetScheduledOutboundHeight

`func (o *OutboundSignedStage) GetScheduledOutboundHeight() int64`

GetScheduledOutboundHeight returns the ScheduledOutboundHeight field if non-nil, zero value otherwise.

### GetScheduledOutboundHeightOk

`func (o *OutboundSignedStage) GetScheduledOutboundHeightOk() (*int64, bool)`

GetScheduledOutboundHeightOk returns a tuple with the ScheduledOutboundHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheduledOutboundHeight

`func (o *OutboundSignedStage) SetScheduledOutboundHeight(v int64)`

SetScheduledOutboundHeight sets ScheduledOutboundHeight field to given value.

### HasScheduledOutboundHeight

`func (o *OutboundSignedStage) HasScheduledOutboundHeight() bool`

HasScheduledOutboundHeight returns a boolean if a field has been set.

### GetBlocksSinceScheduled

`func (o *OutboundSignedStage) GetBlocksSinceScheduled() int64`

GetBlocksSinceScheduled returns the BlocksSinceScheduled field if non-nil, zero value otherwise.

### GetBlocksSinceScheduledOk

`func (o *OutboundSignedStage) GetBlocksSinceScheduledOk() (*int64, bool)`

GetBlocksSinceScheduledOk returns a tuple with the BlocksSinceScheduled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlocksSinceScheduled

`func (o *OutboundSignedStage) SetBlocksSinceScheduled(v int64)`

SetBlocksSinceScheduled sets BlocksSinceScheduled field to given value.

### HasBlocksSinceScheduled

`func (o *OutboundSignedStage) HasBlocksSinceScheduled() bool`

HasBlocksSinceScheduled returns a boolean if a field has been set.

### GetCompleted

`func (o *OutboundSignedStage) GetCompleted() bool`

GetCompleted returns the Completed field if non-nil, zero value otherwise.

### GetCompletedOk

`func (o *OutboundSignedStage) GetCompletedOk() (*bool, bool)`

GetCompletedOk returns a tuple with the Completed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompleted

`func (o *OutboundSignedStage) SetCompleted(v bool)`

SetCompleted sets Completed field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


