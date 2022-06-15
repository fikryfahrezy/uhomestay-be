package blog

import (
	"golang.org/x/sync/errgroup"
)

func ValidateAddBlogIn(i AddBlogIn) error {
	g := new(errgroup.Group)

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditBlogIn(i EditBlogIn) error {
	g := new(errgroup.Group)

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
