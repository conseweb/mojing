/*
新浪微博操作接口，需要先OAuth登录
*/
package sina

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	SINA_WEIBO_API_URL = "https://api.weibo.com/2"
)

/*
Sina API client
*/
type SaeClient struct {
	clientId     string
	clientSecret string
	accessToken  string
	refreshToken string
}

func NewSaeClient(clientID, clientSecret, accessToken, refreshToken string) *SaeClient {
	return &SaeClient{clientID, clientSecret, accessToken, refreshToken}
}

/*
post a sina microblog
*/
func (c *SaeClient) Update(status, lat, long, annotations string) error {
	v := url.Values{}
	v.Add("access_token", c.accessToken)
	v.Add("status", status)
	if lat == "" {
		lat = "0.0"
	}
	if long == "" {
		long = "0.0"
	}
	v.Add("lat", lat)
	v.Add("long", long)
	v.Add("annotations", annotations)

	res, err := http.PostForm(SINA_WEIBO_API_URL+"/statuses/update.json", v)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var f interface{}
	json.NewDecoder(res.Body).Decode(&f)
	m := f.(map[string]interface{})
	if _, ok := m["error"]; ok {
		return errors.New(m["error"].(string))
	}

	return nil

}

func (c *SaeClient) Upload(status string, pic_path, lat, long string) error {
	var b bytes.Buffer
	formdata := multipart.NewWriter(&b)
	formdata.WriteField("access_token", c.accessToken)
	formdata.WriteField("status", status)
	picdata, _ := formdata.CreateFormFile("pic", pic_path)
	if strings.HasPrefix(pic_path, "http") {
		res, err := http.Get(pic_path)
		if err != nil {
			return errors.New("pic load error")
		}
		io.Copy(picdata, res.Body)
		res.Body.Close()
	} else {
		fh, _ := os.Open(pic_path)
		io.Copy(picdata, fh)
		fh.Close()
	}
	form_type := formdata.FormDataContentType()

	formdata.Close()
	if lat == "" {
		lat = "0.0"
	}
	if long == "" {
		long = "0.0"
	}
	formdata.WriteField("lat", lat)
	formdata.WriteField("long", long)

	res, err := http.Post(SINA_WEIBO_API_URL+"/statuses/upload.json", form_type, &b)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var j interface{}

	json.NewDecoder(res.Body).Decode(&j)

	m := j.(map[string]interface{})

	if _, ok := m["error"]; ok {
		return errors.New(m["error"].(string))
	}
	return nil
}
