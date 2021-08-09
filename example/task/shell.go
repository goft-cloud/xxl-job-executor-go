package task

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/xxl-job/xxl-job-executor-go"
)

func Shell(cxt context.Context, param *xxl.RunReq) (msg string) {
	fmt.Println(param.GlueSource)

	var d1 = []byte(param.GlueSource)
	err2 := ioutil.WriteFile("./runtime/t.sh", d1, 0777)
	fmt.Println(err2)


	cmd := exec.Command("./runtime/t.sh", param.ExecutorParams)

	data, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("data" + string(data))
	return
}
