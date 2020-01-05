/*
爬虫模块：爬取中国政府新闻官网
包括了爬虫和保存成数据库，只爬取当天的新闻，根据文章的题目来查重
包括crawl的统筹功能，定时重启爬虫去检查最新的新闻

查重：每次启动爬虫，先保存最新一条的新闻的标题到内存，同时与数据库中的标题做对比，当天的新闻不会超过200条，可以直接用数据库查重

当天时间判断：每次爬取HTML完成后，每条新闻都与当天的日期做一个比较，如果不是当天的新闻，爬虫退出，不再爬取后面的新闻

不在外部使用死循环，直接在爬虫启动中使用循环爬取，外部只需要调用一个启动文件即可

最后一次修改的时间：2020-1-5

TODO : 	查重还没有写，从数据库查取最新的时间，然后和爬取到最新的新闻的时间做对比
*/
package main

import "net/http"
import "fmt"
import "strings"
import "strconv"
import "sync"
import "time"

import "github.com/PuerkitoBio/goquery"
import "gopkg.in/mgo.v2"

var govNewsCrollAddress string = "http://sousuo.gov.cn/column/30611/"     // 257.htm
var govNewsCrollLastDate string = "0"                                     // 最后一篇文章的时间，如果没有则从数据库读取
var allPages int = 5                                                      // 一般一天的新闻不会超过5个页面
var savePath string = "localhost:27017"                                   // 如果是保存成文件就是文件地址，如果是数据库就是数据库地址
var latest string = ""                                                    // 当天的日期
var go_sync sync.WaitGroup

func govNewsCrollHTMLString(address string) (*http.Response ,error){
	// function : 获取html的代码
	// param address : 网址地址
	// return : 返回html代码 类型：http.bodyEOFSignal
	resp,err:=http.Get(address)
	if err != nil{
		fmt.Println("获取中华人民共和国人民网回应失败 :",err)
	}

	return resp,err
}

func govNewsCrollTile(resp *http.Response)(string,error){
	// function : 传入resp的Body内容，然后获取文章的题目和时间
	// param respBody : HTML对象
	// return : 文章列表和时间列表
	var title string
	doc,err:=goquery.NewDocumentFromReader(resp.Body)
	if err!=nil{
		fmt.Println("解析中华人民共和国人民网HTML错误 文章标题",err)
		return "",err
	}
	
	// 爬取文章的题目和时间
	doc.Find("div").Each(func(i int,s *goquery.Selection){
		tempTitle := s.Find("h1").Text()
		if len(tempTitle)!=0{
			title=string([]byte(tempTitle[9:]))
		}
	})
	fmt.Println("新闻标题 : ",title)
    return title,err
}

func govNewsCrollHrefContent(resp *http.Response)([]string){
	// function : 具体内容的href链接
	// param resp : 请求的返回
	// return : 返回链接的字符串列表
	var hrefList []string
	doc,err:=goquery.NewDocumentFromReader(resp.Body)
	if err!=nil{
		fmt.Println("解析中华人民共和国人民网HTML错误",err)
	}

	// 爬取文章的链接地址
    doc.Find("a").Each(func(i int,s *goquery.Selection){
        href,isExist := s.Attr("href")
        if isExist==true{
            if "javascript:void(0)"==href || "http://www.gov.cn"==href || strings.Index(href,"http://sousuo.gov.cn/column")!=-1{
                return
            }
			//fmt.Printf("网址 : %s\n",href)
			hrefList=append(hrefList,href)
        }
	})

	return hrefList
}

func govNewsCrollContent(resp *http.Response)(string,string,error){
	// function : 获取具体的文章内容
	// param address : 具体文章地址链接的HTML对象
	// return : 文章内容string，这部分和前面的标题都要存进数据库
	var content string
	var date string
	doc,err:=goquery.NewDocumentFromReader(resp.Body)
	if err!=nil{
		fmt.Println("解析中华人民共和国人民网 新闻滚动 HTML错误",err)
		return "","",err
	}
	
	// 爬取文章的内容和时间
	doc.Find(".pages_content").Each(func(i int,s *goquery.Selection){
		// 爬取逻辑
		title := s.Find("p").Text()
		// 检查无效打印
        if strings.Index(title,"下一页")!=-1 || strings.Index(title,"上一页")!=-1{
            return
        }
		content=title
	})
	
	// 文章的时间
	doc.Find("div[class=pages-date]").Each(func(i int,s *goquery.Selection){
		date=string([]byte(s.Text()[:16]))
	})
	if err!=nil{
		fmt.Println("错误出现 : ",date,content)
	}
	return date,content,err
}

func checkSame(session *mgo.Session,identity string)(bool,error){
	// function : 内容查重并检查是否是最后一个新闻，如果是最后一条或者是重复内容消息则终止
	// param identity : 用于查重的字符串内容或者用时间,为true表示没有重复，为false表示内容已存在
	type TempStruct struct{
		Date string `bson:"date"`
		Content string `bson:"content`
		Title string `bson:"title"`
		Id string `bson:"id"`
		From int `bson:"from"`
	}
	var result []TempStruct
	err:=session.Find(nil).All(&result)
	for _,i:=range result{
		if i.Content==identity{
			return false,nil
		}
	}
	return true,err
}

func nowTime(m string)(string){
	// function : 返回需要的当前时间字符串，用来判断是否是当天的新闻
	// return : 返回需要的时间字符串
	now:=time.Now().Format("2006-01-02 15:04:05")
	if m=="day"{
		return string([]byte(now[:10]))              // 返回的示例：2019-12-30
	}
	fmt.Println("留意nowTime函数的传参错误，将要返回x")
	return "x"                                       // 注意在外部留意这个x返回的情况，说明传参不对
}

func checkDay(date string)(bool){
	// function : 新闻的发布时间与本地时间对比，如果时间不一致，说明该新闻不是当天的
	// return :当天的新闻返回true否则返回false
	nowDay:=nowTime("day")
	if date==nowDay{
		return true
	}else{
		return false
	}
}

func ctrlDataset(session *mgo.Session,date string)(error){
	// function : 控制数据库，爬虫的话，只保留过去一天的数据，也就是昨天的数据保存之后，当今天过完后也就是12点时，把昨天的数据删除
	//            当本地时间和数据库的最新时间不一致时候
	// param sessioin : 数据库指针
	// param date : 最新一条新闻时间的日期
	// return : 返回查询数据库的错误
	type TempStruct struct{
		Date string `bson:"date"`
		Content string `bson:"content`
		Title string `bson:"title"`
		Id string `bson:"id"`
		From int `bson:"from"`
	}
	var tempS []TempStruct
	c:=session.DB("crawl").C("govNews")
	err:=c.Find(nil).Sort("date").Limit(1).All(&tempS)
	// 判断数据库里面的数据是否与当天的时间一致
	if strings.Index(tempS[0].date,date)!=-1{
		c.RemoveAll(bson.M{"date": date})
	}
	
	return err
}

func saveAsMongoDB(session *mgo.Session ,title string,content string ,time string ,dataFrom int,id string){
	// function : 保存数据到mongo数据库
	// 读表
	c:=session.DB("crawl").C("govNews")
	c.Insert(map[string]interface{}{"title":title,"content":content,"date":time,"id":id,"from":dataFrom})   // 插入
	
}

func main(){
	// function : 总运行启动函数
	session,err:=mgo.Dial(savePath)
	if err!=nil{
		fmt.Println("连接数据库报错 : ",err)
	}
	// 无限爬虫循环
	for {
		fmt.Println("开始爬取中华人民共和国新闻滚动",time.Now())
		for i:=0;i<allPages;i++{
			var tempAddress string = strings.Join([]string{govNewsCrollAddress,".htm"},strconv.Itoa(i))
			go_sync.Add(1)
			go func(tempAddress string,i int,wg *sync.WaitGroup){
				defer wg.Done()
				respAll ,err:= govNewsCrollHTMLString(tempAddress)
				if err!=nil{
					return
				}
				hrefList := govNewsCrollHrefContent(respAll)

				for _,href := range hrefList{
					respOne ,err:= govNewsCrollHTMLString(href)
					if err!=nil{
						continue
					}
					title,err:=govNewsCrollTile(respOne)

					if err!=nil{
						continue
					}

					newsDate,newsContent,err := govNewsCrollContent(respOne)
					// 过滤无效消息
					if len(newsContent)<4 || err!=nil{
						continue
					}

					if !checkDay(newsDate) || !checkSame(title){
						// 不是当天的新闻或者数据库有重复，都会退出循环
						break
					}

					err:=ctrlDataset(session,nowTime("day"))
					if err!=nil{
						// 查询数据库失败
						continue
					}
					saveAsMongoDB(session,title,newsContent,newsDate,6,strconv.Itoa(i))
				}
			}(tempAddress,i,&go_sync)

			time.Sleep(time.Second * 1)
			fmt.Println("准备关闭线程",i)
			
		}
		go_sync.Wait()
		fmt.Println("爬取中华人民共和国新闻滚动结束",time.Now())
		time.Sleep(time.Second * 1800)
	}
	
}