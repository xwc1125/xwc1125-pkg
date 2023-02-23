package safecurl

import (
	"net/http"
	"testing"
)

// TestSafeClient 测试是否能创建安全的HTTP请求客户端
func TestSafeClient(t *testing.T) {

	urlValid := "https://www.xwc1125.com"
	urlInValid := "http://tsrc-event.xwc1125.cn@xwc1125.com/event_flag.html"

	safeClient := NewSafeClient()
	httpReqValid, err := http.NewRequest("GET", urlValid, nil)

	if err != nil {
		return
	}

	respValid, err := safeClient.Do(httpReqValid)

	if respValid == nil {
		t.Fatalf("TestSafeClient `urlValid` tested failed.")
	}

	httpReqInValid, err := http.NewRequest("GET", urlInValid, nil)

	if err != nil {
		return
	}

	respInValid, err := safeClient.Do(httpReqInValid)

	if respInValid != nil {
		t.Fatalf("TestSafeClient `respInValid` tested failed.")
	}

}

// TestSafeCurl 测试安全HTTP GET请求函数
func TestSafeCurl(t *testing.T) {
	hostValid := "https://www.xwc1125.com"
	hostInValid := "https://tst.xwc1125.com"

	checkResultValid := SafeCurl(hostValid)
	checkResultInValid := SafeCurl(hostInValid)

	if checkResultValid == nil {
		t.Fatalf("TestSafeCurl `checkResultValid` tested failed.")
	}

	if checkResultInValid != nil {
		t.Fatalf("TestSafeCurl `checkResultInValid` tested failed.")
	}
}
