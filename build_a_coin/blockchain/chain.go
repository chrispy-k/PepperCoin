package blockchain

import (
	"sync"

	"github.com/chrispy-k/build_a_coin/db"
	"github.com/chrispy-k/build_a_coin/utils"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
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

func (b *blockchain) UTxOutByAddress(address string) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)

	for _, block := range b.Blocks() {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}

	return uTxOuts
}

func (b *blockchain) BalanceByAddresss(address string) int {
	txOuts := b.UTxOutByAddress(address)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}
