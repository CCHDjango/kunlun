/*
公用方法：文件读写，类型转换，字符串判断,字符串拼接
开发者：陈常鸿
创建时间：2019-12-21
最后一次修改时间：2019-12-30

注意事项：
*/
package crawl

import "time"
import "strings"
import "strconv"

func dateJudge(firstDate string,lastDate string)(bool){
	/*
	function : 判断前后时间是否一致，如果一致则返回真，否则假
	param firstDate : 第一个时间 格式：2019-12-21 20:05
	param lastDate : 第二个时间
	return : 如果一致则返回真，否则假
	*/
	if firstDate==lastDate{
		return true
	}else{
		return false
	}
}

func readJson(path string)([]map[string]interface{}){
	/*
	function : 读取json文件并转换成[]map类型
	param path : 需要读取的json文件的路径
	return : 返回
	*/
	var temp []map[string]interface{}
	return temp
}

func sleep(n int){
	/*
	function : 睡眠函数
	param n : 睡眠的秒数
	*/
	if n<0{
		return
	}

	for i:=0;i<=n;i++{
		time.Sleep(time.Second)
	}
}

func nowTime(m string)(string){
	// function : 返回需要的当前时间字符串
	// return : 返回需要的时间字符串
	now:=time.Now().Format("2006-01-02 15:04:05")
	if m=="day"{
		return string([]byte(now[:10]))              // 返回的示例：2019-12-30
	}
	return "x"
}

func strJoin(first string,last string,mid string)(string){
	/*
	function : 拼接合成字符串
	param first : 第一段字符串
	param mid : 中间分割的字符串
	param last : 最后一段字符串
	return : 返回拼接好的字符串
	*/
	return strings.Join([]string{first,last},mid)
}

func intToStr(num int)(string){
	/*
	function : 整型转换成字符串
	param num : 整型
	return : 。。
	*/
	return strconv.Itoa(num)
}