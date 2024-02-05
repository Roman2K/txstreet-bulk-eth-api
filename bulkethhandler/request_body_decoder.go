package bulkethhandler

import (
	"encoding/json"
	"net/http"
)

const maxBodySize = 50 * 1024 * 1024

type requestBodyDecoder[T any] struct {
	w http.ResponseWriter
	r *http.Request
}

func (dec requestBodyDecoder[T]) Decode() (T, error) {
	reader := http.MaxBytesReader(dec.w, dec.r.Body, maxBodySize)
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var result T

	err := decoder.Decode(&result)

	return result, err
}
