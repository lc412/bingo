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

package node

import (
	"github.com/snippetor/bingo/rpc"
	"github.com/snippetor/bingo/net"
	"github.com/valyala/fasthttp"
	"github.com/snippetor/bingo/proto"
	"github.com/snippetor/bingo/codec"
	"strings"
)

var (
	models = make(map[string]IModel)
)

type IModel interface {
	OnInit()
	OnServiceClientConnected(string, net.IConn)
	OnServiceClientDisconnected(string, net.IConn)
	OnReceiveServiceMessage(net.IConn, net.MessageId, body net.MessageBody)
	OnReceiveHttpServiceRequest(*fasthttp.RequestCtx)

	init()
	setNodeName(string)
	setRPCServer(*rpc.Server)
	putRPCClient(string, *rpc.Client)
	putService(string, net.IServer)
}

// rpc
type RPCModule struct {
	nodeName   string
	rpcClients map[string]*rpc.Client
	rpcServer  *rpc.Server
}

func (m *RPCModule) CallMethod(nodeName string, methodName string, args *rpc.Args, callback rpc.RPCCallback) {
	if m.rpcServer != nil && m.rpcServer.HasEndName(nodeName) {
		m.rpcServer.Call(nodeName, methodName, args, callback)
	} else if m.rpcClients != nil {
		if c, ok := m.rpcClients[nodeName]; ok && c != nil {
			c.Call(methodName, args, callback)
		}
	}
}

func (m *RPCModule) CallMethodNoReturn(nodeName string, methodName string, args *rpc.Args) {
	if m.rpcServer != nil && m.rpcServer.HasEndName(nodeName) {
		m.rpcServer.CallNoReturn(nodeName, methodName, args)
	} else if m.rpcClients != nil {
		if c, ok := m.rpcClients[nodeName]; ok && c != nil {
			c.CallNoReturn(methodName, args)
		}
	}
}

func (m *RPCModule) RegisterMethod(methodName string, method rpc.RPCMethod) {
	rpc.RegisterMethod(m.nodeName, methodName, method)
}

func (m *RPCModule) GetServer() *rpc.Server {
	return m.rpcServer
}

func (m *RPCModule) GetClient(name string) (*rpc.Client, bool) {
	if m.rpcClients != nil {
		if c, ok := m.rpcClients[name]; ok {
			return c, ok
		}
	}
	return nil, false
}

func (m *RPCModule) GetClientsWithPrefix(nodeNamePrefix string) []*rpc.Client {
	clients := make([]*rpc.Client, 0)
	if m.rpcClients != nil {
		for n, c := range m.rpcClients {
			if strings.HasPrefix(n, nodeNamePrefix) {
				clients = append(clients, c)
			}
		}
	}
	return clients
}

// service
type ServiceModule struct {
	servers map[string]net.IServer
}

func (m *ServiceModule) GetService(name string) (net.IServer, bool) {
	if m.servers != nil {
		if c, ok := m.servers[name]; ok {
			return c, ok
		}
	}
	return nil, false
}

// proto
type ProtoModule struct {
	proto          *proto.MessageProtocol
	clientProtoVer string
}

func (m *ProtoModule) SetDefaultMessageProtocol(c codec.CodecType) {
	m.proto = proto.NewMessageProtocol(c)
}

func (m *ProtoModule) GetDefaultMessageProtocol() *proto.MessageProtocol {
	return m.proto
}

func (m *ProtoModule) SetClientProtoVersion(version string) {
	m.clientProtoVer = version
}

func (m *ProtoModule) GetClientProtoVersion() string {
	return m.clientProtoVer
}

func (m *ProtoModule) RegisterProto(msgId net.MessageId, v interface{}, protoVersion string) {
	m.proto.GetProtoCollection().Put(msgId, v, protoVersion)
}

func (m *ProtoModule) RegisterProtoDefault(msgId net.MessageId, v interface{}) {
	m.proto.GetProtoCollection().PutDefault(msgId, v)
}

func (m *ProtoModule) Marshal(v interface{}) (net.MessageBody, bool) {
	return m.proto.Marshal(v)
}

func (m *ProtoModule) Unmarshal(data net.MessageBody, v interface{}) bool {
	return m.proto.Unmarshal(data, v)
}

func (m *ProtoModule) UnmarshalTo(msgId net.MessageId, data net.MessageBody, clientProtoVersion string) (interface{}, bool) {
	return m.proto.UnmarshalTo(msgId, data, clientProtoVersion)
}

func (m *ProtoModule) UnmarshalToDefault(msgId net.MessageId, data net.MessageBody) (interface{}, bool) {
	return m.proto.UnmarshalToDefault(msgId, data)
}

// model
type Model struct {
	nodeName string
	RPC      *RPCModule
	Service  *ServiceModule
	Proto    *ProtoModule
}

func (m *Model) init() {
	// init RPC module
	m.RPC = &RPCModule{nodeName: m.nodeName}
	// init Service module
	m.Service = &ServiceModule{}
	// init Proto module
	m.Proto = &ProtoModule{proto: proto.NewMessageProtocol(codec.Protobuf), clientProtoVer: ""}
	m.OnInit()
}

func (m *Model) setNodeName(name string) {
	m.nodeName = name
}

func (m *Model) putService(name string, s net.IServer) {
	if m.Service.servers == nil {
		m.Service.servers = make(map[string]net.IServer)
	}
	m.Service.servers[name] = s
}

func (m *Model) setRPCServer(serv *rpc.Server) {
	m.RPC.rpcServer = serv
}

func (m *Model) putRPCClient(name string, client *rpc.Client) {
	if m.RPC.rpcClients == nil {
		m.RPC.rpcClients = make(map[string]*rpc.Client)
	}
	m.RPC.rpcClients[name] = client
}

func (m *Model) OnInit() {
}

func (m *Model) OnServiceClientConnected(serviceName string, conn net.IConn) {
}

func (m *Model) OnServiceClientDisconnected(serviceName string, conn net.IConn) {
}

func (m *Model) OnReceiveServiceMessage(conn net.IConn, msgId net.MessageId, msgBody net.MessageBody) {
}

func (m *Model) OnReceiveHttpServiceRequest(ctx *fasthttp.RequestCtx) {
}

func (m *Model) GetNodeName() string {
	return m.nodeName
}

func BindNodeModel(modelName string, m IModel) {
	models[modelName] = m
}

func getNodeModel(modelName string) (IModel, bool) {
	m, ok := models[modelName]
	return m, ok
}
