package db

import (

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gormopentracing "gorm.io/plugin/opentracing"
	"log"
	"os"
	"time"
)

func Init(dsn string) (db *gorm.DB, err error) {

	newLogger := logger.New (
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config {
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // 日志级别Silent、Error、Warn、Info
			IgnoreRecordNotFoundError: true,   // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:      true,         // 禁用彩色打印
		},
	)

	//全局模式
	db, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config {
			Logger: newLogger,
			PrepareStmt:            true,
			SkipDefaultTransaction: false,
			DisableNestedTransaction: true,
			DisableForeignKeyConstraintWhenMigrating: false,//默认开启数据库外建约束false，可以关闭设为true
		},
	)
	if err !=nil {
		klog.Error(err)
		return nil, err
	}
	sqlDb, _ := db.DB()
	// 关闭链接// 对于中小型 web 应用程序，我通常使用以下设置作为起点，然后根据负载测试结果和实际吞吐量级别进行优化。
	//SetMaxIdleConns: 设置空闲连接池中链接的最大数量sqlDb.SetMaxIdleConns(25)
	//SetMaxOpenConns: 设置打开数据库链接的最大数量sqlDb.SetMaxOpenConns(25)
	//SetConnMaxLifetime: 设置链接可复用的最大时间sqlDb.SetConnMaxLifetime(5 * time.Minute)
	// 将defer放入延迟调用栈
	//defer sqlDb.Close()
	//可用show variables like 'max_connections'; 查看服务器当前设置的最大连接数
	sqlDb.SetMaxOpenConns(150)//连接池最多同时打开的连接数，这个maxOpenConns理应要设置得比mysql服务器的max_connections值要小
	sqlDb.SetMaxIdleConns(50)//连接池里最大空闲连接数，比SetMaxOpenConns小
	//当连接持续空闲时长达到maxIdleTime后，该连接就会被关闭并从连接池移除，哪怕当前空闲连接数已经小于SetMaxIdleConns(maxIdleConns)设置的值
	//连接每次被使用后，持续空闲时长会被重置，从0开始从新计算
	//连接池里面的连接最大空闲时长
	//用show processlist; 可用查看mysql服务器上的连接信息，Command表示连接的当前状态，Command为Sleep时表示休眠、空闲状态，Time表示此状态的已持续时长:单位为秒
	sqlDb.SetConnMaxIdleTime(5 * time.Minute)
	//maxLifeTime必须要比mysql服务器设置的wait_timeout小，否则会导致golang侧连接池依然保留已被mysql服务器关闭了的连接
	//show variables like 'wait_timeout';
	sqlDb.SetConnMaxLifetime(1 * time.Hour)//连接池里面的连接最大存活时长
	data, _ := sonic.Marshal(sqlDb.Stats()) //获得当前的SQL配置情况
	klog.Debug("DB conf: ", string(data))

	//AutoMigrate 用于自动迁移您的 schema，保持您的 schema 是最新的。
	//自动迁移
	//var models = []interface{}{&User{}, &Note{}}
	//db.AutoMigrate(models)
	//db.AutoMigrate(&User{}, &Note{})

	//会话模式
	/*tx := DB.Session(&gorm.Session {
		PrepareStmt: true,
		SkipDefaultTransaction:   false,//默认开启默认事务false
		DisableNestedTransaction: true,//默认开启嵌套事务false
	})*/

	//db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().Create("")
	//手动迁移
	//db.AutoMigrate()
	//建表就放在各自模块了
	if err = db.Use(gormopentracing.New()); err != nil {
		hlog.Error(err)
		return nil, err
	}
	return db,nil
}
