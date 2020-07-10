package handler

import (
	"context"
	"regexp"
	"time"

	"github.com/micro/go-micro/v2/errors"
	"github.com/zzsds/micro-sms-service/conf"
	"github.com/zzsds/micro-sms-service/consts"
	"github.com/zzsds/micro-sms-service/models"
	"github.com/zzsds/micro-sms-service/modules/provider"
	"github.com/zzsds/micro-sms-service/proto/yunpian"
	"github.com/zzsds/micro-sms-service/service"
)

// Send ...
type YunPian struct {
	Repo         *service.SendRepo
	SmsProvider  provider.Driver
	templateRepo *service.TemplateRepo
}

func NewYunPian() *YunPian {
	return &YunPian{
		Repo:         service.NewSendRepo(),
		templateRepo: service.NewTemplateRepo(),
		SmsProvider:  provider.NewYunPian(conf.Conf.Sms.YunPian.Name, conf.Conf.Sms.YunPian.APIKey, conf.Conf.Sms.YunPian.Debug),
	}
}

// BizType is a single request handler called via client.Call or the generated client code
func (e *YunPian) BizType(ctx context.Context, req *yunpian.BizTypeResponse, rsp *yunpian.BizTypeResponse) error {
	var typeList = make([]*yunpian.MapType, 0, len(consts.SmsBizType_name))
	for key, val := range consts.SmsBizType_name {
		typeList = append(typeList, &yunpian.MapType{Key: key, Value: val})
	}
	rsp.List = typeList
	return nil
}

// Code is a single request handler called via client.Call or the generated client code
func (e *YunPian) Code(ctx context.Context, req *yunpian.CodeResource, rsp *yunpian.CodeResponse) error {
	var (
		expire      time.Duration
		now         = time.Now()
		smsProvider = e.SmsProvider
	)

	if req.Expires <= 0 {
		expire = 5 * time.Minute
	}
	expiresAt := now.Add(expire)
	if !conf.ValidateMobile(req.Mobile) {
		return errors.BadRequest("go.micro.srv.sms Code", "手机号格式错误")
	}

	templateModel, err := e.templateRepo.GetCodeTemplateFirst(smsProvider.String(), req.BizType)
	if err != nil {
		return errors.BadRequest("go.micro.srv.sms Code", err.Error())
	}

	if err := e.Repo.SendCode(&models.Send{
		Content:    provider.SignContent(templateModel.Sign, templateModel.Content),
		Mobile:     req.Mobile,
		Provider:   smsProvider.String(),
		BizType:    req.BizType,
		TemplateID: templateModel.ID,
		Code:       &models.Code{ExpiresAt: expiresAt}}, smsProvider); err != nil {
		return errors.BadRequest("go.micro.srv.sms", err.Error())
	}
	rsp.ExpiresAt = expiresAt.Format("2006-01-02 15:04:05")
	rsp.Success = true
	return nil
}

// Validate is a server side stream handler called via client.Stream or the generated client code
func (e *YunPian) Validate(ctx context.Context, req *yunpian.ValidateRequest, rsp *yunpian.ValidateResponse) (err error) {

	if ok, _ := regexp.MatchString(`([0-9]\d{5})$`, req.Code); !ok {
		return errors.BadRequest("go.micro.srv.sms Validate", "验证码格式错误")
	}
	if !conf.ValidateMobile(req.Mobile) {
		return errors.BadRequest("go.micro.srv.sms Validate", "手机号格式错误")
	}
	if _, err = e.Repo.Validate(req.Mobile, req.Code, req.BizType); err != nil {
		return errors.BadRequest("go.micro.srv.sms Validate", "%s", err.Error())
	}
	rsp.Success = true
	return nil
}

// Notice ...
func (e *YunPian) Notice(ctx context.Context, req *yunpian.EventResource, rsp *yunpian.EventResource) error {

	if !conf.ValidateMobile(req.Mobile) {
		return errors.BadRequest("go.micro.srv.sms Code", "手机号格式错误")
	}

	if len(req.Value) == 0 {
		return errors.BadRequest("go.micro.srv.sms Code", "参数不能为空")
	}

	templateModel, err := e.templateRepo.GetCodeTemplateFirst(e.SmsProvider.String(), req.BizType)
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
		return errors.BadRequest("go.micro.srv.sms", err.Error())
	}
	rsp.Success = true
	return nil
}
