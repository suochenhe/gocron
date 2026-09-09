package models

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"xorm.io/xorm"
)

type TaskProtocol int8

const (
	TaskHTTP TaskProtocol = iota + 1 // HTTP协议
	TaskRPC                          // RPC方式执行命令
)

type TaskLevel int8

const (
	TaskLevelParent TaskLevel = 1 // 父任务
	TaskLevelChild  TaskLevel = 2 // 子任务(依赖任务)
)

type TaskDependencyStatus int8

const (
	TaskDependencyStatusStrong TaskDependencyStatus = 1 // 强依赖
	TaskDependencyStatusWeak   TaskDependencyStatus = 2 // 弱依赖
)

type TaskHTTPMethod int8

const (
	TaskHTTPMethodGet  TaskHTTPMethod = 1
	TaskHttpMethodPost TaskHTTPMethod = 2
)

// 任务
type Task struct {
	Id               int                  `json:"id" xorm:"int pk autoincr"`
	Name             string               `json:"name" xorm:"varchar(128) notnull"`                           // 任务名称
	Level            TaskLevel            `json:"level" xorm:"tinyint notnull index default 1"`               // 任务等级 1: 主任务 2: 依赖任务
	DependencyTaskId string               `json:"dependency_task_id" xorm:"varchar(64) notnull default ''"`   // 依赖任务ID,多个ID逗号分隔
	DependencyStatus TaskDependencyStatus `json:"dependency_status" xorm:"tinyint notnull default 1"`         // 依赖关系 1:强依赖 主任务执行成功, 依赖任务才会被执行 2:弱依赖
	Spec             string               `json:"spec" xorm:"varchar(64) notnull"`                            // crontab
	Protocol         TaskProtocol         `json:"protocol" xorm:"tinyint notnull index"`                      // 协议 1:http 2:系统命令
	Command          string               `json:"command" xorm:"varchar(256) notnull"`                        // URL地址或shell命令
	HttpMethod       TaskHTTPMethod       `json:"http_method" xorm:"tinyint notnull default 1"`               // http请求方法
	Timeout          int                  `json:"timeout" xorm:"mediumint notnull default 0"`                 // 任务执行超时时间(单位秒),0不限制
	Multi            int8                 `json:"multi" xorm:"tinyint notnull default 1"`                     // 是否允许多实例运行
	RetryTimes       int8                 `json:"retry_times" xorm:"tinyint notnull default 0"`               // 重试次数
	RetryInterval    int16                `json:"retry_interval" xorm:"smallint notnull default 0"`           // 重试间隔时间
	NotifyStatus     int8                 `json:"notify_status" xorm:"tinyint notnull default 1"`             // 任务执行结束是否通知 0: 不通知 1: 失败通知 2: 执行结束通知 3: 任务执行结果关键字匹配通知
	NotifyType       int8                 `json:"notify_type" xorm:"tinyint notnull default 0"`               // 通知类型 1: 邮件 2: slack 3: webhook
	NotifyReceiverId string               `json:"notify_receiver_id" xorm:"varchar(256) notnull default '' "` // 通知接受者ID, setting表主键ID，多个ID逗号分隔
	NotifyKeyword    string               `json:"notify_keyword" xorm:"varchar(128) notnull default '' "`
	Tag              string               `json:"tag" xorm:"varchar(32) notnull default ''"`
	Remark           string               `json:"remark" xorm:"varchar(100) notnull default ''"` // 备注
	Status           Status               `json:"status" xorm:"tinyint notnull index default 0"` // 状态 1:正常 0:停止
	Created          time.Time            `json:"created" xorm:"datetime notnull created"`       // 创建时间
	Deleted          time.Time            `json:"deleted" xorm:"datetime deleted"`               // 删除时间
	BaseModel        `json:"-" xorm:"-"`
	Hosts            []TaskHostDetail `json:"hosts" xorm:"-"`
	NotifyReceivers  []string         `json:"notify_receivers" xorm:"-"`
	NextRunTime      time.Time        `json:"next_run_time" xorm:"-"`
}

func taskHostTableName() []string {
	return []string{TablePrefix + "task_host", "th"}
}

// 新增
func (task *Task) Create() (insertId int, err error) {
	_, err = Db.Insert(task)
	if err == nil {
		insertId = task.Id
	}

	return
}

func (task *Task) UpdateBean(id int) (int64, error) {
	return Db.ID(id).
		Cols(`name,spec,protocol,command,timeout,multi,
			retry_times,retry_interval,remark,notify_status,
			notify_type,notify_receiver_id, dependency_task_id, dependency_status, tag,http_method, notify_keyword`).
		Update(task)
}

// 更新
func (task *Task) Update(id int, data CommonMap) (int64, error) {
	return Db.Table(task).ID(id).Update(data)
}

// 删除
func (task *Task) Delete(id int) (int64, error) {
	return Db.ID(id).Delete(task)
}

// 禁用
func (task *Task) Disable(id int) (int64, error) {
	return task.Update(id, CommonMap{"status": Disabled})
}

// 激活
func (task *Task) Enable(id int) (int64, error) {
	return task.Update(id, CommonMap{"status": Enabled})
}

// 获取所有激活任务
func (task *Task) ActiveList(page, pageSize int) ([]Task, error) {
	params := CommonMap{"Page": page, "PageSize": pageSize}
	task.parsePageAndPageSize(params)
	list := make([]Task, 0)
	err := Db.Where("status = ? AND level = ?", Enabled, TaskLevelParent).Limit(task.PageSize, task.pageLimitOffset()).
		Find(&list)

	if err != nil {
		return list, err
	}

	return task.setHostsForTasks(list)
}

// 获取某个主机下的所有激活任务
func (task *Task) ActiveListByHostId(hostId int16) ([]Task, error) {
	taskHostModel := new(TaskHost)
	taskIds, err := taskHostModel.GetTaskIdsByHostId(hostId)
	if err != nil {
		return nil, err
	}
	if len(taskIds) == 0 {
		return nil, nil
	}
	list := make([]Task, 0)
	err = Db.Where("status = ?  AND level = ?", Enabled, TaskLevelParent).
		In("id", taskIds...).
		Find(&list)
	if err != nil {
		return list, err
	}

	return task.setHostsForTasks(list)
}

func (task *Task) setHostsForTasks(tasks []Task) ([]Task, error) {
	taskHostModel := new(TaskHost)
	var err error
	for i, value := range tasks {
		taskHostDetails, err := taskHostModel.GetHostIdsByTaskId(value.Id)
		if err != nil {
			return nil, err
		}
		tasks[i].Hosts = taskHostDetails
	}

	return tasks, err
}

// 判断任务名称是否存在
func (task *Task) NameExist(name string, id int) (bool, error) {
	if id > 0 {
		count, err := Db.Where("name = ? AND status = ? AND id != ?", name, Enabled, id).Count(task)
		return count > 0, err
	}
	count, err := Db.Where("name = ? AND status = ?", name, Enabled).Count(task)

	return count > 0, err
}

func (task *Task) GetStatus(id int) (Status, error) {
	exist, err := Db.ID(id).Get(task)
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, errors.New("not exist")
	}

	return task.Status, nil
}

func (task *Task) Detail(id int) (Task, error) {
	t := Task{}
	_, err := Db.Where("id=?", id).Get(&t)

	if err != nil {
		return t, err
	}

	taskHostModel := new(TaskHost)
	t.Hosts, err = taskHostModel.GetHostIdsByTaskId(id)

	return t, err
}

func (task *Task) List(params CommonMap) ([]Task, error) {
	task.parsePageAndPageSize(params)
	list := make([]Task, 0)
	session := Db.Alias("t").Join("LEFT", taskHostTableName(), "t.id = th.task_id")
	task.parseWhere(session, params)
	err := session.GroupBy("t.id").Desc("t.id").Cols("t.*").Limit(task.PageSize, task.pageLimitOffset()).Find(&list)

	if err != nil {
		return nil, err
	}

	return task.setTaskRelations(list)
}

func (task *Task) setTaskRelations(tasks []Task) ([]Task, error) {
	if tasks, err := task.setHostsForTasks(tasks); err != nil {
		return nil, err
	} else {
		return task.setNotifyReceivers(tasks)
	}
}

type NotifyReceiverOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (task *Task) NotifyReceiverOptions() ([]NotifyReceiverOption, error) {
	settings := make([]Setting, 0)
	if err := Db.Find(&settings); err != nil {
		return nil, err
	}
	options := []NotifyReceiverOption{
		{Value: "不通知", Label: "不通知"},
		{Value: "WebHook", Label: "WebHook"},
		{Value: "所有人", Label: "所有人"},
		{Value: "无", Label: "无"},
	}
	seen := map[string]bool{"不通知": true, "WebHook": true, "所有人": true, "无": true}
	for _, setting := range settings {
		label := notifyReceiverName(setting)
		if label != "" && !seen[label] {
			options = append(options, NotifyReceiverOption{Value: label, Label: label})
			seen[label] = true
		}
	}
	return options, nil
}

func notifyReceiverName(setting Setting) string {
	switch {
	case setting.Code == DingCode && setting.Key == DingUserKey:
		return getDingUserName(setting.Value)
	case setting.Code == MailCode && setting.Key == MailUserKey:
		var user MailUser
		if json.Unmarshal([]byte(setting.Value), &user) == nil {
			return user.Username
		}
	case setting.Code == SlackCode && setting.Key == SlackChannelKey:
		return setting.Value
	}
	return ""
}

func (task *Task) setNotifyReceivers(tasks []Task) ([]Task, error) {
	settings := make([]Setting, 0)
	if err := Db.Find(&settings); err != nil {
		return nil, err
	}
	for i := range tasks {
		tasks[i].NotifyReceivers = formatNotifyReceivers(tasks[i], settings)
	}
	return tasks, nil
}

func formatNotifyReceivers(task Task, settings []Setting) []string {
	if task.NotifyStatus == 0 {
		return []string{"不通知"}
	}
	if task.NotifyType == 3 {
		return []string{"WebHook"}
	}
	names := make([]string, 0)
	for _, id := range strings.Split(task.NotifyReceiverId, ",") {
		if id == "" {
			continue
		}
		if task.NotifyType == 4 {
			if id == DingNotifyNone {
				names = append(names, "无")
				continue
			}
			if id == DingNotifyAll {
				names = append(names, "所有人")
				continue
			}
		}
		for _, setting := range settings {
			if strconv.Itoa(setting.Id) != id {
				continue
			}
			if task.NotifyType == 4 && setting.Code == DingCode && setting.Key == DingUserKey {
				names = append(names, getDingUserName(setting.Value))
			} else if task.NotifyType == 1 && setting.Code == MailCode && setting.Key == MailUserKey {
				var user MailUser
				if json.Unmarshal([]byte(setting.Value), &user) == nil {
					names = append(names, user.Username)
				}
			} else if task.NotifyType == 2 && setting.Code == SlackCode && setting.Key == SlackChannelKey {
				names = append(names, setting.Value)
			}
		}
	}
	return names
}

// 获取依赖任务列表
func (task *Task) GetDependencyTaskList(ids string) ([]Task, error) {
	list := make([]Task, 0)
	if ids == "" {
		return list, nil
	}
	idList := strings.Split(ids, ",")
	taskIds := make([]interface{}, len(idList))
	for i, v := range idList {
		taskIds[i] = v
	}
	fields := "t.*"
	err := Db.Alias("t").
		Where("t.level = ?", TaskLevelChild).
		In("t.id", taskIds).
		Cols(fields).
		Find(&list)

	if err != nil {
		return list, err
	}

	return task.setHostsForTasks(list)
}

func (task *Task) Total(params CommonMap) (int64, error) {
	session := Db.Alias("t").Join("LEFT", taskHostTableName(), "t.id = th.task_id")
	task.parseWhere(session, params)
	list := make([]Task, 0)

	err := session.GroupBy("t.id").Find(&list)

	return int64(len(list)), err
}

func (task *Task) AllTags() ([]string, error) {
	tags := make([]string, 0)
	err := Db.Table(task).Distinct("tag").Where("tag != ''").Cols("tag").Find(&tags)
	return tags, err
}

// 解析where
func (task *Task) parseWhere(session *xorm.Session, params CommonMap) {
	if len(params) == 0 {
		return
	}
	id, ok := params["Id"]
	if ok && id.(int) > 0 {
		session.And("t.id = ?", id)
	}
	hostId, ok := params["HostId"]
	if ok && hostId.(int) > 0 {
		session.And("th.host_id = ?", hostId)
	}
	name, ok := params["Name"]
	if ok && name.(string) != "" {
		session.And("t.name LIKE ?", "%"+name.(string)+"%")
	}
	command, ok := params["Command"]
	if ok && command.(string) != "" {
		session.And("t.command LIKE ?", "%"+command.(string)+"%")
	}
	protocol, ok := params["Protocol"]
	if ok && protocol.(int) > 0 {
		session.And("protocol = ?", protocol)
	}
	status, ok := params["Status"]
	if ok && status.(int) > -1 {
		session.And("status = ?", status)
	}

	tag, ok := params["Tag"]
	if ok && tag.(string) != "" {
		session.And("tag = ? ", tag)
	}
	notifyUser, ok := params["NotifyUser"]
	if ok && notifyUser.(string) != "" {
		task.applyNotifyUserFilter(session, notifyUser.(string))
	}
}

func (task *Task) applyNotifyUserFilter(session *xorm.Session, keyword string) {
	if keyword == "不通知" {
		session.And("t.notify_status = ?", 0)
		return
	}
	if keyword == "WebHook" {
		session.And("t.notify_type = ?", 3)
		return
	}
	settings := make([]Setting, 0)
	if err := Db.Find(&settings); err != nil {
		return
	}
	ids := make([]interface{}, 0)
	for _, setting := range settings {
		if notifyReceiverName(setting) == keyword || strings.Contains(strings.ToLower(notifyReceiverName(setting)), strings.ToLower(keyword)) {
			ids = append(ids, setting.Id)
		}
	}
	if keyword == "所有人" || keyword == "无" {
		session.And("t.notify_receiver_id LIKE ?", "%"+map[string]string{"所有人": DingNotifyAll, "无": DingNotifyNone}[keyword]+"%")
		return
	}
	if len(ids) == 0 {
		session.And("1 = 0")
		return
	}
	conditions := make([]string, 0, len(ids))
	args := make([]interface{}, 0, len(ids)*4)
	for _, id := range ids {
		idValue := strconv.Itoa(id.(int))
		conditions = append(conditions, "t.notify_receiver_id = ? OR t.notify_receiver_id LIKE ? OR t.notify_receiver_id LIKE ? OR t.notify_receiver_id LIKE ?")
		args = append(args, idValue, idValue+",%", "%,"+idValue+",%", "%,"+idValue)
	}
	session.And("("+strings.Join(conditions, " OR ")+")", args...)
}
