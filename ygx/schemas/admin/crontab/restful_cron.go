package crontab

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/robfig/cron"

	"github.com/go-chassis/ygx/libs/common"
	"github.com/go-chassis/ygx/libs/i18n"
	"github.com/go-chassis/ygx/libs/model/admin"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	rf "github.com/go-chassis/go-chassis/server/restful"
	AdminModel "github.com/go-chassis/ygx/libs/model/admin"

	"github.com/go-mesh/openlogging"
)

type RestFulCrontab struct {
}

var crontab *cron.Cron

var once sync.Once

var cmd *exec.Cmd

// Root Sayhello is a method used to reply user with hello
func (r *RestFulCrontab) Root(b *rf.Context) {
	b.Write([]byte(fmt.Sprintf("x-forwarded-host %s", b.ReadRequest().Host)))
}

//
// call   127.0.0.1:15999/addrule  window 调用
//{"cron_name":"test06","type":"0","rule":"*/5 * * * *",
//"program":"C:\\Windows\\system32\\cmd.exe",
//"parameters":" /C, start, D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\admin\\crontab\\rest-server\\sh\\test\\test.bat"}
func (r *RestFulCrontab) AddRule(b *rf.Context) {
	crontabData := AdminModel.XfcCrontab{}
	reponse := common.Response{
		Status:  int(i18n.SUCCESS),
		Data:    crontabData,
		Message: i18n.LoadMessage(i18n.SUCCESS),
	}
	// 更新用户信息
	b.ReadEntity(&crontabData)

	if common.DB.NewRecord(crontabData) == true {
		common.DB.Create(&crontabData)

		//openlogging.GetLogger().Info("1111111111111111111111111111111111111111")
		crontab.AddFunc(crontabData.Rule, func() {
			//openlogging.GetLogger().Info(fmt.Sprint("crontabData.Rule: %s", crontabData.Rule))
			r.RunCmd(crontabData.Program, crontabData.Parameters)
		})
		reponse.Data = crontabData
	} else {
		reponse.Status = int(i18n.FAILURE)
		reponse.Message = i18n.LoadMessage(i18n.FAILURE)
	}

	b.WriteJSON(reponse, "application/json", "")
	return
}

// 必须传递 cron_name 或者 cron_id  或者 entry_id
func (r *RestFulCrontab) DelRule(b *rf.Context) {
	crontabData := AdminModel.XfcCrontab{}
	reponse := common.Response{
		Status:  int(i18n.SUCCESS),
		Data:    crontabData,
		Message: i18n.LoadMessage(i18n.SUCCESS),
	}
	// 更新用户信息
	b.ReadEntity(&crontabData)

	if string(crontabData.EntryID) == "" || string(crontabData.CronId) == "" || crontabData.CronName == "" {
		reponse.Status = int(i18n.PARAMETEREOOR)
		reponse.Message = i18n.LoadMessage(i18n.PARAMETEREOOR)
		b.WriteJSON(reponse, "application/json", "")
	}
	// 数据库删除
	common.DB.Where("cron_id = ? or entry_id=? or cron_name=?", crontabData.CronId, crontabData.EntryID, crontabData.CronName).Delete(crontabData)
	// 删除  crontab
	//entryId := cron.EntryID{crontabData.Entry}
	crontab.Remove(*(*cron.EntryID)(unsafe.Pointer(&crontabData.EntryID)))

	b.WriteJSON(reponse, "application/json", "")
}

// 获取所有 定时任务的列表
func (r *RestFulCrontab) GetRuleList(b *rf.Context) {
	var (
		cronList      []admin.XfcCrontab
		count         int
		ruleListQuest common.SearchListQuest
	)

	ruleListQuest = common.SearchListQuest{
		Page:     1,
		PageSize: 10,
		Data:     crontab,
	}
	b.ReadEntity(&ruleListQuest)

	offset := ruleListQuest.PageSize * ruleListQuest.Page
	common.DB.Limit(ruleListQuest.PageSize).Offset(offset).Find(&cronList)
	common.DB.Find(&cronList).Count(&count)

	listResponse := common.SearchListResponse{
		count,
		cronList,
	}

	reponse := common.Response{
		Status:  int(i18n.SUCCESS),
		Data:    listResponse,
		Message: i18n.LoadMessage(i18n.SUCCESS),
	}

	b.WriteJSON(reponse, "application/json", "")
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestFulCrontab) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/", ResourceFunc: r.Root,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/getrulelist",
			ResourceFunc: r.GetRuleList,
			Metadata: map[string]interface{}{
				"tags": []string{"rule", "name", "programme", "parameters"},
			},
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/addrule",
			ResourceFunc: r.AddRule,
			Metadata: map[string]interface{}{
				"tags": []string{"rule", "name", "programme", "parameters"},
			},
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/delrule",
			ResourceFunc: r.DelRule,
			Metadata: map[string]interface{}{
				"tags": []string{"EntryID"},
			},
			Returns: []*rf.Returns{{Code: 200}}},
	}
}

func (r *RestFulCrontab) RunCron() {
	var cronList []admin.XfcCrontab
	common.DB.Find(&cronList)

	if len(cronList) > 0 {
		for _, cronEx := range cronList {
			cid, err := crontab.AddFunc(cronEx.Rule, func() {
				r.RunCmd(cronEx.Program, cronEx.Parameters)
			})

			openlogging.GetLogger().Error("EntryID:", cid)
			//更新表中运行结果
			if int(cronEx.EntryID) != int(cid) {
				cronEx.RunError = fmt.Sprint("eroor:%s", err)
				cronEx.EntryID = int(cid)
				common.DB.Save(&cronEx)
			}
		}
	}
	crontab.Start()
}

func (r *RestFulCrontab) RunCmd(program string, parameters string) {
	parametersSets := strings.Split(parameters, ",")

	openlogging.GetLogger().Info(fmt.Sprint("program: %+s", program))
	openlogging.GetLogger().Info(fmt.Sprint("parametersSets: %+v  ", parametersSets))

	cmd = exec.Command(program, parametersSets...)
	if runtime.GOOS == "windows" {
		//cmd = exec.Command(program, "/c", "start ", parameters)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	err := cmd.Run()
	if err != nil {
		openlogging.GetLogger().Error("run parameters failure:%s", err)
	} else {
		openlogging.GetLogger().Info(fmt.Sprint("run parameters success:. %s %s ", program, parameters))
	}
}

func init() {
	//国际化语言设置
	i18n.SetPlatLang("admin", "cn")

	//只初始化一次
	once.Do(func() {
		crontab = cron.New()
		cronRest := new(RestFulCrontab)
		cronRest.RunCron()
	})
}
