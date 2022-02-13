package service

import (
	"fmt"
	"github.com/google/uuid"
	jobs2 "github.com/sairson/gotask/service/jobs"
	"github.com/sairson/gotask/service/logger"
	"time"
)


var pool int

// InitAsync 队列初始化
func InitAsync(num int, work int,host string,port int,pass string) error {
	err := jobs2.InitRedisPool(host,port,pass)
	if err != nil {
		return err
	}
	// 先初始化数据库连接,新建队列
	err = jobs2.InitAsync(num)
	if err != nil {
		return err
	}
	pool = work
	return nil
}

// Debug 日志管理
func Debug(i bool){
	if i == true {
		logger.Islog = true
	}else{
		logger.Islog = false
	}
}

// Wait 队列运行
func Wait(tick time.Duration){
	var done = make(chan bool, pool)
	// 做死循环
	for range time.Tick(tick) {
		if len(jobs2.Jobs) == 0 {
			logger.Info("任务队列为空")
			continue
		} else {
			length := len(jobs2.Jobs)
			logger.Info(fmt.Sprintf("当前任务队列长度为%v", length))
			// 每次任务队列不为空时，判断队列长度
			if pool >= length {
				pool = length
			}
			for i := 0; i < pool; i++ {
				j := <-jobs2.Jobs
				go func() {
					logger.Info(fmt.Sprintf("开始执行任务ID:%v", j.UUID))
					// 更改任务执行状态
					j.Status = 2 // 2为执行状态
					err := jobs2.RedisSet(j.UUID, j) // 设置redis的key和vaule
					if err != nil {
						logger.Fatal(err.Error())
						return
					}
					// 开始执行函数方法
					_, err = j.Start() // 执行任务
					done <- true       // 执行完成后才给与true
					if err != nil {
						logger.Fatal(err.Error())
						return
					}
				}()
			}
			// 控制执行数量，阻塞，直到任务执行完毕
			for i := 0; i < pool; i++ {
				<-done
			}
		}
	}
}


// Invoke 函数调用
func Invoke(FuncName string, note string ,param ...map[string]interface{}) (string, error) {
	var j jobs2.Job
	j.DateTime = time.Now().Format("2006-01-02 15:04:05")
	j.Note = note
	j.FuncParams = param
	// 将任务添加进去
	if jobs2.RegisterFunc[FuncName] == nil {
		// 没有找到这个方法
		return "", fmt.Errorf("not register function")
	}
	j.FuncName = FuncName
	j.UUID = uuid.New().String()
	j.Status = 1    // 1为等待状态
	jobs2.AddJob(j) // 添加一个任务到队列当中

	return j.UUID, nil // 返回任务id
}


func Register(FuncName string,FuncMethod interface{}){
	jobs2.RegisterFunc[FuncName] = FuncMethod
}

// GetTaskStatus 通过uuid获取任务状态
func GetTaskStatus(uuid string) int {
	return jobs2.StatusJob(uuid)
}

// GetTaskResult 通过uuid获取任务结果
func GetTaskResult(uuid string) []interface{} {
	return jobs2.GetJob(uuid).FuncResult
}

// GetAllTask 获取全部任务,包括完成和未完成任务
func GetAllTask() ([]jobs2.Job,error) {
	return jobs2.AllJobs()
}

// GetTask 通过uuid来获取对应job
func GetTask(uuid string) (jobs2.Job) {
	return jobs2.GetJob(uuid)
}

func RemoveTask(uuid string)error {
	return jobs2.RemoveJob(uuid)
}