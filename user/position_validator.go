package user

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func ValidateAddPositionIn(i AddPositionIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return errors.New("name required")
		}
		return nil
	})
	g.Go(func() error {
		if i.Level <= 0 {
			return errors.New("level required")
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
			return errors.New("name required")
		}
		return nil
	})
	g.Go(func() error {
		if i.Level <= 0 {
			return errors.New("level required")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
