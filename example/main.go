package main

import (
	"fmt"
	"github.com/sairson/gotask"
	"time"
)

func AddFunc(a,b int) int {
	return a + b
}

func main(){
	//service.Debug(false)

	err := TaskConntroller.InitAsync(100, 3, "192.168.248.219", 6379, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(6 * time.Second)
	TaskConntroller.Register("add",AddFunc) // 注册函数
	send, err := TaskConntroller.Invoke("add", "这是测试函数add",[]map[string]interface{}{
		{"int": 1},
		{"int": 2},
	}...) // 调用函数
	if err != nil {
		return
	}
	fmt.Println(send)
	TaskConntroller.Wait(3 * time.Second)
}