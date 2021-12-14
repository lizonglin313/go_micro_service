package model

type UserDetails struct {
	UserId      int64    // 用户ID
	Username    string   // 用户名
	Password    string   // 用户密码
	Authorities []string // 用户具备的权限
}
