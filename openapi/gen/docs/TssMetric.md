# TssMetric

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Address** | Pointer to **string** |  | [optional] 
**TssTime** | Pointer to **int64** |  | [optional] 

## Methods

### NewTssMetric

`func NewTssMetric() *TssMetric`

NewTssMetric instantiates a new TssMetric object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTssMetricWithDefaults

`func NewTssMetricWithDefaults() *TssMetric`

NewTssMetricWithDefaults instantiates a new TssMetric object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAddress

`func (o *TssMetric) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *TssMetric) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *TssMetric) SetAddress(v string)`

SetAddress sets Address field to given value.

### HasAddress

`func (o *TssMetric) HasAddress() bool`

HasAddress returns a boolean if a field has been set.

### GetTssTime

`func (o *TssMetric) GetTssTime() int64`

GetTssTime returns the TssTime field if non-nil, zero value otherwise.

### GetTssTimeOk

`func (o *TssMetric) GetTssTimeOk() (*int64, bool)`

GetTssTimeOk returns a tuple with the TssTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTssTime

`func (o *TssMetric) SetTssTime(v int64)`

SetTssTime sets TssTime field to given value.

### HasTssTime

`func (o *TssMetric) HasTssTime() bool`

HasTssTime returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


