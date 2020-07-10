package subscriber

import (
	"context"

	"github.com/zzsds/micro-sms-service/conf"
	. "github.com/zzsds/micro-sms-service/conf"

	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/zzsds/micro-sms-service/consts"
	"github.com/zzsds/micro-sms-service/models"
	"github.com/zzsds/micro-sms-service/modules/provider"
	"github.com/zzsds/micro-sms-service/proto/yunpian"
	"github.com/zzsds/micro-sms-service/service"
)

// Send ...
type Send struct {
	Repo         *service.SendRepo
	SmsProvider  provider.Driver
	templateRepo *service.TemplateRepo
}

func NewSend() *Send {
	return &Send{
		templateRepo: service.NewTemplateRepo(),
		SmsProvider:  provider.NewYunPian(Conf.Sms.YunPian.Name, Conf.Sms.YunPian.APIKey, Conf.Sms.YunPian.Debug),
		Repo:         service.NewSendRepo(),
	}
}

// Handle ...
func (e *Send) Handle(ctx context.Context, req *yunpian.EventResource) error {
	log.Info("Handler Received message: ", req.Mobile)
	if req.Mode == consts.SmsMode_Code {
		return errors.BadRequest("go.micro.srv.sms Code", "code function not open")
	}
	if !conf.ValidateMobile(req.Mobile) {
		return errors.BadRequest("go.micro.srv.sms Code", "手机号格式错误")
	}

	templateModel, err := e.templateRepo.GetNoticeTemplateFirst(e.SmsProvider.String(), req.BizType)
	if err != nil {
		return errors.BadRequest("go.micro.srv.sms Code", err.Error())
	}
	if err := e.Repo.SendNotice(&models.Send{
		Content:    provider.SignContent(templateModel.Sign, templateModel.Content),
		Mobile:     req.Mobile,
		Provider:   e.SmsProvider.String(),
		BizType:    req.BizType,
		TemplateID: templateModel.ID,
	}, e.SmsProvider, req.Value); err != nil {
		return err
	}
	return nil
}
