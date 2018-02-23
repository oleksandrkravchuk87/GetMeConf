package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	pb "github.com/YAWAL/GetMeConf/api"
	"github.com/YAWAL/GetMeConf/database"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestGetConfigByName(t *testing.T) {
	old := databaseGetConfigByNameFromDB
	defer func() { databaseGetConfigByNameFromDB = old }()

	databaseGetConfigByNameFromDB = func(confName string, confType string, db *gorm.DB) (database.ConfigInterface, error) {
		return database.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}, nil
	}

	tServer := new(configServer)
	cache := cache.New(5*time.Minute, 10*time.Minute)
	tServer = &configServer{configCache: cache, mut: &sync.Mutex{}}

	res, err := tServer.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "mongodb", ConfigName: "testName"})

	if err != nil {
		fmt.Println(err)
	}

	var expectedConfig []byte
	expectedConfig, err = json.Marshal(database.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"})
	if err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, expectedConfig, res.Config)

}

func TestGetConfigsByType(t *testing.T) {
	old := databaseGetConfigByNameFromDB
	defer func() { databaseGetConfigByNameFromDB = old }()

	databaseGetConfigByNameFromDB = func(confName string, confType string, db *gorm.DB) (database.ConfigInterface, error) {
		return database.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}, nil
	}

	tServer := new(configServer)
	cache := cache.New(5*time.Minute, 10*time.Minute)
	tServer = &configServer{configCache: cache, mut: &sync.Mutex{}}

	res, err := tServer.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "mongodb", ConfigName: "testName"})

	if err != nil {
		fmt.Println(err)
	}

	var expectedConfig []byte
	expectedConfig, err = json.Marshal(database.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"})
	if err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, expectedConfig, res.Config)

}

type mockConfigServer struct {
	grpc.ServerStream
	Results []*pb.GetConfigResponce
}

func (_m *mockConfigServer) Send(response *pb.GetConfigResponce) error {
	_m.Results = append(_m.Results, response)
	return nil
}

func TestMarshalAndSend(t *testing.T) {
	mock := &mockConfigServer{}
	testConfig := database.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}
	_ = marshalAndSend(testConfig, mock)
	assert.Equal(t, 1, len(mock.Results), "expected to contain 1 item")
}

func (s *configServer) GetConfigsByType(typeRequest *pb.GetConfigsByTypeRequest, stream pb.ConfigService_GetConfigsByTypeServer) error {

	switch typeRequest.ConfigType {
	case "mongodb":
		res, err := database.GetMongoDBConfigs(s.db)
		if err != nil {
			return err
		}
		for _, v := range res {
			if err = marshalAndSend(v, stream); err != nil {
				return err
			}
		}
	case "tempconfig":
		res, err := database.GetTempConfigs(s.db)
		if err != nil {
			return err
		}
		for _, v := range res {
			if err = marshalAndSend(v, stream); err != nil {
				return err
			}
		}
	case "tsconfig":
		res, err := database.GetTsconfigs(s.db)
		if err != nil {
			return err
		}
		for _, v := range res {
			if err = marshalAndSend(v, stream); err != nil {
				return err
			}
		}
	default:
		log.Print("unexpacted type")
		return errors.New("unexpacted type")
	}
	return nil
}
