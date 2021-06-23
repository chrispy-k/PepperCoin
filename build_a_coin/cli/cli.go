package cli

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/chrispy-k/build_a_coin/explorer"
	"github.com/chrispy-k/build_a_coin/rest"
)

func usage() {
	fmt.Printf("Welcome to PepperCoin\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port=4000:		Sets the PORT of the server\n")
	fmt.Printf("-mode=rest:		Chooses betweeen 'html' and 'rest\n\n")
	runtime.Goexit()
}

func Start() {
	// we put in default values so we dont care if there is no value or not
	// if len(os.Args) < 2 {
	// 	usage()
	// }

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose betweeen 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}

}
