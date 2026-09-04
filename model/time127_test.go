//go:build go1.27

// Copyright 2026 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestTimeJSONV2(t *testing.T) {
	cases := []int64{
		math.MinInt64,
		-9999999999995500,
		-9999999999994500,
		-9007199254740992, // smallest integer with float precision
		smallestWithMilliPrecision,
		-8123456789012345,
		math.MinInt32,
		-10000,
		-9000,
		-8000,
		-7000,
		-6000,
		-5000,
		-4000,
		-3000,
		-2000,
		-1000,
		-900,
		-800,
		-700,
		-600,
		-500,
		-400,
		-300,
		-200,
		-100,
		-10,
		-1,
		0,
		1,
		10,
		100,
		200,
		300,
		400,
		500,
		600,
		700,
		800,
		900,
		1000,
		2000,
		3000,
		4000,
		5000,
		6000,
		7000,
		8000,
		9000,
		10000,
		math.MaxInt32,
		8123456789012345,
		largestWithMilliPrecision,
		9007199254740992, // largest integer with float precision
		math.MaxInt64,
	}

	for _, i := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if i != math.MinInt64 {
				testV1V2Marshal(t, "-1", Time(i-1))
			}
			testV1V2Marshal(t, "=", Time(i))
			if i != math.MaxInt64 {
				testV1V2Marshal(t, "+1", Time(i+1))
			}
		})
	}
	testV1V2Marshal(t, "-1", Time(-1))
	testV1V2Marshal(t, "0", Time(0))
	testV1V2Marshal(t, "1", Time(1))
}

func testV1V2Marshal(t *testing.T, name string, v jsonv1.Marshaler) {
	t.Run("v1v2_"+name, func(t *testing.T) {
		b1, err := v.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		b2, err := jsonv2.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if string(b1) != string(b2) {
			t.Logf("v1=%s", string(b1))
			t.Logf("v2=%s", string(b2))
			t.Fatalf("v1/v2 mismatch")
		}
	})
}

func testRoundTrip(t *testing.T, name string, v any) {
	t.Run("roundtrip_"+name, func(t *testing.T) {
		b, err := jsonv1.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))

		outPtr := reflect.New(reflect.TypeOf(v)).Interface()
		if err := jsonv1.Unmarshal(b, outPtr); err != nil {
			t.Fatal(err)
		}

		out := reflect.ValueOf(outPtr).Elem().Interface()
		if !reflect.DeepEqual(v, out) {
			t.Fatalf("did not round-trip\nwant: %#v\ngot:  %#v", v, out)
		}
	})
}
