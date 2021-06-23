package blockchain

import (
	"sync"

	"github.com/chrispy-k/build_a_coin/db"
	"github.com/chrispy-k/build_a_coin/utils"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height`
}

var b *blockchain
var once sync.Once

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
}

func (b *blockchain) FromBytes(data []byte) {

}

func Blockchain() *blockchain {
	// want to run this if block of code ONCE
	// no matter how many program, threads are running
	// that can be achieved with GO Sync
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			// search for checkpoint on db
			// restore blockchain from bytes
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				b.FromBytes(checkpoint)
			}

		})
	}
	return b
}
