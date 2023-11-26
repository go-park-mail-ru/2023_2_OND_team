package search

import errPkg "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/errors"

type ErrNoUsers struct {
}

func (e *ErrNoUsers) Error() string {
	return "Can't find any user"
}

func (e *ErrNoUsers) Type() errPkg.Type {
	return errPkg.ErrNotFound
}

type ErrNoBoards struct {
}

func (e *ErrNoBoards) Error() string {
	return "Can't find any board"
}

func (e *ErrNoBoards) Type() errPkg.Type {
	return errPkg.ErrNotFound
}

type ErrNoPins struct {
}

func (e *ErrNoPins) Error() string {
	return "Can't find any pin"
}

func (e *ErrNoPins) Type() errPkg.Type {
	return errPkg.ErrNotFound
}
