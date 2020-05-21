package conf

import (
	"log"
)

// ItemCfg 文章分类配置文件1
type ItemCfg struct {
	BaseConf
	Items []string //分类
}

// Item 全局声明
var Item ItemCfg

func init() {
	log.Println("ItemCfg init ENTER")
	Item.BaseConf.conf = &Item
	Item.Items = make([]string, 0)

}
