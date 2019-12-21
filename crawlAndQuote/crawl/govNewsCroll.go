/*
爬取中华人民共和国人民网，新闻滚动
开发人：陈常鸿
创建时间：2019-12-17
最后一次修改时间：2019-12-21

功能：
网址示例：http://sousuo.gov.cn/column/30611/251.htm
第一次运行，爬取所有文章，之后的运行，从第一页开始做对比，直到匹配到数据库中最新的新闻标题
爬取新闻到数据库新闻最新的位置。
注意：每个新闻的内容页面结构可能都不一样，有些新闻没有文字只有图片
*/
package crawl
import "net/http"
import "fmt"
import "strings"
import "github.com/PuerkitoBio/goquery"

var govNewsCrollAddress string = "http://sousuo.gov.cn/column/30611/257.htm"

func govNewsCrollHTMLString(address string) (*http.Response){
	// function : 获取html的代码
	// param address : 网址地址
	// return : 返回html代码 类型：http.bodyEOFSignal
	resp,err:=http.Get(address)
	if err != nil{
		fmt.Println("获取中华人民共和国人民网回应失败 :",err)
		panic("获取中华人民共和国人民网回应失败")
	}

	return resp
}

func govNewsCrollAllTile(resp *http.Response)([]string){
	// function : 传入resp的Body内容，然后获取文章的题目和时间
	// param respBody : HTML对象
	// return : 文章列表和时间列表
	var titleList []string
	doc,err:=goquery.NewDocumentFromReader(resp.Body)
	if err!=nil{
		fmt.Println("解析中华人民共和国人民网HTML错误",err)
	}
	
	// 爬取文章的题目和时间
	doc.Find("li").Each(func(i int,s *goquery.Selection){
		title := s.Find("a").Text()
		date := s.Find("span").Text()
		if i<19{
			fmt.Printf("Review %d: %s - %s\n", i, title, date)
			titleList=append(titleList,title)
		}
    })
    return titleList
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
            if "javascript:void(0)"==href || "http://www.gov.cn"==href{
                return
            }
			fmt.Printf("网址 : %s\n",href)
			hrefList=append(hrefList,href)
        }
	})

	return hrefList
}

func govNewsCrollContent(resp *http.Response)(string){
	// function : 获取具体的文章内容
	// param address : 具体文章地址链接的HTML对象
	// return : 文章内容string，这部分和前面的标题都要存进数据库
	var content string
	var date string
	doc,err:=goquery.NewDocumentFromReader(resp.Body)
	if err!=nil{
		fmt.Println("解析中华人民共和国人民网 新闻滚动 HTML错误",err)
	}
	
	// 爬取文章的题目和时间
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
	return strings.Join([]string{date,content}," | ")
}

func (c *crawlListStruct)govNewsCrollRun(){
	// function : 总运行启动函数
	fmt.Println("开始爬取中华人民共和国新闻滚动")
	respAll := govNewsCrollHTMLString(govNewsCrollAddress)
	//titleList := govNewsCrollAllTile(respAll)
	hrefList := govNewsCrollHrefContent(respAll)

	for _,href := range hrefList{
		respOne := govNewsCrollHTMLString(href)
		newsContent := govNewsCrollContent(respOne)
		// 过滤无效消息
        if len(newsContent)<4{
            continue
        }
		saveDefault(newsContent)                      // TODO : 正式运行需要更改保存函数
	}
	fmt.Println("爬取中华人民共和国新闻滚动结束")
}
