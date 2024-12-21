package main

import (
	"blockchain2/nodes"        // Importing the node package
	"blockchain2/transactions" // Importing the transactions package
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func sendTransactionToAllMiners(transactions []transactions.Transaction, minerConnections []net.Conn) {
	for _, transaction := range transactions {
		for _, conn := range minerConnections {
			// Serialize the transaction to JSON
			transactionData, err := json.Marshal(transaction)
			if err != nil {
				log.Printf("Error marshalling transaction: %v", err)
				continue
			}

			// Send the serialized data over the connection
			_, err = conn.Write(transactionData)
			if err != nil {
				log.Printf("Error sending transaction to miner: %v", err)
				continue
			}
			fmt.Printf("Sent transaction: %+v to miner\n", transaction)
		}

	}

}

func main() {
	// Define the number of clusters (k) for KMeans algorithm
	k := 3

	// Create 4 nodes and create transactions
	node1 := &nodes.Node{ID: 1}
	node2 := &nodes.Node{ID: 2}
	node3 := &nodes.Node{ID: 3}
	node4 := &nodes.Node{ID: 4}

	node1.CreateTransaction(k)
	node2.CreateTransaction(k)
	node3.CreateTransaction(k)
	node4.CreateTransaction(k)

	// Print the transactions for each node after all nodes have completed
	fmt.Println("\nTransactions for all nodes:")
	fmt.Printf("Node 1 Transaction: %+v\n", node1.Transaction)
	fmt.Printf("Node 2 Transaction: %+v\n", node2.Transaction)
	fmt.Printf("Node 3 Transaction: %+v\n", node3.Transaction)
	fmt.Printf("Node 4 Transaction: %+v\n", node4.Transaction)

	// Miner addresses (example)
	minerAddresses := []string{"localhost:6000", "localhost:6001", "localhost:6002", "localhost:6003", "localhost:6004"}
	var minerConnections []net.Conn

	// Establish connections to miners
	for _, address := range minerAddresses {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Fatalf("Error connecting to miner %s: %v", address, err)
		}
		minerConnections = append(minerConnections, conn)
		defer conn.Close()
	}

	// Send transactions to miners
	sendTransactionToAllMiners([]transactions.Transaction{node1.Transaction[0], node2.Transaction[0], node3.Transaction[0], node4.Transaction[0]}, minerConnections)
}
