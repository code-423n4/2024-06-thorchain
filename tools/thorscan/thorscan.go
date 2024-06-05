package thorscan

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/thorchain/thornode/constants"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// -------------------------------------------------------------------------------------
// Config
// -------------------------------------------------------------------------------------

const (
	// ---------- environment keys ----------

	EnvRPCEndpoint = "RPC_ENDPOINT"
	EnvAPIEndpoint = "API_ENDPOINT"
	EnvParallelism = "PARALLELISM"
)

// -------------------------------------------------------------------------------------
// HTTP
// -------------------------------------------------------------------------------------

// Transport sets the X-Client-ID header on all requests.
type Transport struct {
	Transport http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Client-ID", "thorscan-go")
	return t.Transport.RoundTrip(req)
}

var httpClient *http.Client

// -------------------------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------------------------

func height() (int64, error) {
	res, err := httpClient.Get(RPCEndpoint + "/status")
	if err != nil {
		return 0, err
	}

	// decode response
	var statusResp struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}
	err = json.NewDecoder(res.Body).Decode(&statusResp)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to decode status response")
	}

	res.Body.Close()

	return strconv.ParseInt(statusResp.Result.SyncInfo.LatestBlockHeight, 10, 64)
}

func getBlock(height int64) (*openapi.BlockResponse, error) {
	url := APIEndpoint + "/thorchain/block?height=" + strconv.FormatInt(height, 10)

	// build request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// accept gzip
	req.Header.Set("Accept-Encoding", "gzip")

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// close body
	defer res.Body.Close()

	// wrap response body in a gzip reader
	if strings.Contains(res.Header.Get("Content-Encoding"), "gzip") {
		res.Body, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
	}

	// check status code
	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("block not found")
	default:
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	// decode response
	var blockResp openapi.BlockResponse
	err = json.NewDecoder(res.Body).Decode(&blockResp)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode block response")
	}

	return &blockResp, nil
}

// -------------------------------------------------------------------------------------
// Init
// -------------------------------------------------------------------------------------

var (
	Parallelism = 4
	RPCEndpoint = "https://rpc-v1.ninerealms.com"
	APIEndpoint = "https://thornode-v1.ninerealms.com"
)

func init() {
	var err error

	// set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// set config from env
	if e := os.Getenv(EnvRPCEndpoint); e != "" {
		log.Info().Str("endpoint", e).Msg("setting rpc endpoint")
		RPCEndpoint = e
	}
	if e := os.Getenv(EnvAPIEndpoint); e != "" {
		log.Info().Str("endpoint", e).Msg("setting api endpoint")
		APIEndpoint = e
	}
	if e := os.Getenv(EnvParallelism); e != "" {
		log.Info().Str("prefetch", e).Msg("setting prefetch blocks")
		Parallelism, err = strconv.Atoi(e)
		if err != nil {
			log.Fatal().Err(err).Msg("bad prefetch value")
		}
	}

	// use our own transport to set the client id
	transport := &Transport{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     false,
			MaxIdleConns:          Parallelism * 2,
			MaxIdleConnsPerHost:   Parallelism * 2,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// create new client with better connection reuse
	httpClient = &http.Client{Transport: transport}
}

// -------------------------------------------------------------------------------------
// Exported
// -------------------------------------------------------------------------------------

func Scan(startHeight, stopHeight int) <-chan *openapi.BlockResponse {
	// get current height if start was not provided
	if startHeight <= 0 || stopHeight < 0 {
		height, err := height()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get current height")
		}

		// set start height
		if startHeight <= 0 {
			startHeight = int(height) + startHeight
		}
		if stopHeight < 0 { // zero height means tail indefinitely
			stopHeight = int(height) + stopHeight
		}
	}

	// create queue for block heights to fetch
	queue := make(chan int64)
	go func() {
		for height := int64(startHeight); stopHeight == 0 || int(height) <= stopHeight; height++ {
			queue <- height
		}
	}()

	// setup ring buffer for block prefetching with routine per slot
	ring := make([]chan *openapi.BlockResponse, Parallelism)
	shutdown := make(chan struct{}, Parallelism-1)
	for i := 0; i < Parallelism; i++ {
		ring[i] = make(chan *openapi.BlockResponse)
		go func(i int) {
			for height := range queue {
				for {
					b, err := getBlock(height)
					if err != nil {
						if !strings.Contains(err.Error(), "block not found") {
							log.Error().Err(err).Int64("height", height).Msg("failed to fetch block")
						}
						time.Sleep(constants.ThorchainBlockTime)
						continue
					}
					ring[int(height)%Parallelism] <- b

					// allow all but one routine to exit once we near tip
					blockTime, err := time.Parse(time.RFC3339, b.Header.Time)
					if err != nil {
						log.Fatal().Err(err).Msg("failed to parse block time")
					}
					near := time.Now().Add(-constants.ThorchainBlockTime * time.Duration(Parallelism))
					if err == nil && blockTime.After(near) {
						select {
						case shutdown <- struct{}{}:
							log.Debug().Int64("height", height).Msg("shutting down extra worker")
							return
						default:
						}
					}

					break
				}
			}
		}(i)
	}

	// start sequential reader to send to blocks channel
	out := make(chan *openapi.BlockResponse)
	go func() {
		for height := int64(startHeight); stopHeight == 0 || int(height) <= stopHeight; height++ {
			out <- <-ring[int(height)%Parallelism]
		}
		close(out)
	}()

	return out
}
