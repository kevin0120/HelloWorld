package utils

import (
	"io"
	"log"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLCache struct {
	db *gorm.DB
}

var glbSQLCache = SQLCache{nil}

func logLevelFromName(lvl string) logger.LogLevel {
	var level logger.LogLevel = logger.Silent
	switch lvl {
	case "INFO", "info":
		level = logger.Info
	case "ERROR", "error":
		level = logger.Error
	case "DEBUG", "debug":
		level = logger.Silent
	}

	return level
}

func OpenNewSQLCache(out io.Writer, prefix string, level string) (cache *SQLCache, err error) {
	lvl := logger.Silent
	if level == "" {
		lvl = logLevelFromName(level)
	}
	newLogger := logger.New(
		log.New(out, prefix, log.Ldate|log.Lmicroseconds), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             1 * time.Second, // 慢 SQL 阈值
			LogLevel:                  lvl,             // 日志级别
			IgnoreRecordNotFoundError: false,           // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,           // 禁用彩色打印
		},
	)

	cacheDNS := "file::memory:?cache=shared" //生产环境用缓存模式

	cache = &glbSQLCache

	if cache.db != nil {
		return
	}

	if CheckRuntimeEnvIsDev() {
		cacheDNS = "./cache.db"
	}
	if eng, e := gorm.Open(sqlite.Open(cacheDNS), &gorm.Config{
		Logger: newLogger,
	}); e == nil {
		cache.db = eng
	} else {
		err = e
		return
	}

	if cache.db == nil {
		err = errors.New("openCache Cache Engine Is Empty")
	}

	return
}

func (cache *SQLCache) CacheUpdateSchema(bean interface{}) (err error) {
	if cache == nil || cache.db == nil {
		err = errors.New("CacheUpdateSchema Cache Is Empty")
		return
	}
	eng := cache.db
	err = eng.AutoMigrate(bean)
	return
}

func (cache *SQLCache) CloseCache() (err error) {
	return nil
}

func (cache *SQLCache) ExecRawSQL(sql string) (err error, rowsAffected int64) {
	if cache == nil || cache.db == nil {
		err = errors.New("ExecRawSQL Cache Is Empty")
		return
	}
	eng := cache.db
	e := eng.Exec(sql)
	err = e.Error
	if err != nil {
		return
	}
	return e.Error, e.RowsAffected
}

func (cache *SQLCache) RawDBObj() *gorm.DB {
	if cache == nil || cache.db == nil {
		return nil
	}
	return cache.db
}
