package conf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	//"strings"
)

// IConf 组合模式，父类 定义一个 接口
type IConf interface {
	Save()
	Load()
}

// BaseConf 配置文件操作的基础封装
type BaseConf struct {
	conf IConf
}

// Save ...
func (cfg *BaseConf) Save() {
	log.Println("Save ENTER")
	t := reflect.TypeOf(cfg.conf).Elem().Name()
	fname := fmt.Sprintf("%s.json", t)
	//fname = strings.Replace(fname, "*", "r", -1)
	//log.Println(t)
	log.Println(fname)
	//fmt.Println(cfg.Conf)
	//fmt.Printf("type of is:%v\n", t)

	jsond, err := json.Marshal(cfg.conf)
	if err != nil {
		log.Println("生成json字符串错误")
	}
	SaveFile(fname, jsond)
	log.Println(cfg.conf)

}

// Load ...
func (cfg *BaseConf) Load() {
	log.Println("Load ENTER")
	t := reflect.TypeOf(cfg.conf).Elem().Name()
	fname := fmt.Sprintf("%s.json", t)
	log.Println(fname)
	//fname = strings.Replace(fname, "*", "r", -1)
	//fmt.Println(cfg.conf)
	data, err := ReadFile(fname)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(data, cfg.conf)
	if err != nil {
		log.Println(err)
	}
	log.Println(cfg.conf)
}

// SaveFile 存文件
func SaveFile(filename string, data []byte) error {
	fi, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fi.Close()
	w := bufio.NewWriter(fi)
	w.WriteString(string(data))
	w.Flush()
	return nil
}

// ReadFile 读文件
func ReadFile(filename string) ([]byte, error) {
	fi, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer fi.Close()
	r := bufio.NewReader(fi)
	var data []byte
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			return nil, err
		}
		if 0 == n {
			break
		}
		data = append(data, buf[:n]...)
	}
	return data, nil

}

/*
var sysCfg SysCfg

func main() {

	fmt.Println("Hello Go")
	//sysCfg.BaseConf.Conf = &sysCfg
	//sysCfg.Load()
	//fmt.Println(sysCfg)
	sysCfg.BaseConf.Conf = &sysCfg
	sysCfg.Name = "yang"
	sysCfg.Age = 188
	sysCfg.Save()
	sysCfg.Load()
	fmt.Println(sysCfg)
	//sysCfg.Save(&sysCfg)

	// oth := OthCfg{}
	// oth.Conf = &oth
	// oth.IP = "192.168.1.97"
	// oth.Port = 5057

	// oth.Save()
}
*/
