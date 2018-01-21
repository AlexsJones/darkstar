package modules

import (
	"encoding/json"
	"fmt"
	"time"
)

//Scavange ...
type Scavange struct {
}

//Execute ...
func (*Scavange) Execute() string {

	b, err := json.Marshal(struct{ Time time.Time }{time.Now()})
	if err != nil {
		return string(b)
	}
	fmt.Println("Error marshalling data into JSON")
	return ""
}
