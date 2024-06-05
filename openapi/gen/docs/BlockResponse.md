# BlockResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**BlockResponseId**](BlockResponseId.md) |  | 
**Header** | [**BlockResponseHeader**](BlockResponseHeader.md) |  | 
**BeginBlockEvents** | **[]map[string]string** |  | 
**EndBlockEvents** | **[]map[string]string** |  | 
**Txs** | [**[]BlockTx**](BlockTx.md) |  | 

## Methods

### NewBlockResponse

`func NewBlockResponse(id BlockResponseId, header BlockResponseHeader, beginBlockEvents []map[string]string, endBlockEvents []map[string]string, txs []BlockTx, ) *BlockResponse`

NewBlockResponse instantiates a new BlockResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockResponseWithDefaults

`func NewBlockResponseWithDefaults() *BlockResponse`

NewBlockResponseWithDefaults instantiates a new BlockResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *BlockResponse) GetId() BlockResponseId`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *BlockResponse) GetIdOk() (*BlockResponseId, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *BlockResponse) SetId(v BlockResponseId)`

SetId sets Id field to given value.


### GetHeader

`func (o *BlockResponse) GetHeader() BlockResponseHeader`

GetHeader returns the Header field if non-nil, zero value otherwise.

### GetHeaderOk

`func (o *BlockResponse) GetHeaderOk() (*BlockResponseHeader, bool)`

GetHeaderOk returns a tuple with the Header field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeader

`func (o *BlockResponse) SetHeader(v BlockResponseHeader)`

SetHeader sets Header field to given value.


### GetBeginBlockEvents

`func (o *BlockResponse) GetBeginBlockEvents() []map[string]string`

GetBeginBlockEvents returns the BeginBlockEvents field if non-nil, zero value otherwise.

### GetBeginBlockEventsOk

`func (o *BlockResponse) GetBeginBlockEventsOk() (*[]map[string]string, bool)`

GetBeginBlockEventsOk returns a tuple with the BeginBlockEvents field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBeginBlockEvents

`func (o *BlockResponse) SetBeginBlockEvents(v []map[string]string)`

SetBeginBlockEvents sets BeginBlockEvents field to given value.


### GetEndBlockEvents

`func (o *BlockResponse) GetEndBlockEvents() []map[string]string`

GetEndBlockEvents returns the EndBlockEvents field if non-nil, zero value otherwise.

### GetEndBlockEventsOk

`func (o *BlockResponse) GetEndBlockEventsOk() (*[]map[string]string, bool)`

GetEndBlockEventsOk returns a tuple with the EndBlockEvents field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndBlockEvents

`func (o *BlockResponse) SetEndBlockEvents(v []map[string]string)`

SetEndBlockEvents sets EndBlockEvents field to given value.


### GetTxs

`func (o *BlockResponse) GetTxs() []BlockTx`

GetTxs returns the Txs field if non-nil, zero value otherwise.

### GetTxsOk

`func (o *BlockResponse) GetTxsOk() (*[]BlockTx, bool)`

GetTxsOk returns a tuple with the Txs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTxs

`func (o *BlockResponse) SetTxs(v []BlockTx)`

SetTxs sets Txs field to given value.


### SetTxsNil

`func (o *BlockResponse) SetTxsNil(b bool)`

 SetTxsNil sets the value for Txs to be an explicit nil

### UnsetTxs
`func (o *BlockResponse) UnsetTxs()`

UnsetTxs ensures that no value is present for Txs, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


