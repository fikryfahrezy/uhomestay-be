package document

import (
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	ErrDirNameRequired       = errors.New("nama folder tidak boleh kosong")
	ErrFileRequired          = errors.New("file tidak boleh kosong")
	ErrParentDirRequired     = errors.New("folder induk tidak boleh kosong")
	ErrStatusPrivateRequired = errors.New("status privasi tidak boleh kosong")
)

func ValidateAddDirDocumentIn(i AddDirDocumentIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.Name == "" {
			return ErrDirNameRequired
		}
		return nil
	})

	g.Go(func() error {
		if !i.DirId.Valid {
			return ErrParentDirRequired
		}
		return nil
	})

	g.Go(func() error {
		if !i.IsPrivate.Valid {
			return ErrStatusPrivateRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateAddFileDocumentIn(i AddFileDocumentIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if !i.DirId.Valid {
			return ErrParentDirRequired
		}
		return nil
	})
	g.Go(func() error {
		if i.File.File == nil || i.File.Filename == "" {
			return ErrFileRequired
		}
		return nil
	})
	g.Go(func() error {
		if !i.IsPrivate.Valid {
			return ErrStatusPrivateRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditDirDocumentIn(i EditDirDocumentIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.Name == "" {
			return ErrDirNameRequired
		}
		return nil
	})

	g.Go(func() error {
		if !i.IsPrivate.Valid {
			return ErrStatusPrivateRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditFileDocumentIn(i EditFileDocumentIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.File.File == nil || i.File.Filename == "" {
			return ErrFileRequired
		}
		return nil
	})

	g.Go(func() error {
		if !i.IsPrivate.Valid {
			return ErrStatusPrivateRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
