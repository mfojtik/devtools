package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

var historyFn = filepath.Join(os.TempDir(), ".analize_pr_history")

func ReadCommand(ctx Context) (string, error) {
	line := liner.NewLiner()
	defer func() {
		if err := line.Close(); err != nil {
			log.Printf("unable to close: %v", err)
		}
	}()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(historyFn); err == nil {
		if _, err := line.ReadHistory(f); err != nil {
			return "", err
		}
		if err := f.Close(); err != nil {
			return "", err
		}
	}

	rawInput := ""
	if name, err := line.Prompt(fmt.Sprintf("[#%d]> ", *ctx.PullRequest.Number)); err == nil {
		line.AppendHistory(name)
		rawInput = strings.TrimSpace(name)
	} else if err == liner.ErrPromptAborted || err == io.EOF {
		log.Print("Bye!")
		os.Exit(0)
		return "", err
	} else {
		return "", err
	}

	if len(rawInput) == 0 {
		return "", nil
	}

	if f, err := os.Create(historyFn); err != nil {
		return "", err
	} else {
		defer f.Close()
		_, historyErr := line.WriteHistory(f)
		return rawInput, historyErr
	}
}
