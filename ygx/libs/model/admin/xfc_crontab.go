package admin

// XfcAdmin model  表  Xfc_Admin
type XfcCrontab struct {
	CronId     int    `json:"cron_id"`
	CronName   string `json:"cron_name"`
	Type       int    `json:"type"`
	Rule       string `json:"rule"`
	Program    string `json:"program"`
	Parameters string `json:"parameters"`
	EntryID    int    `json:"entry_id"`
	RunError   string `json:"run_error"`
	Creater    string `json:"creater"`
	CreateTime string `json:"create_time"`
}

//`json:"user_info"`
type CronInfo struct {
	rule XfcCrontab `json:"crontab_info"`
}
