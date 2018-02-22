package database

import (
	"log"

	"fmt"
	"regexp"

	"testing"

	"errors"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func newDB() (sqlmock.Sqlmock, *gorm.DB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("can not create sql mock %v", err)
		return nil, nil, err
	}
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		log.Fatalf("can not open gorm connection %v", err)
		return nil, nil, err
	}
	gormDB.LogMode(true)
	return mock, gormDB, nil
}

func formatRequest(s string) string {
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(s))
}

func TestGetMongoDBConfigs(t *testing.T) {
	m, db, _ := newDB()
	var fieldNames = []string{"domain", "mongodb", "host", "port"}
	rows := sqlmock.NewRows(fieldNames)
	mongodbConfig := Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "testPort"}
	rows = rows.AddRow(mongodbConfig.Domain, mongodbConfig.Mongodb, mongodbConfig.Host, mongodbConfig.Port)
	expConfig := []Mongodb{mongodbConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\"")).WillReturnRows(rows)
	returnedMongoConfigs, err := GetMongoDBConfigs(db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expConfig, returnedMongoConfigs)
}

func TestGetTsconfigs(t *testing.T) {
	m, db, _ := newDB()
	var fieldNames = []string{"module", "target", "source_map", "excluding"}
	rows := sqlmock.NewRows(fieldNames)
	tsConfig := Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	rows = rows.AddRow(tsConfig.Module, tsConfig.Target, tsConfig.SourceMap, tsConfig.Excluding)
	expConfig := []Tsconfig{tsConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\"")).WillReturnRows(rows)
	returnedTsConfigs, err := GetTsconfigs(db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expConfig, returnedTsConfigs)
}

func TestGetTempConfigs(t *testing.T) {
	m, db, _ := newDB()
	var fieldNames = []string{"rest_api_root", "host", "port", "remoting", "legasy_explorer"}
	rows := sqlmock.NewRows(fieldNames)
	tempConfig := Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	rows = rows.AddRow(tempConfig.RestApiRoot, tempConfig.Host, tempConfig.Port, tempConfig.Remoting, tempConfig.LegasyExplorer)
	expConfig := []Tempconfig{tempConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\"")).WillReturnRows(rows)
	returnedTempConfigs, err := GetTempConfigs(db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expConfig, returnedTempConfigs)
}

func TestGetConfigByNameFromDB(t *testing.T) {
	m, db, _ := newDB()
	initConfigDataMap()
	testType := "mongodb"
	testName := "testDomain"
	anotherTestType := "someType"
	var fieldNames = []string{"domain", "mongodb", "host", "port"}
	rows := sqlmock.NewRows(fieldNames)
	expConfig := Mongodb{Domain: testName, Mongodb: true, Host: "testHost", Port: "testPort"}
	rows = rows.AddRow(expConfig.Domain, expConfig.Mongodb, expConfig.Host, expConfig.Port)
	m.ExpectQuery(formatRequest("SELECT * FROM \"" + testType + "s\" WHERE (domain = $1)")).WillReturnRows(rows)
	returnedConfig, err := GetConfigByNameFromDB(testName, testType, db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, &expConfig, returnedConfig)
	_, err = GetConfigByNameFromDB(testName, anotherTestType, db)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpected config type"), err)
	}
}
