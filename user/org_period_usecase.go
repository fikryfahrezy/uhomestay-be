package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/timediff"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

func (d *UserDeps) SaveOrgStructure(ctx context.Context, periodId uint64, ps []PositionIn) error {
	if len(ps) == 0 {
		return nil
	}

	var memberIds []string
	var positionIds []uint64
	for _, p := range ps {
		positionId, err := strconv.ParseUint(strconv.FormatInt(p.Id, 10), 10, 64)
		if err != nil {
			err = errors.Wrap(err, "parse period id")
			return err
		}

		positionIds = append(positionIds, positionId)
		for _, m := range p.Members {
			memberIds = append(memberIds, m.Id)
		}
	}

	cacher := make(map[interface{}]bool)

	members, err := d.MemberRepository.QueryInId(ctx, memberIds)
	if err != nil {
		err = errors.Wrap(err, "query member in id")
		return err
	}

	for _, m := range members {
		cacher[m.Id.UUID.String()] = true
	}

	positions, err := d.PositionRepository.QueryInId(ctx, positionIds)
	if err != nil {
		err = errors.Wrap(err, "query position in id")
		return err
	}

	for _, p := range positions {
		cacher[p.Id] = true
	}

	var structures []OrgStructureModel
	for _, p := range ps {
		orgStructure := OrgStructureModel{
			OrgPeriodId: periodId,
		}
		positionId, err := strconv.ParseUint(strconv.FormatInt(p.Id, 10), 10, 64)
		if err != nil {
			err = errors.Wrap(err, "parse period id")
			return err
		}

		if _, ok := cacher[positionId]; ok {
			orgStructure.PositionId = positionId
			for _, m := range p.Members {
				if _, ok = cacher[m.Id]; ok {
					orgStructure.MemberId = m.Id
					structures = append(structures, orgStructure)
				}
			}
		}
	}

	if len(structures) != 0 {
		if err := d.OrgStructureRepository.BulkSave(ctx, structures); err != nil {
			err = errors.Wrap(err, "save structure")
			return err
		}
	}

	return nil
}

type (
	MemberIn struct {
		Id string `json:"id"`
	}
	PositionIn struct {
		Id      int64      `json:"id"`
		Members []MemberIn `json:"members"`
	}
	AddPeriodIn struct {
		StartDate   string       `json:"start_date"`
		EndDate     string       `json:"end_date"`
		Positions   []PositionIn `json:"positions"`
		Vision      string       `json:"vision"`
		VisionText  string       `json:"vision_text"`
		Mission     string       `json:"mission"`
		MissionText string       `json:"mission_text"`
	}
	AddPeriodRes struct {
		Id uint64 `json:"id"`
	}
	AddPeriodOut struct {
		resp.Response
		Res AddPeriodRes
	}
)

func (d *UserDeps) AddPeriod(ctx context.Context, in AddPeriodIn) (out AddPeriodOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddPeriodIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "add period validation"))
		return
	}

	startDate, err := time.Parse("2006-01-02", in.StartDate)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "parse start date"))
		return
	}

	endDate, err := time.Parse("2006-01-02", in.EndDate)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "parse end date"))
		return
	}

	if isPast := timediff.Compare(endDate, startDate); isPast {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.New("end_date cannot lower than start_date"))
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

	period := OrgPeriodModel{
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  true,
	}
	if err = d.OrgPeriodRepository.DisableAll(ctx); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "disable all period"))
		return
	}

	if period, err = d.OrgPeriodRepository.Save(ctx, period); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save period"))
		return
	}

	goal := GoalModel{
		Vision:      nvV,
		VisionText:  in.VisionText,
		Mission:     nmV,
		MissionText: in.MissionText,
		OrgPeriodId: uint64(period.Id),
	}

	if goal, err = d.GoalRepository.Save(ctx, goal); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save goal"))
		return
	}

	if err = d.SaveOrgStructure(ctx, period.Id, in.Positions); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "build structure"))
		return
	}

	out.Res.Id = period.Id

	return
}

type (
	PeriodRes struct {
		IsActive  bool   `json:"is_active"`
		Id        uint64 `json:"id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	QueryPeriodRes struct {
		Cursor  int64       `json:"cursor"`
		Periods []PeriodRes `json:"periods"`
	}
	QueryPeriodOut struct {
		resp.Response
		Res QueryPeriodRes
	}
)

func (d *UserDeps) QueryPeriod(ctx context.Context, cursor string) (out QueryPeriodOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	periods, err := d.OrgPeriodRepository.Query(ctx, fromCursor, 25)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query period"))
		return
	}

	periodsLen := len(periods)

	var nextCursor int64
	if periodsLen != 0 {
		nextCursor = int64(periods[periodsLen-1].Id)
	}

	outPeriods := make([]PeriodRes, periodsLen)
	for i, p := range periods {
		outPeriods[i] = PeriodRes{
			Id:        p.Id,
			StartDate: p.StartDate.Format("2006-01-02"),
			EndDate:   p.EndDate.Format("2006-01-02"),
			IsActive:  p.IsActive,
		}
	}

	out.Res = QueryPeriodRes{
		Cursor:  nextCursor,
		Periods: outPeriods,
	}

	return
}

type (
	EditPeriodIn struct {
		StartDate   string       `json:"start_date"`
		EndDate     string       `json:"end_date"`
		Positions   []PositionIn `json:"positions"`
		Vision      string       `json:"vision"`
		VisionText  string       `json:"vision_text"`
		Mission     string       `json:"mission"`
		MissionText string       `json:"mission_text"`
	}
	EditPeriodRes struct {
		Id uint64 `json:"id"`
	}
	EditPeriodOut struct {
		resp.Response
		Res EditPeriodRes
	}
)

func (d *UserDeps) EditPeriod(ctx context.Context, pid string, in EditPeriodIn) (out EditPeriodOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	if err = ValidateEditPeriodIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit period validation"))
		return
	}

	period, err := d.OrgPeriodRepository.FindActiveBydId(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find period by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find period by id"))
		return
	}

	startDate, err := time.Parse("2006-01-02", in.StartDate)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "parse start date"))
		return
	}

	endDate, err := time.Parse("2006-01-02", in.EndDate)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "parse end date"))
		return
	}

	if isPast := timediff.Compare(endDate, startDate); isPast {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.New("end_date cannot lower than start_date"))
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

	period.StartDate = startDate
	period.EndDate = endDate

	err = d.OrgPeriodRepository.UpdateById(ctx, id, period)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update period by id"))
		return
	}

	if len(in.Positions) != 0 {
		if err = d.OrgStructureRepository.DeleteByPeriodId(ctx, id); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete org structure by period id"))
			return
		}
	}

	if err = d.SaveOrgStructure(ctx, period.Id, in.Positions); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "build structure"))
		return
	}

	if nv != nil || nm != nil {
		goal := GoalModel{
			Vision:      nvV,
			VisionText:  in.VisionText,
			Mission:     nmV,
			MissionText: in.MissionText,
			OrgPeriodId: uint64(period.Id),
		}

		if goal, err = d.GoalRepository.Save(ctx, goal); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save goal"))
			return
		}
	}

	out.Res.Id = id

	return
}

type (
	RemovePeriodRes struct {
		Id uint64 `json:"id"`
	}
	RemovePeriodOut struct {
		resp.Response
		Res RemovePeriodRes
	}
)

func (d *UserDeps) RemovePeriod(ctx context.Context, pid string) (out RemovePeriodOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	period, err := d.OrgPeriodRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find period by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find period by id"))
		return
	}

	isActiveBefore := period.IsActive
	out.Res.Id = id

	if err = d.OrgPeriodRepository.DeleteById(ctx, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete period by id"))
		return
	}

	if err = d.OrgPeriodRepository.UpdateStatusById(ctx, id, period); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update period status by id"))
		return
	}

	if err = d.OrgPeriodRepository.EnableOtherLastInactiveTx(ctx, id, isActiveBefore); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "enable other last period"))
		return
	}

	return
}

type (
	SwitchPeriodStatusIn struct {
		IsActive null.Bool `json:"is_active"`
	}
	SwitchPeriodStatusRes struct {
		Id uint64 `json:"id"`
	}
	SwitchPeriodStatusOut struct {
		resp.Response
		Res SwitchPeriodStatusRes
	}
)

func (d *UserDeps) SwitchPeriodStatus(ctx context.Context, pid string, in SwitchPeriodStatusIn) (out SwitchPeriodStatusOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	period, err := d.OrgPeriodRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find period by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find period by id"))
		return
	}

	period.IsActive = in.IsActive.Bool
	out.Res.Id = id

	if period.IsActive {
		if err = d.OrgPeriodRepository.DisableAll(ctx); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "enable other last period"))
			return
		}

		if err = d.OrgPeriodRepository.UpdateStatusById(ctx, id, period); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update period status by id"))
			return
		}

		return
	}

	if err = d.OrgPeriodRepository.UpdateStatusById(ctx, id, period); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "enable other last period"))
		return
	}

	isActiveBefore := period.IsActive

	if err = d.OrgPeriodRepository.EnableOtherLastInactiveTx(ctx, id, isActiveBefore); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "enable other last period"))
		return
	}

	return
}

type (
	FindActivePeriodOut struct {
		resp.Response
		Res PeriodRes
	}
)

func (d *UserDeps) FindActivePeriod(ctx context.Context) (out FindActivePeriodOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	period, err := d.OrgPeriodRepository.QueryActive(ctx)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query active period"))
		return
	}

	var startDate, endDate string
	if period.Id != 0 {
		startDate = period.StartDate.Format("2006-01-02")
		endDate = period.StartDate.Format("2006-01-02")
	}

	outPeriod := PeriodRes{
		Id:        period.Id,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  period.IsActive,
	}

	out.Res = outPeriod

	return
}

type (
	StructureMemberOut struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		ProfilePicUrl string `json:"profile_pic_url"`
	}
	StructurePositionOut struct {
		Id      uint64               `json:"id"`
		Name    string               `json:"name"`
		Level   int16                `json:"level"`
		Members []StructureMemberOut `json:"members"`
	}
	StructureRes struct {
		Id        uint64                 `json:"id"`
		StartDate string                 `json:"start_date"`
		EndDate   string                 `json:"end_date"`
		Positions []StructurePositionOut `json:"positions"`
		Vision    string                 `json:"vision"`
		Mission   string                 `json:"mission"`
	}
	StructureOut struct {
		resp.Response
		Res StructureRes
	}
)

func (d *UserDeps) QueryPeriodStructure(ctx context.Context, pid string) (out StructureOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find period by id"))
		return
	}

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	var period OrgPeriodModel
	if id == 0 {
		period, err = d.OrgPeriodRepository.QueryActive(ctx)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
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

	if period.Id == 0 {
		out.Res = StructureRes{
			Id:        0,
			StartDate: "",
			EndDate:   "",
			Positions: make([]StructurePositionOut, 0, 0),
			Vision:    "",
			Mission:   "",
		}
		return
	}

	structures, err := d.OrgStructureRepository.FindByPeriodId(ctx, period.Id)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query structure by period id"))
		return
	}

	goal, err := d.GoalRepository.FindByOrgPeriodId(ctx, period.Id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find goal by org period id"))
		return
	}

	tn := time.Now()

	// For save member id as key and the time of member id added to org structure as value
	// Because there are possibly multiple member ids in one structure
	// Used as a Set or Map of unique member id by the most recent
	membCIds := make(map[string]time.Duration)

	// For save member id as key and the org structure where it is in as value
	// Used as a Set or Map of unique member org structure by that member id
	newStruct := make(map[string]OrgStructureModel)

	// Grouping by latest member id
	// Fill the `membCIds` and `newStruct`
	for _, s := range structures {
		memId := s.MemberId
		ts := tn.Sub(s.CreatedAt)
		_, ok := membCIds[memId]
		if !ok {
			membCIds[memId] = ts
		}

		if ts < membCIds[memId] {
			membCIds[memId] = ts
		}

		newStruct[memId] = s
	}

	// Collect member ids and position ids from new struct
	// for query `in` to database, to check if all member ids and position ids
	// are exist in database
	memberIds := make([]string, 0, len(structures))
	positionIds := make([]uint64, 0, len(structures))
	for _, p := range newStruct {
		positionIds = append(positionIds, p.PositionId)
		memberIds = append(memberIds, p.MemberId)
	}

	positions, err := d.PositionRepository.QueryUndeletedInId(ctx, positionIds)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query position in id"))
		return
	}

	// Build the org strcuture
	// make map of object org strcuture
	// by use position id as key and org structure as value
	// Used array of position object for in org structure
	outPostMap := make(map[uint64]StructurePositionOut)
	for _, p := range positions {
		outPostMap[p.Id] = StructurePositionOut{
			Id:      p.Id,
			Name:    p.Name,
			Level:   p.Level,
			Members: make([]StructureMemberOut, 0, len(membCIds)),
		}
	}

	members, err := d.MemberRepository.QueryInId(ctx, memberIds)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query member in id"))
		return
	}

	// Make object of object
	// use member id as key and needed member data as value
	// used for append to members in position object
	outMemMap := make(map[string]StructureMemberOut)
	for _, m := range members {
		id := m.Id.UUID.String()
		outMemMap[id] = StructureMemberOut{
			Id:            id,
			Name:          m.Name,
			ProfilePicUrl: m.ProfilePicUrl,
		}
	}

	// Append the member objects to each position
	// based on what that member position
	for _, s := range newStruct {
		p := outPostMap[s.PositionId]
		p.Members = append(p.Members, outMemMap[s.MemberId])
		outPostMap[s.PositionId] = p
	}

	// Assemble all the datas
	outPos := make([]StructurePositionOut, 0, len(outPostMap))
	for _, v := range outPostMap {
		outPos = append(outPos, v)
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

	res := StructureRes{
		Id:        period.Id,
		StartDate: period.StartDate.Format("2006-01-02"),
		EndDate:   period.EndDate.Format("2006-01-02"),
		Positions: outPos,
		Vision:    string(vV),
		Mission:   string(mV),
	}

	out.Res = res

	return
}
