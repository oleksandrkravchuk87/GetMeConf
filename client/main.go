package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"google.golang.org/grpc"
	"github.com/YAWAL/GetMeConf/api"
	"golang.org/x/net/context"
	"io/ioutil"
)

const address = "localhost:8081"
const outputPath = "/home/vya/go/src/github.com/YAWAL/GetMeConf/config"

func main() {

	configId := flag.String("config-id", "", "config id")
	configPath := flag.String("config-path", "", "config file path")
	//serverPort := flag.String("server", "localhost:50111", "port for connection to server")
	flag.Parse()

	log.Println("Start checking input data")

	if err := CheckPath(*configPath); err != nil {
		log.Println("Path to config wrong: ", err)
	}

	if err := CheckFile(*configId, *configPath); err != nil {
		log.Println("File does not exist: ", err)
	}

	log.Printf("Start to prepare data about config in: %v with name %s\n", *configPath, *configId)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatalf("DialContext error has occurred: %v", err)
	}

	client := api.NewConfigServiceClient(conn)

	 //config information for receiving to the server
	cnfgInfo := api.ConfigInfo{ConfigId: *configId, ConfigPath: *configPath}

	log.Print(cnfgInfo)

	preparedInfo, err := client.SearchConfig(context.Background(), &cnfgInfo)
	if err != nil {
		log.Fatalf("Error has occured: %v", err)
	}

	log.Print(preparedInfo)

	WriteFile(nil, outputPath, *configId)

}

func CheckPath(path string) error {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		log.Printf("Path: %v exists", path)
		return nil
	} else {
		return err
	}
}

func CheckFile(configId, configPath string) error {
	filePath := filepath.Join(configPath, configId)
	_, err := os.Stat(filePath);
	if  err != nil {
		if os.IsNotExist(err) {
			log.Printf("File %v does not exist", configId)
			return err
		}
	}
	log.Printf("File %v  exist in directory %v", configId, configPath)
	return nil
}

func WriteFile(data []byte, outputPath,fileName string) error{
	if err := ioutil.WriteFile(fileName, data, 0666); err != nil{
		log.Printf("File %v has been created in %v", fileName, outputPath)
		return nil
	}else {
		log.Fatalf("Error during file creation: %v", err)
		return err
	}
}

