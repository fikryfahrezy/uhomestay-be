package cashflow

import (
	"errors"
	"strings"

	"golang.org/x/sync/errgroup"
)

func ValidateAddCashflowIn(i AddCashflowIn, ct CashflowType) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if ct == Unknown {
			return errors.New("type unknown")
		}

		return nil
	})
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

func ValidateEditCashflowIn(i EditCashflowIn, ct CashflowType) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if ct == Unknown {
			return errors.New("type unknown")
		}

		return nil
	})
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
