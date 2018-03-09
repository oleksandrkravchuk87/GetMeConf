package entities

//Mongodb is an random config example
type Mongodb struct {
	Domain  string `json:"domain"`
	Mongodb bool   `json:"mongodb"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

//Tsconfig is an random config example
type Tsconfig struct {
	Module    string `json:"module"`
	Target    string `json:"target"`
	SourceMap bool   `json:"sourceMap"`
	Excluding int    `json:"excluding"`
}

//Tempconfig is an random config example
type Tempconfig struct {
	RestApiRoot    string `json:"restApiRoot"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	Remoting       string `json:"remoting"`
	LegasyExplorer bool   `json:"legasyExplorer"`
}

//ConfigInterface is an interface for all config structures
type ConfigInterface interface {
}

//PersistedData stores the information about all config types in database and is used during searching for a config by name and type
type PersistedData struct {
	ConfigType ConfigInterface
	IDField    string
}
