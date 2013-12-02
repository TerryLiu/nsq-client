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
	ids := [6]string{"29830845", "18926950", "44070964", "286258251", "368117356", "369036345"}
	nicks := [6]string{"豆子/circle", "柯男", "nix", "顽石", "nic", "鹄"}

	m := &UsrModel{items: make([]UsrItem, len(ids))}
	for i := 0; i < len(ids); i++ {
		m.items[i] = UsrItem{ids[i], nicks[i]}
	}
	return m
}

func (m *UsrModel) ItemCount() int {
	return len(m.items)
}

func (m *UsrModel) Value(index int) interface{} {
	return m.items[index].nick
}
