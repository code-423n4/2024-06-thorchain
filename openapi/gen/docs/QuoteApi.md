# \QuoteApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Quoteloanclose**](QuoteApi.md#Quoteloanclose) | **Get** /thorchain/quote/loan/close | 
[**Quoteloanopen**](QuoteApi.md#Quoteloanopen) | **Get** /thorchain/quote/loan/open | 
[**Quotesaverdeposit**](QuoteApi.md#Quotesaverdeposit) | **Get** /thorchain/quote/saver/deposit | 
[**Quotesaverwithdraw**](QuoteApi.md#Quotesaverwithdraw) | **Get** /thorchain/quote/saver/withdraw | 
[**Quoteswap**](QuoteApi.md#Quoteswap) | **Get** /thorchain/quote/swap | 



## Quoteloanclose

> QuoteLoanCloseResponse Quoteloanclose(ctx).Height(height).FromAsset(fromAsset).RepayBps(repayBps).ToAsset(toAsset).LoanOwner(loanOwner).MinOut(minOut).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    height := int64(789) // int64 | optional block height, defaults to current tip (optional)
    fromAsset := "ETH.ETH" // string | the asset used to repay the loan (optional)
    repayBps := int64(100) // int64 | the basis points of the existing position to repay (optional)
    toAsset := "BTC.BTC" // string | the collateral asset of the loan (optional)
    loanOwner := "BTC.BTC" // string | the owner of the loan collateral (optional)
    minOut := "1234" // string | the minimum amount of the target asset to accept (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.QuoteApi.Quoteloanclose(context.Background()).Height(height).FromAsset(fromAsset).RepayBps(repayBps).ToAsset(toAsset).LoanOwner(loanOwner).MinOut(minOut).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `QuoteApi.Quoteloanclose``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Quoteloanclose`: QuoteLoanCloseResponse
    fmt.Fprintf(os.Stdout, "Response from `QuoteApi.Quoteloanclose`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiQuoteloancloseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **height** | **int64** | optional block height, defaults to current tip | 
 **fromAsset** | **string** | the asset used to repay the loan | 
 **repayBps** | **int64** | the basis points of the existing position to repay | 
 **toAsset** | **string** | the collateral asset of the loan | 
 **loanOwner** | **string** | the owner of the loan collateral | 
 **minOut** | **string** | the minimum amount of the target asset to accept | 

### Return type

[**QuoteLoanCloseResponse**](QuoteLoanCloseResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Quoteloanopen

> QuoteLoanOpenResponse Quoteloanopen(ctx).Height(height).FromAsset(fromAsset).Amount(amount).ToAsset(toAsset).Destination(destination).MinOut(minOut).AffiliateBps(affiliateBps).Affiliate(affiliate).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    height := int64(789) // int64 | optional block height, defaults to current tip (optional)
    fromAsset := "BTC.BTC" // string | the collateral asset (optional)
    amount := int64(1000000) // int64 | the collateral asset amount in 1e8 decimals (optional)
    toAsset := "ETH.ETH" // string | the target asset to receive (loan denominated in TOR regardless) (optional)
    destination := "0x1c7b17362c84287bd1184447e6dfeaf920c31bbe" // string | the destination address, required to generate memo (optional)
    minOut := "1234" // string | the minimum amount of the target asset to accept (optional)
    affiliateBps := int64(100) // int64 | the affiliate fee in basis points (optional)
    affiliate := "t" // string | the affiliate (address or thorname) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.QuoteApi.Quoteloanopen(context.Background()).Height(height).FromAsset(fromAsset).Amount(amount).ToAsset(toAsset).Destination(destination).MinOut(minOut).AffiliateBps(affiliateBps).Affiliate(affiliate).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `QuoteApi.Quoteloanopen``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Quoteloanopen`: QuoteLoanOpenResponse
    fmt.Fprintf(os.Stdout, "Response from `QuoteApi.Quoteloanopen`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiQuoteloanopenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **height** | **int64** | optional block height, defaults to current tip | 
 **fromAsset** | **string** | the collateral asset | 
 **amount** | **int64** | the collateral asset amount in 1e8 decimals | 
 **toAsset** | **string** | the target asset to receive (loan denominated in TOR regardless) | 
 **destination** | **string** | the destination address, required to generate memo | 
 **minOut** | **string** | the minimum amount of the target asset to accept | 
 **affiliateBps** | **int64** | the affiliate fee in basis points | 
 **affiliate** | **string** | the affiliate (address or thorname) | 

### Return type

[**QuoteLoanOpenResponse**](QuoteLoanOpenResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Quotesaverdeposit

> QuoteSaverDepositResponse Quotesaverdeposit(ctx).Height(height).Asset(asset).Amount(amount).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    height := int64(789) // int64 | optional block height, defaults to current tip (optional)
    asset := "BTC.BTC" // string | the asset to deposit (optional)
    amount := int64(1000000) // int64 | the source asset amount in 1e8 decimals (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.QuoteApi.Quotesaverdeposit(context.Background()).Height(height).Asset(asset).Amount(amount).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `QuoteApi.Quotesaverdeposit``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Quotesaverdeposit`: QuoteSaverDepositResponse
    fmt.Fprintf(os.Stdout, "Response from `QuoteApi.Quotesaverdeposit`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiQuotesaverdepositRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **height** | **int64** | optional block height, defaults to current tip | 
 **asset** | **string** | the asset to deposit | 
 **amount** | **int64** | the source asset amount in 1e8 decimals | 

### Return type

[**QuoteSaverDepositResponse**](QuoteSaverDepositResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Quotesaverwithdraw

> QuoteSaverWithdrawResponse Quotesaverwithdraw(ctx).Height(height).Asset(asset).Address(address).WithdrawBps(withdrawBps).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    height := int64(789) // int64 | optional block height, defaults to current tip (optional)
    asset := "BTC.BTC" // string | the asset to withdraw (optional)
    address := "bc1qd45uzetakjvdy5ynjjyp4nlnj89am88e4e5jeq" // string | the address for the position (optional)
    withdrawBps := int64(100) // int64 | the basis points of the existing position to withdraw (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.QuoteApi.Quotesaverwithdraw(context.Background()).Height(height).Asset(asset).Address(address).WithdrawBps(withdrawBps).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `QuoteApi.Quotesaverwithdraw``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Quotesaverwithdraw`: QuoteSaverWithdrawResponse
    fmt.Fprintf(os.Stdout, "Response from `QuoteApi.Quotesaverwithdraw`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiQuotesaverwithdrawRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **height** | **int64** | optional block height, defaults to current tip | 
 **asset** | **string** | the asset to withdraw | 
 **address** | **string** | the address for the position | 
 **withdrawBps** | **int64** | the basis points of the existing position to withdraw | 

### Return type

[**QuoteSaverWithdrawResponse**](QuoteSaverWithdrawResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Quoteswap

> QuoteSwapResponse Quoteswap(ctx).Height(height).FromAsset(fromAsset).ToAsset(toAsset).Amount(amount).Destination(destination).RefundAddress(refundAddress).StreamingInterval(streamingInterval).StreamingQuantity(streamingQuantity).ToleranceBps(toleranceBps).AffiliateBps(affiliateBps).Affiliate(affiliate).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    height := int64(789) // int64 | optional block height, defaults to current tip (optional)
    fromAsset := "BTC.BTC" // string | the source asset (optional)
    toAsset := "ETH.ETH" // string | the target asset (optional)
    amount := int64(1000000) // int64 | the source asset amount in 1e8 decimals (optional)
    destination := "0x1c7b17362c84287bd1184447e6dfeaf920c31bbe" // string | the destination address, required to generate memo (optional)
    refundAddress := "0x1c7b17362c84287bd1184447e6dfeaf920c31bbe" // string | the refund address, refunds will be sent here if the swap fails (optional)
    streamingInterval := int64(10) // int64 | the interval in which streaming swaps are swapped (optional)
    streamingQuantity := int64(10) // int64 | the quantity of swaps within a streaming swap (optional)
    toleranceBps := int64(100) // int64 | the maximum basis points from the current feeless swap price to set the limit in the generated memo (optional)
    affiliateBps := int64(100) // int64 | the affiliate fee in basis points (optional)
    affiliate := "t" // string | the affiliate (address or thorname) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.QuoteApi.Quoteswap(context.Background()).Height(height).FromAsset(fromAsset).ToAsset(toAsset).Amount(amount).Destination(destination).RefundAddress(refundAddress).StreamingInterval(streamingInterval).StreamingQuantity(streamingQuantity).ToleranceBps(toleranceBps).AffiliateBps(affiliateBps).Affiliate(affiliate).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `QuoteApi.Quoteswap``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Quoteswap`: QuoteSwapResponse
    fmt.Fprintf(os.Stdout, "Response from `QuoteApi.Quoteswap`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiQuoteswapRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **height** | **int64** | optional block height, defaults to current tip | 
 **fromAsset** | **string** | the source asset | 
 **toAsset** | **string** | the target asset | 
 **amount** | **int64** | the source asset amount in 1e8 decimals | 
 **destination** | **string** | the destination address, required to generate memo | 
 **refundAddress** | **string** | the refund address, refunds will be sent here if the swap fails | 
 **streamingInterval** | **int64** | the interval in which streaming swaps are swapped | 
 **streamingQuantity** | **int64** | the quantity of swaps within a streaming swap | 
 **toleranceBps** | **int64** | the maximum basis points from the current feeless swap price to set the limit in the generated memo | 
 **affiliateBps** | **int64** | the affiliate fee in basis points | 
 **affiliate** | **string** | the affiliate (address or thorname) | 

### Return type

[**QuoteSwapResponse**](QuoteSwapResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

