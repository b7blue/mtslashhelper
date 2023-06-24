package sub

import (
	"fmt"

	conn "mtslashhelper/connector"
	"mtslashhelper/models"

	"golang.org/x/sync/errgroup"
)

// 根据文章链接和发出订阅的用户uid，在订阅列表中添加新文章
// 根据链接解析出：aid、标题、最近一次更新时间
func AddSub(uid int, url string) error {
	tid, err := conn.GetTid(url)
	if err != nil {
		return err
	}

	// 首先判断该文章是否已经被别人订阅过
	if !models.TidExist(tid) {
		// 首先请求该页面，然后解析获得文章编号、标题、链接、最后更新时间
		href := fmt.Sprintf("http://www.mtslash.me/thread-%d-1-1.html", tid)

		// 获得新文章的信息,标题和最后更新时间通过解析收到的html页面得到
		newa, err := GetNewArticle(tid, href)
		if err != nil {
			return err
		}

		// 插入Article表
		models.AddArticle(newa)
	}

	// 在sub数据库中添加订阅关系
	return models.AddSub(tid, uid)
}

// 0：检查所有受到订阅的文章是否更新，假如查到发生更新，
// 1:更新该文章的最后更新时间，这个可以几种所有列表
// 2：就给所有订阅了该文章的用户发消息（在该用户的更新信息表中新增一条）
func CheckUpdate() error {
	// 首先在Article中getall获得文章列表
	subList := models.GetAllArticle()

	// 然后一个个查找该文章的最近一次更新时间，看看是否晚于数据库中的记录
	// 一个线程查100个文章吧，反正也不赶时间。。。
	n := 100
	gr_nums := (len(subList) + n - 1) / n
	updateList := make(chan models.Article, 100) //update list
	c := conn.NewClient()

	// 生产者消费者模型：检查到发生更新——给订阅该文章的所有用户发消息
	var eg errgroup.Group
	for i := 0; i < gr_nums; i++ {

		// 闭包
		i := i
		eg.Go(func() error {
			for j := i * 100; j < (i+1)*100 && j < len(subList); j++ {
				root, err := conn.GetRoot(subList[j].Href, c)
				if err != nil {
					return err
				}
				newupdate := parseupdatetime(root)
				if newupdate.After(subList[j].Lastupdate) {
					subList[j].Lastupdate = newupdate
					updateList <- subList[j]
				}
			}
			return nil
		})
	}

	// updateList Closer
	go func() {
		err := eg.Wait()
		if err != nil {

		}
		close(updateList)
	}()

	// 假如有文更新，全部收集起来放在切片里，等下全部一起更新效率高
	// 并发地给所有订阅该文章的人的消息数据库中插入一条
	var eg1 errgroup.Group
	allUpdate := make([]models.Article, 0)
	for a := range updateList {
		allUpdate = append(allUpdate, a)

		// 闭包
		a := a
		eg1.Go(func() error {
			if err := updateReminder(a); err != nil {
				return err
			}
			return nil
		})
	}

	// 更新数据库
	models.SetUpdateTime(allUpdate)

	err := eg1.Wait()
	return err

}

// 给订阅了该文章的所有用户消息库中插入一条
// 因为同时订阅同一个文章的用户也不会太多，就循环插入好了
func updateReminder(a models.Article) error {
	newmesssage := models.Message{
		ArticleInfo: a,
	}
	users := models.GetSubListByTid(a.Tid)

	// 给订阅了该文章的所有用户消息库中插入一条
	for _, u := range users {
		if err := models.NewMsg2User(newmesssage, u); err != nil {
			return err
		}
	}

	return nil
}
