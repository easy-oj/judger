package settings

import (
	"io/ioutil"
	"runtime"

	"gopkg.in/yaml.v2"
)

type TSettings struct {
	MySQL  TMySQL  `yaml:"mysql" json:"mysql"`
	Redis  TRedis  `yaml:"redis" json:"redis"`
	Core   TCore   `yaml:"core" json:"core"`
	Repos  TRepos  `yaml:"repos" json:"repos"`
	OSS    TOSS    `yaml:"oss" json:"oss"`
	Queue  TQueue  `yaml:"queue" json:"queue"`
	Judger TJudger `yaml:"judger" json:"judger"`
	Path   TPath   `yaml:"path" json:"path"`
}

type TMySQL struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"-"`
	Password string `yaml:"password" json:"-"`
	Database string `yaml:"database" json:"database"`
}

type TRedis struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
}

type TCore struct {
	HTTP TCoreHTTP `yaml:"http" json:"http"`
}

type TCoreHTTP struct {
	Port int `yaml:"port" json:"port"`
}

type TRepos struct {
	Hosts []string   `yaml:"hosts" json:"hosts"`
	Port  int        `yaml:"port" json:"port"`
	HTTP  TReposHTTP `yaml:"http" json:"http"`
	Path  string     `yaml:"path" json:"path"`
}

type TReposHTTP struct {
	Port int  `yaml:"port" json:"port"`
	Log  bool `yaml:"log" json:"log"`
}

type TOSS struct {
	Hosts   []string `yaml:"hosts" json:"hosts"`
	Port    int      `yaml:"port" json:"port"`
	Backend string   `yaml:"backend" json:"backend"`
}

type TQueue struct {
	Hosts []string `yaml:"hosts" json:"hosts"`
	Port  int      `yaml:"port" json:"port"`
}

type TJudger struct {
	Worker int    `yaml:"worker" json:"worker"`
	Path   string `yaml:"path" json:"path"`
	Docker string `yaml:"docker" json:"docker"`
}

type TPath struct {
	Git      string `yaml:"git" json:"git"`
	Tar      string `yaml:"tar" json:"tar"`
	Executor string `yaml:"executor" json:"executor"`
}

var (
	Settings TSettings
	MySQL    TMySQL
	Redis    TRedis
	Core     TCore
	Repos    TRepos
	OSS      TOSS
	Queue    TQueue
	Judger   TJudger
	Path     TPath
)

func InitSettings(confPath string) {
	Settings = newDefaultSettings()

	if bs, err := ioutil.ReadFile(confPath); err != nil {
		panic(err)
	} else if err = yaml.Unmarshal(bs, &Settings); err != nil {
		panic(err)
	}

	MySQL = Settings.MySQL
	Redis = Settings.Redis
	Core = Settings.Core
	Repos = Settings.Repos
	OSS = Settings.OSS
	Queue = Settings.Queue
	Judger = Settings.Judger
	Path = Settings.Path
}

func newDefaultSettings() TSettings {
	return TSettings{
		MySQL: TMySQL{
			Host:     "127.0.0.1",
			Port:     3306,
			Database: "eoj",
		},
		Redis: TRedis{
			Host: "127.0.0.1",
			Port: 6379,
		},
		Core: TCore{
			HTTP: TCoreHTTP{
				Port: 8080,
			},
		},
		Repos: TRepos{
			Port: 3000,
			HTTP: TReposHTTP{
				Port: 8090,
				Log:  false,
			},
			Path: "/dev/shm/eoj/repos",
		},
		OSS: TOSS{
			Port: 3100,
		},
		Queue: TQueue{
			Port: 3200,
		},
		Judger: TJudger{
			Worker: runtime.NumCPU(),
			Path:   "/dev/shm/eoj/judger",
			Docker: "unix:///var/run/docker.sock",
		},
		Path: TPath{
			Git:      "/usr/bin/git",
			Tar:      "/bin/tar",
			Executor: "/opt/eoj/bin/executor",
		},
	}
}
