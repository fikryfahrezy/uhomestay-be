package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type (
	AddPositionIn struct {
		Name  string `json:"name"`
		Level int64  `json:"level"`
	}
	AddPositionRes struct {
		Id uint64 `json:"id"`
	}
	AddPositionOut struct {
		resp.Response
		Res AddPositionRes
	}
)

func (d *UserDeps) AddPosition(ctx context.Context, in AddPositionIn) (out AddPositionOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddPositionIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "add position validation"))
		return
	}

	position := PositionModel{
		Name:  in.Name,
		Level: int16(in.Level),
	}
	if position, err = d.PositionRepository.Save(ctx, position); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save position"))
		return
	}

	out.Res.Id = position.Id

	return
}

type (
	PositionOut struct {
		Level int16  `json:"level"`
		Id    uint64 `json:"id"`
		Name  string `json:"name"`
	}
	QueryPositionRes struct {
		Cursor    int64         `json:"cursor"`
		Positions []PositionOut `json:"positions"`
	}
	QueryPositionOut struct {
		resp.Response
		Res QueryPositionRes
	}
)

func (d *UserDeps) QueryPosition(ctx context.Context, cursor, limit string) (out QueryPositionOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	nlimit, _ := strconv.ParseInt(limit, 10, 64)
	if nlimit == 0 {
		nlimit = 25
	}

	positions, err := d.PositionRepository.Query(ctx, fromCursor, nlimit)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query positions"))
		return
	}

	posLen := len(positions)

	var nextCursor int64
	if posLen != 0 {
		nextCursor = int64(positions[posLen-1].Id)
	}

	outPoisitons := make([]PositionOut, posLen)
	for i, p := range positions {
		outPoisitons[i] = PositionOut{
			Id:    p.Id,
			Name:  p.Name,
			Level: p.Level,
		}
	}

	out.Res = QueryPositionRes{
		Cursor:    nextCursor,
		Positions: outPoisitons,
	}

	return
}

type (
	PositionLevelOut struct {
		Level int64 `json:"level"`
	}
	QueryPositionLevelOut struct {
		resp.Response
		Res []PositionLevelOut
	}
)

func (d *UserDeps) QueryPositionLevel(ctx context.Context) (out QueryPositionLevelOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	ml, err := d.PositionRepository.MaxLevel(ctx)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query position"))
		return
	}

	ml += 1
	outPLs := make([]PositionLevelOut, ml)
	for i := int64(0); i < ml; i++ {
		outPLs[i] = PositionLevelOut{
			Level: i + 1,
		}
	}

	out.Res = outPLs

	return
}

type (
	EditPositionIn struct {
		Name  string `json:"name"`
		Level int64  `json:"level"`
	}
	EditPositionRes struct {
		Id uint64 `json:"id"`
	}
	EditPositionOut struct {
		resp.Response
		Res EditPositionRes
	}
)

func (d *UserDeps) EditPosition(ctx context.Context, pid string, in EditPositionIn) (out EditPositionOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	if err = ValidateEditPositionIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit position validation"))
		return
	}

	level, err := strconv.ParseUint(strconv.FormatInt(in.Level, 10), 10, 16)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "parse period id"))
		return
	}

	position, err := d.PositionRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find position by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find position by id"))
		return
	}

	position.Name = in.Name
	position.Level = int16(level)

	if err = d.PositionRepository.UpdateById(ctx, position.Id, position); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update position by id"))
		return
	}

	out.Res.Id = id

	return
}

type (
	RemovePositionRes struct {
		Id uint64 `json:"id"`
	}
	RemovePositionOut struct {
		resp.Response
		Res RemovePositionRes
	}
)

func (d *UserDeps) RemovePosition(ctx context.Context, pid string) (out RemovePositionOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", errors.Wrap(err, "update position by id"))

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	_, err = d.PositionRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find position by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find position by id"))
		return
	}

	if err = d.PositionRepository.DeleteById(ctx, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete position by id"))
		return
	}

	out.Res.Id = id

	return
}
