/*
爬取到的内容持久化到硬盘
开发人：陈常鸿
创建时间：2019-12-15
最后一次修改时间：2019-12-15

功能：
持久化内容到数据库
持久化内容到文件
*/
package crawl

func saveAsJSON(path string ,time string ,content string){
	// function : 持久化到json文件
	// param path : 保存文件的路径
	// param time : 爬取内容的时间,不是持久的时间，而是检测到内容的时间
	// param content : 具体的内容
}

func saveAsTxt(path string ,time string ,content string){
	// function : 持久化到txt文件
	// param path : 保存文件的路径
	// param time : 爬取内容的时间,不是持久的时间，而是检测到内容的时间
	// param content : 具体的内容
}

func saveAsMongoDB(address string ,port string ,password int ,content string ,dataFrom string){
	// function : 保存数据到mongo数据库
	// param address : 数据库地址
	// param port : 连接端口
	// param password : 数据库密码
	// param content : 需要保存到数据库的内容
	// param dataFrom : 数据来源内容
}

func saveAsMySQL(address string ,port string ,password int ,content string ,dataFrom string){
	// function : 保存数据到mysql数据库
	// param address : 数据库地址
	// param port : 连接端口
	// param password : 数据库密码
	// param content : 需要保存到数据库的内容
	// param dataFrom : 数据来源内容
}

func saveAsRedis(address string ,port string ,password int ,content string ,dataFrom string){
	// function : 保存数据到redis数据库
	// param address : 数据库地址
	// param port : 连接端口
	// param password : 数据库密码
	// param content : 需要保存到数据库的内容
	// param dataFrom : 数据来源内容
}