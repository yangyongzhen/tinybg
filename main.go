package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"tinybg/conf"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID      string //文章唯一ID 取 filename的md5
	Title   string
	Date    string
	Summary string
	Body    string
	File    string
	ImgFile string // 文章前的图片展示
	Item    string //分类
	Author  string //作者

	Cmts    []conf.Comment //评论
	CmtCnt  int            //评论数量
	VistCnt int            //浏览量
}

// Notice 跳转提示
type Notice struct {
	Mess    string
	IsSucc  bool
	TimeOut int
	Href    string
}

var articlesMap map[string]string /*创建集合 */

// NewArticles 最新文章按日期排序
type NewArticles []conf.Articles

// HotArticles 热门文章按访问量排序
type HotArticles []conf.Articles

//  NewArts  最新文章
var NewArts NewArticles

// HotArts 热门文章
var HotArts HotArticles

// NewPosts 总的文章，按时间排过序的
var NewPosts NewArticles

// NewCmts 最新评论
var NewCmts []conf.Comment

//文章排序的实现
//文章排序
//Len()
func (a NewArticles) Len() int {
	return len(a)
}

//Less():将有低到高排序
func (a NewArticles) Less(i, j int) bool {
	//fmt.Println(a[i].Date)
	//fmt.Println(a[j].Date)
	return a[i].Date > a[j].Date
}

//Swap()
func (a NewArticles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

//Len()
func (a HotArticles) Len() int {
	return len(a)
}

//Less():将有低到高排序
func (a HotArticles) Less(i, j int) bool {
	//fmt.Println(a[i].Date)
	//fmt.Println(a[j].Date)
	return a[i].VistCnt > a[j].VistCnt
}

//Swap()
func (a HotArticles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// md5Str 计算MD5值
func md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// RefreshData 后台异步更新数据
func RefreshData() {
	//全局的最新评论，只显示最新的三条
	curcount := len(conf.Cmt.Comts) //总的评论数量
	NewCmts = conf.Cmt.Comts
	if curcount >= 3 {
		NewCmts = conf.Cmt.Comts[(curcount - 3):curcount] //最新的3条评论
	}
	//post := conf.Art.ArticlesMap[conf.Item.Items[id]]
	NewPosts = NewArticles{}
	HotArts = HotArticles{}
	for key, value := range conf.Art.ArticlesMap {
		fmt.Println(key)
		for _, value1 := range value {
			NewPosts = append(NewPosts, value1)
			HotArts = append(HotArts, value1)
		}
	}
	sort.Sort(NewPosts) //日期排序
	sort.Sort(HotArts)  //访问量排序
	//fmt.Println("IS Sorted?", sort.IsSorted(post))
	NewArts = NewPosts
	num := len(NewPosts)
	if num > 9 {
		NewArts = NewPosts[0:9]
		HotArts = HotArts[0:9]
	}
	conf.Art.Save()
	conf.Stat.Save()
}

func handleIndex(c *gin.Context) {
	fmt.Printf("%#v\n", NewPosts)
	spage := c.DefaultQuery("page", "1")
	fmt.Printf("cur page:%s\n", spage)
	page, _ := strconv.Atoi(spage)
	nums := len(NewPosts)
	allpage := nums / 5
	if nums%5 != 0 {
		allpage = nums/5 + 1
	}
	fmt.Printf("all page num:%d\n", allpage)
	posts := NewPosts
	if (page * 5) < nums {
		posts = NewPosts[(page-1)*5 : page*5]
	} else {
		posts = NewPosts[(page-1)*5 : nums]
	}
	tabs := make([]int, allpage+2) //分页表
	if (page - 1) == 0 {
		tabs[0] = 1
	} else {
		tabs[0] = page - 1
	}
	for i := 1; i <= allpage; i++ {
		tabs[i] = i
	}
	if page+1 <= allpage {
		tabs[allpage+1] = page + 1
	} else {
		tabs[allpage+1] = 1
	}
	fmt.Printf("tabs:%#v\n", tabs)
	conf.Stat.ToCnt++
	ip := c.ClientIP()
	fmt.Printf("client ip:%s\n", ip)
	conf.Stat.IPs = append(conf.Stat.IPs, ip)
	c.HTML(http.StatusOK, "index.html", gin.H{"post": posts, "items": conf.Item.Items, "about": conf.Abt, "newcmts": NewCmts, "newart": NewArts, "hotart": HotArts, "vistcnt": conf.Stat.ToCnt, "curpage": page, "tabs": tabs})
}

// strTrip 去除字符串的空格和换行
func strTrip(src string) string {
	str := strings.Replace(src, " ", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return str
}

func getPosts() []Post {
	a := []Post{}
	files, _ := filepath.Glob("posts/*")
	for i, f := range files {
		fmt.Println(i)
		fmt.Println(f)

		file := strings.Replace(f, "posts\\", "", -1)
		fmt.Println(file)
		file = strings.Replace(file, ".md", "", -1)
		fileread, _ := ioutil.ReadFile(f)
		lines := strings.Split(string(fileread), "\n")
		title := string(lines[0])
		date := strTrip(string(lines[1]))
		summary := string(lines[2])
		imgfile := string(lines[3])
		item := strTrip(string(lines[4]))
		author := string(lines[5])

		imgfile = strTrip(imgfile)
		body := ""
		//body := strings.Join(lines[4:len(lines)], "\n")
		fmt.Println(imgfile)
		id := md5Str(file)
		articlesMap[id] = f
		itemcount := len(conf.Item.Items)
		//查找分类看是否已存在
		if itemcount == 0 {
			conf.Item.Items = append(conf.Item.Items, item)
		} else {
			k := 0
			for k = 0; k < itemcount; k++ {
				if conf.Item.Items[k] == item {
					break
				}
			}
			if k >= itemcount {
				//分类之前未存在，添加分类
				conf.Item.Items = append(conf.Item.Items, item)
				itemcount = len(conf.Item.Items)
			}
		}

		for j := 0; j < itemcount; j++ {
			if conf.Item.Items[j] == item {
				art := conf.Articles{id, item, title, date, summary, file, imgfile, author, 0, 0}
				if conf.Art.ArticlesMap[item] == nil {
					conf.Art.ArticlesMap[item] = make(map[string]conf.Articles)
				}
				_, exists := conf.Art.ArticlesMap[item][id]
				if !exists {
					conf.Art.ArticlesMap[item][id] = art
				}
				break
			}
		}

		a = append(a, Post{id, title, date, summary, body, file, imgfile, item, author, nil, conf.Art.ArticlesMap[item][id].CmtCnt, conf.Art.ArticlesMap[item][id].VistCnt})
	}
	conf.Art.Save()
	conf.Item.Save()
	fmt.Printf("%#v\n", conf.Art.ArticlesMap)
	fmt.Printf("%#v\n", articlesMap)
	return a
}

func handleArticles(c *gin.Context) {
	id := c.Param("id")
	file := articlesMap[id]
	fmt.Println(id)
	fmt.Println(file)
	fileread, _ := ioutil.ReadFile(file)
	lines := strings.Split(string(fileread), "\n")
	title := string(lines[0])
	date := string(lines[1])
	summary := string(lines[2])
	imgfile := string(lines[3])
	item := strTrip(string(lines[4]))
	author := string(lines[5])
	imgfile = strTrip(imgfile)
	body := strings.Join(lines[5:len(lines)], "\n")
	//fmt.Println(body)
	//fmt.Printf("%#v\n", conf.Cmt.Comts)
	cmts := []conf.Comment{}
	count := 0 //评论数量
	for i, v := range conf.Cmt.Comts {
		//fmt.Printf("%#v\n", c)
		if v.ID == id {
			cmts = append(cmts, conf.Cmt.Comts[i])
			count++
		}
	}
	//文章浏览量++
	art := conf.Art.ArticlesMap[item][id]
	art.VistCnt++ //浏览数量加一
	conf.Art.ArticlesMap[item][id] = art

	go RefreshData()

	p := Post{id, title, date, summary, body, file, imgfile, item, author, cmts, conf.Art.ArticlesMap[item][id].CmtCnt, conf.Art.ArticlesMap[item][id].VistCnt}
	c.HTML(http.StatusOK, "article.html", gin.H{"post": p, "items": conf.Item.Items, "cmtcounts": count, "newcmts": NewCmts, "newart": NewArts, "hotart": HotArts, "vistcnt": conf.Stat.ToCnt})
}

func handleItems(c *gin.Context) {
	sid := c.Param("id")
	id, _ := strconv.Atoi(sid)
	all := conf.Art.ArticlesMap[conf.Item.Items[id]]
	fmt.Printf("%#v\n", all)
	posts := []conf.Articles{}
	// 遍历 hash
	for _, value := range all {
		posts = append(posts, value)
	}
	//支持分页
	spage := c.DefaultQuery("page", "1")
	fmt.Printf("cur page:%s\n", spage)
	page, _ := strconv.Atoi(spage)
	nums := len(posts)
	allpage := nums / 5
	if nums%5 != 0 {
		allpage = nums/5 + 1
	}
	fmt.Printf("all page num:%d\n", allpage)
	curposts := posts
	if (page * 5) < nums {
		curposts = posts[(page-1)*5 : page*5]
	} else {
		curposts = posts[(page-1)*5 : nums]
	}
	tabs := make([]int, allpage+2) //分页表
	if (page - 1) == 0 {
		tabs[0] = 1
	} else {
		tabs[0] = page - 1
	}
	for i := 1; i <= allpage; i++ {
		tabs[i] = i
	}
	if page+1 <= allpage {
		tabs[allpage+1] = page + 1
	} else {
		tabs[allpage+1] = 1
	}
	fmt.Printf("tabs:%#v\n", tabs)
	c.HTML(http.StatusOK, "items.html", gin.H{"post": curposts, "items": conf.Item.Items, "newcmts": NewCmts, "newart": NewArts, "hotart": HotArts, "vistcnt": conf.Stat.ToCnt, "curitem": id, "curpage": page, "tabs": tabs})
}

func handlePostComment(c *gin.Context) {

	sid := c.PostForm("id")    //文章ID
	item := c.PostForm("item") //文章分类
	item = strTrip(item)
	title := c.PostForm("title")   //文章标题
	text := c.PostForm("text")     //评论内容
	author := c.PostForm("author") //评论者
	email := c.PostForm("email")   //邮箱
	url := c.PostForm("url")       //评论者网址
	time := time.Now().Format("2006-01-02 15:04:05")

	fmt.Println(sid, item, text, author, email, url, time)

	href := "/article/" + sid
	notice := Notice{"提交成功", true, 3, href}
	art := conf.Art.ArticlesMap[item][sid]
	art.CmtCnt++ //评论数量加一
	conf.Art.ArticlesMap[item][sid] = art
	//fmt.Printf("%#v\n", conf.Art.ArticlesMap[item][sid])
	fmt.Println(conf.Art.ArticlesMap[item][sid].CmtCnt)
	//fmt.Printf("%#v\n", conf.Art.ArticlesMap[item])
	cmt := conf.Comment{sid, item, title, author, email, text, time, url}
	fmt.Printf("%#v\n", cmt)
	conf.Cmt.Comts = append(conf.Cmt.Comts, cmt)
	conf.Cmt.Save()
	//conf.Art.Save()

	go RefreshData()

	c.HTML(http.StatusOK, "success.html", gin.H{"notice": notice})
}

// Add ...模板里使用加
func Add(a, b int) int {
	return a + b
}

// Dec ...模板里使用减
func Dec(a, b int) int {
	return a - b
}
func main() {
	router := gin.Default()
	//增加几个模板自定义函数Add Dec
	router.SetFuncMap(template.FuncMap{
		"add": Add,
		"dec": Dec,
	})

	//关于
	conf.Abt.Name = "一米阳光"
	conf.Abt.Jobs = "嵌入式 linux | Android | go web"
	conf.Abt.WX = "yongzhen1111"
	conf.Abt.QQ = "534117529"
	conf.Abt.Email = "534117529@qq.com"
	conf.Abt.Save()
	//分类
	// conf.Item.Items = append(conf.Item.Items, "随笔")
	// fmt.Printf("%d\n", len(conf.Item.Items))
	//加载分类
	conf.Item.Load()

	//加载评论
	conf.Cmt.Load()

	//加载文章
	conf.Art.Load()
	//加载统计信息
	conf.Stat.Load()

	articlesMap = make(map[string]string)
	articleIndex := getPosts()
	fmt.Printf("%#v\n", articleIndex)

	//更新文章排序和浏览量
	RefreshData()
	//静态文件
	router.Static("/assets", "./static")
	//渲染html页面
	router.LoadHTMLGlob("views/*")
	router.GET("/", handleIndex)
	router.GET("/article/:id", handleArticles)
	router.GET("/items/:id", handleItems)
	router.POST("/comment", handlePostComment)
	//运行的端口
	router.Run(":8000")
}
