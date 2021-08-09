package betterog

import (
	"encoding/base64"
	"testing"

	"github.com/julianshen/text2img"
	"github.com/otiai10/gosseract"
	"github.com/stretchr/testify/assert"
)

func TestDrawText(t *testing.T) {
	server, err := NewServer(":8888", text2img.Params{
		FontPath: "../../fonts/SourceHanSansTC-VF.ttf",
	})

	assert.NoError(t, err)

	client := gosseract.NewClient()
	defer client.Close()

	buf, err := server.drawText(base64.URLEncoding.EncodeToString([]byte("For testing")))
	assert.NoError(t, err)
	client.SetImageFromBytes(buf.Bytes())
	text, _ := client.Text()
	assert.Equal(t, "For testing", text)
}

func TestEncoding(t *testing.T) {
	code := base64.RawURLEncoding.EncodeToString([]byte("[筆記]利用axios-mock-adapter為axios提供測試用的假資料"))
	assert.Equal(t, "W-ethuiomF3liKnnlKhheGlvcy1tb2NrLWFkYXB0ZXLngrpheGlvc-aPkOS-m-a4rOippueUqOeahOWBh-izh-aWmQ", code)

	_, err := base64.RawURLEncoding.DecodeString("W-ethuiomF3liKnnlKhheGlvcy1tb2NrLWFkYXB0ZXLngrpheGlvc-aPkOS-m-a4rOippueUqOeahOWBh-izh-aWmQ")
	assert.NoError(t, err)
}
