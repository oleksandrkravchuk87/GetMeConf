package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	"encoding/json"
	"io"

	"github.com/YAWAL/GetMeConf/api"
	"github.com/YAWAL/GetMeConf/database"
	"golang.org/x/net/context"
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

func retrieveConfig(fileName, outputPath *string, client api.ConfigServiceClient) error {
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

func retrieveConfigs(client api.ConfigServiceClient) error {
	stream, err := client.GetConfigsByType(context.Background(), &api.GetConfigsByTypeRequest{ConfigType: *configType})
	if err != nil {
		log.Fatalf("Error during retrieving stream configs has occurred:%v", err)
	}
	for {
		config, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error during streaming has occurred: %v", err)
			return err
		}

		switch *configType {
		case "mongodb":
			var mongodb database.Mongodb
			err := json.Unmarshal(config.Config, &mongodb)
			if err != nil {
				log.Fatalf("Unmarshal mongodb err: %v", err)
			}
			fileName := mongodb.Domain
			WriteFile(config.Config, fileName, *outPath)

		case "tempconfig":
			var tempconfig database.Tempconfig
			err := json.Unmarshal(config.Config, &tempconfig)
			if err != nil {
				log.Fatalf("Unmarshal tempconfig err: %v", err)
			}
			fileName := tempconfig.RestApiRoot
			WriteFile(config.Config, fileName, *outPath)

		case "tsconfig":
			var tsconfig database.Tsconfig
			err := json.Unmarshal(config.Config, &tsconfig)
			if err != nil {
				log.Fatalf("Unmarshal tsconfig err: %v", err)
			}
			fileName := tsconfig.Module
			WriteFile(config.Config, fileName, *outPath)

		default:
			log.Fatalf("Such config: %v does not exist", *configType)
		}
	}
	return nil
}

func WriteFile(data []byte, fileName, outPath string) error {
	fileName = fileName + ".json"
	if err := ioutil.WriteFile(filepath.Join(outPath, fileName), data, 0666); err != nil {
		log.Fatalf("Error during file creation: %v", err)
		return err
	} else {
		log.Printf("File %v has been created in %v", fileName, outPath)
		return nil
	}
}
