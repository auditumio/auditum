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

// Package sql contains infrastructure layer SQL implementations.
//
// This package does not import any specific driver. The code must work on any
// supported dialect, so it should be as agnostic as possible. However, when
// there is a need for some specific dialect features, the code must be protected
// with constructions like:
//
//	if tx.Dialect().Name() == dialect.PG {
//	   // do something specific to PostgreSQL
//	}
package sql
