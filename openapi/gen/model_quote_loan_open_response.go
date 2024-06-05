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

// QuoteLoanOpenResponse struct for QuoteLoanOpenResponse
type QuoteLoanOpenResponse struct {
	// the inbound address for the transaction on the source chain
	InboundAddress *string `json:"inbound_address,omitempty"`
	// the approximate number of source chain blocks required before processing
	InboundConfirmationBlocks *int64 `json:"inbound_confirmation_blocks,omitempty"`
	// the approximate seconds for block confirmations required before processing
	InboundConfirmationSeconds *int64 `json:"inbound_confirmation_seconds,omitempty"`
	// the number of thorchain blocks the outbound will be delayed
	OutboundDelayBlocks int64 `json:"outbound_delay_blocks"`
	// the approximate seconds for the outbound delay before it will be sent
	OutboundDelaySeconds int64 `json:"outbound_delay_seconds"`
	Fees QuoteFees `json:"fees"`
	// Deprecated - migrate to fees object.
	SlippageBps *int64 `json:"slippage_bps,omitempty"`
	// Deprecated - migrate to fees object.
	StreamingSlippageBps *int64 `json:"streaming_slippage_bps,omitempty"`
	// the EVM chain router contract address
	Router *string `json:"router,omitempty"`
	// expiration timestamp in unix seconds
	Expiry int64 `json:"expiry"`
	// static warning message
	Warning string `json:"warning"`
	// chain specific quote notes
	Notes string `json:"notes"`
	// Defines the minimum transaction size for the chain in base units (sats, wei, uatom). Transactions with asset amounts lower than the dust_threshold are ignored.
	DustThreshold *string `json:"dust_threshold,omitempty"`
	// The recommended minimum inbound amount for this transaction type & inbound asset. Sending less than this amount could result in failed refunds.
	RecommendedMinAmountIn *string `json:"recommended_min_amount_in,omitempty"`
	// the recommended gas rate to use for the inbound to ensure timely confirmation
	RecommendedGasRate string `json:"recommended_gas_rate"`
	// the units of the recommended gas rate
	GasRateUnits string `json:"gas_rate_units"`
	// generated memo for the loan open
	Memo *string `json:"memo,omitempty"`
	// the amount of the target asset the user can expect to receive after fees in 1e8 decimals
	ExpectedAmountOut string `json:"expected_amount_out"`
	// the expected collateralization ratio in basis points
	ExpectedCollateralizationRatio string `json:"expected_collateralization_ratio"`
	// the expected amount of collateral increase on the loan
	ExpectedCollateralDeposited string `json:"expected_collateral_deposited"`
	// the expected amount of TOR debt increase on the loan
	ExpectedDebtIssued string `json:"expected_debt_issued"`
	// The number of blocks involved in the streaming swaps during the open loan process.
	StreamingSwapBlocks int64 `json:"streaming_swap_blocks"`
	// The approximate number of seconds taken by the streaming swaps involved in the open loan process.
	StreamingSwapSeconds int64 `json:"streaming_swap_seconds"`
	// The total expected duration for a open loan, measured in seconds, which includes the time for inbound confirmation, the duration of streaming swaps, and any outbound delays.
	TotalOpenLoanSeconds int64 `json:"total_open_loan_seconds"`
}

// NewQuoteLoanOpenResponse instantiates a new QuoteLoanOpenResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewQuoteLoanOpenResponse(outboundDelayBlocks int64, outboundDelaySeconds int64, fees QuoteFees, expiry int64, warning string, notes string, recommendedGasRate string, gasRateUnits string, expectedAmountOut string, expectedCollateralizationRatio string, expectedCollateralDeposited string, expectedDebtIssued string, streamingSwapBlocks int64, streamingSwapSeconds int64, totalOpenLoanSeconds int64) *QuoteLoanOpenResponse {
	this := QuoteLoanOpenResponse{}
	this.OutboundDelayBlocks = outboundDelayBlocks
	this.OutboundDelaySeconds = outboundDelaySeconds
	this.Fees = fees
	this.Expiry = expiry
	this.Warning = warning
	this.Notes = notes
	this.RecommendedGasRate = recommendedGasRate
	this.GasRateUnits = gasRateUnits
	this.ExpectedAmountOut = expectedAmountOut
	this.ExpectedCollateralizationRatio = expectedCollateralizationRatio
	this.ExpectedCollateralDeposited = expectedCollateralDeposited
	this.ExpectedDebtIssued = expectedDebtIssued
	this.StreamingSwapBlocks = streamingSwapBlocks
	this.StreamingSwapSeconds = streamingSwapSeconds
	this.TotalOpenLoanSeconds = totalOpenLoanSeconds
	return &this
}

// NewQuoteLoanOpenResponseWithDefaults instantiates a new QuoteLoanOpenResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewQuoteLoanOpenResponseWithDefaults() *QuoteLoanOpenResponse {
	this := QuoteLoanOpenResponse{}
	return &this
}

// GetInboundAddress returns the InboundAddress field value if set, zero value otherwise.
func (o *QuoteLoanOpenResponse) GetInboundAddress() string {
	if o == nil || o.InboundAddress == nil {
		var ret string
		return ret
	}
	return *o.InboundAddress
}

// GetInboundAddressOk returns a tuple with the InboundAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetInboundAddressOk() (*string, bool) {
	if o == nil || o.InboundAddress == nil {
		return nil, false
	}
	return o.InboundAddress, true
}

// HasInboundAddress returns a boolean if a field has been set.
func (o *QuoteLoanOpenResponse) HasInboundAddress() bool {
	if o != nil && o.InboundAddress != nil {
		return true
	}

	return false
}

// SetInboundAddress gets a reference to the given string and assigns it to the InboundAddress field.
func (o *QuoteLoanOpenResponse) SetInboundAddress(v string) {
	o.InboundAddress = &v
}

// GetInboundConfirmationBlocks returns the InboundConfirmationBlocks field value if set, zero value otherwise.
func (o *QuoteLoanOpenResponse) GetInboundConfirmationBlocks() int64 {
	if o == nil || o.InboundConfirmationBlocks == nil {
		var ret int64
		return ret
	}
	return *o.InboundConfirmationBlocks
}

// GetInboundConfirmationBlocksOk returns a tuple with the InboundConfirmationBlocks field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetInboundConfirmationBlocksOk() (*int64, bool) {
	if o == nil || o.InboundConfirmationBlocks == nil {
		return nil, false
	}
	return o.InboundConfirmationBlocks, true
}

// HasInboundConfirmationBlocks returns a boolean if a field has been set.
func (o *QuoteLoanOpenResponse) HasInboundConfirmationBlocks() bool {
	if o != nil && o.InboundConfirmationBlocks != nil {
		return true
	}

	return false
}

// SetInboundConfirmationBlocks gets a reference to the given int64 and assigns it to the InboundConfirmationBlocks field.
func (o *QuoteLoanOpenResponse) SetInboundConfirmationBlocks(v int64) {
	o.InboundConfirmationBlocks = &v
}

// GetInboundConfirmationSeconds returns the InboundConfirmationSeconds field value if set, zero value otherwise.
func (o *QuoteLoanOpenResponse) GetInboundConfirmationSeconds() int64 {
	if o == nil || o.InboundConfirmationSeconds == nil {
		var ret int64
		return ret
	}
	return *o.InboundConfirmationSeconds
}

// GetInboundConfirmationSecondsOk returns a tuple with the InboundConfirmationSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetInboundConfirmationSecondsOk() (*int64, bool) {
	if o == nil || o.InboundConfirmationSeconds == nil {
		return nil, false
	}
	return o.InboundConfirmationSeconds, true
}

// HasInboundConfirmationSeconds returns a boolean if a field has been set.
func (o *QuoteLoanOpenResponse) HasInboundConfirmationSeconds() bool {
	if o != nil && o.InboundConfirmationSeconds != nil {
		return true
	}

	return false
}

// SetInboundConfirmationSeconds gets a reference to the given int64 and assigns it to the InboundConfirmationSeconds field.
func (o *QuoteLoanOpenResponse) SetInboundConfirmationSeconds(v int64) {
	o.InboundConfirmationSeconds = &v
}

// GetOutboundDelayBlocks returns the OutboundDelayBlocks field value
func (o *QuoteLoanOpenResponse) GetOutboundDelayBlocks() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.OutboundDelayBlocks
}

// GetOutboundDelayBlocksOk returns a tuple with the OutboundDelayBlocks field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetOutboundDelayBlocksOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutboundDelayBlocks, true
}

// SetOutboundDelayBlocks sets field value
func (o *QuoteLoanOpenResponse) SetOutboundDelayBlocks(v int64) {
	o.OutboundDelayBlocks = v
}

// GetOutboundDelaySeconds returns the OutboundDelaySeconds field value
func (o *QuoteLoanOpenResponse) GetOutboundDelaySeconds() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.OutboundDelaySeconds
}

// GetOutboundDelaySecondsOk returns a tuple with the OutboundDelaySeconds field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetOutboundDelaySecondsOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OutboundDelaySeconds, true
}

// SetOutboundDelaySeconds sets field value
func (o *QuoteLoanOpenResponse) SetOutboundDelaySeconds(v int64) {
	o.OutboundDelaySeconds = v
}

// GetFees returns the Fees field value
func (o *QuoteLoanOpenResponse) GetFees() QuoteFees {
	if o == nil {
		var ret QuoteFees
		return ret
	}

	return o.Fees
}

// GetFeesOk returns a tuple with the Fees field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetFeesOk() (*QuoteFees, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Fees, true
}

// SetFees sets field value
func (o *QuoteLoanOpenResponse) SetFees(v QuoteFees) {
	o.Fees = v
}

// GetSlippageBps returns the SlippageBps field value if set, zero value otherwise.
func (o *QuoteLoanOpenResponse) GetSlippageBps() int64 {
	if o == nil || o.SlippageBps == nil {
		var ret int64
		return ret
	}
	return *o.SlippageBps
}

// GetSlippageBpsOk returns a tuple with the SlippageBps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetSlippageBpsOk() (*int64, bool) {
	if o == nil || o.SlippageBps == nil {
		return nil, false
	}
	return o.SlippageBps, true
}

// HasSlippageBps returns a boolean if a field has been set.
func (o *QuoteLoanOpenResponse) HasSlippageBps() bool {
	if o != nil && o.SlippageBps != nil {
		return true
	}

	return false
}

// SetSlippageBps gets a reference to the given int64 and assigns it to the SlippageBps field.
func (o *QuoteLoanOpenResponse) SetSlippageBps(v int64) {
	o.SlippageBps = &v
}

// GetStreamingSlippageBps returns the StreamingSlippageBps field value if set, zero value otherwise.
func (o *QuoteLoanOpenResponse) GetStreamingSlippageBps() int64 {
	if o == nil || o.StreamingSlippageBps == nil {
		var ret int64
		return ret
	}
	return *o.StreamingSlippageBps
}

// GetStreamingSlippageBpsOk returns a tuple with the StreamingSlippageBps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetStreamingSlippageBpsOk() (*int64, bool) {
	if o == nil || o.StreamingSlippageBps == nil {
		return nil, false
	}
	return o.StreamingSlippageBps, true
}

// HasStreamingSlippageBps returns a boolean if a field has been set.
func (o *QuoteLoanOpenResponse) HasStreamingSlippageBps() bool {
	if o != nil && o.StreamingSlippageBps != nil {
		return true
	}

	return false
}

// SetStreamingSlippageBps gets a reference to the given int64 and assigns it to the StreamingSlippageBps field.
func (o *QuoteLoanOpenResponse) SetStreamingSlippageBps(v int64) {
	o.StreamingSlippageBps = &v
}

// GetRouter returns the Router field value if set, zero value otherwise.
func (o *QuoteLoanOpenResponse) GetRouter() string {
	if o == nil || o.Router == nil {
		var ret string
		return ret
	}
	return *o.Router
}

// GetRouterOk returns a tuple with the Router field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetRouterOk() (*string, bool) {
	if o == nil || o.Router == nil {
		return nil, false
	}
	return o.Router, true
}

// HasRouter returns a boolean if a field has been set.
func (o *QuoteLoanOpenResponse) HasRouter() bool {
	if o != nil && o.Router != nil {
		return true
	}

	return false
}

// SetRouter gets a reference to the given string and assigns it to the Router field.
func (o *QuoteLoanOpenResponse) SetRouter(v string) {
	o.Router = &v
}

// GetExpiry returns the Expiry field value
func (o *QuoteLoanOpenResponse) GetExpiry() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Expiry
}

// GetExpiryOk returns a tuple with the Expiry field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetExpiryOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Expiry, true
}

// SetExpiry sets field value
func (o *QuoteLoanOpenResponse) SetExpiry(v int64) {
	o.Expiry = v
}

// GetWarning returns the Warning field value
func (o *QuoteLoanOpenResponse) GetWarning() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Warning
}

// GetWarningOk returns a tuple with the Warning field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetWarningOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Warning, true
}

// SetWarning sets field value
func (o *QuoteLoanOpenResponse) SetWarning(v string) {
	o.Warning = v
}

// GetNotes returns the Notes field value
func (o *QuoteLoanOpenResponse) GetNotes() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Notes
}

// GetNotesOk returns a tuple with the Notes field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetNotesOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Notes, true
}

// SetNotes sets field value
func (o *QuoteLoanOpenResponse) SetNotes(v string) {
	o.Notes = v
}

// GetDustThreshold returns the DustThreshold field value if set, zero value otherwise.
func (o *QuoteLoanOpenResponse) GetDustThreshold() string {
	if o == nil || o.DustThreshold == nil {
		var ret string
		return ret
	}
	return *o.DustThreshold
}

// GetDustThresholdOk returns a tuple with the DustThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetDustThresholdOk() (*string, bool) {
	if o == nil || o.DustThreshold == nil {
		return nil, false
	}
	return o.DustThreshold, true
}

// HasDustThreshold returns a boolean if a field has been set.
func (o *QuoteLoanOpenResponse) HasDustThreshold() bool {
	if o != nil && o.DustThreshold != nil {
		return true
	}

	return false
}

// SetDustThreshold gets a reference to the given string and assigns it to the DustThreshold field.
func (o *QuoteLoanOpenResponse) SetDustThreshold(v string) {
	o.DustThreshold = &v
}

// GetRecommendedMinAmountIn returns the RecommendedMinAmountIn field value if set, zero value otherwise.
func (o *QuoteLoanOpenResponse) GetRecommendedMinAmountIn() string {
	if o == nil || o.RecommendedMinAmountIn == nil {
		var ret string
		return ret
	}
	return *o.RecommendedMinAmountIn
}

// GetRecommendedMinAmountInOk returns a tuple with the RecommendedMinAmountIn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetRecommendedMinAmountInOk() (*string, bool) {
	if o == nil || o.RecommendedMinAmountIn == nil {
		return nil, false
	}
	return o.RecommendedMinAmountIn, true
}

// HasRecommendedMinAmountIn returns a boolean if a field has been set.
func (o *QuoteLoanOpenResponse) HasRecommendedMinAmountIn() bool {
	if o != nil && o.RecommendedMinAmountIn != nil {
		return true
	}

	return false
}

// SetRecommendedMinAmountIn gets a reference to the given string and assigns it to the RecommendedMinAmountIn field.
func (o *QuoteLoanOpenResponse) SetRecommendedMinAmountIn(v string) {
	o.RecommendedMinAmountIn = &v
}

// GetRecommendedGasRate returns the RecommendedGasRate field value
func (o *QuoteLoanOpenResponse) GetRecommendedGasRate() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RecommendedGasRate
}

// GetRecommendedGasRateOk returns a tuple with the RecommendedGasRate field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetRecommendedGasRateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RecommendedGasRate, true
}

// SetRecommendedGasRate sets field value
func (o *QuoteLoanOpenResponse) SetRecommendedGasRate(v string) {
	o.RecommendedGasRate = v
}

// GetGasRateUnits returns the GasRateUnits field value
func (o *QuoteLoanOpenResponse) GetGasRateUnits() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GasRateUnits
}

// GetGasRateUnitsOk returns a tuple with the GasRateUnits field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetGasRateUnitsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GasRateUnits, true
}

// SetGasRateUnits sets field value
func (o *QuoteLoanOpenResponse) SetGasRateUnits(v string) {
	o.GasRateUnits = v
}

// GetMemo returns the Memo field value if set, zero value otherwise.
func (o *QuoteLoanOpenResponse) GetMemo() string {
	if o == nil || o.Memo == nil {
		var ret string
		return ret
	}
	return *o.Memo
}

// GetMemoOk returns a tuple with the Memo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetMemoOk() (*string, bool) {
	if o == nil || o.Memo == nil {
		return nil, false
	}
	return o.Memo, true
}

// HasMemo returns a boolean if a field has been set.
func (o *QuoteLoanOpenResponse) HasMemo() bool {
	if o != nil && o.Memo != nil {
		return true
	}

	return false
}

// SetMemo gets a reference to the given string and assigns it to the Memo field.
func (o *QuoteLoanOpenResponse) SetMemo(v string) {
	o.Memo = &v
}

// GetExpectedAmountOut returns the ExpectedAmountOut field value
func (o *QuoteLoanOpenResponse) GetExpectedAmountOut() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ExpectedAmountOut
}

// GetExpectedAmountOutOk returns a tuple with the ExpectedAmountOut field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetExpectedAmountOutOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExpectedAmountOut, true
}

// SetExpectedAmountOut sets field value
func (o *QuoteLoanOpenResponse) SetExpectedAmountOut(v string) {
	o.ExpectedAmountOut = v
}

// GetExpectedCollateralizationRatio returns the ExpectedCollateralizationRatio field value
func (o *QuoteLoanOpenResponse) GetExpectedCollateralizationRatio() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ExpectedCollateralizationRatio
}

// GetExpectedCollateralizationRatioOk returns a tuple with the ExpectedCollateralizationRatio field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetExpectedCollateralizationRatioOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExpectedCollateralizationRatio, true
}

// SetExpectedCollateralizationRatio sets field value
func (o *QuoteLoanOpenResponse) SetExpectedCollateralizationRatio(v string) {
	o.ExpectedCollateralizationRatio = v
}

// GetExpectedCollateralDeposited returns the ExpectedCollateralDeposited field value
func (o *QuoteLoanOpenResponse) GetExpectedCollateralDeposited() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ExpectedCollateralDeposited
}

// GetExpectedCollateralDepositedOk returns a tuple with the ExpectedCollateralDeposited field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetExpectedCollateralDepositedOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExpectedCollateralDeposited, true
}

// SetExpectedCollateralDeposited sets field value
func (o *QuoteLoanOpenResponse) SetExpectedCollateralDeposited(v string) {
	o.ExpectedCollateralDeposited = v
}

// GetExpectedDebtIssued returns the ExpectedDebtIssued field value
func (o *QuoteLoanOpenResponse) GetExpectedDebtIssued() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ExpectedDebtIssued
}

// GetExpectedDebtIssuedOk returns a tuple with the ExpectedDebtIssued field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetExpectedDebtIssuedOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExpectedDebtIssued, true
}

// SetExpectedDebtIssued sets field value
func (o *QuoteLoanOpenResponse) SetExpectedDebtIssued(v string) {
	o.ExpectedDebtIssued = v
}

// GetStreamingSwapBlocks returns the StreamingSwapBlocks field value
func (o *QuoteLoanOpenResponse) GetStreamingSwapBlocks() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.StreamingSwapBlocks
}

// GetStreamingSwapBlocksOk returns a tuple with the StreamingSwapBlocks field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetStreamingSwapBlocksOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.StreamingSwapBlocks, true
}

// SetStreamingSwapBlocks sets field value
func (o *QuoteLoanOpenResponse) SetStreamingSwapBlocks(v int64) {
	o.StreamingSwapBlocks = v
}

// GetStreamingSwapSeconds returns the StreamingSwapSeconds field value
func (o *QuoteLoanOpenResponse) GetStreamingSwapSeconds() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.StreamingSwapSeconds
}

// GetStreamingSwapSecondsOk returns a tuple with the StreamingSwapSeconds field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetStreamingSwapSecondsOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.StreamingSwapSeconds, true
}

// SetStreamingSwapSeconds sets field value
func (o *QuoteLoanOpenResponse) SetStreamingSwapSeconds(v int64) {
	o.StreamingSwapSeconds = v
}

// GetTotalOpenLoanSeconds returns the TotalOpenLoanSeconds field value
func (o *QuoteLoanOpenResponse) GetTotalOpenLoanSeconds() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.TotalOpenLoanSeconds
}

// GetTotalOpenLoanSecondsOk returns a tuple with the TotalOpenLoanSeconds field value
// and a boolean to check if the value has been set.
func (o *QuoteLoanOpenResponse) GetTotalOpenLoanSecondsOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TotalOpenLoanSeconds, true
}

// SetTotalOpenLoanSeconds sets field value
func (o *QuoteLoanOpenResponse) SetTotalOpenLoanSeconds(v int64) {
	o.TotalOpenLoanSeconds = v
}

func (o QuoteLoanOpenResponse) MarshalJSON_deprecated() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.InboundAddress != nil {
		toSerialize["inbound_address"] = o.InboundAddress
	}
	if o.InboundConfirmationBlocks != nil {
		toSerialize["inbound_confirmation_blocks"] = o.InboundConfirmationBlocks
	}
	if o.InboundConfirmationSeconds != nil {
		toSerialize["inbound_confirmation_seconds"] = o.InboundConfirmationSeconds
	}
	if true {
		toSerialize["outbound_delay_blocks"] = o.OutboundDelayBlocks
	}
	if true {
		toSerialize["outbound_delay_seconds"] = o.OutboundDelaySeconds
	}
	if true {
		toSerialize["fees"] = o.Fees
	}
	if o.SlippageBps != nil {
		toSerialize["slippage_bps"] = o.SlippageBps
	}
	if o.StreamingSlippageBps != nil {
		toSerialize["streaming_slippage_bps"] = o.StreamingSlippageBps
	}
	if o.Router != nil {
		toSerialize["router"] = o.Router
	}
	if true {
		toSerialize["expiry"] = o.Expiry
	}
	if true {
		toSerialize["warning"] = o.Warning
	}
	if true {
		toSerialize["notes"] = o.Notes
	}
	if o.DustThreshold != nil {
		toSerialize["dust_threshold"] = o.DustThreshold
	}
	if o.RecommendedMinAmountIn != nil {
		toSerialize["recommended_min_amount_in"] = o.RecommendedMinAmountIn
	}
	if true {
		toSerialize["recommended_gas_rate"] = o.RecommendedGasRate
	}
	if true {
		toSerialize["gas_rate_units"] = o.GasRateUnits
	}
	if o.Memo != nil {
		toSerialize["memo"] = o.Memo
	}
	if true {
		toSerialize["expected_amount_out"] = o.ExpectedAmountOut
	}
	if true {
		toSerialize["expected_collateralization_ratio"] = o.ExpectedCollateralizationRatio
	}
	if true {
		toSerialize["expected_collateral_deposited"] = o.ExpectedCollateralDeposited
	}
	if true {
		toSerialize["expected_debt_issued"] = o.ExpectedDebtIssued
	}
	if true {
		toSerialize["streaming_swap_blocks"] = o.StreamingSwapBlocks
	}
	if true {
		toSerialize["streaming_swap_seconds"] = o.StreamingSwapSeconds
	}
	if true {
		toSerialize["total_open_loan_seconds"] = o.TotalOpenLoanSeconds
	}
	return json.Marshal(toSerialize)
}

type NullableQuoteLoanOpenResponse struct {
	value *QuoteLoanOpenResponse
	isSet bool
}

func (v NullableQuoteLoanOpenResponse) Get() *QuoteLoanOpenResponse {
	return v.value
}

func (v *NullableQuoteLoanOpenResponse) Set(val *QuoteLoanOpenResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableQuoteLoanOpenResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableQuoteLoanOpenResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableQuoteLoanOpenResponse(val *QuoteLoanOpenResponse) *NullableQuoteLoanOpenResponse {
	return &NullableQuoteLoanOpenResponse{value: val, isSet: true}
}

func (v NullableQuoteLoanOpenResponse) MarshalJSON_deprecated() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableQuoteLoanOpenResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

