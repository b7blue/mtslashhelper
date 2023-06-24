package download

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"

	conn "mtslashhelper/connector"
)

var urlxpath = `(//*[@class="authi"]//a)[1]`

func init() {
	os.MkdirAll("txt", 0664)
}

// 将该帖子转换为txt
func GetTXT(tid int) (filename string, err error) {
	client := conn.NewClient()
	// 获取楼主uid    (//*[@class="authi"])[1]
	url := fmt.Sprintf("http://www.mtslash.me/thread-%d-1-1.html", tid)
	urlnode, err := conn.GetNode(url, urlxpath, client)
	if err != nil {
		return "", err
	}

	spaceurl := htmlquery.SelectAttr(urlnode[0], "href")
	authoridreg := regexp.MustCompile(`\d+`)
	authorid := authoridreg.FindString(spaceurl)

	// 然后get到只看楼主的文章链接，然后解析第一页，然后获得pages
	authviewurl := fmt.Sprintf("http://www.mtslash.me/forum.php?mod=viewthread&tid=%d&page=1&authorid=%s", tid, authorid)
	pages, title, pageonenode, err := Parsefirstpage(authviewurl, client)

	if err != nil {
		return "", err
	}

	// 新建strings.builder
	builder := &strings.Builder{}

	// 循环pages次，解析页面
	for i := 1; i <= pages; i++ {
		if i == 1 {
			page2TXT(pageonenode, builder)
		} else {
			eachpageurl := fmt.Sprintf("http://www.mtslash.me/forum.php?mod=viewthread&tid=%d&page=%d&authorid=%s", tid, i, authorid)
			eachpage, err := conn.GetRoot(eachpageurl, client)
			if err != nil {
				return "", err
			}
			page2TXT(eachpage, builder)
		}
	}

	// 写入文件
	filename = "txt/" + title + ".txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	file.WriteString(builder.String())

	return title, nil
}

// 解析第一页
// 1 获取标题
// 2 获取总页数
// 3 获取该页的root node
func Parsefirstpage(url string, client *http.Client) (int, string, *html.Node, error) {
	pages := 1
	root, err := conn.GetRoot(url, client)
	if err != nil {
		return -1, "", nil, err
	}
	// 1
	utf8decoder := simplifiedchinese.GBK.NewDecoder()
	title, _ := utf8decoder.String(htmlquery.InnerText(htmlquery.FindOne(root, `//*[@id="thread_subject"]`)))
	special := []string{`\`, `/`, `:`, `*`, `?`, `"`, `<`, `>`, `|`}
	for i := range special {
		title = strings.ReplaceAll(title, special[i], "&")
	}
	// 2
	if pagenode := htmlquery.FindOne(root, `(//*[@class="pg"])[1]//input`); pagenode != nil {
		pages, _ = strconv.Atoi(htmlquery.SelectAttr(pagenode, "size"))
	}

	return pages, title, root, nil
}

// 获得该页的文章内容
func page2TXT(node *html.Node, builder *strings.Builder) {
	// enc := mahonia.NewDecoder("gbk")
	// utf8decoder := simplifiedchinese.GBK.NewDecoder()
	textpath := `//*[starts-with(@id,'postmessage') and not(div[@class="quote"])]`
	textnodes := htmlquery.Find(node, textpath)
	for _, tn := range textnodes {
		// s := enc.ConvertString((htmlquery.InnerText(tn)))
		// s, err := utf8decoder.String((htmlquery.InnerText(tn)))
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		builder.WriteString(htmlquery.InnerText(tn))
	}
	// 去除标签全部化为文本
	// html = strings.ReplaceAll(html, "<br>", "\n")
}
