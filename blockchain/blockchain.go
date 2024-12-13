package blockchain

import (
	"blockchain/block"
	"blockchain/transaction"
	"fmt"
)

type Blockchain struct {
	Chain               []block.Block
	CurrentTransactions []transaction.Transaction
}

func NewBlockchain() *Blockchain {
	// Create a genesis block
	genesisBlock := block.NewGenesisBlock()
	return &Blockchain{
		Chain:               []block.Block{genesisBlock},
		CurrentTransactions: []transaction.Transaction{},
	}
}

func (bc *Blockchain) AddTransaction(t transaction.Transaction) {
	bc.CurrentTransactions = append(bc.CurrentTransactions, t)
}

func (bc *Blockchain) AddBlock(b block.Block) {
	bc.Chain = append(bc.Chain, b)
}

func (bc *Blockchain) VerifyTransaction(t transaction.Transaction) bool {
	// Simulate transaction verification (e.g., compare hashes of dataset)
	return true // Simplified for now
}

func (bc *Blockchain) GetLastBlock() block.Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) PrintBlockchain() {
	fmt.Println("Blockchain:")
	for i, b := range bc.Chain {
		fmt.Printf("Block #%d: %v\n", i, b)
	}
}
