package blockchain

import (
	"sync"

	"github.com/chrispy-k/build_a_coin/db"
	"github.com/chrispy-k/build_a_coin/utils"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var b *blockchain
var once sync.Once

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp - lastRecalculatedBlock.Timestamp) / 60
	expectedTime := difficultyInterval * blockInterval

	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func Blockchain() *blockchain {
	// want to run this if block of code ONCE
	// no matter how many program, threads are running
	// that can be achieved with GO Sync
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			// search for checkpoint on db
			// restore blockchain from bytes
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock()
			} else {
				b.restore(checkpoint)
			}

		})
	}
	return b
}

func (b *blockchain) txOuts() []*TxOut {
	blocks := b.Blocks()
	var txOuts []*TxOut
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
}

func (b *blockchain) TxOutByAddress(address string) []*TxOut {
	var ownedTxOuts []*TxOut
	txOuts := b.txOuts()
	for _, txOut := range txOuts {
		if txOut.Owner == address {
			ownedTxOuts = append(ownedTxOuts, txOut)
		}
	}
	return ownedTxOuts
}

func (b *blockchain) BalanceByAddresss(address string) int {
	txOuts := b.TxOutByAddress(address)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}
