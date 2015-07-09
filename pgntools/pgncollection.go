/* 
  pgncollection.go
  Description: Simple tools for handling collections of games in PGN
               format
  ----------------------------------------------------------------------------- 

  Started on  <Sat May  9 16:50:49 2015 Carlos Linares Lopez>
  Last update <jueves, 09 julio 2015 17:52:37 Carlos Linares Lopez (clinares)>
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
	"log"			// logging services
	"fmt"			// printing msgs
	"regexp"                // pgn files are parsed with a regexp
	"strconv"		// to convert integers into strings

	// import a user package to manage paths
	"bitbucket.org/clinares/pgnparser/fstools"
)

// global variables
// ----------------------------------------------------------------------------

// the following regexp is used to match the different sorting criteria

// sorting criteria consists of a sorting direction and a particular variable
// used as a key for sorting games. The direction is specified with either < or
// > meaning increasing and decreasing order respectively; the variable to use
// is preceded by '%' (there is no need actually to use that prefix and this is
// done only for the sake of consistency across different commands of pgnparser)
var reSortingCriteria = regexp.MustCompile (`^\s*(<|>)\s*%([A-Za-z]+)\s*`)

// the following regexp is used to process histogram command lines
var reHistogramCmd = regexp.MustCompile (`^\s*([A-Za-z0-9]+)\s*:\s*%([A-Za-z]+)\s*`)

// the following regexps are used just to locate the main body of the
// LaTeX template

// It is used to locate the beginning of the document of the LaTeX template
var reBeginDocument = regexp.MustCompile (`\\begin{document}`)

// It is used to locate the beginning of the document of the LaTeX template
var reEndDocument = regexp.MustCompile (`\\end{document}`)

// typedefs
// ----------------------------------------------------------------------------

// PGN games can be sorted either in ascending or descending order. The
// direction is then defined as an integer
type sortingDirection int

// A pgnSorting consists of two items: a constant value for distinguishing
// ascending from descending order and a variable name used as a key for sorting
// pgn games
type pgnSorting struct {
	direction sortingDirection
	variable string
}

// A PgnCollection consists of an arbitrary number of PgnGames along with a
// count of the number of games stored in it ---this is given to check for
// consistency so that the difference between nbGames and len (slice) shall be
// always null.

// In addition, a PGN collection contains a sort descriptor which consists of a
// slice of pairs that contain for each variable whether PGN games should be
// sorted in increasing or decreasing order
type PgnCollection struct {

	slice []PgnGame
	sortDescriptor []pgnSorting
	nbGames int;
}

// A histogram is indexed by keys. Keys are either variables (represented as a
// string) or a slice of cases (each defined with a string as well). Both
// variables and cases can be qualified with a title

// A variable just consists of an association of a title and the name of a
// variable
type pgnKeyVar struct {
	title string
	variable string
}

// A case consists of a slice of structs similar to variables but, instead of
// variables, they store propositional expressions
type pgnKeyCase struct {
	title string
	expression string
}

// A full specification of cases consists just of a slice of cases. The whole
// collection of cases can be also qualified with a title
type pgnKeyCases struct {
	title string
	expressions []string
}

// So far, pgn histogram registers are any structs that support the following
// operations: get title, get subtitle, get key and get value. Note: first,
// while titles are always strings, subtitles can be of any type and they should
// be indeed sorted according to their type; second, values should be supported
// by the current implementation of histograms
type pgnHistogramRegister interface {
	GetTitle () string;
	GetSubtitle (game *PgnGame) dataInterface;
	GetKey (game *PgnGame) string;
	GetValue (game *PgnGame) dataHistValue;
}

// consts
// ----------------------------------------------------------------------------

// PGN games can be sorted either in ascending or descending order
const (
	increasing sortingDirection = 1 << iota		// increasing order
	decreasing					// decreasing order
)

// Methods
// ----------------------------------------------------------------------------

// -- Accessors

// the following are getters over the attributes of a PgnCollection

// Return all games as instances of PgnGame that are stored in this particular
// collection
func (games *PgnCollection) GetGames () []PgnGame {
	return games.slice
}

// Return the index-th game stored in this particular collection
func (games *PgnCollection) GetGame (index int) PgnGame {
	return games.slice [index]
}

// Return the number of items in the collection
func (games PgnCollection) Len () int {
	return games.nbGames
}

// -- Sorting

// The following methods ease the task of sorting games in a collection

// Swap two games within the same collection
func (games PgnCollection) Swap (i, j int) {
	games.slice[i], games.slice[j] = games.slice[j], games.slice[i]
}

// This method creates a valid PGN descriptor used for sorting games from a
// string specification. The string contains pairs of the form (<|>) and
// %variable and there can be an arbitrary number of them. The first item is
// used to decide whether to sort games in ascending or descending order; the
// second one is used to decide what variable to use as a key.
func (games *PgnCollection) GetSortDescriptor (sortString string) []pgnSorting {

	// extract all sorting criteria given in the string
	for ;reSortingCriteria.MatchString (sortString); {

		// extract the two groups in the sorting criteria: the direction
		// and the key
		tag := reSortingCriteria.FindStringSubmatchIndex (sortString)
		direction, key := sortString[tag[2]:tag[3]], sortString[tag[4]:tag[5]]

		// and move forward in the string
		sortString = sortString[tag[1]:]

		// store the direction and key in this collection
		var newSorting pgnSorting
		if direction == "<" {
			newSorting = pgnSorting {increasing, key}

		} else if direction == ">" {
			newSorting = pgnSorting {decreasing, key}
		} else {
			log.Fatalf (" An unknown sorting direction has been found: '%v'", direction)
		}
		games.sortDescriptor = append (games.sortDescriptor, newSorting)
	}

	// make sure here that the full sort descriptor was successfully processed
	if len (sortString) > 0 {
		log.Fatalf (" There was an error in the sort string at point '%v'", sortString)
	}

	// and return the descriptor
	return games.sortDescriptor
}

// Return true if the i-th game should be before the j-th game and false
// otherwise
func (games PgnCollection) Less (i, j int) bool {

	// go over all items of the sort descriptor stored in this collection
	// until either the slice is over or a decision has been made whether
	// the i-th game should be before or after the j-th game
	for _, descriptor := range games.sortDescriptor {

		// first of all, check this variable exists in both games
		icontent, ok := games.slice[i].tags[descriptor.variable]; if !ok {
			log.Fatalf ("'%v' is not a variable and can not be used for sorting games",
				descriptor.variable)
		}
		jcontent, ok := games.slice[j].tags[descriptor.variable]; if !ok {
			log.Fatalf ("'%v' is not a variable and can not be used for sorting games",
				descriptor.variable)
		}
		
		// check the direction and then the variable to use
		if descriptor.direction == increasing {
			if icontent.Less (jcontent) {
				return true
			}
			if icontent.Greater (jcontent) {
				return false
			}
		} else if descriptor.direction == decreasing {
			if icontent.Greater (jcontent) {
				return true
			}
			if icontent.Less (jcontent) {
				return false
			}
		} else {
			log.Fatalf (" Unknown sorting direction '%v'", descriptor.direction)
		}
	}

	// if the sorting descriptor was exhausted, then return true by default
	return true
}

// -- Histograms

// The following methods are defined to qualify various types to be used for
// generating histograms

// Just return the value of the title of this key variable
func (key pgnKeyVar) GetTitle () string {
	return key.title
}

// Just return the subtitle to use with this variable. In the case of key
// variables, the subtitle is just the value for a specific game of the given
// variable
// 
// Importantly, this method returns instances of a generic type so that
// subtitles can be more naturally sorted
func (key pgnKeyVar) GetSubtitle (game *PgnGame) dataInterface {

	// Key variables are expected to be found in the tags of a chess game
	value, ok := game.GetTagValue (key.variable); if ok != nil {
		log.Fatalf (" It was not possible to access the subtitle of key '%v'\n", key.variable)
	}

	// in case the value was successfully accessed, just return it
	return value
}

// Just return the key to use for storing this entry in a histogram. For key
// variables, keys are exactly the same than subtitles with a slight
// difference. They are always stored as strings.
func (keyin pgnKeyVar) GetKey (game *PgnGame) (keyout string) {

	// just compute the subtitle corresponding to this key variable and
	// return it as a string
	subtitle := keyin.GetSubtitle (game)

	// go through various type assertions to convert the subtitle to the
	// right type and, from it, to a string
	value, ok := subtitle.(constInteger); if ok {
		keyout = strconv.Itoa (int (value))
	} else {

		value, ok := subtitle.(constString); if ok {
			keyout = string (value)
		} else {
			log.Fatalf (" Unknown type of '%v'\n", keyin.variable)
		}
	}

	// and now return the key
	return
}

// The value of a key variable is always equal to one. This results from the
// fact that histograms store the association (key, value) so that a value=1
// means that one occurrence of a specific key has been observed
func (key pgnKeyVar) GetValue (game *PgnGame) dataHistValue {
	return 1
}

// This function processes the histogram command line provided by the user and
// returns a slice of pgn histogram registers that can then be used to generate
// the histogram of any collection of chess games.
func parseHistCommandLine (histCommandLine string) (histDirective []pgnHistogramRegister) {

	// extract all histogram directives given in the histogram command line
	for ;reHistogramCmd.MatchString (histCommandLine); {

		// extract the two groups in a variable: the title and the
		// variable name
		tag := reHistogramCmd.FindStringSubmatchIndex (histCommandLine)
		title, variable := histCommandLine[tag[2]:tag[3]], histCommandLine[tag[4]:tag[5]]

		// as this has been recognized to be a key variable, a new
		// instance of key variables is created and its fields are
		// filled in
		newRegister := pgnKeyVar {title, variable}
		histDirective = append (histDirective, newRegister)
		
		// and move forward in the string
		histCommandLine = histCommandLine[tag[1]:]
	}

	// verify that the whole histogram command line was processed
	if histCommandLine != "" {
		log.Fatalf (" Syntax error in the histogram directive: '%v'\n", histCommandLine)
	}

	// and return the slice of registers computed so far
	return
}

// Compute a histogram with the information given in the specified histogram
// command line. It returns an instance of a histogram
func (games *PgnCollection) ComputeHistogram (histCommandLine string) Histogram {

	// create a new histogram
	hist := NewHistogram ()

	// process the histogram command line to get the registers with the
	// information of every directive provided by the user
	histRegisters := parseHistCommandLine (histCommandLine)
	
	// process all games in the current collection
	for _, game := range (games.slice) {

		// create a key to be used to access the histogram. Make it
		// initially empty
		var key []string
		
		// and now, for every register
		for _, register := range (histRegisters) {

			// retrieve the value of the variable in this register
			key = append (key, register.GetKey (&game))
		}

		// and now, annotate that one sample was observed for this
		// particular key
		hist.Increment (key, 1)
	}

	return hist
}

// Returns a string with a summary of the information of all games stored in
// this collection. The summary is shown as an ASCII table with heading and
// bottom lines.
//
// In case any required data is not found, a fatal error is raised
func (games *PgnCollection) ShowHeaders () string {

	// show the header
	output := " |  DBGameNo  | Date                | White                     | Black                     | ECO | Time  | Moves | Result |\n +------------+---------------------+---------------------------+---------------------------+-----+-------+-------+--------+\n"

	// and now, add to output information of every single game in the given
	// collection
	for _, game := range games.slice {
		output += game.showHeader () + "\n"
	}

	// and add a bottom line
	output += " +------------+---------------------+---------------------------+---------------------------+-----+-------+-------+--------+"

	// and return the string
	return output
}

// Produces LaTeX code using the specified template with information
// of all games in this collection. The string acknowledges various
// placeholders which have the format '%<name>'. All tag names
// specified in this game are acknowledged. Additionally, '%moves' is
// substituted by the list of moves.
//
// To generate various games in the same latex file, it processes the
// template and separates the main body of the document from the
// preamble. The preamble is then shown only once and the main body is
// repeated with as many games are found in this collection, all of
// them separated by `\clearpage`
func (games *PgnCollection) GamesToLaTeXFromString (template string) string {

	// locate the begin of the document
	if !reBeginDocument.MatchString (template) {
		log.Fatalf (" The begin of the document has not been found")
	}
	tagBegin := reBeginDocument.FindStringSubmatchIndex (template)

	// and the end of the document
	if !reEndDocument.MatchString (template) {
		log.Fatalf (" The end of the document has not been found")
	}
	tagEnd := reEndDocument.FindStringSubmatchIndex (template)

	// now, initialize the latex file with the preamble and make a
	// copy of the body
	output := template[:tagBegin[1]]
	mainBody := template[tagBegin[1]:tagEnd[0]]
	
	// now, process each game in succession
	for _, game := range games.slice {

		// just performing substitutions in the main body and
		// adding '\clearpage' after every game
		output += fmt.Sprintf(`%v
\clearpage
`, game.replacePlaceholders (mainBody))
	}

	// and end the document
	output += template[tagEnd[0]:]
	
	// and return the string
	return output
}

// Produces LaTeX code using the template stored in the specified file
// with information of all games in this collection. The string
// acknowledges various placeholders which have the format
// '%<name>'. All tag names specified in this game are
// acknowledged. Additionally, '%moves' is substituted by the list of
// moves
func (games *PgnCollection) GamesToLaTeXFromFile (templateFile string) string {

	// Open and read the given file and retrieve its contents
	contents := fstools.Read (templateFile, -1)
	template := string (contents[:len (contents)])

	// and now, just return the results of parsing these contents
	return games.GamesToLaTeXFromString (template)
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
