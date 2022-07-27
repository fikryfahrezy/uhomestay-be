package dues

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/sync/errgroup"
)

var (
	ErrFileRequired   = errors.New("file tidak boleh kosong")
	ErrIsPaidRequired = errors.New("status persetujuan tidak boleh kosong")
	ErrMaxFilename    = errors.New("nama file tidak dapat lebih dari 200 karakter")
)

func ValidatePayMemberDuesIn(i PayMemberDuesIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.File.File == nil || i.File.Filename == "" {
			return ErrFileRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.File.Filename) > 200 {
			return ErrMaxFilename
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditMemberDuesIn(i EditMemberDuesIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.File.File == nil || i.File.Filename == "" {
			return ErrFileRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.File.Filename) > 200 {
			return ErrMaxFilename
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidatePaidMemberDuesIn(i PaidMemberDuesIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if !i.IsPaid.Valid {
			return ErrIsPaidRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
