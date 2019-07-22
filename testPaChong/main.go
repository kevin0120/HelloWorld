package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
网页
//https://www.pengfue.com/index_1.html
*/
func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	//读取网页内容
	buf := make([]byte, 1024*4)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 { //表示读取结束或者出问题
			//fmt.Println("resp.Body.Read err =", err)
			break
		}
		result += string(buf[:n])
	}
	return
}

func SpiderOneJoy(url string) (title, content string, err error) {
	result, err1 := HttpGet(url)
	if err1 != nil {
		//fmt.Println("HttpGet err=",err)
		err = err1
		return
	}
	//取标题

	re1 := regexp.MustCompile(`<h1>(?s:(.*?))</h1>`)
	if re1 == nil {
		err = fmt.Errorf("%s", "regexp.MustCompile err")
		return
	}

	//取关键信息
	temTitle := re1.FindAllStringSubmatch(result, 1) //最后一个参数为1,表示只过滤第一个
	for _, data := range temTitle {
		title = data[1]
		title = strings.Replace(title, "\t", "", -1)
		break
	}
	//取内容

	re2 := regexp.MustCompile(`<div class="content-txt pt10">(?s:(.*?))<a id="prev" href="`)
	if re2 == nil {
		err = fmt.Errorf("%s", "regexp.MustCompile err")
		return
	}
	//取关键信息
	temContent := re2.FindAllStringSubmatch(result, 1) //最后一个参数为1,表示只过滤第一个
	for _, data := range temContent {
		content = data[1]
		content = strings.Replace(content, "\t", "", -1)
		break
	}
	return
}

func SpiderPage(i int, page chan<- int) {
	url := "https://www.pengfue.com/index_" + strconv.Itoa(i) + ".html" //将整型转化成字符串
	fmt.Printf("正在爬去第%d个网页,地址为%s:\n", i, url)
	//SpiderPage()
	//将所有网页的内容爬下来
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("SpiderPage HttpGet err=", err)
		return
	}
	//fmt.Println("####################")
	//取内容
	//<h1 class="dp-b"><a href="https://www.pengfue.com/content_1856812_1.html"
	//正则表达式
	re := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))"`)
	if re == nil {
		fmt.Println("regexp.MustCompile err")
		return
	}
	//取关键信息
	joyUrls := re.FindAllStringSubmatch(result, -1) //-1表示过滤所有

	fileTitle := make([]string, 0)
	fileContent := make([]string, 0)

	//以下代码在爬取单个笑话的时候为非并发的
	/*
		for _, data := range joyUrls {
			title, content, err := SpiderOneJoy(data[1])
			if err != nil {
				fmt.Println("SpiderOneJoy err=")
				return
			}
			fileTitle = append(fileTitle, title)
			fileContent = append(fileContent, content)
		}
	*/

	//以下代码在爬取单个笑话的时候为并发的
	fileTitlechan := make(chan string)
	fileContentchan := make(chan string)

	for _, data := range joyUrls {
		//爬取每一个笑话
		//这里注意一定要将data作为参数传给协程,否则的话其他协程有可能修改data的值
		go func(data1 []string) {
			title, content, err := SpiderOneJoy(data1[1])
			if err != nil {
				fmt.Println("SpiderOneJoy err=")
				return
			}
			//	fmt.Println(title)
			fileTitlechan <- title
			fileContentchan <- content

		}(data)
		//fmt.Printf("title=#%v#\n",title)
		//fmt.Printf("content=#%v#",content)
	}
	for range joyUrls {
		fileTitle = append(fileTitle, <-fileTitlechan)
		fileContent = append(fileContent, <-fileContentchan)
		//fmt.Println(fileTitle[i])
	}
	StoreJoyFile(i, fileTitle, fileContent)
	page <- i

}

///写入文件

func StoreJoyFile(i int, fileTitle, fileContent []string) {

	//把内容写入文件,以页码作为文件名
	fileName := "./testPaChong/" + strconv.Itoa(i) + ".html"
	f, err1 := os.Create(fileName)
	if err1 != nil {
		fmt.Println("os.Create err1", err1)
		return
	}
	n := len(fileTitle)
	for i := 0; i < n; i++ {
		f.WriteString(fileTitle[i] + "\n")                                                                        //写内容
		f.WriteString(fileContent[i] + "\n")                                                                      //写内容
		f.WriteString("\n####################################################################################\n") //写内容
	}
}

func DoWork(start, end int) {
	fmt.Printf("正在爬取%d到%d的页面\n", start, end)
	page := make(chan int)
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}

	for i := start; i <= end; i++ {
		fmt.Printf("第%d个页面爬取完成\n", <-page)
	}

}

func main() {
	var start, end int
	fmt.Printf("请输入起始页(>=1):")
	fmt.Scan(&start)
	fmt.Printf("请输入终止页(>=起始页):")
	fmt.Scan(&end)
	DoWork(start, end)
}
