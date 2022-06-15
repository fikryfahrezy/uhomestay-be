package history

import "golang.org/x/sync/errgroup"

func ValidateAddHistoryIn(i AddHistoryIn) error {
	g := new(errgroup.Group)

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
