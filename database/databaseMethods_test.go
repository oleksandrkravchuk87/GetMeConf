package database

import (
	"log"

	"fmt"
	"regexp"

	"testing"

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
	expMongodb := []Mongodb{mongodbConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\"")).WillReturnRows(rows)
	returnedMongoConfigs, _ := GetMongoDBConfigs(db)
	assert.Equal(t, returnedMongoConfigs, expMongodb)
}

func TestGetTsconfigs(t *testing.T) {
	m, db, _ := newDB()
	var fieldNames = []string{"module", "target", "source_map", "excluding"}
	rows := sqlmock.NewRows(fieldNames)
	tsConfig := Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	rows = rows.AddRow(tsConfig.Module, tsConfig.Target, tsConfig.SourceMap, tsConfig.Excluding)
	expMongodb := []Tsconfig{tsConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\"")).WillReturnRows(rows)
	returnedTsConfigs, _ := GetTsconfigs(db)
	assert.Equal(t, returnedTsConfigs, expMongodb)
}

func TestGetTempConfigs(t *testing.T) {
	m, db, _ := newDB()
	var fieldNames = []string{"rest_api_root", "host", "port", "remoting", "legasy_explorer"}
	rows := sqlmock.NewRows(fieldNames)
	tempConfig := Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	rows = rows.AddRow(tempConfig.RestApiRoot, tempConfig.Host, tempConfig.Port, tempConfig.Remoting, tempConfig.LegasyExplorer)
	expMongodb := []Tempconfig{tempConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\"")).WillReturnRows(rows)
	returnedTempConfigs, _ := GetTempConfigs(db)
	assert.Equal(t, returnedTempConfigs, expMongodb)
}
