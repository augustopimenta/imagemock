package main

import (
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/lucasb-eyer/go-colorful"
	"runtime/debug"
	"image/png"
	"net/http"
	"strconv"
	"regexp"
	"errors"
	"fmt"
	"os"
	"math"
	"image"
	"bytes"
	"time"
)

const (
	maxImageSize                    = 5000
	minImageSize                    = 50
)

func main() {
	go clearCache()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/:size", func(c *gin.Context) {
		w, h, err := extractSize(c)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
			return
		}

		bg, fg, err := extractColors(c)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
			return
		}

		text := c.DefaultQuery("t", fmt.Sprintf("%d x %d", w, h))

		keyMap := makeKey(w, h, bg, fg, text)

		var isLock bool = false
		for _, isLock = creating[keyMap]; isLock; {
			time.Sleep(5 * time.Millisecond)
			_, isLock = creating[keyMap]
		}

		if img, ok := GetCache(keyMap); ok {
			sendImage(c, img)
			return
		}

		creating[keyMap] = true
		img, err := generateImage(w, h, text, bg, fg)
		delete(creating, keyMap)

		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
			return
		}

		sendImage(c, img)
	})

	r.Run(getPort())
}

func sendImage(c *gin.Context, image *bytes.Buffer) {
	c.Header("Pragma", "public")
	c.Header("Cache-Control", "max-age=86400")
	c.Header("Expires", time.Now().AddDate(60, 0, 0).Format(http.TimeFormat))

	c.Data(http.StatusOK, "image/png", image.Bytes())
}

func clearCache() {
	for {
		for k, v := range cache {
			if v.lifeTime < time.Now().Unix() {
				delete(cache, k)
				fmt.Println("Removendo cache, index: ", k)
			}
		}
		time.Sleep(cacheRemoveRoutineTimeInSeconds * time.Second)
		debug.FreeOSMemory()
	}
}

func extractSize(c *gin.Context) (int, int, error) {
	r := regexp.MustCompile(`^(\d+)([xX](\d+))?$`)
	matches := r.FindStringSubmatch(c.Param("size"))

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

func extractColors(c *gin.Context) (colorful.Color, colorful.Color, error) {
	bg, err := colorful.Hex("#" + c.DefaultQuery("bg", "666666"))
	if err != nil {
		return colorful.Color{}, colorful.Color{}, errors.New("Cor de fundo incorreta")
	}

	fg, err := colorful.Hex("#" + c.DefaultQuery("fg", "FFFFFF"))
	if err != nil {
		return colorful.Color{}, colorful.Color{}, errors.New("Cor do texto incorreta")
	}

	return bg, fg, nil
}

func generateImage(width, height int, text string, bg, fg colorful.Color) (*bytes.Buffer, error) {
	dc := gg.NewContext(width, height)
	dc.DrawRectangle(0, 0, float64(width), float64(height))

	dc.SetColor(bg)
	dc.Fill()

	points := math.Min(float64(width), float64(height)) / 6

	dc.SetColor(fg)
	if err := dc.LoadFontFace("Roboto-Medium.ttf", points); err != nil {
		return nil, errors.New("Ocorreu um erro para encontrar a fonte")
	}

	dc.DrawStringWrapped(text, float64(width/2), float64(height/2), 0.5, 0.5, float64(width-10), 1.3, gg.AlignCenter)
	data, err := imageToBytes(dc.Image())

	if err != nil {
		return nil, errors.New("Ocorreu um erro para processar a imagem")
	}

	PutCache(makeKey(width, height, bg, fg, text), data)

	return data, nil
}

func makeKey(w, h int, bg, fg colorful.Color, text string) string {
	return fmt.Sprintf("%d;%d;%s;%s;%s", w, h, bg.Hex(), fg.Hex(), text)
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
