/* 
  pgncollection.go
  Description: Simple tools for handling collections of games in PGN
               format
  ----------------------------------------------------------------------------- 

  Started on  <Sat May  9 16:50:49 2015 Carlos Linares Lopez>
  Last update <sábado, 09 mayo 2015 17:00:07 Carlos Linares Lopez (clinares)>
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
	"fmt"			// printing msgs	
)

// typedefs
// ----------------------------------------------------------------------------
type PgnCollection struct {

	slice []PgnGame                                  // collection of games
	nbGames int;                                  // number of games stored
}

// global variables
// ----------------------------------------------------------------------------

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

// ShowHeaders summarizes the main information stored in the tags of all games
// in the given collection
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

// GameToLaTeX produces LaTeX code that uses package skak to show all games in a
// given collection
func (games *PgnCollection) GameToLaTeX () string {

	// start with the preamble of the document
	output := `\documentclass{article}
\usepackage[utf8]{inputenc}
\usepackage[english]{babel}
\usepackage{mathpazo}
\usepackage{nicefrac}
\usepackage{skak}
\def\hrulefill{\leavevmode\leaders\hrule height 10pt\hfill\kern\z@}
\begin{document}`

	// now, process each game in succession
	for _, game := range games.slice {

		output += fmt.Sprintf(`%v
\clearpage
`, game.getLaTeXbody ())
	}

	// and end the document
	output += `
\end{document}`
	
	// and return the string
	return output
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
