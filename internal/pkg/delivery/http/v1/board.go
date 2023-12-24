package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	"github.com/mailru/easyjson"

	errHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/errors"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/structs"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var TimeFormat = "2006-01-02"

func (h *HandlerHTTP) CreateNewBoard(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	if contentType := r.Header.Get("Content-Type"); contentType != ApplicationJson {
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadContentType)
		responseError(w, code, message)
		return
	}

	var newBoard structs.BoardData
	err := easyjson.UnmarshalFromReader(r.Body, &newBoard)
	defer r.Body.Close()
	if err != nil {
		logger.Info("create board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadBody)
		responseError(w, code, message)
		return
	}

	err = newBoard.Validate()
	if err != nil {
		logger.Info("create board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	tagTitles := make([]string, 0)
	if newBoard.Tags != nil {
		tagTitles = append(tagTitles, newBoard.Tags...)
	}
	authorID := r.Context().Value(auth.KeyCurrentUserID).(int)

	newBoardID, err := h.boardCase.CreateNewBoard(r.Context(), entity.Board{
		Title:       *newBoard.Title,
		Description: *newBoard.Description,
		Public:      *newBoard.Public,
		AuthorID:    authorID,
	}, tagTitles)

	if err != nil {
		logger.Info("create board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusCreated, w, "new board was created successfully", map[string]int{"new_board_id": newBoardID})
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errHTTP.ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) GetUserBoards(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	username := chi.URLParam(r, "username")
	if !isValidUsername(username) {
		logger.Info("update board", log.F{"message", errHTTP.ErrInvalidUsername.Error()})
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrInvalidUsername)
		responseError(w, code, message)
		return
	}

	boards, err := h.boardCase.GetBoardsByUsername(r.Context(), username)
	if err != nil {
		logger.Info("get user boards", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	userBoards := make([]structs.CertainBoard, 0, len(boards))
	for _, board := range boards {
		userBoards = append(userBoards, h.converter.ToCertainBoardFromService(&board))
	}
	err = responseOk(http.StatusOK, w, "got user boards successfully", userBoards)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errHTTP.ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) GetCertainBoard(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		logger.Info("get certain board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadUrlParam)
		responseError(w, code, message)
		return
	}

	board, username, err := h.boardCase.GetCertainBoard(r.Context(), int(boardID))
	if err != nil {
		logger.Info("get certain board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusOK, w, "got certain board successfully", h.converter.ToCertainBoardUsernameFromService(&board, username))
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errHTTP.ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) GetBoardInfoForUpdate(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		logger.Info("get certain board info for update", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadUrlParam)
		responseError(w, code, message)
		return
	}

	board, tagTitles, err := h.boardCase.GetBoardInfoForUpdate(r.Context(), int(boardID))
	if err != nil {
		logger.Info("get certain board info for update", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusOK, w, "got certain board successfully", map[string]interface{}{"board": h.converter.ToBoardFromService(&board), "tags": tagTitles})
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errHTTP.ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) UpdateBoardInfo(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	if contentType := r.Header.Get("Content-Type"); contentType != ApplicationJson {
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadContentType)
		responseError(w, code, message)
		return
	}

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadUrlParam)
		responseError(w, code, message)
		return
	}

	var updatedData structs.BoardData
	err = easyjson.UnmarshalFromReader(r.Body, &updatedData)
	defer r.Body.Close()
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadBody)
		responseError(w, code, message)
		return
	}

	err = updatedData.Validate()
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	tagTitles := make([]string, 0)
	if updatedData.Tags != nil {
		tagTitles = append(tagTitles, updatedData.Tags...)
	}

	updatedBoard := entity.Board{
		ID:          int(boardID),
		Title:       *updatedData.Title,
		Description: *updatedData.Description,
		Public:      *updatedData.Public,
	}
	err = h.boardCase.UpdateBoardInfo(r.Context(), updatedBoard, tagTitles)
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusOK, w, "updated certain board successfully", nil)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errHTTP.ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadUrlParam)
		responseError(w, code, message)
		return
	}

	err = h.boardCase.DeleteCertainBoard(r.Context(), int(boardID))
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusOK, w, "deleted board successfully", nil)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errHTTP.ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) AddPinsToBoard(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	boardIDStr := chi.URLParam(r, "boardID")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		logger.Error("parse board id from query params")
		err = responseError(w, "parse_url", "internal error")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	pins := make(map[string][]int)
	err = json.NewDecoder(r.Body).Decode(&pins)
	defer r.Body.Close()
	if err != nil {
		logger.Info("bad decode body")
		err = responseError(w, "bad_body", "failed to parse the request body")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}
	pinIds, ok := pins["pins"]
	if !ok {
		logger.Info("the request does not specify pins")
		err = responseError(w, "bad_body", "the request does not specify pins")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err = h.pinCase.IsAvailableBatchPinForFixOnBoard(r.Context(), pinIds, userID)
	if err != nil {
		logger.Warn(err.Error(), log.F{"action", "check availability pins for fixed on board"})
		err = responseError(w, "not_access", "there are pins in the batch that are not available for the user to add")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err = h.boardCase.FixPinsOnBoard(r.Context(), int(boardID), pinIds, userID)
	if err != nil {
		logger.Warn(err.Error(), log.F{"action", "fix pins on board"})
		err = responseError(w, "not_access", "there are pins in the batch that are not available for the user to add")
	} else {
		err = responseOk(http.StatusCreated, w, "pins have been successfully added to the board", nil)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) DeletePinFromBoard(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		logger.Info("delete pin from board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadUrlParam)
		responseError(w, code, message)
		return
	}

	if contentType := r.Header.Get("Content-Type"); contentType != ApplicationJson {
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadContentType)
		responseError(w, code, message)
		return
	}

	delPinFromBoard := structs.DeletePinFromBoard{}
	err = easyjson.UnmarshalFromReader(r.Body, &delPinFromBoard)
	defer r.Body.Close()
	if err != nil {
		code, message := errHTTP.GetErrCodeMessage(errHTTP.ErrBadBody)
		responseError(w, code, message)
		return
	}

	err = h.boardCase.DeletePinFromBoard(r.Context(), int(boardID), delPinFromBoard.PinID)
	if err != nil {
		logger.Info("delete pin from board", log.F{"message", err.Error()})
		code, message := errHTTP.GetErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusOK, w, "deleted pin from board successfully", nil)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errHTTP.ErrInternalError.Error()))
	}
}
