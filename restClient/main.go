package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"os"

	"github.com/YAWAL/GetMeConf/api"
	"github.com/YAWAL/GetMeConf/database"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serviceHost := os.Getenv("SERVICEHOST")
	if port == "" {
		serviceHost = "localhost"
	}
	servicePort := os.Getenv("SERVICEPORT")
	if servicePort == "" {
		servicePort = "3000"
	}

	address := fmt.Sprintf("%s:%s", serviceHost, servicePort)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	conn.GetState()
	log.Printf("State: %v", conn.GetState())
	defer conn.Close()
	if err != nil {
		log.Fatalf("DialContext error has occurred: %v", err)
	}

	client := api.NewConfigServiceClient(conn)
	log.Printf("Processing client...")

	//http server
	router := gin.Default()
	router.GET("/getConfig/:type/:name", func(c *gin.Context) {
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

	router.GET("/getConfig/:type", func(c *gin.Context) {
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

	router.POST("/createConfig/:type", func(c *gin.Context) {
		configType := c.Param("type")
		confTypeStruct, _ := selectType(configType)
		var bytes []byte
		c.Bind(&confTypeStruct)
		bytes, err = json.Marshal(confTypeStruct)
		result, err := client.CreateConfig(context.Background(), &api.Config{ConfigType: configType, Config: bytes})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"config": result,
		})
	})

	router.DELETE("/deleteConfig/:type/:name", func(c *gin.Context) {
		configType := c.Param("type")
		configName := c.Param("name")
		deleteResult, err := client.DeleteConfig(context.Background(), &api.DeleteConfigRequest{configName, configType})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"config": deleteResult,
		})
	})

	router.PUT("/updateConfig/:type", func(c *gin.Context) {
		configType := c.Param("type")
		confTypeStruct, _ := selectType(configType)
		var bytes []byte
		c.Bind(&confTypeStruct)
		bytes, err = json.Marshal(confTypeStruct)
		updateResult, err := client.UpdateConfig(context.Background(), &api.Config{ConfigType: configType, Config: bytes})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"config": updateResult,
		})
	})

	router.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusText(http.StatusOK),
		})
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	defer srv.Shutdown(context.Background())
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("filed to run server: %v", err)
	}

}

func selectType(cType string) (database.ConfigInterface, error) {
	switch cType {
	case "mongodb":
		return new(database.Mongodb), nil
	case "tempconfig":
		return new(database.Tempconfig), nil
	case "tsconfig":
		return new(database.Tsconfig), nil
	default:
		log.Printf("Such config: %v does not exist", cType)
		return nil, errors.New("config does not exist")
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
		return nil, errors.New("config does not exist")
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
