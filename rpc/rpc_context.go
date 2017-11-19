// Copyright 2017 bingo Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"github.com/snippetor/bingo/net"
	"strconv"
	"github.com/snippetor/bingo/log/fwlogger"
)

type Args map[string]string

func (a Args) Put(key, value string) {
	a[key] = value
}

func (a Args) Get(key string) (string, bool) {
	v, ok := a[key]
	return v, ok
}

func (a Args) MustGet(key, def string) string {
	if v, ok := a[key]; ok {
		return v
	}
	return def
}

func (a Args) PutInt(key string, value int) {
	a[key] = strconv.Itoa(value)
}

func (a Args) GetInt(key string) (int, bool) {
	if v, ok := a[key]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i, true
		}
	}
	return 0, false
}

func (a Args) MustGetInt(key string, def int) int {
	if v, ok := a.GetInt(key); ok {
		return v
	}
	return def
}

func (a Args) PutInt32(key string, value int32) {
	a[key] = strconv.Itoa(int(value))
}

func (a Args) GetInt32(key string) (int32, bool) {
	if v, ok := a[key]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			return int32(i), true
		}
	}
	return 0, false
}

func (a Args) MustGetInt32(key string, def int32) int32 {
	if v, ok := a.GetInt32(key); ok {
		return v
	}
	return def
}

func (a Args) PutInt64(key string, value int64) {
	a[key] = strconv.FormatInt(value, 10)
}

func (a Args) GetInt64(key string) (int64, bool) {
	if v, ok := a[key]; ok {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i, true
		}
	}
	return 0, false
}

func (a Args) MustGetInt64(key string, def int64) int64 {
	if v, ok := a.GetInt64(key); ok {
		return v
	}
	return def
}

func (a Args) PutFloat32(key string, value float32) {
	a[key] = strconv.FormatFloat(float64(value), 'f', -1, 32)
}

func (a Args) GetFloat32(key string) (float32, bool) {
	if v, ok := a[key]; ok {
		if i, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(i), true
		}
	}
	return 0.0, false
}

func (a Args) MustGetFloat32(key string, def float32) float32 {
	if v, ok := a.GetFloat32(key); ok {
		return v
	}
	return def
}

func (a Args) PutBool(key string, value bool) {
	a[key] = strconv.FormatBool(value)
}

func (a Args) GetBool(key string) (bool, bool) {
	if v, ok := a[key]; ok {
		if i, err := strconv.ParseBool(v); err == nil {
			return i, true
		}
	}
	return false, false
}

func (a Args) MustGetBool(key string, def bool) bool {
	if v, ok := a.GetBool(key); ok {
		return v
	}
	return def
}

func (a Args) PutBytes(key string, value []byte) {
	a[key] = string(value)
}

func (a Args) GetBytes(key string) ([]byte, bool) {
	if v, ok := a[key]; ok {
		return []byte(v), true
	}
	return nil, false
}

func (a Args) MustGetBytes(key string, def []byte) []byte {
	if v, ok := a.GetBytes(key); ok {
		return v
	}
	return def
}

type Context struct {
	callSeq        int32
	RemoteNodeName string
	conn           net.IConn
	Method         string
	Args
}

func (c *Context) Return(r *Result) {
	if body, err := defaultCodec.Marshal(&RPCMethodReturn{CallSeq: c.callSeq, Method: c.Method, Returns: r.Args}); err == nil {
		if !c.conn.Send(net.MessageId(RPC_MSGID_RETURN), body) {
			fwlogger.E("-- return rpc method %s failed! send message failed --", c.Method)
		}
	} else {
		fwlogger.E("-- return rpc method %s failed! marshal message failed --", c.Method)
	}
}

func (c *Context) ReturnNil() {
	if body, err := defaultCodec.Marshal(&RPCMethodReturn{CallSeq: c.callSeq, Method: c.Method, Returns: nil}); err == nil {
		if !c.conn.Send(net.MessageId(RPC_MSGID_RETURN), body) {
			fwlogger.E("-- return rpc method %s failed! send message failed --", c.Method)
		}
	} else {
		fwlogger.E("-- return rpc method %s failed! marshal message failed --", c.Method)
	}
}
