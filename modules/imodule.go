package modules

//IModule definition
type IModule interface {
	Execute() string
}

func Execute(m IModule) string {
	return m.Execute()
}
