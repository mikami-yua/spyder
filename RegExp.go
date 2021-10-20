package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

/*
并发爬取思路

 */





//定义正则
var (
	reQQEmail =`\d+@qq.com`
	reEmail=`\w+@\w+\.\w+?`
	reLink=`href="(https?://[\s\S]+?)"`
	rePhone  = `1[3456789]\d\s?\d{4}\s?\d{4}`
	reIdcard = `[123456789]\d{5}((19\d{2})|(20[01]\d))((0[1-9])|(1[012]))((0[1-9])|([12]\d)|(3[01]))\d{3}[\dXx]`
	reImg    = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

//爬取邮箱
func GetEmail(url string)  {
	pageStr := GetPageStr(url)

	//3.过滤数据，得到qq邮箱
	re := regexp.MustCompile(reQQEmail)               //根据正则定义得到一个正则对象
	results := re.FindAllStringSubmatch(pageStr, -1) //取这个页面的全部符合正则的内容

	//4.遍历结果(一个slice)
	for _ ,result:=range results{
		fmt.Println("email:",result[0])//只使用result相当于从大切片里取小切片
	}
}

//爬取链接
func GetLink(url string)  {
	pageStr := GetPageStr(url)
	fmt.Println(pageStr)
	//3.过滤数据
	re := regexp.MustCompile(reLink)
	results := re.FindAllStringSubmatch(pageStr, -1)

	//4.遍历结果(一个slice)
	for _ ,result:=range results{
		fmt.Println("link:",result[1])
	}
}

//爬取手机号
func GetPhone(url string)  {
	pageStr := GetPageStr(url)
	fmt.Println(pageStr)
	//3.过滤数据
	re := regexp.MustCompile(rePhone)
	results := re.FindAllStringSubmatch(pageStr, -1)

	//4.遍历结果(一个slice)
	for _ ,result:=range results{
		fmt.Println("phone:",result)
	}
}

//爬取身份证
func GetIdCard(url string)  {
	pageStr := GetPageStr(url)
	fmt.Println(pageStr)
	//3.过滤数据
	re := regexp.MustCompile(reIdcard)
	results := re.FindAllStringSubmatch(pageStr, -1)

	//4.遍历结果(一个slice)
	for _ ,result:=range results{
		fmt.Println("IdCard:",result)
	}
}

//爬取图片
func GetImg(url string)  {
	pageStr := GetPageStr(url)
	//fmt.Println(pageStr)
	//3.过滤数据
	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)

	//4.遍历结果(一个slice)
	for _ ,result:=range results{
		fmt.Println("Img:",result[0])//只留图片的链接
	}
}


//处理异常
func HandleError(err error,why string)  {
	if err != nil {
		fmt.Println(why,err)
	}
}

//根据url获取内容
func GetPageStr(url string) (pageStr string) {
	//1.请求页面拿到响应
	resp, err := http.Get(url)//访问这地址
	HandleError(err,"http.Get url")
	//get也是资源使用defer关闭资源
	defer resp.Body.Close()

	//2.读取返回的信息
	pageBytes, err := ioutil.ReadAll(resp.Body)//readfile只能从文件读，readall读一个流.读出来是字节数据
	HandleError(err,"ioutil.ReadAll")
	//字节转字符串
	pageStr = string(pageBytes)
	return pageStr
}


func main() {
	//url:="http://www.baidu.com/s?wd=%E8%B4%B4%E5%90%A7%20%E7%95%99%E4%B8%8B%E9%82%AE%E7%AE%B1&rsv_spt=1&rsv_iqid=0x98ace53400003985&issp=1&f=8&rsv_bp=1&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_enter=1&rsv_dl=ib&rsv_sug2=0&inputT=5197&rsv_sug4=6345"
	//GetEmail(url)
	//GetLink(url)
	//GetPhone(url)
	//GetIdCard("https://henan.qq.com/a/20171107/069413.htm")

	//GetImg("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/1.html")


	//并发爬取图片，每一页开启一个协程爬取
	//myTestPage()
	//DownLoadFileTest("https://uploadfile.bizhizu.cn/up/ac/6a/82/ac6a825d23c550326a233a16e96c6fd4.jpg","1.jpg")
}

//测试获取页面是否ok
func myTestPage()  {
	pageStr := GetPageStr("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/1.html")
	fmt.Println(pageStr)
	//获取链接
	GetImg("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/1.html")
}

//测试下载是否ok
func DownLoadFileTest(url string,figname string) (ok bool) {
	resp, err := http.Get(url)
	HandleError(err,"http.get.url")
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err,"res.body")
	figname="C:\\Users\\jia\\go\\src\\spyder\\img\\"+figname
	//println(figname)
	//写出数据
	err=ioutil.WriteFile(figname,bytes,0666)//最后一个是权限
	if err != nil {
		return false
	}else {
		return true
	}
}