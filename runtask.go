package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/capnspacehook/taskmaster"
)

//将string类型的用户输入转换为time类型
func stringTotime(timemee string) (t time.Time) {
	timetemp := "2006-01-02 1a5:04:05"
	tm, err := time.ParseInLocation(timetemp, timemee, time.Local)
	if err != nil {
		panic(err)
	}
	return tm
}

//检查用户输入的时间并返回结果
func checkTime(str string) (t time.Time) {
	if str != "" {
		t := stringTotime(str)
		return t
	}
	tm := time.Now().Add(time.Minute * +1)
	return tm
}

func runWinTask(path string, args string, enable bool, tm time.Time) {
	//创建初始化计划任务
	taskService, err := taskmaster.Connect()
	if err != nil {
		fmt.Println(err)
	}

	defer taskService.Disconnect()
	//定义新的计划任务
	newTaskDef := taskService.NewTaskDefinition()
	//添加执行程序的路径和参数
	newTaskDef.AddAction(taskmaster.ExecAction{
		Path: path,
		Args: args,
	})
	//定义计划任务程序的执行时间等
	newTaskDef.AddTrigger(taskmaster.DailyTrigger{
		DayInterval: 1,
		TaskTrigger: taskmaster.TaskTrigger{
			Enabled:       enable,
			StartBoundary: tm,
		},
	})

	//创建计划任务
	resp, _, err := taskService.CreateTask("\\windows\\update", newTaskDef, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

}

func main() {
	// 命令行启动
	taskPath := flag.String("path", "", "指定计划任务启动程序的路径(必选)")
	taskArgs := flag.String("args", "", "添加计划任务启动程序的参数（可选）")
	flagTaskTime := flag.String("time", "", "添加计划任务启动的时间（可选，格式：\"2006-01-02 15:04:05\"），默认当前时间一分钟后执行")
	taskEnable := flag.Bool("enable", true, "计划任务是否启用（可选），默认为true")
	flag.Parse()

	if *taskPath == "" {
		fmt.Println("请使用-h寻求帮助")
		os.Exit(1)
	}
	// 检查时间
	taskTime := checkTime(*flagTaskTime)

	runWinTask(*taskPath, *taskArgs, *taskEnable, taskTime)
}
