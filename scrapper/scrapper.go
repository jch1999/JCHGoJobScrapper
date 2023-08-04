package scrapper

import (
	// "encoding/csv"
	"fmt"
	"log"
	"net/http"
	// "os"
	"github.com/tsak/concurrent-csv-writer"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	// salary   string // salary가 있는 것도 없는 것도 존재.
}

func Scrape(term string) {
	var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=" + term
	var jobs []extractedJob
	totalPages := getPages(baseURL)

	mainChannel := make(chan []extractedJob)
	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, mainChannel)
	}
	for i := 0; i < totalPages; i++ {
		extractedJobs := <-mainChannel
		// extractedJobs... 는 extractedJobs의 content를 가져온다는 의미?
		//배열을 추가하는게 아닌 배열의 내용물을 추가하기 위해 배열... 사용
		jobs = append(jobs, extractedJobs...)
	}
	// fmt.Println(jobs)
	// writeJobs(jobs)
	writeJobsRoutine(jobs)
	fmt.Println("Done, extracted ", len(jobs))
}

func getPage(page int, baseURL string, mainChannel chan []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
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
		// job := extractJob(card, c)
		// jobs = append(jobs, job)
		go extractJob(card, c)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mainChannel <- jobs
}

func getPages(baseURL string) int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	//res.Body는 byte인데, 입력과 출력 IO라고 한다... c#의 streamreader같은 거 같다.
	//따라서 닫아줄 필요가 있다.
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	//.pagintion이 하나라서 Each를 써도 문제가 없는 건가보다
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// fmt.Println(s.Find("a").Length())
		pages = s.Find("a").Length()
	})
	return pages
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("value")
	title := CleanString(card.Find(".area_job>.job_tit>a>span").Text())
	location := CleanString(card.Find(".area_job>.job_condition>span>a").Text())

	c <- extractedJob{id: id,
		title:    title,
		location: location}
}

// CleanString cleans a string
func CleanString(str string) string {
	//Fields는 문자열을 분리시킨다.
	//stringTokenizer와 같은 것?
	//TrimSpace로 양쪽 끝에서의 공백을 제거
	//Join은 배열을 separater을 사용해 합친다.
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func writeJobsRoutine(jobs []extractedJob) {
	file, err := ccsv.NewCsvWriter("jobs.csv")
	checkErr(err)

	defer file.Close()

	headers := []string{"ID", "Title", "Location"}

	wErr := file.Write(headers)
	checkErr(wErr)

	done := make(chan bool)

	for _, job := range jobs {
		go func(job extractedJob) {
			jwErr := file.Write([]string{"https://www.saramin.co.kr/zf_user/jobs/relay/view?isMypage=no&rec_idx=" + job.id, job.title, job.location})
			checkErr(jwErr)
			// fmt.Println("ah")
			done <- true
		}(job)
	}

	for i := 0; i < len(jobs); i++ {
		<-done
	}
}

// func writeJobs(jobs []extractedJob) {
// 	file, err := os.Create("jobs.csv")
// 	checkErr(err)

// 	w := csv.NewWriter(file)
// 	// Flush는 파일에 데이터를 입력하는 함수이다.
// 	//defer은 함수가 끝나는 시점에 실행된다.
// 	defer w.Flush()

// 	headers := []string{"ID", "Title", "Location"}

// 	wErr := w.Write(headers)
// 	checkErr(wErr)

// 	for _, job := range jobs {
// 		jobSlice := []string{"https://www.saramin.co.kr/zf_user/jobs/relay/view?isMypage=no&rec_idx=" + job.id, job.title, job.location}
// 		jwErr := w.Write(jobSlice)
// 		checkErr(jwErr)
// 	}
// }

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
