package homestay

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/sync/errgroup"
)

var (
	ErrHomestayImageRequired = errors.New("foto atau gambar tidak boleh kosong")
	ErrMaxHomestayImageName  = errors.New("nama foto atau gambar tidak dapat lebih dari 200 karakter")
)

func ValidateAddHomestayImageIn(i AddHomestayImageIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.File.File == nil || i.File.Filename == "" {
			return ErrHomestayImageRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.File.Filename) > 200 {
			return ErrMaxHomestayImageName
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
