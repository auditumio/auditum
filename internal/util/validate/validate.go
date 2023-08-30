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

package validate

type Validatable interface {
	Validate() error
}

func Each(v ...Validatable) error {
	for _, validatable := range v {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func ErrorAsValidatable(err error) Validatable {
	return errorValidatable{err: err}
}

type errorValidatable struct {
	err error
}

func (e errorValidatable) Validate() error {
	return e.err
}
