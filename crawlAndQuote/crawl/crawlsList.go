/*
每个爬虫函数的启动的列表
开发人：陈常鸿
创建时间：2019-12-15
最后一次修改时间：2019-12-15

功能：就是每个爬虫启动的方法列表
*/
package crawl

// 每个爬虫启动方法的列表，每次增加新的爬虫任务，先在工程目录写好对应的下载和解析的代码，然后把最后的执行函数
// 添加进这个列表之内
type crawlList interface{

}

// 实例化crawlList的结构体
type crawlListStruct struct{
	funcList []func()
}

// 把每个站点的爬虫启动代码添加到crawlListStruct的funcList中
func addFuncToCrawlList(cls *crawlListStruct){
	// function : 把每个站点的爬虫启动代码添加到crawlListStruct的funcList中
	// param cls : 爬虫启动函数的结构体对象
	cls.funcList=append(cls.funcList,cls.hasakiTestfunc)   // 这是一个展示用法的用例
}

// 展示用例
func (c *crawlListStruct)hasakiTestfunc(){

}