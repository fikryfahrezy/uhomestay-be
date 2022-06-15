package user

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func ValidateAddMemberIn(i AddMemberIn) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		if strings.Trim(i.Name, " ") == "" {
			return errors.New("name required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayName, " ") == "" {
			return errors.New("homestay_name required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Username, " ") == "" {
			return errors.New("username required")
		}
		return nil
	})
	g.Go(func() error {
		if i.PositionId == 0 {
			return errors.New("position_id required")
		}
		return nil
	})
	g.Go(func() error {
		if i.PeriodId == 0 {
			return errors.New("period_id required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.WaPhone, " ") == "" {
			return errors.New("wa_phone required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.OtherPhone, " ") == "" {
			return errors.New("other_phone required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayAddress, " ") == "" {
			return errors.New("homestay_address required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLatitude, " ") == "" {
			return errors.New("homestay_latitude required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLongitude, " ") == "" {
			return errors.New("homestay_longitude required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Password, " ") == "" {
			return errors.New("password required")
		}
		return nil
	})
	g.Go(func() error {
		if !i.IsAdmin.Valid {
			return errors.New("is_admin required")
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
			return errors.New("name required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayName, " ") == "" {
			return errors.New("homestay_name required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Username, " ") == "" {
			return errors.New("username required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.WaPhone, " ") == "" {
			return errors.New("wa_phone required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.OtherPhone, " ") == "" {
			return errors.New("other_phone required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayAddress, " ") == "" {
			return errors.New("homestay_address required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLatitude, " ") == "" {
			return errors.New("homestay_latitude required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLongitude, " ") == "" {
			return errors.New("homestay_longitude required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Password, " ") == "" {
			return errors.New("password required")
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
			return errors.New("password required")
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
			return errors.New("name required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayName, " ") == "" {
			return errors.New("homestay_name required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Username, " ") == "" {
			return errors.New("username required")
		}
		return nil
	})
	g.Go(func() error {
		if i.PositionId == 0 {
			return errors.New("position_id required")
		}
		return nil
	})
	g.Go(func() error {
		if i.PeriodId == 0 {
			return errors.New("period_id required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.WaPhone, " ") == "" {
			return errors.New("wa_phone required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.OtherPhone, " ") == "" {
			return errors.New("other_phone required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayAddress, " ") == "" {
			return errors.New("homestay_address required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLatitude, " ") == "" {
			return errors.New("homestay_latitude required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLongitude, " ") == "" {
			return errors.New("homestay_longitude required")
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
			return errors.New("is_admin required")
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
			return errors.New("name required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayName, " ") == "" {
			return errors.New("homestay_name required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.Username, " ") == "" {
			return errors.New("username required")
		}
		return nil
	})

	g.Go(func() error {
		if strings.Trim(i.WaPhone, " ") == "" {
			return errors.New("wa_phone required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.OtherPhone, " ") == "" {
			return errors.New("other_phone required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayAddress, " ") == "" {
			return errors.New("homestay_address required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLatitude, " ") == "" {
			return errors.New("homestay_latitude required")
		}
		return nil
	})
	g.Go(func() error {
		if strings.Trim(i.HomestayLongitude, " ") == "" {
			return errors.New("homestay_longitude required")
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
