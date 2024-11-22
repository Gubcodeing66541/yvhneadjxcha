package Logic

import (
	"server/Base"
)

type System struct{}

func (System) getClearCache() []string {
	return []string{
		WeChat{}.GetKey(),
	}
}

func (s System) ClearCache() {
	conn := Base.RedisPool.Get()
	defer conn.Close()
	_, _ = conn.Do("flushall")
}
