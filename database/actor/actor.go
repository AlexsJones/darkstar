package actor

import "github.com/jinzhu/gorm"

//Payload ...
type Payload struct {
	Data string
}

//Instruction ...
type Instruction struct {
	ModuleName    string
	ModulePayload *Payload
}

//Actor ORM from the message DTO
type Actor struct {
	gorm.Model
	ActorID            string
	IPAddress          string
	GoOS               string
	Kernel             string
	Core               string
	Platform           string
	OS                 string
	Hostname           string
	CPUs               int32
	CurrentInstruction *Instruction
}
