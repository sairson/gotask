# gotask
通过redis实现的简单任务队列（带有持久化），用于任务下发
![image](https://user-images.githubusercontent.com/74412075/153556999-b7489265-796b-43a0-a3d0-787457c429ac.png)

目前支持`int，string,[]string,[]int,bool，float64，[]float64`类型作为函数参数


使用文档：
```
go get github.com/sairson/gotask/service
```
example：

```
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
	//service.Debug(false)

	err := service.InitAsync(100, 3, "192.168.248.219", 6379, "")
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
```

函数支持
```
// GetTaskStatus 通过uuid获取任务状态
func GetTaskStatus(uuid string) int 

// GetTaskResult 通过uuid获取任务结果
func GetTaskResult(uuid string) []interface{} 

// GetAllTask 获取全部任务,包括完成和未完成任务
func GetAllTask() ([]jobs2.Job,error) 

// GetTask 通过uuid来获取对应任务
func GetTask(uuid string) jobs2.Job

// RemoveTask 通过uuid来移除执行完毕或还未执行的任务
func RemoveTask(uuid string)error
```
