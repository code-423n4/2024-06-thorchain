# QueueResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Swap** | **int64** |  | 
**Outbound** | **int64** | number of signed outbound tx in the queue | 
**Internal** | **int64** |  | 
**ScheduledOutboundValue** | **string** | scheduled outbound value in RUNE | 
**ScheduledOutboundClout** | **string** | scheduled outbound clout in RUNE | 

## Methods

### NewQueueResponse

`func NewQueueResponse(swap int64, outbound int64, internal int64, scheduledOutboundValue string, scheduledOutboundClout string, ) *QueueResponse`

NewQueueResponse instantiates a new QueueResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQueueResponseWithDefaults

`func NewQueueResponseWithDefaults() *QueueResponse`

NewQueueResponseWithDefaults instantiates a new QueueResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSwap

`func (o *QueueResponse) GetSwap() int64`

GetSwap returns the Swap field if non-nil, zero value otherwise.

### GetSwapOk

`func (o *QueueResponse) GetSwapOk() (*int64, bool)`

GetSwapOk returns a tuple with the Swap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwap

`func (o *QueueResponse) SetSwap(v int64)`

SetSwap sets Swap field to given value.


### GetOutbound

`func (o *QueueResponse) GetOutbound() int64`

GetOutbound returns the Outbound field if non-nil, zero value otherwise.

### GetOutboundOk

`func (o *QueueResponse) GetOutboundOk() (*int64, bool)`

GetOutboundOk returns a tuple with the Outbound field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutbound

`func (o *QueueResponse) SetOutbound(v int64)`

SetOutbound sets Outbound field to given value.


### GetInternal

`func (o *QueueResponse) GetInternal() int64`

GetInternal returns the Internal field if non-nil, zero value otherwise.

### GetInternalOk

`func (o *QueueResponse) GetInternalOk() (*int64, bool)`

GetInternalOk returns a tuple with the Internal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternal

`func (o *QueueResponse) SetInternal(v int64)`

SetInternal sets Internal field to given value.


### GetScheduledOutboundValue

`func (o *QueueResponse) GetScheduledOutboundValue() string`

GetScheduledOutboundValue returns the ScheduledOutboundValue field if non-nil, zero value otherwise.

### GetScheduledOutboundValueOk

`func (o *QueueResponse) GetScheduledOutboundValueOk() (*string, bool)`

GetScheduledOutboundValueOk returns a tuple with the ScheduledOutboundValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheduledOutboundValue

`func (o *QueueResponse) SetScheduledOutboundValue(v string)`

SetScheduledOutboundValue sets ScheduledOutboundValue field to given value.


### GetScheduledOutboundClout

`func (o *QueueResponse) GetScheduledOutboundClout() string`

GetScheduledOutboundClout returns the ScheduledOutboundClout field if non-nil, zero value otherwise.

### GetScheduledOutboundCloutOk

`func (o *QueueResponse) GetScheduledOutboundCloutOk() (*string, bool)`

GetScheduledOutboundCloutOk returns a tuple with the ScheduledOutboundClout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheduledOutboundClout

`func (o *QueueResponse) SetScheduledOutboundClout(v string)`

SetScheduledOutboundClout sets ScheduledOutboundClout field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


