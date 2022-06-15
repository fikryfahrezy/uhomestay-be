package user

import (
	"errors"

	"golang.org/x/sync/errgroup"
)

func ValidateAddGoalIn(i AddGoalIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.OrgPeriodId < 1 {
			return errors.New("org period id required")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
