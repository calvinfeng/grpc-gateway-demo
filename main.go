package main

import (
	"github.com/calvinfeng/grpc-gateway-demo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
