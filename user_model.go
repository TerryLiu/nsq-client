package main

import (
	"github.com/lxn/walk"
)

type UsrModel struct {
	walk.ListModelBase
	items []User
}

func NewUsrModel() *UsrModel {
	m := &UsrModel{items: make([]User, len(UserMgr.Users))}

	for i, usr := range UserMgr.Users {
		m.items[i] = usr
	}
	return m
}

func (m *UsrModel) ItemCount() int {
	return len(m.items)
}

func (m *UsrModel) Value(index int) interface{} {
	return m.items[index].Nick
}
