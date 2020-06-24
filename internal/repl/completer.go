package repl

import "github.com/c-bata/go-prompt"

func NewCompleter() func(prompt.Document) []prompt.Suggest {
	return func(document prompt.Document) []prompt.Suggest {
		return nil
	}
}
