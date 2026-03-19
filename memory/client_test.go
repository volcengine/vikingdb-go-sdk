// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package memory

import (
	"encoding/json"
	"testing"

	mmodel "github.com/volcengine/vikingdb-go-sdk/memory/model"
)

func TestParseErrorCode(t *testing.T) {
	cases := []struct {
		name string
		in   interface{}
		want int
	}{
		{name: "nil", in: nil, want: 0},
		{name: "json_number", in: json.Number("1000001"), want: 1000001},
		{name: "string_number", in: "1000003", want: 1000003},
		{name: "string_invalid", in: "abc", want: 0},
		{name: "int", in: 42, want: 42},
		{name: "int64", in: int64(7), want: 7},
	}
	for _, tc := range cases {
		if got := parseErrorCode(tc.in); got != tc.want {
			t.Fatalf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestPromoteException(t *testing.T) {
	base := &mmodel.VikingMemException{
		Code:       1000001,
		Message:    "unauthorized",
		StatusCode: 401,
		RequestID:  "rid",
	}
	err := mmodel.PromoteException(base)
	if _, ok := err.(*mmodel.UnauthorizedException); !ok {
		t.Fatalf("expected UnauthorizedException, got %T", err)
	}
}

func TestCollectionDefaultProjectName(t *testing.T) {
	c, err := New(AuthNone())
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	coll, err := c.Collection(mmodel.CollectionMeta{CollectionName: "c1"})
	if err != nil {
		t.Fatalf("collection: %v", err)
	}
	if coll.meta.ProjectName != "default" {
		t.Fatalf("project_name: got %q, want %q", coll.meta.ProjectName, "default")
	}
}
