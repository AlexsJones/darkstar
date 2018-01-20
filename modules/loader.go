package modules

import modules "github.com/AlexsJones/darkstar/modules/scavange"

//LoadModule ...
func LoadModule(name string) (IModule, error) {

	switch name {
	case "scavange":
		return &modules.Scavange{}, nil
	}

	panic("Unable to load module from darkstar server")
}
