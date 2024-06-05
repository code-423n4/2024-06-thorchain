# NodeJail

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ReleaseHeight** | Pointer to **int64** |  | [optional] 
**Reason** | Pointer to **string** |  | [optional] 

## Methods

### NewNodeJail

`func NewNodeJail() *NodeJail`

NewNodeJail instantiates a new NodeJail object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeJailWithDefaults

`func NewNodeJailWithDefaults() *NodeJail`

NewNodeJailWithDefaults instantiates a new NodeJail object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetReleaseHeight

`func (o *NodeJail) GetReleaseHeight() int64`

GetReleaseHeight returns the ReleaseHeight field if non-nil, zero value otherwise.

### GetReleaseHeightOk

`func (o *NodeJail) GetReleaseHeightOk() (*int64, bool)`

GetReleaseHeightOk returns a tuple with the ReleaseHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReleaseHeight

`func (o *NodeJail) SetReleaseHeight(v int64)`

SetReleaseHeight sets ReleaseHeight field to given value.

### HasReleaseHeight

`func (o *NodeJail) HasReleaseHeight() bool`

HasReleaseHeight returns a boolean if a field has been set.

### GetReason

`func (o *NodeJail) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *NodeJail) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *NodeJail) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *NodeJail) HasReason() bool`

HasReason returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


