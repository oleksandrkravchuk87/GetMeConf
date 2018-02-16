package dataStructs

type Mongodb struct {
	Domain  string `json:"domain"`
	Mongodb bool   `json:"mongodb"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

type Tsconfig struct {
	Module    string `json:"module"`
	Target    string `json:"target"`
	SourseMap bool   `json:"sourseMap"`
	Exclude   int    `json:"exclude"`
}

type Tempconfig struct {
	RestApiRoot    string `json:"restApiRoot"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	Remoting       string `json:"remoting"`
	LegasyExplorer bool   `json:"legasyExplorer"`
}

type ConfigInterface interface {
}

//PersistedData stores the information about all config types in database and is used during searching for a config by name and type
type PersistedData struct {
	ConfigType ConfigInterface
	IDField    string
}
