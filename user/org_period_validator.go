package user

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func ValidateAddPeriodIn(i AddPeriodIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.StartDate, " ") == "" {
			return errors.New("start_date required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.EndDate, " ") == "" {
			return errors.New("end_date required")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditPeriodIn(i EditPeriodIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.StartDate, " ") == "" {
			return errors.New("start_date required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.EndDate, " ") == "" {
			return errors.New("end_date required")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
