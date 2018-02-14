package database

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
	Excluding int    `json:"excluding"`
}

type TempConfig struct {
	RestApiRoot    string `json:"restApiRoot"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	Remoting       string `json:"remoting"`
	LegasyExplorer bool   `json:"legasyExplorer"`
}
