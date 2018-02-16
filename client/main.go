package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/YAWAL/GetMeConf/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const address = "localhost:8081"

var (
	//client     api.ConfigServiceClient
	configName = flag.String("config-name", "", "config name")
	configType = flag.String("config-type", "", "config type")
	outPath = flag.String("outpath", "", "output path for config file")
)

func main() {


	//serverPort := flag.String("server", "localhost:50111", "port for connection to server")
	flag.Parse()

	if *configName == "" && *configType == "" {
		log.Fatal("Can't proccess => config name and config type are empty")
	}

	log.Printf("Start checking input data:\n Config name: %v\n Config typy : %v\n Output path: %v", *configName, *configType, *outPath)
	log.Printf("Processing ...")

	conn, err := grpc.Dial(address, grpc.WithInsecure(),grpc.WithBlock())
	defer conn.Close()
	log.Printf("State: %v", conn.GetState())

	if err != nil {
		log.Fatalf("DialContext error has occurred: %v", err)
	}

	client := api.NewConfigServiceClient(conn)
	log.Printf("Processing client...")

	if *configName != "" && *configType != "" {
		log.Printf("Processing retrieveConfig...")

		err := retrieveConfig(configName, outPath, client)
		if err != nil {
			log.Fatalf("!!!: %v", err)
		}
	}

	log.Printf("End retrieveConfig...")

	//for true {
	//}
}

func retrieveConfig(fileName, outputPath *string, client api.ConfigServiceClient) error{
	conf, err := client.GetConfigByName(context.Background(), &api.GetConfigByNameRequest{ConfigName: *configName, ConfigType: *configType})
	if err != nil {
		log.Fatalf("Error during retrieving config has occurred: %v", err)
		return err
	}
	if err := WriteFile(conf.Config, *fileName, *outputPath); err != nil {
		log.Fatalf("Error during writing file in retrieving config: %v", err)
		return err
	}
	return nil
}

func WriteFile(data []byte, fileName, outPath string) error {
	fileName = fileName + ".json"
	if err := ioutil.WriteFile(filepath.Join(outPath, fileName ), data, 0666); err != nil {
		log.Fatalf("Error during file creation: %v", err)
		return err
	} else {
		log.Printf("File %v has been created in %v", fileName, outPath)
		return nil
	}
}
