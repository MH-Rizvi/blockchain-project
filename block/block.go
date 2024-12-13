package block

import (
	"blockchain/transaction"
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Hash         string
	PreviousHash string
	Transactions []transaction.Transaction
	Nonce        int
}

func NewGenesisBlock() Block {
	return Block{
		Hash:         "genesis-block-hash",
		PreviousHash: "",
		Transactions: nil,
		Nonce:        0,
	}
}

func NewBlock(previousHash string, transactions []transaction.Transaction) Block {
	return Block{
		Hash:         "",
		PreviousHash: previousHash,
		Transactions: transactions,
		Nonce:        0,
	}
}

func (b *Block) CalculateHash(nonce int) string {
	data := fmt.Sprintf("%v-%v-%d", b.PreviousHash, b.Transactions, nonce)
	hash := sha256.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
