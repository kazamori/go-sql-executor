package query

import (
	"io/ioutil"

	"github.com/gohxs/readline"
	"github.com/xo/usql/rline"
)

func NewNopRline() (rline.IO, error) {
	stdout := ioutil.Discard
	stderr := ioutil.Discard

	l, err := readline.NewEx(&readline.Config{
		HistoryFile:            "",
		DisableAutoSaveHistory: true,
		InterruptPrompt:        "^C",
		HistorySearchFold:      true,
		Stdin:                  readline.Stdin,
		Stdout:                 stdout,
		Stderr:                 stderr,
		FuncIsTerminal: func() bool {
			return false
		},
		FuncFilterInputRune: func(r rune) (rune, bool) {
			if r == readline.CharCtrlZ {
				return r, false
			}
			return r, true
		},
	})
	if err != nil {
		return nil, err
	}

	rline := &rline.Rline{
		Inst: l,
		N:    nil,
		C:    nil,
		Out:  stdout,
		Err:  stderr,
		Int:  false,
		Cyg:  false,
		P:    l.SetPrompt,
		S:    l.SaveHistory,
		Pw:   nil,
	}
	return rline, nil
}
