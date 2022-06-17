package dashboard

import (
	"context"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
)

type (
	PrivateRes struct {
		Documents     []DocumentOut        `json:"documents"`
		Members       []MemberOut          `json:"members"`
		Cashflows     CashflowRes          `json:"cashflow"`
		Dues          DuesOut              `json:"dues"`
		Blogs         []BlogOut            `json:"blogs"`
		Positions     []PositionOut        `json:"positions"`
		LatestHistory LatestHistoryRes     `json:"latest_history"`
		MemberDues    []MembersDuesOut     `json:"member_dues"`
		OrgPeriodGoal FindOrgPeriodGoalRes `json:"org_period_goal"`
		ActivePeriod  PeriodRes            `json:"active_period"`
	}
	PrivateOut struct {
		resp.Response
		Res PrivateRes
	}
)

func (d *DashboardDeps) GetPrivate(ctx context.Context) (out PrivateOut) {
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	c := make(chan CashflowRes)
	cr := make(chan resp.Response)
	go func(ctx context.Context, c chan CashflowRes, res chan resp.Response) {
		out := d.QueryCashflow(ctx, "")

		l := len(out.Res.Cashflows)
		if l > 5 {
			l = 5
		}

		cs := make([]CashflowOut, l)
		for i := range cs {
			cs[i] = CashflowOut(out.Res.Cashflows[i])
		}

		cc := CashflowRes{
			TotalCash:   out.Res.TotalCash,
			IncomeCash:  out.Res.IncomeCash,
			OutcomeCash: out.Res.OutcomeCash,
			Cashflows:   cs,
		}

		c <- cc
		res <- out.Response
	}(ctx, c, cr)

	m := make(chan []MemberOut)
	mr := make(chan resp.Response)
	go func(ctx context.Context, m chan []MemberOut, res chan resp.Response) {
		out := d.QueryMember(ctx, "")

		l := len(out.Res.Members)
		if l > 5 {
			l = 5
		}

		ms := make([]MemberOut, l)
		for i := range ms {
			ms[i] = MemberOut(out.Res.Members[i])
		}

		m <- ms
		res <- out.Response
	}(ctx, m, mr)

	do := make(chan []DocumentOut)
	dr := make(chan resp.Response)
	go func(ctx context.Context, do chan []DocumentOut, res chan resp.Response) {
		out := d.QueryDocument(ctx, "", "999")

		dc := make([]DocumentOut, 0, 0)
		for _, v := range out.Res.Documents {
			if v.Type == Filetype.String {
				dc = append(dc, DocumentOut(v))
			}

			if len(dc) == 5 {
				break
			}
		}

		do <- dc
		res <- out.Response
	}(ctx, do, dr)

	ds := make(chan DuesOut)
	dsr := make(chan resp.Response)
	go func(ctx context.Context, ds chan DuesOut, res chan resp.Response) {
		out := d.QueryDues(ctx, "")

		var d DuesOut
		if len(out.Res.Dues) > 1 {
			d = DuesOut(out.Res.Dues[0])
		}

		ds <- d
		res <- out.Response
	}(ctx, ds, dsr)

	b := make(chan []BlogOut)
	br := make(chan resp.Response)
	go func(ctx context.Context, b chan []BlogOut, res chan resp.Response) {
		out := d.QueryBlog(ctx, "")

		l := len(out.Res.Blogs)
		if l > 5 {
			l = 5
		}

		bs := make([]BlogOut, l)
		for i := range bs {
			bs[i] = BlogOut(out.Res.Blogs[i])
		}

		b <- bs
		res <- out.Response
	}(ctx, b, br)

	p := make(chan []PositionOut)
	pr := make(chan resp.Response)
	go func(ctx context.Context, p chan []PositionOut, res chan resp.Response) {
		out := d.QueryPosition(ctx, "")

		l := len(out.Res.Positions)
		if l > 5 {
			l = 5
		}

		ps := make([]PositionOut, l)
		for i := range ps {
			ps[i] = PositionOut(out.Res.Positions[i])
		}

		p <- ps
		res <- out.Response
	}(ctx, p, pr)

	h := make(chan LatestHistoryRes)
	hr := make(chan resp.Response)
	go func(ctx context.Context, h chan LatestHistoryRes, res chan resp.Response) {
		out := d.FindLatestHistory(ctx)

		h <- LatestHistoryRes(out.Res)
		res <- out.Response
	}(ctx, h, hr)

	md := make(chan []MembersDuesOut)
	mdr := make(chan resp.Response)
	go func(ctx context.Context, md chan []MembersDuesOut, res chan resp.Response) {
		out := d.QueryMembersDues(ctx, "0", "")

		l := len(out.Res.MemberDues)
		if l > 5 {
			l = 5
		}

		mds := make([]MembersDuesOut, l)
		for i := range mds {
			mds[i] = MembersDuesOut(out.Res.MemberDues[i])
		}

		md <- mds
		res <- out.Response
	}(ctx, md, mdr)

	op := make(chan FindOrgPeriodGoalRes)
	opr := make(chan resp.Response)
	go func(ctx context.Context, op chan FindOrgPeriodGoalRes, res chan resp.Response) {
		out := d.FindOrgPeriodGoal(ctx, "0")

		var o FindOrgPeriodGoalRes
		if out.Error == nil {
			o = FindOrgPeriodGoalRes(out.Res)
		}

		op <- o
		res <- out.Response
	}(ctx, op, opr)

	pe := make(chan PeriodRes)
	per := make(chan resp.Response)
	go func(ctx context.Context, pe chan PeriodRes, res chan resp.Response) {
		out := d.FindActivePeriod(ctx)

		var p PeriodRes
		if out.Error == nil {
			p = PeriodRes(out.Res)
		}

		pe <- p
		res <- out.Response
	}(ctx, pe, per)

	cV := <-c
	// crV := <-cr
	_ = <-cr
	mV := <-m
	// mrV := <-mr
	_ = <-mr
	doV := <-do
	// drV := <-dr
	_ = <-dr
	dsV := <-ds
	// dsrV := <-dsr
	_ = <-dsr
	bV := <-b
	// brV := <-br
	_ = <-br
	pV := <-p
	// prV := <-pr
	_ = <-pr
	hV := <-h
	// hrV := <-hr
	_ = <-hr
	mdV := <-md
	// mdrV := <-mdr
	_ = <-mdr
	opV := <-op
	// oprV := <-opr
	_ = <-opr
	peV := <-pe
	// perV := <-per
	_ = <-per

	// if crV.Error != nil {
	// 	out.Response = crV
	// 	return
	// }

	// if mrV.Error != nil {
	// 	out.Response = mrV
	// 	return
	// }

	// if drV.Error != nil {
	// 	out.Response = drV
	// 	return
	// }

	// if dsrV.Error != nil {
	// 	out.Response = dsrV
	// 	return
	// }

	// if brV.Error != nil {
	// 	out.Response = brV
	// 	return
	// }

	// if prV.Error != nil {
	// 	out.Response = prV
	// 	return
	// }

	// if hrV.Error != nil {
	// 	out.Response = hrV
	// 	return
	// }

	// if mdrV.Error != nil {
	// 	out.Response = mdrV
	// 	return
	// }

	// if oprV.Error != nil {
	// 	out.Response = oprV
	// 	return
	// }

	// if perV.Error != nil {
	// 	out.Response = perV
	// 	return
	// }

	out.Res = PrivateRes{
		Documents:     doV,
		Members:       mV,
		Cashflows:     cV,
		Dues:          dsV,
		Blogs:         bV,
		Positions:     pV,
		LatestHistory: hV,
		MemberDues:    mdV,
		OrgPeriodGoal: opV,
		ActivePeriod:  peV,
	}

	return
}

type (
	PublicRes struct {
		Documents     []DocumentOut    `json:"documents"`
		Blogs         []BlogOut        `json:"blogs"`
		LatestHistory LatestHistoryRes `json:"latest_history"`
		Structure     StructureRes     `json:"org_period_structures"`
	}
	PublicOut struct {
		resp.Response
		Res PublicRes
	}
)

func (d *DashboardDeps) GetPublic(ctx context.Context) (out PublicOut) {
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	do := make(chan []DocumentOut)
	dr := make(chan resp.Response)
	go func(ctx context.Context, do chan []DocumentOut, res chan resp.Response) {
		out := d.QueryDocument(ctx, "", "999")

		dc := make([]DocumentOut, 0, 0)
		for _, v := range out.Res.Documents {
			if v.IsPrivate == false && v.Type == Filetype.String {
				dc = append(dc, DocumentOut(v))
			}

			if len(dc) == 5 {
				break
			}
		}

		do <- dc
		res <- out.Response
	}(ctx, do, dr)

	b := make(chan []BlogOut)
	br := make(chan resp.Response)
	go func(ctx context.Context, b chan []BlogOut, res chan resp.Response) {
		out := d.QueryBlog(ctx, "")

		l := len(out.Res.Blogs)
		if l > 8 {
			l = 8
		}

		bs := make([]BlogOut, l)
		for i := range bs {
			bs[i] = BlogOut(out.Res.Blogs[i])
		}

		b <- bs
		res <- out.Response
	}(ctx, b, br)

	h := make(chan LatestHistoryRes)
	hr := make(chan resp.Response)
	go func(ctx context.Context, h chan LatestHistoryRes, res chan resp.Response) {
		out := d.FindLatestHistory(ctx)

		h <- LatestHistoryRes(out.Res)
		res <- out.Response
	}(ctx, h, hr)

	op := make(chan StructureRes)
	ops := make(chan resp.Response)
	go func(ctx context.Context, op chan StructureRes, res chan resp.Response) {
		out := d.QueryPeriodStructure(ctx, "0")

		o := StructureRes{
			Positions: make([]StructurePositionOut, 0, 0),
		}

		if out.Error == nil {
			ps := make([]StructurePositionOut, len(out.Res.Positions))
			for i, v := range out.Res.Positions {
				m := make([]StructureMemberOut, len(v.Members))
				for j, q := range v.Members {
					m[j] = StructureMemberOut(q)
				}

				ps[i] = StructurePositionOut{
					Id:      v.Id,
					Name:    v.Name,
					Level:   v.Level,
					Members: m,
				}
			}

			o = StructureRes{
				Id:        out.Res.Id,
				StartDate: out.Res.StartDate,
				EndDate:   out.Res.EndDate,
				Positions: ps,
				Vision:    out.Res.Vision,
				Mission:   out.Res.Mission,
			}
		}

		op <- o
		res <- out.Response
	}(ctx, op, ops)

	doV := <-do
	// drV := <-dr
	_ = <-dr
	bV := <-b
	// brV := <-br
	_ = <-br
	hV := <-h
	// hrV := <-hr
	_ = <-hr
	opV := <-op
	// opsV := <-ops
	_ = <-ops

	// if drV.Error != nil {
	// 	out.Response = drV
	// 	return
	// }

	// if brV.Error != nil {
	// 	out.Response = brV
	// 	return
	// }

	// if hrV.Error != nil {
	// 	out.Response = hrV
	// 	return
	// }

	// if opsV.Error != nil {
	// 	out.Response = opsV
	// 	return
	// }

	out.Res = PublicRes{
		Documents:     doV,
		Blogs:         bV,
		LatestHistory: hV,
		Structure:     opV,
	}

	return
}