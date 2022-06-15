package document

import (
	"errors"

	"golang.org/x/sync/errgroup"
)

func ValidateAddDirDocumentIn(i AddDirDocumentIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if i.Name == "" {
			return errors.New("name required")
		}
		return nil
	})

	g.Go(func() error {
		if !i.DirId.Valid {
			return errors.New("dir_id required")
		}
		return nil
	})

	g.Go(func() error {
		if !i.IsPrivate.Valid {
			return errors.New("is_private required")
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
			return errors.New("dir_id required")
		}
		return nil
	})
	g.Go(func() error {
		if i.File.File == nil || i.File.Filename == "" {
			return errors.New("file required")
		}
		return nil
	})
	g.Go(func() error {
		if !i.IsPrivate.Valid {
			return errors.New("is_private required")
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
			return errors.New("name required")
		}
		return nil
	})

	g.Go(func() error {
		if !i.IsPrivate.Valid {
			return errors.New("is_private required")
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
			return errors.New("file required")
		}
		return nil
	})

	g.Go(func() error {
		if !i.IsPrivate.Valid {
			return errors.New("is_private required")
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
