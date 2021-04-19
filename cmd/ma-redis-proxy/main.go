package main

import (
	"fmt"
	"os"

	"github.com/moonactive/ma-redis-proxy/pkg/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
