package aud_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/auditumio/auditum/internal/aud"
)

func TestDecodePageToken(t *testing.T) {
	id := aud.MustParseID("00000000-0000-0000-0000-000000000002")
	opTime := time.Date(2023, 1, 1, 0, 2, 0, 0, time.UTC)

	tests := []struct {
		name   string
		token  string
		cursor any
		want   any
		err    error
	}{
		{
			token:  "eyJsb3QiOiIyMDIzLTAxLTAxVDAwOjAyOjAwWiIsImxpZCI6WzAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDJdfQ",
			cursor: &aud.RecordCursor{},
			want: &aud.RecordCursor{
				LastOperationTime: &opTime,
				LastID:            &id,
			},
			err: nil,
		},
		{
			token:  "e30",
			cursor: &aud.RecordCursor{},
			want:   &aud.RecordCursor{},
			err:    nil,
		},
		{
			token:  "",
			cursor: &aud.RecordCursor{},
			want:   &aud.RecordCursor{},
			err:    nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := aud.DecodePageToken(test.token, test.cursor)
			assert.Equal(t, test.want, test.cursor)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestEncodePageToken(t *testing.T) {
	id1 := aud.MustParseID("00000000-0000-0000-0000-000000000001")
	opTime1 := time.Date(2023, 1, 1, 0, 1, 0, 0, time.UTC)
	id2 := aud.MustParseID("00000000-0000-0000-0000-000000000002")
	opTime2 := time.Date(2023, 1, 1, 0, 2, 0, 0, time.UTC)

	tests := []struct {
		name   string
		cursor aud.Cursor
		want   string
		err    error
	}{
		{
			cursor: aud.NewRecordCursor(
				[]aud.Record{
					{
						ID: id1,
						Operation: aud.Operation{
							Time: opTime1,
						},
					},
					{
						ID: id2,
						Operation: aud.Operation{
							Time: opTime2,
						},
					},
				},
				2,
			),
			want: "eyJsb3QiOiIyMDIzLTAxLTAxVDAwOjAyOjAwWiIsImxpZCI6WzAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDAsMCwwLDJdfQ",
			err:  nil,
		},
		{
			cursor: aud.NewRecordCursor(
				[]aud.Record{
					{
						ID: id1,
						Operation: aud.Operation{
							Time: opTime1,
						},
					},
					{
						ID: id2,
						Operation: aud.Operation{
							Time: opTime2,
						},
					},
				},
				20,
			),
			want: "",
			err:  nil,
		},
		{
			cursor: aud.NewRecordCursor(
				[]aud.Record{},
				20,
			),
			want: "",
			err:  nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := aud.EncodePageToken(test.cursor)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.err, err)
		})
	}
}
