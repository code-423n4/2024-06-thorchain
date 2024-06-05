# TssKeysignMetric

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TxId** | Pointer to **string** |  | [optional] 
**NodeTssTimes** | [**[]TssMetric**](TssMetric.md) |  | 

## Methods

### NewTssKeysignMetric

`func NewTssKeysignMetric(nodeTssTimes []TssMetric, ) *TssKeysignMetric`

NewTssKeysignMetric instantiates a new TssKeysignMetric object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTssKeysignMetricWithDefaults

`func NewTssKeysignMetricWithDefaults() *TssKeysignMetric`

NewTssKeysignMetricWithDefaults instantiates a new TssKeysignMetric object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTxId

`func (o *TssKeysignMetric) GetTxId() string`

GetTxId returns the TxId field if non-nil, zero value otherwise.

### GetTxIdOk

`func (o *TssKeysignMetric) GetTxIdOk() (*string, bool)`

GetTxIdOk returns a tuple with the TxId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTxId

`func (o *TssKeysignMetric) SetTxId(v string)`

SetTxId sets TxId field to given value.

### HasTxId

`func (o *TssKeysignMetric) HasTxId() bool`

HasTxId returns a boolean if a field has been set.

### GetNodeTssTimes

`func (o *TssKeysignMetric) GetNodeTssTimes() []TssMetric`

GetNodeTssTimes returns the NodeTssTimes field if non-nil, zero value otherwise.

### GetNodeTssTimesOk

`func (o *TssKeysignMetric) GetNodeTssTimesOk() (*[]TssMetric, bool)`

GetNodeTssTimesOk returns a tuple with the NodeTssTimes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeTssTimes

`func (o *TssKeysignMetric) SetNodeTssTimes(v []TssMetric)`

SetNodeTssTimes sets NodeTssTimes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


