package database

import (
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/jinzhu/gorm"
)

//AutoMigrate creates tables
func AutoMigrate(databaseConnection *gorm.DB) {
	databaseConnection.AutoMigrate(&message.Message{})
	databaseConnection.AutoMigrate(&Module{})
}
