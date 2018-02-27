package main

import (
	"encoding/csv"
	"os"
	"log"
	"github.com/YAWAL/GetMeConf/database"
	"strconv"
	"google.golang.org/grpc"
	"github.com/YAWAL/GetMeConf/api"
	"golang.org/x/net/context"
	"encoding/json"
)

func readConfig(fileName string) ([][]string) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		log.Fatalf("error during opening file has occurred: %v", err)
		return nil
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("error during reading file has occurred: %v", err)
		return nil
	}
	return records
}

func createByteConfig(fileName string) []byte {
	records := readConfig(fileName)
	switch fileName {
	case "mongo.csv":
		var mongocnf database.Mongodb
		mongocnf.Domain = records[0][0]

		if records[0][1] != "true" && records[0][1] != "false" {
			log.Fatalf("field Mongodb should be true or false, but is: %v", records[0][1])
			return nil
		}
		if records[0][1] == "true" {
			mongocnf.Mongodb = true
		}
		if records[0][1] == "false" {
			mongocnf.Mongodb = false
		}
		mongocnf.Host = records[0][2]
		mongocnf.Port = records[0][3]
		if bytesMongo, err := json.Marshal(mongocnf); err == nil {
			return bytesMongo
		} else {
			log.Printf("Error during converting Mongodb structure to []byte has occurred: %v", err)
			return nil
		}

	case "tempcnf.csv":
		var tempcnf database.Tempconfig
		tempcnf.RestApiRoot = records[0][0]
		tempcnf.Host = records[0][1]
		tempcnf.Port = records[0][2]
		tempcnf.Remoting = records[0][3]
		if records[0][4] != "true" && records[0][4] != "false" {
			log.Fatalf("field legasyExplorer should be true or false, but is: %v", records[0][4])
			return nil
		}
		if records[0][4] == "true" {
			tempcnf.LegasyExplorer = true
		}
		if records[0][4] == "false" {
			tempcnf.LegasyExplorer = false
		}
		if bytesTempcnf, err := json.Marshal(tempcnf); err == nil {
			return bytesTempcnf
		} else {
			log.Printf("Error during converting Tempconfig structure to []byte has occurred: %v", err)
			return nil
		}
	case "tscnf.csv":
		var tscnf database.Tsconfig
		tscnf.Module = records[0][0]
		tscnf.Target = records[0][1]
		if records[0][2] != "true" && records[0][2] != "false" {
			log.Fatalf("field sourceMap should be true or false, but is: %v", records[0][2])
			return nil
		}
		if records[0][2] == "true" {
			tscnf.SourceMap = true
		}
		if records[0][2] == "false" {
			tscnf.SourceMap = false
		}
		excluding, err := strconv.Atoi(records[0][3])
		if err != nil {
			log.Fatalf("field Excluding should be integer, but is: %t", records[0][3])
		}
		tscnf.Excluding = excluding
		if bytesTscnf, err := json.Marshal(tscnf); err == nil {
			return bytesTscnf
		} else {
			log.Printf("Error during converting Tsconfig structure to []byte has occurred: %v", err)
			return nil
		}
	default:
		log.Fatalf("Cant find file: %v", fileName)
	}
	return nil
}

func sentConfigToServer(fileName string) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	log.Printf("State: %v", conn.GetState())
	if err != nil {
		log.Fatalf("Dial error has occurred: %v", err)
	}
	client := api.NewConfigServiceClient(conn)
	config := createByteConfig(fileName)
	resp, err := client.CreateConfig(context.Background(), &api.Config{Config: config, ConfigType: *configType})
	if err != nil {
		log.Printf("Error during client.CreateConfig has occurred: %v", err)
	}
	if resp.Status != "OK" {
		log.Printf("Error during creating config has occurred: %v responce status: %v", err, resp.Status)
	}
}
