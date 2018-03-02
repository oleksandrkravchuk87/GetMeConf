package main

import (
	"log"
	"google.golang.org/grpc"
	"github.com/YAWAL/GetMeConf/api"
	"golang.org/x/net/context"
)

func deleteConfig(configType, configName string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	log.Printf("State: %v", conn.GetState())
	if err != nil {
		log.Fatalf("Dial error has occurred: %v", err)
	}
	defer conn.Close()
	client := api.NewConfigServiceClient(conn)
	resp, err := client.DeleteConfig(context.Background(), &api.DeleteConfigRequest{ConfigType: configType, ConfigName: configName})
	if err != nil {
		log.Printf("Error during client.DeleteConfig has occurred: %v", err)
	}
	if resp.Status == "" {
		log.Printf("Responce status is empty: %v", resp.Status)
	}
}
