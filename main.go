package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	ethereum "github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethClient "github.com/ethereum/go-ethereum/ethclient"
)

const (
	usdtAddress         = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	erc20TransferTopic0 = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

func main() {
	websocketAddress := os.Getenv("ERIGON_WEBSOCKET_URL")
	if websocketAddress == "" {
		fmt.Println("Please set ERIGON_WEBSOCKET_URL environment variable")
		return
	}
	client, err := ethClient.Dial(websocketAddress)
	if err != nil {
		fmt.Printf("Failed to connect to websocket: %v\n", err)
		return
	}

	erc20, uniswap := parseCliArgs()

	var filter ethereum.FilterQuery
	if erc20 {
		filter.Topics = [][]ethCommon.Hash{{ethCommon.HexToHash(erc20TransferTopic0)}}
	}
	if uniswap {
		filter.Addresses = []ethCommon.Address{ethCommon.HexToAddress(usdtAddress)}
	}

	receivedLogs := make(chan types.Log)
	subscription, err := client.SubscribeFilterLogs(context.Background(), filter, receivedLogs)
	if err != nil {
		fmt.Printf("Unable to create subscription: %v\n", err)
		return
	}
	defer subscription.Unsubscribe()
	fmt.Println("Subscription successful")
	errChan := subscription.Err()

	for {
		select {
		case log := <-receivedLogs:
			fmt.Printf("Received log %d of block %d from address %s and topic0 %s\n",
				log.Index,
				log.BlockNumber,
				log.Address.String(),
				log.Topics[0].String(),
			)
		case err := <-errChan:
			fmt.Printf("Subscription error encountered: %v\n", err)
			return
		}
	}
}

// Parses command line flags, and returns whether or not to filter by the ERC20
// transfer topic 0 hash, and whether or not to filter by USDT address
func parseCliArgs() (bool, bool) {
	erc20TransferFlag := flag.Bool("erc20", false, "Whether to filter logs by topic 0 hash "+erc20TransferTopic0)
	usdtFlag := flag.Bool("usdt", false, "Whether to filter logs by address "+usdtAddress)
	flag.Parse()
	return *erc20TransferFlag, *usdtFlag
}
