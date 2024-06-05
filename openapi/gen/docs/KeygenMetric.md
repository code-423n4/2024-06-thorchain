# KeygenMetric

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PubKey** | Pointer to **string** |  | [optional] 
**NodeTssTimes** | [**[]NodeKeygenMetric**](NodeKeygenMetric.md) |  | 

## Methods

### NewKeygenMetric

`func NewKeygenMetric(nodeTssTimes []NodeKeygenMetric, ) *KeygenMetric`

NewKeygenMetric instantiates a new KeygenMetric object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeygenMetricWithDefaults

`func NewKeygenMetricWithDefaults() *KeygenMetric`

NewKeygenMetricWithDefaults instantiates a new KeygenMetric object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPubKey

`func (o *KeygenMetric) GetPubKey() string`

GetPubKey returns the PubKey field if non-nil, zero value otherwise.

### GetPubKeyOk

`func (o *KeygenMetric) GetPubKeyOk() (*string, bool)`

GetPubKeyOk returns a tuple with the PubKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPubKey

`func (o *KeygenMetric) SetPubKey(v string)`

SetPubKey sets PubKey field to given value.

### HasPubKey

`func (o *KeygenMetric) HasPubKey() bool`

HasPubKey returns a boolean if a field has been set.

### GetNodeTssTimes

`func (o *KeygenMetric) GetNodeTssTimes() []NodeKeygenMetric`

GetNodeTssTimes returns the NodeTssTimes field if non-nil, zero value otherwise.

### GetNodeTssTimesOk

`func (o *KeygenMetric) GetNodeTssTimesOk() (*[]NodeKeygenMetric, bool)`

GetNodeTssTimesOk returns a tuple with the NodeTssTimes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeTssTimes

`func (o *KeygenMetric) SetNodeTssTimes(v []NodeKeygenMetric)`

SetNodeTssTimes sets NodeTssTimes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


