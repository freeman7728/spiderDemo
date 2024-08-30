/*
 * @Description:
 * @author: freeman7728
 * @Date: 2024-08-29 19:28:02
 * @LastEditTime: 2024-08-30 14:32:07
 * @LastEditors: freeman7728
 */
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/PuerkitoBio/goquery"
)

var idx int

const (
	USERNAME = "root"
	PASSWORD = "root"
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "douban_movie"
)

var DB *sql.DB

type movieData struct {
	Title    string `json:"title"`
	Rank     string `json:"rank"`
	ImgUrl   string `json:"imgUrl"`
	Score    string `json:"score"`
	Quote    string `json:"quote"`
	Year     string `json:"year"`
	Director string `json:"director"`
	Actor    string `json:"actor"`
}

func (m *movieData) PrintToScreen() {
	fmt.Println("Title", m.Title)
	fmt.Println("Rank", m.Rank)
	fmt.Println("ImgUrl", m.ImgUrl)
	fmt.Println("Score", m.Score)
	fmt.Println("Quote", m.Quote)
	fmt.Println("Year", m.Year)
	fmt.Println("Director", m.Director)
	fmt.Println("Actor", m.Actor)
}

func main() {
	idx = 1
	InitDB()
	// for i := 0; i < 10; i++ {
	// 	Spider(strconv.Itoa(i * 25))
	// }

}

func Spider(page string) {
	//TODO 发送请求
	client := http.Client{}
	//TODO 分页
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start="+page, nil)
	if err != nil {
		fmt.Println("err", err)
	}
	//添加请求头使其符合浏览器访问的形式
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	defer resp.Body.Close()
	//TODO 解析网页
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("body解析失败", err)
	}

	//#content > div > div.article > ol > li
	//#content > div > div.article > ol > li:nth-child(1) > div > div.pic > a > img
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p:nth-child(1)
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > div > span.rating_num
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p.quote > span
	//TODO 获取节点
	docDetail.Find("#content > div > div.article > ol > li").
		Each(func(i int, s *goquery.Selection) {
			title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text()
			img := s.Find("div > div.pic > a > img")
			imgTmp, ok := img.Attr("src")
			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
			score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
			quote := s.Find("div > div.info > div.bd > p.quote > span").Text()
			if ok {
				d, a, y := InfoSplit(info)
				curMovie := &movieData{
					Title:    title,
					Rank:     strconv.Itoa(idx),
					ImgUrl:   imgTmp,
					Score:    score,
					Quote:    quote,
					Year:     y,
					Director: d,
					Actor:    a,
				}
				fmt.Println(*curMovie)
				//curMovie.PrintToScreen()
				idx += 1
			}
		})

	//TODO 保存信息
}

//TODO 数据的处理
/*
info:
                            导演: 李·昂克里奇 Lee Unkrich / 阿德里安·莫利纳 Adrian Molina   主演: ...
                            2017 / 美国 / 喜剧 动画 奇幻 音乐
*/
//使用在线正则生成表达式
func InfoSplit(info string) (director, actor, year string) {
	yearReg, _ := regexp.Compile(`(\d+)`)
	directorReg, _ := regexp.Compile(`导演:(.*)主演:`)
	actorReg, _ := regexp.Compile(`主演: (.*)`)
	actor = actorReg.FindString(info)
	year = yearReg.FindString(info)
	director = directorReg.FindString(info)
	return
}

// TODO 数据库的初始化
func InitDB() {
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	DB, _ := sql.Open("mysql", path)
	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)
	if err := DB.Ping(); err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println("connect success => ", path)
}
