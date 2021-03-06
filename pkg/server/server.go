package betterog

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/julianshen/betterog/pkg/page"
	"github.com/julianshen/text2img"

	ua "github.com/mileusna/useragent"
)

type BetterOG struct {
	Addr   string
	drawer text2img.Drawer
}

func NewServer(port string, params text2img.Params) (*BetterOG, error) {
	drawer, err := text2img.NewDrawer(params)

	if err != nil {
		return nil, err
	}

	return &BetterOG{
		Addr:   fmt.Sprintf(":%s", port),
		drawer: drawer,
	}, nil
}

func (bog *BetterOG) drawText(text string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if decoded, err := base64.RawURLEncoding.DecodeString(text); err == nil {
		text = string(decoded)
	} else {
		return nil, err
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

func isBot(c *gin.Context) bool {
	useragent := ua.Parse(c.Request.Header.Get("User-Agent"))

	return useragent.Bot
}

func cacheControl(c *gin.Context, ttl int) {
	c.Writer.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, s-maxage=%d", ttl, ttl))
	c.Writer.Header().Set("CDN-Cache-Control", fmt.Sprintf("%d", ttl))
	c.Writer.Header().Set("Cloudflare-CDN-Cache-Control", fmt.Sprintf("max-age=%d", ttl))
}

func (bog *BetterOG) Start() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/**/*")
	// for health check
	r.GET("/__ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/t/:text", func(c *gin.Context) {
		if !isBot(c) {
			c.AbortWithStatus(403)
			return
		}

		text := c.Param("text")

		if buf, err := bog.drawText(text); err == nil {
			c.Header("content-length", strconv.Itoa(len(buf.Bytes())))
			cacheControl(c, 86400)
			c.Data(200, "image/jpeg", buf.Bytes())
		} else {
			c.AbortWithError(500, err)
		}
	})

	r.GET("/c/:encodedurl", func(c *gin.Context) {
		if !isBot(c) {
			c.AbortWithStatus(403)
			return
		}

		encodedurl := c.Param("encodedurl")

		if buf, err := page.Capture(encodedurl); err == nil {
			c.Header("content-length", strconv.Itoa(len(buf)))
			cacheControl(c, 10800)
			c.Data(200, "image/jpeg", buf)
		} else {
			log.Println(err)
			c.AbortWithError(500, err)
		}
	})

	r.GET("/tt/:template/:encoded_title", func(c *gin.Context) {
		template := c.Param("template")
		title := ""
		encodedTitle := c.Param("encoded_title")

		if titleBytes, err := base64.RawURLEncoding.DecodeString(encodedTitle); err == nil {
			title = string(titleBytes)
			log.Println(title)
		} else {
			log.Printf("encoded :'%s'\n", encodedTitle)
			log.Println(err)
		}

		c.HTML(http.StatusOK, fmt.Sprintf("%s/index.tmpl", template), gin.H{
			"title": title,
		})
	})

	r.GET("/jlns1/:encoded_title", func(c *gin.Context) {
		if !isBot(c) {
			c.AbortWithStatus(403)
			return
		}

		encodedTitle := c.Param("encoded_title")
		url := fmt.Sprintf("http://localhost%s/tt/jlns1/%s", bog.Addr, encodedTitle)
		log.Println(url)

		if buf, err := page.Capture(base64.RawURLEncoding.EncodeToString([]byte(url))); err == nil {
			c.Header("content-length", strconv.Itoa(len(buf)))
			cacheControl(c, 10800)
			c.Data(200, "image/jpeg", buf)
		} else {
			log.Println(err)
			c.AbortWithError(500, err)
		}
	})

	r.StaticFS("/test", http.Dir("static/test"))
	r.StaticFS("/images", http.Dir("static/images"))

	r.Run(bog.Addr)
}
