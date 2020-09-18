package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//Servercfg 控制伺服器的設定
type Servercfg struct {
	ForwordName        string `json:"ForwordName"`
	ForwordServerIP    string `json:"ForwordServerIP"`
	ForwordCtrlAPIPort string `json:"ForwordCtrlAPIPort"`
	ListenAPIPort      string `json:"ListenAPIPort"`
	AESKey             string `json:"AESKey"`
}

var cfg []Servercfg
var timetag int64

func gettime() {
	for {
		timetag = time.Now().Unix()
		//fmt.Println(timetag)
		time.Sleep(100 * time.Millisecond)
	}

}

//GetConfig 讀取配置檔
func GetConfig() {
	tmp, err := ioutil.ReadFile("./checkinit.conf")
	if err != nil {
		fmt.Println("GetConfig error")
	}
	config := string(tmp)
	err = json.Unmarshal([]byte(config), &cfg)
	if err != nil {
		fmt.Println(err)
		fmt.Println("cfg to json error")
	}
	fmt.Println(cfg)
}

//StartServer 開始監聽轉發
func StartServer(configdata Servercfg) {
	r := gin.Default() // 使用默认中间件
	//（logger和recovery）
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
			"message": "pong",
		})
	})
	r.GET("/check", CheckClientFunc(configdata))
	r.Run(":" + configdata.ListenAPIPort) // 启动服务，并默认监听8080端口
}

func main() {
	//CtrlIPTable("192.168.10.10", "2000")
	go gettime()
	GetConfig()
	for i := 0; i < len(cfg); i++ {
		go StartServer(cfg[i])
	}
	for {
	}
}

//CheckClientFunc 確認config
func CheckClientFunc(configdata Servercfg) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		data := c.Query("data")
		//fmt.Println(data)
		keystring := configdata.AESKey
		data = AESDecode(data, keystring)
		dataspl := strings.Split(data, "|")
		if len(dataspl) != 2 {

			c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
				"message": "err",
			})
			return
		}
		ip := dataspl[1]
		inttime := dataspl[0]
		reqip := c.ClientIP()
		fmt.Println(ip)
		fmt.Println(inttime)
		fmt.Println(reqip)
		timeint, err := strconv.ParseInt(inttime, 10, 64)
		if err != nil {
			resp, err := http.Get("http://" + configdata.ForwordServerIP + ":" + configdata.ForwordCtrlAPIPort + "/open?ip=" + ip)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
				"message": "err",
			})
			return
		}
		//開始檢驗來源跟取得ip
		if reqip == ip {
			//開始檢驗時間
			if timetag-timeint < 30 {
				c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
					"message": "OK",
				})
				return
			}
			c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
				"message": "err",
			})
			return

		}

		c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
			"message": "err",
		})
		return
	}
	return gin.HandlerFunc(fn)
}

//CheckClient 確認客戶端是否合法
func CheckClient(c *gin.Context, index int) {
	//fmt.Println(index)
	data := c.Query("data")
	//fmt.Println(data)
	keystring := "962EE76B443BC11BBD7A4800DDEACA43"
	data = AESDecode(data, keystring)
	dataspl := strings.Split(data, "|")
	if len(dataspl) != 2 {

		c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
			"message": "err",
		})
		return
	}
	ip := dataspl[1]
	inttime := dataspl[0]
	reqip := c.ClientIP()
	fmt.Println(ip)
	fmt.Println(inttime)
	fmt.Println(reqip)
	timeint, err := strconv.ParseInt(inttime, 10, 64)
	if err != nil {

		c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
			"message": "err",
		})
		return
	}
	//開始檢驗來源跟取得ip
	if reqip == ip {
		//開始檢驗時間
		if timetag-timeint < 30 {

			c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
				"message": "OK",
			})
			return
		}
		c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
			"message": "err",
		})
		return

	}

	c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
		"message": "err",
	})
	return

}

//RunCmd 運行指令
func RunCmd(cmdstr string) {
	cmd := exec.Command("/bin/sh", "-c", cmdstr) ///檢視當前目錄下檔案
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
