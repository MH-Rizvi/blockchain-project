package nodes

import (
	"blockchain2/transactions"
)

// Node structure representing each user node
type Node struct {
	ID          int
	Transaction []transactions.Transaction
}

// CreateTransaction generates a transaction for the node
func (n *Node) CreateTransaction(k int) {
	// Generate transaction and store it in the node's Transaction field

	transaction := transactions.CreateTransaction(k)
	n.Transaction = append(n.Transaction, transaction)
}
