package connector

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

// 假如要对该页面进一步解析，用这个
func GetRoot(url string, client *http.Client) (*html.Node, error) {
	html, err := GetHTML(url, client)

	if err != nil {
		return nil, err
	} else {
		defer html.Close()
		root, err := htmlquery.Parse(html)
		if err != nil {
			return nil, err
		}
		if root == nil {
			return nil, errors.New("html == nil")
		}
		return root, nil
	}
}

// 假如直接获得需要的节点就行了，用这个
func GetNode(url, xp string, client *http.Client) (nodes []*html.Node, err error) {
	html, err := GetHTML(url, client)

	if err != nil {
		return nil, err
	} else {
		top, err := htmlquery.Parse(html)
		if err != nil {
			return nil, err
		}
		if top == nil {
			return nil, errors.New("html == nil")
		}

		// 不知道为啥，但是一定要用转换成xpath.Expr能连续用多次
		// ioutil.WriteFile("test.txt", []byte(htmlquery.InnerText(top)), 0666)
		selector := xpath.MustCompile(xp)
		if nodes = htmlquery.QuerySelectorAll(top, selector); nodes == nil {
			return nil, fmt.Errorf("cannot find node by xpath")
		}
		html.Close()
		return nodes, nil
	}
}
