# BlockTx

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Hash** | **string** |  | 
**Tx** | **map[string]interface{}** |  | 
**Result** | [**BlockTxResult**](BlockTxResult.md) |  | 

## Methods

### NewBlockTx

`func NewBlockTx(hash string, tx map[string]interface{}, result BlockTxResult, ) *BlockTx`

NewBlockTx instantiates a new BlockTx object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockTxWithDefaults

`func NewBlockTxWithDefaults() *BlockTx`

NewBlockTxWithDefaults instantiates a new BlockTx object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHash

`func (o *BlockTx) GetHash() string`

GetHash returns the Hash field if non-nil, zero value otherwise.

### GetHashOk

`func (o *BlockTx) GetHashOk() (*string, bool)`

GetHashOk returns a tuple with the Hash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHash

`func (o *BlockTx) SetHash(v string)`

SetHash sets Hash field to given value.


### GetTx

`func (o *BlockTx) GetTx() map[string]interface{}`

GetTx returns the Tx field if non-nil, zero value otherwise.

### GetTxOk

`func (o *BlockTx) GetTxOk() (*map[string]interface{}, bool)`

GetTxOk returns a tuple with the Tx field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTx

`func (o *BlockTx) SetTx(v map[string]interface{})`

SetTx sets Tx field to given value.


### GetResult

`func (o *BlockTx) GetResult() BlockTxResult`

GetResult returns the Result field if non-nil, zero value otherwise.

### GetResultOk

`func (o *BlockTx) GetResultOk() (*BlockTxResult, bool)`

GetResultOk returns a tuple with the Result field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResult

`func (o *BlockTx) SetResult(v BlockTxResult)`

SetResult sets Result field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


