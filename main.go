package main

import (
	"net/http"
)

var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword="

func main() {
	pages := getPages()
}

func getPages() int {
	res, err := http.Get(baseURL)
	return 0
}

// youjinlee19979 months ago
// 2022.10.24 기준으로 영상 속 사이트는 크롤링이 불가능합니다.
// 그래서 저는 일단 사람인으로 대체해서 진행했어요!

// 크롤링 주소는 https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python 이구요,

// id, title, location을 찾을 때는 아래와 같은 코드를 사용했습니다.
// title := cleanString(card.Find(".area_job>.job_tit>a").Text())
// location := cleanString(card.Find(".area_job>.job_condition>span>a").Text())

// summary와 연봉정보는 메인에 안나와서 별도로 하진 않았지만, 강의 듣는데는 큰 문제 없습니다.
// 우선은 3.7강 까지 학습하는데는 큰 문제 없는걸로 확인했습니다.

// 재밌는 강의 만들어주신 니꼬쌤 감사하고, 강의 들으시는 분들도 화이팅입니다!

//감사합니다!
