/*
Thornode API

Thornode REST API.

Contact: devs@thorchain.org
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)


// QuoteApiService QuoteApi service
type QuoteApiService service

type ApiQuoteloancloseRequest struct {
	ctx context.Context
	ApiService *QuoteApiService
	height *int64
	fromAsset *string
	repayBps *int64
	toAsset *string
	loanOwner *string
	minOut *string
}

// optional block height, defaults to current tip
func (r ApiQuoteloancloseRequest) Height(height int64) ApiQuoteloancloseRequest {
	r.height = &height
	return r
}

// the asset used to repay the loan
func (r ApiQuoteloancloseRequest) FromAsset(fromAsset string) ApiQuoteloancloseRequest {
	r.fromAsset = &fromAsset
	return r
}

// the basis points of the existing position to repay
func (r ApiQuoteloancloseRequest) RepayBps(repayBps int64) ApiQuoteloancloseRequest {
	r.repayBps = &repayBps
	return r
}

// the collateral asset of the loan
func (r ApiQuoteloancloseRequest) ToAsset(toAsset string) ApiQuoteloancloseRequest {
	r.toAsset = &toAsset
	return r
}

// the owner of the loan collateral
func (r ApiQuoteloancloseRequest) LoanOwner(loanOwner string) ApiQuoteloancloseRequest {
	r.loanOwner = &loanOwner
	return r
}

// the minimum amount of the target asset to accept
func (r ApiQuoteloancloseRequest) MinOut(minOut string) ApiQuoteloancloseRequest {
	r.minOut = &minOut
	return r
}

func (r ApiQuoteloancloseRequest) Execute() (*QuoteLoanCloseResponse, *http.Response, error) {
	return r.ApiService.QuoteloancloseExecute(r)
}

/*
Quoteloanclose Method for Quoteloanclose

Provide a quote estimate for the provided loan close.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiQuoteloancloseRequest
*/
func (a *QuoteApiService) Quoteloanclose(ctx context.Context) ApiQuoteloancloseRequest {
	return ApiQuoteloancloseRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return QuoteLoanCloseResponse
func (a *QuoteApiService) QuoteloancloseExecute(r ApiQuoteloancloseRequest) (*QuoteLoanCloseResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *QuoteLoanCloseResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QuoteApiService.Quoteloanclose")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/quote/loan/close"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
	}
	if r.fromAsset != nil {
		localVarQueryParams.Add("from_asset", parameterToString(*r.fromAsset, ""))
	}
	if r.repayBps != nil {
		localVarQueryParams.Add("repay_bps", parameterToString(*r.repayBps, ""))
	}
	if r.toAsset != nil {
		localVarQueryParams.Add("to_asset", parameterToString(*r.toAsset, ""))
	}
	if r.loanOwner != nil {
		localVarQueryParams.Add("loan_owner", parameterToString(*r.loanOwner, ""))
	}
	if r.minOut != nil {
		localVarQueryParams.Add("min_out", parameterToString(*r.minOut, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiQuoteloanopenRequest struct {
	ctx context.Context
	ApiService *QuoteApiService
	height *int64
	fromAsset *string
	amount *int64
	toAsset *string
	destination *string
	minOut *string
	affiliateBps *int64
	affiliate *string
}

// optional block height, defaults to current tip
func (r ApiQuoteloanopenRequest) Height(height int64) ApiQuoteloanopenRequest {
	r.height = &height
	return r
}

// the collateral asset
func (r ApiQuoteloanopenRequest) FromAsset(fromAsset string) ApiQuoteloanopenRequest {
	r.fromAsset = &fromAsset
	return r
}

// the collateral asset amount in 1e8 decimals
func (r ApiQuoteloanopenRequest) Amount(amount int64) ApiQuoteloanopenRequest {
	r.amount = &amount
	return r
}

// the target asset to receive (loan denominated in TOR regardless)
func (r ApiQuoteloanopenRequest) ToAsset(toAsset string) ApiQuoteloanopenRequest {
	r.toAsset = &toAsset
	return r
}

// the destination address, required to generate memo
func (r ApiQuoteloanopenRequest) Destination(destination string) ApiQuoteloanopenRequest {
	r.destination = &destination
	return r
}

// the minimum amount of the target asset to accept
func (r ApiQuoteloanopenRequest) MinOut(minOut string) ApiQuoteloanopenRequest {
	r.minOut = &minOut
	return r
}

// the affiliate fee in basis points
func (r ApiQuoteloanopenRequest) AffiliateBps(affiliateBps int64) ApiQuoteloanopenRequest {
	r.affiliateBps = &affiliateBps
	return r
}

// the affiliate (address or thorname)
func (r ApiQuoteloanopenRequest) Affiliate(affiliate string) ApiQuoteloanopenRequest {
	r.affiliate = &affiliate
	return r
}

func (r ApiQuoteloanopenRequest) Execute() (*QuoteLoanOpenResponse, *http.Response, error) {
	return r.ApiService.QuoteloanopenExecute(r)
}

/*
Quoteloanopen Method for Quoteloanopen

Provide a quote estimate for the provided loan open.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiQuoteloanopenRequest
*/
func (a *QuoteApiService) Quoteloanopen(ctx context.Context) ApiQuoteloanopenRequest {
	return ApiQuoteloanopenRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return QuoteLoanOpenResponse
func (a *QuoteApiService) QuoteloanopenExecute(r ApiQuoteloanopenRequest) (*QuoteLoanOpenResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *QuoteLoanOpenResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QuoteApiService.Quoteloanopen")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/quote/loan/open"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
	}
	if r.fromAsset != nil {
		localVarQueryParams.Add("from_asset", parameterToString(*r.fromAsset, ""))
	}
	if r.amount != nil {
		localVarQueryParams.Add("amount", parameterToString(*r.amount, ""))
	}
	if r.toAsset != nil {
		localVarQueryParams.Add("to_asset", parameterToString(*r.toAsset, ""))
	}
	if r.destination != nil {
		localVarQueryParams.Add("destination", parameterToString(*r.destination, ""))
	}
	if r.minOut != nil {
		localVarQueryParams.Add("min_out", parameterToString(*r.minOut, ""))
	}
	if r.affiliateBps != nil {
		localVarQueryParams.Add("affiliate_bps", parameterToString(*r.affiliateBps, ""))
	}
	if r.affiliate != nil {
		localVarQueryParams.Add("affiliate", parameterToString(*r.affiliate, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiQuotesaverdepositRequest struct {
	ctx context.Context
	ApiService *QuoteApiService
	height *int64
	asset *string
	amount *int64
}

// optional block height, defaults to current tip
func (r ApiQuotesaverdepositRequest) Height(height int64) ApiQuotesaverdepositRequest {
	r.height = &height
	return r
}

// the asset to deposit
func (r ApiQuotesaverdepositRequest) Asset(asset string) ApiQuotesaverdepositRequest {
	r.asset = &asset
	return r
}

// the source asset amount in 1e8 decimals
func (r ApiQuotesaverdepositRequest) Amount(amount int64) ApiQuotesaverdepositRequest {
	r.amount = &amount
	return r
}

func (r ApiQuotesaverdepositRequest) Execute() (*QuoteSaverDepositResponse, *http.Response, error) {
	return r.ApiService.QuotesaverdepositExecute(r)
}

/*
Quotesaverdeposit Method for Quotesaverdeposit

Provide a quote estimate for the provided saver deposit.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiQuotesaverdepositRequest
*/
func (a *QuoteApiService) Quotesaverdeposit(ctx context.Context) ApiQuotesaverdepositRequest {
	return ApiQuotesaverdepositRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return QuoteSaverDepositResponse
func (a *QuoteApiService) QuotesaverdepositExecute(r ApiQuotesaverdepositRequest) (*QuoteSaverDepositResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *QuoteSaverDepositResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QuoteApiService.Quotesaverdeposit")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/quote/saver/deposit"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
	}
	if r.asset != nil {
		localVarQueryParams.Add("asset", parameterToString(*r.asset, ""))
	}
	if r.amount != nil {
		localVarQueryParams.Add("amount", parameterToString(*r.amount, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiQuotesaverwithdrawRequest struct {
	ctx context.Context
	ApiService *QuoteApiService
	height *int64
	asset *string
	address *string
	withdrawBps *int64
}

// optional block height, defaults to current tip
func (r ApiQuotesaverwithdrawRequest) Height(height int64) ApiQuotesaverwithdrawRequest {
	r.height = &height
	return r
}

// the asset to withdraw
func (r ApiQuotesaverwithdrawRequest) Asset(asset string) ApiQuotesaverwithdrawRequest {
	r.asset = &asset
	return r
}

// the address for the position
func (r ApiQuotesaverwithdrawRequest) Address(address string) ApiQuotesaverwithdrawRequest {
	r.address = &address
	return r
}

// the basis points of the existing position to withdraw
func (r ApiQuotesaverwithdrawRequest) WithdrawBps(withdrawBps int64) ApiQuotesaverwithdrawRequest {
	r.withdrawBps = &withdrawBps
	return r
}

func (r ApiQuotesaverwithdrawRequest) Execute() (*QuoteSaverWithdrawResponse, *http.Response, error) {
	return r.ApiService.QuotesaverwithdrawExecute(r)
}

/*
Quotesaverwithdraw Method for Quotesaverwithdraw

Provide a quote estimate for the provided saver withdraw.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiQuotesaverwithdrawRequest
*/
func (a *QuoteApiService) Quotesaverwithdraw(ctx context.Context) ApiQuotesaverwithdrawRequest {
	return ApiQuotesaverwithdrawRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return QuoteSaverWithdrawResponse
func (a *QuoteApiService) QuotesaverwithdrawExecute(r ApiQuotesaverwithdrawRequest) (*QuoteSaverWithdrawResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *QuoteSaverWithdrawResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QuoteApiService.Quotesaverwithdraw")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/quote/saver/withdraw"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
	}
	if r.asset != nil {
		localVarQueryParams.Add("asset", parameterToString(*r.asset, ""))
	}
	if r.address != nil {
		localVarQueryParams.Add("address", parameterToString(*r.address, ""))
	}
	if r.withdrawBps != nil {
		localVarQueryParams.Add("withdraw_bps", parameterToString(*r.withdrawBps, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiQuoteswapRequest struct {
	ctx context.Context
	ApiService *QuoteApiService
	height *int64
	fromAsset *string
	toAsset *string
	amount *int64
	destination *string
	refundAddress *string
	streamingInterval *int64
	streamingQuantity *int64
	toleranceBps *int64
	affiliateBps *int64
	affiliate *string
}

// optional block height, defaults to current tip
func (r ApiQuoteswapRequest) Height(height int64) ApiQuoteswapRequest {
	r.height = &height
	return r
}

// the source asset
func (r ApiQuoteswapRequest) FromAsset(fromAsset string) ApiQuoteswapRequest {
	r.fromAsset = &fromAsset
	return r
}

// the target asset
func (r ApiQuoteswapRequest) ToAsset(toAsset string) ApiQuoteswapRequest {
	r.toAsset = &toAsset
	return r
}

// the source asset amount in 1e8 decimals
func (r ApiQuoteswapRequest) Amount(amount int64) ApiQuoteswapRequest {
	r.amount = &amount
	return r
}

// the destination address, required to generate memo
func (r ApiQuoteswapRequest) Destination(destination string) ApiQuoteswapRequest {
	r.destination = &destination
	return r
}

// the refund address, refunds will be sent here if the swap fails
func (r ApiQuoteswapRequest) RefundAddress(refundAddress string) ApiQuoteswapRequest {
	r.refundAddress = &refundAddress
	return r
}

// the interval in which streaming swaps are swapped
func (r ApiQuoteswapRequest) StreamingInterval(streamingInterval int64) ApiQuoteswapRequest {
	r.streamingInterval = &streamingInterval
	return r
}

// the quantity of swaps within a streaming swap
func (r ApiQuoteswapRequest) StreamingQuantity(streamingQuantity int64) ApiQuoteswapRequest {
	r.streamingQuantity = &streamingQuantity
	return r
}

// the maximum basis points from the current feeless swap price to set the limit in the generated memo
func (r ApiQuoteswapRequest) ToleranceBps(toleranceBps int64) ApiQuoteswapRequest {
	r.toleranceBps = &toleranceBps
	return r
}

// the affiliate fee in basis points
func (r ApiQuoteswapRequest) AffiliateBps(affiliateBps int64) ApiQuoteswapRequest {
	r.affiliateBps = &affiliateBps
	return r
}

// the affiliate (address or thorname)
func (r ApiQuoteswapRequest) Affiliate(affiliate string) ApiQuoteswapRequest {
	r.affiliate = &affiliate
	return r
}

func (r ApiQuoteswapRequest) Execute() (*QuoteSwapResponse, *http.Response, error) {
	return r.ApiService.QuoteswapExecute(r)
}

/*
Quoteswap Method for Quoteswap

Provide a quote estimate for the provided swap.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiQuoteswapRequest
*/
func (a *QuoteApiService) Quoteswap(ctx context.Context) ApiQuoteswapRequest {
	return ApiQuoteswapRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return QuoteSwapResponse
func (a *QuoteApiService) QuoteswapExecute(r ApiQuoteswapRequest) (*QuoteSwapResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *QuoteSwapResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "QuoteApiService.Quoteswap")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/quote/swap"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
	}
	if r.fromAsset != nil {
		localVarQueryParams.Add("from_asset", parameterToString(*r.fromAsset, ""))
	}
	if r.toAsset != nil {
		localVarQueryParams.Add("to_asset", parameterToString(*r.toAsset, ""))
	}
	if r.amount != nil {
		localVarQueryParams.Add("amount", parameterToString(*r.amount, ""))
	}
	if r.destination != nil {
		localVarQueryParams.Add("destination", parameterToString(*r.destination, ""))
	}
	if r.refundAddress != nil {
		localVarQueryParams.Add("refund_address", parameterToString(*r.refundAddress, ""))
	}
	if r.streamingInterval != nil {
		localVarQueryParams.Add("streaming_interval", parameterToString(*r.streamingInterval, ""))
	}
	if r.streamingQuantity != nil {
		localVarQueryParams.Add("streaming_quantity", parameterToString(*r.streamingQuantity, ""))
	}
	if r.toleranceBps != nil {
		localVarQueryParams.Add("tolerance_bps", parameterToString(*r.toleranceBps, ""))
	}
	if r.affiliateBps != nil {
		localVarQueryParams.Add("affiliate_bps", parameterToString(*r.affiliateBps, ""))
	}
	if r.affiliate != nil {
		localVarQueryParams.Add("affiliate", parameterToString(*r.affiliate, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
