/*
常驻服务，接受外部请求，然后从数据库打包数据并返回数据
外部请求有两种，一种是高权限的请求，返回所有数据
一种是低权限的请求，返回低密度的数据
高权限的请求来自于开发者本地的请求，用于每天保存数据到本地的数据库
低权限的请求来自前端请求，用于数据展示，数据展示不需要太多数据，所以这个请求是过滤过的数据
同时定时启动爬虫去爬取数据
实时行情为一直连接，如果行情模块报错，只要在此主线程中重启行情模块即可
接受到请求，判断请求类型，然后去读取数据库，最后返回数据
消息中心，如果模块报错，发送消息到指定的地方，比如钉钉，微信，邮箱等，如果是已知情况则可以自动重启重新运行，位置情况则发消息手动重启服务
服务接口有两个，一个是下载一天的数据，包括一天的行情和一天的新闻
另一个接口是前端请求数据展示的接口，前端展示请求的接口请求的是一天内30分钟的K线，以及全部的新闻舆情数据

请求返回的数据是：
{
	"news":[
		{"date":"","content":"","title":""},
		{"date":"","content":"","title":""}
	],
	"quote":[
		{"date":"","open":"","high":"","low":"","close":"","frequency":""},
		{"date":"","open":"","high":"","low":"","close":"","frequency":""}
	]
}
最后一次修改时间：2020-1-9

注意：此代码没有做过多的设计，在业务没有确定的情况下，不要做过多的设计，经济最优原则
在实现功能的条件下，用最短的时间，最简单的实现方法
本来应该是从该程序启动新闻爬虫的，但是在服务器用screen方法另起线程让新闻爬虫运行也一样，所以暂时就使用screen的方式启动爬虫和数据库，服务程序另外再用
screen启动即可
*/
package main

import "fmt"
import "sync"
import "github.com/gin-gonic/gin"
import "gopkg.in/mgo.v2"

// 全局变量们 Global Vars
var mgoURL string = "0.0.0.0:27017"
var session *mgo.Session
var go_sync sync.WaitGroup

type newsStruct struct{
	Date string `bson:"date"`
	Content string `bson:"content"`
	Title string `bson:"title"`
	Id string `bson:"id"`
	From int `bson:"from"`
}
var frontChan chan []newsStruct
var dataChan chan []newsStruct

func frontServer(context *gin.Context){
	// function : 前端服务接口服务函数
	fmt.Println("前端服务接口启动")
	fc:=make(chan []newsStruct)
	go_sync.Add(1)
	go getFrontData(fc,&go_sync)
	c:=<-fc
	context.JSON(200,gin.H{
		"code":200,
		"success":true,
		"date":c[0].Date,
		"title":c[0].Title,
	})
	go_sync.Wait()
}

func dataServer(context *gin.Context){
	// function : 数据请求服务
	fmt.Println("全数据请求")
	context.JSON(200,gin.H{
		"code":200,
		"success":true,
		"data":"hasaki",
	})
}

func getFrontData(fc chan []newsStruct,wg *sync.WaitGroup){
	// function :从数据库查询前端展示需要的数据，一天比特币的半小时行情，以及当天的新闻
	defer wg.Done()
	var result []newsStruct
	err:=session.DB("crawl").C("govNews").Find(nil).All(&result)
	fmt.Println("get front data from database : ",result)
	if err!=nil{
		fmt.Println("查询数据库报错 : ",err)
	}else{
		fc<-result
	}
}

func getSaveData(){
	// function : 从数据库查询所有历史数据需要的数据，一天的比特币一分钟行情，以及当天的新闻
	// param SDChan : 数据管道，把保存历史数据从管道传出去到客户端
	var result []newsStruct
	err:=session.DB("crawl").C("govNews").Find(nil).All(&result)
	if err!=nil{
		fmt.Println("查询数据库报错 : ",err)
	}else{
		dataChan<-result
	}
}

func main(){
	fmt.Println("启动数据服务")
	// 连接数据库
	session,_=mgo.Dial(mgoURL)
	
	// Engin指针
    router := gin.Default()

	router.GET("/frontServer", frontServer)
	router.GET("/dataServer",dataServer)
    // 指定地址和端口号
	router.Run("0.0.0.0:8888")                    // 如果是云服务改成0.0.0.0:8888
	
	// 启动爬虫线程
	//go crawlRun()
	
	// 行情启动方法
}