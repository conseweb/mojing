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
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	// md5("mojing")
	TOKEN     = "2cfb8f0539915d7df8c9139c5564cf83"
	APPID     = "wx5b0dcca6246fad3e"
	APPSECRET = "5d9afef1608641a63137e7094e2b8a36"
	Text      = "text"
	Location  = "location"
	Image     = "image"
	Link      = "link"
	Event     = "event"
	Music     = "music"
	News      = "news"

	getTokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET"
)

type msgBase struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
}

type Request struct {
	XMLName                xml.Name `xml:"xml"`
	msgBase                         // base struct
	Location_X, Location_Y float32
	Scale                  int
	Label                  string
	PicUrl                 string
	MsgId                  int
}

type Response struct {
	XMLName xml.Name `xml:"xml"`
	msgBase
	ArticleCount int     `xml:",omitempty"`
	Articles     []*item `xml:"Articles>item,omitempty"`
	FuncFlag     int
}

type item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}

type MsgHeader struct {
	To      string `xml: "ToUserName"`
	From    string `xml: "FromUserName"`
	Time    int    `xml: "CreateTime"`
	MsgType string `xml: "MsgType"`
	Id      string `xml: "MsgId"`
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
*/
type TextMessage struct {
	MsgHeader
	Content string `xml: "Content"`
}

/*
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
*/
type ImageMessage struct {
	MsgHeader
	Url     string `xml: "PicUrl"`
	MediaId string `xml: "MediaId"`
}

/*
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
type VoiceMessage struct {
	MsgHeader
	MediaId string `xml: "MediaId"`
	Format  string `xml: "Format"`
}

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
	// if len(req.Form) > 0 {
	// 	for k, v := range req.Form {
	// 		x := fmt.Sprintf("%s, %s", k, v[0])
	// 		fmt.Println(x)
	// 	}
	// }
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

func parseMsgBase(data string) (base *msgBase, doc dom.Document) {
	d, err := dom.ParseString(data)
	if err != nil {
		fmt.Println("Parse xml string failed.", err)
		return nil, nil
	}

	base = new(msgBase)
	base.ToUserName = d.GetElementsByTagName("ToUserName").Item(0).FirstChild().NodeValue()
	base.FromUserName = d.GetElementsByTagName("FromUserName").Item(0).FirstChild().NodeValue()
	seconds, _ := strconv.Atoi(d.GetElementsByTagName("CreateTime").Item(0).FirstChild().NodeValue())
	base.CreateTime = time.Duration(seconds) * time.Second
	base.MsgType = d.GetElementsByTagName("MsgType").Item(0).FirstChild().NodeValue()
	return base, d

	// fmt.Println(d.GetElementsByTagName("ToUserName").Item(0).NodeName())
	// fmt.Println(d.GetElementsByTagName("ToUserName").Item(0).FirstChild().NodeValue())
	// fmt.Println(d.GetElementsByTagName("FromUserName").Item(0).FirstChild().NodeValue())
	// fmt.Println(d.GetElementsByTagName("CreateTime").Item(0).FirstChild().NodeValue())
	// fmt.Println(root.ChildNodes().Item(2).NodeValue())
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
	// 	return "Hello world!"
	// })

	// render html templates from templates directory
	m.Use(render.Renderer(render.Options{
		Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "hello", "魔镜")
	})

	// This will set the Content-Type header to "application/json; charset=UTF-8"
	m.Get("/api", func(r render.Render) {
		r.JSON(200, map[string]interface{}{"hello": "world"})
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
		// 	return "not wechat message"
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
				resptxt := "Hello, i received a text from " + base.FromUserName
				resp := EncodeTextRespMsg(base.FromUserName, base.ToUserName, resptxt)
				d, _ := dom.ParseString(resp)
				xml := dom.ToXml(d)
				return xml
			} else if base.MsgType == Image {
				resptxt := "Hello, i received a image from " + base.FromUserName
				resp := EncodeTextRespMsg(base.FromUserName, base.ToUserName, resptxt)
				d, _ := dom.ParseString(resp)
				xml := dom.ToXml(d)
				return xml
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

// func DecodeRequest(data []byte) (req *Request, err error) {
// 	req = &Request{}
// 	if err = xml.Unmarshal(data, req); err != nil {
// 		return
// 	}
// 	req.CreateTime *= time.Second
// 	return
// }

// func NewResponse() (resp *Response) {
// 	resp = &Response{}
// 	resp.CreateTime = time.Duration(time.Now().Unix())
// 	return
// }

// func (resp Response) Encode() (data []byte, err error) {
// 	resp.CreateTime = time.Duration(time.Now().Unix())
// 	data, err = xml.Marshal(resp)
// 	return
// }
