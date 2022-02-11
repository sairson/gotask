package logger

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

// 设置常用颜色
const (
	Cyan   = color.FgCyan
	Blue   = color.FgBlue
	Red    = color.FgRed
	Green  = color.FgGreen
	Yellow = color.FgYellow
)

var Islog bool = true

// 设置日志调试等级
const (
	LevelInfo  = 1
	LevelWarn  = 2
	LevelFatal = 3
)

// 设置输出函数
func log(level int, details string) {
	if level < 1 || level > 3 {
		return
	}
	if Islog == true {
		fmt.Println(details)
	}

}

// 返回时间字符串
func unixTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Info 标准调试函数info
func Info(details string) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", unixTime(), color.GreenString("INFO"), details))
}

// Fatal 标准调试函数Fatal
func Fatal(details string) {
	log(LevelFatal, fmt.Sprintf("[%s] [%s] %s", unixTime(), color.RedString("FATAL"), details))
}
