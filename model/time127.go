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
	"encoding/json/jsontext"
	"strconv"
)

const (
	smallestWithMilliPrecision int64 = -8796093022208000
	largestWithMilliPrecision  int64 = 8796093022208000
)

func (t *Time) MarshalJSONTo(enc *jsontext.Encoder) error {
	i := int64(*t)

	// outside the range where floats can maintain thousandths precision, match MarshalJSON to get identical float truncation/rounding
	if i < smallestWithMilliPrecision || i > largestWithMilliPrecision {
		return enc.WriteToken(jsontext.Float(float64(*t) / float64(second)))
	}

	buf := make([]byte, 0, 18) // length of min/max int and optional space for - and .

	// Write out the timestamp as a float divided by 1000.
	// This is faster than converting to a float.

	var fraction int64
	if i < 0 {
		buf = append(buf, '-')
		// negating i is safe from overflowing because we can't get here when i is MinInt64
		fraction = -i % 1000
		i = -i / 1000
	} else {
		fraction = i % 1000
		i = i / 1000
	}

	buf = strconv.AppendInt(buf, i, 10)
	if fraction > 0 {
		buf = append(buf, '.')
		for _, place := range []int64{100, 10, 1} {
			if fraction == 0 {
				break
			}
			buf = append(buf, byte('0'+fraction/place))
			fraction %= place
		}
	}

	return enc.WriteValue(jsontext.Value(buf))
}

// // UnmarshalJSONFrom implements the json.Unmarshaler interface.
// func (t *Time) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
// 	b, err := dec.ReadValue()
// 	if err != nil {
// 		return err
// 	}

// 	base, frac, found := bytes.Cut(b, []byte{'.'})
// 	if !found {
// 		v, err := strconv.ParseInt(string(base), 10, 64)
// 		if err != nil {
// 			return err
// 		}
// 		*t = Time(v * second)
// 	} else {
// 		v, err := strconv.ParseInt(string(base), 10, 64)
// 		if err != nil {
// 			return err
// 		}

// 		prec := dotPrecision - len(frac)
// 		if prec < 0 {
// 			frac = frac[:dotPrecision]
// 		}
// 		va, err := strconv.ParseInt(string(frac), 10, 32)
// 		if err != nil {
// 			return err
// 		}
// 		switch prec {
// 		case 1:
// 			va *= 10
// 		case 2:
// 			va *= 100
// 		}

// 		if len(base) > 0 && base[0] == '-' {
// 			va = -va
// 		}
// 		*t = Time(v*second + va)
// 	}
// 	return nil
// }
