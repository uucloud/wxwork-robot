package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/uucloud/wework-robot/utils"
	"github.com/uucloud/wework-robot/wxprotoc/receiver"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var M *Manager

func init() {
	M = NewManager()
}

type Manager struct {
	receiver.Receiver

	rLock  sync.RWMutex
	robots map[string]*Robot

	router *mux.Router

	ctx context.Context
}

func NewManager() *Manager {
	m := &Manager{
		Receiver: nil,
		rLock:    sync.RWMutex{},
		robots:   map[string]*Robot{},
		router:   mux.NewRouter(),
		ctx:      context.Background(),
	}
	m.init()
	return m
}

func (m *Manager) Handler() http.Handler {
	return m.router
}

func (m *Manager) AddRobot(name, sendKey, receiveToken, receiveEncodingAESKey string) (*Robot, error) {
	m.rLock.Lock()
	defer m.rLock.Unlock()
	_, e := m.robots[name]
	if e {
		return nil, fmt.Errorf("robot %s has exist", name)
	}
	m.robots[name] = newRobot(m.ctx, name, sendKey, receiveToken, receiveEncodingAESKey, time.Second*10)
	return m.robots[name], nil
}

func (m *Manager) init() {
	m.router.HandleFunc("/", m.callback)
}

func (m *Manager) callback(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.verify(w, r)
	case http.MethodPost:
		m.msg(w, r)
	default:
		utils.Logger.Error("callback method error", zap.String("method", r.Method))
	}
}

func (m *Manager) verify(w http.ResponseWriter, r *http.Request) {
	msgSign := parseURLQuery(r, "msg_signature")
	timestamp := parseURLQuery(r, "timestamp")
	nonce := parseURLQuery(r, "nonce")
	echostr := parseURLQuery(r, "echostr")

	for _, robot := range m.robots {
		echo, err := robot.Rec.Verify(msgSign, timestamp, nonce, echostr)
		if err == nil {
			_, _ = w.Write(echo)
			return
		}
	}
	utils.Logger.Error("unknown verify",
		zap.String("msgSign", msgSign),
		zap.String("timestamp", timestamp),
		zap.String("nonce", nonce),
		zap.String("echo", echostr),
	)

	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("error"))
}

func (m *Manager) msg(w http.ResponseWriter, r *http.Request) {
	msgSign := parseURLQuery(r, "msg_signature")
	timestamp := parseURLQuery(r, "timestamp")
	nonce := parseURLQuery(r, "nonce")

	utils.Logger.Info("params", zap.String("timestamp", timestamp), zap.String("nonce", nonce), zap.String("msgSign", msgSign))

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	}
	defer r.Body.Close()

	for i, robot := range m.robots {
		d, dErr := robot.Rec.DecryptMsg(msgSign, timestamp, nonce, data)
		utils.Logger.Info("decryptMsg", zap.String("receive msg", string(d)))
		if dErr == nil {
			var event receiver.CallEvent
			err = xml.Unmarshal(d, &event)
			if err != nil {
				zap.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			rbt := m.robots[i]

			b, err := rbt.Rec.Reply(&event)
			if err != nil {
				utils.Logger.Error("reply error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			_, _ = w.Write(b)
			return
		}
	}

	utils.Logger.Error("robot not found",
		zap.String("msgSign", msgSign),
		zap.String("timestamp", timestamp),
		zap.String("nonce", nonce),
	)

	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("robot not found"))
}

func parseURLQuery(r *http.Request, key string) string {
	v := r.URL.Query()[key]
	if len(v) > 0 {
		return v[0]
	}
	return ""
}
