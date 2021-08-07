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
	buf, err := server.drawText("For testing")
	assert.NoError(t, err)

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(buf.Bytes())
	text, _ := client.Text()
	assert.Equal(t, "For testing", text)

	buf, err = server.drawText(base64.URLEncoding.EncodeToString([]byte("For testing")))
	assert.NoError(t, err)
	client.SetImageFromBytes(buf.Bytes())
	text, _ = client.Text()
	assert.Equal(t, "For testing", text)
}
