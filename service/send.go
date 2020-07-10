package service

import (
	"errors"
	"time"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/zzsds/micro-sms-service/consts"
	"github.com/zzsds/micro-sms-service/models"
	"github.com/zzsds/micro-sms-service/modules/provider"
	"gitlab.bft.pub/welfare/common/utils"
)

const (
	defaultCode = "123456"
)

// SendRepo ...
type SendRepo struct {
}

func NewSendRepo() *SendRepo {
	return &SendRepo{}
}

// Create ...
func (repo *SendRepo) Create(sendModel *models.Send) error {
	tx := Db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	log.Logf(log.InfoLevel, "%v", sendModel)
	if err := tx.Create(sendModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CreateCode ...
func (repo *SendRepo) CreateCode(sendModel *models.Code) error {
	tx := Db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(sendModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Validate ...
func (repo *SendRepo) Validate(mobile, code string, bizType int32) (*models.Code, error) {
	var (
		codeModel models.Code
		sendModel models.Send
		now       = time.Now()
	)

	if err := Db.Where(models.Send{Mobile: mobile, BizType: bizType}).Last(&sendModel).Error; err != nil {
		return nil, err
	}

	if err := Db.Where(models.Code{SendID: sendModel.ID}).Last(&codeModel).Error; err != nil {
		return nil, err
	}

	if codeModel.IsUse == consts.Enabled_Yes {
		return &codeModel, errors.New("验证码已失效")
	}

	if codeModel.ExpiresAt.Before(now) {
		return &codeModel, errors.New("验证码已过期")
	}

	if err := repo.UseCode(codeModel.ID); err != nil {
		return nil, err
	} else {
		codeModel.IsUse = consts.Enabled_Yes
	}

	return &codeModel, nil
}

// UseCode ...
func (repo *SendRepo) UseCode(id uint) error {
	var codeModel models.Code
	codeModel.ID = id
	return Db.Model(codeModel).Update("is_use", uint(consts.Enabled_Yes)).Error
}

// Mobile ...
func (repo *SendRepo) Mobile(mobile string) (*models.Code, error) {
	var (
		codeModel models.Code
	)
	if err := Db.Where(codeModel).First(&codeModel).Error; err != nil {
		return nil, err
	}
	return &codeModel, nil
}

// SendCode ...
func (repo *SendRepo) SendCode(sendModel *models.Send, driver provider.Driver) error {
	code := utils.GenValidateCode(6)
	if driver.GetDebug() {
		code = defaultCode
	}
	sendModel.Content = driver.Assignment(sendModel.Content, code)
	sendModel.Code.Code = code
	if err := Db.Create(sendModel).Error; err != nil {
		return err
	}

	sid, err := driver.Send(sendModel.Mobile, sendModel.Content)
	if err != nil {
		Db.Model(models.Send{Model: sendModel.Model}).Update("message", err.Error())
		return err
	}

	if err := Db.Model(models.Send{Model: sendModel.Model}).Update(&models.Send{
		SID:     sid,
		Success: consts.Enabled_Yes,
	}).Error; err != nil {
		return err
	}
	sendModel.SID = sid
	sendModel.Success = consts.Enabled_Yes
	return nil
}

// SendNotice ...
func (repo *SendRepo) SendNotice(sendModel *models.Send, driver provider.Driver, value []string) error {
	for _, val := range value {
		sendModel.Content = driver.Assignment(sendModel.Content, val)
	}

	if err := repo.Create(sendModel); err != nil {
		return err
	}

	sid, err := driver.Send(sendModel.Mobile, sendModel.Content)
	if err != nil {
		Db.Model(sendModel).Update("message", err.Error())
		return err
	}

	if err := Db.Model(sendModel).Update(&models.Send{
		SID:     sid,
		Success: consts.Enabled_Yes,
	}).Error; err != nil {
		return err
	}

	return nil
}
