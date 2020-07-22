package orm

import (
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"naga/config"
	nagalog "naga/logger"
	"time"
)

var (
	Gorm        *gorm.DB
	err         error
	RedisClient *redis.Client
	Cache       *cache.Cache
)

var c = config.Config

type MyLogger struct {
}

func (logger MyLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		nagalog.OrmLog.WithFields(
			logrus.Fields{
				"module":    "gorm",
				"values":    v[4],
				"rows":      v[5],
				"src_ref":   v[1],
				"exec_time": fmt.Sprintf("%gms", float64(v[2].(time.Duration).Nanoseconds()/1e4)/100.0),
			},
		).Info(v[3])
	case "log":
		nagalog.OrmLog.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Info(v[2])
	}
}

func Start() {
	// mysql连接池
	a := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", c.Mysql.User, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Database, c.Mysql.Charset)
	Gorm, err = gorm.Open("mysql", a)
	if err != nil {
		logrus.Panic("connect mysql err:" + err.Error())
		return
	}
	Gorm.LogMode(c.Mysql.Logdebug)
	Gorm.SetLogger(&MyLogger{})
	Gorm.DB().SetMaxOpenConns(c.Mysql.MaxActive) // 用于设置最大打开的连接数，默认值为0表示不限制。
	Gorm.DB().SetMaxIdleConns(c.Mysql.MaxIdle)   // 最大空闲数
	_ = Gorm.DB().Ping()

	// goredis
	RedisClient = redis.NewClient(&redis.Options{
		DB:           c.Redis.Database,
		Password:     c.Redis.Password,
		MinIdleConns: c.Redis.MaxIdle,
		PoolSize:     c.Redis.MaxActive,
		Addr:         c.Redis.Host + ":" + c.Redis.Port,
	})

	// go-cache
	Cache = cache.New(1*time.Second, 5*time.Second)
}
