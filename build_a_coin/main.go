package main

import (
	"github.com/chrispy-k/build_a_coin/cli"
	"github.com/chrispy-k/build_a_coin/db"
)

func main() {
	// defer lets the defered line get executed after the wrapping function finishes
	defer db.Close()
	cli.Start()
}
