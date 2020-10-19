package utils

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/bitly/go-simplejson"
)

func HttpPutJson(url, body string) error {
	response, err := httplib.Put(url).Header("Content-Type", "application/json").SetTimeout(10*time.Second, 10*time.Second).Body(body).Response()
	if err == nil {
		defer response.Body.Close()
		if response.StatusCode >= 300 || response.StatusCode < 200 {
			body, _ := ioutil.ReadAll(response.Body)
			err = errors.New(response.Status + "；" + string(body))
		}
	}
	return err
}

func HttpPostJson(url, body string) (*simplejson.Json, error) {
	response, err := httplib.Post(url).Header("Content-Type", "application/json").SetTimeout(10*time.Second, 10*time.Second).Body(body).Response()
	var sj *simplejson.Json
	if err == nil {
		defer response.Body.Close()
		if response.StatusCode >= 300 || response.StatusCode < 200 {
			body, _ := ioutil.ReadAll(response.Body)
			err = errors.New(response.Status + "；" + string(body))
		} else {
			bodyBytes, _ := ioutil.ReadAll(response.Body)
			sj, err = simplejson.NewJson(bodyBytes)
		}
	}

	return sj, err
}
