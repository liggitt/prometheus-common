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
	"encoding/json/v2"
	"errors"
	"fmt"
)

func (s *SampleHistogramPair) MarshalJSONTo(enc *jsontext.Encoder) error {
	if s.Histogram == nil {
		return errors.New("histogram is nil")
	}
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}
	if err := s.Timestamp.MarshalJSONTo(enc); err != nil {
		return err
	}
	if err := s.Histogram.MarshalJSONTo(enc); err != nil {
		return err
	}
	return enc.WriteToken(jsontext.EndArray)
}

func (s *SampleHistogramPair) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
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

	if err := json.UnmarshalDecode(dec, &s.Histogram); err != nil {
		return err
	} else if s.Histogram == nil {
		return errors.New("histogram is nil")
	}

	if t, err := dec.ReadToken(); err != nil {
		return err
	} else if t.Kind() != jsontext.EndArray.Kind() {
		return fmt.Errorf("expected ]")
	}
	return nil
}

func (s *SampleHistogram) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}

	if err := enc.WriteToken(jsontext.String("count")); err != nil {
		return err
	}
	if err := s.Count.MarshalJSONTo(enc); err != nil {
		return err
	}

	if err := enc.WriteToken(jsontext.String("sum")); err != nil {
		return err
	}
	if err := s.Sum.MarshalJSONTo(enc); err != nil {
		return err
	}

	if err := enc.WriteToken(jsontext.String("buckets")); err != nil {
		return err
	}
	if s.Buckets == nil {
		// write empty buckets as null for json v1 compatibility
		if err := enc.WriteToken(jsontext.Null); err != nil {
			return err
		}
	} else {
		if err := enc.WriteToken(jsontext.BeginArray); err != nil {
			return err
		}
		for i := range s.Buckets {
			if err := s.Buckets[i].MarshalJSONTo(enc); err != nil {
				return err
			}
		}
		if err := enc.WriteToken(jsontext.EndArray); err != nil {
			return err
		}
	}

	if err := enc.WriteToken(jsontext.EndObject); err != nil {
		return err
	}

	return nil
}

func (s *HistogramBucket) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}
	if err := enc.WriteToken(jsontext.Int(int64(s.Boundaries))); err != nil {
		return err
	}
	if err := s.Lower.MarshalJSONTo(enc); err != nil {
		return err
	}
	if err := s.Upper.MarshalJSONTo(enc); err != nil {
		return err
	}
	if err := s.Count.MarshalJSONTo(enc); err != nil {
		return err
	}
	return enc.WriteToken(jsontext.EndArray)
}

func (s *HistogramBucket) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if t, err := dec.ReadToken(); err != nil {
		return err
	} else if t.Kind() != jsontext.BeginArray.Kind() {
		return fmt.Errorf("expected [")
	}

	if err := json.UnmarshalDecode(dec, &s.Boundaries); err != nil {
		return nil
	}
	if err := json.UnmarshalDecode(dec, &s.Lower); err != nil {
		return nil
	}
	if err := json.UnmarshalDecode(dec, &s.Upper); err != nil {
		return nil
	}
	if err := json.UnmarshalDecode(dec, &s.Count); err != nil {
		return nil
	}

	if t, err := dec.ReadToken(); err != nil {
		return err
	} else if t.Kind() != jsontext.EndArray.Kind() {
		return fmt.Errorf("expected ]")
	}
	return nil
}

func (v *FloatString) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalFloatAsString(enc, float64(*v))
}
