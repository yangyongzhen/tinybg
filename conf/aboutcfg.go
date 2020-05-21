package conf

import (
	"log"
)

// AbtCfg 配置文件1
type AbtCfg struct {
	BaseConf
	Name  string
	Jobs  string
	WX    string
	QQ    string
	Email string
}

// Oth 全局声明
var Abt AbtCfg

func init() {
	log.Println("AbtCfg init ENTER")
	Abt.BaseConf.conf = &Abt

}
