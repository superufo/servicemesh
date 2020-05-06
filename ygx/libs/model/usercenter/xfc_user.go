package usercenter

// XfcAdmin model  表  Xfc_Admin
type XfcUser struct {
	UserId        int32   `json:"user_id"`
	UserName      string  `json:"user_name"`
	Passwd        string  `json:"passwd"`
	Nick          string  `json:"nick"`
	Mobile        string  `json:"mobile"`
	Email         string  `json:"email"`
	Post          string  `json:"post"`
	TeamId        int32   `json:"team_id"`
	TeamName      string  `json:"post"`
	Introduce     string  `json:"introduce"`
	Department    string  `json:"department"`
	DepartmentId  int32   `json:"department_id"`
	Balance       float64 `json:"balance"`
	FreezeBalance float64 `json:"freeze_balance"`
	RealBalance   float64 `json:"real_balance"`
	RoleId        int32   `json:"role_id"`
}

//`json:"user_info"`
type UserInfo struct {
	User XfcUser `json:"user_info"`
}
