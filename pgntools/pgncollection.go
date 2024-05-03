/*
  pgncollection.go
  Description: Simple tools for handling collections of games in PGN
               format
  -----------------------------------------------------------------------------

  Started on  <Sat May  9 16:50:49 2015 Carlos Linares Lopez>
  Last update <sábado, 07 mayo 2016 17:26:19 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package pgntools

import (
	"fmt"
	"io"  // io streams
	"log" // logging services

	// access to file mgmt functions
	// pgn files are parsed with a regexp
	// to convert integers into strings
	"text/template" // go facility for processing templates

	// import the parser of propositional formulae

	// import my favourite package to automatically create tables
	"github.com/clinaresl/table"
)

// typedefs
// ----------------------------------------------------------------------------

// // PGN games can be sorted either in ascending or descending order. The
// // direction is then defined as an integer
// type sortingDirection int

// // A pgnSorting consists of two items: a constant value for distinguishing
// // ascending from descending order and a variable name used as a key for sorting
// // pgn games
// type pgnSorting struct {
// 	direction sortingDirection
// 	variable  string
// }

// A PgnCollection consists of an arbitrary number of PgnGames along with a
// count of the number of games stored in it ---this is given to check for
// consistency so that the difference between nbGames and len (slice) shall be
// always null.
//
// In addition, a PGN collection contains a sort descriptor which consists of a
// slice of pairs that contain for each variable whether PGN games should be
// sorted in increasing or decreasing order
type PgnCollection struct {
	slice []PgnGame
	// sortDescriptor []pgnSorting
	nbGames int
}

// // A histogram is indexed by keys. Keys are either variables (represented as a
// // string) or a slice of cases (each defined with a string as well). Both
// // variables and cases can be qualified with a title

// // A variable just consists of an association of a title and the name of a
// // variable
// type pgnKeyVar struct {
// 	title    string
// 	variable string
// }

// // A case consists of a slice of structs similar to variables but, instead of
// // variables, they store propositional expressions
// type pgnKeyCase struct {
// 	title      string
// 	expression string
// }

// // A full specification of cases consists just of a slice of cases. The whole
// // collection of cases can be also qualified with a title
// type pgnKeyCases struct {
// 	title       string
// 	expressions []pgnKeyCase
// }

// // So far, pgn histogram registers are any structs that support the following
// // operations: get title, get subtitle, get key and get value. Note: first,
// // while titles are always strings, subtitles can be of any type and they should
// // be indeed sorted according to their type; second, values should be supported
// // by the current implementation of histograms
// type pgnHistogramRegister interface {
// 	GetTitle() string
// 	GetSubtitle(game *PgnGame) dataInterface
// 	GetKey(game *PgnGame) string
// 	GetValue(game *PgnGame) dataHistValue
// }

// consts
// ----------------------------------------------------------------------------

// // PGN games can be sorted either in ascending or descending order
// const (
// 	increasing sortingDirection = 1 << iota // increasing order
// 	decreasing                              // decreasing order
// )

// Methods
// ----------------------------------------------------------------------------

// -- Accessors

// the following are getters over the attributes of a PgnCollection

// Return all games as instances of PgnGame that are stored in this particular
// collection
func (games *PgnCollection) GetGames() []PgnGame {
	return games.slice
}

// Return the index-th game stored in this particular collection. The result is
// undefined in case an index out of bounds is passed
func (games *PgnCollection) GetGame(index int) PgnGame {
	return games.slice[index]
}

// Return the number of items in the collection
func (games PgnCollection) Len() int {
	return games.nbGames
}

// Methods
// ----------------------------------------------------------------------------

// Play this collection of games on the given writer showing the board every
// number of given plies
func (c PgnCollection) Play(plies int, writer io.Writer) {

	// Just play each game
	for _, igame := range c.slice {

		// First, show the tags of this game only in case a strictly positive
		// value is given
		if plies > 0 {
			for name, value := range igame.GetTags() {
				io.WriteString(writer, fmt.Sprintf(" %v: %v\n", name, value))
			}
			io.WriteString(writer, "\n")
		}

		// Show the game
		igame.Play(plies, writer)
	}
}

// -- Sorting

// The following methods ease the task of sorting games in a collection

// // Swap two games within the same collection
// func (games PgnCollection) Swap(i, j int) {
// 	games.slice[i], games.slice[j] = games.slice[j], games.slice[i]
// }

// // This method creates a valid PGN descriptor used for sorting games from a
// // string specification. The string contains pairs of the form (<|>) and
// // %variable and there can be an arbitrary number of them. The first item is
// // used to decide whether to sort games in ascending or descending order; the
// // second one is used to decide what variable to use as a key.
// func (games *PgnCollection) GetSortDescriptor(sortString string) []pgnSorting {

// 	// extract all sorting criteria given in the string
// 	for reSortingCriteria.MatchString(sortString) {

// 		// extract the two groups in the sorting criteria: the direction
// 		// and the key
// 		tag := reSortingCriteria.FindStringSubmatchIndex(sortString)
// 		direction, key := sortString[tag[2]:tag[3]], sortString[tag[4]:tag[5]]

// 		// and move forward in the string
// 		sortString = sortString[tag[1]:]

// 		// store the direction and key in this collection
// 		var newSorting pgnSorting
// 		if direction == "<" {
// 			newSorting = pgnSorting{increasing, key}

// 		} else if direction == ">" {
// 			newSorting = pgnSorting{decreasing, key}
// 		} else {
// 			log.Fatalf(" An unknown sorting direction has been found: '%v'", direction)
// 		}
// 		games.sortDescriptor = append(games.sortDescriptor, newSorting)
// 	}

// 	// make sure here that the full sort descriptor was successfully processed
// 	if len(sortString) > 0 {
// 		log.Fatalf(" There was an error in the sort string at point '%v'", sortString)
// 	}

// 	// and return the descriptor
// 	return games.sortDescriptor
// }

// // Return true if the i-th game should be before the j-th game and false
// // otherwise
// func (games PgnCollection) Less(i, j int) bool {

// 	// go over all items of the sort descriptor stored in this collection
// 	// until either the slice is over or a decision has been made whether
// 	// the i-th game should be before or after the j-th game
// 	for _, descriptor := range games.sortDescriptor {

// 		// first of all, check this variable exists in both games
// 		icontent, ok := games.slice[i].tags[descriptor.variable]
// 		if !ok {
// 			log.Fatalf("'%v' is not a variable and can not be used for sorting games",
// 				descriptor.variable)
// 		}
// 		jcontent, ok := games.slice[j].tags[descriptor.variable]
// 		if !ok {
// 			log.Fatalf("'%v' is not a variable and can not be used for sorting games",
// 				descriptor.variable)
// 		}

// 		// check the direction and then the variable to use
// 		if descriptor.direction == increasing {
// 			if icontent.Less(jcontent) {
// 				return true
// 			}
// 			if icontent.Greater(jcontent) {
// 				return false
// 			}
// 		} else if descriptor.direction == decreasing {
// 			if icontent.Greater(jcontent) {
// 				return true
// 			}
// 			if icontent.Less(jcontent) {
// 				return false
// 			}
// 		} else {
// 			log.Fatalf(" Unknown sorting direction '%v'", descriptor.direction)
// 		}
// 	}

// 	// if the sorting descriptor was exhausted, then return true by default
// 	return true
// }

// // -- Histograms

// // The following methods are defined to qualify various types to be used for
// // generating histograms

// // Just return the value of the title of this key variable
// func (key pgnKeyVar) GetTitle() string {
// 	return key.title
// }

// // Just return the title of this case specification
// func (key pgnKeyCases) GetTitle() string {
// 	return key.title
// }

// // Just return the subtitle to use with this variable. In the case of key
// // variables, the subtitle is just the value for a specific game of the given
// // variable
// //
// // Importantly, this method returns instances of a generic type so that
// // subtitles can be more naturally sorted
// func (key pgnKeyVar) GetSubtitle(game *PgnGame) dataInterface {

// 	// Key variables are expected to be found in the tags of a chess game
// 	value, ok := game.GetTagValue(key.variable)
// 	if ok != nil {
// 		log.Fatalf(" It was not possible to access the subtitle of key '%v'\n", key.variable)
// 	}

// 	// in case the value was successfully accessed, just return it
// 	return value
// }

// // Return the subtitle to use for this case specification. The importance of
// // cases to be disjoint and to fully enumerate all cases stems from the fact
// // that this service should have one and only one case whose expression
// // evaluates to true, with all the other evaluating to false. Thus, this service
// // simply returns the title of the expression that is verified for this specific
// // game
// //
// // This method always return a string which is the title of the only case that
// // is verified by this specific game
// func (key pgnKeyCases) GetSubtitle(game *PgnGame) dataInterface {

// 	var ititle string

// 	// first, start by creating a symbol table with all the information
// 	// appearing in the headers of this game
// 	symtable := make(map[string]pfparser.RelationalInterface)
// 	for key, content := range game.tags {

// 		// first, verify whether this is an integer
// 		value, ok := content.(constInteger)
// 		if ok {

// 			symtable[key] = pfparser.ConstInteger(value)
// 		} else {

// 			// if not, check if it is a string
// 			value, ok := content.(constString)
// 			if ok {
// 				symtable[key] = pfparser.ConstString(value)
// 			} else {
// 				log.Fatal(" Unknown type")
// 			}
// 		}
// 	}

// 	// for all cases in this specification
// 	for _, icase := range key.expressions {

// 		// verify whether this expression is verified
// 		var err error
// 		var logEvaluator pfparser.LogicalEvaluator

// 		// parse and evaluate this particular case
// 		iexpression := icase.expression
// 		logEvaluator, err = pfparser.Parse(&iexpression, 0)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if logEvaluator.Evaluate(symtable) == pfparser.TypeBool(true) {
// 			ititle = icase.title
// 		}
// 	}

// 	return constString(ititle)
// }

// // Just return the key to use for storing this entry in a histogram. For key
// // variables, keys are exactly the same than subtitles with a slight
// // difference. They are always stored as strings.
// func (keyin pgnKeyVar) GetKey(game *PgnGame) (keyout string) {

// 	// just compute the subtitle corresponding to this key variable and
// 	// return it as a string
// 	subtitle := keyin.GetSubtitle(game)

// 	// go through various type assertions to convert the subtitle to the
// 	// right type and, from it, to a string
// 	value, ok := subtitle.(constInteger)
// 	if ok {
// 		keyout = strconv.Itoa(int(value))
// 	} else {

// 		value, ok := subtitle.(constString)
// 		if ok {
// 			keyout = string(value)
// 		} else {
// 			log.Fatalf(" Unknown type of '%v'\n", keyin.variable)
// 		}
// 	}

// 	// and now return the key
// 	return
// }

// // Return the key to use for storing this entry in a histogram. For key cases,
// // keys are exactly the same than subtitles.
// func (keyin pgnKeyCases) GetKey(game *PgnGame) (keyout string) {

// 	// just compute the subtitle corresponding to this key case and return
// 	// it as a string
// 	subtitle := keyin.GetSubtitle(game)
// 	value, ok := subtitle.(constString)
// 	if !ok {
// 		log.Fatalf(" A subtitle of a type different than string was returned!")
// 	}
// 	return string(value)
// }

// // The value of a key variable is always equal to one. This results from the
// // fact that histograms store the association (key, value) so that a value=1
// // means that one occurrence of a specific key has been observed
// func (key pgnKeyVar) GetValue(game *PgnGame) dataHistValue {
// 	return 1
// }

// // The value of a key case is always equal to one. This results from the fact
// // that histograms store the association (key, value) so that a value=1 means
// // that one occurrence of a specific key has been observed
// func (key pgnKeyCases) GetValue(game *PgnGame) dataHistValue {
// 	return 1
// }

// // This function processes the histogram command line provided by the user and
// // returns a slice of pgn histogram registers that can then be used to generate
// // the histogram of any collection of chess games.
// func parseHistCommandLine(histCommandLine string) (histDirective []pgnHistogramRegister) {

// 	// extract all histogram directives given in the histogram command line
// 	for reHistogramCmdVar.MatchString(histCommandLine) ||
// 		reHistogramCmdCase.MatchString(histCommandLine) {

// 		// in case the following directive is recognized as a variable
// 		if reHistogramCmdVar.MatchString(histCommandLine) {

// 			// extract the two groups in a variable: the title and
// 			// the variable name
// 			tag := reHistogramCmdVar.FindStringSubmatchIndex(histCommandLine)
// 			title := histCommandLine[tag[2]:tag[3]]
// 			variable := histCommandLine[tag[4]:tag[5]]

// 			// as this has been recognized to be a key variable, a new
// 			// instance of key variables is created and its fields are
// 			// filled in
// 			newRegister := pgnKeyVar{title, variable}
// 			histDirective = append(histDirective, newRegister)

// 			// and move forward in the string
// 			histCommandLine = histCommandLine[tag[1]:]
// 		} else if reHistogramCmdCase.MatchString(histCommandLine) {

// 			tag := reHistogramCmdCase.FindStringSubmatchIndex(histCommandLine)

// 			// extract the title and the definition with all cases
// 			title := histCommandLine[tag[2]:tag[3]]
// 			cases := histCommandLine[tag[4]:tag[5]]

// 			// create an empty slice of cases
// 			var expressions []pgnKeyCase

// 			// process each case separately
// 			for reHistogramCmdSubcase.MatchString(cases) {

// 				subtag := reHistogramCmdSubcase.FindStringSubmatchIndex(cases)

// 				// create a case and add it to the slice of
// 				// cases specifying this title and this case
// 				expressions = append(expressions,
// 					pgnKeyCase{cases[subtag[2]:subtag[3]],
// 						cases[subtag[4]:subtag[5]]})

// 				// and move forward in the string
// 				cases = cases[subtag[1]:]
// 			}

// 			// create a new instance of cases to store this
// 			// specification
// 			histDirective = append(histDirective,
// 				pgnKeyCases{title, expressions})

// 			// and move forward in the original string
// 			histCommandLine = histCommandLine[tag[1]:]
// 		} else {
// 			log.Fatalf(" Syntax error in the histogram directive: '%v'\n",
// 				histCommandLine)
// 		}
// 	}

// 	// and return the slice of registers computed so far
// 	return
// }

// // Compute a histogram with the information given in the specified histogram
// // command line using the games stored in the receiver. It returns an instance
// // of a histogram
// func (games *PgnCollection) ComputeHistogram(histCommandLine string) Histogram {

// 	// create a new histogram
// 	hist := NewHistogram()

// 	// process the histogram command line to get the registers with the
// 	// information of every directive provided by the user
// 	histRegisters := parseHistCommandLine(histCommandLine)

// 	// process all games in the current collection
// 	for _, game := range games.slice {

// 		// create a key to be used to access the histogram. Make it
// 		// initially empty
// 		var key []string

// 		// and now, for every register
// 		for _, register := range histRegisters {

// 			// retrieve the value of the variable in this register
// 			key = append(key, register.GetKey(&game))
// 		}

// 		// and now, annotate that one sample was observed for this
// 		// particular key
// 		hist.Increment(key, 1)
// 	}

// 	return hist
// }

// Templates
//
// All the following methods are used to handle templates both for generating
// ascii and LaTeX output
// ----------------------------------------------------------------------------

// This function is used in text/templates and it is the equivalent to the
// homonym function defined for PgnGame.
//
// It returns the empty string if the given name is not defined as a tag in the
// first game of this collection and the value of the tag (casted to a string)
// otherwise. This function assumes that all games within the same collection
// share the same tags.
func (games *PgnCollection) GetTagValue(name string) string {

	// first, attempt at reading the specified tag from the first game of
	// the collection
	val, err := games.slice[0].GetTagValue(name)

	// if such tag does not exist, then return the empty string
	if err != nil {
		return ""
	}

	// otherwise, return its value
	return val
}

// This is an auxiliary function used in text/templates to generate slices of
// strings to be given as argument to other methods
func (games *PgnCollection) GetSlice(fields ...any) []any {
	return fields
}

// Returns a table according to the specification given in first place. Columns
// are populated with the tags given in fields. It is intended to be used in
// ascii table templates
func (games *PgnCollection) GetTable(specline string, fields []any) table.Table {

	// Create a table according to the given specification
	table, err := table.NewTable(specline)
	if err != nil {
		log.Fatal(" Fatal error while constructing the table in PgnCollection.GetTable")
	}

	// Add the header
	table.AddThickRule()
	table.AddRow(fields...)
	table.AddDoubleRule()

	// Now, add a row per game
	for idx, game := range games.slice {

		// show a separator every ten lines to make the table easier to
		// read
		if idx > 0 && idx%10 == 0 {
			table.AddSingleRule()
		}

		// and show here the information from the specified fields for
		// this game
		table.AddRow(game.getFields(fields)...)
	}

	// End the table and return the table as a string
	table.AddThickRule()
	return *table
}

// Writes into the specified writer the result of instantiating the given
// template file with information of all games in this collection. The template
// acknowledges all tags of a pgngame plus others. For a full description, see
// the manual.
func (games *PgnCollection) GamesToWriterFromTemplate(dst io.Writer, templateFile string) {

	// access a template and parse its contents
	template, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatal(err)
	}

	// and now execute the template
	err = template.Execute(dst, games)
	if err != nil {
		log.Fatal(err)
	}
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
