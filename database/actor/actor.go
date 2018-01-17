package actor

import "github.com/jinzhu/gorm"

//Actor ORM from the message DTO
type Actor struct {
	gorm.Model
	ActorID   string
	IPAddress string
	GoOS      string
	Kernel    string
	Core      string
	Platform  string
	OS        string
	Hostname  string
	CPUs      int32
}
