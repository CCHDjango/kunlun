/*
该爬虫模块的公共对象
开发人：陈常鸿
创建时间：2019-12-19
最后一次修改时间：2019-12-19

功能与注意事项
*/
package crawl

type dataType struct{
	Date string `json:"date"`
	Title string `json:"title"`
	Content string `json:"content"`
	From int `json:"from"`
	Id string
}