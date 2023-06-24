package connector

import (
	"strings"
	"sync"
)

var COOKIE string = "ivGn_2132_lastrequest=c0c70ILCo2mbzB8yIfVnuwlBUmnYXTLWJPOSAUwBZDK%2BXXpotYFo;ivGn_2132_smile=1D1; ivGn_2132_saltkey=uUtZM7Yf; ivGn_2132_lastvisit=1661751641; ivGn_2132_seccode=333368.3fef72c7ab764e8b86; ivGn_2132_ulastactivity=8d39sKMOJqxvosIkqbie0UWgCqB3JV82iiR6R4CMxVBCYSZ%2FMA7o; ivGn_2132_auth=87d1Qw%2BRwGGxfBkDHDBYRpI1m3t2JbIKOuRL94lj6xi8pWzYhTt5C%2BUfYHHhthU9FVX%2F1vmCcrrB9Br3imAaz0auCDU; ivGn_2132_lastcheckfeed=951446%7C1662912506; ivGn_2132_visitedfid=19; ivGn_2132_viewid=tid_242685; ivGn_2132_sid=pX2xu2; ivGn_2132_lip=223.72.73.110%2C1662913398; ivGn_2132_sendmail=1; ivGn_2132_noticeTitle=1;ivGn_2132_lastact=1662913837%09forum.php%09viewthread; ivGn_2132_st_p=951446%7C1662913837%7Cc03d9cbc9b8bbefa534e476fe00682cb"
var rwlock sync.RWMutex

// 获取本地存储的cookie
func Getcookie() string {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return COOKIE
}

// 根据http响应中的lastrequest更新cookie
func Setlastrequest(oldcookie, newrq string) string {
	newcookie := newrq + oldcookie[strings.IndexRune(oldcookie, ';')+1:]
	rwlock.Lock()
	defer rwlock.Unlock()
	COOKIE = newcookie
	return newcookie
}
