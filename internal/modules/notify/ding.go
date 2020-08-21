package notify

// 发送消息到ding

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/suochenhe/gocron/internal/models"
	"github.com/suochenhe/gocron/internal/modules/httpclient"
	"github.com/suochenhe/gocron/internal/modules/logger"
	"github.com/suochenhe/gocron/internal/modules/utils"
)

type Ding struct{}

func (ding *Ding) Send(msg Message) {
	model := new(models.Setting)
	dingSetting, err := model.DingNotify()
	if err != nil {
		logger.Error("#ding#从数据库获取钉钉配置失败", err)
		return
	}
	if dingSetting.Url == "" {
		logger.Error("#ding#url为空")
		return
	}
	if len(dingSetting.Users) == 0 {
		logger.Error("#ding#用户配置为空")
		return
	}
	logger.Debugf("%+v", dingSetting)
	ding.send(msg, dingSetting)
}

func (ding *Ding) send(msg Message, dingSetting models.Ding) {
	formatBody := ding.format(msg, dingSetting)
	dingUrl := dingSetting.Url
	timeout := 30
	maxTimes := 3
	i := 0
	for i < maxTimes {
		resp := httpclient.PostJson(dingUrl, formatBody, timeout)
		logger.Debugf("status code: %d", resp.StatusCode)
		if resp.StatusCode == 200 {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			logger.Errorf("ding#发送消息失败#%s#消息内容-%s", resp.Body, msg["content"])
		}
	}
}

func map2Json(body map[string]interface{}) string {
	jsonString, _ := json.Marshal(body)
	return string(jsonString)
}

//格式化消息内容
func (ding *Ding) format(msg Message, dingSetting models.Ding) string {
	resultText := parseNotifyTemplate(dingSetting.Template, msg)
	body := map[string]interface{}{"msgtype": "text"}
	body["text"] = map[string]string{"content": "【定时任务通知】\n" + resultText}

	taskReceiverId := msg["task_receiver_id"].(string)

	switch taskReceiverId {

	case models.DingNotifyNone:
	case models.DingNotifyAll:
		body["at"] = map[string]bool{"isAtAll": true}
	default:
		taskReceiverIds := strings.Split(taskReceiverId, ",")
		var users []string
		for _, v := range dingSetting.Users {
			if utils.InStringSlice(taskReceiverIds, strconv.Itoa(v.Id)) {
				users = append(users, v.Name)
			}
		}
		body["at"] = map[string][]string{"atMobiles": users}
	}
	return map2Json(body)
}
