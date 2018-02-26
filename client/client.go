package main

import (
	"flag"
	"log"

	"github.com/YAWAL/GetMeConf/api"

	"google.golang.org/grpc"
)

const address = "localhost:3000"

var (
	configName = flag.String("config-name", "", "config name")
	configType = flag.String("config-type", "", "config type")
	outPath    = flag.String("outpath", "", "output path for config file")
)

func main() {

	flag.Parse()

	if *configName == "" && *configType == "" {
		log.Fatal("Can't proccess => config name and config type are empty")
	}

	log.Printf("Start checking input data:\n Config name: %v\n Config type : %v\n Output path: %v\nProcessing ...", *configName, *configType, *outPath)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	log.Printf("State: %v", conn.GetState())

	if err != nil {
		log.Fatalf("Dial error has occurred: %v", err)
	}

	client := api.NewConfigServiceClient(conn)
	log.Printf("Processing client...")

	if *configName != "" && *configType != "" {
		log.Printf("Processing retrieving config...")

		err := retrieveConfig(configName, outPath, client)
		if err != nil {
			log.Fatalf("retrieveConfig err: %v", err)
		}
	}

	if *configName == "" && *configType != "" {
		err := retrieveConfigs(client)
		if err != nil {
			log.Fatalf("retrieveConfigs err : %v", err)
		}
	}
	log.Printf("End retrieving config.")

}
