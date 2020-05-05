package main

type Config struct {
	ServerPort string         `yaml:"service_port"`
	Database   DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Task struct {
	Id               int64
	Url              string `json:"url"`
	NumberOfRequests int    `json:"number_of_requests"`
}

type Tasks struct {
	Total    int64 `json:"total"`
	InQueue  int64 `json:"in_queue"`
	Complete int64 `json:"complete"`
	Error    int64 `json:"error"`
}

type Urls struct {
	Total    int64 `json:"total"`
	InQueue  int64 `json:"in_queue"`
	Complete int64 `json:"complete"`
	Error    int64 `json:"error"`
}

type Status struct {
	Tasks Tasks `json:"tasks"`
	Urls  Urls  `json:"urls"`
}
