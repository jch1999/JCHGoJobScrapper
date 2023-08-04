package main

import (
	// "net/http"

	// "fmt"
	"os"
	"strings"

	"main/scrapper"

	"github.com/labstack/echo"
)

const fileNAME string = "jobs.csv"

//go echo go 언어 기반의 서버를 만들어주는 패키지?

// echo handler
func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	//요청이 다를 수 있으므로 서버에 파일을 남겨놓는 것은 좋지 못하니 삭제
	defer os.Remove(fileNAME)
	// fmt.Println(c.FormValue("term"))
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	//Attachment() 첨부파일을 리턴
	return c.Attachment(fileNAME, "job.csv")
}
func main() {
	// Echo instance
	e := echo.New()
	//Set URL
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
	// scrapper.Scrape("python")
}
