package page

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapture(t *testing.T) {
	buf, e := Capture(base64.URLEncoding.EncodeToString([]byte("https://www.google.com")))
	assert.NoError(t, e)
	assert.NotZero(t, len(buf))
}
