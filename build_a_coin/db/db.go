package db

import (
	"github.com/boltdb/bolt"
	"github.com/chrispy-k/build_a_coin/utils"
)

// this is not exported
var db *bolt.DB

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
)

// this is the function that is exported and that will initialize the database
func DB() *bolt.DB {
	if db == nil {
		// 0600 for handling read/write permission
		dbPointer, err := bolt.Open("blockchain.db", 0600, nil)
		db = dbPointer
		utils.HandleError(err)
		err = db.Update(func(t *bolt.Tx) error {
			// one bucker for data, one bucket for blocks
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleError(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleError(err)
	}
	return db
}
