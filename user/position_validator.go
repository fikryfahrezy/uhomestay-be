package user

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	ErrPositionNameRequired  = errors.New("nama jabatan tidak boleh kosong")
	ErrPositionLevelRequired = errors.New("level jabatan tidak boleh kosong")
)

func ValidateAddPositionIn(i AddPositionIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return ErrPositionNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if i.Level <= 0 {
			return ErrPositionLevelRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditPositionIn(i EditPositionIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return ErrPositionNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if i.Level <= 0 {
			return ErrPositionLevelRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
