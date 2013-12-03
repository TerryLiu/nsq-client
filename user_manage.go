package main

import (
	"log"
	"strings"
)

type userManage struct {
	Users []User
}

type User struct {
	Id   string
	Nick string
	Pwd  string
}

var (
	UserMgr *userManage
)

func init() {
	log.Println("users.init...")
	UserMgr = &userManage{}

	UserMgr.Users = []User{
		{
			"29830845",
			"豆子/circle",
			"123456",
		},
		{
			"18926950",
			"柯男",
			"123456",
		},
		{
			"44070964",
			"nix",
			"123456",
		},
		{
			"286258251",
			"顽石",
			"123456",
		},
		{
			"368117356",
			"nic",
			"123456",
		},
		{
			"369036345",
			"鹄",
			"123456",
		},
	}
}

func (usrManage *userManage) IsUserValid(usrId, pwd string) (User, bool) {
	if strings.EqualFold(usrId, "") || strings.EqualFold(pwd, "") {
		return User{}, false
	}
	for _, usr := range usrManage.Users {
		if strings.EqualFold(usr.Id, usrId) && strings.EqualFold(usr.Pwd, pwd) {
			return usr, true
		}
	}
	return User{}, false
}
