/*
爬取中华人民共和国人民网，新闻滚动
开发人：陈常鸿
创建时间：2019-12-17
最后一次修改时间：2019-12-29

功能：
网址示例：http://sousuo.gov.cn/column/30611/251.htm
第一次运行，爬取所有文章，之后的运行，从第一页开始做对比，直到匹配到数据库中最新的新闻标题
爬取新闻到数据库新闻最新的位置。
注意：每个新闻的内容页面结构可能都不一样，有些新闻没有文字只有图片
2019-12-21该新闻滚动页数达到5K+，一定要用多线程实现
保存文件按照时间切分保存成文件

*/
package crawl
import "net/http"
import "fmt"
import "strings"
import "strconv"
import "sync"
import "github.com/PuerkitoBio/goquery"

var govNewsCrollAddress string = "http://sousuo.gov.cn/column/30611/"
var govNewsCrollLastDate string = "0"                                     // 最后一篇文章的时间，如果没有则从数据库读取
var govNewsCrollAllPages int = 5705

func govNewsCrollHTMLString(address string) (*http.Response,error){
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
	// return : 文章大标题
	var title string
	doc,err:=goquery.NewDocumentFromReader(resp.Body)
	if err!=nil{
		fmt.Println("解析中华人民共和国人民网HTML错误",err)
		return "",err
	}
	
	// 爬取文章的题目和时间
	doc.Find("div").Each(func(i int,s *goquery.Selection){
		tempTitle := s.Find("h1").Text()
		if len(tempTitle)!=0{
			title=string([]byte(tempTitle[9:]))
		} 
    })
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
			fmt.Printf("网址 : %s\n",href)
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
	return date,content,err
}

func govNewsCrollCheckSame(identity string)(bool){
	// function : 内容查重并检查是否是最后一个新闻，如果是最后一条或者是重复内容消息则终止
	// param identity : 用于查重的字符串内容或者用时间
	return true
}

func (c *crawlListStruct)govNewsCrollRun(){
	// function : 总运行启动函数
	fmt.Println("开始爬取中华人民共和国新闻滚动")
	var go_sync sync.WaitGroup
	session ,err:= settingMongo("127.0.0.1","27017","")
	if err!=nil{
		fmt.Println("连接数据库报错 : ",err)
		return
	}
	for i:=0;i<govNewsCrollAllPages;i++{
		go_sync.Add(1)
		var tempAddress string = strings.Join([]string{govNewsCrollAddress,".htm"},strconv.Itoa(i))
		var tempCheckSame bool

		go func(tempAddress string,i int,wg *sync.WaitGroup){
			defer wg.Done()
			respAll,err := govNewsCrollHTMLString(tempAddress)
			if err!=nil{
				return
			}
			//titleList := govNewsCrollAllTile(respAll)
			hrefList := govNewsCrollHrefContent(respAll)

			for _,href := range hrefList{
				respOne ,err:= govNewsCrollHTMLString(href)
				if err!=nil{
					continue
				}
				title,err :=govNewsCrollTile(respOne)
				if err!=nil{
					continue
				}
				date,newsContent,err := govNewsCrollContent(respOne)
				tempCheckSame= govNewsCrollCheckSame(newsContent)
				// 检查内容查重,如果检查到查重则推出循环，然后等待还没有入库的数据线程结束
				if tempCheckSame==true{
					break
				}
				// 过滤无效消息
				if len(newsContent)<4 || err!=nil{
					continue
				}
				saveAsMongoDB(session,"crawl","govNews",title,newsContent ,date,6,"id")
			}
		}(tempAddress,i,&go_sync)
		

	}
	
	go_sync.Wait()
	fmt.Println("爬取中华人民共和国新闻滚动结束")
}