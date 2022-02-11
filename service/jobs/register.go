package jobs

import (
	"fmt"
	"strings"
)

// RegisterFunc 注册的函数，用于任务队列查找,函数调用方法用小写实现
var RegisterFunc = map[string]interface{}{}

func interfaceTypeConversion(value []map[string]interface{}) ([]map[string]interface{},error) {
	for l,i := range value{
		for k,v := range i {
			switch strings.ToLower(k) {
			case "[]string":
				var t []string
				for _,n := range v.([]interface{}) {
					t = append(t,n.(string))
				}
				value[l][k] = t
				break
			case "[]int":
				var t []int
				for _,n := range v.([]interface{}) {
					t = append(t,int(n.(float64)))
				}
				value[l][k] = t
				break
			case "[]float64":
				var t []float64
				for _,n := range v.([]interface{}) {
					t = append(t,n.(float64))
				}
				value[l][k] = t
				break
			case "string":
				value[l][k] = v.(string)
				break
			case "int":
				value[l][k] = int(v.(float64))
				break
			case "float64":
				value[l][k] = v.(float64)
				break
			case "bool":
				value[l][k] = v.(bool)
				break
			default:
				return nil,fmt.Errorf("not support type")
			}
		}
	}
	return value,nil
}