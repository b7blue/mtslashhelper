package sub

import (
	conn "mtslashhelper/connector"
	"mtslashhelper/models"
	"regexp"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// 这里要不要做错误处理？还是说可以保证这里拿到的一定是正确的页面？

// 获得登记新文章所需的两个信息
func GetNewArticle(tid int, href string) (newa models.Article, err error) {
	root, err := conn.GetRoot(href, conn.NewClient())
	if err != nil {
		return models.Article{}, err
	}
	return models.Article{
		Tid:        tid,
		Href:       href,
		Title:      parsetitle(root),
		Lastupdate: parseupdatetime(root),
	}, err
}

// 获得上次更新日期
func GetLastUpdate(href string) (lastupdate time.Time, err error) {
	root, err := conn.GetRoot(href, conn.NewClient())
	if err != nil {
		return time.Time{}, err
	}
	return parseupdatetime(root), err
}

func parseupdatetime(root *html.Node) time.Time {
	// 发帖日期和最后回复日期要格式化
	form := "2006-1-2 15:04"
	datereg := regexp.MustCompile(`[1-9][0-9]*-[1-9][0-2]*-[1-9][0-2]* [0-5][0-9]*:[0-5][0-9]*`)
	t, _ := time.Parse(form, datereg.FindString(htmlquery.InnerText(htmlquery.FindOne(root, `//*[@class="pstatus"]`))))

	return t
}

func parsetitle(root *html.Node) string {
	utf8decoder := simplifiedchinese.GBK.NewDecoder()
	title, _ := utf8decoder.String(htmlquery.InnerText(htmlquery.FindOne(root, `//*[@id="thread_subject"]`)))

	// 对标题做处理特殊字符替换处理，存数据库的话应该不用吧
	// special := []string{`\`, `/`, `:`, `*`, `?`, `"`, `<`, `>`, `|`}
	// for i := range special {
	// 	title = strings.ReplaceAll(title, special[i], "&")
	// }

	return title
}
