package homestay

import (
	"errors"
	"strings"
	"unicode/utf8"

	"golang.org/x/sync/errgroup"
)

var (
	ErrHomestayNameRequired   = errors.New("nama homestay tidak boleh kosong")
	ErrMaxHomestayName        = errors.New("nama homestay tidak dapat lebih dari 100 karakter")
	ErrHomestayAddresRequired = errors.New("alamat homestay tidak boleh kosong")
	ErrMaxHomestayAddress     = errors.New("alamat homestay tidak dapat lebih dari 200 karakter")
	ErrLongitudeRequired      = errors.New("titik garis bujur map homestay tidak boleh kosong")
	ErrLatitudeRequired       = errors.New("titik garis lintang map homestay tidak boleh kosong")
	ErrMaxHomestayLat         = errors.New("titik garis bujur map homestay tidak dapat lebih dari 50 karakter")
	ErrMaxHomestayLng         = errors.New("titik garis lintang map homestay tidak dapat lebih dari 50 karakter")
	ErrHomestayPhotosRequired = errors.New("foto homestay tidak boleh kosong")
)

func ValidateAddMemberHomestayIn(i AddMemberHomestayIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return ErrHomestayNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Name) > 100 {
			return ErrMaxHomestayName
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Address, " ") == "" {
			return ErrHomestayAddresRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Address) > 200 {
			return ErrMaxHomestayAddress
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Latitude, " ") == "" {
			return ErrLatitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Latitude) > 50 {
			return ErrMaxHomestayLat
		}
		return nil
	})

	g.Go(func() error {
		if strings.Trim(i.Longitude, " ") == "" {
			return ErrLongitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Longitude) > 50 {
			return ErrMaxHomestayLng
		}
		return nil
	})
	g.Go(func() error {
		if len(i.ImageIds) == 0 {
			return ErrHomestayImageRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditMemberHomestayIn(i EditMemberHomestayIn) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return ErrHomestayNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Name) > 100 {
			return ErrMaxHomestayName
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Address, " ") == "" {
			return ErrHomestayAddresRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Address) > 200 {
			return ErrMaxHomestayAddress
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Latitude, " ") == "" {
			return ErrLatitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Latitude) > 50 {
			return ErrMaxHomestayLat
		}
		return nil
	})

	g.Go(func() error {
		if strings.Trim(i.Longitude, " ") == "" {
			return ErrLongitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Longitude) > 50 {
			return ErrMaxHomestayLng
		}
		return nil
	})
	g.Go(func() error {
		if len(i.ImageIds) == 0 {
			return ErrHomestayImageRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
