package main

import (
	"github.com/lxn/walk"
)

type UsrItem struct {
	id   string
	nick string
}

type UsrModel struct {
	walk.ListModelBase
	items []UsrItem
}

func NewUsrModel() *UsrModel {
	m := &UsrModel{items: make([]UsrItem, len(UserMgr.Users))}

	for i, usr := range UserMgr.Users {
		m.items[i] = UsrItem{usr.Id, usr.Nick}
	}
	return m
}

func (m *UsrModel) ItemCount() int {
	return len(m.items)
}

func (m *UsrModel) Value(index int) interface{} {
	return m.items[index].nick
}
