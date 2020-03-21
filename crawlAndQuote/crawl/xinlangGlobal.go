/*
新浪新闻，国际新闻爬虫
爬虫的地址不是PC端的新浪地址，而是用手机的新浪新闻地址：
https://news.sina.cn/gj?vt=1&pos=8       vt=1的时候是简化版

能直接拿到题目，但是要进去到页面才能拿到时间，链接的标签是
a class="f_card_m_f_a_r"

进入链接之后，新闻标题是: h1 class="art_tit_h1"
时间的标签是：time class="art_time"

只要拿到标题和时间即可，文章内容不需要爬取，注意是一个页面，
爬取到内容直接保存到数据库即可
*/
package main

import "time"
import "fmt"
import "strings"
import "net/http"
import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "github.com/PuerkitoBio/goquery"

var xinlangGlobalNewsTitle string = ""   // 新浪国际新闻新闻标题
var xinlangGlobalNewsTime string = ""    // 新浪国际新闻文章时间
var xinlangGlobalNewsContent string = "" // 新浪国际新闻文章内容
var mgoSession *mgo.Session

func main(){
	// function : 启动新浪国际新闻爬虫
	// 半个小时更新一次
	session,err:=mgo.Dial("0.0.0.0:27017")
	mgoSession=session
	if err!=nil{
		fmt.Println("链接数据库报错")
	}

	for {
		fmt.Println("爬虫开始")
		xinlangGlobalLink()
		time.Sleep(time.Minute * 10)
		// 删除昨天的数据
		err=ctrlDataset(session,nowTime("day"))
		if err!=nil{
			fmt.Println("删除昨天的新浪国际新闻报错",err)
		}
	}
}

func xinlangGlobalLink(){
	// function : 在新闻滚动页面获取到每个新闻的链接
	var xlGlobalNewAddress string = "https://news.sina.cn/gj?vt=1&pos=8"
	var linkList []string      // 新浪国际新闻滚动新闻的链接列表
	resp,err:=getHTMLResponse(xlGlobalNewAddress)
	if err!=nil{
		fmt.Println("get xinlangGlobalHTML list failed")
	}
	doc,errd:=goquery.NewDocumentFromReader(resp.Body)
	if errd!=nil{
		fmt.Println("explain html list doc failed")
	}
	doc.Find("a").Each(func(i int,s *goquery.Selection){
		href,isExist := s.Attr("href")
		if isExist==true{
			// 过滤http
			if strings.Index(href,"https")!=-1{
				linkList=append(linkList,href)
			}
		}
	})
	
	for _,linkValue:=range(linkList){
		xinlangGlobalHTML(linkValue)

		fmt.Println("获取到新浪国际新闻",xinlangGlobalNewsTitle)
		// 查重,如果查重发现有重复就进入等待时间
		sameResult:=checkSame(mgoSession,xinlangGlobalNewsTitle)
		if sameResult==false{
			continue
		}
		// 把数据保存进数据库
		saveAsMongoDB(mgoSession,xinlangGlobalNewsTitle,xinlangGlobalNewsContent,xinlangGlobalNewsTime,xinlangGlobalNewsTime)
		
	}
	
}

func xinlangGlobalHTML(link string){
	// function : 通过链接获取新闻滚动的html
	resp,err:=getHTMLResponse(link)
	if err!=nil{
		fmt.Println("get xinlangGlobalHTML failed")
	}
	doc,errd:=goquery.NewDocumentFromReader(resp.Body)
	if errd!=nil{
		fmt.Println("explain html doc failed")
	}
	
	doc.Find("div").Each(func(i int,s *goquery.Selection){
		title := s.Find("h1").Text()
		content := s.Find("p").Text()
		time:=nowTime("x")
		if len(title)>4 && len(content)>4{
			// 把新闻内容赋值给全局变量
			xinlangGlobalNewsTitle=title
			xinlangGlobalNewsTime=time
			xinlangGlobalNewsContent=content
		}
		
	})
	
}

func getHTMLResponse(address string)(*http.Response,error){
	// function : 通过一个网址获取该网页的html对象
	// param address : 网页链接
	// return : 一个http的返回对象
	resp,err:=http.Get(address)
	if err != nil{
		fmt.Println("请求获取html失败 :",err)
	}

	return resp,err
}

func nowTime(m string)(string){
	// function : 返回需要的当前时间字符串
	// return : 返回需要的时间字符串
	now:=time.Now().Format("2006-01-02 15:04:05")
	if m=="day"{
		return string([]byte(now[:10]))              // 返回的示例：2019-12-30
	}
	return now
}

func saveAsMongoDB(session *mgo.Session ,title string,content string ,time string ,id string){
	// function : 保存数据到mongo数据库
	// 读表
	c:=session.DB("crawl").C("xinlangGlobal")
	c.Insert(map[string]interface{}{"title":title,"content":content,"date":time,"id":id,"from":3})   // 插入
	fmt.Println("插入数据到数据库 : ",title)
}

func checkSame(session *mgo.Session,identity string)(bool){
	// function : 内容查重并检查是否是最后一个新闻，如果是最后一条或者是重复内容消息则终止
	// param identity : 用于查重的字符串内容或者用时间,为true表示没有重复，为false表示内容已存在
	type TempStruct struct{
		Date string `bson:"date"`
		Content string `bson:"content"`
		Title string `bson:"title"`
		Id string `bson:"id"`
		From int `bson:"from"`
	}
	var result []TempStruct
	err:=session.DB("crawl").C("xinlangGlobal").Find(nil).All(&result)
	if err!=nil{
		fmt.Println("数据库查询错误报错")
		return false
	}
	for _,i:=range result{
		if i.Title==identity{
			return false
		}
	}
	return true
}

func ctrlDataset(session *mgo.Session,date string)(error){
	// function : 控制数据库，爬虫的话，只保留过去一天的数据，也就是昨天的数据保存之后，当今天过完后也就是12点时，把昨天的数据删除
	//            当本地时间和数据库的最新时间不一致时候
	// param sessioin : 数据库指针
	// param date : 最新一条新闻时间的日期
	// return : 返回查询数据库的错误
	type TempStruct struct{
		Date string `bson:"date"`
		Content string `bson:"content"`
		Title string `bson:"title"`
		Id string `bson:"id"`
		From int `bson:"from"`
	}
	var tempS []TempStruct
	c:=session.DB("crawl").C("xinlangGlobal")
	err:=c.Find(nil).Limit(1).All(&tempS)
	if len(tempS)==0{
		// 如果一开始数据库就没有，那么就跳过
		return nil
	}
	// 判断数据库里面的数据是否与当天的时间一致
	if strings.Index(tempS[0].Date,date)==-1{
		fmt.Println("----------------",tempS[0].Date,date)
		c.RemoveAll(bson.M{"date": tempS[0].Date})
	}
	
	return err
}