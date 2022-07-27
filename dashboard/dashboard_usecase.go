package dashboard

import (
	"context"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dues"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
)

type (
	PrivateRes struct {
		MemberTotal     int64                `json:"member_total"`
		DocumentTotal   int64                `json:"document_total"`
		BlogTotal       int64                `json:"blog_total"`
		PositionTotal   int64                `json:"position_total"`
		MemberDuesTotal int64                `json:"member_dues_total"`
		ImageTotal      int64                `json:"image_total"`
		Documents       []DocumentOut        `json:"documents"`
		Members         []MemberOut          `json:"members"`
		Cashflows       CashflowRes          `json:"cashflow"`
		Dues            DuesOut              `json:"dues"`
		Blogs           []BlogOut            `json:"blogs"`
		Positions       []PositionOut        `json:"positions"`
		LatestHistory   LatestHistoryRes     `json:"latest_history"`
		MemberDues      []MembersDuesOut     `json:"member_dues"`
		OrgPeriodGoal   FindOrgPeriodGoalRes `json:"org_period_goal"`
		ActivePeriod    PeriodRes            `json:"active_period"`
		Images          []ImageOut           `json:"images"`
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
		cashflowStatsOut := d.CalculateCashflow(ctx)

		cc := CashflowRes{
			TotalCash:   cashflowStatsOut.Res.TotalCash,
			IncomeCash:  cashflowStatsOut.Res.IncomeCash,
			OutcomeCash: cashflowStatsOut.Res.OutcomeCash,
		}

		c <- cc
		res <- out.Response
	}(ctx, c, cr)

	m := make(chan []MemberOut)
	mt := make(chan int64)
	mr := make(chan resp.Response)
	go func(ctx context.Context, m chan []MemberOut, mt chan int64, res chan resp.Response) {
		out := d.QueryMember(ctx, "", "", "5")

		l := len(out.Res.Members)
		if l > 5 {
			l = 5
		}

		ms := make([]MemberOut, l)
		for i := range ms {
			ms[i] = MemberOut(out.Res.Members[i])
		}

		m <- ms
		mt <- out.Res.Total
		res <- out.Response
	}(ctx, m, mt, mr)

	do := make(chan []DocumentOut)
	dt := make(chan int64)
	dr := make(chan resp.Response)
	go func(ctx context.Context, do chan []DocumentOut, dt chan int64, res chan resp.Response) {
		out := d.QueryDocument(ctx, "", "", "999")

		dc := make([]DocumentOut, 0)
		for _, v := range out.Res.Documents {
			if v.Type == Filetype.String {
				dc = append(dc, DocumentOut(v))
			}

			if len(dc) == 5 {
				break
			}
		}

		do <- dc
		dt <- out.Res.Total
		res <- out.Response
	}(ctx, do, dt, dr)

	ds := make(chan DuesOut)
	dsr := make(chan resp.Response)
	go func(ctx context.Context, ds chan DuesOut, res chan resp.Response) {
		out := d.FindLatestDues(ctx)

		ds <- DuesOut(out.Res)
		res <- out.Response
	}(ctx, ds, dsr)

	b := make(chan []BlogOut)
	bt := make(chan int64)
	br := make(chan resp.Response)
	go func(ctx context.Context, b chan []BlogOut, bt chan int64, res chan resp.Response) {
		out := d.QueryBlog(ctx, "", "")

		l := len(out.Res.Blogs)
		if l > 5 {
			l = 5
		}

		bs := make([]BlogOut, l)
		for i := range bs {
			bs[i] = BlogOut(out.Res.Blogs[i])
		}

		b <- bs
		bt <- out.Res.Total
		res <- out.Response
	}(ctx, b, bt, br)

	p := make(chan []PositionOut)
	pt := make(chan int64)
	pr := make(chan resp.Response)
	go func(ctx context.Context, p chan []PositionOut, pt chan int64, res chan resp.Response) {
		out := d.QueryPosition(ctx, "", "5")

		l := len(out.Res.Positions)
		if l > 5 {
			l = 5
		}

		ps := make([]PositionOut, l)
		for i := range ps {
			ps[i] = PositionOut(out.Res.Positions[i])
		}

		p <- ps
		pt <- out.Res.Total
		res <- out.Response
	}(ctx, p, pt, pr)

	h := make(chan LatestHistoryRes)
	hr := make(chan resp.Response)
	go func(ctx context.Context, h chan LatestHistoryRes, res chan resp.Response) {
		out := d.FindLatestHistory(ctx)

		h <- LatestHistoryRes(out.Res)
		res <- out.Response
	}(ctx, h, hr)

	md := make(chan []MembersDuesOut)
	mdT := make(chan int64)
	mdr := make(chan resp.Response)
	go func(ctx context.Context, md chan []MembersDuesOut, mdT chan int64, res chan resp.Response) {
		out := d.QueryMembersDues(ctx, "0", dues.QueryMembersDuesQIn{
			Limit: "3",
		})

		l := len(out.Res.MemberDues)
		if l > 3 {
			l = 3
		}

		mds := make([]MembersDuesOut, l)
		for i := range mds {
			mds[i] = MembersDuesOut(out.Res.MemberDues[i])
		}

		md <- mds
		mdT <- out.Res.Total
		res <- out.Response
	}(ctx, md, mdT, mdr)

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

	imgs, imgT, _ := func(ctx context.Context) (imgs []ImageOut, imgT int64, res resp.Response) {
		out := d.QueryImage(ctx, "", "5")

		l := len(out.Res.Images)
		if l > 5 {
			l = 5
		}

		ps := make([]ImageOut, l)
		for i := range ps {
			ps[i] = ImageOut(out.Res.Images[i])
		}

		imgs = ps
		imgT = out.Res.Total
		res = out.Response
		return
	}(ctx)

	cV := <-c
	<-cr
	mV := <-m
	mtV := <-mt
	<-mr
	doV := <-do
	dtV := <-dt
	<-dr
	dsV := <-ds
	<-dsr
	bV := <-b
	btV := <-bt
	<-br
	pV := <-p
	ptV := <-pt
	<-pr
	hV := <-h
	<-hr
	mdV := <-md
	mdtV := <-mdT
	<-mdr
	opV := <-op
	<-opr
	peV := <-pe
	<-per

	out.Res = PrivateRes{
		MemberTotal:     mtV,
		DocumentTotal:   dtV,
		BlogTotal:       btV,
		PositionTotal:   ptV,
		MemberDuesTotal: mdtV,
		ImageTotal:      imgT,
		Documents:       doV,
		Members:         mV,
		Cashflows:       cV,
		Dues:            dsV,
		Blogs:           bV,
		Positions:       pV,
		LatestHistory:   hV,
		MemberDues:      mdV,
		OrgPeriodGoal:   opV,
		ActivePeriod:    peV,
		Images:          imgs,
	}

	return
}

type (
	PublicRes struct {
		Documents     []DocumentOut    `json:"documents"`
		Blogs         []BlogOut        `json:"blogs"`
		LatestHistory LatestHistoryRes `json:"latest_history"`
		Structure     StructureRes     `json:"org_period_structures"`
		Images        []ImageOut       `json:"images"`
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
		out := d.QueryDocument(ctx, "", "", "999")

		dc := make([]DocumentOut, 0)
		for _, v := range out.Res.Documents {
			if !v.IsPrivate && v.Type == Filetype.String {
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
		out := d.QueryBlog(ctx, "", "")

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
			Id:        out.Res.Id,
			StartDate: out.Res.StartDate,
			EndDate:   out.Res.EndDate,
			Positions: make([]StructurePositionOut, 0),
			Vision:    out.Res.Vision,
			Mission:   out.Res.Mission,
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

			o.Positions = ps
		}

		op <- o
		res <- out.Response
	}(ctx, op, ops)

	imgs := make(chan []ImageOut)
	imgR := make(chan resp.Response)
	go func(ctx context.Context, imgs chan []ImageOut, res chan resp.Response) {
		out := d.QueryImage(ctx, "", "5")

		l := len(out.Res.Images)
		if l > 5 {
			l = 5
		}

		ps := make([]ImageOut, l)
		for i := range ps {
			ps[i] = ImageOut(out.Res.Images[i])
		}

		imgs <- ps
		res <- out.Response
	}(ctx, imgs, imgR)

	doV := <-do
	<-dr
	bV := <-b
	<-br
	hV := <-h
	<-hr
	opV := <-op
	<-ops
	imgsV := <-imgs
	<-imgR

	out.Res = PublicRes{
		Documents:     doV,
		Blogs:         bV,
		LatestHistory: hV,
		Structure:     opV,
		Images:        imgsV,
	}

	return
}
