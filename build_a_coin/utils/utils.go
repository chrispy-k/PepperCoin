package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// with the parameter like i interface{}, the function will receive anything
func ToBytes(i interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	HandleError(encoder.Encode(i))
	return buffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleError(decoder.Decode(i))
}
