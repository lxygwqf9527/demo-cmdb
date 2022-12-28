package impl_test

import (
	"context"
	"os"
	"testing"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger/zap"

	// 注册所有对象
	_ "github.com/lxygwqf9527/demo-cmdb/apps/all"
	"github.com/lxygwqf9527/demo-cmdb/apps/resource"
	"github.com/lxygwqf9527/demo-cmdb/apps/task"
	"github.com/lxygwqf9527/demo-cmdb/conf"
)

var (
	ins task.ServiceServer
)

func TestCreateTask(t *testing.T) {
	req := task.NewCreateTaskRequst()
	req.Type = task.Type_RESOURCE_SYNC
	req.Region = os.Getenv("TX_CLOUD_REGION")
	req.ResourceType = resource.Type_HOST
	req.SecretId = "cebvfibmiom0325ciso0"
	taskIns, err := ins.CreateTask(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(taskIns)
}

func init() {
	// 通过环境变量加载测试配置
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}

	// 全局日志对象初始化
	zap.DevelopmentSetup()

	// 初始化所有实例
	if err := app.InitAllApp(); err != nil {
		panic(err)
	}

	ins = app.GetGrpcApp(task.AppName).(task.ServiceServer)
}
