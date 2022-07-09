package dues

import (
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	ErrFileRequired   = errors.New("file tidak boleh kosong")
	ErrIsPaidRequired = errors.New("status persetujuan tidak boleh kosong")
)

func ValidatePayMemberDuesIn(i PayMemberDuesIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.File.File == nil || i.File.Filename == "" {
			return ErrFileRequired
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
