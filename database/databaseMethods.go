package database

import (
	"fmt"
	"time"

	"log"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

var (
	once sync.Once
)

func InitPostgresDB(cfg PostgresConfig) (db *gorm.DB, err error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Dbhost, cfg.Dbport, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	once.Do(func() {
		db, err = gorm.Open("postgres", dbInfo)
	})
	if err != nil {
		log.Printf("error during connection to postgres database has occurred: %v", err)
		return nil, err
	}
	db.DB().SetMaxOpenConns(cfg.MaxOpenedConnectionsToDb)
	db.DB().SetMaxIdleConns(cfg.MaxIdleConnectionsToDb)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(cfg.MbConnMaxLifetimeMinutes))
	log.Printf("connection to postgres database has been established")
	return db, nil
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

func GetAll(db *gorm.DB) ([]Mongodb, error) {

	var results []Mongodb
	err := db.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
