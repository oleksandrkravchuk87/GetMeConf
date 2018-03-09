package repository

import (
	"log"

	"fmt"
	"regexp"

	"testing"

	"github.com/YAWAL/GetMeConf/entities"

	"github.com/jinzhu/gorm"

	"errors"

	"strconv"

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

func TestFind(t *testing.T) {
	m, db, _ := newDB()
	mongoRepo := MongoDBConfigRepoImpl{DB: db}
	mongodbConfig := entities.Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "testPort"}
	mongoRows := getMongoDBRows(mongodbConfig.Domain)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).WithArgs("testDomain").WillReturnRows(mongoRows)
	returnedMongoConfigs, err := mongoRepo.Find("testDomain")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, &mongodbConfig, returnedMongoConfigs)

	configName := "notExistingConfig"
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).WithArgs("notExistingConfig").WillReturnError(expectedError)
	_, returnedErr := mongoRepo.Find(configName)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

	tsRepo := TsConfigRepoImpl{DB: db}
	tsConfig := entities.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	tsRows := getTsConfigRows(tsConfig.Module)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).WithArgs("testModule").WillReturnRows(tsRows)
	returnedTsConfigs, err := tsRepo.Find("testModule")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	fmt.Println(returnedTsConfigs)
	assert.Equal(t, &tsConfig, returnedTsConfigs)

	configName = "notExistingConfig"
	expectedError = errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).WithArgs("notExistingConfig").WillReturnError(expectedError)
	_, returnedErr = tsRepo.Find(configName)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

	tempRepo := TempConfigRepoImpl{DB: db}
	tempConfig := entities.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	tempRows := getTempConfigRows(tempConfig.RestApiRoot)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).WithArgs("testRestApiRoot").WillReturnRows(tempRows)
	returnedTempConfigs, err := tempRepo.Find("testRestApiRoot")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	fmt.Println(returnedTempConfigs)
	assert.Equal(t, &tempConfig, returnedTempConfigs)

	configName = "notExistingConfig"
	expectedError = errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).WithArgs("notExistingConfig").WillReturnError(expectedError)
	_, returnedErr = tempRepo.Find(configName)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

}

func TestFindAll(t *testing.T) {
	m, db, _ := newDB()
	mongoRepo := MongoDBConfigRepoImpl{DB: db}
	mongodbConfig := entities.Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "testPort"}
	mongoRows := getMongoDBRows(mongodbConfig.Domain)
	expConfigs := []entities.Mongodb{mongodbConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\"")).WillReturnRows(mongoRows)
	returnedMongoConfigs, err := mongoRepo.FindAll()
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expConfigs, returnedMongoConfigs)

	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\"")).WillReturnError(expectedError)
	_, returnedErr := mongoRepo.FindAll()
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

	tsRepo := TsConfigRepoImpl{DB: db}
	tsConfig := entities.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	tsRows := getTsConfigRows(tsConfig.Module)
	expTsConfigs := []entities.Tsconfig{tsConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\"")).WillReturnRows(tsRows)
	returnedTsConfigs, err := tsRepo.FindAll()
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expTsConfigs, returnedTsConfigs)

	expectedError = errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\"")).WillReturnError(expectedError)
	_, returnedErr = tsRepo.FindAll()
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

	tempRepo := TempConfigRepoImpl{DB: db}
	tempConfig := entities.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	tempRows := getTempConfigRows(tempConfig.RestApiRoot)
	expTempConfigs := []entities.Tempconfig{tempConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\"")).WillReturnRows(tempRows)
	returnedTempConfigs, err := tempRepo.FindAll()
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	fmt.Println(returnedTempConfigs)
	assert.Equal(t, expTempConfigs, returnedTempConfigs)

	expectedError = errors.New("db error")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\"")).WillReturnError(expectedError)
	_, returnedErr = tempRepo.FindAll()
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

}

func TestSave(t *testing.T) {
	m, db, _ := newDB()
	mockRepo := MongoDBConfigRepoImpl{DB: db}
	mongodbConfig := entities.Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "testPort"}
	m.ExpectExec(formatRequest("INSERT INTO \"mongodbs\" (\"domain\",\"mongodb\",\"host\",\"port\") VALUES ($1,$2,$3,$4) RETURNING \"mongodbs\".*")).
		WithArgs("testDomain", true, "testHost", "testPort").
		WillReturnResult(sqlmock.NewResult(0, 1))
	result, err := mockRepo.Save(&mongodbConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)

	mongodbConfigErr := entities.Mongodb{Domain: "testDomainError", Mongodb: true, Host: "testHost", Port: "testPort"}
	expectedError := errors.New("db error")
	m.ExpectExec(formatRequest("INSERT INTO \"mongodbs\" (\"domain\",\"mongodb\",\"host\",\"port\") VALUES ($1,$2,$3,$4) RETURNING \"mongodbs\".*")).
		WithArgs("testDomainError", true, "testHost", "testPort").
		WillReturnError(expectedError)
	_, returnedErr := mockRepo.Save(&mongodbConfigErr)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

	tsRepo := TsConfigRepoImpl{DB: db}
	tsConfig := entities.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	m.ExpectExec(formatRequest("INSERT INTO \"tsconfigs\" (\"module\",\"target\",\"source_map\",\"excluding\") VALUES ($1,$2,$3,$4) RETURNING \"tsconfigs\".*")).
		WithArgs("testModule", "testTarget", true, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	result, err = tsRepo.Save(&tsConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)

	tsConfigErr := entities.Tsconfig{Module: "testModuleError", Target: "testTarget", SourceMap: true, Excluding: 1}
	expectedError = errors.New("db error")
	m.ExpectExec(formatRequest("INSERT INTO \"tsconfigs\" (\"module\",\"target\",\"source_map\",\"excluding\") VALUES ($1,$2,$3,$4) RETURNING \"tsconfigs\".*")).
		WithArgs("testModuleError", "testTarget", true, 1).
		WillReturnError(expectedError)
	_, returnedErr = tsRepo.Save(&tsConfigErr)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

	tempRepo := TempConfigRepoImpl{DB: db}
	tempConfig := entities.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	m.ExpectExec(formatRequest("INSERT INTO \"tempconfigs\" (\"rest_api_root\",\"host\",\"port\",\"remoting\",\"legasy_explorer\") VALUES ($1,$2,$3,$4,$5) RETURNING \"tempconfigs\".*")).
		WithArgs("testApiRoot", "testHost", "testPort", "testRemoting", true).
		WillReturnResult(sqlmock.NewResult(0, 1))
	result, err = tempRepo.Save(&tempConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)

	tempConfigErr := entities.Tempconfig{RestApiRoot: "testApiRootError", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	expectedError = errors.New("db error")
	m.ExpectExec(formatRequest("INSERT INTO \"tempconfigs\" (\"rest_api_root\",\"host\",\"port\",\"remoting\",\"legasy_explorer\") VALUES ($1,$2,$3,$4,$5) RETURNING \"tempconfigs\".*")).
		WithArgs("testApiRootError", "testHost", "testPort", "testRemoting", true).
		WillReturnError(expectedError)
	_, returnedErr = tempRepo.Save(&tempConfigErr)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestDelete(t *testing.T) {
	m, db, _ := newDB()
	mockRepo := MongoDBConfigRepoImpl{DB: db}
	testType := "mongodb"
	testID := "testID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (domain = $1)")).
		WithArgs("testID").WillReturnResult(sqlmock.NewResult(0, 1))
	res, err := mockRepo.Delete(testID)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "deleted 1 row(s)", res)

	testID = "notExistingTestID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (domain = $1)")).
		WithArgs("notExistingTestID").WillReturnError(errors.New("could not delete from database"))
	_, returnedErr := mockRepo.Delete(testID)
	expectedError := errors.New("could not delete from database")
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

	tsRepo := TsConfigRepoImpl{DB: db}
	testType = "tsconfig"
	testID = "testID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (module = $1)")).
		WithArgs("testID").WillReturnResult(sqlmock.NewResult(0, 1))
	res, err = tsRepo.Delete(testID)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "deleted 1 row(s)", res)

	testID = "notExistingTestID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (module = $1)")).
		WithArgs("notExistingTestID").WillReturnError(errors.New("could not delete from database"))
	_, returnedErr = tsRepo.Delete(testID)
	expectedError = errors.New("could not delete from database")
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

	tempRepo := TempConfigRepoImpl{DB: db}
	testType = "tempconfig"
	testID = "testID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (rest_api_root = $1)")).
		WithArgs("testID").WillReturnResult(sqlmock.NewResult(0, 1))
	res, err = tempRepo.Delete(testID)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "deleted 1 row(s)", res)

	testID = "notExistingTestID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (rest_api_root = $1)")).
		WithArgs("notExistingTestID").WillReturnError(errors.New("could not delete from database"))
	_, returnedErr = tempRepo.Delete(testID)
	expectedError = errors.New("could not delete from database")
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestUpdate(t *testing.T) {
	m, db, _ := newDB()
	mockRepo := MongoDBConfigRepoImpl{DB: db}
	config := entities.Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "8080"}
	rows := getMongoDBRows(config.Domain)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("testDomain").
		WillReturnRows(rows)
	m.ExpectExec(formatRequest("UPDATE mongodbs SET mongodb = $1, port = $2, host = $3 WHERE domain = $4")).
		WithArgs(strconv.FormatBool(config.Mongodb), config.Port, config.Host, config.Domain).
		WillReturnResult(sqlmock.NewResult(0, 1))
	result, err := mockRepo.Update(&config)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)

	configErrOne := entities.Mongodb{Domain: "errOneConfig", Mongodb: true, Host: "testHost", Port: "8080"}
	rows = getMongoDBRows(configErrOne.Domain)
	expectedErrorOne := errors.New("record not found")
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("errOneConfig").
		WillReturnError(expectedErrorOne)
	_, returnedErr := mockRepo.Update(&configErrOne)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedErrorOne, returnedErr)
	}

	expectedErrorTwo := errors.New("db error")
	configErrTwo := entities.Mongodb{Domain: "errTwoConfig", Mongodb: true, Host: "testHost", Port: "8080"}
	rows = getMongoDBRows(configErrTwo.Domain)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("errTwoConfig").
		WillReturnRows(rows)
	m.ExpectExec(formatRequest("UPDATE mongodbs SET mongodb = $1, port = $2, host = $3 WHERE domain = $4")).
		WithArgs(strconv.FormatBool(configErrTwo.Mongodb), configErrTwo.Port, configErrTwo.Host, configErrTwo.Domain).
		WillReturnError(expectedErrorTwo)
	_, returnedErr = mockRepo.Update(&configErrTwo)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedErrorTwo, returnedErr)
	}

	expectedErrorThree := errors.New("fields are empty")
	configErrThree := entities.Mongodb{Domain: "errThreeConfig", Mongodb: true, Host: "", Port: ""}
	rows = getMongoDBRows(configErrThree.Domain)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("errThreeConfig").
		WillReturnRows(rows)
	m.ExpectExec(formatRequest("UPDATE mongodbs SET mongodb = $1, port = $2, host = $3 WHERE domain = $4")).
		WithArgs(strconv.FormatBool(configErrThree.Mongodb), configErrThree.Port, configErrThree.Host, configErrThree.Domain).
		WillReturnError(expectedErrorThree)
	_, returnedErr = mockRepo.Update(&configErrThree)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedErrorThree, returnedErr)
	}

	m, db, _ = newDB()

	tsRepo := TsConfigRepoImpl{DB: db}
	tsConfig := entities.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	tsRows := getTsConfigRows(tsConfig.Module)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("testModule").
		WillReturnRows(tsRows)
	m.ExpectExec(formatRequest("UPDATE tsconfigs SET target = $1, source_map = $2, excluding = $3 WHERE module = $4")).
		WithArgs(tsConfig.Target, strconv.FormatBool(tsConfig.SourceMap), strconv.Itoa(tsConfig.Excluding), tsConfig.Module).
		WillReturnResult(sqlmock.NewResult(0, 1))
	tsResult, err := tsRepo.Update(&tsConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", tsResult)

	tsConfigErrOne := entities.Tsconfig{Module: "errOneConfig", Target: "testTarget", SourceMap: true, Excluding: 1}
	expectedTsErrorOne := errors.New("record not found")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("errOneConfig").
		WillReturnError(expectedTsErrorOne)
	_, tsReturnedErr := tsRepo.Update(&tsConfigErrOne)
	if assert.Error(t, tsReturnedErr) {
		assert.Equal(t, expectedTsErrorOne, tsReturnedErr)
	}

	expectedTsErrorTwo := errors.New("db error")
	tsConfigErrTwo := entities.Tsconfig{Module: "errTwoConfig", Target: "testTarget", SourceMap: true, Excluding: 1}
	tsRows = getTsConfigRows(tsConfigErrTwo.Module)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("errTwoConfig").
		WillReturnRows(tsRows)
	m.ExpectExec(formatRequest("UPDATE tsconfigs SET target = $1, source_map = $2, excluding = $3 WHERE module = $4")).
		WithArgs(tsConfigErrTwo.Target, strconv.FormatBool(tsConfigErrTwo.SourceMap), strconv.Itoa(tsConfigErrTwo.Excluding), tsConfigErrTwo.Module).
		WillReturnError(expectedTsErrorTwo)
	_, tsReturnedErrTwo := tsRepo.Update(&tsConfigErrTwo)
	if assert.Error(t, tsReturnedErrTwo) {
		assert.Equal(t, expectedTsErrorTwo, tsReturnedErrTwo)
	}

	expectedTsErrorThree := errors.New("fields are empty")
	tsConfigErrThree := entities.Tsconfig{Module: "errThreeConfig", Target: "", SourceMap: true, Excluding: 1}
	tsRows = getTsConfigRows(tsConfigErrThree.Module)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("errThreeConfig").
		WillReturnRows(tsRows)
	m.ExpectExec(formatRequest("UPDATE tsconfigs SET target = $1, source_map = $2, excluding = $3 WHERE module = $4")).
		WithArgs(tsConfigErrThree.Target, strconv.FormatBool(tsConfigErrThree.SourceMap), strconv.Itoa(tsConfigErrThree.Excluding), tsConfigErrThree.Module).
		WillReturnError(expectedTsErrorThree)
	_, tsReturnedErrThree := tsRepo.Update(&tsConfigErrThree)
	if assert.Error(t, tsReturnedErrThree) {
		assert.Equal(t, expectedTsErrorThree, tsReturnedErrThree)
	}

	m, db, _ = newDB()

	tempRepo := TempConfigRepoImpl{DB: db}
	tempConfig := entities.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	tempRows := getTempConfigRows(tempConfig.RestApiRoot)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("testApiRoot").
		WillReturnRows(tempRows)
	m.ExpectExec(formatRequest("UPDATE tempconfigs SET remoting = $1, port = $2, host = $3, legasy_explorer = $4 WHERE rest_api_root = $5")).
		WithArgs(tempConfig.Remoting, tempConfig.Port, tempConfig.Host, strconv.FormatBool(tempConfig.LegasyExplorer), tempConfig.RestApiRoot).
		WillReturnResult(sqlmock.NewResult(0, 1))
	tempResult, err := tempRepo.Update(&tempConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", tempResult)

	tempConfigErrOne := entities.Tempconfig{RestApiRoot: "errOneConfig", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	expectedTempErrorOne := errors.New("record not found")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("errOneConfig").
		WillReturnError(expectedTempErrorOne)
	_, tempReturnedErr := tempRepo.Update(&tempConfigErrOne)
	if assert.Error(t, tempReturnedErr) {
		assert.Equal(t, expectedTempErrorOne, tempReturnedErr)
	}

	expectedTempErrorTwo := errors.New("db error")
	tempConfigErrTwo := entities.Tempconfig{RestApiRoot: "errTwoConfig", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	tempRows = getTempConfigRows(tempConfigErrTwo.RestApiRoot)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("errTwoConfig").
		WillReturnRows(tempRows)
	m.ExpectExec(formatRequest("UPDATE tempconfigs SET remoting = $1, port = $2, host = $3, legasy_explorer = $4 WHERE rest_api_root = $5")).
		WithArgs(tempConfigErrTwo.Remoting, tempConfigErrTwo.Port, tempConfigErrTwo.Host, strconv.FormatBool(tempConfigErrTwo.LegasyExplorer), tempConfigErrTwo.RestApiRoot).
		WillReturnError(expectedTempErrorTwo)
	_, tempReturnedErrTwo := tempRepo.Update(&tempConfigErrTwo)
	if assert.Error(t, tempReturnedErrTwo) {
		assert.Equal(t, expectedTempErrorTwo, tempReturnedErrTwo)
	}

	expectedTempErrorThree := errors.New("fields are empty")
	tempConfigErrThree := entities.Tempconfig{RestApiRoot: "errThreeConfig", Host: "", Port: "", Remoting: "", LegasyExplorer: true}
	tempRows = getTempConfigRows(tempConfigErrThree.RestApiRoot)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("errThreeConfig").
		WillReturnRows(tempRows)
	m.ExpectExec(formatRequest("UPDATE tempconfigs SET remoting = $1, port = $2, host = $3, legasy_explorer = $4 WHERE rest_api_root = $5")).
		WithArgs(tempConfigErrThree.Remoting, tempConfigErrThree.Port, tempConfigErrThree.Host, strconv.FormatBool(tempConfigErrThree.LegasyExplorer), tempConfigErrThree.RestApiRoot).
		WillReturnError(expectedTempErrorThree)
	_, tempReturnedErrThree := tempRepo.Update(&tempConfigErrThree)
	if assert.Error(t, tempReturnedErrThree) {
		assert.Equal(t, expectedTempErrorThree, tempReturnedErrThree)
	}
}

func getMongoDBRows(configID string) *sqlmock.Rows {
	var fieldNames = []string{"domain", "mongodb", "host", "port"}
	rows := sqlmock.NewRows(fieldNames)
	mongodbConfig := entities.Mongodb{Domain: configID, Mongodb: true, Host: "testHost", Port: "testPort"}
	rows = rows.AddRow(mongodbConfig.Domain, mongodbConfig.Mongodb, mongodbConfig.Host, mongodbConfig.Port)
	return rows
}

func getTsConfigRows(configID string) *sqlmock.Rows {
	var fieldNames = []string{"module", "target", "source_map", "excluding"}
	rows := sqlmock.NewRows(fieldNames)
	tsConfig := entities.Tsconfig{Module: configID, Target: "testTarget", SourceMap: true, Excluding: 1}
	rows = rows.AddRow(tsConfig.Module, tsConfig.Target, tsConfig.SourceMap, tsConfig.Excluding)
	return rows
}

func getTempConfigRows(configID string) *sqlmock.Rows {
	var fieldNames = []string{"rest_api_root", "host", "port", "remoting", "legasy_explorer"}
	rows := sqlmock.NewRows(fieldNames)
	tempConfig := entities.Tempconfig{RestApiRoot: configID, Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	rows = rows.AddRow(tempConfig.RestApiRoot, tempConfig.Host, tempConfig.Port, tempConfig.Remoting, tempConfig.LegasyExplorer)
	return rows
}
