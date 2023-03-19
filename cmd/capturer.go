package main

import (
	"coin-capturer/internal/erc20"
	"coin-capturer/internal/listener"
	"coin-capturer/internal/util"
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
			event, err := erc20.ParseTransfer(vLog)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("[%s] Received transfer event - Tx Hash: %s, From: %s, To: %s, Value: %s USDT\n",
				time.Now().Format("2006-01-02 15:04:05"), vLog.TxHash.Hex(), event.From.Hex(), event.To.Hex(), util.ToDecimal(event.Tokens, 18))
		}
	}
}
