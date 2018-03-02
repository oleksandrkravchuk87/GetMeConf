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

const mongoConf = "mongodb"
const tsConf = "tsconfig"
const tempConf = "tempconfig"

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
	if err != nil {
		log.Fatalf("dialContext error has occurred: %v", err)
	}
	conn.GetState()
	log.Printf("State: %v", conn.GetState())
	defer conn.Close()

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
		createResult, err := createConfig(c, client)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"config": createResult,
		})
	})

	router.DELETE("/deleteConfig/:type/:name", func(c *gin.Context) {
		configType := c.Param("type")
		configName := c.Param("name")
		deleteResult, err := client.DeleteConfig(context.Background(), &api.DeleteConfigRequest{ConfigName: configName, ConfigType: configType})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"config": deleteResult,
		})
	})

	router.PUT("/updateConfig/:type", func(c *gin.Context) {
		updateResult, err := updateConfig(c, client)
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
	case mongoConf:
		return new(database.Mongodb), nil
	case tempConf:
		return new(database.Tempconfig), nil
	case tsConf:
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
	case mongoConf:
		var mongodb database.Mongodb
		err := json.Unmarshal(config.Config, &mongodb)
		if err != nil {
			log.Printf("Unmarshal mongodb err: %v", err)
			return nil, err
		}
		return mongodb, err
	case tempConf:
		var tempconfig database.Tempconfig
		err := json.Unmarshal(config.Config, &tempconfig)
		if err != nil {
			log.Printf("Unmarshal tempconfig err: %v", err)
			return nil, err
		}
		return tempconfig, err
	case tsConf:
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
		case mongoConf:
			var mongodb database.Mongodb
			err := json.Unmarshal(config.Config, &mongodb)
			if err != nil {
				log.Printf("Unmarshal mongodb err: %v", err)
				return nil, err
			}
			resultConfigs = append(resultConfigs, mongodb)
		case tempConf:
			var tempconfig database.Tempconfig
			err := json.Unmarshal(config.Config, &tempconfig)
			if err != nil {
				log.Printf("Unmarshal tempconfig err: %v", err)
				return nil, err
			}
			resultConfigs = append(resultConfigs, tempconfig)
		case tsConf:
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

func updateConfig(c *gin.Context, client api.ConfigServiceClient) (*api.Responce, error) {
	configType := c.Param("type")
	confTypeStruct, err := selectType(configType)
	if err != nil {
		return nil, err
	}
	var bytes []byte
	if err = c.Bind(&confTypeStruct); err != nil {
		return nil, err
	}
	bytes, err = json.Marshal(confTypeStruct)
	if err != nil {
		return nil, err
	}
	updateResult, err := client.UpdateConfig(context.Background(), &api.Config{ConfigType: configType, Config: bytes})
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func createConfig(c *gin.Context, client api.ConfigServiceClient) (*api.Responce, error) {
	configType := c.Param("type")
	confTypeStruct, err := selectType(configType)
	if err != nil {
		return nil, err
	}
	var bytes []byte
	if err = c.Bind(&confTypeStruct); err != nil {
		return nil, err
	}
	bytes, err = json.Marshal(confTypeStruct)
	if err != nil {
		return nil, err
	}
	result, err := client.CreateConfig(context.Background(), &api.Config{ConfigType: configType, Config: bytes})
	if err != nil {
		return nil, err
	}
	return result, err
}
