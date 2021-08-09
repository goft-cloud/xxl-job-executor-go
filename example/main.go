package main

import (
	"log"

	xxl "github.com/xxl-job/xxl-job-executor-go"
	"github.com/xxl-job/xxl-job-executor-go/example/task"
)

func main() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://127.0.0.1:8686/xxl-job-admin/"),
		xxl.AccessToken(""),         //请求令牌(默认为空)
		xxl.ExecutorIp("127.0.0.1"), //可自动获取
		xxl.ExecutorPort("9999"),    //默认9999（非必填）
		xxl.RegistryKey("golang-job"),
		xxl.SetLogger(xxl.NewDefaultLogger()), //自定义日志
	)
	exec.Init()
	//设置日志查看handler
	exec.LogHandler(func(req *xxl.LogReq) *xxl.LogRes {
		return &xxl.LogRes{Code: 200, Msg: "", Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "这个是自定义日志handler",
			IsEnd:       true,
		}}
	})
	//注册任务handler
	exec.RegTask("task.test", task.Test)
	exec.RegTask("task.test2", task.Test2)
	exec.RegTask("task.panic", task.Panic)
	exec.RegTask("task.shell", task.Shell)
	log.Fatal(exec.Run())
}
