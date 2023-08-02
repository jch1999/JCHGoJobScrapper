package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	// salary   string // salary가 있는 것도 없는 것도 존재. 위치가 비정규적
}

//3:30
var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python"

func main() {
	var jobs []extractedJob
	totalPages := getPages()
	// fmt.Println(totalPages)
	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		// extractedJobs... 는 extractedJobs의 content를 가져온다는 의미?
    //slice 크기가 정해지지 않은 배열
		jobs = append(jobs, extractedJobs...)
	}
	fmt.Println(jobs)
}
func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&recruitPage=" + strconv.Itoa((page))
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	searchCards := doc.Find(".item_recruit")
	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}
func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	//res.Body는 byte인데, 입력과 출력 IO라고 한다... c#의 streamreader같은 거 같다.
	//따라서 닫아줄 필요가 있다.
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// fmt.Println(doc)
	//.pagintion이 하나라서 Each를 써도 문제가 없는 건가보다
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// fmt.Println(s.Find("a").Length())
		pages = s.Find("a").Length()
	})
	return pages
}

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("value")
	// fmt.Println(id)
	title := cleanString(card.Find(".area_job>.job_tit>a>span").Text())
	// fmt.Println(title)
	location := cleanString(card.Find(".area_job>.job_condition>span>a").Text())
	// fmt.Println(location)
	// fmt.Println(id, title, location)
	return extractedJob{id: id,
		title:    title,
		location: location}
}

func cleanString(str string) string {
	//Fields는 문자열을 분리시킨다.
	//stringTokenizer와 같은 것?
	//TrimSpace로 양쪽 끝에서의 공백을 제거
	//Join은 배열을 separater을 상요해 합친다.
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Statuis: ", res.StatusCode)
	}
}
