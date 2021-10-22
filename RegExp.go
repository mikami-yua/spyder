package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
并发爬取思路
1.初始化数据管道
2.爬虫写出：每页创建一个协程，向管道中添加图片链接（一共六页创建6个协程）
3.任务统计协程：检测任务是否完成，完成则关闭数据管道
4.下载协程：从管道里读取链接，并下载

*/

//定义正则
var (
	reQQEmail = `\d+@qq.com`
	reEmail   = `\w+@\w+\.\w+?`
	reLink    = `href="(https?://[\s\S]+?)"`
	rePhone   = `1[3456789]\d\s?\d{4}\s?\d{4}`
	reIdcard  = `[123456789]\d{5}((19\d{2})|(20[01]\d))((0[1-9])|(1[012]))((0[1-9])|([12]\d)|(3[01]))\d{3}[\dXx]`
	reImg     = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

//管道定义
var (
	//存放图片链接的数据管道
	chanImageUrls chan string
	//任务统计管道
	chanTask  chan string
	waitGroup sync.WaitGroup
)

//爬取邮箱
func GetEmail(url string) {
	pageStr := GetPageStr(url)

	//3.过滤数据，得到qq邮箱
	re := regexp.MustCompile(reQQEmail)              //根据正则定义得到一个正则对象
	results := re.FindAllStringSubmatch(pageStr, -1) //取这个页面的全部符合正则的内容

	//4.遍历结果(一个slice)
	for _, result := range results {
		fmt.Println("email:", result[0]) //只使用result相当于从大切片里取小切片
	}
}

//爬取链接
func GetLink(url string) {
	pageStr := GetPageStr(url)
	fmt.Println(pageStr)
	//3.过滤数据
	re := regexp.MustCompile(reLink)
	results := re.FindAllStringSubmatch(pageStr, -1)

	//4.遍历结果(一个slice)
	for _, result := range results {
		fmt.Println("link:", result[1])
	}
}

//爬取手机号
func GetPhone(url string) {
	pageStr := GetPageStr(url)
	fmt.Println(pageStr)
	//3.过滤数据
	re := regexp.MustCompile(rePhone)
	results := re.FindAllStringSubmatch(pageStr, -1)

	//4.遍历结果(一个slice)
	for _, result := range results {
		fmt.Println("phone:", result)
	}
}

//爬取身份证
func GetIdCard(url string) {
	pageStr := GetPageStr(url)
	fmt.Println(pageStr)
	//3.过滤数据
	re := regexp.MustCompile(reIdcard)
	results := re.FindAllStringSubmatch(pageStr, -1)

	//4.遍历结果(一个slice)
	for _, result := range results {
		fmt.Println("IdCard:", result)
	}
}

//爬取图片
func GetImg(url string) {
	pageStr := GetPageStr(url)
	//fmt.Println(pageStr)
	//3.过滤数据
	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)

	//4.遍历结果(一个slice)
	for _, result := range results {
		fmt.Println("Img:", result[0]) //只留图片的链接
	}
}

//处理异常
func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}

//根据url获取内容
func GetPageStr(url string) (pageStr string) {
	//1.请求页面拿到响应
	resp, err := http.Get(url) //访问这地址
	HandleError(err, "http.Get url")
	//get也是资源使用defer关闭资源
	defer resp.Body.Close()

	//2.读取返回的信息
	pageBytes, err := ioutil.ReadAll(resp.Body) //readfile只能从文件读，readall读一个流.读出来是字节数据
	HandleError(err, "ioutil.ReadAll")
	//字节转字符串
	pageStr = string(pageBytes)
	return pageStr
}

//爬取图片链接到管道
//url传递的整页链接
func getImgUrls(url string) {
	urls := getCurPageImg(url) //获得所有图片的urls
	//遍历切片里所有链接，存入数据管道
	for _, url := range urls {
		chanImageUrls <- url //把url塞入管道
	}
	chanTask <- url //标识当前协程完成，监控协程仅仅监控爬数据的协程
	//每完成一个任务，写一条数据，用于监控协程知道已经完成了几个任务
	waitGroup.Done() //监控主协程是否结束
}

//获取当前页的图片链接
func getCurPageImg(url string) (urls []string) {
	pageStr := GetPageStr(url) //获得页面数据
	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("共找到%d条结果\n", len(results))
	for _, result := range results {
		url := result[0]
		urls = append(urls, url) //向返回值切片加入数据
	}
	return urls
}

//任务统计协程
func checkOk() {
	var count int
	for {
		url := <-chanTask
		fmt.Printf("%s 路径完成了爬取任务", url)
		count++
		if count == 6 {
			close(chanImageUrls) //关闭数据管道
			break
		}
	}
	waitGroup.Done()
}

//下载图片协程
func downloadImg() {
	for url := range chanImageUrls {
		name := getNameFromUrl(url) //获取图片名
		ok := DownLoadFileTest(url, name)
		if ok {
			fmt.Printf("%s 下载成功\n", name)
		} else {
			fmt.Printf("%s 下载失败\n", name)
		}
	}
	waitGroup.Done()
}

//获取图片名的方法
func getNameFromUrl(url string) (name string) {
	//取最后一截
	index := strings.LastIndex(url, "/") //返回最后一个杠的未知
	//切出来
	name = url[index+1:]
	//时间戳解决重名
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	name = timePrefix + "_" + name
	return name
}

func main() {
	/*简单功能
	//url:="http://www.baidu.com/s?wd=%E8%B4%B4%E5%90%A7%20%E7%95%99%E4%B8%8B%E9%82%AE%E7%AE%B1&rsv_spt=1&rsv_iqid=0x98ace53400003985&issp=1&f=8&rsv_bp=1&rsv_idx=2&ie=utf-8&tn=baiduhome_pg&rsv_enter=1&rsv_dl=ib&rsv_sug2=0&inputT=5197&rsv_sug4=6345"
	//GetEmail(url)
	//GetLink(url)
	//GetPhone(url)
	//GetIdCard("https://henan.qq.com/a/20171107/069413.htm")

	//GetImg("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/1.html")


	//并发爬取图片，每一页开启一个协程爬取
	//myTestPage()
	//DownLoadFileTest("https://uploadfile.bizhizu.cn/up/ac/6a/82/ac6a825d23c550326a233a16e96c6fd4.jpg","1.jpg")

	*/
	oringUrl := "https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/"
	//1.初始化管道
	chanImageUrls = make(chan string, 1000000) //图片大小未知，按照1000000记录
	chanTask = make(chan string, 6)            //一共六页，需要6个协程，大小是6
	//2.爬虫协程
	for i := 1; i < 7; i++ {
		waitGroup.Add(1) //增加一个 go routine
		go getImgUrls(oringUrl + strconv.Itoa(i) + ".html")
	}
	//3.任务统计协程（任务是否都完成，完成后关闭管道）
	waitGroup.Add(1) //增加协程就需要+1
	go checkOk()
	//4.下载协程，从管道读取链接并下载
	//爬完之后，图片链接都在一个管道中（不要开太多协程，可能会封ip）
	for i := 0; i < 5; i++ { //5个协程从管道中下载
		waitGroup.Add(1)
		go downloadImg()
	}
	waitGroup.Wait()
}

//测试获取页面是否ok
func myTestPage() {
	pageStr := GetPageStr("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/1.html")
	fmt.Println(pageStr)
	//获取链接
	GetImg("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/1.html")
}

//测试下载是否ok
func DownLoadFileTest(url string, figname string) (ok bool) {
	resp, err := http.Get(url)
	HandleError(err, "http.get.url")
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "res.body")
	figname = "C:\\Users\\jia\\go\\src\\spyder\\img\\" + figname
	//println(figname)
	//写出数据
	err = ioutil.WriteFile(figname, bytes, 0666) //最后一个是权限
	if err != nil {
		return false
	} else {
		return true
	}
}
