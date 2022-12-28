package cvm_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	"github.com/lxygwqf9527/demo-cmdb/apps/host"
	"github.com/lxygwqf9527/demo-cmdb/provider/txyun/connectivity"
	"github.com/lxygwqf9527/demo-cmdb/provider/txyun/cvm"
	tx_cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var (
	op *cvm.CVMOperator
)

func TestQuery(t *testing.T) {
	var (
		offset int64 = 1
		limit  int64 = 20
	)
	req := tx_cvm.NewDescribeInstancesRequest()
	req.Offset = &offset
	req.Limit = &limit
	set, err := op.Query(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestPaggerQuery(t *testing.T) {
	p := cvm.NewPagger(5, op)

	for p.Next() {
		set := host.NewHostSet()
		if err := p.Scan(context.Background(), set); err != nil {
			panic(err)
		}
		t.Log(set)
	}
}

func init() {
	err := connectivity.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
	zap.DevelopmentSetup()
	op = cvm.NewCVMOperator(connectivity.C().CvmClient())
}
