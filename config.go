package main

import (
	"flag"
	"io/ioutil"
	"log"
)

var ENABLE_PROFILER bool
var serviceCfg Config

type Config struct {
	APIPort                int    `json:"apiport"`
	WebHookURL             string `json:"webhook_url"`
	InitBitcoinFee         uint   `json:"init_bitcoin_fee"`
	BitcoinRefreshInterval uint   `json:"bitcoin_fee_refresh_interval"`
}

func readConfigAndArg() {
	data, err := ioutil.ReadFile("./cfg.json")
	if err != nil {
		log.Println(err)
		// return
	}
	var tempCfg Config
	if data != nil {
		err = json.Unmarshal(data, &tempCfg)
		if err != nil {
			panic(err)
		}
	}

	argProfiler := flag.Bool("profiler", false, "set profiler")
	flag.Parse()
	if tempCfg.APIPort == 0 {
		tempCfg.APIPort = DefaultAPIPort
	}
	if tempCfg.InitBitcoinFee == 0 {
		tempCfg.InitBitcoinFee = DefaultBitcoinFee
	}
	if tempCfg.BitcoinRefreshInterval == 0 {
		tempCfg.BitcoinRefreshInterval = DefaultBitcoinFeeRefreshInterval
	}

	ENABLE_PROFILER = *argProfiler
	serviceCfg = tempCfg
}
