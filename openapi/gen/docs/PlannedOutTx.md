# PlannedOutTx

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Chain** | **string** |  | 
**ToAddress** | **string** |  | 
**Coin** | [**Coin**](Coin.md) |  | 
**Refund** | **bool** | returns true if the planned transaction has a refund memo | 

## Methods

### NewPlannedOutTx

`func NewPlannedOutTx(chain string, toAddress string, coin Coin, refund bool, ) *PlannedOutTx`

NewPlannedOutTx instantiates a new PlannedOutTx object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlannedOutTxWithDefaults

`func NewPlannedOutTxWithDefaults() *PlannedOutTx`

NewPlannedOutTxWithDefaults instantiates a new PlannedOutTx object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChain

`func (o *PlannedOutTx) GetChain() string`

GetChain returns the Chain field if non-nil, zero value otherwise.

### GetChainOk

`func (o *PlannedOutTx) GetChainOk() (*string, bool)`

GetChainOk returns a tuple with the Chain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChain

`func (o *PlannedOutTx) SetChain(v string)`

SetChain sets Chain field to given value.


### GetToAddress

`func (o *PlannedOutTx) GetToAddress() string`

GetToAddress returns the ToAddress field if non-nil, zero value otherwise.

### GetToAddressOk

`func (o *PlannedOutTx) GetToAddressOk() (*string, bool)`

GetToAddressOk returns a tuple with the ToAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToAddress

`func (o *PlannedOutTx) SetToAddress(v string)`

SetToAddress sets ToAddress field to given value.


### GetCoin

`func (o *PlannedOutTx) GetCoin() Coin`

GetCoin returns the Coin field if non-nil, zero value otherwise.

### GetCoinOk

`func (o *PlannedOutTx) GetCoinOk() (*Coin, bool)`

GetCoinOk returns a tuple with the Coin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCoin

`func (o *PlannedOutTx) SetCoin(v Coin)`

SetCoin sets Coin field to given value.


### GetRefund

`func (o *PlannedOutTx) GetRefund() bool`

GetRefund returns the Refund field if non-nil, zero value otherwise.

### GetRefundOk

`func (o *PlannedOutTx) GetRefundOk() (*bool, bool)`

GetRefundOk returns a tuple with the Refund field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefund

`func (o *PlannedOutTx) SetRefund(v bool)`

SetRefund sets Refund field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


