# InvariantResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Invariant** | **string** | The name of the invariant. | 
**Broken** | **bool** | Returns true if the invariant is broken. | 
**Msg** | **[]string** | Informative message about the invariant result. | 

## Methods

### NewInvariantResponse

`func NewInvariantResponse(invariant string, broken bool, msg []string, ) *InvariantResponse`

NewInvariantResponse instantiates a new InvariantResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInvariantResponseWithDefaults

`func NewInvariantResponseWithDefaults() *InvariantResponse`

NewInvariantResponseWithDefaults instantiates a new InvariantResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInvariant

`func (o *InvariantResponse) GetInvariant() string`

GetInvariant returns the Invariant field if non-nil, zero value otherwise.

### GetInvariantOk

`func (o *InvariantResponse) GetInvariantOk() (*string, bool)`

GetInvariantOk returns a tuple with the Invariant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInvariant

`func (o *InvariantResponse) SetInvariant(v string)`

SetInvariant sets Invariant field to given value.


### GetBroken

`func (o *InvariantResponse) GetBroken() bool`

GetBroken returns the Broken field if non-nil, zero value otherwise.

### GetBrokenOk

`func (o *InvariantResponse) GetBrokenOk() (*bool, bool)`

GetBrokenOk returns a tuple with the Broken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBroken

`func (o *InvariantResponse) SetBroken(v bool)`

SetBroken sets Broken field to given value.


### GetMsg

`func (o *InvariantResponse) GetMsg() []string`

GetMsg returns the Msg field if non-nil, zero value otherwise.

### GetMsgOk

`func (o *InvariantResponse) GetMsgOk() (*[]string, bool)`

GetMsgOk returns a tuple with the Msg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMsg

`func (o *InvariantResponse) SetMsg(v []string)`

SetMsg sets Msg field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


