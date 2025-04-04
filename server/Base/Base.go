package Base

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server/App/Constant"
	"server/App/Model/Common"
	"server/App/Model/Count"
	"server/App/Model/Message"
	"server/App/Model/Service"
	"server/App/Model/ServiceManager"
	"server/App/Model/User"
	"server/Base/Config"
	"server/Base/Nsq"
	"server/Base/WebSocket"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nsqio/go-nsq"
)

type Base struct{}

var AppConfig Config.App

var MysqlConn *gorm.DB

var WebsocketHub WebSocket.Hub

var RedisPool *redis.Pool //创建redis连接池

var Producer *nsq.Producer

func (b Base) Init() {

	b.initConfig()
	//b.InitConsumer()

	b.initMysql()
	b.initWebSocketHub()
	b.initRedis()
}

func (b Base) InitConsumer() {
	Producer = Nsq.NsqConsumer{}.CreateProducer(AppConfig.Mq.Nsq.Host)
	Nsq.NsqConsumer{}.InitConsumer(Constant.Topic, Constant.Channel, AppConfig.Mq.Nsq.Host)
}

func (b Base) initWebSocketHub() {
	WebsocketHub = WebSocket.Hub{
		UserListMap:        map[string]map[string]WebSocket.Connect{},
		UserConnGroupList:  map[string]map[string]WebSocket.Connect{},
		UserConnIdGroupMap: map[string]map[string]int{},
		ServiceBindUser:    map[string]int{},
		ServiceBindGroup:   map[string]int{},
	}
	WebsocketHub.Run()
}

// 配置初始化
func (b Base) initConfig() {
	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	res, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Print(err.Error())
		panic("json配置文件打开失败")
	}

	err = json.Unmarshal(res, &AppConfig)
	if err != nil {
		fmt.Print(err.Error())
		panic("json 配置解析异常")
	}
}

// mysql 初始化
func (b Base) initMysql() {

	var err error
	c := AppConfig.Database.Mysql
	connStr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Database)
	MysqlConn, err = gorm.Open("mysql", connStr)
	if err != nil {
		fmt.Print(err.Error())
		panic("mysql 初始异常")
	}

	auto := MysqlConn.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

	auto.AutoMigrate(&Common.Domain{}, &Common.WeChatAuth{})
	auto.AutoMigrate(&Message.Message{})
	auto.AutoMigrate(&Service.Service{}, &Service.ServiceRoomDetail{}, &Service.ServiceAuth{}, &Service.ServiceBlack{}, &Service.ServiceMessage{}, &Service.ServiceRoom{})
	auto.AutoMigrate(&User.User{}, &User.UserLoginLog{})
	auto.AutoMigrate(&Count.CountServiceRoom{})
	auto.AutoMigrate(&ServiceManager.ServiceManager{}, &ServiceManager.ServiceManagerMessage{})
	auto.AutoMigrate(&ServiceManager.ServiceManagerBot{}, &ServiceManager.ServiceManagerBotMessage{})
	auto.AutoMigrate(Common.SystemConfig{}, Common.Setting{}, &Common.Rename{})
	auto.AutoMigrate(&Service.Order{})
	auto.AutoMigrate(&Service.ServiceMenuSetting{})
	auto.AutoMigrate(&Service.ServiceNoticeSetting{})
	auto.AutoMigrate(&ServiceManager.ServiceManagerRenewRecorder{})
	auto.AutoMigrate(&ServiceManager.ServiceManagerAuth{})
	auto.AutoMigrate(&Common.Cookie{})
	auto.AutoMigrate(&User.UserAuthMap{})

	MysqlConn.DB().SetMaxIdleConns(10)
	MysqlConn.DB().SetMaxOpenConns(100)

	MysqlConn.DB().SetConnMaxLifetime(20 * time.Second)
	MysqlConn.DB().SetMaxOpenConns(100)

	MysqlConn.LogMode(AppConfig.Debug)
}

func (b Base) initRedis() {
	redisConfig := AppConfig.Database.Redis
	config := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	RedisPool = &redis.Pool{ //实例化一个连接池
		MaxIdle:     16,  //最初的连接数量
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			option := redis.DialPassword(redisConfig.Password)
			return redis.Dial("tcp", config, option)
		},
	}
}
