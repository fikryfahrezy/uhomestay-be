package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type (
	AddGoalIn struct {
		Vision      string `json:"vision"`
		VisionText  string `json:"vision_text"`
		Mission     string `json:"mission"`
		MissionText string `json:"mission_text"`
		OrgPeriodId int64  `json:"org_period_id"`
	}
	AddGoalRes struct {
		Id int64 `json:"id"`
	}
	AddGoalOut struct {
		resp.Response
		Res AddGoalRes
	}
)

func (d *UserDeps) AddGoal(ctx context.Context, in AddGoalIn) (out AddGoalOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddGoalIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "add goal validaion"))
		return
	}

	_, err = d.OrgPeriodRepository.FindUndeletedById(ctx, uint64(in.OrgPeriodId))
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no rows find period by id"))
		return
	}

	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find period by id"))
		return
	}

	unmarshal := func(title, content string, m chan map[string]interface{}, res chan resp.Response) {
		var mv map[string]interface{}
		var r resp.Response
		if content != "" {
			err = json.Unmarshal([]byte(content), &mv)
			if err != nil {
				r = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "unmarshal "+title+" json"))
			}
		}

		m <- mv
		res <- r
	}

	nv := make(chan map[string]interface{})
	nerr := make(chan resp.Response)
	go unmarshal("vision", in.Vision, nv, nerr)

	nm := make(chan map[string]interface{})
	merr := make(chan resp.Response)
	go unmarshal("mission", in.Mission, nm, merr)

	nvV := <-nv
	nRV := <-nerr
	nmV := <-nm
	mRV := <-merr

	if nRV.Error != nil {
		out.Response = nRV
		return
	}

	if mRV.Error != nil {
		out.Response = mRV
		return
	}

	goal := GoalModel{
		Vision:      nvV,
		VisionText:  in.VisionText,
		Mission:     nmV,
		MissionText: in.MissionText,
		OrgPeriodId: uint64(in.OrgPeriodId),
	}

	if goal, err = d.GoalRepository.Save(ctx, goal); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save goal"))
		return
	}

	out.Res.Id = int64(goal.Id)

	return
}

type (
	FindOrgPeriodGoalRes struct {
		Id          int64  `json:"id"`
		Vision      string `json:"vision"`
		VisionText  string `json:"vision_text"`
		Mission     string `json:"mission"`
		MissionText string `json:"mission_text"`
	}
	FindOrgPeriodGoalOut struct {
		resp.Response
		Res FindOrgPeriodGoalRes
	}
)

func (d *UserDeps) FindOrgPeriodGoal(ctx context.Context, pid string) (out FindOrgPeriodGoalOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	var period OrgPeriodModel
	if id == 0 {
		period, err = d.OrgPeriodRepository.QueryActive(ctx)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query active period"))
			return
		}
	}

	if id != 0 {
		period, err = d.OrgPeriodRepository.FindUndeletedById(ctx, id)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query structure by period id"))
			return
		}
	}

	goal, err := d.GoalRepository.FindByOrgPeriodId(ctx, period.Id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find goal by org period id"))
		return
	}

	marshal := func(title string, content map[string]interface{}, mc chan []byte, res chan resp.Response) {
		var r resp.Response
		var m []byte
		if content != nil && len(content) != 0 {
			m, err = json.Marshal(content)
			if err != nil {
				r = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, title+" json marshal"))
			}
		}

		mc <- m
		res <- r
	}

	v := make(chan []byte)
	vR := make(chan resp.Response)
	go marshal("vision", goal.Vision, v, vR)

	m := make(chan []byte)
	mR := make(chan resp.Response)
	go marshal("mission", goal.Mission, m, mR)

	vV := <-v
	vRV := <-vR
	mV := <-m
	mRV := <-mR

	if vRV.Error != nil {
		out.Response = vRV
		return
	}

	if mRV.Error != nil {
		out.Response = mRV
		return
	}

	out.StatusCode = http.StatusOK
	out.Res = FindOrgPeriodGoalRes{
		Id:          int64(goal.Id),
		Vision:      string(vV),
		Mission:     string(mV),
		VisionText:  goal.VisionText,
		MissionText: goal.MissionText,
	}

	return
}
