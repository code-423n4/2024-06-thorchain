# VersionResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Current** | **string** | current version | 
**Next** | **string** | next version (minimum version for a node to become Active) | 
**NextSinceHeight** | Pointer to **int64** | height at which the minimum joining version last changed | [optional] 
**Querier** | **string** | querier version | 

## Methods

### NewVersionResponse

`func NewVersionResponse(current string, next string, querier string, ) *VersionResponse`

NewVersionResponse instantiates a new VersionResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVersionResponseWithDefaults

`func NewVersionResponseWithDefaults() *VersionResponse`

NewVersionResponseWithDefaults instantiates a new VersionResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCurrent

`func (o *VersionResponse) GetCurrent() string`

GetCurrent returns the Current field if non-nil, zero value otherwise.

### GetCurrentOk

`func (o *VersionResponse) GetCurrentOk() (*string, bool)`

GetCurrentOk returns a tuple with the Current field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrent

`func (o *VersionResponse) SetCurrent(v string)`

SetCurrent sets Current field to given value.


### GetNext

`func (o *VersionResponse) GetNext() string`

GetNext returns the Next field if non-nil, zero value otherwise.

### GetNextOk

`func (o *VersionResponse) GetNextOk() (*string, bool)`

GetNextOk returns a tuple with the Next field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNext

`func (o *VersionResponse) SetNext(v string)`

SetNext sets Next field to given value.


### GetNextSinceHeight

`func (o *VersionResponse) GetNextSinceHeight() int64`

GetNextSinceHeight returns the NextSinceHeight field if non-nil, zero value otherwise.

### GetNextSinceHeightOk

`func (o *VersionResponse) GetNextSinceHeightOk() (*int64, bool)`

GetNextSinceHeightOk returns a tuple with the NextSinceHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextSinceHeight

`func (o *VersionResponse) SetNextSinceHeight(v int64)`

SetNextSinceHeight sets NextSinceHeight field to given value.

### HasNextSinceHeight

`func (o *VersionResponse) HasNextSinceHeight() bool`

HasNextSinceHeight returns a boolean if a field has been set.

### GetQuerier

`func (o *VersionResponse) GetQuerier() string`

GetQuerier returns the Querier field if non-nil, zero value otherwise.

### GetQuerierOk

`func (o *VersionResponse) GetQuerierOk() (*string, bool)`

GetQuerierOk returns a tuple with the Querier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuerier

`func (o *VersionResponse) SetQuerier(v string)`

SetQuerier sets Querier field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


