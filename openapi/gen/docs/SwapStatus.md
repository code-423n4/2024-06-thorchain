# SwapStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Pending** | **bool** | true when awaiting a swap | 
**Streaming** | Pointer to [**StreamingStatus**](StreamingStatus.md) |  | [optional] 

## Methods

### NewSwapStatus

`func NewSwapStatus(pending bool, ) *SwapStatus`

NewSwapStatus instantiates a new SwapStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwapStatusWithDefaults

`func NewSwapStatusWithDefaults() *SwapStatus`

NewSwapStatusWithDefaults instantiates a new SwapStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPending

`func (o *SwapStatus) GetPending() bool`

GetPending returns the Pending field if non-nil, zero value otherwise.

### GetPendingOk

`func (o *SwapStatus) GetPendingOk() (*bool, bool)`

GetPendingOk returns a tuple with the Pending field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPending

`func (o *SwapStatus) SetPending(v bool)`

SetPending sets Pending field to given value.


### GetStreaming

`func (o *SwapStatus) GetStreaming() StreamingStatus`

GetStreaming returns the Streaming field if non-nil, zero value otherwise.

### GetStreamingOk

`func (o *SwapStatus) GetStreamingOk() (*StreamingStatus, bool)`

GetStreamingOk returns a tuple with the Streaming field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreaming

`func (o *SwapStatus) SetStreaming(v StreamingStatus)`

SetStreaming sets Streaming field to given value.

### HasStreaming

`func (o *SwapStatus) HasStreaming() bool`

HasStreaming returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


