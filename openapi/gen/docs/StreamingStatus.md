# StreamingStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Interval** | **int64** | how often each swap is made, in blocks | 
**Quantity** | **int64** | the total number of swaps in a streaming swaps | 
**Count** | **int64** | the amount of swap attempts so far | 

## Methods

### NewStreamingStatus

`func NewStreamingStatus(interval int64, quantity int64, count int64, ) *StreamingStatus`

NewStreamingStatus instantiates a new StreamingStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStreamingStatusWithDefaults

`func NewStreamingStatusWithDefaults() *StreamingStatus`

NewStreamingStatusWithDefaults instantiates a new StreamingStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInterval

`func (o *StreamingStatus) GetInterval() int64`

GetInterval returns the Interval field if non-nil, zero value otherwise.

### GetIntervalOk

`func (o *StreamingStatus) GetIntervalOk() (*int64, bool)`

GetIntervalOk returns a tuple with the Interval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterval

`func (o *StreamingStatus) SetInterval(v int64)`

SetInterval sets Interval field to given value.


### GetQuantity

`func (o *StreamingStatus) GetQuantity() int64`

GetQuantity returns the Quantity field if non-nil, zero value otherwise.

### GetQuantityOk

`func (o *StreamingStatus) GetQuantityOk() (*int64, bool)`

GetQuantityOk returns a tuple with the Quantity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuantity

`func (o *StreamingStatus) SetQuantity(v int64)`

SetQuantity sets Quantity field to given value.


### GetCount

`func (o *StreamingStatus) GetCount() int64`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *StreamingStatus) GetCountOk() (*int64, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *StreamingStatus) SetCount(v int64)`

SetCount sets Count field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


