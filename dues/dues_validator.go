package dues

import (
	"errors"
	"strings"

	"golang.org/x/sync/errgroup"
)

func ValidateAddDuesIn(i AddDuesIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if strings.Trim(i.Date, " ") == "" {
			return errors.New("date required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.IdrAmount, " ") == "" {
			return errors.New("idr_amount required")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditDuesIn(i EditDuesIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if strings.Trim(i.Date, " ") == "" {
			return errors.New("date required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.IdrAmount, " ") == "" {
			return errors.New("idr_amount required")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
