package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Transaction structure
type Transaction struct {
	DatasetHash   string `json:"dataset_hash"`
	AlgorithmHash string `json:"algorithm_hash"`
	OutputHash    string `json:"output_hash"`
}

// LoadTransaction fetches dataset and algorithm, runs them, and creates an output hash
func LoadTransaction() Transaction {
	// Simulate fetching data from IPFS
	datasetHash := "ipfs_dataset_hash_placeholder"
	algorithmHash := "ipfs_algorithm_hash_placeholder"

	// Simulate running the algorithm on the dataset to produce an output hash
	outputHash := calculateOutputHash(datasetHash, algorithmHash)

	return Transaction{
		DatasetHash:   datasetHash,
		AlgorithmHash: algorithmHash,
		OutputHash:    outputHash,
	}
}

// calculateOutputHash simulates hashing dataset and algorithm to create an output
func calculateOutputHash(dataset, algorithm string) string {
	// Simple hash calculation placeholder (use actual cryptographic hash in production)
	return fmt.Sprintf("%x", (len(dataset)+len(algorithm))^0xabc)
}

// sendTransactionToMiner sends a transaction to a miner node
func sendTransactionToMiner(transaction Transaction, minerAddress string) error {
	url := fmt.Sprintf("http://%s/addTransaction", minerAddress)

	// Serialize the transaction
	data, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// Send the transaction via HTTP POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("miner responded with status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Fetch user node environment variable
	userID := os.Getenv("USER_ID")
	if userID == "" {
		userID = "default-user"
	}

	// Define the miner addresses (replace with actual miner IPs/ports in Docker)
	minerAddresses := []string{
		"miner-node-1:8081", // Miner 1
		"miner-node-2:8082", // Miner 2
		"miner-node-3:8083", // Miner 3
	}

	// User node lifecycle: Send exactly 4 transactions
	for i := 0; i < 4; i++ { // Send exactly 4 transactions
		// Create a new transaction
		transaction := LoadTransaction()
		fmt.Printf("User %s created transaction: %+v\n", userID, transaction)

		// Send the transaction to all miners
		for _, minerAddress := range minerAddresses {
			err := sendTransactionToMiner(transaction, minerAddress)
			if err != nil {
				fmt.Printf("Error sending transaction to miner %s: %v\n", minerAddress, err)
			} else {
				fmt.Printf("Transaction successfully sent to miner %s\n", minerAddress)
			}
		}

		// Wait before creating the next transaction
		time.Sleep(10 * time.Second)
	}

	// After sending 4 transactions, stop the user node
	fmt.Println("Sent 4 transactions. Shutting down user node.")
}
