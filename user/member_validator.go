package user

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	ErrMemberNameRequired     = errors.New("nama anggota tidak boleh kosong")
	ErrHomestayNameRequired   = errors.New("nama homestay anggota tidak boleh kosong")
	ErrUsernameRequired       = errors.New("username anggota tidak boleh kosong")
	ErrPositionRequired       = errors.New("jabatan tidak boleh kosong")
	ErrOrgPeriodRequired      = errors.New("periode organisasi tidak boleh kosong")
	ErrWaPhoneRequired        = errors.New("nomor whats app anggota tidak boleh kosong")
	ErrOtherPhoneRequired     = errors.New("nomor lainnya anggota tidak boleh kosong")
	ErrHomestayAddresRequired = errors.New("alamat homestay anggota tidak boleh kosong")
	ErrLongitudeRequired      = errors.New("titik garis bujur map homestay anggota tidak boleh kosong")
	ErrLatitudeRequired       = errors.New("titik garis lintang map homestay anggota tidak boleh kosong")
	ErrPasswordRequired       = errors.New("password anggota tidak boleh kosong")
	ErrIsAdminRequired        = errors.New("status admin tidak boleh kosong")
)

func ValidateAddMemberIn(i AddMemberIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return ErrMemberNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayName, " ") == "" {
			return ErrHomestayNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Username, " ") == "" {
			return ErrUsernameRequired
		}
		return nil
	})
	g.Go(func() error {
		if len(i.PositionIds) == 0 {
			return ErrPositionRequired
		}
		return nil
	})
	g.Go(func() error {
		if i.PeriodId == 0 {
			return ErrOrgPeriodRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.WaPhone, " ") == "" {
			return ErrWaPhoneRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.OtherPhone, " ") == "" {
			return ErrOtherPhoneRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayAddress, " ") == "" {
			return ErrHomestayAddresRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLatitude, " ") == "" {
			return ErrLatitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLongitude, " ") == "" {
			return ErrLongitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Password, " ") == "" {
			return ErrPasswordRequired
		}
		return nil
	})
	g.Go(func() error {
		if !i.IsAdmin.Valid {
			return ErrIsAdminRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateRegisterIn(i RegisterIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return ErrMemberNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayName, " ") == "" {
			return ErrHomestayNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Username, " ") == "" {
			return ErrUsernameRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.WaPhone, " ") == "" {
			return ErrWaPhoneRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.OtherPhone, " ") == "" {
			return ErrOtherPhoneRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayAddress, " ") == "" {
			return ErrHomestayAddresRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLatitude, " ") == "" {
			return ErrLatitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLongitude, " ") == "" {
			return ErrLongitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Password, " ") == "" {
			return ErrPasswordRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateLoginIn(i LoginIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.Identifier, " ") == "" {
			return errors.New("identifier required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Password, " ") == "" {
			return ErrPasswordRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateEditMemberIn(i EditMemberIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return ErrMemberNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayName, " ") == "" {
			return ErrHomestayNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Username, " ") == "" {
			return ErrUsernameRequired
		}
		return nil
	})
	g.Go(func() error {
		if len(i.PositionIds) == 0 {
			return ErrPositionRequired
		}
		return nil
	})
	g.Go(func() error {
		if i.PeriodId == 0 {
			return ErrOrgPeriodRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.WaPhone, " ") == "" {
			return ErrWaPhoneRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.OtherPhone, " ") == "" {
			return ErrOtherPhoneRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayAddress, " ") == "" {
			return ErrHomestayAddresRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLatitude, " ") == "" {
			return ErrLatitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLongitude, " ") == "" {
			return ErrLongitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Password, " ") != "" {
			// TODO: Add Validation rule
		}
		return nil
	})
	g.Go(func() error {
		if !i.IsAdmin.Valid {
			return ErrIsAdminRequired
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func ValidateUpdateProfileIn(i UpdateProfileIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return ErrMemberNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayName, " ") == "" {
			return ErrHomestayNameRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Username, " ") == "" {
			return ErrUsernameRequired
		}
		return nil
	})

	g.Go(func() error {
		if strings.Trim(i.WaPhone, " ") == "" {
			return ErrWaPhoneRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.OtherPhone, " ") == "" {
			return ErrOtherPhoneRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayAddress, " ") == "" {
			return ErrHomestayAddresRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLatitude, " ") == "" {
			return ErrLatitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLongitude, " ") == "" {
			return ErrLongitudeRequired
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Password, " ") != "" {
			// TODO: Add Validation rule
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
