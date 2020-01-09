/*
爬虫模块：爬取中国政府新闻官网
包括了爬虫和保存成数据库，只爬取当天的新闻，根据文章的题目来查重
包括crawl的统筹功能，定时重启爬虫去检查最新的新闻

查重：每次启动爬虫，先保存最新一条的新闻的标题到内存，同时与数据库中的标题做对比，当天的新闻不会超过200条，可以直接用数据库查重

当天时间判断：每次爬取HTML完成后，每条新闻都与当天的日期做一个比较，如果不是当天的新闻，爬虫退出，不再爬取后面的新闻

不在外部使用死循环，直接在爬虫启动中使用循环爬取，外部只需要调用一个启动文件即可

最后一次修改的时间：2020-1-9

注意 : goquery分析HTML对象，只能分析一次，导致新闻文章的时间和内容无法分析出来,就算把把HTML对象分成多个对象也没有用
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
import "gopkg.in/mgo.v2/bson"

var govNewsCrollAddress string = "http://sousuo.gov.cn/column/30611/"     // 257.htm
var govNewsCrollLastDate string = "0"                                     // 最后一篇文章的时间，如果没有则从数据库读取
var allPages int = 5                                                      // 一般一天的新闻不会超过5个页面
var savePath string = "0.0.0.0:27017"                                   // 如果是保存成文件就是文件地址，如果是数据库就是数据库地址
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

func getMultiSession(resp *http.Response)(*goquery.Document,*goquery.Document,*goquery.Document){
	// function : 通过网页HTML的对象返回多个goquery的doc对象
	doc1,err1:=goquery.NewDocumentFromReader(resp.Body)
	doc2,err2:=goquery.NewDocumentFromReader(resp.Body)
	doc3,err3:=goquery.NewDocumentFromReader(resp.Body)
	if err1!=nil || err2!=nil || err3!=nil{
		fmt.Println("分析goquery的错误 : ",err1,err2,err3)
		panic("分析goquery的错误")
	}
	return doc1,doc2,doc3
}

func govNewsCrollTile(doc *goquery.Document)(string,error){
	// function : 传入resp的Body内容，然后获取文章的题目和时间
	// param respBody : HTML对象
	// return : 文章列表和时间列表
	var title string
	
	// 爬取文章的题目和时间
	doc.Find("div").Each(func(i int,s *goquery.Selection){
		tempTitle := s.Find("h1").Text()
		if len(tempTitle)!=0{
			title=string([]byte(tempTitle[9:]))
		}
	})
	fmt.Println("新闻标题 : ",title)
    return title,nil
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

func govNewsCrollContent(doc2 *goquery.Document,doc3 *goquery.Document)(string,string,error){
	// function : 获取具体的文章内容
	// param address : 具体文章地址链接的HTML对象
	// return : 文章内容string，这部分和前面的标题都要存进数据库
	var content string = "x"
	var date string = nowTime("day")

	// 文章的时间
	doc2.Find("div[class=pages-date]").Each(func(i int,s *goquery.Selection){
		date=string([]byte(s.Text()[:16]))
		fmt.Println("文章的发表时间 : ",s.Text())
	})

	// 爬取文章的内容和时间
	doc3.Find(".pages_content").Each(func(i int,s *goquery.Selection){
		fmt.Println("文章内容 : ",s.Text())
		// 爬取逻辑
		title := s.Find("p").Text()
		// 检查无效打印
        if strings.Index(title,"下一页")!=-1 || strings.Index(title,"上一页")!=-1{
            
        }else{
			content=title
		}
		
	})

	return date,content,nil
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
	err:=session.DB("crawl").C("govNews").Find(nil).All(&result)
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
		Content string `bson:"content"`
		Title string `bson:"title"`
		Id string `bson:"id"`
		From int `bson:"from"`
	}
	var tempS []TempStruct
	c:=session.DB("crawl").C("govNews")
	err:=c.Find(nil).Sort("date").Limit(1).All(&tempS)
	if len(tempS)==0{
		// 如果一开始数据库就没有，那么就跳过
		return nil
	}
	// 判断数据库里面的数据是否与当天的时间一致
	if strings.Index(tempS[0].Date,date)!=-1{
		c.RemoveAll(bson.M{"date": date})
	}
	
	return err
}

func saveAsMongoDB(session *mgo.Session ,title string,content string ,time string ,dataFrom int,id string){
	// function : 保存数据到mongo数据库
	// 读表
	c:=session.DB("crawl").C("govNews")
	c.Insert(map[string]interface{}{"title":title,"content":content,"date":time,"id":id,"from":dataFrom})   // 插入
	fmt.Println("插入数据到数据库 : ",title)
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
					doc1,doc2,doc3:=getMultiSession(respOne)

					title,_:=govNewsCrollTile(doc1)
					
					newsDate,newsContent,err := govNewsCrollContent(doc2,doc3)
					// 过滤无效消息
					if len(newsDate)<4 || err!=nil{
						fmt.Println("无效消息被过滤 : ",newsDate,newsContent,err)
						continue
					}

					if !checkDay(newsDate) || !checkSame(session,title){
						// 不是当天的新闻或者数据库有重复，都会退出循环
						fmt.Println("不是当天的新闻或者数据库有重复")
						break
					}

					err=ctrlDataset(session,nowTime("day"))
					if err!=nil{
						// 查询数据库失败
						fmt.Println("查询数据库失败")
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