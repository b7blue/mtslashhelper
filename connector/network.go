package connector

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// 有一个重新加载页面，要用那个set-cookie的ivGn_2132_lastrequest来替换
// html用gzip压缩过

// "http://www.mtslash.me/home.php?mod=space&do=favorite&view=me"

func GetHTML(href string, client *http.Client) (html io.ReadCloser, err error) {
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return nil, err
	}
	cookie := Getcookie()
	req.Header = map[string][]string{
		"Host":                      {"www.mtslash.me"},
		"Proxy-Connection":          {"keep-alive"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Cache-Control":             {"max-age=0"},
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36"},
		"Upgrade-Insecure-Requests": {"1"},
		// "Referer":                   {"http://www.mtslash.me/home.php?mod=space&do=favorite&view=me"},
		"Accept-Encoding": {"gzip", "deflate"},
		"Accept-Language": {"zh-CN", "zh;q=0.9"},
		"Cookie":          {cookie},
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("ohno!网络出现了问题！")
	}

	// 事情是这样的，可能会出现那个“页面重新加载”界面，然后必须根据那个setcookie来改变last request再重新请求
	if res.Header.Get("content-encoding") == "gzip" {
		// fmt.Println("一次成功")

	} else {
		// 从http响应中获取lastrequest
		rawlq := res.Header.Get("set-cookie")
		newlq := rawlq[:strings.IndexRune(rawlq, ';')+1]

		// 生成新cookie，同时更新本地cookie文件
		cookie = Setlastrequest(cookie, newlq)

		req.Header.Set("Cookie", cookie)
		res, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("ohno!网络出现了问题！")
		}
		// fmt.Println("两次成功")
	}

	if res.Header.Get("content-encoding") == "gzip" {
		// 假如遇到添加_dsign的反扒机制，获取js代码，执行js，获得新连接，然后再用新连接调用GetHTML
		// 根据响应头有没有set cookie ivGn_2132_viewid 来判断是否成功获取到了帖子页面！！！
		if len(res.Header.Values("Set-Cookie")) == 4 {
			jsreg := regexp.MustCompile(`<script type="text/javascript">(.+)</script>`)
			raw, _ := ioutil.ReadAll(gzip2bytes(res.Body))
			js := jsreg.FindAllStringSubmatch(string(raw), -1)[0][1]
			// fmt.Println(js)
			righturl, err := getrighturl(js)
			if err != nil {
				return nil, err
			}
			html, err = GetHTML(righturl, client)
			if err != nil {
				return nil, err
			}
		} else {
			html = gzip2bytes(res.Body)
		}

	} else {

		return nil, fmt.Errorf("获取页面失败，cookie出错，链接为：%s", href)
	}
	// b, err := ioutil.ReadAll(html)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	return html, nil
}

// 发来正式的html文件的时候，会gzip格式压缩，所以要将其解码
func gzip2bytes(zipb io.Reader) io.ReadCloser {
	gzipreader, err := gzip.NewReader(zipb)
	if err != nil {
		fmt.Println(err)
	}

	return gzipreader
}

func NewClient() *http.Client {
	return &http.Client{}
}

// 检测链接是否有问题
func Checknetwork(url string) {
	client := NewClient()
	html, err := GetHTML(url, client)
	if err != nil {
		fmt.Println("network err:", err)
	} else {
		defer html.Close()

		b, err := ioutil.ReadAll(html)
		if err != nil {
			fmt.Println(err)
		}
		ioutil.WriteFile("test.html", b, 0666)
	}
}
