package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	stats "github.com/semihalev/gin-stats"
)

func startGinService() {
	log.Println("initiating api-service...")

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(stats.RequestStats())

	r.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, stats.Report())
	})
	r.GET("/health", API_HealthCheck)
	r.GET("/getbitcoinfee", API_GetBitcoinFee)

	err := r.Run("0.0.0.0:" + strconv.Itoa(serviceCfg.APIPort))
	if err != nil {
		panic(err)
	}
}

func API_GetBitcoinFee(c *gin.Context) {
	feeRWLock.RLock()
	defer feeRWLock.RUnlock()

	c.JSON(http.StatusOK, API_respond{
		Result: bitcoinFee,
		Error:  nil,
	})
}

func API_HealthCheck(c *gin.Context) {
	status := "healthy"
	bitcoinFeeStatus := "updated"

	if lastBitcoinFee < 0 {
		status = "unhealthy"
		bitcoinFeeStatus = "outdated"
	}
	c.JSON(http.StatusOK, gin.H{
		"status":           status,
		"bitcoinFeeStatus": bitcoinFeeStatus,
	})
}

func buildGinErrorRespond(err error) *API_respond {
	errStr := err.Error()
	respond := API_respond{
		Result: nil,
		Error:  &errStr,
	}
	return &respond
}
