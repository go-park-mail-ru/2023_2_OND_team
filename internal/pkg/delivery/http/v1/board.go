package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	bCase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var TimeFormat = "2006-01-02"

var (
	ErrEmptyTitle        = errors.New("empty or null board title has been provided")
	ErrEmptyPubOpt       = errors.New("null public option has been provided")
	ErrInvalidBoardTitle = errors.New("invalid or empty board title has been provided")
	ErrInvalidTagTitles  = errors.New("invalid tag titles have been provided")
	ErrInvalidUsername   = errors.New("invalid username has been provided")
)

var (
	wrappedErrors      = map[error]string{ErrInvalidTagTitles: "bad_Tagtitles"}
	errCodeCompability = map[error]string{
		ErrInvalidBoardTitle:     "bad_boardTitle",
		ErrEmptyTitle:            "empty_boardTitle",
		ErrEmptyPubOpt:           "bad_pubOpt",
		ErrInvalidUsername:       "bad_username",
		bCase.ErrInvalidUsername: "non_existingUser",
		bCase.ErrNoSuchBoard:     "no_board",
		bCase.ErrNoAccess:        "no_access",
	}
)

// data for board creation/update
type BoardData struct {
	Title       *string  `json:"title" example:"new board"`
	Description *string  `json:"description" example:"long desc"`
	Public      *bool    `json:"public" example:"true"`
	Tags        []string `json:"tags" example:"['blue', 'car']"`
}

// board view for delivery layer
type CertainBoard struct {
	ID          int      `json:"board_id" example:"22"`
	Title       string   `json:"title" example:"new board"`
	Description string   `json:"description" example:"long desc"`
	CreatedAt   string   `json:"created_at" example:"07-11-2023"`
	PinsNumber  int      `json:"pins_number" example:"12"`
	Pins        []string `json:"pins" example:"['/pic1', '/pic2']"`
	Tags        []string `json:"tags" example:"['love', 'green']"`
}

func ToCertainBoardFromService(board entity.BoardWithContent) CertainBoard {
	return CertainBoard{
		ID:          board.BoardInfo.ID,
		Title:       board.BoardInfo.Title,
		Description: board.BoardInfo.Description,
		CreatedAt:   board.BoardInfo.CreatedAt.Format(TimeFormat),
		PinsNumber:  board.PinsNumber,
		Pins:        board.Pins,
		Tags:        board.TagTitles,
	}
}

func (data *BoardData) Validate() error {
	if data.Title == nil || *data.Title == "" {
		return ErrInvalidBoardTitle
	}
	if data.Description == nil {
		data.Description = new(string)
		*data.Description = ""
	}
	if data.Public == nil {
		return ErrEmptyPubOpt
	}
	if !isValidBoardTitle(*data.Title) {
		return ErrInvalidBoardTitle
	}
	if err := checkIsValidTagTitles(data.Tags); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrInvalidTagTitles)
	}
	return nil
}

func getErrCodeMessage(err error) (string, string) {
	var (
		code              string
		general, specific bool
	)

	code, general = generalErrCodeCompability[err]
	if general {
		return code, err.Error()
	}

	code, specific = errCodeCompability[err]
	if !specific {
		for wrappedErr, code_ := range wrappedErrors {
			if errors.Is(err, wrappedErr) {
				specific = true
				code = code_
			}
		}
	}
	if specific {
		return code, err.Error()
	}

	return ErrInternalError.Error(), generalErrCodeCompability[ErrInternalError]
}

func (h *HandlerHTTP) CreateNewBoard(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	if contentType := w.Header().Get("Content-Type"); contentType != ApplicationJson {
		code, message := getErrCodeMessage(ErrBadContentType)
		responseError(w, code, message)
		return
	}

	var newBoard BoardData
	err := json.NewDecoder(r.Body).Decode(&newBoard)
	defer r.Body.Close()
	if err != nil {
		logger.Info("create board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(ErrBadBody)
		responseError(w, code, message)
		return
	}

	err = newBoard.Validate()
	if err != nil {
		logger.Info("create board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(err)
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
		code, message := getErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusCreated, w, "new board was created successfully", map[string]int{"new_board_id": newBoardID})
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) GetUserBoards(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	username := chi.URLParam(r, "username")
	if !isValidUsername(username) {
		logger.Info("update board", log.F{"message", ErrInvalidUsername.Error()})
		code, message := getErrCodeMessage(ErrInvalidUsername)
		responseError(w, code, message)
		return
	}

	boards, err := h.boardCase.GetBoardsByUsername(r.Context(), username)
	if err != nil {
		logger.Info("get user boards", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	userBoards := make([]CertainBoard, 0, len(boards))
	for _, board := range boards {
		userBoards = append(userBoards, ToCertainBoardFromService(board))
	}
	err = responseOk(http.StatusOK, w, "got user boards successfully", userBoards)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) GetCertainBoard(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		logger.Info("get certain board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(ErrBadUrlParam)
		responseError(w, code, message)
		return
	}

	board, err := h.boardCase.GetCertainBoard(r.Context(), int(boardID))
	if err != nil {
		logger.Info("get certain board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusOK, w, "got certain board successfully", ToCertainBoardFromService(board))
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) GetBoardInfoForUpdate(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		logger.Info("get certain board info for update", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(ErrBadUrlParam)
		responseError(w, code, message)
		return
	}

	board, tagTitles, err := h.boardCase.GetBoardInfoForUpdate(r.Context(), int(boardID))
	if err != nil {
		logger.Info("get certain board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusOK, w, "got certain board successfully", map[string]interface{}{"board": board, "tags": tagTitles})
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) UpdateBoardInfo(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	if contentType := w.Header().Get("Content-Type"); contentType != ApplicationJson {
		code, message := getErrCodeMessage(ErrBadContentType)
		responseError(w, code, message)
		return
	}

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(ErrBadUrlParam)
		responseError(w, code, message)
		return
	}

	var updatedData BoardData
	err = json.NewDecoder(r.Body).Decode(&updatedData)
	defer r.Body.Close()
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(ErrBadBody)
		responseError(w, code, message)
		return
	}

	err = updatedData.Validate()
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(err)
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
		code, message := getErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusOK, w, "updated certain board successfully", nil)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrInternalError.Error()))
	}
}

func (h *HandlerHTTP) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(ErrBadUrlParam)
		responseError(w, code, message)
		return
	}

	err = h.boardCase.DeleteCertainBoard(r.Context(), int(boardID))
	if err != nil {
		logger.Info("update certain board", log.F{"message", err.Error()})
		code, message := getErrCodeMessage(err)
		responseError(w, code, message)
		return
	}

	err = responseOk(http.StatusOK, w, "deleted board successfully", nil)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrInternalError.Error()))
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
