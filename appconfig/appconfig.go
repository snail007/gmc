package appconfig

import (
	"fmt"
	"os"
	"time"

	"github.com/snail007/gmc/session"
	"github.com/snail007/gmc/session/filestore"
	"github.com/snail007/gmc/session/memorystore"
	"github.com/snail007/gmc/session/redisstore"
	"github.com/snail007/gmc/template"
	"github.com/spf13/viper"
)

var (
	//Parsed config
	SessionStore      session.Store
	SessionCookieName string
	Tpl               *template.Template
	FileConfig        *viper.Viper
	isParsed          bool
	Config            *APPConfig
)

//APPConfig configs all settings of gmc
type APPConfig struct {
	ViewsDir   string
	ConfigFile string
}

func NewAPPConfig() *APPConfig {
	return &APPConfig{
		ViewsDir:   "views",
		ConfigFile: "",
	}
}

func Parse(cfg *APPConfig) (err error) {
	if isParsed {
		return
	}
	defer func() {
		if err == nil {
			isParsed = true
		}
	}()
	Config = cfg

	//create views template object
	Tpl, err = template.New(Config.ViewsDir)
	if err != nil {
		return
	}

	//create config file object
	FileConfig = viper.New()
	if cfg.ConfigFile == "" {
		FileConfig.SetConfigName("app")
		FileConfig.AddConfigPath(".")
		FileConfig.AddConfigPath("config")
		FileConfig.AddConfigPath("conf")
		FileConfig.AddConfigPath("../appconfig")
	} else {
		FileConfig.SetConfigFile(cfg.ConfigFile)
	}
	err = FileConfig.ReadInConfig()
	if err != nil {
		return
	}
	//create session store object
	err = initSessionStore()
	if err != nil {
		return
	}
	return
}

func initSessionStore() (err error) {
	if !FileConfig.GetBool("session.enable") {
		return
	}
	typ := FileConfig.GetString("session.store")
	if typ == "" {
		typ = "memory"
	}
	ttl := FileConfig.GetInt64("session.ttl")
	SessionCookieName = FileConfig.GetString("session.cookie_name")

	switch typ {
	case "file":
		cfg := filestore.NewConfig()
		cfg.TTL = ttl
		dir := FileConfig.GetString("session.file.dir")
		if dir == "" {
			dir = os.TempDir()
		}
		cfg.Dir = dir
		cfg.GCtime = FileConfig.GetInt("session.file.gctime")
		cfg.Prefix = FileConfig.GetString("session.file.prefix")
		SessionStore, err = filestore.New(cfg)
	case "memory":
		cfg := memorystore.NewConfig()
		cfg.TTL = ttl
		cfg.GCtime = FileConfig.GetInt("session.memory.gctime")
		SessionStore, err = memorystore.New(cfg)
	case "redis":
		cfg := redisstore.NewRedisStoreConfig()
		cfg.RedisCfg.Addr = FileConfig.GetString("session.redis.address")
		cfg.RedisCfg.Password = FileConfig.GetString("session.redis.password")
		cfg.RedisCfg.Prefix = FileConfig.GetString("session.redis.prefix")
		cfg.RedisCfg.Debug = FileConfig.GetBool("session.redis.debug")
		cfg.RedisCfg.Timeout = time.Second * FileConfig.GetDuration("session.redis.timeout")
		cfg.RedisCfg.DBNum = FileConfig.GetInt("session.redis.dbnum")
		cfg.RedisCfg.MaxIdle = FileConfig.GetInt("session.redis.maxidle")
		cfg.RedisCfg.MaxActive = FileConfig.GetInt("session.redis.maxactive")
		cfg.RedisCfg.MaxConnLifetime = time.Second * FileConfig.GetDuration("session.redis.maxconnlifetime")
		cfg.RedisCfg.Wait = FileConfig.GetBool("session.redis.wait")
		cfg.TTL = ttl
		SessionStore, err = redisstore.New(cfg)
	default:
		err = fmt.Errorf("unknown session store type %s", typ)
	}
	return
}
