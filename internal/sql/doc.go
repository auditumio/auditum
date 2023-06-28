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
