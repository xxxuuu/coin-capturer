package main

import (
	"coin-capturer/internal/listener"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"time"
)

func main() {
	config, err := listener.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded config: %+v", config)

	l := listener.New(config)
	logs := make(chan types.Log)
	sub, err := l.Run(logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Printf("[%s] Received transfer event - Tx Hash:%s, From:%s, To:%s, Value:%s\n",
				time.Now().Format("2006-01-02 15:04:05"), vLog.TxHash.Hex(), vLog.Topics[1].Hex(), vLog.Topics[2].Hex(), vLog.Data)
		}
	}
}
