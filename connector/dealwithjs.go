package connector

import (
	"fmt"

	"github.com/robertkrimen/otto"
	// "github.com/dop251/goja"
)

func getrighturl(js string) (righturl string, err error) {
	head := `var url;var geturl=function(n){url = n.href;if (url == undefined){url = location;}};var location={href:"",assign:function(n){this.href=n;},replace:function(n){this.href=n;}};var window={href:""};`
	tail := `geturl(location);`

	js = head + js + tail

	vm := otto.New()
	_, err = vm.Run(js)
	if err != nil {
		return "", err
	}
	suburl, err := vm.Get("url")
	if err != nil {
		return "", err
	}
	// fmt.Println(suburl)
	return fmt.Sprintf("http://www.mtslash.me%v", suburl), nil
}
