// Copyright 2022 The ILLA Authors.
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

package redis

type Options struct {
	Host             string `validate:"required"`
	Port             string `validate:"required"`
	DatabaseIndex    int    `validate:"gte=0,lte=16"`
	DatabaseUsername string `validate:"required"`
	DatabasePassword string `validate:"required"`
	SSL              bool
}

type Command struct {
	Mode  string `validate:"required,oneof=select raw"`
	Query string
}
