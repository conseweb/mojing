/*
Text:
<xml>
 <ToUserName><![CDATA[toUser]]></ToUserName>
 <FromUserName><![CDATA[fromUser]]></FromUserName>
 <CreateTime>1348831860</CreateTime>
 <MsgType><![CDATA[text]]></MsgType>
 <Content><![CDATA[this is a test]]></Content>
 <MsgId>1234567890123456</MsgId>
</xml>


*/

package main

import (
	"bitbucket.org/qiyi/godom"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// md5("mojing")
	TOKEN       = "2cfb8f0539915d7df8c9139c5564cf83"
	APPID       = "wx5c913a59135446c1"
	APPSECRET   = "c2a3b1b715daba03bb899c1851c749cd"
	TEST_APPID  = "wx91f60e2b363a36df"
	TEST_SECRET = "f20b0df0280e6b0817772f5ea618f3c5"
	Text        = "text"
	Location    = "location"
	Image       = "image"
	Link        = "link"
	Event       = "event"
	Music       = "music"
	News        = "news"

	getTokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET"
)

var ACCESS_TOKEN string = ""

type msgBase struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
}

var MENU string = `{
	"button":[
	{
	       "name":"魔镜",
	       "sub_button":[
	        {
	           "type":"click",
	           "name":"介绍",
	           "key":"m1_jieshao"
	        },
	        {
	           "type":"click",
	           "name":"计划",
	           "key":"m1_jihua"
	        },
	        {
	           "type":"click",
	           "name":"合作",
	           "key":"m1_hezuo"
	        },
	        {
	           "type":"click",
	           "name":"联系",
	           "key":"m1_lianxi"
	        }]
	  },
	  {
	       "name":"案例",
	       "sub_button":[
	        {
	           "type":"view",
	           "name":"Restaurant Picker",
	           "key":"m2_restaurant",
	           "url":"http://jindou.io/demos/restaurant_picker"
	        },
	        {
	           "type":"view",
	           "name":"租车",
	           "key":"m2_zuche",
	           "url":"http://m.zuche.com/html5/newversion/index.html"
	        }]
	  },
	  {
	       "name":"招聘",
	       "sub_button":[
	        {
	           "type":"click",
	           "name":"职位",
	           "key":"m3_zhiwei"
	        }]
	  }]
	}
`

func LoadMenu(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		//Do something
	}

	return string(content)
}

/*
<xml>
 <ToUserName><![CDATA[toUser]]></ToUserName>
 <FromUserName><![CDATA[fromUser]]></FromUserName>
 <CreateTime>1348831860</CreateTime>
 <MsgType><![CDATA[text]]></MsgType>
 <Content><![CDATA[this is a test]]></Content>
 <MsgId>1234567890123456</MsgId>
</xml>

Image:
<xml>
 <ToUserName><![CDATA[toUser]]></ToUserName>
 <FromUserName><![CDATA[fromUser]]></FromUserName>
 <CreateTime>1348831860</CreateTime>
 <MsgType><![CDATA[image]]></MsgType>
 <PicUrl><![CDATA[this is a url]]></PicUrl>
 <MediaId><![CDATA[media_id]]></MediaId>
 <MsgId>1234567890123456</MsgId>
</xml>

Voice:
<xml>
    <ToUserName><![CDATA[toUser]]></ToUserName>
    <FromUserName><![CDATA[fromUser]]></FromUserName>
    <CreateTime>1357290913</CreateTime>
    <MsgType><![CDATA[voice]]></MsgType>
    <MediaId><![CDATA[media_id]]></MediaId>
    <Format><![CDATA[Format]]></Format>
    <MsgId>1234567890123456</MsgId>
</xml>
*/

// 对字符串进行SHA1哈希
func sha1hash(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func checkSignature(req *http.Request) bool {
	// test url:
	// http://localhost:3000/wx/?signature=8b90bc413e1090141e1e1755c7febbe849236582&timestamp=3436&nonce=23435&echostr=testok
	// get 3 params
	req.ParseForm()
	sig := req.Form.Get("signature")
	ts := req.Form.Get("timestamp")
	nonce := req.Form.Get("nonce")
	// echostr := req.Form.Get("echostr")

	// fmt.Printf("1: %s, 2: %s, 3: %s\n", TOKEN, ts, nonce)
	// sort
	sa := []string{TOKEN, ts, nonce}
	// sort.Strings(sa)
	sort.Sort(sort.StringSlice(sa))
	// hash
	data := fmt.Sprintf("%s%s%s", sa[0], sa[1], sa[2])
	sekret := sha1hash(data)
	if sekret == sig {
		return true
	} else {
		return false
	}
}

func main() {
	logfile, err := os.OpenFile("test.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	defer logfile.Close()
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime)
	logger.Println("Started ...")

	m := martini.Classic()
	// m.Get("/", func() string {
	//  return "Hello world!"
	// })

	// render html templates from templates directory
	m.Use(render.Renderer(render.Options{
		Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "hello", "魔镜")
	})

	m.Get("/ip", func(r render.Render, req *http.Request) {
		// if checkSignature(req) == false {
		// 	return "not wechat message"
		// }
		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		proxy, _, _ := net.SplitHostPort(req.Header.Get("X-FORWARDED-FOR"))
		r.JSON(200, map[string]interface{}{"ip": ip, "proxy": proxy})
	})

	// This will set the Content-Type header to "application/json; charset=UTF-8"
	m.Get("/api", func(r render.Render) {
		r.JSON(200, map[string]interface{}{"hello": "world"})
	})

	m.Get("/wx/menu/create", func(req *http.Request) string {
		// if checkSignature(req) == false {
		// 	return "not wechat message"
		// }
		_, msg := CreateMenu(MENU)
		return msg + MENU
	})

	m.Get("/wx/menu/get", func(req *http.Request) string {
		// if checkSignature(req) == false {
		// 	return "not wechat message"
		// }
		_, msg := GetMenu()
		return msg
	})

	m.Get("/wx/**", func(req *http.Request) string {
		if checkSignature(req) == false {
			return "not wechat message"
		}
		echostr := req.Form.Get("echostr")
		return echostr
	})

	m.Post("/wx/**", func(req *http.Request) string {
		// if checkSignature(req) == false {
		//  return "not wechat message"
		// }
		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
			return "error."
		}

		fmt.Println(string(body))
		base, doc := parseMsgBase(string(body))
		if (base != nil) && (doc != nil) {
			if base.MsgType == Text {
				content := doc.GetElementsByTagName("Content").Item(0).FirstChild().NodeValue()
				fmt.Println(content)
				resptxt := "Hello, i received a text "
				resp := EncodeTextRespMsg(base.FromUserName, base.ToUserName, resptxt)
				d, _ := dom.ParseString(resp)
				xml := dom.ToXml(d)
				return xml
			} else if base.MsgType == Image {
				picUrl := getStrValueByFieldName(doc, "PicUrl")
				mediaId := getStrValueByFieldName(doc, "MediaId")

				getMediaFile(mediaId)

				resptxt := "Hello, i received a image " + picUrl
				resp := EncodeTextRespMsg(base.FromUserName, base.ToUserName, resptxt)
				d, _ := dom.ParseString(resp)
				xml := dom.ToXml(d)
				return xml
			} else if base.MsgType == Event {
				event := getStrValueByFieldName(doc, "Event")
				if event == "subscribe" {
					resptxt := "欢迎订阅魔镜订阅号"
					resp := EncodeTextRespMsg(base.FromUserName, base.ToUserName, resptxt)
					d, _ := dom.ParseString(resp)
					xml := dom.ToXml(d)
					return xml
				} else if event == "unsubscribe" {
					log.Println(base.FromUserName + " unsubscribed!")
				} else if event == "CLICK" {
					eventkey := getStrValueByFieldName(doc, "EventKey")
					if eventkey == "m1_jieshao" {
						// resptxt := ""
					}
					resptxt := fmt.Sprintf("Menu click event: %s", eventkey)
					resp := EncodeTextRespMsg(base.FromUserName, base.ToUserName, resptxt)
					d, _ := dom.ParseString(resp)
					xml := dom.ToXml(d)
					return xml
				}
			}

		}
		return "parse message failed."

	})

	m.Run()
	// log.Fatal(http.ListenAndServe(":3000", m))
}

func EncodeTextRespMsg(to string, from string, content string) string {
	var txtmsg string = `<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`

	result := fmt.Sprintf(txtmsg, to, from, time.Duration(time.Now().Unix()), content)
	return result
}

func getStrValueByFieldName(doc dom.Document, fn string) string {
	return doc.GetElementsByTagName(fn).Item(0).FirstChild().NodeValue()
}

func getIntValueByFieldName(doc dom.Document, fn string) int {
	val := doc.GetElementsByTagName(fn).Item(0).FirstChild().NodeValue()
	ival, _ := strconv.Atoi(val)
	return ival
}

func parseMsgBase(data string) (base *msgBase, doc dom.Document) {
	d, err := dom.ParseString(data)
	if err != nil {
		fmt.Println("Parse xml string failed.", err)
		return nil, nil
	}

	base = new(msgBase)
	base.ToUserName = getStrValueByFieldName(d, "ToUserName")
	base.FromUserName = getStrValueByFieldName(d, "FromUserName")
	seconds := getIntValueByFieldName(d, "CreateTime")
	base.CreateTime = time.Duration(seconds) * time.Second
	base.MsgType = d.GetElementsByTagName("MsgType").Item(0).FirstChild().NodeValue()
	return base, d
}

func GetAccessToken() string {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	x := fmt.Sprintf(url, APPID, APPSECRET)
	// http get
	res, err := http.Get(x)
	if err != nil {
		log.Fatal(err)
	}
	data, err2 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err2)
	}

	type Message struct {
		Access_token string
		Expires_in   int
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))

	var m Message
	if err := dec.Decode(&m); err != nil {
		log.Fatal(err)
	}

	return m.Access_token
}

// ftype: image, voice, video, thumb
func postFile(filename string, ftype string) error {
	tUrlTemp := "http://file.api.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s"
	targetUrl := fmt.Sprintf(tUrlTemp, ACCESS_TOKEN, ftype)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// 关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

func getMediaFile(mid string) {
	tUrlTemp := "http://file.api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
	targetUrl := fmt.Sprintf(tUrlTemp, ACCESS_TOKEN, mid)

	// http get
	res, err := http.Get(targetUrl)
	if err != nil {
		log.Fatal(err)
	}
	data, err2 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err2)
	}

	fn := "/home/ubuntu/public/images/" + mid + ".jpg"
	// save to file
	f, _ := os.Create(fn)
	defer f.Close()
	f.Write(data)
}

func CreateMenu(menu string) (int, string) {
	access_token := GetAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s"

	x := fmt.Sprintf(url, access_token)
	res, err := http.Post(x, "application/json", strings.NewReader(menu))
	if err != nil {
		log.Fatal(err)
	}
	data, err2 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err2)
	}

	type Message struct {
		Errcode int
		Errmsg  string
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))

	var m Message
	if err := dec.Decode(&m); err != nil {
		log.Fatal(err)
	}

	return m.Errcode, m.Errmsg
}

func GetMenu() (int, string) {
	access_token := GetAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s"

	x := fmt.Sprintf(url, access_token)
	res, err := http.Get(x)
	if err != nil {
		log.Fatal(err)
	}
	data, err2 := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err2)
	}

	type Message struct {
		Errcode int
		Errmsg  string
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))

	var m Message
	if err := dec.Decode(&m); err != nil {
		log.Fatal(err)
	}

	return m.Errcode, m.Errmsg
}
