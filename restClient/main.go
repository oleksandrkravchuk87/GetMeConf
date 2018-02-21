package main

import (
	"encoding/json"
	"log"

	"net/http"

	"io"

	"github.com/YAWAL/GetMeConf/api"
	"github.com/YAWAL/GetMeConf/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const address = "localhost:3000"

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	conn.GetState()
	defer conn.Close()
	log.Printf("State: %v", conn.GetState())

	if err != nil {
		log.Fatalf("DialContext error has occurred: %v", err)
	}

	client := api.NewConfigServiceClient(conn)
	log.Printf("Processing client...")

	//http server
	hs := gin.Default()
	hs.GET("/getConfig/:type/:name", func(c *gin.Context) {
		configType := c.Param("type")
		configName := c.Param("name")
		resultConfig, err := retrieveConfig(&configName, &configType, client)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"config": resultConfig,
		})
	})

	hs.GET("/getConfig/:type/", func(c *gin.Context) {
		configType := c.Param("type")
		resultConfig, err := retrieveConfigs(&configType, client)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"config": resultConfig,
		})
	})

	if err := hs.Run(":8080"); err != nil {
		log.Fatalf("filed to run server: %v", err)
	}

}

func retrieveConfig(configName, configType *string, client api.ConfigServiceClient) (database.ConfigInterface, error) {
	config, err := client.GetConfigByName(context.Background(), &api.GetConfigByNameRequest{ConfigName: *configName, ConfigType: *configType})
	if err != nil {
		log.Printf("Error during retrieving config has occurred: %v", err)
		return nil, err
	}
	switch *configType {
	case "mongodb":
		var mongodb database.Mongodb
		err := json.Unmarshal(config.Config, &mongodb)
		if err != nil {
			log.Printf("Unmarshal mongodb err: %v", err)
			return nil, err
		}
		return mongodb, err
	case "tempconfig":
		var tempconfig database.Tempconfig
		err := json.Unmarshal(config.Config, &tempconfig)
		if err != nil {
			log.Printf("Unmarshal tempconfig err: %v", err)
			return nil, err
		}
		return tempconfig, err
	case "tsconfig":
		var tsconfig database.Tsconfig
		err := json.Unmarshal(config.Config, &tsconfig)
		if err != nil {
			log.Printf("Unmarshal tsconfig err: %v", err)
			return nil, err
		}
		return tsconfig, err
	default:
		log.Printf("Such config: %v does not exist", *configType)
		return nil, err
	}
}

func retrieveConfigs(configType *string, client api.ConfigServiceClient) ([]database.ConfigInterface, error) {
	stream, err := client.GetConfigsByType(context.Background(), &api.GetConfigsByTypeRequest{ConfigType: *configType})
	if err != nil {
		log.Printf("Error during retrieving stream configs has occurred:%v", err)
		return nil, err
	}
	var resultConfigs []database.ConfigInterface
	for {
		config, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error during streaming has occurred: %v", err)
			return nil, err
		}
		switch *configType {
		case "mongodb":
			var mongodb database.Mongodb
			err := json.Unmarshal(config.Config, &mongodb)
			if err != nil {
				log.Printf("Unmarshal mongodb err: %v", err)
				return nil, err
			}
			resultConfigs = append(resultConfigs, mongodb)
		case "tempconfig":
			var tempconfig database.Tempconfig
			err := json.Unmarshal(config.Config, &tempconfig)
			if err != nil {
				log.Printf("Unmarshal tempconfig err: %v", err)
				return nil, err
			}
			resultConfigs = append(resultConfigs, tempconfig)
		case "tsconfig":
			var tsconfig database.Tsconfig
			err := json.Unmarshal(config.Config, &tsconfig)
			if err != nil {
				log.Printf("Unmarshal tsconfig err: %v", err)
				return nil, err
			}
			resultConfigs = append(resultConfigs, tsconfig)
		default:
			log.Printf("Such config: %v does not exist", *configType)
			return nil, err
		}
	}
	return resultConfigs, nil
}
