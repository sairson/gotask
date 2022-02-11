package main

import (
	"fmt"
	"github.com/sairson/gotask/service"
	"time"
)

func AddFunc(a,b int) int {
	return a + b
}

func main(){
	service.Debug(false)

	err := service.InitAsync(100, 5*time.Second, 3, "192.168.248.219", 6379, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(6 * time.Second)
	service.Register("add",AddFunc) // 注册函数
	send, err := service.Invoke("add", []map[string]interface{}{
		{"int": 1},
		{"int": 2},
	}...) // 调用函数
	if err != nil {
		return
	}
	fmt.Println(send)
	service.Wait(3 * time.Second)
}