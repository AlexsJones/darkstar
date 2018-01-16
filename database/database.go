package database

import (
	"github.com/AlexsJones/darkstar/database/actor"
	"github.com/jinzhu/gorm"
)

//AutoMigrate creates tables
func AutoMigrate(databaseConnection *gorm.DB) {
	databaseConnection.AutoMigrate(&actor.Actor{})
}
