package impl

import (
	"context"
	"fmt"

	"github.com/lxygwqf9527/demo-cmdb/apps/host"
	"github.com/lxygwqf9527/demo-cmdb/apps/resource"
	"github.com/lxygwqf9527/demo-cmdb/apps/secret"
	"github.com/lxygwqf9527/demo-cmdb/apps/task"
	"github.com/lxygwqf9527/demo-cmdb/conf"
	"github.com/lxygwqf9527/demo-cmdb/provider/txyun/connectivity"
	"github.com/lxygwqf9527/demo-cmdb/provider/txyun/cvm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskService struct {
}

func (i *impl) CreateTask(ctx context.Context, req *task.CreateTaskRequst) (*task.Task, error) {
	// 查询secret
	s, err := i.secret.DescribeSecret(ctx, secret.NewDescribeSecretRequest(req.SecretId))
	if err != nil {
		return nil, err
	}
	// 解密api secret
	if s.Data.DecryptAPISecret(conf.C().App.EncryptKey); err != nil {
		return nil, err
	}
	// t := task.NewDefaultTask()
	// t.Run()

	switch req.Type {
	case task.Type_RESOURCE_SYNC:
		// 利用secret的信息，初始化一个operater
		switch s.Data.Vendor {
		case resource.Vendor_TENCENT:
			switch req.ResourceType {
			case resource.Type_HOST:
				// 初始化腾讯云cvm operator
				txConn := connectivity.NewTencentCloudClient(s.Data.ApiKey, s.Data.ApiSecret, req.Region)
				cvmOp := cvm.NewCVMOperator(txConn.CvmClient())
				// 分页查询
				pagger := cvm.NewPagger(float64(s.Data.RequestRate), cvmOp)
				for pagger.Next() {
					set := host.NewHostSet()
					// 查询分页有错误就立即返回
					if pagger.Scan(ctx, set); err != nil {
						return nil, err
					}
					// 保存该页数据,同步时，记录下日志
					for index := range set.Items {
						_, err := i.host.SyncHost(ctx, set.Items[index])
						if err != nil {
							i.log.Errorf("sync host error,%s", err)
							continue
						}
					}
				}
			case resource.Type_RDS:
			case resource.Type_BILL:
			}
		case resource.Vendor_ALIYUN:
		case resource.Vendor_HUAWEI:
		case resource.Vendor_AMAZON:
		case resource.Vendor_VSPHERE:
		default:
			return nil, fmt.Errorf("unknow resource type: %s", req.Type)
		}
		//
	case task.Type_RESOURCE_RELEASE:
	default:
		return nil, fmt.Errorf("unknow task type: %s", req.Type)
	}
	return nil, nil
}
func (i *TaskService) QueryBook(context.Context, *task.QueryTaskRequest) (*task.TaskSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryBook not implemented")
}
func (i *TaskService) DescribeBook(context.Context, *task.DescribeTaskRequest) (*task.Task, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeBook not implemented")
}
