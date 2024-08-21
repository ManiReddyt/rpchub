package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"rpchub/rest"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/robfig/cron/v3"
)

type Config struct {
	Chains map[string][]string `json:"chains"`
}

var (
	config      Config
	HealthyRPCs map[string][]string
	HealthLock  sync.Mutex
)

func ReadConfigFromJson() Config {
	jsonData, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func main() {

	config = ReadConfigFromJson()

	c := cron.New(cron.WithLocation(time.UTC))
	c.AddFunc("@every 10s", func() {
		//health check should happen here and update the healthyRPC's
		HealthLock.Lock()
		defer HealthLock.Unlock()
		HealthyRPCs = make(map[string][]string)

		for chainId, rpcs := range config.Chains {
			HealthyRPCs[chainId] = []string{}
			for _, rpc := range rpcs {
				//check health of rpc
				ethClient, err := ethclient.Dial(rpc)
				if err != nil {
					fmt.Printf("rpc %s is unhealthy\n", rpc)
					continue
				}

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
				defer cancel()

				cId, err := ethClient.ChainID(ctx)
				if err != nil {
					fmt.Printf("rpc %s is unhealthy\n", rpc)
					continue
				}
				if chainId != cId.String() {
					fmt.Printf("ChainId %s is not matching\n", rpc)
					continue
				}

				ethClient.Close()
				HealthyRPCs[chainId] = append(HealthyRPCs[chainId], rpc)
			}
		}

		fmt.Printf("healthyRPCs: %#v\n", HealthyRPCs)
	})

	//for every 10minutes read the config file because we don't want to restart the server if there is a change in config
	c.AddFunc("@every 10m", func() {
		config = ReadConfigFromJson()
	})

	go c.Start()

	server := rest.NewServer(&HealthyRPCs, &HealthLock)
	server.Start()

}
