package database

import (
	"fmt"
	"time"

	"log"

	"errors"
	"strings"

	"reflect"

	"github.com/YAWAL/GetMeConf/dataStructs"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

func InitPostgresDB(cfg PostgresConfig) (err error) {
	if db != nil {
		log.Printf("connection to postgres database already exists")
		return nil
	}
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Dbhost, cfg.Dbport, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	db, err = gorm.Open("postgres", dbInfo)

	if err != nil {
		log.Printf("error during connection to postgres database has occurred: %v", err)
		return err
	}
	db.DB().SetMaxOpenConns(cfg.MaxOpenedConnectionsToDb)
	db.DB().SetMaxIdleConns(cfg.MaxIdleConnectionsToDb)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(cfg.MbConnMaxLifetimeMinutes))
	log.Printf("connection to postgres database has been established")
	return nil
}

//
//func (db *dbConnection) GetAllRuleDefinition() ([]RuleDefinition, error) {
//	ruleDefinitions := make([]RuleDefinition, 0)
//	var ruleDefinition RuleDefinition
//	rows, err := db.Table("RULE_DEFINITION").Rows()
//	if err != nil {
//		return nil, err
//	}
//	for rows.Next() {
//		err = rows.Scan(&ruleDefinition)
//		if err != nil {
//			return nil, err
//		}
//		ruleDefinitions = append(ruleDefinitions, ruleDefinition)
//	}
//	return ruleDefinitions, nil
//}
//
//func (db *dbConnection) GetRuleDefinitionByID(ruleDefinitionID string) (*RuleDefinition, error) {
//	var ruleDefinition RuleDefinition
//	err := db.Where("RULE_PCD = ?", ruleDefinitionID).First(&ruleDefinition).Error
//	if err != nil {
//		return nil, err
//	}
//	return &ruleDefinition, nil
//}

func GetAll(db *gorm.DB) ([]dataStructs.Mongodb, error) {

	var results []dataStructs.Mongodb
	err := db.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

type configInterface interface {
}

func GetConfigByNameFromDB(confName string, confType string) (configInterface, error) {
	type persistedData struct {
		configType configInterface
		idField    string
	}
	var factory = map[string]persistedData{
		"mongodb":    persistedData{new(dataStructs.Mongodb), "domain"},
		"tempconfig": persistedData{new(dataStructs.TempConfig), "host"},
		"tsconfig":   persistedData{new(dataStructs.Tsconfig), "module"},
	}

	cType := strings.ToLower(confType)
	configStruct, ok := factory[cType]
	if !ok {
		return nil, errors.New("unexpected config type")
	}
	if !db.HasTable(configStruct.configType) {
		return nil, errors.New("could not find table " + cType)
	}

	result := configStruct.configType
	err := db.Where(configStruct.idField+" = ?", confName).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetConfigsByTypeFromDB(confType string) ([]configInterface, error) {
	type persistedData struct {
		configType configInterface
		idField    string
	}
	var factory = map[string]persistedData{
		"mongodb":    persistedData{new(dataStructs.Mongodb), "domain"},
		"tempconfig": persistedData{new(dataStructs.TempConfig), "host"},
		"tsconfig":   persistedData{new(dataStructs.Tsconfig), "module"},
	}

	cType := strings.ToLower(confType)
	configStruct, ok := factory[cType]
	if !ok {
		return nil, errors.New("unexpected config type")
	}
	if !db.HasTable(configStruct.configType) {
		return nil, errors.New("could not find table " + cType)
	}

	structType := reflect.TypeOf(configStruct.configType)
	value := reflect.MakeSlice(reflect.SliceOf(structType), 0, 0)
	structSlice := reflect.New(value.Type())
	structSlice.Elem().Set(value)

	rows, err := db.Find(structSlice.Interface()).Rows()
	rows.Scan()
	//err := db.Find(structSlice.Interface()).Error

	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(structSlice)
	return structSlice.Interface().([]configInterface), nil
}
