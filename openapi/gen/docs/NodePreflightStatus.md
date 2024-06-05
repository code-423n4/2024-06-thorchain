# NodePreflightStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | the next status of the node | 
**Reason** | **string** | the reason for the transition to the next status | 
**Code** | **int64** |  | 

## Methods

### NewNodePreflightStatus

`func NewNodePreflightStatus(status string, reason string, code int64, ) *NodePreflightStatus`

NewNodePreflightStatus instantiates a new NodePreflightStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodePreflightStatusWithDefaults

`func NewNodePreflightStatusWithDefaults() *NodePreflightStatus`

NewNodePreflightStatusWithDefaults instantiates a new NodePreflightStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *NodePreflightStatus) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *NodePreflightStatus) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *NodePreflightStatus) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetReason

`func (o *NodePreflightStatus) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *NodePreflightStatus) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *NodePreflightStatus) SetReason(v string)`

SetReason sets Reason field to given value.


### GetCode

`func (o *NodePreflightStatus) GetCode() int64`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *NodePreflightStatus) GetCodeOk() (*int64, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *NodePreflightStatus) SetCode(v int64)`

SetCode sets Code field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


