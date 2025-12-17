// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/websocket/common"
	"chatwiki/internal/app/websocket/define"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type WsMessage struct {
	openid  string
	message []byte
}

type QueryClient struct {
	openid string
	count  int
	over   chan struct{}
}

var (
	AllOpenidMap   = make(map[string]map[*Client]struct{})
	EventOpenChan  = make(chan *Client)
	EventCloseChan = make(chan *Client)
	EventPullChan  = make(chan *WsMessage, 1000)
	EventPushChan  = make(chan *WsMessage, 1000)
	EventQueryChan = make(chan *QueryClient)
	wsUpgrader     = websocket.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}
)

func Running() {
	for {
		select {
		case client := <-EventOpenChan:
			go common.OpenPush(client.openid, client.ip, client.stime)
			if _, ok := AllOpenidMap[client.openid]; !ok {
				AllOpenidMap[client.openid] = make(map[*Client]struct{})
			}
			AllOpenidMap[client.openid][client] = struct{}{}
		case client := <-EventCloseChan:
			client.etime = tool.Time2Int()
			if len(AllOpenidMap[client.openid]) > 0 {
				go common.ClosePush(client.openid, client.ip, client.stime, client.etime)
				delete(AllOpenidMap[client.openid], client)
				close(client.send)
			}
			if len(AllOpenidMap[client.openid]) == 0 {
				delete(AllOpenidMap, client.openid)
			}
		case message := <-EventPullChan:
			if define.IsDev {
				logs.Debug(`receive:%s/%s`, message.openid, message.message)
			}
		case message := <-EventPushChan:
			for _, openid := range strings.Split(message.openid, `,`) {
				for client := range AllOpenidMap[openid] {
					select {
					case client.send <- message.message:
					default:
					}
				}
			}
		case query := <-EventQueryChan:
			query.count = len(AllOpenidMap[query.openid])
			query.over <- struct{}{}
			close(query.over)
		}
	}
}

func InitWs() {
	go Running()
	http.HandleFunc(`/ping`, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
		fmt.Println(fmt.Sprintf(`[INFO] %s | 200 | %s | %s "%s"`, tool.Date(), r.Host, r.Method, r.RequestURI))
		return
	})
	http.HandleFunc(`/ws`, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf(`[INFO] %s | 200 | %s | %s "%s"`, tool.Date(), common.GetRequestIP(r), r.Method, r.RequestURI))
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			logs.Error(err.Error())
			return
		}
		openid := r.URL.Query().Get(`openid`)
		if len(openid) == 0 {
			_ = conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		client := &Client{openid: openid, ip: common.GetRequestIP(r), conn: conn, send: make(chan []byte, 256), stime: tool.Time2Int()}
		EventOpenChan <- client
		go client.PushMessage()
		go client.PullMessage()
	})
	http.HandleFunc(`/isOnLine`, func(w http.ResponseWriter, r *http.Request) {
		openid := r.URL.Query().Get(`openid`)
		if len(openid) == 0 {
			_, _ = w.Write([]byte(lib_web.FmtJson(nil, errors.New(`openid empty`))))
			return
		}
		query := QueryClient{openid: openid, over: make(chan struct{})}
		EventQueryChan <- &query //wait query
		<-query.over             //wait over
		_, _ = w.Write([]byte(lib_web.FmtJson(query.count, nil)))
		return
	})
	addr := fmt.Sprintf(`:%d`, cast.ToUint(define.Config.WebService[`port`]))
	if err := http.ListenAndServe(addr, nil); err != nil {
		logs.Error(err.Error())
		panic(err)
	}
}
