package handler

import (
	"context"

	"github.com/zzsds/micro-sms-service/service"

	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/errors"
	"github.com/zzsds/micro-sms-service/models"
	"github.com/zzsds/micro-sms-service/modules/provider"
	"github.com/zzsds/micro-sms-service/proto/sms"
)

// Sms ...
type Sms struct {
	Repo         *service.SendRepo
	SmsProvider  provider.Driver
	templateRepo *service.TemplateRepo
}

func NewSms() *Sms {
	return &Sms{templateRepo: service.NewTemplateRepo()}
}

//
// SendList ...
func (e *Sms) SendList(ctx context.Context, req *sms.SendPage, rsp *sms.SendPage) error {
	return nil
}

// TemplateList ...
func (e *Sms) TemplateList(ctx context.Context, req *sms.TemplatePage, rsp *sms.TemplatePage) error {
	return nil
}

// Template ...
func (e *Sms) Template(ctx context.Context, req *sms.SmsStruct, rsp *sms.SmsStruct) (err error) {
	result, err := e.templateRepo.FirstModelTemplate(models.Template{
		Model:    gorm.Model{ID: uint(req.Id)},
		Provider: req.Provider,
		Mode:     req.Mode,
		BizType:  req.BizType,
	})
	if err != nil {
		return errors.BadRequest("go.micro.srv.sms", err.Error())
	}
	rsp = e.DownloadObj(result)
	return nil
}

func (e *Sms) List(ctx context.Context, req *sms.ListRequest, rsp *sms.ListResponse) error {
	return nil
}

func (e *Sms) Create(ctx context.Context, in *sms.CreateTempRequest, rsp *sms.CreateTempResponse) error {
	return nil
}

func (e *Sms) DownloadObj(m *models.Template) (obj *sms.SmsStruct) {
	obj = new(sms.SmsStruct)
	obj.Id = int64(m.ID)
	obj.Provider = m.Provider
	obj.Sign = m.Sign
	obj.Content = m.Content
	obj.Mode = m.Mode
	return
}
