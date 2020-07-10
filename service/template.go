package service

import (
	"github.com/jinzhu/gorm"
	"github.com/zzsds/micro-sms-service/models"
	"gitlab.bft.pub/welfare/common/proto/define"
)

// TemplateRepo ...
type TemplateRepo struct {
}

func NewTemplateRepo() *TemplateRepo {
	return &TemplateRepo{}
}

// Create ...
func (r *TemplateRepo) Create(templateModel *models.Template) error {
	tx := Db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&templateModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetCodeTemplateFirst ...
func (r *TemplateRepo) GetCodeTemplateFirst(provider string, bizType int32) (*models.Template, error) {
	var (
		templateModel models.Template
		err           error
	)
	templateModel.Provider = provider
	dbx := Db.Scopes(Scope("enabled", define.Enabled_Yes), Scope("mode", define.SmsMode_Code), Scope("biz_type", bizType))
	if err = dbx.Where(&templateModel).First(&templateModel).Error; err != nil {
		return nil, err
	}
	return &templateModel, nil
}

// GetNoticeTemplateFirst ...
func (r *TemplateRepo) GetNoticeTemplateFirst(provider string, bizType int32) (*models.Template, error) {
	var (
		templateModel models.Template
		err           error
	)
	templateModel.Provider = provider
	dbx := Db.Scopes(Scope("enabled", define.Enabled_Yes), Scope("mode", define.SmsMode_Notice), Scope("biz_type", bizType))
	if err = dbx.Where(&templateModel).First(&templateModel).Error; err != nil {
		return nil, err
	}
	return &templateModel, nil
}

func (r *TemplateRepo) FirstModelTemplate(m models.Template) (result *models.Template, err error) {
	result = new(models.Template)
	if err = Db.Where(models.Template{
		Model:    gorm.Model{ID: m.ID},
		Provider: m.Provider,
		Mode:     m.Mode,
		BizType:  m.BizType,
	}).First(result).Error; err != nil {
		return
	}
	return
}
