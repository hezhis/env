package env

import (
	"encoding/json"
	"os"
	"sync"
)

type RedisConn struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type DBConnections struct {
	Redis map[string]*RedisConn `json:"redis"`
}

func (c *DBConnections) GetRedisConn(connName string) (*RedisConn, bool) {
	conn, ok := c.Redis[connName]
	return conn, ok
}

type OBS struct {
	Bucket  string `json:"bucket"`
	Backup  int    `json:"backup"`
	Expires int    `json:"expires"`
}

type Cloud struct {
	Obs map[string]*OBS
}

type HuaWeiObs struct {
	Endpoint   string
	Ak         string
	Sk         string
	BucketName string
	PathStyle  bool
}

type Env struct {
	DBConnections `json:"db"`
	Cloud         `json:"cloud"`
	HuaWeiObs     `json:"huawei_obs"`
}

var (
	env        *Env
	hasInitEnv bool
	initEnvMu  sync.RWMutex
)

func InitEnv(cfgFilePath string) error {
	if cfgFilePath == "" {
		cfgFilePath = defaultCfgFilePath
	}

	if hasInitEnv {
		return nil
	}

	initEnvMu.Lock()
	defer initEnvMu.Unlock()

	if hasInitEnv {
		return nil
	}

	env = &Env{}

	buf, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, env)
	if err != nil {
		return err
	}

	hasInitEnv = true

	return nil
}

func MustEnv() *Env {
	if !hasInitEnv {
		panic("env not init")
	}

	return env
}
