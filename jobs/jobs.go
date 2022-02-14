package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sairson/gotask/logger"
	"reflect"
)

// 定义每一个作业需要的结构体



var Jobs chan Job

type Job struct {
	DateTime string // 任务创建时间
	Note 	 string // 任务备注
	UUID     string // 任务名称
	FuncName string
	// FuncMethod interface{} // 函数方法
	Status int //作业状态
	/*
		1 --->  等待
		2 --->  执行中
		3 --->  完成
	*/
	FuncParams []map[string]interface{} // 调用参数
	FuncResult []interface{}
}



// Start  开始执行一个方法
func (j *Job) Start() ([]reflect.Value, error) {
	defer func() {
		j.Status = 3                // 执行完毕，准备返回结果
		err := RedisSet(j.UUID, *j) // 数据库读写，更改任务执行状态
		if err != nil {
			return
		}
	}()
	// 找到名为name的函数
	//var f reflect.Value
	f := reflect.ValueOf(RegisterFunc[j.FuncName])
	if len(j.FuncParams) != f.Type().NumIn() { // 如果参数的值不等于函数所需要的值
		return nil, fmt.Errorf("call params lengths is not func need")
	}
	args := make([]reflect.Value, len(j.FuncParams)) // 获取参数
	// 添加参数
	for k, param := range j.FuncParams {
		for _,v := range param {
			if v == "" || v == 0 {
				continue
			}
			args[k] = reflect.ValueOf(v)
		}
	}
	r := f.Call(args) //调用函数
	//err := recover()
	//fmt.Println(r)
	for _, k := range r {
		j.FuncResult = append(j.FuncResult, k.Interface())
	}
	//j.FuncResult = append(j.FuncResult,) // 将执行结果填入
	return r, nil
}

func InitAsync(num int) error {
	Jobs = make(chan Job, num)
	logger.Info(fmt.Sprintf("初始化队列成功,队列长度%v", num))
	r, err := RedisDB.Do(context.Background(), "keys", "*").Result()
	if err != nil {
		return err
	}
	// 将没有执行完和等待中的任务从新加入到队列中,重新执行
	for _, i := range r.([]interface{}) {
		rs, err := RedisDB.Get(context.Background(), i.(string)).Result()
		if err != nil {
			return err
		}
		var job Job
		// 新建一个编码器
		decoder := json.NewDecoder(bytes.NewReader([]byte(rs)))
		//.UseNumber()
		err = decoder.Decode(&job)
		if err != nil {
			return err
		}
		job.FuncParams,err = InterfaceTypeConversion(job.FuncParams)
		if err != nil {
			return err
		}
		// 将等待的或者未执行完的重新加入队列执行
		if job.Status == 1 || job.Status == 2 {
			job.Status = 1 // 将job的status重新赋值为1，再次加入队列
			AddJob(job)
		}
	}
	return nil
}

func AddJob(job Job) {
	// 数据库添加失败
	err := RedisSet(job.UUID, job)
	if err != nil {
		logger.Fatal("任务添加失败")
		return
	}
	Jobs <- job // 将任务添加到队列当中
	logger.Info("任务添加成功")

}

// StatusJob 获取uuid对应job的执行状态
func StatusJob(uuid string) int {
	// 通过uuid获取执行状态
	return RedisGet(uuid).Status
}

func AllJobs() ([]Job,error) {
	// 返回所有jobs，包括完成和未完成jobs
	var ExistJobs []Job
	r, err := RedisDB.Do(context.Background(), "keys", "*").Result()
	if err != nil {
		return []Job{},err
	}
	// 将没有执行完和等待中的任务重新加入到队列中,重新执行
	for _, i := range r.([]interface{}) {
		rs, err := RedisDB.Get(context.Background(), i.(string)).Result()
		if err != nil {
			return []Job{},err
		}
		var job Job
		// 新建一个编码器
		decoder := json.NewDecoder(bytes.NewReader([]byte(rs)))

		err = decoder.Decode(&job)
		if err != nil {
			return []Job{},err
		}
		// 获取存储的json格式参数，以方便后续使用
		job.FuncParams,err = InterfaceTypeConversion(job.FuncParams)
		if err != nil {
			return []Job{},err
		}
		// 将等待的或者未执行完的重新加入队列执行
		ExistJobs = append(ExistJobs,job)
	}
	return ExistJobs,nil
}

// GetJob 返回一个job
func GetJob(uuid string) Job {
	return RedisGet(uuid)
}


func RemoveJob(uuid string) error {
	// 移除还未执行的job
	if len(Jobs) != 0 {
		for j := range Jobs {
			if j.UUID == uuid {
				err := RedisDel(uuid)
				if err != nil {
					return err
				}
				break
			}else{
				Jobs <- j
				continue
			}
		}
	}

	err := RedisDel(uuid)
	if err != nil {
		return err
	}
	return nil
}