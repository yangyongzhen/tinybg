package conf

import (
	"log"
)

type Articles struct {
	ID      string
	Item    string //分类
	Title   string
	Date    string
	Summary string
	File    string
	ImgFile string // 文章前的图片展示
	Author  string
	CmtCnt  int //评论数量
	VistCnt int //浏览数量
}

// ArtCfg 配置文件2
type ArtCfg struct {
	BaseConf
	Ver         int
	ArticlesMap map[string]map[string]Articles
}

// Art 全局声明
var Art ArtCfg

func init() {
	log.Println("ArtCfg init ENTER")
	Art.BaseConf.conf = &Art
	Art.ArticlesMap = make(map[string]map[string]Articles)
}
