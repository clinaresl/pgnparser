/*
  pgngame.go
  Description: Simple tools to handle a single game in PGN format
  -----------------------------------------------------------------------------

  Started on  <Sat May  9 16:59:21 2015 Carlos Linares Lopez>
  Last update <sábado, 07 mayo 2016 16:44:27 Carlos Linares Lopez (clinares)>
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
	"errors"  // for signaling errors
	"fmt"     // printing msgs
	"log"     // logging services
	"strconv" // to conver int to string
)

// typedefs
// ----------------------------------------------------------------------------

// tags tables and symbol tables store any data that supports the following
// boolean comparisons: <, = and >
type dataInterface interface {
	Less(right dataInterface) bool
	Equal(right dataInterface) bool
	Greater(right dataInterface) bool
}

// a number of types that qualify for the dataInterface are integers and strings
type constInteger int32
type constString string

// A PGN move consist of a single ply. For each move the move number, color and
// actual move value (in algebraic form) is stored. Additionally, in case that
// the elapsed move time was present in the PGN file, it is also stored
// here.
//
// Finally, any combination of moves after the move are combined into the
// same field (comments). In case various comments were given they are then
// separated by '\n'.
type PgnMove struct {
	number    int
	color     int
	moveValue string
	emt       float32
	comments  string
}

// The outcome of a chess game consists of the score obtained by every player as
// two float32 numbers such that their sum equals 1. Plausible outcomes are (0,
// 1), (1, 0) and (0.5, 0.5)
type PgnOutcome struct {
	scoreWhite, scoreBlack float32
}

// A game consists just of a map that stores information of all PGN tags, the
// sequence of moves and finally the outcome.
type PgnGame struct {
	tags    map[string]dataInterface
	moves   []PgnMove
	outcome PgnOutcome
}

// Methods
// ----------------------------------------------------------------------------

// Return true if both the receiver and the argument are integers and the
// receiver is less than the argument
func (constant constInteger) Less(right dataInterface) bool {

	var value constInteger
	var ok bool

	// verify both types are compatible
	value, ok = right.(constInteger)
	if !ok {
		log.Fatal(" Type mismatch in pgngame.go::Less (constInteger)")
	}

	return int32(constant) < int32(value)
}

// Return true if both the receiver and the argument are integers holding the
// same value
func (constant constInteger) Equal(right dataInterface) bool {

	var value constInteger
	var ok bool

	// verify both types are compatible
	value, ok = right.(constInteger)
	if !ok {
		log.Fatal(" Type mismatch in pgngame.go::Equal (constInteger)")
	}

	return int32(constant) == int32(value)
}

// Return true if both the receiver and the argument are integers and the
// receiver is greater than the argument
func (constant constInteger) Greater(right dataInterface) bool {

	var value constInteger
	var ok bool

	// verify both types are compatible
	value, ok = right.(constInteger)
	if !ok {
		log.Fatal(" Type mismatch in pgngame.go::Greater (constInteger)")
	}

	return int32(constant) > int32(value)
}

// Return true if both the receiver and the argument are strings and the
// receiver is less than the argument
func (constant constString) Less(right dataInterface) bool {

	var value constString
	var ok bool

	// verify both types are compatible
	value, ok = right.(constString)
	if !ok {
		log.Fatal(" Type mismatch in pgngame.go::Less (constString)")
	}

	return string(constant) < string(value)
}

// Return true if both the receiver and the argument are integers holding the
// same value
func (constant constString) Equal(right dataInterface) bool {

	var value constString
	var ok bool

	// verify both types are compatible
	value, ok = right.(constString)
	if !ok {
		log.Fatal(" Type mismatch in pgngame.go::Equal (constString)")
	}

	return string(constant) == string(value)
}

// Return true if both the receiver and the argument are integers and the
// receiver is greater than the argument
func (constant constString) Greater(right dataInterface) bool {

	var value constString
	var ok bool

	// verify both types are compatible
	value, ok = right.(constString)
	if !ok {
		log.Fatal(" Type mismatch in pgngame.go::Greater (constString)")
	}

	return string(constant) > string(value)
}

// Produces a string with the actual content of this move
func (move PgnMove) String() string {
	var output string

	// first, show the ply
	if move.color == 1 {
		output += fmt.Sprintf("%v.", move.number)
	} else {
		output += fmt.Sprintf("%v. ... ", move.number)
	}

	output += fmt.Sprintf("%v ", move.moveValue)
	return output
}

// Produces a string with information of this outcome as a pair of
// floating-point numbers
func (outcome PgnOutcome) String() string {
	return fmt.Sprintf("%v - %v", outcome.scoreWhite, outcome.scoreBlack)
}

// Return the tags of this game as a map from tag names to tag values. Although
// tag values are given between double quotes, these are not shown.
func (game *PgnGame) GetTags() map[string]dataInterface {
	return game.tags
}

// Return a list of the moves of this game as a slice of PgnMove
func (game *PgnGame) GetMoves() []PgnMove {
	return game.moves
}

// Return an instance of PgnOutcome with the result of this game
func (game *PgnGame) GetOutcome() PgnOutcome {
	return game.outcome
}

// Parse all moves of this game. Show the board between showboard consecutive
// plies
func (game *PgnGame) ParseMoves(plies int) {

	nrplies := 0
	board := NewPgnBoard()

	for _, move := range game.moves {
		board.UpdateBoard(move, plies > 0)

		// show the board on the console?
		nrplies += 1 // incremente the number of plies processed
		if plies > 0 && nrplies%plies == 0 {
			fmt.Printf("%v\n\n", board)
		}
	}

	// finally, if plies is positive, show the end position unless it was
	// incidentally shown within the previous loops
	if plies > 0 && nrplies%plies != 0 {
		fmt.Printf("%v\n\n", board)
	}
}

// Templates
//
// All the following methods are used to handle templates both for generating
// ascii and LaTeX output
// ----------------------------------------------------------------------------

// getColorPrefix is a helper function that returns the prefix of the color of
// the receiving move. In case it is white's turn then '.' is returned;
// otherwise '...' is returned
func (move PgnMove) getColorPrefix() (prefix string) {
	if move.color == 1 {
		prefix = "."
	} else if move.color == -1 {
		prefix = "..."
	} else {
		log.Fatalf(fmt.Sprintf(" Unknown color in move '%v'", move))
	}
	return
}

// Produces a LaTeX string with a plain list of the moves of this game. It is
// intended to be used in LaTeX templates
func (game *PgnGame) GetLaTeXMoves() (output string) {

	// Initialization
	output = `\mainline{`

	// Iterate over all moves
	for _, move := range game.moves {

		// in case it is white's turn then precede this move by the move
		// counter and the prefix of the color
		if move.color == 1 {
			output += fmt.Sprintf("%v. %v", move.number, move)
		} else {

			// otherwise, just show the actual move
			output += fmt.Sprintf(" %v", move)
		}
	}

	// close the mainline
	output += `}`

	// and return the string
	return
}

// Produces a LaTeX string with the list of moves of this game along with the
// different annotations.
//
// This method successively processes the moves in this PgnGame until a comment
// is found. If a "literal" command is found, it is just added to the
// output. Other "special" comments are:
//
// 1. %emt which show the elapsed move time
//
// 2. %show which generates a LaTeX command for showing the current board
//
// It is intended to be used in LaTeX templates
func (game *PgnGame) GetLaTeXMovesWithComments() (output string) {

	// the variable newMainLine is used to determine whether the next move
	// should start with a LaTeX command \mainline. Obviously, this is
	// initially true
	newMainLine := true

	// Iterate over all moves
	for _, move := range game.moves {

		// before printing this move, check if a new mainline has to be
		// started (e.g., because the previous move ended with a
		// comment
		if newMainLine {
			output += `\mainline{ `
		}

		// now in case either we are starting a new mainline or it is
		// white's move, then show all the details of the move including
		// counter and color prefix
		if newMainLine || move.color == 1 {

			// now, show the actual move with all details
			output += fmt.Sprintf("%v%v %v ", move.number, move.getColorPrefix(), move.moveValue)
		} else {

			// otherwise, just show the actual move
			output += fmt.Sprintf("%v ", move.moveValue)
		}

		// if this move contains either a comment or the emt
		if move.emt != -1 || move.comments != "" {

			output += "} "

			// now, in case emt is present, show it
			if move.emt != -1 {
				output += fmt.Sprintf(`({\it %v}) `, move.emt)
			}

			// if a comment is present, show it as well
			if move.comments != "" {

				output += fmt.Sprintf("%v ", move.comments)
			}
		}

		// and check whether a new mainline has to be started in the
		// next iteration
		newMainLine = (move.emt != -1 || move.comments != "")
	}

	// and return the string computed so far
	return
}

// Return the value of a specific tag and nil if it exists or any value and err
// in case it does not exist. It is intended to be used in LaTeX templates
func (game *PgnGame) GetTagValue(name string) (value dataInterface, err error) {

	if value, ok := game.tags[name]; ok {
		return value, nil
	}

	// when getting here, the required tag has not been found
	return constString(""), errors.New(fmt.Sprintf("tag '%s' not found!", name))
}

// getAndCheckTag is a helper function whose purpose is just to retrieve the
// value of a given tag. In case an error happened (most likely because it does
// not exist) then a fatal error is issued and execution is stopped
func (game *PgnGame) getAndCheckTag(tagname string) dataInterface {

	value, err := game.GetTagValue(tagname)

	// in an error was found, then issue a fatal error
	if err != nil {
		log.Fatalf(fmt.Sprintf("'%v' not found in game\n '%v'!", tagname, game))
	}

	// otherwise, return the value of this tagname
	return value
}

// A field is either a tag of the receiver game or a value computed from the
// tags. Fields which are computed from tags are:
//
//	Moves: number of moves (two plies each)
//	Result: consists of a utf-8 string which contains the final result of the
//	game
//
// This method is used to compute arbitrary fields to be shown in ascii tables
func (game *PgnGame) getField(field string) string {

	// -- Moves
	if field == "Moves" {

		// get the ply count of this game
		plyCount := game.getAndCheckTag("PlyCount")

		// now, compute the number of moves from the number of plies. If
		// the number of plies is even, then the number of moves is half
		// the number of plies, otherwise, add 1
		moves, ok := plyCount.(constInteger)
		if !ok {
			log.Fatalf(fmt.Sprintf(" It was not possible to convert the PlyCount ('%v') into an integer", plyCount))
		}
		if 2*(moves/2) < moves {
			moves = moves/2 + 1
		} else {
			moves /= 2
		}

		// and return the number of moves
		return strconv.Itoa(int(moves))
	}

	// -- Moves
	if field == "Result" {

		var scoreWhite, scoreBlack string
		if game.outcome.scoreWhite == 0.5 {
			scoreWhite, scoreBlack = "½", "½"
		} else if game.outcome.scoreWhite == 1 {
			scoreWhite, scoreBlack = "1", "0"
		} else {
			scoreWhite, scoreBlack = "0", "1"
		}
		return scoreWhite + "-" + scoreBlack
	}

	// -- tags

	// after trying special fields, then tags defined in this game are
	// tried. In case they do not exist, an error is automatically raisedx
	return fmt.Sprintf("%v", game.getAndCheckTag(field))
}

// Return a slice of strings with the values of all given fields. This method is
// used to compute the fields of a game to be shown on an ascii table.
//
// Note that the returned slice contain instances of any. This is necessary
// because []string is not a subtype of []any, i.e., slices are type-invariant
// in Go (and other langs).
func (game *PgnGame) getFields(fields []any) (result []any) {

	// iterate over all fields
	for _, field := range fields {

		// compute the value of the next field and add it to the slice
		// to return
		field_str, ok := field.(string)
		if !ok {
			log.Fatalf(fmt.Sprintf(" It was not possible to convert the field '%v' into a string", field))
		}
		result = append(result, game.getField(field_str))
	}

	// return the slice of strings computed so far
	return
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
