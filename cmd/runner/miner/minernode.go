package main

import (
	"blockchain2/transactions"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Miner struct {
	ID          int
	Block       *Block
	BlockMutex  sync.Mutex
	StopMining  chan bool
	TargetHash  string
	Verified    bool
	FoundHash   string
	FoundNonce  int
	Transaction []transactions.Transaction
	OtherMiners []*Miner
	Blockchain  []*Block
}

type Block struct {
	Index        int
	Transactions []transactions.Transaction
	PrevHash     string
	Nonce        int
	Timestamp    time.Time
	Hash         string
}

var stopMiningSignal = make(chan bool, 1)
var miners []*Miner

func (m *Miner) AddBlockToBlockchain() {
	m.Block.Hash = m.FoundHash
	m.Block.Timestamp = time.Now()
	m.Blockchain = append(m.Blockchain, m.Block)
	fmt.Printf("Miner %d added block to blockchain: %+v\n", m.ID, m.Block)
}

func (m *Miner) AddTransactionToBlock(transaction transactions.Transaction) int {
	m.BlockMutex.Lock()
	defer m.BlockMutex.Unlock()

	m.Block.Transactions = append(m.Block.Transactions, transaction)
	fmt.Printf("Miner %d added transaction: %+v\n", m.ID, transaction)
	length := len(m.Block.Transactions)
	if length == 4 {
		fmt.Printf("Miner %d has 4 transactions. Ready to mine!\n", m.ID)
	}
	return length
}

func CalculateHash(block *Block) string {
	blockData := fmt.Sprintf("%v%v%v%v", block.Transactions, block.PrevHash, block.Nonce, block.Timestamp)
	hash := sha256.New()
	hash.Write([]byte(blockData))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (m *Miner) MineBlock() {
	fmt.Printf("Miner %d started mining...\n", m.ID)

	var hash string
	var nonce int
	for {
		select {
		case <-stopMiningSignal:
			fmt.Printf("Miner %d stopped mining due to a valid hash found by another miner.\n", m.ID)
			return
		default:
			nonce++
			m.Block.Nonce = nonce
			hash = CalculateHash(m.Block)
			if hash[0] == '0' {
				m.FoundHash = hash
				m.FoundNonce = nonce
				fmt.Printf("Miner %d found hash: %s with nonce %d\n", m.ID, m.FoundHash, m.FoundNonce)
				m.BroadcastFoundHashAndNonce()
				stopMiningSignal <- true
				return
			}
		}
	}
}

func (m *Miner) BroadcastFoundHashAndNonce() {
	fmt.Printf("Miner %d broadcasting found hash: %s and nonce: %d\n", m.ID, m.FoundHash, m.FoundNonce)
	for _, otherMiner := range m.OtherMiners {
		go otherMiner.VerifyBlock(m.FoundHash, m.FoundNonce)
	}
}

func (m *Miner) VerifyBlock(foundHash string, foundNonce int) {
	m.Block.Nonce = foundNonce
	calculatedHash := CalculateHash(m.Block)
	if calculatedHash == foundHash && foundHash[0] == '0' {
		fmt.Printf("Miner %d successfully verified the block with hash: %s\n", m.ID, foundHash)
		m.Verified = true
		m.AddBlockToBlockchain()
	} else {
		fmt.Printf("Miner %d failed to verify the block\n", m.ID)
	}
}

func (m *Miner) HandleTransactions(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Miner %d: Error reading from connection: %v", m.ID, err)
			return
		}

		var transaction transactions.Transaction
		err = json.Unmarshal(buffer[:n], &transaction)
		if err != nil {
			log.Printf("Miner %d: Error unmarshalling transaction: %v", m.ID, err)
			continue
		}

		length := m.AddTransactionToBlock(transaction)
		if length == 4 {
			fmt.Printf("Miner %d returned to HandleTransactions\n\n", m.ID)
			m.MineBlock()
		}
	}
}

func (m *Miner) StartMiner(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Miner %d: Error starting server on port %s: %v", m.ID, port, err)
	}
	defer listener.Close()

	fmt.Printf("Miner %d is listening for transactions on port %s...\n", m.ID, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Miner %d: Error accepting connection: %v", m.ID, err)
			continue
		}
		go m.HandleTransactions(conn)
	}
}

func main() {
	ports := []string{"6000", "6001", "6002", "6003", "6004"}
	targetHash := "0"

	for i := 1; i <= 5; i++ {
		miner := &Miner{
			ID: i,
			Block: &Block{
				Index:        0,
				Transactions: []transactions.Transaction{},
				Timestamp:    time.Now(),
				PrevHash:     "",
			},
			StopMining: make(chan bool),
			TargetHash: targetHash,
		}
		miners = append(miners, miner)
	}

	for i := 0; i < len(miners); i++ {
		for j := 0; j < len(miners); j++ {
			if i != j {
				miners[i].OtherMiners = append(miners[i].OtherMiners, miners[j])
			}
		}
	}

	started := make(chan bool)

	var wg sync.WaitGroup
	for i, miner := range miners {
		wg.Add(1)
		go func(m *Miner, port string) {
			defer wg.Done()
			m.StartMiner("localhost:" + port)
			started <- true
		}(miner, ports[i])
	}

	for i := 0; i < len(miners); i++ {
		<-started
	}

	fmt.Printf("Back here\n")
}
