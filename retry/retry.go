package retry

import (
	"time"
)

// 重试方法 attempts 调用次数， sleep 睡眠时间， callback 需要重试的方法
func Retry(attempts int, sleep time.Duration, callback func(curr int) error) (err error) {
	for curr := 0; ; curr++ {
		err = callback(curr + 1)
		if err == nil {
			return
		}
		if curr >= (attempts - 1) {
			break
		}
		time.Sleep(time.Second * time.Duration(curr+1) * sleep)
	}
	return err
}

// 重试方法 attempts 调用次数， sleep 睡眠时间， callback 需要重试的方法， resp 其他返回值
func RetryResp(attempts int, sleep time.Duration, callback func(curr int) (interface{}, error)) (resp interface{}, err error) {
	for curr := 0; ; curr++ {
		resp, err = callback(curr + 1)
		if err == nil {
			return
		}
		if curr >= (attempts - 1) {
			break
		}
		time.Sleep(time.Second * time.Duration(curr+1) * sleep)
	}
	return resp, err
}
