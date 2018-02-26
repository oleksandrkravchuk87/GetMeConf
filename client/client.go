package main

import (
	"flag"
	"log"
)

const address = "localhost:3000"

var (
	configName = flag.String("config-name", "", "config name")
	configType = flag.String("config-type", "", "config type")
	outPath    = flag.String("outpath", "", "output path for config file")
	command    = flag.String("command", "", "1 of CRUD operations")
)

func main() {

	flag.Parse()

	log.Printf("Processing client...Reading flags")

	switch *command {
	case "create":

	case "read":
		ConfigRetriever()
	case "update":

		//TODO implement

	case "delete":

	default:
		log.Fatalf("Cant parse propper flag. Valid flags are: create, read, update, delete")

	}

}
