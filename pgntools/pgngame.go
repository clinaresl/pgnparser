/* 
  pgngame.go
  Description: Simple tools to handle a single game in PGN format
  ----------------------------------------------------------------------------- 

  Started on  <Sat May  9 16:59:21 2015 Carlos Linares Lopez>
  Last update <domingo, 10 mayo 2015 01:09:44 Carlos Linares Lopez (clinares)>
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

// the following regexp matches any placeholder appearing in a LaTeX
// file
var reGroupPlaceholder = regexp.MustCompile (`%[\w\d]+`)


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

// String
//
// produces a string with information of this tag
// ----------------------------------------------------------------------------
func (tag PgnTag) String () string {
	return fmt.Sprintf ("%v: %v", tag.name, tag.value)
}

// String
// 
// produces a string with information of this move
// ----------------------------------------------------------------------------
func (move PgnMove) String () string {
	if move.color == 1 {
		return fmt.Sprintf ("%v. %v", move.moveNumber, move.moveValue)
	}
	return fmt.Sprintf (" %v ", move.moveValue)
}

// String
//
// produces a string with information of this outcome
// ----------------------------------------------------------------------------
func (outcome PgnOutcome) String () string {
	return fmt.Sprintf ("%v - %v", outcome.scoreWhite, outcome.scoreBlack)
}

// String
//
// produces a string with the list of moves of this game
// ----------------------------------------------------------------------------
func (game *PgnGame) String () string {
	output := ""
	for _, move := range game.moves {
		output += fmt.Sprintf ("%v", move)
	}
	return output
}

// GetTags
//
// Return the tags of this game
// ----------------------------------------------------------------------------
func (game *PgnGame) GetTags () map[string]string {
	return game.tags
}

// GetMoves
//
// Return a list of the moves of this game
// ----------------------------------------------------------------------------
func (game *PgnGame) GetMoves () []PgnMove {
	return game.moves
}

// GetOutcome
//
// Return an instance of PgnOutcome with the result of this game
// ----------------------------------------------------------------------------
func (game *PgnGame) GetOutcome () PgnOutcome {
	return game.outcome
}

// GetTagValue
// 
// return the value of a specific tag and nil if it exists or any value and err
// in case it does not exist
// ----------------------------------------------------------------------------
func (game *PgnGame) GetTagValue (name string) (value string, err error) {

	if value, ok := game.tags[name]; ok {
		return value, nil
	}
	
	// when getting here, the required tag has not been found
	return "", errors.New ("tag not found!")
}

// ShowHeader
// 
// return a string with a summary of the main information stored in this game
// ----------------------------------------------------------------------------
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

// replacePlaceholders
//
// returns the result of replacing all placeholders in template with their
// value. Placeholders are identified with the string '%<name>'. All tag names
// specified in this game are acknowledged. Additionally, '%moves' is
// substituted by the list of moves
// ----------------------------------------------------------------------------
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

// GameToLaTeXFromString
//
// produces LaTeX code using the specified template with information of this
// game. The string acknowledges various placeholders which have the format
// '%<name>'. All tag names specified in this game are
// acknowledged. Additionally, '%moves' is substituted by the list of moves
// ----------------------------------------------------------------------------
func (game *PgnGame) GameToLaTeXFromString (template string) string {

	// just substitute values over the given template and return the result
	return game.replacePlaceholders (template)
}

// GameToLaTeXFromFile
//
// produces LaTeX code using the template stored in the specified file with
// information of this game. The string acknowledges various placeholders which
// have the format '%<name>'. All tag names specified in this game are
// acknowledged. Additionally, '%moves' is substituted by the list of moves
// ----------------------------------------------------------------------------
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
