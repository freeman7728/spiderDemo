/*
 * @Description:
 * @author: freeman7728
 * @Date: 2024-08-29 19:28:02
 * @LastEditTime: 2024-08-30 11:15:46
 * @LastEditors: freeman7728
 */
package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	Spider()
}

func Spider() {
	//TODO 发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250", nil)
	if err != nil {
		fmt.Println("err", err)
	}
	//添加请求头使其符合浏览器访问的形式
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	// req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	// req.Header.Set("Cache-Control", "max-age=0")
	// req.Header.Set("Priority", "u=0, i")
	// req.Header.Set("Sec-Ch-Ua", "")
	// req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	// req.Header.Set("Sec-Ch-Ua-Platform", "Windows")
	// req.Header.Set("Sec-Fetch-Dest", "document")
	// req.Header.Set("Sec-Fetch-Mode", "navigate")
	// req.Header.Set("Sec-Fetch-Site", "none")
	// req.Header.Set("Sec-Fetch-User", "?1")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
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
	//TODO 获取节点
	//
	//#content > div > div.article > ol > li
	//#content > div > div.article > ol > li:nth-child(1) > div > div.pic > a > img
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p:nth-child(1)
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > div > span.rating_num
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p.quote > span
	docDetail.Find("#content > div > div.article > ol > li").
		Each(func(i int, s *goquery.Selection) {
			title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text()
			img := s.Find("div > div.pic > a > img")
			imgTmp, ok := img.Attr("src")
			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
			score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
			quote := s.Find("div > div.info > div.bd > p.quote > span").Text()
			if ok {
				fmt.Println("rank:", i)
				fmt.Println("title:", title)
				fmt.Println("imgTmp:", imgTmp)
				fmt.Println("info:", info)
				fmt.Println("score:", score)
				fmt.Println("quote:", quote)
				d, a, y := InfoSplit(info)
				fmt.Println("year:", y)
				fmt.Println("director:", d)
				fmt.Println("actor:", a)
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
