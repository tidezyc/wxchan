package weixin

import (
	"os"
	"testing"

	"github.com/tuotoo/qrcode"
)

func TestQrcode(t *testing.T) {
	fi, err := os.Open("x.png")
	if err != nil {
		t.Fatal(err)
	}
	matrix, err := qrcode.Decode(fi)
	if err != nil {
		t.Fatal(err)
	}
	PrintQrcode(matrix)
}
