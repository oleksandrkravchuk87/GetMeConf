package main

import (
	"flag"
	"log"
)

const address = "localhost:3000"

var (
	command    = flag.String("command", "", "1 of CRUD operations")
	configName = flag.String("config-name", "", "config name")
	configType = flag.String("config-type", "", "config type")
	outPath    = flag.String("outpath", "", "output path for config file")
	fileName   = flag.String("file-name", "", "config file's name")
)

func main() {

	flag.Parse()

	log.Println("Processing client...Reading flags")

	switch *command {
	case "create":
		log.Printf("Command: %v /n Config name: %v /n Config type: %v", *command, *configName, *configType)
		sentConfigToServer(*fileName)

	case "read":
		ConfigRetriever()
	case "update":

		//TODO implement

	case "delete":

	default:
		log.Fatalf("Cant parse command flag. Valid command flags are: create, read, update, delete")

	}

}
