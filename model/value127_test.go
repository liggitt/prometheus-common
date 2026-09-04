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
	"fmt"
	"testing"
)

func TestSampleStreamJSONV2(t *testing.T) {
	cases := []SampleStream{
		{},
		{Metric: Metric{}},
		{Metric: Metric{"a": "", "b": "", "c": "", "d": "", "e": "", "f": "", "g": ""}},
		{Values: []SamplePair{}},
		{Values: []SamplePair{{}}},
		{Values: []SamplePair{{Timestamp: Time(1), Value: SampleValue(1)}}},
		{Histograms: []SampleHistogramPair{{Timestamp: Time(1), Histogram: &SampleHistogram{}}}},
		{Metric: Metric{}, Values: []SamplePair{}},
		{Metric: Metric{}, Values: []SamplePair{{}}, Histograms: []SampleHistogramPair{{Timestamp: Time(1), Histogram: &SampleHistogram{}}}},
	}

	for i, v := range cases {
		testV1V2Marshal(t, fmt.Sprintf("%d_%#v", i, v), v)
		testRoundTrip(t, fmt.Sprintf("%d_%#v", i, v), v)
	}
}
