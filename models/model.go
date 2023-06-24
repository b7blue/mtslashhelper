package models

import "time"

type User struct {
	// 自增长主键
	Uid   int    `gorm:"type:bigint auto_increment;primarykey;"`
	PW    string `gorm:"type:varchar(100);NOT NULL"` //哈希值
	Email string `gorm:"type:varchar(30) Unique;NOT NULL"`
}

func (User) TableName() string {
	return "User"
}

// 受到订阅的所有文章的信息
// 当某个用户新增订阅文章，且该文章没有被订阅过的时候，在此表新增条目
type Article struct {
	Tid        int       `gorm:"type:bigint;primarykey;NOT NULL"`
	Title      string    `gorm:"type:varchar(100);NOT NULL"`
	Href       string    `gorm:"type:varchar(100);NOT NULL"`
	Lastupdate time.Time `gorm:"type:datetime;NOT NULL"`
	// Subscriber string    `gorm:"type:int;NOT NULL"`
}

func (Article) TableName() string {
	return "Article"
}

// 文章订阅情况的表
// 当某个用户新增订阅文章的时候，就要在此表中增加用户和文章的订阅关系
type Subscription struct {
	Uid int `gorm:"primaryKey;NOT NULL"`
	Tid int `gorm:"primaryKey;NOT NULL"`
}

func (Subscription) TableName() string {
	return "Subscription"
}

// 注册-后端创建新用户的表项-新建一个uid_updates表
type Message struct {
	// 每一条消息的编号作为主键，自增长，作为消息删除的根据
	Id          int     `gorm:"primarykey;auto_increment"`
	ArticleInfo Article `gorm:"embedded"`
}

// 专门存放用户cookie的表
// type Cookie struct {
// 	Mtstempid string    `gorm:"NOT NULL"`
// 	Uid       int       `gorm:"primarykey;NOT NULL"`
// 	LastLogin time.Time `gorm:"NOT NULL"`
// }

// func (Cookie) TableName() string {
// 	return "Cookie"
// }
