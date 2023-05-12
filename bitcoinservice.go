package main

import (
	"fmt"
	"sync"
	"time"

	resty "github.com/go-resty/resty/v2"
)

var feeRWLock sync.RWMutex
var bitcoinFee float64 // satoshi / vbyte
var lastBitcoinFee float64

type BlockCypherFeeResponse struct {
	HighFee   uint `json:"high_fee_per_kb"`
	MediumFee uint `json:"medium_fee_per_kb"`
	LowFee    uint `json:"low_fee_per_kb"`
}

func initBitcoinService() {
	go func() {
		bitcoinFee = float64(serviceCfg.InitBitcoinFee)
		lastBitcoinFee = -1
		for {
			func() {
				feeRWLock.Lock()
				defer func() {
					feeRWLock.Unlock()
					time.Sleep(time.Duration(serviceCfg.BitcoinRefreshInterval) * time.Minute)
				}()

				client := resty.New()

				response, err := client.R().
					Get(BlockCypherBTCHost)

				if err != nil {
					lastBitcoinFee = -1
					msg := fmt.Sprintf("[BitcoinFee][WRN] Error get response: %v", err)
					fmt.Printf("%v\n", msg)
					sendSlackNotification(msg, serviceCfg.WebHookURL)
					return
				}
				if response.StatusCode() != 200 {
					lastBitcoinFee = -1
					msg := fmt.Sprintf("[BitcoinFee][WRN] Response status code: %v\n", response.StatusCode())
					fmt.Printf("%v\n", msg)
					sendSlackNotification(msg, serviceCfg.WebHookURL)
					return
				}
				var responseBody BlockCypherFeeResponse
				err = json.Unmarshal(response.Body(), &responseBody)
				if err != nil {
					lastBitcoinFee = -1
					msg := fmt.Sprintf("[BitcoinFee][WRN] Error parse body: %v\n", err)
					fmt.Printf("%v\n", msg)
					sendSlackNotification(msg, serviceCfg.WebHookURL)
					return
				}
				lastBitcoinFee = float64(responseBody.LowFee) / 1024
				bitcoinFee = lastBitcoinFee
			}()
		}
	}()
}
