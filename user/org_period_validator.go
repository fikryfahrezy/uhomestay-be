package user

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	ErrPeriodStartDateRequired = errors.New("tanggal mulai periode tidak boleh kosong")
	ErrPeriodEndDateRequired   = errors.New("tanggal mulai periode tidak boleh kosong")
)

func ValidateAddPeriodIn(i AddPeriodIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.StartDate, " ") == "" {
			return ErrPeriodStartDateRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.EndDate, " ") == "" {
			return ErrPeriodEndDateRequired
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
			return ErrPeriodStartDateRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.EndDate, " ") == "" {
			return ErrPeriodEndDateRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
