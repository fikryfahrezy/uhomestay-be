package user

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/pagination"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/fikryfahrezy/crypt/agron2"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/jwt"
	"github.com/gofrs/uuid"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
	"gopkg.in/guregu/null.v4"
)

type (
	AddMemberIn struct {
		PositionId        int64                 `mapstructure:"position_id"`
		PeriodId          int64                 `mapstructure:"period_id"`
		Name              string                `mapstructure:"name"`
		Username          string                `mapstructure:"username"`
		Password          string                `mapstructure:"password"`
		WaPhone           string                `mapstructure:"wa_phone"`
		OtherPhone        string                `mapstructure:"other_phone"`
		HomestayName      string                `mapstructure:"homestay_name"`
		HomestayAddress   string                `mapstructure:"homestay_address"`
		HomestayLatitude  string                `mapstructure:"homestay_latitude"`
		HomestayLongitude string                `mapstructure:"homestay_longitude"`
		IsAdmin           null.Bool             `mapstructure:"is_admin"`
		File              httpdecode.FileHeader `mapstructure:"profile"`
	}
	AddMemberRes struct {
		Id string `json:"id"`
	}
	AddMemberOut struct {
		resp.Response
		Res AddMemberRes
	}
)

func (d *UserDeps) MemberSaver(ctx context.Context, in AddMemberIn, isApproved bool) (out AddMemberOut) {
	var err error

	_, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if !ok {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(errors.New("trx required"), "trx"))
		return
	}

	var periodId uint64
	if in.PeriodId != 0 {
		periodId, err = strconv.ParseUint(strconv.FormatInt(in.PeriodId, 10), 10, 64)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find parse by id"))
			return
		}

		_, err = d.OrgPeriodRepository.FindUndeletedById(ctx, periodId)
		if errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find period by id"))
			return
		}
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find period by id"))
			return
		}
	}

	var positionId uint64
	if in.PositionId != 0 {
		positionId, err = strconv.ParseUint(strconv.FormatInt(in.PositionId, 10), 10, 64)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "parse period id"))
			return
		}

		_, err = d.PositionRepository.FindUndeletedById(ctx, positionId)
		if errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find position by id"))
			return
		}
		if err != nil {
			out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "find position by id"))
			return
		}
	}

	member := MemberModel{
		Name:              in.Name,
		OtherPhone:        in.OtherPhone,
		WaPhone:           in.WaPhone,
		HomestayName:      in.HomestayName,
		HomestayAddress:   in.HomestayAddress,
		HomestayLatitude:  in.HomestayLatitude,
		HomestayLongitude: in.HomestayLongitude,
		Username:          in.Username,
		IsAdmin:           in.IsAdmin.Bool,
		IsApproved:        isApproved,
	}
	existingMember, err := d.MemberRepository.CheckUniqueField(ctx, member)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "check unique field"))
		return
	}

	if !existingMember.Id.UUID.IsNil() {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(errors.New("unique field already exist"), "duplication"))
		return
	}

	uid, err := uuid.NewV6()
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "generate uuid"))
		return
	}

	memberId := uid.String()
	if err = member.Id.Scan(memberId); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "scaning uid"))
		return
	}

	hash, err := agron2.Argon2Hash(in.Password, d.Argon2Salt, 1, 64*1024, 4, 32, argon2.Version, agron2.Argon2Id)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "hashing password"))
		return
	}

	member.Password = hash

	file := in.File.File

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileUrl string
	if file != nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + strings.Trim(in.File.Filename, " ")
		if fileUrl, err = d.Upload(filename, file); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "upload file"))
			return
		}
	}

	member.ProfilePicUrl = fileUrl

	if err = d.MemberRepository.Save(ctx, member); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save member"))
		return
	}

	if positionId != 0 && periodId != 0 {
		orgStructure := OrgStructureModel{
			MemberId:    memberId,
			PositionId:  positionId,
			OrgPeriodId: periodId,
		}

		if err = d.OrgStructureRepository.Save(ctx, orgStructure); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save sturcture"))
			return
		}
	}

	out.Res.Id = memberId

	return
}

func (d *UserDeps) AddMember(ctx context.Context, in AddMemberIn) (out AddMemberOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddMemberIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "member register validation"))
		return
	}

	saverOut := d.MemberSaver(ctx, in, true)
	if saverOut.Error != nil {
		out.Response = saverOut.Response
		return
	}

	out.Res.Id = saverOut.Res.Id

	return
}

type (
	RegisterIn struct {
		Name              string `json:"name"`
		Username          string `json:"username"`
		Password          string `json:"password"`
		WaPhone           string `json:"wa_phone"`
		OtherPhone        string `json:"other_phone"`
		HomestayName      string `json:"homestay_name"`
		HomestayAddress   string `json:"homestay_address"`
		HomestayLatitude  string `json:"homestay_latitude"`
		HomestayLongitude string `json:"homestay_longitude"`
	}
	RegisterRes struct {
		Token string `json:"token"`
	}
	RegisterOut struct {
		resp.Response
		Res RegisterRes
	}
)

func (d *UserDeps) MemberRegister(ctx context.Context, in RegisterIn) (out RegisterOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateRegisterIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "member register validation"))
		return
	}

	saverIn := AddMemberIn{
		Name:              in.Name,
		HomestayName:      in.HomestayName,
		HomestayAddress:   in.HomestayAddress,
		HomestayLatitude:  in.HomestayLatitude,
		HomestayLongitude: in.HomestayLongitude,
		WaPhone:           in.WaPhone,
		OtherPhone:        in.OtherPhone,
		Username:          in.Username,
		Password:          in.Password,
		IsAdmin:           null.BoolFrom(false),
	}

	saverOut := d.MemberSaver(ctx, saverIn, false)
	if saverOut.Error != nil {
		out.Response = saverOut.Response
		return
	}

	jwtToken, err := jwt.Sign(
		"",
		"token",
		d.JwtIssuerUrl,
		d.JwtKey,
		d.JwtAudiences,
		time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Time{},
		time.Time{},
		jwt.JwtPrivateClaim{
			Uid: saverOut.Res.Id,
		})
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "jwt signer"))
		return
	}

	out.Res.Token = jwtToken

	return
}

type (
	LoginIn struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	LoginRes struct {
		Token string `json:"token"`
	}
	LoginOut struct {
		resp.Response
		Res LoginRes
	}
)

func (d *UserDeps) MemberLogin(ctx context.Context, in LoginIn) (out LoginOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	if err = ValidateLoginIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "member login validation"))
		return
	}

	member, err := d.MemberRepository.FindByUsername(in.Identifier)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member by username"))
		return
	}

	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by username"))
		return
	}

	if !member.IsApproved {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.New("member not approved yet"))
		return
	}

	if err = agron2.Argon2Verify(member.Password, in.Password, agron2.Argon2Id); err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "verify password"))
		return
	}

	jwtToken, err := jwt.Sign(
		"",
		"token",
		d.JwtIssuerUrl,
		d.JwtKey,
		d.JwtAudiences,
		time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Time{},
		time.Time{},
		jwt.JwtPrivateClaim{
			Uid: member.Id.UUID.String(),
		})
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "jwt signer"))
		return
	}

	out.Res.Token = jwtToken

	return
}

func (d *UserDeps) AdminLogin(ctx context.Context, in LoginIn) (out LoginOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	if err = ValidateLoginIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "member login validation"))
		return
	}

	member, err := d.MemberRepository.FindAdminByByUsername(in.Identifier)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find admin member by username"))
		return
	}

	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by username"))
		return
	}

	if err = agron2.Argon2Verify(member.Password, in.Password, agron2.Argon2Id); err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "verify password"))
		return
	}

	jwtToken, err := jwt.Sign(
		"",
		"token",
		d.JwtIssuerUrl,
		d.JwtKey,
		d.JwtAudiences,
		time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Time{},
		time.Time{},
		jwt.JwtPrivateAdminClaim{
			Uid:     member.Id.UUID.String(),
			IsAdmin: true,
		})
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "jwt signer"))
		return
	}

	out.Res.Token = jwtToken

	return
}

type (
	EditMemberIn struct {
		Name              string                `mapstructure:"name"`
		Username          string                `mapstructure:"username"`
		Password          string                `mapstructure:"password"`
		WaPhone           string                `mapstructure:"wa_phone"`
		OtherPhone        string                `mapstructure:"other_phone"`
		HomestayName      string                `mapstructure:"homestay_name"`
		HomestayAddress   string                `mapstructure:"homestay_address"`
		HomestayLatitude  string                `mapstructure:"homestay_latitude"`
		HomestayLongitude string                `mapstructure:"homestay_longitude"`
		IsAdmin           null.Bool             `mapstructure:"is_admin"`
		PositionId        int64                 `mapstructure:"position_id"`
		PeriodId          int64                 `mapstructure:"period_id"`
		File              httpdecode.FileHeader `mapstructure:"profile"`
	}
	EditMemberRes struct {
		Id string `json:"id"`
	}
	EditMemberOut struct {
		resp.Response
		Res EditMemberRes
	}
)

func (d *UserDeps) EditMember(ctx context.Context, uid string, in EditMemberIn) (out EditMemberOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "member id not valid"))
		return
	}

	if err = ValidateEditMemberIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit member validation"))
		return
	}

	member, err := d.MemberRepository.FindById(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member by id"))
		return
	}

	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by id"))
		return
	}

	orgStructure, err := d.OrgStructureRepository.FindLatestByMemberId(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find user org structure by member id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find user org structure by member id"))
		return
	}

	var periodId uint64
	periodId, err = strconv.ParseUint(strconv.FormatInt(in.PeriodId, 10), 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "parse period id"))
		return
	}

	if periodId != orgStructure.OrgPeriodId {
		_, err = d.OrgPeriodRepository.FindById(ctx, periodId)
		if errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find period by id"))
			return
		}
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find period by id"))
			return
		}
	}

	var positionId uint64
	positionId, err = strconv.ParseUint(strconv.FormatInt(in.PositionId, 10), 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "parse period id"))
		return
	}

	if positionId != orgStructure.PositionId {
		_, err = d.PositionRepository.FindById(ctx, positionId)
		if errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find position by id"))
			return
		}
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find position by id"))
			return
		}
	}

	member.Name = in.Name
	member.OtherPhone = in.OtherPhone
	member.WaPhone = in.WaPhone
	member.HomestayName = in.HomestayName
	member.HomestayAddress = in.HomestayAddress
	member.HomestayLatitude = in.HomestayLatitude
	member.HomestayLongitude = in.HomestayLongitude
	member.Username = in.Username
	member.IsAdmin = in.IsAdmin.Bool

	existingMember, err := d.MemberRepository.CheckOtherUniqueField(ctx, uid, member)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "check other unique field"))
		return
	}

	if !existingMember.Id.UUID.IsNil() {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(errors.New("unique field already exist"), "duplication"))
		return
	}

	if in.Password != "" {
		hash, err := agron2.Argon2Hash(in.Password, d.Argon2Salt, 1, 64*1024, 4, 32, argon2.Version, agron2.Argon2Id)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "hashing password"))
			return
		}

		member.Password = hash
	}

	file := in.File.File

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileUrl string
	if file != nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + strings.Trim(in.File.Filename, " ")
		if fileUrl, err = d.Upload(filename, file); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "upload file"))
			return
		}
	}

	if fileUrl != "" {
		member.ProfilePicUrl = fileUrl
	}

	if err = d.MemberRepository.Update(ctx, uid, member); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update member"))
		return
	}

	if periodId != orgStructure.OrgPeriodId || positionId != orgStructure.PositionId {
		orgStructure = OrgStructureModel{
			MemberId:    member.Id.UUID.String(),
			PositionId:  positionId,
			OrgPeriodId: periodId,
		}
		if err = d.OrgStructureRepository.Save(ctx, orgStructure); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save structure"))
			return
		}
	}

	out.Res.Id = uid

	return
}

type (
	RemoveMemberRes struct {
		Id string `json:"id"`
	}
	RemoveMemberOut struct {
		resp.Response
		Res RemoveMemberRes
	}
)

func (d *UserDeps) RemoveMember(ctx context.Context, uid string) (out RemoveMemberOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "member id not valid"))
		return
	}

	_, err = d.MemberRepository.FindById(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "find member by id"))
		return
	}

	if err = d.MemberRepository.DeleteById(ctx, uid); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete member by id"))
		return
	}

	out.Res.Id = uid

	return
}

type (
	MemberOut struct {
		Id                string `json:"id"`
		Username          string `json:"username"`
		Name              string `json:"name"`
		WaPhone           string `json:"wa_phone"`
		OtherPhone        string `json:"other_phone"`
		HomestayName      string `json:"homestay_name"`
		HomestayAddress   string `json:"homestay_address"`
		HomestayLatitude  string `json:"homestay_latitude"`
		HomestayLongitude string `json:"homestay_longitude"`
		ProfilePicUrl     string `json:"profile_pic_url"`
		IsAdmin           bool   `json:"is_admin"`
		IsApproved        bool   `json:"is_approved"`
	}
	QueryMemberRes struct {
		Cursor  string      `json:"cursor"`
		Members []MemberOut `json:"members"`
	}
	QueryMemberOut struct {
		resp.Response
		Res QueryMemberRes
	}
)

func (d *UserDeps) QueryMember(ctx context.Context, q, cursor, limit string) (out QueryMemberOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	s, t, err := pagination.DecodeSIDCursor(cursor)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "decode sid cursor"))
		return
	}

	var uid pgtypeuuid.UUID
	uid.Scan(s)

	nlimit, _ := strconv.ParseInt(limit, 10, 64)
	if nlimit == 0 {
		nlimit = 25
	}

	members, err := d.MemberRepository.Query(ctx, uid, q, t, nlimit)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query member"))
		return
	}

	mLen := len(members)

	var nextCursor string
	if mLen != 0 {
		md := members[mLen-1]
		nextCursor = pagination.EncodeSIDCursor(md.Id.UUID.String(), md.CreatedAt)
	}

	outMembers := make([]MemberOut, mLen)
	for i, m := range members {
		outMembers[i] = MemberOut{
			Id:                m.Id.UUID.String(),
			Name:              m.Name,
			WaPhone:           m.WaPhone,
			OtherPhone:        m.OtherPhone,
			HomestayName:      m.HomestayName,
			HomestayAddress:   m.HomestayAddress,
			HomestayLatitude:  m.HomestayLatitude,
			HomestayLongitude: m.HomestayLongitude,
			ProfilePicUrl:     m.ProfilePicUrl,
			Username:          m.Username,
			IsAdmin:           m.IsAdmin,
			IsApproved:        m.IsApproved,
		}
	}

	out.Res = QueryMemberRes{
		Cursor:  nextCursor,
		Members: outMembers,
	}

	return
}

type (
	MemberDetailRes struct {
		Id                string `json:"id"`
		Name              string `json:"name"`
		Username          string `json:"username"`
		WaPhone           string `json:"wa_phone"`
		OtherPhone        string `json:"other_phone"`
		HomestayName      string `json:"homestay_name"`
		HomestayAddress   string `json:"homestay_address"`
		HomestayLatitude  string `json:"homestay_latitude"`
		HomestayLongitude string `json:"homestay_longitude"`
		ProfilePicUrl     string `json:"profile_pic_url"`
		IsAdmin           bool   `json:"is_admin"`
		IsApproved        bool   `json:"is_approved"`
		PeriodId          uint64 `json:"period_id"`
		Period            string `json:"period"`
		PositionId        uint64 `json:"position_id"`
		Position          string `json:"position"`
	}
	FindMemberDetailOut struct {
		resp.Response
		Res MemberDetailRes
	}
)

func (d *UserDeps) FindMemberDetail(ctx context.Context, uid string) (out FindMemberDetailOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "member id not valid"))
		return
	}

	mc := make(chan MemberModel)
	mR := make(chan resp.Response)
	go func(ctx context.Context, uid string, m chan MemberModel, res chan resp.Response) {
		var member MemberModel
		var r resp.Response

		member, r = func(ctx context.Context, uid string) (MemberModel, resp.Response) {
			member, err := d.MemberRepository.FindById(ctx, uid)
			if errors.Is(err, pgx.ErrNoRows) {
				return MemberModel{}, resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member by id"))
			}

			if err != nil {
				return MemberModel{}, resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by id"))
			}

			return member, resp.Response{}
		}(ctx, uid)

		m <- member
		res <- r
	}(ctx, uid, mc, mR)

	pc := make(chan OrgPeriodModel)
	psc := make(chan PositionModel)
	rs := make(chan resp.Response)
	go func(ctx context.Context, uid string, pc chan OrgPeriodModel, psc chan PositionModel, res chan resp.Response) {
		var period OrgPeriodModel
		var position PositionModel
		var r resp.Response

		period, position, r = func(ctx context.Context, uid string) (OrgPeriodModel, PositionModel, resp.Response) {
			var period OrgPeriodModel
			var position PositionModel

			orgStructure, err := d.OrgStructureRepository.FindLatestByMemberId(ctx, uid)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return OrgPeriodModel{}, PositionModel{}, resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find user org structure by member id"))
			}

			if errors.Is(err, pgx.ErrNoRows) {
				err = nil
			}

			if orgStructure.PositionId != 0 {
				position, err = d.PositionRepository.FindById(ctx, orgStructure.PositionId)
				if err != nil && !errors.Is(err, pgx.ErrNoRows) {
					return OrgPeriodModel{}, PositionModel{}, resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find position by id"))
				}
			}

			if orgStructure.OrgPeriodId != 0 {
				period, err = d.OrgPeriodRepository.FindById(ctx, orgStructure.OrgPeriodId)
				if err != nil && !errors.Is(err, pgx.ErrNoRows) {
					return OrgPeriodModel{}, PositionModel{}, resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find period by id"))
				}
			}

			return period, position, resp.Response{}
		}(ctx, uid)

		pc <- period
		psc <- position
		res <- r
	}(ctx, uid, pc, psc, rs)

	member := <-mc
	mRV := <-mR
	period := <-pc
	position := <-psc
	rsV := <-rs

	if mRV.Error != nil {
		out.Response = mRV
		return
	}

	if rsV.Error != nil {
		out.Response = rsV
		return
	}

	out.Res = MemberDetailRes{
		Id:                member.Id.UUID.String(),
		Name:              member.Name,
		WaPhone:           member.WaPhone,
		OtherPhone:        member.OtherPhone,
		HomestayName:      member.HomestayName,
		HomestayAddress:   member.HomestayAddress,
		HomestayLatitude:  member.HomestayLatitude,
		HomestayLongitude: member.HomestayLongitude,
		ProfilePicUrl:     member.ProfilePicUrl,
		Username:          member.Username,
		IsAdmin:           member.IsAdmin,
		IsApproved:        member.IsApproved,
		PeriodId:          period.Id,
		Period:            period.StartDate.Format("2006-01-02") + "|" + period.EndDate.Format("2006-01-02"),
		PositionId:        position.Id,
		Position:          position.Name,
	}

	return
}

type (
	MemberApprovalRes struct {
		Id string `json:"id"`
	}
	MemberApprovalOut struct {
		resp.Response
		Res MemberApprovalRes
	}
)

func (d *UserDeps) ApproveMember(ctx context.Context, uid string) (out MemberApprovalOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "member id not valid"))
		return
	}

	member, err := d.MemberRepository.FindById(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member by id"))
		return
	}

	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find unapproved member by id"))
		return
	}

	if member.IsApproved {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.New("member already approved"))
		return
	}

	member.IsApproved = true
	if err = d.MemberRepository.Update(ctx, uid, member); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update member"))
		return
	}

	out.Res.Id = uid

	return
}

type (
	UpdateProfileIn struct {
		Name              string                `mapstructure:"name"`
		Username          string                `mapstructure:"username"`
		Password          string                `mapstructure:"password"`
		WaPhone           string                `mapstructure:"wa_phone"`
		OtherPhone        string                `mapstructure:"other_phone"`
		HomestayName      string                `mapstructure:"homestay_name"`
		HomestayAddress   string                `mapstructure:"homestay_address"`
		HomestayLatitude  string                `mapstructure:"homestay_latitude"`
		HomestayLongitude string                `mapstructure:"homestay_longitude"`
		File              httpdecode.FileHeader `mapstructure:"profile"`
	}
	UpdateProfileRes struct {
		Id string `json:"id"`
	}
	UpdateProfileOut struct {
		resp.Response
		Res UpdateProfileRes
	}
)

func (d *UserDeps) UpdatProfile(ctx context.Context, uid string, in UpdateProfileIn) (out UpdateProfileOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "member id not valid"))
		return
	}

	if err = ValidateUpdateProfileIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit member validation"))
		return
	}

	member, err := d.MemberRepository.FindById(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member by id"))
		return
	}

	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by id"))
		return
	}

	member.Name = in.Name
	member.OtherPhone = in.OtherPhone
	member.WaPhone = in.WaPhone
	member.HomestayName = in.HomestayName
	member.HomestayAddress = in.HomestayAddress
	member.HomestayLatitude = in.HomestayLatitude
	member.HomestayLongitude = in.HomestayLongitude
	member.Username = in.Username

	existingMember, err := d.MemberRepository.CheckOtherUniqueField(ctx, uid, member)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "check other unique field"))
		return
	}

	if !existingMember.Id.UUID.IsNil() {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(errors.New("unique field already exist"), "duplication"))
		return
	}

	if in.Password != "" {
		hash, err := agron2.Argon2Hash(in.Password, d.Argon2Salt, 1, 64*1024, 4, 32, argon2.Version, agron2.Argon2Id)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "hashing password"))
			return
		}

		member.Password = hash
	}

	file := in.File.File

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileUrl string
	if file != nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + strings.Trim(in.File.Filename, " ")
		if fileUrl, err = d.Upload(filename, file); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "upload file"))
			return
		}
	}

	if fileUrl != "" {
		member.ProfilePicUrl = fileUrl
	}

	if err = d.MemberRepository.Update(ctx, uid, member); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update member"))
		return
	}

	out.Res.Id = uid

	return
}
