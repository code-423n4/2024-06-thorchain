# BlockTxResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **int64** |  | [optional] 
**Data** | Pointer to **string** |  | [optional] 
**Log** | Pointer to **string** |  | [optional] 
**Info** | Pointer to **string** |  | [optional] 
**GasWanted** | Pointer to **string** |  | [optional] 
**GasUsed** | Pointer to **string** |  | [optional] 
**Events** | Pointer to **[]map[string]string** |  | [optional] 
**Codespace** | Pointer to **string** |  | [optional] 

## Methods

### NewBlockTxResult

`func NewBlockTxResult() *BlockTxResult`

NewBlockTxResult instantiates a new BlockTxResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockTxResultWithDefaults

`func NewBlockTxResultWithDefaults() *BlockTxResult`

NewBlockTxResultWithDefaults instantiates a new BlockTxResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *BlockTxResult) GetCode() int64`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *BlockTxResult) GetCodeOk() (*int64, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *BlockTxResult) SetCode(v int64)`

SetCode sets Code field to given value.

### HasCode

`func (o *BlockTxResult) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetData

`func (o *BlockTxResult) GetData() string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *BlockTxResult) GetDataOk() (*string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *BlockTxResult) SetData(v string)`

SetData sets Data field to given value.

### HasData

`func (o *BlockTxResult) HasData() bool`

HasData returns a boolean if a field has been set.

### GetLog

`func (o *BlockTxResult) GetLog() string`

GetLog returns the Log field if non-nil, zero value otherwise.

### GetLogOk

`func (o *BlockTxResult) GetLogOk() (*string, bool)`

GetLogOk returns a tuple with the Log field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLog

`func (o *BlockTxResult) SetLog(v string)`

SetLog sets Log field to given value.

### HasLog

`func (o *BlockTxResult) HasLog() bool`

HasLog returns a boolean if a field has been set.

### GetInfo

`func (o *BlockTxResult) GetInfo() string`

GetInfo returns the Info field if non-nil, zero value otherwise.

### GetInfoOk

`func (o *BlockTxResult) GetInfoOk() (*string, bool)`

GetInfoOk returns a tuple with the Info field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfo

`func (o *BlockTxResult) SetInfo(v string)`

SetInfo sets Info field to given value.

### HasInfo

`func (o *BlockTxResult) HasInfo() bool`

HasInfo returns a boolean if a field has been set.

### GetGasWanted

`func (o *BlockTxResult) GetGasWanted() string`

GetGasWanted returns the GasWanted field if non-nil, zero value otherwise.

### GetGasWantedOk

`func (o *BlockTxResult) GetGasWantedOk() (*string, bool)`

GetGasWantedOk returns a tuple with the GasWanted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasWanted

`func (o *BlockTxResult) SetGasWanted(v string)`

SetGasWanted sets GasWanted field to given value.

### HasGasWanted

`func (o *BlockTxResult) HasGasWanted() bool`

HasGasWanted returns a boolean if a field has been set.

### GetGasUsed

`func (o *BlockTxResult) GetGasUsed() string`

GetGasUsed returns the GasUsed field if non-nil, zero value otherwise.

### GetGasUsedOk

`func (o *BlockTxResult) GetGasUsedOk() (*string, bool)`

GetGasUsedOk returns a tuple with the GasUsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGasUsed

`func (o *BlockTxResult) SetGasUsed(v string)`

SetGasUsed sets GasUsed field to given value.

### HasGasUsed

`func (o *BlockTxResult) HasGasUsed() bool`

HasGasUsed returns a boolean if a field has been set.

### GetEvents

`func (o *BlockTxResult) GetEvents() []map[string]string`

GetEvents returns the Events field if non-nil, zero value otherwise.

### GetEventsOk

`func (o *BlockTxResult) GetEventsOk() (*[]map[string]string, bool)`

GetEventsOk returns a tuple with the Events field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvents

`func (o *BlockTxResult) SetEvents(v []map[string]string)`

SetEvents sets Events field to given value.

### HasEvents

`func (o *BlockTxResult) HasEvents() bool`

HasEvents returns a boolean if a field has been set.

### SetEventsNil

`func (o *BlockTxResult) SetEventsNil(b bool)`

 SetEventsNil sets the value for Events to be an explicit nil

### UnsetEvents
`func (o *BlockTxResult) UnsetEvents()`

UnsetEvents ensures that no value is present for Events, not even an explicit nil
### GetCodespace

`func (o *BlockTxResult) GetCodespace() string`

GetCodespace returns the Codespace field if non-nil, zero value otherwise.

### GetCodespaceOk

`func (o *BlockTxResult) GetCodespaceOk() (*string, bool)`

GetCodespaceOk returns a tuple with the Codespace field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodespace

`func (o *BlockTxResult) SetCodespace(v string)`

SetCodespace sets Codespace field to given value.

### HasCodespace

`func (o *BlockTxResult) HasCodespace() bool`

HasCodespace returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


