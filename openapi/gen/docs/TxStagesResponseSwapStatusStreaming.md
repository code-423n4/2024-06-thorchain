# TxStagesResponseSwapStatusStreaming

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Interval** | **int32** | how often each swap is made, in blocks | 
**Quantity** | **int32** | the total number of swaps in a streaming swaps | 
**Count** | **int32** | the amount of swap attempts so far | 

## Methods

### NewTxStagesResponseSwapStatusStreaming

`func NewTxStagesResponseSwapStatusStreaming(interval int32, quantity int32, count int32, ) *TxStagesResponseSwapStatusStreaming`

NewTxStagesResponseSwapStatusStreaming instantiates a new TxStagesResponseSwapStatusStreaming object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTxStagesResponseSwapStatusStreamingWithDefaults

`func NewTxStagesResponseSwapStatusStreamingWithDefaults() *TxStagesResponseSwapStatusStreaming`

NewTxStagesResponseSwapStatusStreamingWithDefaults instantiates a new TxStagesResponseSwapStatusStreaming object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInterval

`func (o *TxStagesResponseSwapStatusStreaming) GetInterval() int32`

GetInterval returns the Interval field if non-nil, zero value otherwise.

### GetIntervalOk

`func (o *TxStagesResponseSwapStatusStreaming) GetIntervalOk() (*int32, bool)`

GetIntervalOk returns a tuple with the Interval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterval

`func (o *TxStagesResponseSwapStatusStreaming) SetInterval(v int32)`

SetInterval sets Interval field to given value.


### GetQuantity

`func (o *TxStagesResponseSwapStatusStreaming) GetQuantity() int32`

GetQuantity returns the Quantity field if non-nil, zero value otherwise.

### GetQuantityOk

`func (o *TxStagesResponseSwapStatusStreaming) GetQuantityOk() (*int32, bool)`

GetQuantityOk returns a tuple with the Quantity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuantity

`func (o *TxStagesResponseSwapStatusStreaming) SetQuantity(v int32)`

SetQuantity sets Quantity field to given value.


### GetCount

`func (o *TxStagesResponseSwapStatusStreaming) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *TxStagesResponseSwapStatusStreaming) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *TxStagesResponseSwapStatusStreaming) SetCount(v int32)`

SetCount sets Count field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


