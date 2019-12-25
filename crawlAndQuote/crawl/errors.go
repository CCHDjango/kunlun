/*
爬虫错误处理方法
开发人：陈常鸿
创建时间：2019-12-15
最后一次修改时间：2019-12-15

爬虫错误信息
程序报错，连接错误，解析错误，下载报错，持久化报错
*/
package crawl

import "errors"

type errorInterface interface{

}

// 视线错误函数的列表
type errorStruct struct{
	//TestError error
}

var testError error = errors.New("测试错误类型")