package httpdecode

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/guregu/null.v4"
)

func Multipart(r *http.Request, in interface{}, maxMemory int64, fs ...mapstructure.DecodeHookFunc) error {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return err
	}

	formValues := make(map[string]interface{})
	formFile := r.MultipartForm
	for k, v := range formFile.Value {
		formValues[k] = v[0]
	}

	for k, v := range formFile.File {
		formValues[k] = v[0]
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &in,
		WeaklyTypedInput: true,
		DecodeHook:       mapstructure.ComposeDecodeHookFunc(fs...),
	})
	if err != nil {
		return err
	}

	err = decoder.Decode(formValues)
	if err != nil {
		return err
	}

	return nil
}

func BoolToNullBoolHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	var s null.Bool

	if t == reflect.TypeOf(s) {
		if v, err := strconv.ParseBool(fmt.Sprint(data)); err == nil {
			s = null.Bool{
				NullBool: sql.NullBool{
					Bool:  v,
					Valid: true,
				},
			}
		}

		return s, nil
	}

	return data, nil
}

func IntToNulIntHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	var s null.Int

	if t == reflect.TypeOf(s) {
		if v, err := strconv.ParseInt(fmt.Sprint(data), 10, 64); err == nil {
			s = null.Int{
				NullInt64: sql.NullInt64{
					Int64: v,
					Valid: true,
				},
			}
		}

		return s, nil
	}

	return data, nil
}

type File interface {
	io.Reader
	io.Closer
}

type FileHeader struct {
	Filename string
	File     File
}

func MultipartToFileHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	var s FileHeader

	if t == reflect.TypeOf(s) {
		if v, ok := data.(*multipart.FileHeader); ok {
			f, err := v.Open()
			if err != nil {
				return nil, err
			}

			s = FileHeader{
				Filename: v.Filename,
				File:     f,
			}
		}

		return s, nil
	}

	return data, nil
}
