package main

import (
	"encoding/json"
	"os"
	Server "proj1/server"
	"strconv"
	"fmt"
)

func main() {
	if len(os.Args) < 2 {
		config := Server.Config{Encoder: json.NewEncoder(os.Stdout), Decoder: json.NewDecoder(os.Stdin), Mode: "s", ConsumersCount: 1}
		Server.Run(config)
	} else if len(os.Args) == 2{
		ConsCnt, _ := strconv.Atoi(os.Args[1])
		config := Server.Config{Encoder: json.NewEncoder(os.Stdout), Decoder: json.NewDecoder(os.Stdin), Mode: "p", ConsumersCount: ConsCnt}
		Server.Run(config)
	} else {
		fmt.Println("Usage: twitter <number of consumers>\n<number of consumers> = the number of goroutines (i.e., consumers) to be part of the parallel version.")
	}
}
