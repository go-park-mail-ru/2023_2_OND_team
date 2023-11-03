package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	boardDTO "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	bCase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func (h *HandlerHTTP) CreateNewBoard(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on create new board:", logger.F{"method", r.Method}, logger.F{"URL", r.URL.Path})
	SetContentTypeJSON(w)

	var newBoard boardDTO.BoardData
	err := json.NewDecoder(r.Body).Decode(&newBoard)
	defer r.Body.Close()
	if err != nil {
		h.log.Info("create board: ", logger.F{"message", err.Error()})
		responseError(w, BadBodyCode, BadBodyMessage)
		return
	}
	fmt.Println(newBoard)
	newBoard.AuthorID = r.Context().Value(auth.KeyCurrentUserID).(int)
	err = h.boardCase.CreateNewBoard(r.Context(), newBoard)
	if err != nil {
		h.log.Info("create board", logger.F{"message", err.Error()})
		switch err {
		case bCase.ErrInvalidBoardTitle:
			responseError(w, "bad_boardTitle", err.Error())
		default:
			if errors.Is(err, bCase.ErrInvalidTagTitles) {
				responseError(w, "bad_tagTitles", err.Error())
				return
			}
			responseError(w, InternalErrorCode, InternalServerErrMessage)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = responseOk(w, "new board was created successfully", nil)
	if err != nil {
		h.log.Error(err.Error())
		w.Write([]byte(InternalServerErrMessage))
		return
	}
	h.log.Info("create new board request: successfull respond")
}

func (h *HandlerHTTP) GetUserBoards(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request to get user boards:", logger.F{"method", r.Method}, logger.F{"URL", r.URL.Path})
	SetContentTypeJSON(w)

	userBoards, err := h.boardCase.GetBoardsByUsername(r.Context(), chi.URLParam(r, "username"))
	if err != nil {
		h.log.Info("get user boards: ", logger.F{"message", err.Error()})
		switch err {
		case bCase.ErrInvalidUsername:
			responseError(w, "bad_username", err.Error())
		default:
			responseError(w, InternalErrorCode, InternalServerErrMessage)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = responseOk(w, "got user boards successfully", userBoards)
	if err != nil {
		h.log.Error(err.Error())
		w.Write([]byte(InternalServerErrMessage))
		return
	}
	h.log.Info("get user boards request: successfull respond")
}

func (h *HandlerHTTP) GetCertainBoard(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request to get certain board:", logger.F{"method", r.Method}, logger.F{"URL", r.URL.Path})
	SetContentTypeJSON(w)

	boardID, err := strconv.ParseInt(chi.URLParam(r, "boardID"), 10, 64)
	if err != nil {
		h.log.Info("get certain board ", logger.F{"message", err.Error()})
		responseError(w, BadQueryParamCode, BadQueryParamMessage)
		return
	}

	board, err := h.boardCase.GetCertainBoardByID(r.Context(), int(boardID))
	if err != nil {
		h.log.Info("get certain board: ", logger.F{"message", err.Error()})
		switch err {
		case bCase.ErrNoSuchBoard:
			responseError(w, "no_board", err.Error())
		default:
			responseError(w, InternalErrorCode, InternalServerErrMessage)
		}
		return
	}

	err = responseOk(w, "got certain board successfully", board)
	if err != nil {
		h.log.Error(err.Error())
		w.Write([]byte(InternalServerErrMessage))
		return
	}

	h.log.Info("get certain board request: successfull respond")
}
