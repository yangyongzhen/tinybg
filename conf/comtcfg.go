package conf

import (
	"log"
)

// Comment 文章评论内容
type Comment struct {
	ID    string //对应的文章ID
	Item  string //对应的文章分类
	Title string //对应的文章标题
	User  string //用户名
	Email string //邮箱
	Text  string //评论内容
	Times string //评论时间
	URL   string //用户的网址

}

// ComtCfg 评论管理
type ComtCfg struct {
	BaseConf
	Comts []Comment
}

// Cmt 全局声明
var Cmt ComtCfg

func init() {
	log.Println("ComtCfg init ENTER")
	Cmt.BaseConf.conf = &Cmt
	Cmt.Comts = make([]Comment, 0)

}
