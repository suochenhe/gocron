package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Setting struct {
	Id    int    `xorm:"int pk autoincr"`
	Code  string `xorm:"varchar(32) notnull"`
	Key   string `xorm:"varchar(64) notnull"`
	Value string `xorm:"varchar(4096) notnull default '' "`
}

const dingTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
`

const slackTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
`
const emailTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
`
const webhookTemplate = `
{
  "task_id": "{{.TaskId}}",
  "task_name": "{{.TaskName}}",
  "status": "{{.Status}}",
  "result": "{{.Result}}"
}
`

const (
	DingCode        = "ding"
	DingUrlKey      = "url"
	DingTemplateKey = "template"
	DingUserKey     = "user"
	DingNotifyNone  = "-"
	DingNotifyAll   = "all"
)

const (
	SlackCode        = "slack"
	SlackUrlKey      = "url"
	SlackTemplateKey = "template"
	SlackChannelKey  = "channel"
)

const (
	MailCode        = "mail"
	MailTemplateKey = "template"
	MailServerKey   = "server"
	MailUserKey     = "user"
)

const (
	WebhookCode        = "webhook"
	WebhookTemplateKey = "template"
	WebhookUrlKey      = "url"
)

// 初始化基本字段 邮件、slack等
func (setting *Setting) InitBasicField() {
	setting.Code = DingCode
	setting.Key = DingUrlKey
	setting.Value = ""
	Db.Insert(setting)
	setting.Id = 0

	setting.Code = DingCode
	setting.Key = DingTemplateKey
	setting.Value = dingTemplate
	Db.Insert(setting)
	setting.Id = 0

	setting.Code = SlackCode
	setting.Key = SlackUrlKey
	setting.Value = ""
	Db.Insert(setting)
	setting.Id = 0

	setting.Code = SlackCode
	setting.Key = SlackTemplateKey
	setting.Value = slackTemplate
	Db.Insert(setting)
	setting.Id = 0

	setting.Code = MailCode
	setting.Key = MailServerKey
	setting.Value = ""
	Db.Insert(setting)
	setting.Id = 0

	setting.Code = MailCode
	setting.Key = MailTemplateKey
	setting.Value = emailTemplate
	Db.Insert(setting)
	setting.Id = 0

	setting.Code = WebhookCode
	setting.Key = WebhookTemplateKey
	setting.Value = webhookTemplate
	Db.Insert(setting)
	setting.Id = 0

	setting.Code = WebhookCode
	setting.Key = WebhookUrlKey
	setting.Value = ""
	Db.Insert(setting)
}

// region ding配置
type Ding struct {
	Url      string    `json:"url"`
	Users []DingUser `json:"users"`
	Template string    `json:"template"`
}

type DingUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func getDingUserMobile(value string) string {
	r := strings.Split(value, "@")
	return r[1]
}

func getDingUserName(value string) string {
	r := strings.Split(value, "@")
	return r[0]
}

func getDingUserValue(name string, mobile string) string {
	value := fmt.Sprintf("%s%s%s", name, "@", mobile)
	return value
}

func formatDingUser(du *DingUser, ds Setting) {
	du.Id = ds.Id
	du.Name = getDingUserName(ds.Value)
}

func formatDingMobile(du *DingUser, ds Setting) {
	du.Id = ds.Id
	du.Name = getDingUserMobile(ds.Value)
}


func (setting *Setting) formatDing(list []Setting, ding *Ding) {
	for _, v := range list {
		switch v.Key {
		case DingUrlKey:
			ding.Url = v.Value
		case DingTemplateKey:
			ding.Template = v.Value
		default:
			du := DingUser{}
			formatDingUser(&du, v)
			ding.Users = append(ding.Users, du)
		}
	}
}

func (setting *Setting) formatDingNotify(list []Setting, ding *Ding) {
	for _, v := range list {
		switch v.Key {
		case DingUrlKey:
			ding.Url = v.Value
		case DingTemplateKey:
			ding.Template = v.Value
		default:
			du := DingUser{}
			formatDingMobile(&du, v)
			ding.Users = append(ding.Users, du)
		}
	}
}

func (setting *Setting) Ding() (Ding, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", DingCode).Find(&list)
	ding := Ding{}
	if err != nil {
		return ding, err
	}

	setting.formatDing(list, &ding)

	return ding, err
}

func (setting *Setting) DingNotify() (Ding, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", DingCode).Find(&list)
	ding := Ding{}
	if err != nil {
		return ding, err
	}

	setting.formatDingNotify(list, &ding)

	return ding, err
}

// 更新ding配置
func (setting *Setting) UpdateDing(url, template string) error {
	setting.Value = url

	Db.Cols("value").Update(setting, Setting{Code: DingCode, Key: DingUrlKey})

	setting.Value = template
	Db.Cols("value").Update(setting, Setting{Code: DingCode, Key: DingTemplateKey})

	return nil
}

// 创建ding用户
func (setting *Setting) CreateDingUser(name string, mobile string) (int64, error) {
	value := getDingUserValue(name, mobile)
	setting.Code = DingCode
	setting.Key = DingUserKey
	setting.Value = value

	return Db.Insert(setting)
}

func (setting *Setting) IsUserNameExist(name string) bool {
	ns := new(Setting)
	count, _ := Db.Where("`code` = ? AND `key` = ? AND `value` LIKE ?", DingCode, DingUserKey, "%"+name+"%").Count(ns)

	return count > 0
}

func (setting *Setting) IsUserMobileExist(mobile string) bool {
	ns := new(Setting)
	count, _ := Db.Where("`code` = ? AND `key` = ? AND `value` LIKE ?", DingCode, DingUserKey, "%"+mobile+"%").Count(ns)

	return count > 0
}

// 删除ding用户
func (setting *Setting) RemoveDingUser(id int) (int64, error) {
	setting.Code = DingCode
	setting.Key = DingUserKey
	setting.Id = id
	return Db.Delete(setting)
}

// region slack配置

type Slack struct {
	Url      string    `json:"url"`
	Channels []Channel `json:"channels"`
	Template string    `json:"template"`
}

type Channel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (setting *Setting) Slack() (Slack, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", SlackCode).Find(&list)
	slack := Slack{}
	if err != nil {
		return slack, err
	}

	setting.formatSlack(list, &slack)

	return slack, err
}

func (setting *Setting) formatSlack(list []Setting, slack *Slack) {
	for _, v := range list {
		switch v.Key {
		case SlackUrlKey:
			slack.Url = v.Value
		case SlackTemplateKey:
			slack.Template = v.Value
		default:
			slack.Channels = append(slack.Channels, Channel{
				v.Id, v.Value,
			})
		}
	}
}

func (setting *Setting) UpdateSlack(url, template string) error {
	setting.Value = url

	Db.Cols("value").Update(setting, Setting{Code: SlackCode, Key: SlackUrlKey})

	setting.Value = template
	Db.Cols("value").Update(setting, Setting{Code: SlackCode, Key: SlackTemplateKey})

	return nil
}

// 创建slack渠道
func (setting *Setting) CreateChannel(channel string) (int64, error) {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Value = channel

	return Db.Insert(setting)
}

func (setting *Setting) IsChannelExist(channel string) bool {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Value = channel

	count, _ := Db.Count(setting)

	return count > 0
}

// 删除slack渠道
func (setting *Setting) RemoveChannel(id int) (int64, error) {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Id = id
	return Db.Delete(setting)
}

// endregion

type Mail struct {
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	User      string     `json:"user"`
	Password  string     `json:"password"`
	MailUsers []MailUser `json:"mail_users"`
	Template  string     `json:"template"`
}

type MailUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// region 邮件配置
func (setting *Setting) Mail() (Mail, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", MailCode).Find(&list)
	mail := Mail{MailUsers: make([]MailUser, 0)}
	if err != nil {
		return mail, err
	}

	setting.formatMail(list, &mail)

	return mail, err
}

func (setting *Setting) formatMail(list []Setting, mail *Mail) {
	mailUser := MailUser{}
	for _, v := range list {
		switch v.Key {
		case MailServerKey:
			json.Unmarshal([]byte(v.Value), mail)
		case MailUserKey:
			json.Unmarshal([]byte(v.Value), &mailUser)
			mailUser.Id = v.Id
			mail.MailUsers = append(mail.MailUsers, mailUser)
		case MailTemplateKey:
			mail.Template = v.Value
		}

	}
}

func (setting *Setting) UpdateMail(config, template string) error {
	setting.Value = config
	Db.Cols("value").Update(setting, Setting{Code: MailCode, Key: MailServerKey})

	setting.Value = template
	Db.Cols("value").Update(setting, Setting{Code: MailCode, Key: MailTemplateKey})

	return nil
}

func (setting *Setting) CreateMailUser(username, email string) (int64, error) {
	setting.Code = MailCode
	setting.Key = MailUserKey
	mailUser := MailUser{0, username, email}
	jsonByte, err := json.Marshal(mailUser)
	if err != nil {
		return 0, err
	}
	setting.Value = string(jsonByte)

	return Db.Insert(setting)
}

func (setting *Setting) RemoveMailUser(id int) (int64, error) {
	setting.Code = MailCode
	setting.Key = MailUserKey
	setting.Id = id
	return Db.Delete(setting)
}

type WebHook struct {
	Url      string `json:"url"`
	Template string `json:"template"`
}

func (setting *Setting) Webhook() (WebHook, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", WebhookCode).Find(&list)
	webHook := WebHook{}
	if err != nil {
		return webHook, err
	}

	setting.formatWebhook(list, &webHook)

	return webHook, err
}

func (setting *Setting) formatWebhook(list []Setting, webHook *WebHook) {
	for _, v := range list {
		switch v.Key {
		case WebhookUrlKey:
			webHook.Url = v.Value
		case WebhookTemplateKey:
			webHook.Template = v.Value
		}

	}
}

func (setting *Setting) UpdateWebHook(url, template string) error {
	setting.Value = url

	Db.Cols("value").Update(setting, Setting{Code: WebhookCode, Key: WebhookUrlKey})

	setting.Value = template
	Db.Cols("value").Update(setting, Setting{Code: WebhookCode, Key: WebhookTemplateKey})

	return nil
}

// endregion
