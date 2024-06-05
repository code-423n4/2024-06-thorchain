# Tx

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Chain** | Pointer to **string** |  | [optional] 
**FromAddress** | Pointer to **string** |  | [optional] 
**ToAddress** | Pointer to **string** |  | [optional] 
**Coins** | [**[]Coin**](Coin.md) |  | 
**Gas** | [**[]Coin**](Coin.md) |  | 
**Memo** | Pointer to **string** |  | [optional] 

## Methods

### NewTx

`func NewTx(coins []Coin, gas []Coin, ) *Tx`

NewTx instantiates a new Tx object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTxWithDefaults

`func NewTxWithDefaults() *Tx`

NewTxWithDefaults instantiates a new Tx object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Tx) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Tx) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Tx) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Tx) HasId() bool`

HasId returns a boolean if a field has been set.

### GetChain

`func (o *Tx) GetChain() string`

GetChain returns the Chain field if non-nil, zero value otherwise.

### GetChainOk

`func (o *Tx) GetChainOk() (*string, bool)`

GetChainOk returns a tuple with the Chain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChain

`func (o *Tx) SetChain(v string)`

SetChain sets Chain field to given value.

### HasChain

`func (o *Tx) HasChain() bool`

HasChain returns a boolean if a field has been set.

### GetFromAddress

`func (o *Tx) GetFromAddress() string`

GetFromAddress returns the FromAddress field if non-nil, zero value otherwise.

### GetFromAddressOk

`func (o *Tx) GetFromAddressOk() (*string, bool)`

GetFromAddressOk returns a tuple with the FromAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFromAddress

`func (o *Tx) SetFromAddress(v string)`

SetFromAddress sets FromAddress field to given value.

### HasFromAddress

`func (o *Tx) HasFromAddress() bool`

HasFromAddress returns a boolean if a field has been set.

### GetToAddress

`func (o *Tx) GetToAddress() string`

GetToAddress returns the ToAddress field if non-nil, zero value otherwise.

### GetToAddressOk

`func (o *Tx) GetToAddressOk() (*string, bool)`

GetToAddressOk returns a tuple with the ToAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToAddress

`func (o *Tx) SetToAddress(v string)`

SetToAddress sets ToAddress field to given value.

### HasToAddress

`func (o *Tx) HasToAddress() bool`

HasToAddress returns a boolean if a field has been set.

### GetCoins

`func (o *Tx) GetCoins() []Coin`

GetCoins returns the Coins field if non-nil, zero value otherwise.

### GetCoinsOk

`func (o *Tx) GetCoinsOk() (*[]Coin, bool)`

GetCoinsOk returns a tuple with the Coins field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCoins

`func (o *Tx) SetCoins(v []Coin)`

SetCoins sets Coins field to given value.


### GetGas

`func (o *Tx) GetGas() []Coin`

GetGas returns the Gas field if non-nil, zero value otherwise.

### GetGasOk

`func (o *Tx) GetGasOk() (*[]Coin, bool)`

GetGasOk returns a tuple with the Gas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGas

`func (o *Tx) SetGas(v []Coin)`

SetGas sets Gas field to given value.


### GetMemo

`func (o *Tx) GetMemo() string`

GetMemo returns the Memo field if non-nil, zero value otherwise.

### GetMemoOk

`func (o *Tx) GetMemoOk() (*string, bool)`

GetMemoOk returns a tuple with the Memo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemo

`func (o *Tx) SetMemo(v string)`

SetMemo sets Memo field to given value.

### HasMemo

`func (o *Tx) HasMemo() bool`

HasMemo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


