package main

import (
	"github.com/lxn/walk"
	"os"
	"strings"
)

type MsgItem struct {
	body string
}

type MsgModel struct {
	walk.ListModelBase
	items []MsgItem
}

func NewMsgModel() *MsgModel {
	env := os.Environ()

	m := &MsgModel{items: make([]MsgItem, len(env))}

	for i, e := range env {
		j := strings.Index(e, "=")
		if j == 0 {
			continue
		}

		name := e[0:j]
		m.items[i] = MsgItem{name}
	}

	return m
}

func (m *MsgModel) ItemCount() int {
	return len(m.items)
}

func (m *MsgModel) Value(index int) interface{} {
	return m.items[index].body
}
