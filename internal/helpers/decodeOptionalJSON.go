package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DecodeOptionalJSON(r *http.Request, dst any) error {

	if r.Body == nil || r.Body == http.NoBody || r.ContentLength == 0 {
		return nil
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	return nil
}
