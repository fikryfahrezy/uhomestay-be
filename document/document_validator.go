package document

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/sync/errgroup"
)

var (
	ErrDirNameRequired       = errors.New("nama folder tidak boleh kosong")
	ErrFileRequired          = errors.New("file tidak boleh kosong")
	ErrParentDirRequired     = errors.New("parent folder tidak boleh kosong")
	ErrStatusPrivateRequired = errors.New("status privasi tidak boleh kosong")
	ErrMaxDirName            = errors.New("nama folder tidak dapat lebih dari 200 karakter")
	ErrMaxFileName           = errors.New("nama file tidak dapat lebih dari 200 karakter")
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
	g.Go(func() error {
		if utf8.RuneCountInString(i.Name) > 200 {
			return ErrMaxDirName
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
	g.Go(func() error {
		if utf8.RuneCountInString(i.File.Filename) > 200 {
			return ErrMaxFileName
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
	g.Go(func() error {
		if utf8.RuneCountInString(i.Name) > 200 {
			return ErrMaxDirName
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
		if !i.IsPrivate.Valid {
			return ErrStatusPrivateRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.File.Filename) > 200 {
			return ErrMaxFileName
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
