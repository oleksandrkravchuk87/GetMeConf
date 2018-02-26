package main

import (
	"log"
	"io"
	"encoding/json"
	"github.com/YAWAL/GetMeConf/database"
	"golang.org/x/net/context"
	"github.com/YAWAL/GetMeConf/api"
)

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
