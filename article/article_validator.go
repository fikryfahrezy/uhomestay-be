package article

import (
	"unicode/utf8"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	ErrMaxTitle     = errors.New("judul tidak dapat lebih dari 200 karakter")
	ErrMaxShortDesc = errors.New("deskripsi singkat tidak dapat lebih dari 200 karakter")
	ErrMaxSlug      = errors.New("slug tidak dapat lebih dari 200 karakter")
)

func ValidateAddArticleIn(i AddArticleIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if utf8.RuneCountInString(i.Title) > 200 {
			return ErrMaxTitle
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.ShortDesc) > 200 {
			return ErrMaxShortDesc
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Slug) > 200 {
			return ErrMaxSlug
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditArticleIn(i EditArticleIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if utf8.RuneCountInString(i.Title) > 200 {
			return ErrMaxTitle
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.ShortDesc) > 200 {
			return ErrMaxShortDesc
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
