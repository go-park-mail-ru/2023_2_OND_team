package v1

import (
	"net/http"
	"strconv"

	errHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/errors"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/search"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
)

var (
	defaultSearchCount  = 20
	maxSearchCount      = 50
	defaultSearchOffset = 0
	maxSearchOffset     = 50
)

var (
	userSortOpts  = []string{"id", "subscribers"}
	boardSortOpts = []string{"id", "pins"}
	pinSortOpts   = []string{"id", "likes"}
)

var (
	defaultSearchSort = "id"
	defaultSortOder   = "desc"
)

func (h *HandlerHTTP) SearchUsers(w http.ResponseWriter, r *http.Request) {
	opts, err := GetSearchOpts(r, userSortOpts, defaultSearchSort)
	if err != nil {
		h.responseErr(w, r, err)
		return
	}

	if users, err := h.searchCase.GetUsers(r.Context(), opts); err != nil {
		h.responseErr(w, r, err)
	} else if err := responseOk(http.StatusOK, w, "got users sucessfully", h.converter.ToUsersForSearchFromService(users)); err != nil {
		h.responseErr(w, r, err)
	}
}

func (h *HandlerHTTP) SearchBoards(w http.ResponseWriter, r *http.Request) {
	opts, err := GetSearchOpts(r, boardSortOpts, defaultSearchSort)
	if err != nil {
		h.responseErr(w, r, err)
		return
	}

	if boards, err := h.searchCase.GetBoards(r.Context(), opts); err != nil {
		h.responseErr(w, r, err)
	} else if err := responseOk(http.StatusOK, w, "got boards sucessfully", h.converter.ToBoardsForSearchFromService(boards)); err != nil {
		h.responseErr(w, r, err)
	}
}

func (h *HandlerHTTP) SearchPins(w http.ResponseWriter, r *http.Request) {
	opts, err := GetSearchOpts(r, pinSortOpts, defaultSearchSort)
	if err != nil {
		h.responseErr(w, r, err)
		return
	}

	if pins, err := h.searchCase.GetPins(r.Context(), opts); err != nil {
		h.responseErr(w, r, err)
	} else if err := responseOk(http.StatusOK, w, "got pins sucessfully", h.converter.ToPinsForSearchFromService(pins)); err != nil {
		h.responseErr(w, r, err)
	}
}

func GetSearchOpts(r *http.Request, sortOpts []string, defaultSortOpt string) (*search.SearchOpts, error) {
	opts := &search.SearchOpts{}
	invalidParams := map[string]string{}

	generalOpts, err := GetGeneralOpts(r, invalidParams)
	if err != nil {
		return nil, err
	}
	opts.General = *generalOpts

	if sortBy := r.URL.Query().Get("sortBy"); sortBy != "" {
		if !isCorrentSortOpt(sortOpts, sortBy) {
			invalidParams["sortBy"] = sortBy
		} else {
			opts.SortBy = sortBy
		}
	} else {
		opts.SortBy = defaultSortOpt
	}

	if len(invalidParams) > 0 {
		return nil, &errHTTP.ErrInvalidQueryParam{Params: invalidParams}
	}

	return opts, nil
}

func isCorrentSortOpt(correctOpts []string, opt string) bool {
	for _, correctOpt := range correctOpts {
		if opt == correctOpt {
			return true
		}
	}
	return false
}

func GetGeneralOpts(r *http.Request, invalidParams map[string]string) (*search.GeneralOpts, error) {
	opts := &search.GeneralOpts{}

	if templateParam := r.URL.Query().Get("template"); templateParam != "" {
		if template := search.Template(templateParam); !template.Validate() {
			invalidParams["template"] = string(template)
		} else {
			opts.Template = template
		}
	} else {
		return nil, &errHTTP.ErrNoData{}
	}

	if sortOrder := r.URL.Query().Get("order"); sortOrder != "" {
		if sortOrder != "asc" && sortOrder != "desc" {
			invalidParams["order"] = sortOrder
		} else {
			opts.SortOrder = sortOrder
		}
	} else {
		opts.SortOrder = defaultSortOder
	}

	if countParam := r.URL.Query().Get("count"); countParam != "" {
		if count, err := strconv.ParseInt(countParam, 10, 64); err != nil || count < 0 {
			invalidParams["count"] = countParam
		} else {
			opts.Count = int(count)
		}
	} else {
		opts.Count = defaultSearchCount
	}

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if offset, err := strconv.ParseInt(offsetParam, 10, 64); err != nil || offset < 0 {
			invalidParams["offset"] = offsetParam
		} else {
			opts.Offset = int(offset)
		}
	} else {
		opts.Offset = defaultSearchOffset
	}

	if opts.Count > maxSearchCount {
		opts.Count = maxSearchCount
	}
	if opts.Offset > maxSearchOffset {
		opts.Offset = maxSearchOffset
	}

	userID, _ := r.Context().Value(auth.KeyCurrentUserID).(int)
	opts.CurrUserID = userID

	return opts, nil
}
