package main

import (
	"context"
	"encoding/json"

	"testing"
	"time"

	pb "github.com/YAWAL/GetMeConf/api"

	"errors"

	"github.com/YAWAL/GetMeConf/entitys"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockMongoDBConfigRepo struct {
}

func (m *mockMongoDBConfigRepo) Find(configName string) (*entitys.Mongodb, error) {
	return &entitys.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}, nil
}

func (m *mockMongoDBConfigRepo) FindAll() ([]entitys.Mongodb, error) {
	return []entitys.Mongodb{{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}}, nil
}

func (m *mockMongoDBConfigRepo) Update(config *entitys.Mongodb) (string, error) {
	return "OK", nil
}

func (m *mockMongoDBConfigRepo) Save(config *entitys.Mongodb) (string, error) {
	return "OK", nil
}

func (m *mockMongoDBConfigRepo) Delete(configName string) (string, error) {
	return "OK", nil
}

type mockErrorMongoDBConfigRepo struct {
}

func (m *mockErrorMongoDBConfigRepo) Find(configName string) (*entitys.Mongodb, error) {
	return nil, errors.New("error from database querying")
}

func (m *mockErrorMongoDBConfigRepo) FindAll() ([]entitys.Mongodb, error) {
	return nil, errors.New("error from database querying")
}

func (m *mockErrorMongoDBConfigRepo) Update(config *entitys.Mongodb) (string, error) {
	return "", errors.New("error from database querying")
}

func (m *mockErrorMongoDBConfigRepo) Save(config *entitys.Mongodb) (string, error) {
	return "", errors.New("error from database querying")
}
func (m *mockErrorMongoDBConfigRepo) Delete(configName string) (string, error) {
	return "", errors.New("error from database querying")
}

type mockTsConfigRepo struct {
}

func (m *mockTsConfigRepo) Find(configName string) (*entitys.Tsconfig, error) {
	return &entitys.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}, nil
}

func (m *mockTsConfigRepo) FindAll() ([]entitys.Tsconfig, error) {
	return []entitys.Tsconfig{{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}}, nil
}

func (m *mockTsConfigRepo) Update(config *entitys.Tsconfig) (string, error) {
	return "OK", nil
}

func (m *mockTsConfigRepo) Save(config *entitys.Tsconfig) (string, error) {
	return "OK", nil
}

func (m *mockTsConfigRepo) Delete(configName string) (string, error) {
	return "OK", nil
}

type mockErrorTsConfigRepo struct {
}

func (m *mockErrorTsConfigRepo) Find(configName string) (*entitys.Tsconfig, error) {
	return nil, errors.New("error from database querying")
}

func (m *mockErrorTsConfigRepo) FindAll() ([]entitys.Tsconfig, error) {
	return nil, errors.New("error from database querying")
}

func (m *mockErrorTsConfigRepo) Update(config *entitys.Tsconfig) (string, error) {
	return "", errors.New("error from database querying")
}

func (m *mockErrorTsConfigRepo) Save(config *entitys.Tsconfig) (string, error) {
	return "", errors.New("error from database querying")
}
func (m *mockErrorTsConfigRepo) Delete(configName string) (string, error) {
	return "", errors.New("error from database querying")
}

type mockTempConfigRepo struct {
}

func (m *mockTempConfigRepo) Find(configName string) (*entitys.Tempconfig, error) {
	return &entitys.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}, nil
}

func (m *mockTempConfigRepo) FindAll() ([]entitys.Tempconfig, error) {
	return []entitys.Tempconfig{{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}}, nil
}

func (m *mockTempConfigRepo) Update(config *entitys.Tempconfig) (string, error) {
	return "OK", nil
}

func (m *mockTempConfigRepo) Save(config *entitys.Tempconfig) (string, error) {
	return "OK", nil
}

func (m *mockTempConfigRepo) Delete(configName string) (string, error) {
	return "OK", nil
}

type mockErrorTempConfigRepo struct {
}

func (m *mockErrorTempConfigRepo) Find(configName string) (*entitys.Tempconfig, error) {
	return nil, errors.New("error from database querying")
}

func (m *mockErrorTempConfigRepo) FindAll() ([]entitys.Tempconfig, error) {
	return nil, errors.New("error from database querying")
}

func (m *mockErrorTempConfigRepo) Update(config *entitys.Tempconfig) (string, error) {
	return "", errors.New("error from database querying")
}

func (m *mockErrorTempConfigRepo) Save(config *entitys.Tempconfig) (string, error) {
	return "", errors.New("error from database querying")
}
func (m *mockErrorTempConfigRepo) Delete(configName string) (string, error) {
	return "", errors.New("error from database querying")
}

func TestGetConfigByName(t *testing.T) {

	configCache := cache.New(5*time.Minute, 10*time.Minute)
	mock := &mockConfigServer{}
	mock.configCache = configCache
	mock.mongoDBConfigRepo = &mockMongoDBConfigRepo{}
	mock.tsConfigRepo = &mockTsConfigRepo{}
	mock.tempConfigRepo = &mockTempConfigRepo{}

	res, err := mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "mongodb", ConfigName: "testNameMongo"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	var expectedConfig []byte
	expectedConfig, err = json.Marshal(entitys.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)

	res, err = mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "tsconfig", ConfigName: "testNameTs"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedConfig, err = json.Marshal(entitys.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)

	res, err = mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "tempconfig", ConfigName: "testNameTemp"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedConfig, err = json.Marshal(entitys.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)

	mock.configCache.Flush()

	mock.mongoDBConfigRepo = &mockErrorMongoDBConfigRepo{}
	expectedError := errors.New("error from database querying")
	_, err = mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "mongodb", ConfigName: "testNameMongo"})
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	mock.tsConfigRepo = &mockErrorTsConfigRepo{}
	_, err = mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "tsconfig", ConfigName: "testNameTs"})
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
	_, err = mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "tempconfig", ConfigName: "testNameTemp"})
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	_, err = mock.GetConfigByName(context.Background(), &pb.GetConfigByNameRequest{ConfigType: "unexpectedConfigType", ConfigName: "testNameTemp"})
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpected type"), err)
	}
}

func TestGetConfigByName_FromCache(t *testing.T) {
	testName := "testName"
	testConf := entitys.Mongodb{Domain: testName, Mongodb: true, Host: "testHost", Port: "testPort"}
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
	expectedConfig, err = json.Marshal(entitys.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)

}

func TestGetConfigsByType(t *testing.T) {

	mock := &mockConfigServer{}
	mock.mongoDBConfigRepo = &mockMongoDBConfigRepo{}
	err := mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "mongodb"}, mock)
	assert.Equal(t, 1, len(mock.Results), "expected to contain 1 item")
	mock.tsConfigRepo = &mockTsConfigRepo{}
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "tsconfig"}, mock)
	assert.Equal(t, 2, len(mock.Results), "expected to contain 1 item")
	mock.tempConfigRepo = &mockTempConfigRepo{}
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "tempconfig"}, mock)
	assert.Equal(t, 3, len(mock.Results), "expected to contain 1 item")
	if err != nil {
		t.Error("error during unit testing of GetConfigsByType function: ", err)
	}
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "unexpectedConfigType"}, mock)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpected type"), err)
	}

	expectedError := errors.New("error from database querying")
	err = nil
	mock.mongoDBConfigRepo = &mockErrorMongoDBConfigRepo{}
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "mongodb"}, mock)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}

	err = nil
	mock.tsConfigRepo = &mockErrorTsConfigRepo{}
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "tsconfig"}, mock)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("error from database querying"), err)
	}
	err = nil
	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "tempconfig"}, mock)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	err = nil
	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
	err = mock.GetConfigsByType(&pb.GetConfigsByTypeRequest{ConfigType: "unexpectedType"}, mock)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpected type"), err)
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

func TestCreateConfig(t *testing.T) {

	configCache := cache.New(5*time.Minute, 10*time.Minute)
	mock := &mockConfigServer{}
	mock.configCache = configCache
	mock.mongoDBConfigRepo = &mockMongoDBConfigRepo{}
	mock.tsConfigRepo = &mockTsConfigRepo{}
	mock.tempConfigRepo = &mockTempConfigRepo{}

	testConfMongo := entitys.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}
	byteRes, err := json.Marshal(testConfMongo)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	res, err := mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "mongodb", Config: byteRes})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedResponse := &pb.Responce{Status: "OK"}
	assert.Equal(t, expectedResponse, res)

	testConfTs := entitys.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	byteRes, err = json.Marshal(testConfTs)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	res, err = mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "tsconfig", Config: byteRes})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedResponse, res)

	testConfTemp := entitys.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	byteRes, err = json.Marshal(testConfTemp)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	res, err = mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "tempconfig", Config: byteRes})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedResponse, res)

	mock.mongoDBConfigRepo = &mockErrorMongoDBConfigRepo{}
	mock.tsConfigRepo = &mockErrorTsConfigRepo{}
	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
	expectedError := errors.New("error from database querying")

	_, resultingErr := mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "mongodb", Config: byteRes})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
	resultingErr = nil
	_, resultingErr = mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "tsconfig", Config: byteRes})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
	resultingErr = nil
	_, resultingErr = mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "tempconfig", Config: byteRes})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
	resultingErr = nil
	_, resultingErr = mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "unexpectedType", Config: byteRes})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, errors.New("unexpected type"), resultingErr)
	}
}

func TestDeleteConfig(t *testing.T) {

	configCache := cache.New(5*time.Minute, 10*time.Minute)
	mock := &mockConfigServer{}
	mock.configCache = configCache
	mock.mongoDBConfigRepo = &mockMongoDBConfigRepo{}
	mock.tsConfigRepo = &mockTsConfigRepo{}
	mock.tempConfigRepo = &mockTempConfigRepo{}

	res, err := mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "mongodb", ConfigName: "testName"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedResponse := &pb.Responce{Status: "OK"}
	assert.Equal(t, expectedResponse, res)

	res, err = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "tsconfig", ConfigName: "testName"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	assert.Equal(t, expectedResponse, res)

	res, err = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "tempconfig", ConfigName: "testName"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	assert.Equal(t, expectedResponse, res)

	mock.mongoDBConfigRepo = &mockErrorMongoDBConfigRepo{}
	mock.tsConfigRepo = &mockErrorTsConfigRepo{}
	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
	expectedError := errors.New("error from database querying")
	_, resultingErr := mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "mongodb", ConfigName: "errorTestName"})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
	resultingErr = nil
	_, resultingErr = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "tsconfig", ConfigName: "errorTestName"})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
	resultingErr = nil
	_, resultingErr = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "tempconfig", ConfigName: "errorTestName"})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
	resultingErr = nil
	_, resultingErr = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "unexpectedType", ConfigName: "errorTestName"})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, errors.New("unexpected type"), resultingErr)
	}
}

func TestUpdateConfig(t *testing.T) {

	configCache := cache.New(5*time.Minute, 10*time.Minute)
	mock := &mockConfigServer{}
	mock.configCache = configCache
	mock.mongoDBConfigRepo = &mockMongoDBConfigRepo{}
	mock.tsConfigRepo = &mockTsConfigRepo{}
	mock.tempConfigRepo = &mockTempConfigRepo{}

	testConfMongo := entitys.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}
	byteResMongo, err := json.Marshal(testConfMongo)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	testConfTs := entitys.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	byteResTs, err := json.Marshal(testConfTs)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	testConfTemp := entitys.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	byteResTemp, err := json.Marshal(testConfTemp)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	resp, err := mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "mongodb", Config: byteResMongo})
	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
	resp, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tsconfig", Config: byteResTs})
	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
	resp, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tempconfig", Config: byteResTemp})
	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "unexpectedConfigType"})
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpected type"), err)
	}

	expectedError := errors.New("error from database querying")
	mock.mongoDBConfigRepo = &mockErrorMongoDBConfigRepo{}
	err = nil
	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "mongodb", Config: byteResMongo})
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}

	err = nil
	mock.tsConfigRepo = &mockErrorTsConfigRepo{}
	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tsconfig", Config: byteResTs})
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}

	err = nil
	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tempconfig", Config: byteResTemp})
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}
