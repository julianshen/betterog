package betterog

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianshen/text2img"
)

type BetterOG struct {
	Addr   string
	drawer text2img.Drawer
}

func NewServer(addr string, params text2img.Params) (*BetterOG, error) {
	drawer, err := text2img.NewDrawer(params)

	if err != nil {
		return nil, err
	}

	return &BetterOG{
		Addr:   addr,
		drawer: drawer,
	}, nil
}

func (bog *BetterOG) drawText(text string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if decoded, err := base64.URLEncoding.DecodeString(text); err == nil {
		text = string(decoded)
	}

	var err error
	var img image.Image
	if img, err = bog.drawer.Draw(text); err == nil {
		if err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err == nil {
			return buf, nil
		}
	}

	return nil, err
}

func (bog *BetterOG) Start() {
	r := gin.Default()

	// for health check
	r.GET("/__ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/t/:text", func(c *gin.Context) {
		log.Println(c.Request.Header)
		text := c.Param("text")

		if buf, err := bog.drawText(text); err == nil {
			c.Data(200, "image/jpeg", buf.Bytes())
		} else {
			c.AbortWithError(500, err)
		}
	})

	r.StaticFS("/test", http.Dir("static/"))

	r.Run(bog.Addr)
}
