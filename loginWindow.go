package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	//"log"
)

type LoginWindow struct {
	*walk.MainWindow
	userEdit *walk.LineEdit
	pwdEdit  *walk.LineEdit
	loginBtn *walk.PushButton
}

func NewLoginWindow() {
	walk.SetPanicOnError(true)
	myWindow, _ := walk.NewMainWindow()
	mw := &LoginWindow{MainWindow: myWindow}
	//mw.SetLayout(walk.NewVBoxLayout())
	mw.SetTitle("nsq client")

	userLabel, _ := walk.NewLabel(mw)
	userLabel.SetText("用户名:")
	userLabel.SetX(20)
	userLabel.SetY(10)
	userLabel.SetSize(walk.Size{40, 20})

	userEdit, _ := walk.NewLineEdit(mw)
	userEdit.SetReadOnly(false)
	userEdit.SetX(70)
	userEdit.SetY(10)
	userEdit.SetSize(walk.Size{200, 20})
	userEdit.KeyDown().Attach(mw.onKeyDown)
	mw.userEdit = userEdit

	pwdLabel, _ := walk.NewLabel(mw)
	pwdLabel.SetText("密码:")
	pwdLabel.SetX(20)
	pwdLabel.SetY(40)
	pwdLabel.SetSize(walk.Size{40, 20})

	pwdEdit, _ := walk.NewLineEdit(mw)
	pwdEdit.SetReadOnly(false)
	pwdEdit.SetX(70)
	pwdEdit.SetY(40)
	pwdEdit.SetSize(walk.Size{200, 20})
	pwdEdit.KeyDown().Attach(mw.onKeyDown)
	mw.pwdEdit = pwdEdit

	loginBtn, _ := walk.NewPushButton(mw)
	loginBtn.SetText("登陆")
	loginBtn.SetX(120)
	loginBtn.SetY(70)
	loginBtn.SetSize(walk.Size{60, 30})
	loginBtn.Clicked().Attach(mw.loginBtn_OnClick)
	mw.loginBtn = loginBtn

	mw.Show()
	mw.userEdit.SetFocus()
	mw.SetMinMaxSize(walk.Size{300, 150}, walk.Size{})
	mw.SetSize(walk.Size{300, 150})
	mw.Run()
}

func (mw *LoginWindow) loginBtn_OnClick() {
	nick := mw.userEdit.Text()
	pwd := mw.pwdEdit.Text()
	if !UserMgr.IsUserAndPwdValid(nick, pwd) {
		mw.onError("用户名或密码不正确！")
		return
	}
	go NewChatWindow()
	mw.MainWindow.SetVisible(false)
}

func (mw *LoginWindow) onKeyDown(key walk.Key) {

	switch key {
	case walk.KeyReturn:
		nick := mw.userEdit.Text()
		pwd := mw.pwdEdit.Text()
		if !UserMgr.IsUserAndPwdValid(nick, pwd) {
			mw.onError("用户名或密码不正确！")
			return
		}
		go NewChatWindow()
		mw.MainWindow.SetVisible(false)
	default:
	}
}

func (mw *LoginWindow) onError(msg string) {
	walk.MsgBox(mw, "错误", msg, walk.MsgBoxIconInformation)
}
