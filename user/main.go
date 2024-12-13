package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// AI deterministic function: simple K-means (just for illustration)
func kMeans(dataset [][]float64) []float64 {
	var sum []float64
	for _, data := range dataset {
		if len(sum) == 0 {
			sum = make([]float64, len(data))
		}
		for i, val := range data {
			sum[i] += val
		}
	}
	for i := range sum {
		sum[i] /= float64(len(dataset))
	}
	return sum
}

// Create Transaction: Dataset, Algorithm, Output Hash
func createTransaction() (string, string, string) {
	// Example dataset (e.g., 2D points for K-means)
	dataset := [][]float64{
		{1.0, 2.0},
		{3.0, 4.0},
		{5.0, 6.0},
	}

	// Algorithm (simple K-means here)
	algorithm := "K-means"

	// Get output from AI function (deterministic)
	output := kMeans(dataset)

	// Hash the dataset, algorithm, and output
	datasetHash := hashData(dataset)
	algorithmHash := hashData([]string{algorithm})
	outputHash := hashData(output)

	// Return the transaction as JSON
	transaction := map[string]string{
		"datasetHash":   datasetHash,
		"algorithmHash": algorithmHash,
		"outputHash":    outputHash,
	}

	transactionJSON, _ := json.Marshal(transaction)
	return string(transactionJSON), datasetHash, outputHash
}

// Hash the data (general-purpose hash function)
func hashData(data interface{}) string {
	dataBytes, _ := json.Marshal(data)
	hash := sha256.New()
	hash.Write(dataBytes)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func sendTransaction(transaction string) {
	conn, err := net.Dial("tcp", "miner-node:8081") // Ensure this is the correct address of the miner node
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer conn.Close()

	// Send the transaction properly
	_, err = conn.Write([]byte(transaction + "\n")) // Adding a newline as delimiter
	if err != nil {
		log.Fatal("Failed to send transaction:", err)
	}

	log.Println("Transaction sent to miner node")
}

func main() {
	// Create a transaction
	transaction, datasetHash, outputHash := createTransaction()

	// Log the transaction for debugging
	log.Printf("Transaction created: %s\n", transaction)
	log.Printf("Dataset Hash: %s\n", datasetHash)
	log.Printf("Output Hash: %s\n", outputHash)

	// Send the transaction to the miner node
	sendTransaction(transaction)
}
