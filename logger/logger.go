package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"naga/config"
	"os"
	"runtime"
	"time"
)

var OrmLog = logrus.New()

func Start() {
	// 定义日志格式
	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	logrus.SetOutput(os.Stdout)
	logrus.AddHook(&DefaultFieldsHook{})
	logSite := NewHook()
	logSite.Field = "line"
	logrus.AddHook(logSite)

	OrmLog.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	OrmLog.SetOutput(os.Stdout)
	OrmLog.AddHook(&DefaultFieldsHook{})
	OrmLog.AddHook(logSite)

	if config.Config.LogPath != "" {
		// check log path exist
		_ = os.MkdirAll(config.Config.LogPath, os.ModePerm)
		// spite log every day
		newLfsHook()
		newLfsHook2()
		// log request params to log in warning
	}
}

func newLfsHook() {
	var writer *rotatelogs.RotateLogs
	var err error
	if runtime.GOOS == "windows" {
		writer, err = rotatelogs.New(
			config.Config.LogPath+"framework"+".%Y%m%d.log",
			// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
			rotatelogs.WithRotationTime(24*time.Hour),
			// WithMaxAge和WithRotationCount二者只能设置一个,
			// WithMaxAge设置文件清理前的最长保存时间,
			// WithRotationCount设置文件清理前最多保存的个数.
			rotatelogs.WithMaxAge(24*time.Hour*7),
			//rotatelogs.WithRotationCount(31),
		)
	} else {
		writer, err = rotatelogs.New(
			config.Config.LogPath+"framework"+".%Y%m%d.log",
			// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
			rotatelogs.WithLinkName(config.Config.LogPath+"framework"),
			// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
			rotatelogs.WithRotationTime(24*time.Hour),
			// WithMaxAge和WithRotationCount二者只能设置一个,
			// WithMaxAge设置文件清理前的最长保存时间,
			// WithRotationCount设置文件清理前最多保存的个数.
			rotatelogs.WithMaxAge(24*time.Hour*7),
			//rotatelogs.WithRotationCount(31),
		)
	}

	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}
	logrus.SetLevel(logrus.InfoLevel)
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.AddHook(lfsHook)
}

func newLfsHook2() {
	writer, err := rotatelogs.New(
		config.Config.LogPath+"orm"+".%Y%m%d.log",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(config.Config.LogPath+"orm"),
		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotatelogs.WithRotationTime(24*time.Hour),
		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		rotatelogs.WithMaxAge(24*time.Hour*7),
		//rotatelogs.WithRotationCount(31),
	)

	if err != nil {
		OrmLog.Errorf("config local file system for logger error: %v", err)
	}
	OrmLog.SetLevel(logrus.InfoLevel)
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	OrmLog.AddHook(lfsHook)
}

type DefaultFieldsHook struct{}

func (df *DefaultFieldsHook) Fire(entry *logrus.Entry) error {
	_, ok := entry.Data["module"].(string)
	if ok == false {
		entry.Data["module"] = "system"
	}
	return nil
}

func (df *DefaultFieldsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
