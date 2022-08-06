package homestay

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

var (
	ErrMemberHomestayNotFound = errors.New("homestay anggota tidak ditemukan")
	ErrMemberNotFound         = errors.New("anggota tidak ditemukan")
)

type (
	AddMemberHomestayIn struct {
		Name      string  `json:"name"`
		Address   string  `json:"address"`
		Latitude  string  `json:"latitude"`
		Longitude string  `json:"longitude"`
		ImageIds  []int64 `json:"image_ids"`
	}
	AddMemberHomestayRes struct {
		Id int64 `json:"id"`
	}
	AddMemberHomestayOut struct {
		resp.Response
		Res AddMemberHomestayRes
	}
)

func (d *HomestayDeps) AddMemberHomestay(ctx context.Context, uid string, in AddMemberHomestayIn) (out AddMemberHomestayOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberNotFound)
		return
	}

	if err = ValidateAddMemberHomestayIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", err)
		return
	}

	_, err = d.MemberRepository.FindById(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by uid"))
		return
	}

	var newImageIds []uint64
	for _, v := range in.ImageIds {
		if v < 1 {
			continue
		}

		newImageIds = append(newImageIds, uint64(v))
	}

	homestayImages, err := d.HomestayImageRepository.QueryInId(ctx, newImageIds)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query homestay in id"))
		return
	}

	if len(homestayImages) == 0 {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", ErrHomestayImageRequired)
		return
	}

	memberHomestay := MemberHomestayModel{
		Name:         in.Name,
		Address:      in.Address,
		Latitude:     in.Latitude,
		Longitude:    in.Longitude,
		ThumbnailUrl: homestayImages[0].Url,
		MemberId:     uid,
	}

	if memberHomestay, err = d.MemberHomestayRepository.Save(ctx, memberHomestay); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save member homestay"))
		return
	}

	var newHomestayImagesIds []uint64
	for _, v := range homestayImages {
		newHomestayImagesIds = append(newHomestayImagesIds, v.Id)
	}

	if err := d.HomestayImageRepository.UpdateHomestayIdInId(ctx, memberHomestay.Id, newHomestayImagesIds); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update homestay images homestay id"))
		return
	}

	out.Res.Id = int64(memberHomestay.Id)

	return
}

type (
	EditMemberHomestayIn struct {
		Name      string  `json:"name"`
		Address   string  `json:"address"`
		Latitude  string  `json:"latitude"`
		Longitude string  `json:"longitude"`
		ImageIds  []int64 `json:"image_ids"`
	}
	EditMemberHomestayRes struct {
		Id int64 `json:"id"`
	}
	EditMemberHomestayOut struct {
		resp.Response
		Res EditMemberHomestayRes
	}
)

func (d *HomestayDeps) EditMemberHomestay(ctx context.Context, pid, uid string, in EditMemberHomestayIn) (out EditMemberHomestayOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberNotFound)
		return
	}

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberHomestayNotFound)
		return
	}

	if err = ValidateEditMemberHomestayIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", err)
		return
	}

	_, err = d.MemberRepository.FindById(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by uid"))
		return
	}

	memberHomestay, err := d.MemberHomestayRepository.FindById(ctx, uid, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberHomestayNotFound)
		return
	}

	var newImageIds []uint64
	for _, v := range in.ImageIds {
		if v < 1 {
			continue
		}

		newImageIds = append(newImageIds, uint64(v))
	}

	homestayImages, err := d.HomestayImageRepository.QueryInId(ctx, newImageIds)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query homestay in id"))
		return
	}

	if len(homestayImages) == 0 {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", ErrHomestayImageRequired)
		return
	}

	memberHomestay.Name = in.Name
	memberHomestay.Address = in.Address
	memberHomestay.Latitude = in.Latitude
	memberHomestay.Longitude = in.Longitude
	memberHomestay.ThumbnailUrl = homestayImages[0].Url

	if err = d.MemberHomestayRepository.UpdateById(ctx, uid, id, memberHomestay); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update member homestay by id"))
		return
	}

	var newHomestayImagesIds []uint64
	for _, v := range homestayImages {
		newHomestayImagesIds = append(newHomestayImagesIds, v.Id)
	}

	if err := d.HomestayImageRepository.UpdateHomestayIdInId(ctx, memberHomestay.Id, newHomestayImagesIds); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update homestay images homestay id"))
	}

	oldHomestayImages, err := d.HomestayImageRepository.FindByMemberHomestayId(ctx, id)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query images"))
		return
	}

	if len(oldHomestayImages) != 0 {
		var homestayImgIds []uint64
		for _, v := range oldHomestayImages {
			var isExist bool
			for _, q := range newImageIds {
				isExist = q == v.Id

				if isExist {
					break
				}
			}

			if !isExist {
				homestayImgIds = append(homestayImgIds, v.Id)
			}
		}

		if err = d.HomestayImageRepository.DeleteInId(ctx, homestayImgIds); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete image by id"))
			return
		}
	}

	out.Res.Id = int64(memberHomestay.Id)

	return
}

type (
	MemberHomestaysRes struct {
		Id           int64  `json:"id"`
		Name         string `json:"name"`
		ThumbnailUrl string `json:"thumbnail_url"`
	}
	QueryMemberHomestayRes struct {
		Cursor          int64                `json:"cursor"`
		Total           int64                `json:"total"`
		MemberHomestays []MemberHomestaysRes `json:"member_homestays"`
	}
	QueryMemberHomestayImageOut struct {
		resp.Response
		Res QueryMemberHomestayRes
	}
)

func (d *HomestayDeps) QueryMemberHomestays(ctx context.Context, uid, cursor, limit string) (out QueryMemberHomestayImageOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	nlimit, _ := strconv.ParseInt(limit, 10, 64)
	if nlimit == 0 {
		nlimit = 25
	}

	memberHomestayNumber, err := d.MemberHomestayRepository.CountMemberHomestay(ctx, uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "count image"))
		return
	}

	memberHomestays, err := d.MemberHomestayRepository.Query(ctx, uid, fromCursor, nlimit)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query member homestays"))
		return
	}

	memberHomeLen := len(memberHomestays)

	var nextCursor int64
	if memberHomeLen != 0 {
		nextCursor = int64(memberHomestays[memberHomeLen-1].Id)
	}

	outMemberHomestays := make([]MemberHomestaysRes, memberHomeLen)
	for i, p := range memberHomestays {
		outMemberHomestays[i] = MemberHomestaysRes{
			Id:           int64(p.Id),
			Name:         p.Name,
			ThumbnailUrl: p.ThumbnailUrl,
		}
	}

	out.Res = QueryMemberHomestayRes{
		Cursor:          nextCursor,
		Total:           memberHomestayNumber,
		MemberHomestays: outMemberHomestays,
	}

	return
}

type (
	RemoveImageRes struct {
		Id int64 `json:"id"`
	}
	RemoveImageOut struct {
		resp.Response
		Res RemoveImageRes
	}
)

func (d *HomestayDeps) RemoveMemberHomestay(ctx context.Context, pid, uid string) (out RemoveImageOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberNotFound)
		return
	}

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberHomestayNotFound)
		return
	}

	_, err = d.MemberRepository.FindById(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by uid"))
		return
	}

	_, err = d.MemberHomestayRepository.FindById(ctx, uid, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberHomestayNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find image by id"))
		return
	}

	if err = d.MemberHomestayRepository.DeleteById(ctx, uid, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete image by id"))
		return
	}

	homestayImages, err := d.HomestayImageRepository.FindByMemberHomestayId(ctx, id)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query images"))
		return
	}

	if len(homestayImages) != 0 {
		var homestayImgIds []uint64
		for _, v := range homestayImages {
			homestayImgIds = append(homestayImgIds, v.Id)
		}

		if err = d.HomestayImageRepository.DeleteInId(ctx, homestayImgIds); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete image by id"))
			return
		}
	}

	out.Res.Id = int64(id)

	return
}

type (
	HomestayImageRes struct {
		Id  int64  `json:"id"`
		Url string `json:"url"`
	}
	MemberHomestayRes struct {
		Id             int64              `json:"id"`
		Name           string             `json:"name"`
		Address        string             `json:"address"`
		Latitude       string             `json:"latitude"`
		Longitude      string             `json:"longitude"`
		HomestayImages []HomestayImageRes `json:"images"`
	}
	MemberHomestayOut struct {
		resp.Response
		Res MemberHomestayRes
	}
)

func (d *HomestayDeps) FindMemberHomestay(ctx context.Context, pid, uid string) (out MemberHomestayOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberNotFound)
		return
	}

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberHomestayNotFound)
		return
	}

	_, err = d.MemberRepository.FindById(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by uid"))
		return
	}

	memberHomestay, err := d.MemberHomestayRepository.FindById(ctx, uid, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrMemberHomestayNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find image by id"))
		return
	}

	homestayImages, err := d.HomestayImageRepository.FindByMemberHomestayId(ctx, id)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query images"))
		return
	}

	newHomestayImages := make([]HomestayImageRes, len(homestayImages))
	for i, v := range homestayImages {
		newHomestayImages[i] = HomestayImageRes{
			Id:  int64(v.Id),
			Url: v.Url,
		}
	}

	out.Res = MemberHomestayRes{
		Id:             int64(memberHomestay.Id),
		Name:           memberHomestay.Name,
		Address:        memberHomestay.Address,
		Latitude:       memberHomestay.Latitude,
		Longitude:      memberHomestay.Longitude,
		HomestayImages: newHomestayImages,
	}

	return
}
