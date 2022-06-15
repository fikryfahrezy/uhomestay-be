package dues

import (
	"errors"

	"golang.org/x/sync/errgroup"
)

func ValidatePayMemberDuesIn(i PayMemberDuesIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.File.File == nil || i.File.Filename == "" {
			return errors.New("file required")
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
			return errors.New("file required")
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
			return errors.New("is_paid required")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
