package models

// 根据用户uid查找订阅的文章列表
func GetSubListByUid(uid int) (article []Article) {
	tids := make([]int, 0)
	db.Where(&Subscription{Uid: uid}).Select("Tid").Find(&tids)
	return GetArticlesbyTids(tids)
}

// 根据tid查找订阅列表
func GetSubListByTid(tid int) (uids []int) {
	db := GetConnection()
	db.Where(&Subscription{Tid: tid}).Select("Uid").Find(&uids)
	return uids
}

// 新增订阅
func AddSub(uid, tid int) error {
	db := GetConnection()
	// .Clauses(clause.Insert{Modifier: "IGNORE"})
	return db.Create(&Subscription{
		Tid: tid,
		Uid: uid,
	}).Error
}

// 取消订阅
func DelSub(uid int, tid []int) error {
	delsub := make([]*Subscription, len(tid))
	for i := range delsub {
		delsub[i] = &Subscription{
			Uid: uid,
			Tid: tid[i],
		}
	}
	re := db.Delete(&delsub)
	// 假如发现已经没有用户订阅这个文章了，就把它删除
	for i := range tid {
		if nil == GetSubListByTid(tid[i]) {
			DelArticle(tid[i])
		}
	}
	return re.Error
}
