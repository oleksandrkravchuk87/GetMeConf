package main

import (
	"log"
	"google.golang.org/grpc"
	"github.com/YAWAL/GetMeConf/api"
	"golang.org/x/net/context"
)

func sentUpdatedConfigToServer(fileName string) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	log.Printf("State: %v", conn.GetState())
	if err != nil {
		log.Fatalf("Dial error has occurred: %v", err)
	}
	defer conn.Close()
	client := api.NewConfigServiceClient(conn)
	config := createByteConfig(fileName)
	resp, err := client.UpdateConfig(context.Background(), &api.Config{Config: config, ConfigType: *configType})
	if err != nil {
		log.Printf("Error during client.CreateConfig has occurred: %v", err)
	}
	if resp.Status != "OK" {
		log.Printf("Error during creating config has occurred: %v responce status: %v", err, resp.Status)
	}
	log.Printf("Responce: %v", resp.Status)
}