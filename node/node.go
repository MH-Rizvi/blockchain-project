package node

import (
	"blockchain/block"
	"blockchain/blockchain"
	"blockchain/transaction"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var targetHashPrefix = "0" // Target prefix for Proof-of-Work (PoW)

// StartUserNode initializes a user node that creates and sends transactions
func StartUserNode(id int, blockchain *blockchain.Blockchain) {
	for {
		// Simulate fetching dataset and algorithm from IPFS
		datasetHash := fmt.Sprintf("dataset-hash-%d", rand.Int())
		algorithmHash := fmt.Sprintf("algorithm-hash-%d", rand.Int())
		outputHash := fmt.Sprintf("output-hash-%d", rand.Int())

		// Create a new transaction
		t := transaction.NewTransaction(datasetHash, algorithmHash, outputHash)

		// Send transaction to miners
		sendTransactionToMiners(t)

		// Simulate a delay before next transaction
		time.Sleep(time.Second * 5)
	}
}

func sendTransactionToMiners(t transaction.Transaction) {
	// Simulate sending a transaction to miner nodes over TCP
	conn, err := net.Dial("tcp", "miner-node-1:8081")
	if err != nil {
		fmt.Println("Error connecting to miner:", err)
		return
	}
	defer conn.Close()

	// Send the transaction
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(t)
	if err != nil {
		fmt.Println("Error sending transaction:", err)
	}
}

// StartMinerNode initializes a miner node that verifies transactions and mines blocks
func StartMinerNode(id int, blockchain *blockchain.Blockchain) {
	for {
		// Listen for incoming transactions
		listenForTransactions(id, blockchain)
	}
}

func listenForTransactions(id int, blockchain *blockchain.Blockchain) {
	// Start TCP server to receive transactions
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", 8081+id)) // Miner ports 8081, 8082, 8083
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Receive transaction
		var t transaction.Transaction
		decoder := json.NewDecoder(conn)
		err = decoder.Decode(&t)
		if err != nil {
			fmt.Println("Error decoding transaction:", err)
			continue
		}

		// Verify transaction
		if blockchain.VerifyTransaction(t) {
			blockchain.AddTransaction(t)
			if len(blockchain.CurrentTransactions) >= 4 {
				// Mine the block
				mineBlock(blockchain)
				// Reset transactions for the next block
				blockchain.CurrentTransactions = nil
			}
		}

		conn.Close()
	}
}

func mineBlock(blockchain *blockchain.Blockchain) {
	// Create a new block from the current transactions
	lastBlock := blockchain.GetLastBlock()
	newBlock := block.NewBlock(lastBlock.Hash, blockchain.CurrentTransactions)

	// Attempt to mine the block using PoW
	nonce := 0
	for {
		hash := newBlock.CalculateHash(nonce)
		if hash[:len(targetHashPrefix)] == targetHashPrefix {
			// Block mined
			newBlock.Nonce = nonce
			blockchain.AddBlock(newBlock)
			fmt.Printf("Miner mined a block: %v\n", newBlock)
			// Propagate the block
			propagateBlock(newBlock)
			break
		}
		nonce++
	}
}

func propagateBlock(b block.Block) {
	// Simulate block propagation via gossip or flooding
	fmt.Println("Propagating block:", b)
	// Send to other miners (simplified)
}
