package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/YAWAL/GetMeConf/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
)

const address = "localhost:8081"

var (
	client     api.ConfigServiceClient
	configName *string
	configType *string
	outPath    *string
)

func main() {

	configName := flag.String("config-name", "", "config name")
	configType := flag.String("config-type", "", "config type")
	outPath := flag.String("outpath", "", "output path for config file")
	//serverPort := flag.String("server", "localhost:50111", "port for connection to server")
	flag.Parse()

	if *configName == "" && *configType == ""{
		log.Fatal("Can't proccess => config name and config type are empty")
	}

	log.Printf("Start checking input data:\n Config name: %v\n Config typy : %v\n Output path: %v", *configName, * configType, * outPath)
	log.Printf("Processing ...")

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatalf("DialContext error has occurred: %v", err)
	}

	client = api.NewConfigServiceClient(conn)

	if *configName != "" && *configType != "" {
		retrieveConfig(*configName, *outPath)
	}

	if *configName == "" && *configType != "" {
		retrieveConfigs()

	}

	for true {
	}
}

func retrieveConfig(fileName, outputPath string) {
	conf, err := client.GetConfigByName(context.Background(), &api.GetConfigByNameRequest{ConfigName: *configName, ConfigType: *configType})
	if err != nil {
		log.Fatalf("Error during retrieving config has occurred: %v", err)
	}
	if err := WriteFile(conf.Config, fileName, outputPath); err != nil {
		log.Fatalf("Error during writing file in retrieving config: %v", err)
	}
}

func retrieveConfigs() {
	stream, err := client.GetConfigsByType(context.Background(), &api.GetConfigsByTypeRequest{ConfigType: *configType})
	if err != nil {
		log.Fatalf("Error during retrieving configs has occurred:%v", err)
	}
	for  {
		config, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error during streaming has occurred: %v", err)
		}

		//procesing configs
	}




}


func WriteFile(data []byte, fileName, outPath string) error {
	if err := ioutil.WriteFile(filepath.Join(outPath, fileName, ".json"), data, 0666); err != nil {
		log.Fatalf("Error during file creation: %v", err)
		return err
	} else {
		log.Printf("File %v has been created in %v", fileName, outPath)
		return nil
	}
}
