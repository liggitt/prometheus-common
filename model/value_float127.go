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
	"fmt"
	"strconv"
)

func (s *SamplePair) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}
	if err := s.Timestamp.MarshalJSONTo(enc); err != nil {
		return err
	}
	if err := s.Value.MarshalJSONTo(enc); err != nil {
		return err
	}
	return enc.WriteToken(jsontext.EndArray)
}

func (s *SamplePair) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if t, err := dec.ReadToken(); err != nil {
		return err
	} else if t.Kind() != jsontext.BeginArray.Kind() {
		return fmt.Errorf("expected [")
	}

	if v, err := dec.ReadValue(); err != nil {
		return err
	} else if err := s.Timestamp.UnmarshalJSON(v); err != nil {
		return err
	}

	if v, err := dec.ReadValue(); err != nil {
		return err
	} else if err := s.Value.UnmarshalJSON(v); err != nil {
		return err
	}

	if t, err := dec.ReadToken(); err != nil {
		return err
	} else if t.Kind() != jsontext.EndArray.Kind() {
		return fmt.Errorf("expected ]")
	}
	return nil
}

func (v *SampleValue) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalFloatAsString(enc, float64(*v))
}

func marshalFloatAsString(enc *jsontext.Encoder, f float64) error {
	buf := make([]byte, 0, 32)
	buf = append(buf, '"')
	buf = strconv.AppendFloat(buf, f, 'f', -1, 64)
	buf = append(buf, '"')
	return enc.WriteValue(jsontext.Value(buf))
}
