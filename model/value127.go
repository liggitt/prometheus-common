//go:build go1.27

package model

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

func (ss *SampleStream) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}

	// "metrics", always written
	if err := enc.WriteToken(jsontext.String("metric")); err != nil {
		return err
	}
	if ss.Metric == nil {
		if err := enc.WriteToken(jsontext.Null); err != nil {
			return err
		}
	} else {
		// deterministic order for json/v1 compatibility
		if err := json.MarshalEncode(enc, ss.Metric, json.Deterministic(true)); err != nil {
			return err
		}
	}

	// "values", written if non-empty or if histograms is empty
	if len(ss.Values) > 0 || len(ss.Histograms) == 0 {
		if err := enc.WriteToken(jsontext.String("values")); err != nil {
			return err
		}
		if ss.Values == nil {
			if err := enc.WriteToken(jsontext.Null); err != nil {
				return err
			}
		} else {
			if err := enc.WriteToken(jsontext.BeginArray); err != nil {
				return err
			}
			for i := range ss.Values {
				if err := ss.Values[i].MarshalJSONTo(enc); err != nil {
					return err
				}
			}
			if err := enc.WriteToken(jsontext.EndArray); err != nil {
				return err
			}
		}
	}

	// "histograms", written if non-empty
	if len(ss.Histograms) > 0 {
		if err := enc.WriteToken(jsontext.String("histograms")); err != nil {
			return err
		}
		if ss.Histograms == nil {
			if err := enc.WriteToken(jsontext.Null); err != nil {
				return err
			}
		} else {
			if err := enc.WriteToken(jsontext.BeginArray); err != nil {
				return err
			}
			for i := range ss.Histograms {
				if err := ss.Histograms[i].MarshalJSONTo(enc); err != nil {
					return err
				}
			}
			if err := enc.WriteToken(jsontext.EndArray); err != nil {
				return err
			}
		}
	}

	return enc.WriteToken(jsontext.EndObject)
}
