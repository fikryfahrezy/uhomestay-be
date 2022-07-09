package dues

import (
	"errors"
	"strings"

	"golang.org/x/sync/errgroup"
)

var (
	ErrDateRequired      = errors.New("tanggal tidak boleh kosong")
	ErrIdrAmountRequired = errors.New("jumlah nominal rupiah tidak boleh kosong")
)

func ValidateAddDuesIn(i AddDuesIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if strings.Trim(i.Date, " ") == "" {
			return ErrDateRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.IdrAmount, " ") == "" {
			return ErrIdrAmountRequired
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
			return ErrDateRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.IdrAmount, " ") == "" {
			return ErrIdrAmountRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
