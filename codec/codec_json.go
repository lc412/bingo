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

package codec

import (
	json1 "encoding/json"
	"github.com/snippetor/bingo/net"
)

// JSON消息协议
type json struct {
}

func (j *json) Marshal(v interface{}) (net.MessageBody, error) {
	return json1.Marshal(v)
}

func (j *json) Unmarshal(data net.MessageBody, v interface{}) error {
	return json1.Unmarshal(data, v)
}
