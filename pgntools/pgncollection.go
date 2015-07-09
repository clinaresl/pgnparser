/* 
  pgncollection.go
  Description: Simple tools for handling collections of games in PGN
               format
  ----------------------------------------------------------------------------- 

  Started on  <Sat May  9 16:50:49 2015 Carlos Linares Lopez>
  Last update <jueves, 09 julio 2015 08:10:29 Carlos Linares Lopez (clinares)>
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

// A histogram is indexed by keys. Keys are either variables (represented as a
// string) or a slice of cases (each defined with a string as well). Both
// variables and cases can be qualified with a title
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

// consts
// ----------------------------------------------------------------------------

// PGN games can be sorted either in ascending or descending order
const (
	increasing sortingDirection = 1 << iota		// increasing order
	decreasing					// decreasing order
)

// Methods
// ----------------------------------------------------------------------------

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
