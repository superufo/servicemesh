package admin

// XfcAdmin model  表  Xfc_Admin
type XfcAdmin struct {
	AdminId     int    `json:"admin_id"`
	UserName    string `json:"user_name"`
	DepId       int    `json:"dep_id"`
	PosId       int    `json:"pos_id"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Ecsalt      string `json:"ec_salt"`
	AddTime     int64  `json:"add_time"`
	Lastlogin   int64  `json:"last_login"`
	Lastip      string `json:"last_ip"`
	NavList     string `json:"nav_list"`
	LangType    string `json:"lang_type"`
	AgencyId    int    `json:"agency_id"`
	SuppliersId int    `json:"suppliers_id"`
	TodoList    string `json:"todo_list"`
	RoleId      int32  `json:"role_id"`
	IsLock      int    `json:"is_lock"`
	IsSales     int    `json:"is_sales"`
	IsGovsales  int    `json:"is_govsales"`
}

//`json:"user_info"`
type UserInfo struct {
	User XfcAdmin `json:"user_info"`
}
