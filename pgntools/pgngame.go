/* 
  pgngame.go
  Description: Simple tools to handle a single game in PGN format
  ----------------------------------------------------------------------------- 

  Started on  <Sat May  9 16:59:21 2015 Carlos Linares Lopez>
  Last update <sábado, 09 mayo 2015 17:31:18 Carlos Linares Lopez (clinares)>
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
	"errors"		// for signaling errors
	"fmt"			// printing msgs	
	"log"			// logging services
	"strconv"		// to convert from strings to other types
)

// typedefs
// ----------------------------------------------------------------------------
type PgnTag struct {

	name, value string;	              // name and value of a single tag
}

type PgnMove struct {

	moveNumber int;                                  // current move number
	color int;                                  // color: 1=white; -1=black
	moveValue string;                           // move value in PGN format
	emt float32;                                       // elapsed move time
	comments string; 	  // comments - in case there are various, each
				      // one is added after a newline character
}

type PgnOutcome struct {

	scoreWhite, scoreBlack float32;                 // score of each player
}

type PgnGame struct {

	tags map[string]string;      // A game consists of a collection of tags
	                                               // index by the tag name
	moves []PgnMove;                      // sequence of moves of this game
	outcome PgnOutcome;                                    // final outcome
}

// Methods
// ----------------------------------------------------------------------------

// the following methods overwrite the string output method
func (tag PgnTag) String () string {
	return fmt.Sprintf ("%v: %v", tag.name, tag.value)
}

func (move PgnMove) String () string {
	if move.color == 1 {
		return fmt.Sprintf ("%v. %v", move.moveNumber, move.moveValue)
	}
	return fmt.Sprintf (" %v ", move.moveValue)
}

func (outcome PgnOutcome) String () string {
	return fmt.Sprintf ("%v - %v", outcome.scoreWhite, outcome.scoreBlack)
}

// the following service just prints all the sequence of moves in the given game
func (game *PgnGame) String () string {
	output := ""
	for _, move := range game.moves {
		output += fmt.Sprintf ("%v", move)
	}
	return output
}

// the following are getters over the attributes of a PgnGame
func (game *PgnGame) GetTags () map[string]string {
	return game.tags
}

func (game *PgnGame) GetMoves () []PgnMove {
	return game.moves
}

func (game *PgnGame) GetOutcome () PgnOutcome {
	return game.outcome
}

// GetTagValue return the value of a specific tag and nil if it exists or any
// value and err in case it does not exist
func (game *PgnGame) GetTagValue (name string) (value string, err error) {

	if value, ok := game.tags[name]; ok {
		return value, nil
	}
	
	// when getting here, the required tag has not been found
	return "", errors.New ("tag not found!")
}

// ShowHeader summarizes the main information stored in the tags of a specific
// game
func (game *PgnGame) ShowHeader () string {

	// first, verify that all necessary tags are available
	dbGameNo, err := game.GetTagValue ("FICSGamesDBGameNo")
	if err != nil {
		log.Fatalf ("FICSGamesDBGameNo not found!")
	}
	
	date, err := game.GetTagValue ("Date")
	if err != nil {
		log.Fatalf ("Date not found!")
	}
	
	time, err := game.GetTagValue ("Time")
	if err != nil {
		log.Fatalf ("Time not found!")
	}
	
	white, err := game.GetTagValue ("White")
	if err != nil {
		log.Fatalf ("White not found!")
	}
	
	whiteELO, err := game.GetTagValue ("WhiteElo")
	if err != nil {
		log.Fatalf ("WhiteElo not found!")
	}
	
	black, err := game.GetTagValue ("Black")
	if err != nil {
		log.Fatalf ("Black not found!")
	}
	
	blackELO, err := game.GetTagValue ("BlackElo")
	if err != nil {
		log.Fatalf ("BlackElo not found!")
	}
	
	ECO, err := game.GetTagValue ("ECO")
	if err != nil {
		log.Fatalf ("ECO not found!")
	}
	
	timeControl, err := game.GetTagValue ("TimeControl")
	if err != nil {
		log.Fatalf ("TimeControl not found!")
	}

	plyCount, err := game.GetTagValue ("PlyCount")
	if err != nil {
		log.Fatalf ("PlyCount not found!")
	}
	moves, err := strconv.Atoi (plyCount)
	if 2*(moves/2) < moves {
		moves = moves/2 + 1
	} else {
		moves /=2
	}

	var scoreWhite, scoreBlack string;
	outcome := game.GetOutcome ()
	if outcome.scoreWhite == 0.5 {
		scoreWhite, scoreBlack = "½", "½"
	} else if outcome.scoreWhite == 1 {
		scoreWhite, scoreBlack = "1", "0"
	} else {
		scoreWhite, scoreBlack = "0", "1"
	}

	return fmt.Sprintf (" | %10v | %v %v | %-18v (%4v) | %-18v (%4v) | %v | %v | %5v |    %v-%-v |", dbGameNo, date, time, white, whiteELO, black, blackELO, ECO, timeControl, moves, scoreWhite, scoreBlack)
}

// getLaTeXbody computes the main part of the LaTeX document that shows
// information of a specific game
func (game *PgnGame) getLaTeXbody () string {

	// first, verify that all necessary tags are available
	event, err := game.GetTagValue ("Event")
	if err != nil {
		log.Fatalf ("Event not found!")
	}
	
	date, err := game.GetTagValue ("Date")
	if err != nil {
		log.Fatalf ("Date not found!")
	}
	
	white, err := game.GetTagValue ("White")
	if err != nil {
		log.Fatalf ("White not found!")
	}
	
	whiteELO, err := game.GetTagValue ("WhiteElo")
	if err != nil {
		log.Fatalf ("WhiteElo not found!")
	}
	
	black, err := game.GetTagValue ("Black")
	if err != nil {
		log.Fatalf ("Black not found!")
	}
	
	blackELO, err := game.GetTagValue ("BlackElo")
	if err != nil {
		log.Fatalf ("BlackElo not found!")
	}
	
	ECO, err := game.GetTagValue ("ECO")
	if err != nil {
		log.Fatalf ("ECO not found!")
	}
	
	timeControl, err := game.GetTagValue ("TimeControl")
	if err != nil {
		log.Fatalf ("TimeControl not found!")
	}

	var scoreWhite, scoreBlack string;
	outcome := game.GetOutcome ()
	if outcome.scoreWhite == 0.5 {
		scoreWhite, scoreBlack = `\textonehalf`, `\textonehalf`
	} else if outcome.scoreWhite == 1 {
		scoreWhite, scoreBlack = "1", "0"
	} else {
		scoreWhite, scoreBlack = "0", "1"
	}
	
	// now, initialize the output with the main contents of the LaTeX body
	output := fmt.Sprintf (`\begin{center}
  {\Large %v (%v)}  
\end{center}

\hrule
\noindent
\WhiteKnightOnWhite %v (%v) \hfill %v\\
\BlackKnightOnWhite %v (%v) \hfill %v
\hrule

\vspace{0.5cm}

\newgame

\mainline{%v}\hfill{\textbf{%v}-\textbf{%v}}

\begin{center}
  \showboard
\end{center}`, event, timeControl, white, whiteELO, date, black, blackELO, ECO, game, scoreWhite, scoreBlack)

	// and return the string computed so far
	return output
}

// GameToLaTeX produces LaTeX code that uses package skak to show the given game
func (game *PgnGame) GameToLaTeX () string {

	// justsubstitute values over a standard template
	output := fmt.Sprintf (`\documentclass{article}
\usepackage[utf8]{inputenc}
\usepackage[english]{babel}
\usepackage{mathpazo}
\usepackage{nicefrac}
\usepackage{skak}
\def\hrulefill{\leavevmode\leaders\hrule height 10pt\hfill\kern\z@}
\begin{document}

%v

\end{document}`, game.getLaTeXbody ())

	// and return the string
	return output
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
