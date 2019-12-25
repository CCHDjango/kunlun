/*
爬取cnBeta.com的新闻
开发者：陈常鸿
创建时间：2019-12-19
最后一次修改时间：2019-12-19

功能与注意事项：
cnBeta新闻只爬取题目和摘要
*/
package crawl
import "net/http"
import "fmt"
import "github.com/PuerkitoBio/goquery"

var cnBetaNewsAddress string = "https://www.cnbeta.com/"

func cnBetaNewsHTML(address string) (*http.Response) {
	// function : 获取cnBeta新闻的HTML对象的方法
	// param : HTML地址
	// return : 返回response
	resp,err:=http.Get(address)
	if err != nil{
		fmt.Println("获取cnBeta HTML回应失败 :",err)
		panic("获取cnBeta HTML回应失败")
	}

	return resp
}

func cnBetaNewsAllTitle(resp *http.Response) ([]string){
	// function : 解析HTML获取所有新闻的题目文字
	// param respBody : HTML对象
	// return : 文章列表和时间列表
	var titleList []string
	doc,err:=goquery.NewDocumentFromReader(resp.Body)
	if err!=nil{
		fmt.Println("解析cnBeta HTML错误",err)
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

func (c *crawlListStruct)cnBetaNewsRun(){
	// function : 总运行启动函数
	fmt.Println("开始爬取cnBeta新闻滚动")
	respAll := cnBetaNewsHTML(cnBetaNewsAddress)
	//titleList := govNewsCrollAllTile(respAll)
	contentList := cnBetaNewsAllTitle(respAll)
	for _,newContent :=range contentList{
		saveDefault(newContent)                      // TODO : 正式运行需要更改保存函数
	}

	fmt.Println("爬取cnBeta新闻滚动结束")
}