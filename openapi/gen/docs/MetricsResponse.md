# MetricsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Keygen** | Pointer to [**[]KeygenMetric**](KeygenMetric.md) |  | [optional] 
**Keysign** | Pointer to [**KeysignMetrics**](KeysignMetrics.md) |  | [optional] 

## Methods

### NewMetricsResponse

`func NewMetricsResponse() *MetricsResponse`

NewMetricsResponse instantiates a new MetricsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMetricsResponseWithDefaults

`func NewMetricsResponseWithDefaults() *MetricsResponse`

NewMetricsResponseWithDefaults instantiates a new MetricsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeygen

`func (o *MetricsResponse) GetKeygen() []KeygenMetric`

GetKeygen returns the Keygen field if non-nil, zero value otherwise.

### GetKeygenOk

`func (o *MetricsResponse) GetKeygenOk() (*[]KeygenMetric, bool)`

GetKeygenOk returns a tuple with the Keygen field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeygen

`func (o *MetricsResponse) SetKeygen(v []KeygenMetric)`

SetKeygen sets Keygen field to given value.

### HasKeygen

`func (o *MetricsResponse) HasKeygen() bool`

HasKeygen returns a boolean if a field has been set.

### GetKeysign

`func (o *MetricsResponse) GetKeysign() KeysignMetrics`

GetKeysign returns the Keysign field if non-nil, zero value otherwise.

### GetKeysignOk

`func (o *MetricsResponse) GetKeysignOk() (*KeysignMetrics, bool)`

GetKeysignOk returns a tuple with the Keysign field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeysign

`func (o *MetricsResponse) SetKeysign(v KeysignMetrics)`

SetKeysign sets Keysign field to given value.

### HasKeysign

`func (o *MetricsResponse) HasKeysign() bool`

HasKeysign returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


