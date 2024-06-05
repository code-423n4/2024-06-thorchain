package utxo

import (
	"net/http"
	"os"
	"testing"
	"time"

	"gitlab.com/thorchain/thornode/bifrost/metrics"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/config"
	. "gopkg.in/check.v1"
)

const (
	bob      = "bob"
	password = "password"
)

var m *metrics.Metrics

func TestPackage(t *testing.T) { TestingT(t) }

func GetMetricForTest(c *C, chain common.Chain) *metrics.Metrics {
	if m != nil {
		return m
	}
	var err error
	m, err = metrics.NewMetrics(config.BifrostMetricsConfiguration{
		Enabled:      false,
		ListenPort:   9000,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Chains:       common.Chains{common.DOGEChain, common.BCHChain, common.LTCChain, common.BTCChain},
	})
	c.Assert(m, NotNil)
	c.Assert(err, IsNil)
	return m
}

func httpTestHandler(c *C, rw http.ResponseWriter, fixture string) {
	content, err := os.ReadFile(fixture)
	if err != nil {
		c.Fatal(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	if _, err = rw.Write(content); err != nil {
		c.Fatal(err)
	}
}
