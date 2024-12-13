package main

import (
	"blockchain/blockchain"
	"blockchain/node"
)

func main() {
	// Create the blockchain and start mining nodes
	blockchainInstance := blockchain.NewBlockchain()

	// Start user nodes
	for i := 1; i <= 4; i++ {
		go node.StartUserNode(i, blockchainInstance)
	}

	// Start miner nodes
	for i := 1; i <= 3; i++ {
		go node.StartMinerNode(i, blockchainInstance)
	}

	// Keep the main thread alive
	select {}
}
