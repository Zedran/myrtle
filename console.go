package main

import (
	"bufio"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

// Current page indicator type
type Page uint8

const (
	EXIT Page = iota
	START_PAGE
	RESULTS_PAGE
	OBJECT_PAGE
)

const (
	LONG_DELAY  = 50 * time.Millisecond
	MED_DELAY   = 25 * time.Millisecond
	SHORT_DELAY = 10 * time.Millisecond

	RES_PER_PAGE int = 20

	TERM_HEIGHT int = 26
	TERM_WIDTH  int = 80
)

type Console struct {
	// Old console dimensions that are restored on exit
	oldH,
	oldW int

	// Currently displayed page
	page Page

	// Current results page
	resPage int

	// Display radius / altitude ASL
	radius bool

	// Display precise / shortened values
	precise bool

	// Search by name / catalog number
	byName bool

	// HTTP client used to fetch data
	client *http.Client

	// Input scanner
	scanner *bufio.Scanner

	// A pointer to current phrase
	phrase *Phrase

	// A list of pointers to found matches
	matches []*Match

	// A pointer to most recently computed object
	curObj *Elements
}

// Main method of struct that launches and drives the interface.
func (c *Console) Run() {
	defer c.clear()
	defer c.restoreTerminalSize()
	c.clear()

	for {
		switch c.page {
		case EXIT:
			return
		case START_PAGE:
			c.showStartPage()
		case RESULTS_PAGE:
			c.showResultsPage()
		case OBJECT_PAGE:
			c.showObjectPage()
		default:
			Logf("Unknown Action value: %d", c.page)
		}
	}
}

// Creates a list of matches according to provided string.
func (c *Console) fetchData(queryString string) {
	var (
		matches []*Match
		err     error
	)
	c.matches = nil

	if c.byName {
		matches, err = Query(c.client, queryString, "")
	} else {
		matches, err = Query(c.client, "", queryString)
	}
	if err != nil {
		pterm.Println(err)
		if err != errShortQuery {
			Log(err)
		}
		return
	}

	if len(matches) > 0 {
		c.matches = matches
		c.resPage = 0
	} else {
		pterm.Println("No matches found.")
	}
}

// Prints the starting page.
func (c *Console) showStartPage() {
	c.clear()

	c.offSetBy(6)

	title, err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("My", pterm.NewStyle(pterm.FgLightGreen)),
		putils.LettersFromStringWithStyle("R", pterm.NewStyle(pterm.FgRed)),
		putils.LettersFromStringWithStyle("TLE", pterm.NewStyle(pterm.FgCyan))).Srender()
	if err != nil {
		Log(err)
	}

	pterm.DefaultCenter.Print(title)

	pterm.DefaultCenter.Print("My Refined TLE Browser")

	c.offSetBy(5)
	pterm.Println(" Enter the name of the object or type '/h' to read help message.")

	c.offSetBy(5)
	c.showSearchDialog()
}

// Prints the page containing the query results.
func (c *Console) showResultsPage() {
	c.clear()
	c.printMatches()
	c.offSetBy(2)

	n := c.pickResult()

	if n == -1 {
		return
	}

	tle := ParseMatch(c.matches[n])
	c.curObj = CalculateElements(tle, M_E, R_E)
	c.nextPage()
}

// Displays the object page containing calculated orbital elements
// for the selected sattellite.
func (c *Console) showObjectPage() {
	c.clear()
	c.printElements(!c.radius, c.precise)
	c.offSetBy(1)

	c.showSearchDialog()
}

// Displays the help page.
func (c *Console) showHelpPage() {
	msg := []string{
		"Commands:\n",
		" /b    - back              |  /e    - exit            |  /f - forward",
		" /a    - display altitude  |  /r    - display radius  |  /p - precise values",
		" /s    - short values      |  /n    - search by name  |  /c - search by cat num",
		" />[n] - next results page |  /<[n] - previous page",
		"\nSymbols:\n",
		" SMa    -  Semi-Major Axis            |  SMi    -  Semi-Minor Axis",
		" PeR/A  -  Periapsis Radius/Altitude  |  ApR/A  -  Apoapsis Radius/Altitude",
		" R/Alt  -  Radius/Altitude            |  Ecc    -  Orbital Eccentricity",
		" T      -  Orbital Period             |  PeT    -  Time to Periapsis",
		" ApT    -  Time to Apoapsis           |  Vel    -  Orbital Velocity",
		" Inc    -  Orbital Inclination        |  LAN    -  Longitude of Ascending Node",
		" LPe    -  Longitude of Periapsis     |  AgP    -  Argument of Periapsis",
		" TrA    -  True Anomaly               |  TrL    -  True Longitude",
		" MnA    -  Mean Anomaly               |  MnL    -  Mean Longitude",
		" EcA    -  Eccentric Anomaly",
		"\nReferences:\n",
		" * CelesTrak - https://celestrak.com",
		" * pterm     - https://github.com/pterm/pterm",
	}

	c.clear()

	for i := range msg {
		pterm.Println(msg[i])
		time.Sleep(SHORT_DELAY)
	}

	c.offSetBy(1)
	c.getInput("Press Enter to continue...")
}

// Gets input from the user and creates a phrase from it.
func (c *Console) getInput(prompt string) *Phrase {
	pterm.Print(prompt + "  ")
	if c.scanner.Scan() {
		return NewPhrase(c.scanner.Text())
	}
	return NewPhrase("/e")
}

// Prompts the user to pick the result and runs commands included inside
// the input.
func (c *Console) pickResult() int {
	length := len(c.matches)

	for {
		phrase := c.getInput("PICK RESULT:")

		if phrase == nil {
			continue
		}

		if !c.runCommands(phrase) {
			return -1
		}

		n, err := strconv.Atoi(phrase.Object)

		if err != nil {
			pterm.Println("Not a number")
			continue
		} else if n <= 0 || n > length {
			pterm.Printfln("Number not in range of the matches pool (1;%d>", length)
			continue
		}

		return n - 1
	}
}

// Displays the search dialog and launches action depending
// on the provided input.
func (c *Console) showSearchDialog() {
	for {
		phrase := c.getInput("SEARCH FOR:")
		if phrase == nil {
			continue
		}

		c.resetFlags()

		if !c.runCommands(phrase) {
			return
		}

		if len(phrase.Object) > 0 {
			c.phrase = phrase
			c.fetchData(c.phrase.Object)
		}

		if c.matches != nil {
			if c.page == OBJECT_PAGE && len(phrase.Object) > 0 {
				c.previousPage()
			} else {
				c.nextPage()
			}
			return
		}
	}
}

// Prints object's orbital elements. If alt is true, the altitude ASL
// is displayed. If acc is true, the values are not shortened
// and are displayed as exact numbers.
func (c *Console) printElements(alt, acc bool) {
	elements := c.curObj.ToString(alt, acc)

	pterm.Println(c.curObj.GetTitle())

	time.Sleep(LONG_DELAY)

	for i := range elements {
		pterm.Println(" ", elements[i])
		time.Sleep(SHORT_DELAY)
	}

	time.Sleep(MED_DELAY)
}

// Displays the current page of found results.
func (c *Console) printMatches() {
	pterm.Printf("RESULTS FOR %s (%d):\n\n", c.phrase.Object, len(c.matches))

	time.Sleep(LONG_DELAY)

	var lastI int

	if len(c.matches) < c.resPage*RES_PER_PAGE+RES_PER_PAGE {
		lastI = len(c.matches)
	} else {
		lastI = c.resPage*RES_PER_PAGE + RES_PER_PAGE
	}

	for i := c.resPage * RES_PER_PAGE; i < lastI; i++ {
		pterm.Printf(
			"%9d |  NAME:%25s     NORAD SIG:%8s  | %2d\n",
			i+1, c.matches[i].Title, c.matches[i].GetCatNum(), i+1,
		)
		time.Sleep(SHORT_DELAY)
	}

	time.Sleep(MED_DELAY)
	pterm.Printf("%66s     %d / %d", "PAGE:", c.resPage+1, int(math.Ceil(float64(len(c.matches)/RES_PER_PAGE)))+1)
	time.Sleep(MED_DELAY)
}

// Skips to the next page.
func (c *Console) nextPage() {
	if c.page < OBJECT_PAGE {
		c.page++
	}
}

// Skips to the previous page.
func (c *Console) previousPage() {
	if c.page > EXIT {
		c.page--
	}
}

// Clear console window.
func (c *Console) clear() {
	pterm.Print("\033[H\033[2J")
}

// Offsets the current cursor position by n lines.
func (c *Console) offSetBy(n int) {
	for i := 0; i < n; i++ {
		pterm.Println()
	}
}

// Restores the remembered terminal window size.
func (c *Console) restoreTerminalSize() {
	c.setWorkingHeight(c.oldH, c.oldW)
}

// Sets the terminal height and width according to the values passed.
func (c *Console) setWorkingHeight(h, w int) {
	pterm.Printfln("\x1b[8;%d;%dt", h, w)
}

// Parses the command that changes the matches currently displayed
// on the result page.
func (c *Console) parseSwitchResPageCommand(cmd string, direction int) {
	if len(cmd) == 1 {
		c.switchResPage(direction)
		return
	}

	n, err := strconv.Atoi(cmd[1:])
	if err != nil {
		return
	}

	c.switchResPage(n * direction)
}

// Sets new results page number.
func (c *Console) switchResPage(by int) {
	newN := c.resPage + by
	if newN >= 0 && newN <= int(math.Ceil(float64(len(c.matches)/RES_PER_PAGE))) {
		c.resPage = newN
	}
}

// Runs commands contained within the Phrase. Returns true if any further
// action should be taken by the calling function (e.g. proceed with query).
func (c *Console) runCommands(phrase *Phrase) bool {
	if Contains(phrase.Commands, "b") {
		c.previousPage()
		return false
	} else if Contains(phrase.Commands, "e") {
		c.page = EXIT
		return false
	} else if Contains(phrase.Commands, "h") {
		c.showHelpPage()
		return false
	}

	switch c.page {
	case START_PAGE:
		if Contains(phrase.Commands, "f") && len(c.matches) > 0 {
			c.page = RESULTS_PAGE
		} else if len(phrase.Object) >= MIN_QLEN {
			c.setFlags(phrase)
			return true
		} else if len(phrase.Object) < MIN_QLEN {
			return false
		}
	case RESULTS_PAGE:
		rightIdx := ContainsPart(phrase.Commands, ">")
		leftIdx := ContainsPart(phrase.Commands, "<")

		if Contains(phrase.Commands, "f") && c.curObj != nil {
			c.page = OBJECT_PAGE
		} else if rightIdx > -1 {
			c.parseSwitchResPageCommand(phrase.Commands[rightIdx], 1)
		} else if leftIdx > -1 {
			c.parseSwitchResPageCommand(phrase.Commands[leftIdx], -1)
		} else {
			return true
		}
	case OBJECT_PAGE:
		if len(phrase.Object) >= MIN_QLEN || len(phrase.Commands) > 0 {
			c.setFlags(phrase)
			return true
		}
	}
	return false
}

// Sets display flags according to commands contained within the Phrase.
func (c *Console) setFlags(phrase *Phrase) {
	if Contains(phrase.Commands, "a") {
		c.radius = false
	} else if Contains(phrase.Commands, "r") {
		c.radius = true
	}

	if Contains(phrase.Commands, "p") {
		c.precise = true
	} else if Contains(phrase.Commands, "s") {
		c.precise = false
	}
}

// Resets display flags.
func (c *Console) resetFlags() {
	c.radius = true
	c.precise = false
	c.byName = true
}

// Sets up the new console interface.
func NewConsole(client *http.Client) *Console {
	var c Console

	c.oldH = pterm.GetTerminalHeight()
	c.oldW = pterm.GetTerminalWidth()

	c.page = START_PAGE

	c.client = client
	c.scanner = bufio.NewScanner(os.Stdin)

	c.resetFlags()

	c.setWorkingHeight(TERM_HEIGHT, TERM_WIDTH)

	return &c
}
