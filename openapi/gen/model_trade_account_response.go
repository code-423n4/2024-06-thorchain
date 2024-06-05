/*
Thornode API

Thornode REST API.

Contact: devs@thorchain.org
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// TradeAccountResponse struct for TradeAccountResponse
type TradeAccountResponse struct {
	// trade account asset with \"~\" separator
	Asset string `json:"asset"`
	// units of trade asset belonging to this owner
	Units string `json:"units"`
	// thor address of trade account owner
	Owner string `json:"owner"`
	// last thorchain height trade assets were added to trade account
	LastAddHeight *int64 `json:"last_add_height,omitempty"`
	// last thorchain height trade assets were withdrawn from trade account
	LastWithdrawHeight *int64 `json:"last_withdraw_height,omitempty"`
}

// NewTradeAccountResponse instantiates a new TradeAccountResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTradeAccountResponse(asset string, units string, owner string) *TradeAccountResponse {
	this := TradeAccountResponse{}
	this.Asset = asset
	this.Units = units
	this.Owner = owner
	return &this
}

// NewTradeAccountResponseWithDefaults instantiates a new TradeAccountResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTradeAccountResponseWithDefaults() *TradeAccountResponse {
	this := TradeAccountResponse{}
	return &this
}

// GetAsset returns the Asset field value
func (o *TradeAccountResponse) GetAsset() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Asset
}

// GetAssetOk returns a tuple with the Asset field value
// and a boolean to check if the value has been set.
func (o *TradeAccountResponse) GetAssetOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Asset, true
}

// SetAsset sets field value
func (o *TradeAccountResponse) SetAsset(v string) {
	o.Asset = v
}

// GetUnits returns the Units field value
func (o *TradeAccountResponse) GetUnits() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Units
}

// GetUnitsOk returns a tuple with the Units field value
// and a boolean to check if the value has been set.
func (o *TradeAccountResponse) GetUnitsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Units, true
}

// SetUnits sets field value
func (o *TradeAccountResponse) SetUnits(v string) {
	o.Units = v
}

// GetOwner returns the Owner field value
func (o *TradeAccountResponse) GetOwner() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Owner
}

// GetOwnerOk returns a tuple with the Owner field value
// and a boolean to check if the value has been set.
func (o *TradeAccountResponse) GetOwnerOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Owner, true
}

// SetOwner sets field value
func (o *TradeAccountResponse) SetOwner(v string) {
	o.Owner = v
}

// GetLastAddHeight returns the LastAddHeight field value if set, zero value otherwise.
func (o *TradeAccountResponse) GetLastAddHeight() int64 {
	if o == nil || o.LastAddHeight == nil {
		var ret int64
		return ret
	}
	return *o.LastAddHeight
}

// GetLastAddHeightOk returns a tuple with the LastAddHeight field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TradeAccountResponse) GetLastAddHeightOk() (*int64, bool) {
	if o == nil || o.LastAddHeight == nil {
		return nil, false
	}
	return o.LastAddHeight, true
}

// HasLastAddHeight returns a boolean if a field has been set.
func (o *TradeAccountResponse) HasLastAddHeight() bool {
	if o != nil && o.LastAddHeight != nil {
		return true
	}

	return false
}

// SetLastAddHeight gets a reference to the given int64 and assigns it to the LastAddHeight field.
func (o *TradeAccountResponse) SetLastAddHeight(v int64) {
	o.LastAddHeight = &v
}

// GetLastWithdrawHeight returns the LastWithdrawHeight field value if set, zero value otherwise.
func (o *TradeAccountResponse) GetLastWithdrawHeight() int64 {
	if o == nil || o.LastWithdrawHeight == nil {
		var ret int64
		return ret
	}
	return *o.LastWithdrawHeight
}

// GetLastWithdrawHeightOk returns a tuple with the LastWithdrawHeight field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TradeAccountResponse) GetLastWithdrawHeightOk() (*int64, bool) {
	if o == nil || o.LastWithdrawHeight == nil {
		return nil, false
	}
	return o.LastWithdrawHeight, true
}

// HasLastWithdrawHeight returns a boolean if a field has been set.
func (o *TradeAccountResponse) HasLastWithdrawHeight() bool {
	if o != nil && o.LastWithdrawHeight != nil {
		return true
	}

	return false
}

// SetLastWithdrawHeight gets a reference to the given int64 and assigns it to the LastWithdrawHeight field.
func (o *TradeAccountResponse) SetLastWithdrawHeight(v int64) {
	o.LastWithdrawHeight = &v
}

func (o TradeAccountResponse) MarshalJSON_deprecated() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["asset"] = o.Asset
	}
	if true {
		toSerialize["units"] = o.Units
	}
	if true {
		toSerialize["owner"] = o.Owner
	}
	if o.LastAddHeight != nil {
		toSerialize["last_add_height"] = o.LastAddHeight
	}
	if o.LastWithdrawHeight != nil {
		toSerialize["last_withdraw_height"] = o.LastWithdrawHeight
	}
	return json.Marshal(toSerialize)
}

type NullableTradeAccountResponse struct {
	value *TradeAccountResponse
	isSet bool
}

func (v NullableTradeAccountResponse) Get() *TradeAccountResponse {
	return v.value
}

func (v *NullableTradeAccountResponse) Set(val *TradeAccountResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableTradeAccountResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableTradeAccountResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTradeAccountResponse(val *TradeAccountResponse) *NullableTradeAccountResponse {
	return &NullableTradeAccountResponse{value: val, isSet: true}
}

func (v NullableTradeAccountResponse) MarshalJSON_deprecated() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTradeAccountResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


