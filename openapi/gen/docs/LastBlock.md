# LastBlock

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Chain** | **string** |  | 
**LastObservedIn** | **int64** |  | 
**LastSignedOut** | **int64** |  | 
**Thorchain** | **int64** |  | 

## Methods

### NewLastBlock

`func NewLastBlock(chain string, lastObservedIn int64, lastSignedOut int64, thorchain int64, ) *LastBlock`

NewLastBlock instantiates a new LastBlock object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLastBlockWithDefaults

`func NewLastBlockWithDefaults() *LastBlock`

NewLastBlockWithDefaults instantiates a new LastBlock object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChain

`func (o *LastBlock) GetChain() string`

GetChain returns the Chain field if non-nil, zero value otherwise.

### GetChainOk

`func (o *LastBlock) GetChainOk() (*string, bool)`

GetChainOk returns a tuple with the Chain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChain

`func (o *LastBlock) SetChain(v string)`

SetChain sets Chain field to given value.


### GetLastObservedIn

`func (o *LastBlock) GetLastObservedIn() int64`

GetLastObservedIn returns the LastObservedIn field if non-nil, zero value otherwise.

### GetLastObservedInOk

`func (o *LastBlock) GetLastObservedInOk() (*int64, bool)`

GetLastObservedInOk returns a tuple with the LastObservedIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastObservedIn

`func (o *LastBlock) SetLastObservedIn(v int64)`

SetLastObservedIn sets LastObservedIn field to given value.


### GetLastSignedOut

`func (o *LastBlock) GetLastSignedOut() int64`

GetLastSignedOut returns the LastSignedOut field if non-nil, zero value otherwise.

### GetLastSignedOutOk

`func (o *LastBlock) GetLastSignedOutOk() (*int64, bool)`

GetLastSignedOutOk returns a tuple with the LastSignedOut field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastSignedOut

`func (o *LastBlock) SetLastSignedOut(v int64)`

SetLastSignedOut sets LastSignedOut field to given value.


### GetThorchain

`func (o *LastBlock) GetThorchain() int64`

GetThorchain returns the Thorchain field if non-nil, zero value otherwise.

### GetThorchainOk

`func (o *LastBlock) GetThorchainOk() (*int64, bool)`

GetThorchainOk returns a tuple with the Thorchain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThorchain

`func (o *LastBlock) SetThorchain(v int64)`

SetThorchain sets Thorchain field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


