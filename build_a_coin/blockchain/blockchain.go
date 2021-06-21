package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

// singleton pattern
// we want to always be sharing only one instance of the blockchain

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

// this is called the receiver function
func (b *Block) calculateHash() {
	// b here is referring to the b that we received
	hash := sha256.Sum256([]byte(b.Hash + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	length := len(GetBlockchain().blocks)
	if length == 0 {
		return ""
	}
	return GetBlockchain().blocks[length-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks) + 1}
	newBlock.calculateHash()
	return &newBlock
}

//
func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

var ErrNotFound = errors.New("block not found")

func (b *blockchain) GetBlock(height int) (*Block, error) {
	if height > len(b.blocks) {
		return nil, ErrNotFound
	}
	return b.blocks[height-1], nil
}

func GetBlockchain() *blockchain {
	// want to run this if block of code ONCE
	// no matter how many program, threads are running
	// that can be achieved with GO Sync
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}
