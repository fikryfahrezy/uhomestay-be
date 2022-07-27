package image

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/sync/errgroup"
)

var (
	ErrImageRequired = errors.New("gambar tidak boleh kosong")
	ErrMaxImagename  = errors.New("nama gambar tidak dapat lebih dari 200 karakter")
)

func ValidateAddImageIn(i AddImageIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.File.File == nil || i.File.Filename == "" {
			return ErrImageRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.File.Filename) > 200 {
			return ErrMaxImagename
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
