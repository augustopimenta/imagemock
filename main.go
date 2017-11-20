package main

import (
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"regexp"
	"os"
	"net/http"
	"github.com/lucasb-eyer/go-colorful"
	"math"
	"image"
	"bytes"
	"image/png"
	"errors"
	"time"
	"runtime/debug"
)

type cacheImage struct{
	image *bytes.Buffer
	lifeTime int64
}

const (
	maxImageSize = 5000
	minImageSize = 50
	cacheTimeInSeconds = 100
	cacheRemoveRoutineTimeInSeconds = 5
)

var cache = make(map[string]*cacheImage)
var creating = make(map[string]bool)

func main() {
	go clearCache()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		html(c, http.StatusOK, renderHelp())
	})

	r.GET("/:size", func(c *gin.Context) {
		width, height, err := extractSize(c.Param("size"))

		if err != nil {
			html(c, http.StatusInternalServerError, renderError(err.Error()))
			return
		}

		background := c.DefaultQuery("bg", "666666")
		fColor := c.DefaultQuery("c", "FFFFFF")
		text := c.DefaultQuery("t", fmt.Sprintf("%d x %d", width, height))

		keyMap := string(width) +";"+ string(height) +";"+ background +";"+ fColor +";"+ text

		var isLock bool = false
		for _, isLock = creating[keyMap]; isLock ; {
			time.Sleep(5 * time.Millisecond)
			_, isLock = creating[keyMap]
		}

		if image, ok := cache[keyMap]; ok {
			sendImage(c,image.image)
			cache[keyMap].lifeTime = time.Now().Add(cacheTimeInSeconds * time.Second).Unix()
			return
		}

		creating[keyMap] = true
		image, err := generateImage(width, height, background, fColor, text)
		delete(creating,keyMap)

		if err != nil {
			html(c, http.StatusInternalServerError, renderError(err.Error()))
			return
		}

		sendImage(c, image)
	})

	r.Run(getPort())
}

func sendImage(c *gin.Context, image *bytes.Buffer){
	c.Header("Pragma","public")
	c.Header("Cache-Control","max-age=86400")
	c.Header("Expires", time.Now().AddDate(60, 0, 0).Format(http.TimeFormat))

	c.Data(http.StatusOK, "image/png", image.Bytes())
}

func clearCache() {
	for {
		for k, v := range cache {
			if v.lifeTime < time.Now().Unix() {
				delete(cache, k)
				fmt.Println("removendo cache, index: ",k)
			}
		}
		time.Sleep(cacheRemoveRoutineTimeInSeconds * time.Second)
		debug.FreeOSMemory()
	}
}

func extractSize(size string) (int, int, error) {
	r := regexp.MustCompile(`^(\d+)([xX](\d+))?$`)
	matches := r.FindStringSubmatch(size)

	if len(matches) == 0 {
		return 0, 0, errors.New("Parâmetros inválidos")
	}

	width, _ := strconv.Atoi(matches[1])
	height, err := strconv.Atoi(matches[3])

	if err != nil {
		height = width
	}

	if width > maxImageSize || height > maxImageSize {
		return width, height, errors.New("A imagem deve ter no máximo 5000 x 5000")
	}

	if width < minImageSize || height < minImageSize {
		return width, height, errors.New("A imagem deve ter no mínimo 50 x 50")
	}

	return width, height, nil;
}

func generateImage(width, height int, background, fColor, text string) (*bytes.Buffer, error) {
	dc := gg.NewContext(width, height)
	dc.DrawRectangle(0,0, float64(width), float64(height))

	bgColor, err := colorful.Hex("#" + background)

	if err != nil {
		return nil, errors.New("Cor de fundo incorreta")
	}

	dc.SetColor(bgColor)
	dc.Fill()

	fontColor, err := colorful.Hex("#" + fColor)

	if err != nil {
		return nil, errors.New("Cor do texto incorreta")
	}

	points := math.Min(float64(width), float64(height)) / 6

	dc.SetColor(fontColor);
	if err := dc.LoadFontFace("Roboto-Medium.ttf", points); err != nil {
		return nil, errors.New("Ocorreu um erro para encontrar a fonte")
	}

	dc.DrawStringWrapped(text, float64(width/2), float64(height/2), 0.5, 0.5, float64(width - 10), 1.3, gg.AlignCenter)
	data, err := imageToBytes(dc.Image())

	if err != nil {
		return nil, errors.New("Ocorreu um erro para processar a imagem")
	}

	cache[string(width) +";"+ string(height) +";"+ background +";"+ fColor +";"+ text] = &cacheImage{
		image:data,lifeTime:time.Now().Add(cacheTimeInSeconds * time.Second).Unix(),
	}

	return data,nil
}

func imageToBytes(image image.Image) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, image)

	return buf, err;
}

func getPort() (string) {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	fmt.Printf("Aplicação iniciada na porta %s\n", port)
	return fmt.Sprintf(":%s", port)
}

func html(c *gin.Context, code int, text []byte) {
	c.Data(code, "text/html; charset=utf-8", text)
}
