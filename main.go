package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/antchfx/htmlquery"
)

func main() {
	// fmt.Println(db.AddSub(1, 2))
	// fmt.Println(db.AddSub(2, 2))
	// fmt.Println(db.AddSub(3, 2))
	// fmt.Println(db.AddSub(3, 2))

	// fmt.Println(db.AddSub(1, 1))
	// fmt.Println(db.AddSub(1, 3))

	// fmt.Println(db.DelSub(1, []int{1, 2}))
	// connector.Checknetwork("http://www.mtslash.me/thread-242685-1-1.html")

	// fmt.Println(db.IsNewUser("df8yewhfjk@qq.com"))
	// allA := models.GetAllArticle()
	// for _, a := range allA {
	// 	models.NewMsg2User(models.Message{ArticleInfo: a}, 1)
	// }
	b, err := os.ReadFile("C:/Users/70408/desktop/test.html")
	if err != nil {
		fmt.Println(err)
	}
	reader := bytes.NewReader(b)
	n, err := htmlquery.Parse(reader)
	s := htmlquery.InnerText(n)
	// dec := simplifiedchinese.GB18030.NewDecoder()
	// b, err = dec.Bytes(b)
	// s, err = dec.Bytes(s)
	// fmt.Println(string(b))
	file, err := os.OpenFile("C:/Users/70408/desktop/testgbk.txt", os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	file.WriteString(s)
}
