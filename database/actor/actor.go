package actor

import "github.com/jinzhu/gorm"

//Actor ...
type Actor struct {
	gorm.Model
	Identifier string
}
