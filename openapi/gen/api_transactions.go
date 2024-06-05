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
	"strings"
)


// TransactionsApiService TransactionsApi service
type TransactionsApiService service

type ApiTxRequest struct {
	ctx context.Context
	ApiService *TransactionsApiService
	hash string
	height *int64
}

// optional block height, defaults to current tip
func (r ApiTxRequest) Height(height int64) ApiTxRequest {
	r.height = &height
	return r
}

func (r ApiTxRequest) Execute() (*TxResponse, *http.Response, error) {
	return r.ApiService.TxExecute(r)
}

/*
Tx Method for Tx

Returns the observed transaction for a provided inbound or outbound hash.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param hash
 @return ApiTxRequest
*/
func (a *TransactionsApiService) Tx(ctx context.Context, hash string) ApiTxRequest {
	return ApiTxRequest{
		ApiService: a,
		ctx: ctx,
		hash: hash,
	}
}

// Execute executes the request
//  @return TxResponse
func (a *TransactionsApiService) TxExecute(r ApiTxRequest) (*TxResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *TxResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.Tx")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/tx/{hash}"
	localVarPath = strings.Replace(localVarPath, "{"+"hash"+"}", url.PathEscape(parameterToString(r.hash, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
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

type ApiTxSignersRequest struct {
	ctx context.Context
	ApiService *TransactionsApiService
	hash string
	height *int64
}

// optional block height, defaults to current tip
func (r ApiTxSignersRequest) Height(height int64) ApiTxSignersRequest {
	r.height = &height
	return r
}

func (r ApiTxSignersRequest) Execute() (*TxDetailsResponse, *http.Response, error) {
	return r.ApiService.TxSignersExecute(r)
}

/*
TxSigners Method for TxSigners

Returns the signers for a provided inbound or outbound hash.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param hash
 @return ApiTxSignersRequest
*/
func (a *TransactionsApiService) TxSigners(ctx context.Context, hash string) ApiTxSignersRequest {
	return ApiTxSignersRequest{
		ApiService: a,
		ctx: ctx,
		hash: hash,
	}
}

// Execute executes the request
//  @return TxDetailsResponse
func (a *TransactionsApiService) TxSignersExecute(r ApiTxSignersRequest) (*TxDetailsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *TxDetailsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.TxSigners")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/tx/details/{hash}"
	localVarPath = strings.Replace(localVarPath, "{"+"hash"+"}", url.PathEscape(parameterToString(r.hash, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
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

type ApiTxSignersOldRequest struct {
	ctx context.Context
	ApiService *TransactionsApiService
	hash string
	height *int64
}

// optional block height, defaults to current tip
func (r ApiTxSignersOldRequest) Height(height int64) ApiTxSignersOldRequest {
	r.height = &height
	return r
}

func (r ApiTxSignersOldRequest) Execute() (*TxSignersResponse, *http.Response, error) {
	return r.ApiService.TxSignersOldExecute(r)
}

/*
TxSignersOld Method for TxSignersOld

Deprecated - migrate to /thorchain/tx/details.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param hash
 @return ApiTxSignersOldRequest
*/
func (a *TransactionsApiService) TxSignersOld(ctx context.Context, hash string) ApiTxSignersOldRequest {
	return ApiTxSignersOldRequest{
		ApiService: a,
		ctx: ctx,
		hash: hash,
	}
}

// Execute executes the request
//  @return TxSignersResponse
func (a *TransactionsApiService) TxSignersOldExecute(r ApiTxSignersOldRequest) (*TxSignersResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *TxSignersResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.TxSignersOld")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/tx/{hash}/signers"
	localVarPath = strings.Replace(localVarPath, "{"+"hash"+"}", url.PathEscape(parameterToString(r.hash, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
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

type ApiTxStagesRequest struct {
	ctx context.Context
	ApiService *TransactionsApiService
	hash string
	height *int64
}

// optional block height, defaults to current tip
func (r ApiTxStagesRequest) Height(height int64) ApiTxStagesRequest {
	r.height = &height
	return r
}

func (r ApiTxStagesRequest) Execute() (*TxStagesResponse, *http.Response, error) {
	return r.ApiService.TxStagesExecute(r)
}

/*
TxStages Method for TxStages

Returns the processing stages of a provided inbound hash.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param hash
 @return ApiTxStagesRequest
*/
func (a *TransactionsApiService) TxStages(ctx context.Context, hash string) ApiTxStagesRequest {
	return ApiTxStagesRequest{
		ApiService: a,
		ctx: ctx,
		hash: hash,
	}
}

// Execute executes the request
//  @return TxStagesResponse
func (a *TransactionsApiService) TxStagesExecute(r ApiTxStagesRequest) (*TxStagesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *TxStagesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.TxStages")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/tx/stages/{hash}"
	localVarPath = strings.Replace(localVarPath, "{"+"hash"+"}", url.PathEscape(parameterToString(r.hash, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
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

type ApiTxStatusRequest struct {
	ctx context.Context
	ApiService *TransactionsApiService
	hash string
	height *int64
}

// optional block height, defaults to current tip
func (r ApiTxStatusRequest) Height(height int64) ApiTxStatusRequest {
	r.height = &height
	return r
}

func (r ApiTxStatusRequest) Execute() (*TxStatusResponse, *http.Response, error) {
	return r.ApiService.TxStatusExecute(r)
}

/*
TxStatus Method for TxStatus

Returns the status of a provided inbound hash.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param hash
 @return ApiTxStatusRequest
*/
func (a *TransactionsApiService) TxStatus(ctx context.Context, hash string) ApiTxStatusRequest {
	return ApiTxStatusRequest{
		ApiService: a,
		ctx: ctx,
		hash: hash,
	}
}

// Execute executes the request
//  @return TxStatusResponse
func (a *TransactionsApiService) TxStatusExecute(r ApiTxStatusRequest) (*TxStatusResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *TxStatusResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TransactionsApiService.TxStatus")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/tx/status/{hash}"
	localVarPath = strings.Replace(localVarPath, "{"+"hash"+"}", url.PathEscape(parameterToString(r.hash, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.height != nil {
		localVarQueryParams.Add("height", parameterToString(*r.height, ""))
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
