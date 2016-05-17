package dao

import (
	"database/sql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path/filepath"
	// "sync"
)

var db *sql.DB

// var dbOnce sync.Once

//启动初始化
func init() {
	dbURL, err := getDBUrl()
	if err != nil {
		beego.Error("Error on parse database config, %s", err.Error())
		os.Exit(1)
	}
	beego.Debug(dbURL)
	// db, err = sql.Open("mysql", "xxx:xxx@tcp(xxx:3306)/xxx?charset=utf8")
	db, err = sql.Open("mysql", dbURL)
	if err != nil {
		beego.Error("Error on initializing database handle, %s", err.Error())
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		beego.Error("Error on opening database connection: %s", err.Error())
		os.Exit(1)
	}

	db.SetMaxIdleConns(50)
}

//The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
// Thus, the Open function should be called just once. I
//t is rarely necessary to close a DB.
//DB is a database handle representing a pool of zero or more underlying connections.
//It's safe for concurrent use by multiple goroutines.
//The sql package creates and frees connections automatically; it also maintains a free pool of idle connections
//because of all of these reasons, i choose the singleton pattern to init object DB.
func GetDB() *sql.DB {
	//单例模式，lazy load
	// dbOnce.Do(func() {
	// 	dbURL, err := getDBUrl()
	// 	if err != nil {
	// 		beego.Error("Error on parse database config, %s", err.Error())
	// 		os.Exit(1)
	// 	}
	// 	beego.Debug(dbURL)
	// 	// db, err = sql.Open("mysql", "haochushi:dev4chushi007.com@tcp(localhost:3306)/app?charset=utf8")
	// 	db, err = sql.Open("mysql", dbURL)
	// 	if err != nil {
	// 		beego.Error("Error on initializing database handle, %s", err.Error())
	// 		os.Exit(1)
	// 	}

	// 	err = db.Ping()
	// 	if err != nil {
	// 		beego.Error("Error on opening database connection: %s", err.Error())
	// 		os.Exit(1)
	// 	}

	// 	db.SetMaxIdleConns(50)
	// })
	return db
}

func getDBUrl() (string, error) {
	mysqlConfigPath := filepath.Join(beego.AppPath, "conf", "mysql.ini")
	mysqlConfig, err := config.NewConfig("ini", mysqlConfigPath)
	if err != nil {
		beego.Error("new mysql config error", err.Error())
		return "", err
	}

	envRunMode := os.Getenv("BEEGO_RUNMODE")
	if envRunMode == "" {
		envRunMode = beego.PROD
	}

	username := mysqlConfig.String(envRunMode + "::username")
	password := mysqlConfig.String(envRunMode + "::password")
	protocol := mysqlConfig.String(envRunMode + "::protocol")
	ip := mysqlConfig.String(envRunMode + "::ip")
	port := mysqlConfig.String(envRunMode + "::port")
	dbname := mysqlConfig.String(envRunMode + "::dbname")
	charset := mysqlConfig.String(envRunMode + "::charset")

	dbURL := username + ":" + password + "@" + protocol + "(" + ip + ":" + port + ")" + "/" + dbname + "?charset=" + charset

	return dbURL, nil
}
