package history

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/pkg/errors"
)

type (
	AddHistoryIn struct {
		Content string `json:"content"`
	}
	AddHistoryRes struct {
		Id string `json:"id"`
	}
	AddHistoryOut struct {
		resp.Response
		Res AddHistoryRes
	}
)

func (d *HistoryDeps) AddHistory(ctx context.Context, in AddHistoryIn) (out AddHistoryOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddHistoryIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "add history validaion"))
		return
	}

	var nc map[string]interface{}
	if in.Content != "" {
		if err = json.Unmarshal([]byte(in.Content), &nc); err != nil {
			out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "unmarshal json"))
			return
		}
	}

	history := HistoryModel{
		Content: nc,
	}

	historyId, err := d.HistoryRepository.Save(ctx, history)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save history"))
		return
	}

	out.StatusCode = http.StatusCreated
	out.Res.Id = historyId

	return
}

type (
	LatestHistoryRes struct {
		Id      string `json:"id"`
		Content string `json:"content"`
	}
	FindLatestHistoryOut struct {
		resp.Response
		Res LatestHistoryRes
	}
)

func (d *HistoryDeps) FindLatestHistory(ctx context.Context) (out FindLatestHistoryOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	history, _ := d.HistoryRepository.FindLatest(ctx)

	c := []byte("")
	if history.Content != nil && len(history.Content) != 0 {
		c, err = json.Marshal(history.Content)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "json marshal"))
			return
		}
	}

	out.StatusCode = http.StatusOK
	out.Res = LatestHistoryRes{
		Id:      history.Id,
		Content: string(c),
	}

	return
}
