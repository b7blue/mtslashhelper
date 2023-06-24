package models

import (
	"fmt"
	"sort"
)

type MessageList []Message

func (list MessageList) Len() int {
	return len(list)
}

func (list MessageList) Less(i, j int) bool {
	return list[i].ArticleInfo.Lastupdate.After(list[j].ArticleInfo.Lastupdate)
}

func (list MessageList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func NewMsg2User(m Message, uid int) error {
	db := GetConnection()
	db.Table(fmt.Sprintf("msg_%d", uid)).Save(&m)
	return db.Error
}

func GetMsgList(uid int) (m []Message) {
	m = make([]Message, 0)
	db.Table(fmt.Sprintf("msg_%d", uid)).Find(&m)
	// 从大到小排序
	sort.Sort(MessageList(m))
	return
}

func NewMsgList(uid int) error {
	db := GetConnection()
	db.Table(fmt.Sprintf("msg_%d", uid)).AutoMigrate(&Message{})
	return db.Error
}

func DelMsg(uid int, id []int) error {
	delsub := make([]*Message, len(id))
	for i := range delsub {
		delsub[i] = &Message{
			Id: id[i],
		}
	}
	re := db.Table(fmt.Sprintf("msg_%d", uid)).Delete(&delsub, "Id")

	return re.Error
}
