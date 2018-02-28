package database

import (
	"log"

	"fmt"
	"regexp"

	"testing"

	"errors"

	"github.com/gin-gonic/gin/json"
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
	mongodbConfig := Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "testPort"}
	rows := getMongoDBRows()
	expConfig := []Mongodb{mongodbConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\"")).WillReturnRows(rows)
	returnedMongoConfigs, err := GetMongoDBConfigs(db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expConfig, returnedMongoConfigs)
}

func TestGetMongoDBConfigs_withDBError(t *testing.T) {
	m, db, _ := newDB()
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\"")).WillReturnError(expectedError)
	_, returnedErr := GetMongoDBConfigs(db)

	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

}

func TestGetTsconfigs(t *testing.T) {
	m, db, _ := newDB()
	rows := getTsConfigRows()
	tsConfig := Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	expConfig := []Tsconfig{tsConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\"")).WillReturnRows(rows)
	returnedTsConfigs, err := GetTsconfigs(db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expConfig, returnedTsConfigs)
}

func TestGetTsconfigs_withDBError(t *testing.T) {
	m, db, _ := newDB()
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\"")).WillReturnError(expectedError)
	_, returnedErr := GetTsconfigs(db)

	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestGetTempConfigs(t *testing.T) {
	m, db, _ := newDB()
	rows := getTempConfigRows()
	tempConfig := Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	expConfig := []Tempconfig{tempConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\"")).WillReturnRows(rows)
	returnedTempConfigs, err := GetTempConfigs(db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expConfig, returnedTempConfigs)
}

func TestGetTempConfigs_withDBError(t *testing.T) {
	m, db, _ := newDB()
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\"")).WillReturnError(expectedError)
	_, returnedErr := GetTempConfigs(db)

	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
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

func TestGetConfigByNameFromDB_withDBError(t *testing.T) {
	m, db, _ := newDB()
	initConfigDataMap()
	testType := "mongodb"
	testName := "testDomain"
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"" + testType + "s\" WHERE (domain = $1)")).WillReturnError(expectedError)
	_, returnedErr := GetConfigByNameFromDB(testName, testType, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestSaveConfigToDB(t *testing.T) {
	m, db, _ := newDB()
	initConfigDataMap()
	testType := "mongodb"
	config := Mongodb{"testDomain", true, "testHost", "8080"}
	configBytes, _ := json.Marshal(config)
	m.ExpectExec(formatRequest("INSERT INTO \"mongodbs\" (\"domain\",\"mongodb\",\"host\",\"port\") VALUES ($1,$2,$3,$4) RETURNING \"mongodbs\".*")).
		WithArgs("testDomain", true, "testHost", "8080").
		WillReturnResult(sqlmock.NewResult(0, 1))
	result, _ := SaveConfigToDB(testType, configBytes, db)
	assert.Equal(t, "OK", result)

	config = Mongodb{"notExisitingConfig", true, "testHost", "8080"}
	configBytes, _ = json.Marshal(config)
	m.ExpectExec(formatRequest("INSERT INTO \"mongodbs\" (\"domain\",\"mongodb\",\"host\",\"port\") VALUES ($1,$2,$3,$4) RETURNING \"mongodbs\".*")).
		WithArgs("notExisitingConfig", true, "testHost", "8080").
		WillReturnError(errors.New("db error"))
	_, returnedErr := SaveConfigToDB(testType, configBytes, db)
	expectedError := errors.New("db error")
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestDeleteConfigFromDB(t *testing.T) {
	m, db, _ := newDB()
	initConfigDataMap()
	testType := "mongodb"
	testID := "testID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (domain = $1)")).
		WithArgs("testID").WillReturnResult(sqlmock.NewResult(0, 1))
	res, err := DeleteConfigFromDB(testID, testType, db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "deleted 1 row(s)", res)

	testID = "notExistingTestID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (domain = $1)")).
		WithArgs("notExistingTestID").WillReturnError(errors.New("could not delete from database"))
	_, returnedErr := DeleteConfigFromDB(testID, testType, db)
	expectedError := errors.New("could not delete from database")
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func getMongoDBRows() *sqlmock.Rows {
	var fieldNames = []string{"domain", "mongodb", "host", "port"}
	rows := sqlmock.NewRows(fieldNames)
	mongodbConfig := Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "testPort"}
	rows = rows.AddRow(mongodbConfig.Domain, mongodbConfig.Mongodb, mongodbConfig.Host, mongodbConfig.Port)
	return rows
}

func getTsConfigRows() *sqlmock.Rows {
	var fieldNames = []string{"module", "target", "source_map", "excluding"}
	rows := sqlmock.NewRows(fieldNames)
	tsConfig := Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	rows = rows.AddRow(tsConfig.Module, tsConfig.Target, tsConfig.SourceMap, tsConfig.Excluding)
	return rows
}

func getTempConfigRows() *sqlmock.Rows {
	var fieldNames = []string{"rest_api_root", "host", "port", "remoting", "legasy_explorer"}
	rows := sqlmock.NewRows(fieldNames)
	tempConfig := Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	rows = rows.AddRow(tempConfig.RestApiRoot, tempConfig.Host, tempConfig.Port, tempConfig.Remoting, tempConfig.LegasyExplorer)
	return rows
}

func TestUpdateMongoDBConfigInDB(t *testing.T) {
	m, db, _ := newDB()
	rows := getMongoDBRows()
	initConfigDataMap()
	config := Mongodb{"testDomain", true, "testHost", "8080"}
	configBytes, _ := json.Marshal(config)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("testDomain").
		WillReturnRows(rows)
	m.ExpectExec("^UPDATE \"mongodbs\" SET ").
		WillReturnResult(sqlmock.NewResult(0, 1))
	result, err := UpdateMongoDBConfigInDB(configBytes, db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)
}

func TestUpdateMongoDBConfigInDB_withFirstDBError(t *testing.T) {
	m, db, _ := newDB()
	initConfigDataMap()
	config := Mongodb{"notExisitingConfig", true, "testHost", "8080"}
	configBytes, _ := json.Marshal(config)
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("notExisitingConfig").
		WillReturnError(expectedError)
	_, returnedErr := UpdateMongoDBConfigInDB(configBytes, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestUpdateMongoDBConfigInDB_withSecondDBError(t *testing.T) {
	m, db, _ := newDB()
	rows := getMongoDBRows()
	initConfigDataMap()
	config := Mongodb{"testDomain", true, "testHost", "8080"}
	configBytes, _ := json.Marshal(config)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("testDomain").
		WillReturnRows(rows)
	expectedError := errors.New("db error")
	m.ExpectExec("UPDATE \"mongodbs\" SET").
		WillReturnError(expectedError)
	_, returnedErr := UpdateMongoDBConfigInDB(configBytes, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestUpdateMongoDBConfigInDB_withThirdDBError(t *testing.T) {
	m, db, _ := newDB()
	rows := getMongoDBRows()
	initConfigDataMap()
	config := Mongodb{"testDomain", true, "", ""}
	configBytes, _ := json.Marshal(config)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("testDomain").
		WillReturnRows(rows)
	expectedError := errors.New("fields are empty")
	m.ExpectExec("UPDATE \"mongodbs\" SET").
		WillReturnError(expectedError)
	_, returnedErr := UpdateMongoDBConfigInDB(configBytes, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestUpdateTempConfigInDB(t *testing.T) {
	m, db, _ := newDB()
	rows := getTempConfigRows()
	initConfigDataMap()
	config := Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	configBytes, _ := json.Marshal(config)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("testApiRoot").
		WillReturnRows(rows)
	m.ExpectExec("UPDATE \"tempconfigs\" SET").
		WillReturnResult(sqlmock.NewResult(0, 1))
	result, err := UpdateTempConfigInDB(configBytes, db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)

}

func TestUpdateTempConfigInDB_withFirstDBError(t *testing.T) {
	m, db, _ := newDB()
	initConfigDataMap()
	config := Tempconfig{RestApiRoot: "notExisitingConfig", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	configBytes, _ := json.Marshal(config)
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("notExisitingConfig").
		WillReturnError(expectedError)
	_, returnedErr := UpdateTempConfigInDB(configBytes, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestUpdateTempConfigInDB_withSecondDBError(t *testing.T) {
	m, db, _ := newDB()
	rows := getTempConfigRows()
	initConfigDataMap()
	config := Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	configBytes, _ := json.Marshal(config)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("testApiRoot").
		WillReturnRows(rows)
	expectedError := errors.New("db error")
	m.ExpectExec("UPDATE \"tempconfigs\" SET").
		WillReturnError(expectedError)
	_, returnedErr := UpdateTempConfigInDB(configBytes, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestUpdateTempConfigInDB_withThirdDBError(t *testing.T) {
	m, db, _ := newDB()
	rows := getTempConfigRows()
	initConfigDataMap()
	config := Tempconfig{RestApiRoot: "testApiRoot", Host: "", Port: "", Remoting: "", LegasyExplorer: true}
	configBytes, _ := json.Marshal(config)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("testApiRoot").
		WillReturnRows(rows)
	expectedError := errors.New("fields are empty")
	m.ExpectExec("UPDATE \"tsconfigs\" SET").
		WillReturnError(expectedError)
	_, returnedErr := UpdateTempConfigInDB(configBytes, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestUpdateTsConfigInDB(t *testing.T) {
	m, db, _ := newDB()
	rows := getTsConfigRows()
	initConfigDataMap()
	config := Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	configBytes, _ := json.Marshal(config)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("testModule").
		WillReturnRows(rows)
	m.ExpectExec("UPDATE \"tsconfigs\" SET").
		WillReturnResult(sqlmock.NewResult(0, 1))
	result, err := UpdateTsConfigInDB(configBytes, db)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)
}

func TestUpdateTsConfigInDB_withFirstDBError(t *testing.T) {
	m, db, _ := newDB()
	initConfigDataMap()
	config := Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	configBytes, _ := json.Marshal(config)
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("testModule").
		WillReturnError(expectedError)
	_, returnedErr := UpdateTsConfigInDB(configBytes, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestUpdateTsConfigInDB_withSecondDBError(t *testing.T) {
	m, db, _ := newDB()
	rows := getTsConfigRows()
	initConfigDataMap()
	config := Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	configBytes, _ := json.Marshal(config)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("testModule").
		WillReturnRows(rows)
	expectedError := errors.New("db error")
	m.ExpectExec("UPDATE \"tsconfigs\" SET ").
		WillReturnError(expectedError)
	_, returnedErr := UpdateTsConfigInDB(configBytes, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestUpdateTsConfigInDB_withThirdDBError(t *testing.T) {
	m, db, _ := newDB()
	rows := getMongoDBRows()
	initConfigDataMap()
	config := Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	configBytes, _ := json.Marshal(config)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("testModule").
		WillReturnRows(rows)
	expectedError := errors.New("fields are empty")
	m.ExpectExec("UPDATE \"tsconfigs\" SET ").
		WillReturnError(expectedError)
	_, returnedErr := UpdateTsConfigInDB(configBytes, db)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}
