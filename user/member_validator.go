package user

import (
	"strings"
	"unicode/utf8"

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
	ErrMaxName                = errors.New("nama anggota tidak dapat lebih dari 100 karakter")
	ErrMaxWaPhone             = errors.New("nomor whats app tidak dapat lebih dari 50 karakter")
	ErrMaxOtherPhone          = errors.New("nomor lainnya tidak dapat lebih dari 50 karakter")
	ErrMaxHomestayName        = errors.New("nama homestay anggota tidak dapat lebih dari 100 karakter")
	ErrMaxHomestayAddress     = errors.New("alamat homestay anggota tidak dapat lebih dari 200 karakter")
	ErrMaxHomestayLat         = errors.New("titik garis bujur map homestay tidak dapat lebih dari 50 karakter")
	ErrMaxHomestayLng         = errors.New("titik garis lintang map homestay tidak dapat lebih dari 50 karakter")
	ErrMaxUsername            = errors.New("username anggota tidak dapat lebih dari 50 karakter")
	ErrMaxPassword            = errors.New("password anggota tidak dapat lebih dari 200 karakter")
	ErrProfileRequired        = errors.New("foto profile tidak boleh kosong")
	ErrProfileFileName        = errors.New("nama file foto profile tidak dapat lebih dari 200 karakter")
	ErrIdCardRequired         = errors.New("ktp tidak boleh kosong")
	ErrIdCardFileName         = errors.New("nama file ktp tidak dapat lebih dari 200 karakter")
	ErrHomestayPhotoRequired  = errors.New("foto homestay tidak boleh kosong")
	ErrHomestayPhotoFileName  = errors.New("nama file foto homestay tidak dapat lebih dari 200 karakter")
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
	g.Go(func() error {
		if utf8.RuneCountInString(i.Name) > 100 {
			return ErrMaxName
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.WaPhone) > 50 {
			return ErrMaxWaPhone
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.OtherPhone) > 50 {
			return ErrMaxOtherPhone
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Username) > 50 {
			return ErrMaxUsername
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Password) > 200 {
			return ErrMaxPassword
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
	g.Go(func() error {
		if utf8.RuneCountInString(i.Name) > 100 {
			return ErrMaxName
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.WaPhone) > 50 {
			return ErrMaxWaPhone
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.OtherPhone) > 50 {
			return ErrMaxOtherPhone
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.HomestayName) > 100 {
			return ErrMaxHomestayName
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.HomestayAddress) > 200 {
			return ErrMaxHomestayAddress
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.HomestayLatitude) > 50 {
			return ErrMaxHomestayLat
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.HomestayLongitude) > 50 {
			return ErrMaxHomestayLng
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Username) > 50 {
			return ErrMaxUsername
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Password) > 200 {
			return ErrMaxPassword
		}
		return nil
	})
	g.Go(func() error {
		if i.Profile.File == nil || i.Profile.Filename == "" {
			return ErrProfileRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Profile.Filename) > 200 {
			return ErrProfileFileName
		}
		return nil
	})
	g.Go(func() error {
		if i.IdCard.File == nil || i.IdCard.Filename == "" {
			return ErrIdCardRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.IdCard.Filename) > 200 {
			return ErrIdCardFileName
		}
		return nil
	})
	g.Go(func() error {
		if i.HomestayPhoto.File == nil || i.HomestayPhoto.Filename == "" {
			return ErrHomestayPhotoRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.HomestayPhoto.Filename) > 200 {
			return ErrHomestayPhotoFileName
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
		if strings.Trim(i.Password, " ") != "" && utf8.RuneCountInString(i.Password) > 200 {
			return ErrMaxPassword
		}
		return nil
	})
	g.Go(func() error {
		if !i.IsAdmin.Valid {
			return ErrIsAdminRequired
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Name) > 100 {
			return ErrMaxName
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.WaPhone) > 50 {
			return ErrMaxWaPhone
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.OtherPhone) > 50 {
			return ErrMaxOtherPhone
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Username) > 50 {
			return ErrMaxUsername
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Profile.Filename) > 200 {
			return ErrProfileFileName
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.IdCard.Filename) > 200 {
			return ErrIdCardFileName
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
		if strings.Trim(i.Password, " ") != "" && utf8.RuneCountInString(i.Password) > 200 {
			return ErrMaxPassword
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Name) > 100 {
			return ErrMaxName
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.WaPhone) > 50 {
			return ErrMaxWaPhone
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.OtherPhone) > 50 {
			return ErrMaxOtherPhone
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Username) > 50 {
			return ErrMaxUsername
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.Profile.Filename) > 200 {
			return ErrProfileFileName
		}
		return nil
	})
	g.Go(func() error {
		if utf8.RuneCountInString(i.IdCard.Filename) > 200 {
			return ErrIdCardFileName
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
