/*
昆仑平台数据爬取模块调度引擎
开发人：陈常鸿
创建时间：2019-12-15
最后一次修改时间：2019-12-15

功能：
定时轮询一次要爬取内容的网站
保存最后一次内容的id到持久化文件
加载webside.json中的网站，并根据name来调取函数对应的爬虫方法
*/
package crawl

type engine interface{
	timer()
	saveContentId(id string)
	loadContentId()
	loadWebsideJson()
}

func timer(){
	// function : 计时器，多久启动轮询一次爬虫
}

func saveContentId(id string){
	// function : 保存内容id到本地，用于断点续爬，如果爬虫程序死掉了，那么重启程序读取id断点续爬
	// param id : 内容的id
}

func loadContentId(){
	// function : 程序启动的时候，查询是否有过往的内容id，然后在内容id的位置开始爬取
}

func loadWebsideJson(){
	// function : 下一个爬虫轮询开始之前，都加载一次webside.json，根据json里面的内容进行爬取内容
}