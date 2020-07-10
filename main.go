package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/zzsds/micro-sms-service/conf"
	"github.com/zzsds/micro-sms-service/handler"
	"github.com/zzsds/micro-sms-service/proto/sms"
	"github.com/zzsds/micro-sms-service/proto/yunpian"
	"github.com/zzsds/micro-sms-service/service"
	"github.com/zzsds/micro-sms-service/subscriber"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.srv.sms"),
		micro.Version("latest"),
	)
	srv.Init()

	conf.InitConfig()
	service.InitConnectionMysql()

	defer service.Db.Close()

	{
		conf.CheckErr(yunpian.RegisterSendHandler(srv.Server(), handler.NewYunPian()))
		conf.CheckErr(sms.RegisterSmsHandler(srv.Server(), handler.NewSms()))
	}

	{
		conf.CheckErr(micro.RegisterSubscriber("go.micro.srv.sms", srv.Server(), subscriber.NewSend()))
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
