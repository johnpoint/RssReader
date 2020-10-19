package apis

import (
	"log"
	"net/http"
	"os"
	"rssreader/src/model"

	echo2 "github.com/labstack/echo/v4"
)

const VERSION = "0.20"

func Accessible(c echo2.Context) error {
	return c.HTML(http.StatusOK, "<h1>RssReader api</h1>(´・ω・`) 运行正常<br><hr>Ver: "+VERSION)
}

func Post(c echo2.Context) error {
	file, err := os.Open("syspost")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, model.Response{Code: 0, Message: err.Error()})
	}

	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, model.Response{Code: 0, Message: err.Error()})
	}

	fileSize := fileinfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, model.Response{Code: 0, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.Response{Code: 200, Message: string(buffer)})
}
