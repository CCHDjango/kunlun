/*
爬取到的内容持久化到硬盘
增加功能，用作操作数据的模块
开发人：陈常鸿
创建时间：2019-12-15
最后一次修改时间：2020-1-3

功能：
持久化内容到数据库
持久化内容到文件
*/
package crawl

import "fmt"
import "encoding/json"
import "io/ioutil"
import "gopkg.in/mgo.v2"
import "strings"
//import "github.com/go-mgo/mgo/bson"

// 全局参数
var mongoURL string
var mySQLURL string
var redisURl string

func saveAsJSON(path string ,time string ,title string ,content string,from int,id string){
	// function : 持久化到json文件
	// param path : 保存文件的路径
	// param time : 爬取内容的时间,不是持久的时间，而是检测到内容的时间
	// param content : 具体的内容
	// param title : 文章的标题
	// param id : id
	// param from : 数据来源
	var js = new(dataType)
	js.Date=time
	js.Title=title
	js.Content=content
	js.From=from
	js.Id=id
	if data,err:=json.Marshal(js);err==nil{
		// data type : byte
		fmt.Println(string(data),err)
		_=ioutil.WriteFile(path, data, 0755)
	}else{
		fmt.Println("写入json文件报错结果错误",err)
	}
}

func saveAsTxt(path string ,time string ,title string,content string,from int,id string){
	// function : 持久化到txt文件
	// param path : 保存文件的路径
	// param time : 爬取内容的时间,不是持久的时间，而是检测到内容的时间
	// param content : 具体的内容
}

func settingMongo(address string,port string,password string)(*mgo.Session,error){
	// function : 传入参数，并返回一个MongoDB的连接
	// param address : 数据库地址
	// param port : 数据库端口
	// param password : 数据库密码
	// return 一个数据库连接对象
	session,err:=mgo.Dial(strings.Join([]string{address,port},":"))
	return session ,err
}

func saveAsMongoDB(session *mgo.Session ,dbName string,tbName string,title string,content string ,time string ,dataFrom int,id string){
	// function : 保存数据到mongo数据库
	// param session : 数据库连接对象
	// param title : 文章标题
	// param time : 时间
	// param content : 需要保存到数据库的内容
	// param dataFrom : 数据来源内容
	// param dbName : 数据库名称
	// param tbName : 表名

	// var mongoURL string = strings.Join([]string{address,port},":")
	// session,err:=mgo.Dial(mongoURL)
	c:=session.DB(dbName).C(tbName)
	c.Insert(map[string]interface{}{"id":id,"date":time,"content":content,"title":title,"from":dataFrom})
}

func saveAsMySQL(address string ,port string ,password string ,title string,content string ,time string ,dataFrom int,id string){
	// function : 保存数据到mysql数据库
	// param address : 数据库地址
	// param port : 连接端口
	// param password : 数据库密码
	// param title : 文章标题
	// param time : 时间
	// param content : 需要保存到数据库的内容
	// param dataFrom : 数据来源内容

	// 如果数据库没有密码的情况，那么passdword为default
	if password=="default"{

	}else{
	
	}
}

func saveAsRedis(address string ,port string ,password string ,title string,content string ,time string ,dataFrom int,id string){
	// function : 保存数据到redis数据库
	// param address : 数据库地址
	// param port : 连接端口
	// param password : 数据库密码
	// param title : 文章标题
	// param time : 时间
	// param content : 需要保存到数据库的内容
	// param dataFrom : 数据来源内容

	// 如果数据库没有密码的情况，那么passdword为default
	if password=="default"{

	}else{
		
	}
}

func saveDefault(content string){
	// fuction : 用来占位置的保存爬虫内容的函数
}

func mgoInsert(data map[string]interface{}){
	// function : mongo插入
	// param data : 需要插入的数据
}

func mgoDelete(id string){
	// function : 根据id删除数据库的某一条
	// param id : id字段
}

func mgoDeleteAll(tableName string){
	// function : 一次删除一个表
	// param tableName : 需要删除的表名
}