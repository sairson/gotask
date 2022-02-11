# gotask
The simple task queue is stripped when the program is written to achieve the task delivery function, which is used together with Redis<br>
利用redis实现的简单任务队列（带有持久化），用于任务下发
![image](https://user-images.githubusercontent.com/74412075/153556999-b7489265-796b-43a0-a3d0-787457c429ac.png)
目前支持int，string,[]string,[]int,bool类型作为函数参数

使用
```
go get github.com/sairson/gotask
```

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
