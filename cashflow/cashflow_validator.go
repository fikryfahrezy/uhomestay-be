package cashflow

import (
	"errors"
	"strings"

	"golang.org/x/sync/errgroup"
)

var (
	ErrUnknownType       = errors.New("tipe cashflow tidak diketahui, tipe yang diperbolehkan 'pemasukan' atau 'pengeluaran'")
	ErrDateRequired      = errors.New("tanggal tidak boleh kosong")
	ErrIdrAmountRequired = errors.New("jumlah nominal rupiah tidak boleh kosong")
)

func ValidateAddCashflowIn(i AddCashflowIn, ct CashflowType) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if ct == Unknown {
			return ErrUnknownType
		}

		return nil
	})
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

func ValidateEditCashflowIn(i EditCashflowIn, ct CashflowType) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if ct == Unknown {
			return ErrUnknownType
		}

		return nil
	})
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
