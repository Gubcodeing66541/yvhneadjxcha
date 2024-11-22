package Logic

import (
	"encoding/json"
	"fmt"
	Common2 "server/App/Common"
	"server/App/Model/Common"
	"server/Base"
)

type WeChat struct{}

func (WeChat) GetKey() string {
	return "REDIS-CONNECT-KEY"
}

func (w WeChat) GetAuth() (Common.WeChatAuth, error) {
	var whChat Common.WeChatAuth
	res := Common2.RedisTools{}.GetString(w.GetKey())
	if res != "" {
		err := json.Unmarshal([]byte(res), &whChat)
		if err != nil {
			fmt.Print("REDIS ERR CONNECT")
		}
		return whChat, nil
	}
	Base.MysqlConn.Find(&whChat, "type = 'auth' and status = 'enable'")
	jsonStr, err := json.Marshal(whChat)
	if err != nil {
		return whChat, err
	}
	Common2.RedisTools{}.SetString(w.GetKey(), string(jsonStr))
	return whChat, nil
}

func (w WeChat) ClearCache() {
	Common2.RedisTools{}.SetString(w.GetKey(), "")
}
