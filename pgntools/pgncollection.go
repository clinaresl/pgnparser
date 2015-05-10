/* 
  pgncollection.go
  Description: Simple tools for handling collections of games in PGN
               format
  ----------------------------------------------------------------------------- 

  Started on  <Sat May  9 16:50:49 2015 Carlos Linares Lopez>
  Last update <domingo, 10 mayo 2015 02:17:07 Carlos Linares Lopez (clinares)>
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

// the following regexps are used just to locate the main body of the
// LaTeX template
var reBeginDocument = regexp.MustCompile (`\\begin{document}`)
var reEndDocument = regexp.MustCompile (`\\end{document}`)

// typedefs
// ----------------------------------------------------------------------------
type PgnCollection struct {

	slice []PgnGame                                  // collection of games
	nbGames int;                                  // number of games stored
}

// Methods
// ----------------------------------------------------------------------------

// the following are getters over the attributes of a PgnCollection
func (games *PgnCollection) GetGames () []PgnGame {
	return games.slice
}

func (games *PgnCollection) GetGame (index int) PgnGame {
	return games.slice [index]
}

func (games *PgnCollection) GetNbGames () int {
	return games.nbGames
}

// ShowHeaders
// 
// returns a string with a summary of the information of all games
// stored in this collection
// ----------------------------------------------------------------------------
func (games *PgnCollection) ShowHeaders () string {

	// show the header
	output := " |  DBGameNo  | Date                | White                     | Black                     | ECO | Time  | Moves | Result |\n +------------+---------------------+---------------------------+---------------------------+-----+-------+-------+--------+\n"

	// and now, add to output information of every single game in the given
	// collection
	for _, game := range games.slice {
		output += game.ShowHeader () + "\n"
	}

	// and add a bottom line
	output += " +------------+---------------------+---------------------------+---------------------------+-----+-------+-------+--------+"

	// and return the string
	return output
}

// GamesToLaTeXFromString
// 
// produces LaTeX code using the specified template with information
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
// ----------------------------------------------------------------------------
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

// GamesToLaTeXFromFile
//
// produces LaTeX code using the template stored in the specified file
// with information of all games in this collection. The string
// acknowledges various placeholders which have the format
// '%<name>'. All tag names specified in this game are
// acknowledged. Additionally, '%moves' is substituted by the list of
// moves
// ----------------------------------------------------------------------------
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
