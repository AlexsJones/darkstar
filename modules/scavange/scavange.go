package modules

import (
	"encoding/json"
	"fmt"
	"time"
)

//Scavange ...
type Scavange struct {
	CurrentTime time.Time
}

//Execute ...
func (s *Scavange) Execute() string {

	s.CurrentTime = time.Now()

	b, err := json.Marshal(s)
	if err == nil {
		return string(b)
	}
	fmt.Println("Error marshalling data into JSON")
	return ""
}
