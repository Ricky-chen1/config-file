// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cloudwego/kitex-examples/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/kitex_gen/api/echo"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kitexserver "github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/config-file/filewatcher"
	"github.com/kitex-contrib/config-file/parser"
	fileserver "github.com/kitex-contrib/config-file/server"
	"github.com/kitex-contrib/config-file/utils"
)

var _ api.Echo = &EchoImpl{}

const (
	filepath    = "kitex_server.json"
	key         = "ServiceName"
	serviceName = "ServiceName"
)

// EchoImpl implements the last service interface defined in the IDL.
type EchoImpl struct{}

// Echo implements the Echo interface.
func (s *EchoImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	klog.Info("echo called")
	return &api.Response{Message: req.Message}, nil
}

// customed by user
type MyParser struct{}

// one example for custom parser
// if the type of server config is json or yaml,just using default parser
func (p *MyParser) Decode(kind parser.ConfigType, data []byte, config interface{}) error {
	return json.Unmarshal(data, config)
}

func main() {
	klog.SetLevel(klog.LevelDebug)

	// create a file watcher object
	fw, err := filewatcher.NewFileWatcher(filepath)
	if err != nil {
		panic(err)
	}
	// start watching file changes
	if err = fw.StartWatching(); err != nil {
		panic(err)
	}
	defer fw.StopWatching()

	// customed by user
	params := &parser.ConfigParam{}
	opts := &utils.Options{
		CustomParser: &MyParser{},
		CustomParams: params,
	}

	svr := echo.NewServer(
		new(EchoImpl),
		kitexserver.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
		kitexserver.WithSuite(fileserver.NewSuite(key, fw, opts)), // add watcher
	)
	if err := svr.Run(); err != nil {
		log.Println("server stopped with error:", err)
	} else {
		log.Println("server stopped")
	}
}
