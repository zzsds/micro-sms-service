package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/micro/go-micro/v2"
	"github.com/zzsds/micro-sms-service/proto/yunpian"
)

func TestSend(t *testing.T) {
	srv := micro.NewService(
		micro.Name("go.micro.srv.wallet"),
	)
	srvClinet := micro.NewEvent("go.micro.srv.sms", srv.Client())
	err := srvClinet.Publish(context.Background(), &yunpian.EventResource{
		Mobile:  "13032368493",
		Mode:    1,
		BizType: 14,
		Value:   []string{fmt.Sprintf("%2.f", 33.2225)},
	})
	// 2.2.0.01
	t.Log(err)
	t.Log(err)
}
