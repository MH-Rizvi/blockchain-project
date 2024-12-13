package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

// Transaction structure
type Transaction struct {
	DatasetHash   string
	AlgorithmHash string
	OutputHash    string
}

// Block structure
type Block struct {
	Index        int
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
	Timestamp    time.Time
}

// Blockchain
var Blockchain []Block

// Utility to calculate hash
func calculateHash(prevBlockHash string, transactions []Transaction, nonce int) string {
	// Concatenate previous block hash, output hashes of transactions, and nonce
	record := prevBlockHash
	for _, tx := range transactions {
		record += tx.OutputHash
	}
	record += fmt.Sprintf("%d", nonce)

	// Compute SHA256 hash
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// Proof-of-Work function
func proofOfWork(prevBlockHash string, transactions []Transaction, difficulty int) (string, int, bool) {
	prefix := strings.Repeat("0", difficulty)
	var nonce int
	for {
		// Calculate hash
		hash := calculateHash(prevBlockHash, transactions, nonce)
		// Check if hash meets the target difficulty
		if strings.HasPrefix(hash, prefix) {
			return hash, nonce, true
		}
		nonce++
	}
}

// Verify a transaction
func verifyTransaction(tx Transaction) bool {
	// Simulate dataset and algorithm fetching from IPFS and compute the output hash
	calculatedOutputHash := tx.OutputHash // Replace with real computation if needed
	return calculatedOutputHash == tx.OutputHash
}

func minerNode(minerID string) {
	ln, err := net.Listen("tcp", ":8081") // Port for the miner node
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Println("Miner node listening on port 8081...")

	var transactions []Transaction
	difficulty := 4
	blockMined := false // Track if a block has been mined

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Read transaction from user
		var tx Transaction
		decoder := json.NewDecoder(conn)
		err = decoder.Decode(&tx)
		if err != nil {
			fmt.Println("Error decoding transaction:", err)
			conn.Close()
			continue
		}
		conn.Close()

		// Verify transaction
		if verifyTransaction(tx) {
			fmt.Println("Transaction verified:", tx)
			transactions = append(transactions, tx)
		}

		// If enough transactions (4), create a new block
		if len(transactions) >= 4 && !blockMined {
			prevHash := ""
			if len(Blockchain) > 0 {
				prevHash = Blockchain[len(Blockchain)-1].Hash
			}

			// Perform Proof-of-Work
			blockHash, nonce, success := proofOfWork(prevHash, transactions, difficulty)
			if success {
				newBlock := Block{
					Index:        len(Blockchain),
					Transactions: transactions,
					PrevHash:     prevHash,
					Hash:         blockHash,
					Nonce:        nonce,
					Timestamp:    time.Now(),
				}

				Blockchain = append(Blockchain, newBlock)
				fmt.Printf("Block mined successfully:\n %+v\n", newBlock)

				// Propagate block to other miners
				propagateBlock(newBlock, minerID)
				blockMined = true // Mark that the block has been mined
			}
			transactions = nil // Reset transactions
		}

		// Exit the loop after mining one block
		if blockMined {
			fmt.Println("One block mined and added to blockchain. Shutting down miner.")
			break
		}
	}
}

// Propagate the mined block to other miners
func propagateBlock(block Block, minerID string) {
	// List of other miners with the correct port numbers
	otherMiners := []string{
		"miner-node-2:8082", // Miner 2
		"miner-node-3:8083", // Miner 3
	}

	// Loop through the miners and propagate the block
	for _, minerAddr := range otherMiners {
		conn, err := net.Dial("tcp", minerAddr)
		if err != nil {
			fmt.Println("Error connecting to other miner:", minerAddr)
			continue
		}

		// Create an encoder and send the block to the miner
		encoder := json.NewEncoder(conn)
		err = encoder.Encode(block)
		if err != nil {
			fmt.Println("Error propagating block to miner:", minerAddr)
		}

		// Close the connection once the block is sent
		conn.Close()
	}
}

// Create the genesis block
func createGenesisBlock() Block {
	genesisTransactions := []Transaction{
		{DatasetHash: "genesis_dataset", AlgorithmHash: "genesis_algo", OutputHash: "genesis_output"},
	}
	genesisHash := calculateHash("0", genesisTransactions, 0) // PrevHash = "0", Nonce = 0
	return Block{
		Index:        0,
		Transactions: genesisTransactions,
		PrevHash:     "0",
		Hash:         genesisHash,
		Nonce:        0,
		Timestamp:    time.Now(),
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	minerID := os.Getenv("MINER_ID")
	if minerID == "" {
		minerID = "default-miner"
	}
	fmt.Println("Miner ID:", minerID)

	// Initialize blockchain with the genesis block
	genesisBlock := createGenesisBlock()
	Blockchain = append(Blockchain, genesisBlock)
	fmt.Printf("Genesis Block Created:\n %+v\n", genesisBlock)

	// Start the miner node
	minerNode(minerID)
}
