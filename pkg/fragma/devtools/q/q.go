// Copyright 2023 Igor Zibarev
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

// Package q wraps github.com/ryboe/q for convenient usage during development
// without having to add/remove the module.
package q

import "github.com/ryboe/q"

func init() {
	q.CallDepth = 3
}

func Q(v ...any) {
	q.Q(v...)
}
