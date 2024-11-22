package Constant

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

// MyTime 自定义时间
type SystemTime time.Time

func (t *SystemTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*t = SystemTime(t1)
	return err
}

func (t SystemTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t SystemTime) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format("2006-01-02 15:04:05"), nil
}

func (t *SystemTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = SystemTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *SystemTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

func (t *SystemTime) ToTime() time.Time {
	return time.Time(*t)
}
