# BlockResponseId

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Hash** | **string** |  | 
**Parts** | [**BlockResponseIdParts**](BlockResponseIdParts.md) |  | 

## Methods

### NewBlockResponseId

`func NewBlockResponseId(hash string, parts BlockResponseIdParts, ) *BlockResponseId`

NewBlockResponseId instantiates a new BlockResponseId object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockResponseIdWithDefaults

`func NewBlockResponseIdWithDefaults() *BlockResponseId`

NewBlockResponseIdWithDefaults instantiates a new BlockResponseId object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHash

`func (o *BlockResponseId) GetHash() string`

GetHash returns the Hash field if non-nil, zero value otherwise.

### GetHashOk

`func (o *BlockResponseId) GetHashOk() (*string, bool)`

GetHashOk returns a tuple with the Hash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHash

`func (o *BlockResponseId) SetHash(v string)`

SetHash sets Hash field to given value.


### GetParts

`func (o *BlockResponseId) GetParts() BlockResponseIdParts`

GetParts returns the Parts field if non-nil, zero value otherwise.

### GetPartsOk

`func (o *BlockResponseId) GetPartsOk() (*BlockResponseIdParts, bool)`

GetPartsOk returns a tuple with the Parts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParts

`func (o *BlockResponseId) SetParts(v BlockResponseIdParts)`

SetParts sets Parts field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


