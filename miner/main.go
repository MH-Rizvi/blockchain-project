package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

// Block structure
type Block struct {
	BlockIndex   int      // Added BlockIndex field
	PrevHash     string   // Previous block's hash
	OutputHash   string   // Hash of the transaction output
	Nonce        int      // Nonce for PoW
	Hash         string   // Hash of the block
	Transactions []string // Transactions in the block
}

// Blockchain structure
var blockchain []Block

// Transaction Queue
var transactionQueue []string

// Genesis block function
func generateGenesisBlock() Block {
	return Block{
		BlockIndex:   0,                  // Genesis block index is 0
		PrevHash:     "0",                // No previous hash for genesis
		OutputHash:   "genesis",          // Arbitrary hash for the first block
		Nonce:        0,                  // No nonce for genesis
		Hash:         "genesisblockhash", // A placeholder for the genesis block hash
		Transactions: []string{},
	}
}

// PoW function to mine a new block
func proofOfWork(prevHash string, transactions []string, blockIndex int) Block {
	nonce := 0
	var hash string
	target := "0" // We want hashes that start with '0'

	// Combine transactions into a single output hash for the block
	outputHash := hashData(transactions)

	for {
		hash = calculateBlockHash(prevHash, outputHash, nonce, blockIndex)

		// Print each attempt during the mining process
		fmt.Printf("Attempting nonce %d: %s\n", nonce, hash)

		if strings.HasPrefix(hash, target) {
			break
		}
		nonce++
	}

	// Create the block with the found nonce
	return Block{
		BlockIndex:   blockIndex,
		PrevHash:     prevHash,
		OutputHash:   outputHash,
		Nonce:        nonce,
		Hash:         hash,
		Transactions: transactions, // Add the 4 transactions to the block
	}
}

// Calculate block hash (using SHA256)
func calculateBlockHash(prevHash, outputHash string, nonce, blockIndex int) string {
	blockString := fmt.Sprintf("%d%s%s%d", blockIndex, prevHash, outputHash, nonce)
	hash := sha256.New()
	hash.Write([]byte(blockString))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Hash the data (general-purpose hash function)
func hashData(data interface{}) string {
	dataBytes, _ := json.Marshal(data)
	hash := sha256.New()
	hash.Write(dataBytes)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Handle transaction from user node
func handleTransaction(conn net.Conn) {
	// Read the transaction from user node
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal("Failed to read from user node:", err)
	}
	transactionData := string(buf[:n]) // Read only the valid part of the buffer

	log.Printf("Received transaction: %s\n", transactionData)

	// Parse the transaction JSON
	var transaction map[string]string
	err = json.Unmarshal([]byte(transactionData), &transaction)
	if err != nil {
		log.Fatal("Failed to parse transaction:", err)
	}

	// Add the transaction to the queue
	transactionQueue = append(transactionQueue, transactionData)

	// Once 4 transactions are received, mine a block
	if len(transactionQueue) == 4 {
		// Perform PoW to create a new block
		prevHash := blockchain[len(blockchain)-1].Hash
		blockIndex := len(blockchain) // Increment the block index for the new block
		block := proofOfWork(prevHash, transactionQueue, blockIndex)
		blockchain = append(blockchain, block)

		// Log the block mined
		log.Printf("Block mined! Hash: %s\n", block.Hash)
		log.Printf("Blockchain: %+v\n", blockchain)

		// Reset the transaction queue after mining
		transactionQueue = nil
	}
}

// Set up the miner node listener
func startMiner() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Error starting miner node:", err)
	}
	defer ln.Close()

	log.Println("Miner node listening on port 8081...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err)
		}
		go handleTransaction(conn)
	}
}

func main() {
	// Start the miner node to listen for incoming connections
	// Initialize the blockchain with the genesis block
	blockchain = append(blockchain, generateGenesisBlock())

	// Start the miner node to handle incoming transactions
	startMiner()
}
