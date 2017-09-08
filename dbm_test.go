package dbm

import "testing"

func Test_connect(t *testing.T) {
	err := New("127.0.0.1", "root", "snyptfx16", "whisper").Port(3306).Set(
		SetTimeout("200"),
		SetCharset("utf-8"),
	).Add("test")
	if err != nil {
		t.Errorf("Test_connect error:%s", err)
	}
}
