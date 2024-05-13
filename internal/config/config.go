package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Cfg struct {
	ollamaModel string

	supervisedModel bool
	debugMode       bool

	pluginsPath string
	logName     string

	malvusCfg MalvusCfg
	AppLogger *logrus.Logger
}

type MalvusCfg struct {
	apiKey         string
	apiEndpoint    string
	collectionName string
}

func New() Cfg {

	malvusCfg := MalvusCfg{
		apiEndpoint:    "localhost:19530",
		collectionName: "CGPTMemory",
	}

	cfg := Cfg{
		ollamaModel:     "llama3",
		supervisedModel: false,
		debugMode:       false,
		pluginsPath:     "./internal/plugins",
		malvusCfg:       malvusCfg,
	}

	err := cfg.InitLogger()
	if err != nil {
		fmt.Println("Error initializing logger: ", err)
		os.Exit(1)
	}

	return cfg
}

func (c Cfg) OllamaModel() string {
	return c.ollamaModel
}

func (c Cfg) SupervisedModel() bool {
	return c.supervisedModel
}

func (c Cfg) DebugMode() bool {
	return c.debugMode
}

func (c Cfg) SetDebugMode(debugMode bool) Cfg {
	c.debugMode = debugMode
	return c
}

func (c Cfg) PluginsPath() string {
	return c.pluginsPath
}

func (c Cfg) SetSupervisedModel(supervisedModel bool) Cfg {
	c.supervisedModel = supervisedModel
	return c
}

func (c Cfg) SetOllamaModel(ollamaModel string) Cfg {
	c.ollamaModel = ollamaModel
	return c
}

func (c Cfg) SetPluginsPath(pluginsPath string) Cfg {
	c.pluginsPath = pluginsPath
	return c
}

func (c Cfg) SetMalvusApiKey(apiKey string) Cfg {
	c.malvusCfg.apiKey = apiKey
	return c
}

func (c Cfg) MalvusApiKey() string {
	return c.malvusCfg.apiKey
}

func (c Cfg) MalvusApiEndpoint() string {
	return c.malvusCfg.apiEndpoint
}

func (c Cfg) SetMalvusApiEndpoint(apiEndpoint string) Cfg {
	c.malvusCfg.apiEndpoint = apiEndpoint
	return c
}

func (c Cfg) MalvusCollectionName() string {
	return c.malvusCfg.collectionName
}

func (c Cfg) SetMalvusCollectionName(collectionName string) Cfg {
	c.malvusCfg.collectionName = collectionName
	return c

}

func (c *Cfg) InitLogger() error {
	if c.logName == "" {
		c.logName = "clara.log"
	}
	file, err := os.OpenFile(c.logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file %s for output: %s", c.logName, err)
	}
	c.AppLogger = logrus.New()
	c.AppLogger.Out = file
	return nil
}
