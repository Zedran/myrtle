package main

import "strings"

const COMMAND_PREFIX = "/"

/* Phrase passed by the user when asked for input. */
type Phrase struct {
	// Name or catalog number of an object
	Object   string
	
	// A slice of commands passed
	Commands []string
}

/* Creates a new Phrase from string input. */
func NewPhrase(queryString string) *Phrase {
	if len(queryString) == 0 {
		return nil
	}
	
	var phrase Phrase

	words := strings.Split(queryString, " ")

	for i := range words {
		if strings.HasPrefix(words[i], COMMAND_PREFIX) {
			phrase.Commands = append(phrase.Commands, extractCommands(words[i])...)
		} else if len(phrase.Object) == 0 {
			phrase.Object = words[i]
		}
	}

	phrase.Commands = RemoveDuplicates(phrase.Commands)

	return &phrase
}

/* Extracts commands from a passed space-separated word. */
func extractCommands(word string) []string {
	commands := make([]string, 0)

	for _, c := range strings.Split(word, COMMAND_PREFIX) {
		if len(strings.TrimSpace(c)) > 0 {
			commands = append(commands, c)
		}
	}

	return commands
}
