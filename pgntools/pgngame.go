/* 
  pgngame.go
  Description: Simple tools to handle a single game in PGN format
  ----------------------------------------------------------------------------- 

  Started on  <Sat May  9 16:59:21 2015 Carlos Linares Lopez>
  Last update <domingo, 10 mayo 2015 15:49:45 Carlos Linares Lopez (clinares)>
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
	"regexp"                // pgn files are parsed with a regexp
	"strconv"		// to convert from strings to other types

	// import a user package to manage paths
	"bitbucket.org/clinares/pgnparser/fstools"
)

// global variables
// ----------------------------------------------------------------------------

// The following regexp matches any placeholder appearing in a LaTeX
// file. Placeholder have the form '%name' where 'name' consists of any
// combination of alpha and numeric characters.
var reGroupPlaceholder = regexp.MustCompile (`%[\w\d]+`)


// typedefs
// ----------------------------------------------------------------------------

// A PGN tag consists of a pair <name, value>
type PgnTag struct {

	name, value string;
}

// A PGN move consist of a single ply. For each move the move number, color and
// actual move value (in algebraic form) is stored. Additionally, in case that
// the elapsed move time was present in the PGN file, it is also stored
// here.
//
// Finally, any combination of moves after the move are combined into the
// same field (comments). In case various comments were given they are then
// separated by '\n'.
type PgnMove struct {

	moveNumber int;
	color int;
	moveValue string;
	emt float32;
	comments string;
}

// The outcome of a chess game consists of the score obtained by every player as
// two float32 numbers such that their sum equals 1. Plausible outcomes are (0,
// 1), (1, 0) and (0.5, 0.5)
type PgnOutcome struct {

	scoreWhite, scoreBlack float32;
}

// A game consists just of a map that stores information of all PGN tags, the
// sequence of moves and finally the outcome.
type PgnGame struct {

	tags map[string]string;
	moves []PgnMove;
	outcome PgnOutcome;
}

// Methods
// ----------------------------------------------------------------------------

// Produces a string with information of this tag
func (tag PgnTag) String () string {
	return fmt.Sprintf ("%v: %v", tag.name, tag.value)
}

// Produces a string with information of this move.
//
// The main rationale behind this method is that when applied successively to
// the moves of a particular game (given as a instance of PgnGame), the full
// sequence has to be shown in algebraic form. This means, that only white moves
// are preceded by the move number, while black's moves are just inserted
// without the move number.
func (move PgnMove) String () string {
	if move.color == 1 {
		return fmt.Sprintf ("%v. %v", move.moveNumber, move.moveValue)
	}
	return fmt.Sprintf (" %v ", move.moveValue)
}

// Produces a string with information of this outcome as a pair of
// floating-point numbers
// ----------------------------------------------------------------------------
func (outcome PgnOutcome) String () string {
	return fmt.Sprintf ("%v - %v", outcome.scoreWhite, outcome.scoreBlack)
}

// Produces a string with the list of moves of this game.
//
// This method successively invokes the String () service provided by PgnMove
// over every move of this particular game. As a result, a full transcription of
// the game is returned in the output string
func (game *PgnGame) String () string {
	output := ""
	for _, move := range game.moves {
		output += fmt.Sprintf ("%v", move)
	}
	return output
}

// Return the tags of this game as a map from tag names to tag values. Although
// tag values are given between double quotes, these are not shown.
func (game *PgnGame) GetTags () map[string]string {
	return game.tags
}

// Return a list of the moves of this game as a slice of PgnMove
func (game *PgnGame) GetMoves () []PgnMove {
	return game.moves
}

// Return an instance of PgnOutcome with the result of this game
func (game *PgnGame) GetOutcome () PgnOutcome {
	return game.outcome
}

// Return the value of a specific tag and nil if it exists or any value and err
// in case it does not exist
func (game *PgnGame) GetTagValue (name string) (value string, err error) {

	if value, ok := game.tags[name]; ok {
		return value, nil
	}
	
	// when getting here, the required tag has not been found
	return "", errors.New ("tag not found!")
}

// Return a string with a summary of the main information stored in this game
//
// In case any required data is not found, a fatal error is raised
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

// returns the result of replacing all placeholders in template with their
// value. Placeholders are identified with the string '%<name>'. All tag names
// specified in this game are acknowledged. Additionally, '%moves' is
// substituted by the list of moves func (game *PgnGame) replacePlaceholders
func (game *PgnGame) replacePlaceholders (template string) string {

	return reGroupPlaceholder.ReplaceAllStringFunc(template,
		func (name string) string {

			// get rid of the leading '%' character
			placeholder := name[1:]
			
			// most placeholders are just tag names. However,
			// 'moves' is also acknowledged
			if placeholder == "moves" {
				return fmt.Sprintf ("%v", game)
			}

			// otherwise, return the value of this tag
			return game.tags [placeholder]
		})
}

// Produces LaTeX code using the specified template with information of this
// game. The string acknowledges various placeholders which have the format
// '%<name>'. All tag names specified in this game are
// acknowledged. Additionally, '%moves' is substituted by the list of moves
func (game *PgnGame) GameToLaTeXFromString (template string) string {

	// just substitute values over the given template and return the result
	return game.replacePlaceholders (template)
}

// Produces LaTeX code using the template stored in the specified file with
// information of this game. The string acknowledges various placeholders which
// have the format '%<name>'. All tag names specified in this game are
// acknowledged. Additionally, '%moves' is substituted by the list of moves
func (game *PgnGame) GameToLaTeXFromFile (templateFile string) string {

	// Open and read the given file and retrieve its contents
	contents := fstools.Read (templateFile, -1)
	template := string (contents[:len (contents)])

	// and now, just return the results of parsing these contents
	return game.GameToLaTeXFromString (template)
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
