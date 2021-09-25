package config

import (
	"strconv"
	"time"

	"github.com/spf13/viper"

	//driver
	"gopkg.in/mgo.v2"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	mongo "github.com/abcdef-id/go-lib/mgo"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	mocket "github.com/selvatico/go-mocket"
)

var (
	//MySqlDB MySqlDB
	MySqlDB *gorm.DB
	//Redis Redis
	Redis *redis.Client
	//Mgo Mgo
	Mgo *mongo.MongoDatabase
)

//Database Database
type Database struct {
	Host              string
	User              string
	Password          string
	DBName            string
	DBNumber          int
	Port              int
	API_URL           string
	ReconnectRetry    int
	ReconnectInterval int64
	DebugMode         bool
}

// LoadDBConfig load database configuration
func LoadDBConfig(name string) Database {
	db := viper.Sub("database." + name)
	conf := Database{
		Host:              db.GetString("host"),
		User:              db.GetString("user"),
		Password:          db.GetString("password"),
		DBName:            db.GetString("db_name"),
		DBNumber:          db.GetInt("db_number"),
		Port:              db.GetInt("port"),
		API_URL:           db.GetString("api_url"),
		ReconnectRetry:    db.GetInt("reconnect_retry"),
		ReconnectInterval: db.GetInt64("reconnect_interval"),
		DebugMode:         db.GetBool("debug"),
	}
	return conf
}

func OpenMySqlPool() {
	if viper.Get("env") != "testing" {
		MySqlDB = MysqlConnect("mysql")
	} else {
		MySqlDB = MysqlConnectTest("mysql")
	}
	pool := viper.Sub("database.mysql.pool")
	MySqlDB.DB().SetMaxOpenConns(pool.GetInt("maxOpenConns"))
	MySqlDB.DB().SetMaxIdleConns(pool.GetInt("maxIdleConns"))
	MySqlDB.DB().SetConnMaxLifetime(pool.GetDuration("maxLifetime") * time.Second)
}

// MysqlConnect connect to mysql using config name. return *gorm.DB incstance
func MysqlConnect(configName string) *gorm.DB {
	mysql := LoadDBConfig(configName)
	connection, err := gorm.Open("mysql", mysql.User+":"+mysql.Password+"@tcp("+mysql.Host+":"+strconv.Itoa(mysql.Port)+")/"+mysql.DBName+"?charset=utf8&parseTime=True&loc=Local")
	// enable debug
	// connection.LogMode(true)
	if err != nil {
		panic(err)
	}

	if mysql.DebugMode {
		return connection.Debug()
	}

	return connection
}

//MysqlConnectTest connect to mysql using config name. return *gorm.DB incstance
func MysqlConnectTest(configName string) *gorm.DB {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	connection, err := gorm.Open(mocket.DriverName, "connection_string")

	if err != nil {
		panic(err)
	}

	return connection
}

//MongoConnect MongoConnect
func MongoConnect() {
	if viper.Get("env") != "testing" {
		conf := LoadDBConfig("mongo")
		mongoConf := &mgo.DialInfo{
			Addrs:    []string{conf.Host},
			Username: conf.User,
			Password: conf.Password,
			Database: conf.DBName,
		}
		session, err := mgo.DialWithInfo(mongoConf)
		if err != nil {
			panic(err)
		}
		s := mongo.MongoSession{Session: session}
		Mgo = &mongo.MongoDatabase{Database: s.Session.DB(conf.DBName)}
	}
}
