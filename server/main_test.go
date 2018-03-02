package main

import (
	"context"
	"encoding/json"
	"errors"
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
	configCache := cache.New(5*time.Minute, 10*time.Minute)
	mock := &mockConfigServer{}
	mock.configCache = configCache
	res, err := mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "mongodb", ConfigName: "testName"})

	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	var expectedConfig []byte
	expectedConfig, err = json.Marshal(database.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)
	expectedError := errors.New("error from database querying")
	databaseGetConfigByNameFromDB = func(confName string, confType string, db *gorm.DB) (database.ConfigInterface, error) {
		return nil, expectedError
	}
	_, err = mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "someType", ConfigName: "SomeName"})
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}

}

func TestGetConfigByName_FromCache(t *testing.T) {
	testName := "testName"
	testConf := database.Mongodb{Domain: testName, Mongodb: true, Host: "testHost", Port: "testPort"}
	configCache := cache.New(5*time.Minute, 10*time.Minute)
	mock := &mockConfigServer{}
	mock.configCache = configCache

	byteRes, err := json.Marshal(testConf)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	configResponse := &pb.GetConfigResponce{Config: byteRes}
	mock.configCache.Set(testName, configResponse, 5*time.Minute)
	res, err := mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "mongodb", ConfigName: "testName"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	var expectedConfig []byte
	expectedConfig, err = json.Marshal(database.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)

}

func TestGetConfigsByType(t *testing.T) {
	oldDatabaseGetMongoDBConfigs := databaseGetMongoDBConfigs
	defer func() { databaseGetMongoDBConfigs = oldDatabaseGetMongoDBConfigs }()

	oldDatabaseGetTempConfigs := databaseGetTempConfigs
	defer func() { databaseGetTempConfigs = oldDatabaseGetTempConfigs }()

	oldDatabaseGetTsconfigs := databaseGetTsconfigs
	defer func() { databaseGetTsconfigs = oldDatabaseGetTsconfigs }()

	databaseGetMongoDBConfigs = func(db *gorm.DB) ([]database.Mongodb, error) {
		return []database.Mongodb{{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}}, nil
	}

	databaseGetTempConfigs = func(db *gorm.DB) ([]database.Tempconfig, error) {
		return []database.Tempconfig{{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}}, nil
	}

	databaseGetTsconfigs = func(db *gorm.DB) ([]database.Tsconfig, error) {
		return []database.Tsconfig{{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}}, nil
	}

	mock := &mockConfigServer{}

	err := mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "mongodb"}, mock)
	assert.Equal(t, 1, len(mock.Results), "expected to contain 1 item")
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "tsconfig"}, mock)
	assert.Equal(t, 2, len(mock.Results), "expected to contain 1 item")
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "tempconfig"}, mock)
	assert.Equal(t, 3, len(mock.Results), "expected to contain 1 item")
	if err != nil {
		t.Error("error during unit testing of GetConfigsByType function: ", err)
	}
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "unexpectedConfigType"}, mock)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpacted type"), err)
	}

	expectedError := errors.New("error from database querying")
	databaseGetMongoDBConfigs = func(db *gorm.DB) ([]database.Mongodb, error) {
		return nil, expectedError
	}

	databaseGetTempConfigs = func(db *gorm.DB) ([]database.Tempconfig, error) {
		return nil, expectedError
	}

	databaseGetTsconfigs = func(db *gorm.DB) ([]database.Tsconfig, error) {
		return nil, expectedError
	}
	err = nil
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "mongodb"}, mock)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("error from database querying"), err)
	}

	err = nil
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "tsconfig"}, mock)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("error from database querying"), err)
	}
	err = nil
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "tempconfig"}, mock)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}

}

type mockConfigServer struct {
	configServer
	grpc.ServerStream
	Results []*pb.GetConfigResponce
}

func (mcs *mockConfigServer) Send(response *pb.GetConfigResponce) error {
	mcs.Results = append(mcs.Results, response)
	return nil
}

func TestMarshalAndSend(t *testing.T) {
	mock := &mockConfigServer{}
	testConfig := database.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}
	_ = marshalAndSend(testConfig, mock)
	assert.Equal(t, 1, len(mock.Results), "expected to contain 1 item")
}

func TestCreateConfig(t *testing.T) {
	old := databaseSaveConfigToDB
	defer func() { databaseSaveConfigToDB = old }()
	databaseSaveConfigToDB = func(confType string, config []byte, db *gorm.DB) (string, error) {
		return "OK", nil

	}
	configCache := cache.New(5*time.Minute, 10*time.Minute)
	mock := &mockConfigServer{}
	mock.configCache = configCache

	res, err := mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "testType", Config: []byte("testConfig")})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedResponse := &pb.Responce{Status: "OK"}
	assert.Equal(t, expectedResponse, res)

	expectedError := errors.New("error from database querying")
	databaseSaveConfigToDB = func(confType string, config []byte, db *gorm.DB) (string, error) {
		return "", expectedError

	}
	_, resultingErr := mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "testType", Config: []byte("errTestConfig")})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
}

func TestDeleteConfig(t *testing.T) {
	old := databaseDeleteConfigFromDB
	defer func() { databaseDeleteConfigFromDB = old }()
	databaseDeleteConfigFromDB = func(configName, confType string, db *gorm.DB) (string, error) {
		return "OK", nil
	}
	configCache := cache.New(5*time.Minute, 10*time.Minute)
	mock := &mockConfigServer{}
	mock.configCache = configCache

	res, err := mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "testType", ConfigName: "testName"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedResponse := &pb.Responce{Status: "OK"}
	assert.Equal(t, expectedResponse, res)

	expectedError := errors.New("error from database querying")
	databaseDeleteConfigFromDB = func(configName, confType string, db *gorm.DB) (string, error) {
		return "", expectedError

	}
	_, resultingErr := mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "testType", ConfigName: "errorTestName"})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
}

func TestUpdateConfig(t *testing.T) {
	oldDatabaseUpdateMongoDBConfigInDB := databaseUpdateMongoDBConfigInDB
	defer func() { databaseUpdateMongoDBConfigInDB = oldDatabaseUpdateMongoDBConfigInDB }()

	oldDatabaseUpdateTempconfigInDB := databaseUpdateTempconfigInDB
	defer func() { databaseUpdateTempconfigInDB = oldDatabaseUpdateTempconfigInDB }()

	oldDatabaseUpdateTsConfigInDB := databaseUpdateTsConfigInDB
	defer func() { databaseUpdateTsConfigInDB = oldDatabaseUpdateTsConfigInDB }()

	databaseUpdateMongoDBConfigInDB = func(config []byte, db *gorm.DB) (string, error) {
		return "OK", nil
	}

	databaseUpdateTempconfigInDB = func(config []byte, db *gorm.DB) (string, error) {
		return "OK", nil
	}

	databaseUpdateTsConfigInDB = func(config []byte, db *gorm.DB) (string, error) {
		return "OK", nil
	}

	configCache := cache.New(5*time.Minute, 10*time.Minute)
	mock := &mockConfigServer{}
	mock.configCache = configCache

	resp, err := mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "mongodb"})
	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
	resp, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tsconfig"})
	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
	resp, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tempconfig"})
	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "unexpectedConfigType"})
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpacted type"), err)
	}

	expectedError := errors.New("error from database querying")
	databaseUpdateMongoDBConfigInDB = func(config []byte, db *gorm.DB) (string, error) {
		return "", expectedError
	}

	databaseUpdateTempconfigInDB = func(config []byte, db *gorm.DB) (string, error) {
		return "", expectedError
	}

	databaseUpdateTsConfigInDB = func(config []byte, db *gorm.DB) (string, error) {
		return "", expectedError
	}
	err = nil
	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "mongodb"})
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("error from database querying"), err)
	}

	err = nil
	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tsconfig"})
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("error from database querying"), err)
	}

	err = nil
	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tempconfig"})
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("error from database querying"), err)
	}
}
