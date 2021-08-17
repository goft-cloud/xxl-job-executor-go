package xxl

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime/debug"
)

// TaskFunc 任务执行函数
type TaskFunc func(cxt context.Context, param *RunReq) string

// Task 任务
type Task struct {
	Id        int64
	Name      string
	Ext       context.Context
	Param     *RunReq
	fn        TaskFunc
	Cancel    context.CancelFunc
	StartTime int64
	EndTime   int64
	//日志
	log Logger
}

// Run 运行任务
func (t *Task) Run(callback func(code int64, msg string)) {
	t.log.InfoJob(t.Param.LogID, TaskLogStart)
	defer func(cancel func()) {
		if err := recover(); err != nil {
			t.log.Info(t.Info()+" panic: %v", err)
			debug.PrintStack() //堆栈跟踪
			callback(500, "task panic:"+fmt.Sprintf("%v", err))
			cancel()
		}
	}(t.Cancel)
	msg := t.fn(t.Ext, t.Param)
	callback(200, msg)
	t.log.InfoJob(t.Param.LogID, TaskLogEnd)
	return
}

// Info 任务信息
func (t *Task) Info() string {
	return "任务ID[" + Int64ToStr(t.Id) + "]任务名称[" + t.Name + "]参数：" + t.Param.ExecutorParams
}

type ShellTask struct {
	RunPath string
	Log     Logger
}

func (st *ShellTask) Shell(cxt context.Context, param *RunReq) (msg string) {
	jobId := param.LogID
	if len(param.GlueSource) == 0 {
		return "GlueSource 执行脚本内容不能为空！"
	}

	if len(st.RunPath) == 0 {
		st.RunPath = DefaultRunPath
	}

	if !IsDir(st.RunPath) {
		CreateDir(st.RunPath)
	}

	// 执行脚本路径
	shellPath := fmt.Sprintf("%s/%d.sh", st.RunPath, param.JobID)

	shellContext := []byte(param.GlueSource)
	err := ioutil.WriteFile(shellPath, shellContext, 0777)
	if err != nil {
		return "生成执行文件失败 file=" + shellPath + " err=" + err.Error()
	}

	cmd := exec.Command(shellPath, param.ExecutorParams)
	data, err2 := cmd.Output()
	if err2 != nil {
		return "脚本执行失败 file=" + shellPath + " err=" + err2.Error()
	}

	st.Log.InfoJob(jobId, "脚本执行成功")
	st.Log.InfoJob(jobId, string(data))
	return
}
