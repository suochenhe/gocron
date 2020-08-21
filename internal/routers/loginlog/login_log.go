package loginlog

import (
	"github.com/suochenhe/gocron/internal/models"
	"github.com/suochenhe/gocron/internal/modules/logger"
	"github.com/suochenhe/gocron/internal/modules/utils"
	"github.com/suochenhe/gocron/internal/routers/base"
	macaron "gopkg.in/macaron.v1"
)

func Index(ctx *macaron.Context) string {
	loginLogModel := new(models.LoginLog)
	params := models.CommonMap{}
	base.ParsePageAndPageSize(ctx, params)
	total, err := loginLogModel.Total()
	if err != nil {
		logger.Error(err)
	}
	loginLogs, err := loginLogModel.List(params)
	if err != nil {
		logger.Error(err)
	}

	jsonResp := utils.JsonResponse{}

	return jsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  loginLogs,
	})
}
