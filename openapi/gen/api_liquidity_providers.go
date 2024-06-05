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


// LiquidityProvidersApiService LiquidityProvidersApi service
type LiquidityProvidersApiService service

type ApiLiquidityProviderRequest struct {
	ctx context.Context
	ApiService *LiquidityProvidersApiService
	asset string
	address string
	height *int64
}

// optional block height, defaults to current tip
func (r ApiLiquidityProviderRequest) Height(height int64) ApiLiquidityProviderRequest {
	r.height = &height
	return r
}

func (r ApiLiquidityProviderRequest) Execute() (*LiquidityProvider, *http.Response, error) {
	return r.ApiService.LiquidityProviderExecute(r)
}

/*
LiquidityProvider Method for LiquidityProvider

Returns the liquidity provider information for an address and asset.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param asset
 @param address
 @return ApiLiquidityProviderRequest
*/
func (a *LiquidityProvidersApiService) LiquidityProvider(ctx context.Context, asset string, address string) ApiLiquidityProviderRequest {
	return ApiLiquidityProviderRequest{
		ApiService: a,
		ctx: ctx,
		asset: asset,
		address: address,
	}
}

// Execute executes the request
//  @return LiquidityProvider
func (a *LiquidityProvidersApiService) LiquidityProviderExecute(r ApiLiquidityProviderRequest) (*LiquidityProvider, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *LiquidityProvider
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LiquidityProvidersApiService.LiquidityProvider")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/pool/{asset}/liquidity_provider/{address}"
	localVarPath = strings.Replace(localVarPath, "{"+"asset"+"}", url.PathEscape(parameterToString(r.asset, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"address"+"}", url.PathEscape(parameterToString(r.address, "")), -1)

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

type ApiLiquidityProvidersRequest struct {
	ctx context.Context
	ApiService *LiquidityProvidersApiService
	asset string
	height *int64
}

// optional block height, defaults to current tip
func (r ApiLiquidityProvidersRequest) Height(height int64) ApiLiquidityProvidersRequest {
	r.height = &height
	return r
}

func (r ApiLiquidityProvidersRequest) Execute() ([]LiquidityProviderSummary, *http.Response, error) {
	return r.ApiService.LiquidityProvidersExecute(r)
}

/*
LiquidityProviders Method for LiquidityProviders

Returns all liquidity provider information for an asset.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param asset
 @return ApiLiquidityProvidersRequest
*/
func (a *LiquidityProvidersApiService) LiquidityProviders(ctx context.Context, asset string) ApiLiquidityProvidersRequest {
	return ApiLiquidityProvidersRequest{
		ApiService: a,
		ctx: ctx,
		asset: asset,
	}
}

// Execute executes the request
//  @return []LiquidityProviderSummary
func (a *LiquidityProvidersApiService) LiquidityProvidersExecute(r ApiLiquidityProvidersRequest) ([]LiquidityProviderSummary, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  []LiquidityProviderSummary
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LiquidityProvidersApiService.LiquidityProviders")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/thorchain/pool/{asset}/liquidity_providers"
	localVarPath = strings.Replace(localVarPath, "{"+"asset"+"}", url.PathEscape(parameterToString(r.asset, "")), -1)

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