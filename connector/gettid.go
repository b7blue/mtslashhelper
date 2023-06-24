package connector

import (
	"fmt"
	"regexp"
	"strconv"
)

func GetTid(url string) (int, error) {
	// 链接有两种形式
	// http://www.mtslash.me/thread-335605-1-1.html?_dsign=78f309f7
	// http://www.mtslash.me/forum.php?mod=viewthread&tid=335605&page=1&authorid=385059&_dsign=78f309f7
	reg1 := regexp.MustCompile(`thread-(\d+)-`)
	reg2 := regexp.MustCompile(`&tid=(\d+)`)
	domainreg := regexp.MustCompile(`(http://www.mtslash.me/).*`)
	if !domainreg.MatchString(url) {
		return -1, fmt.Errorf("解析链接的过程出错:不是来自mtslash的链接")
	}

	var re [][]string
	if re = reg1.FindAllStringSubmatch(url, -1); re == nil {
		re = reg2.FindAllStringSubmatch(url, -1)
	}
	if re == nil {
		return -1, fmt.Errorf("解析链接的过程出错:无效链接")
	}
	tid, err := strconv.Atoi(re[0][1])
	if err != nil {
		return -1, fmt.Errorf("解析链接的过程出错:%v", err)
	}
	return tid, nil
}
