package alarm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jesseincn/utils/log"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const ding_url = "https://oapi.dingtalk.com/robot/send?access_token="

type Ding struct {
	lock    sync.Mutex
	token   string               // 钉钉告警群token
	timeMap map[string]time.Time // 时间字典，记录各个模块最近发送警告的时间
	span    int                  // 同一个模块发送钉钉消息的时间间隔，单位为秒
}

func NewDing(token string, span int) *Ding {
	return &Ding{
		token:   token,
		timeMap: make(map[string]time.Time),
		span:    span,
	}
}

func (d *Ding) Send(module, msg string, phone ...string) error {
	if module == "" {
		return errors.New("module is empty")
	}
	var can bool
	d.lock.Lock()
	v, ok := d.timeMap[module]
	d.lock.Unlock()
	if ok {
		if time.Now().Sub(v).Seconds() > float64(d.span) {
			can = true
		} else {
			log.Logger.Info("this msg is passed", zap.String("module", module), zap.String("msg", msg))
			return nil
		}
	} else {
		can = true
	}
	if can {
		if strings.Contains(module, "【") || strings.Contains(module, "[") {
			msg = fmt.Sprintf("%s%s", module, msg)
		}
		err := d.send(msg, phone, d.token)
		if err != nil {
			log.Logger.Error("d.send", zap.Error(err))
			return err
		}
		d.lock.Lock()
		d.timeMap[module] = time.Now()
		d.lock.Unlock()
	}
	return nil
}

func (d *Ding) send(msg string, phones []string, token string) error {
	ms := Message{
		MsgType: "text",
		Text: Texts{
			Content: msg,
		},
		At: Ats{
			AtMobiles: phones,
			IsAtAll:   false,
		},
	}
	para, _ := json.Marshal(ms)
	req, err := http.NewRequest("POST", ding_url+token, bytes.NewReader([]byte(para)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	var tpt = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tpt,
		Timeout:   time.Duration(10000) * time.Millisecond,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r := &Rsp{}
	_ = json.Unmarshal(data, r)
	if r.Errcode != 0 {
		return errors.New(r.Errmsg)
	}
	return nil
}

type Texts struct {
	Content string `json:"content"`
}

type Ats struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}
type Message struct {
	MsgType string `json:"msgtype"`
	Text    Texts  `json:"text"`
	At      Ats    `json:"at"`
}

type Rsp struct {
	Errmsg  string `json:"errmsg"`
	Errcode int    `json:"errcode"`
}
