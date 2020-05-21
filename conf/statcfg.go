package conf

import (
	"log"
)

// LvMsg 留言
type LvMsg struct {
	User  string //用户名
	Email string //邮箱
	Text  string //评论内容
	Times string //评论时间
	URL   string //用户的网址
}

// StatCfg 统计 配置文件1
type StatCfg struct {
	BaseConf
	ToCnt int      //总访问量
	IPs   []string //访问者IP
	Msgs  []LvMsg
}

// Stat 全局声明
var Stat StatCfg

func init() {
	log.Println("StatCfg init ENTER")
	Stat.BaseConf.conf = &Stat

	Stat.Msgs = make([]LvMsg, 0)
	Stat.IPs = make([]string, 0)
}
