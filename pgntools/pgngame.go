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
	"errors" // for signaling errors
	"fmt"    // printing msgs
	"log"    // logging services

	"github.com/expr-lang/expr"
)

// typedefs
// ----------------------------------------------------------------------------

// A PGN move consist of a single ply. For each move the move number, color
// (with -1 representing black and +1 representing white) and actual move value
// (in algebraic form) is stored. Additionally, in case that the elapsed move
// time was present in the PGN file, it is also stored here.
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
// 1), (1, 0) and (0.5, 0.5). In addition, the pair (-1, -1) is considered for
// those games which are not properly ended
type PgnOutcome struct {
	scoreWhite, scoreBlack float32
}

// A game consists just of a map that stores information of all PGN tags, the
// sequence of moves and finally the outcome.
type PgnGame struct {
	tags    map[string]any
	moves   []PgnMove
	outcome PgnOutcome
}

// Functions
// ----------------------------------------------------------------------------
// Evaluate the given expression in the specified environment and return the
// result
func evaluateExpr(expression string, env map[string]any) (any, error) {

	// compile and run the given expression
	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		return nil, err
	}
	output, err := expr.Run(program, env)
	if err != nil {
		return nil, err
	}

	// and return the result without errors
	return output, nil
}

// Methods
// ----------------------------------------------------------------------------

// Return the number of the given PgnMove
func (move PgnMove) Number() int {
	return move.number
}

// Return the color of the given PgnMove
func (move PgnMove) Color() int {
	return move.color
}

// Return the actual move of the given PgnMove
func (move PgnMove) Move() string {
	return move.moveValue
}

// Return comments of the given PgnMove
func (move PgnMove) Comments() string {
	return move.comments
}

// Produces a string with the actual content of this move
func (move PgnMove) String() string {
	var output string

	// first, show the ply
	if move.color == 1 {
		output += fmt.Sprintf("%v. ", move.number)
	} else {
		output += fmt.Sprintf("%v. ... ", move.number)
	}

	output += fmt.Sprintf("%v ", move.moveValue)
	return output
}

// Produces a string with information of this outcome as a pair of
// floating-point numbers
func (outcome PgnOutcome) String() string {

	// In case this game was not properly ended, show an asterisk
	if outcome.scoreWhite == outcome.scoreBlack &&
		outcome.scoreWhite == -1 {
		return "*"
	}

	// Otherwise, show the result avoiding the usage of floating point numbers
	if outcome.scoreWhite == outcome.scoreBlack &&
		outcome.scoreWhite == 0.5 {
		return "1/2-1/2"
	}
	return fmt.Sprintf("%v-%v", outcome.scoreWhite, outcome.scoreBlack)
}

// return a string showing all moves in the specified interval in vertical mode,
// i.e. from move number from until move number to not included.
func (game *PgnGame) prettyMoves(from, to int) (output string) {

	// in case no moves were given just return the empty string
	if from == to {
		return
	}

	// get the slice of moves to show
	moves := game.moves[from:to]

	// add the first move. This is important because in case it is black to move,
	// an ellipsis should be shown first and, in case it is white's turn
	// everything will get rendered as desired
	output = fmt.Sprintf(" %v", moves[0])

	// process the rest of moves taking care to add a trailing newline after each
	// black's move
	idx := 1
	for idx < len(moves) {

		// first, in case the previous move was black's turn
		if moves[idx-1].Color() == -1 {

			// then add a trailing newline
			output += "\n"

			// and also show the number of the next move
			output += fmt.Sprintf(" %v. ", moves[idx].Number())
		}

		// Add the next move and proceed
		output += fmt.Sprintf("%v ", moves[idx].Move())

		// and proceed to the next move
		idx += 1
	}

	// and return the string computed so far
	return
}

// Return an environment for the evaluation of expressions
func (game *PgnGame) getEnv() (env map[string]any) {

	env = make(map[string]any)

	// Add all variables found in the tags of this game
	for variable, value := range game.Tags() {
		env[variable] = value
	}

	// In addition, create the variable "Moves" representing the number of moves
	// (not plies)
	if len(game.moves)%2 == 0 {
		env["Moves"] = len(game.moves) / 2
	} else {
		env["Moves"] = 1 + len(game.moves)/2
	}

	// and return the environment
	return
}

// Return the result of executing the given criteria as a string with
// information in this game and nil if no error happened.
func (game *PgnGame) getResult(criteria string) (string, error) {

	// execute the ith-criteria of this histogram
	env := game.getEnv()
	output, err := evaluateExpr(criteria, env)
	if err != nil {
		return "", err
	}

	// return the result casted as a string with success
	return fmt.Sprintf("%v", output), nil
}

// return true if the receiver must go before the other game and false otherwise
// according to the given sorting criteria. If the evaluation of any criteria
// produced an error it is returned and the boolean result is invalid
func (game PgnGame) lessGame(other PgnGame, criteria criteriaSorting) (bool, error) {

	// process all criteria given
	for _, icriteria := range criteria {

		// get the result of this criteria both in this game and the other
		iresult, ierr := game.getResult(icriteria.criteria)
		if ierr != nil {
			return false, ierr
		}
		jresult, jerr := other.getResult(icriteria.criteria)
		if jerr != nil {
			return false, jerr
		}

		// The result of an execution could be anything. However sorting is done
		// lexicographically on the given criteria and thus comparisons are done
		// as strings (note that "false" < "true"). Next in case one of the
		// values is either gt or lt than the other a comparison is performed.
		// Otherwise, the next sorting criteria should be visited
		if (iresult < jresult && icriteria.direction == increasing) ||
			(iresult > jresult && icriteria.direction == decreasing) {
			return true, nil
		}
		if (iresult > jresult && icriteria.direction == increasing) ||
			(iresult < jresult && icriteria.direction == decreasing) {
			return false, nil
		}
	}

	// At this point, both games have been proven to be strinctly equal
	// according to the given criteria
	return false, nil
}

// Return the tags of this game
func (game *PgnGame) Tags() (tags map[string]any) {
	return game.tags
}

// Return a list of the moves of this game as a slice of PgnMove
func (game *PgnGame) Moves() []PgnMove {
	return game.moves
}

// Return an instance of PgnOutcome with the result of this game
func (game *PgnGame) Outcome() PgnOutcome {
	return game.outcome
}

// Return whether the given expression is true or not for this specific game
func (game *PgnGame) Filter(expression string) (bool, error) {

	// First of all, create an environment for the evaluation of the given expression
	env := game.getEnv()

	// evaluate the given expression within the environment
	output, err := evaluateExpr(expression, env)
	if err != nil {
		return false, err
	}

	// Verify the result can be expressed as a boolean value
	result, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf(" The expression '%v' does not produced a boolean value!", expression)
	}

	// and return the result
	return result, nil
}

// Return the contents of this game in PGN format
func (game *PgnGame) GetPGN() (output string) {

	// First, show all tags followed by a blank line
	for variable, value := range game.tags {
		output += fmt.Sprintf("[%v \"%v\"]\n", variable, value)
	}
	output += "\n"

	// Next, write all moves of this game in a single line
	idx := 0
	for idx < len(game.moves) {

		// Write the move number and the white's move
		output += fmt.Sprintf("%v. %v ", game.moves[idx].number, game.moves[idx].moveValue)

		// and in case this move has an emt/ comments add them
		if game.moves[idx].emt > 0.0 {
			output += fmt.Sprintf("{[%%emt %v]} ", game.moves[idx].emt)
		}
		if game.moves[idx].comments != "" {
			output += fmt.Sprintf("{ %v } ", game.moves[idx].comments)
		}
		idx += 1

		// in case there is a move for black, then add it immediately after
		if idx < len(game.moves) {
			output += fmt.Sprintf("%v ", game.moves[idx].moveValue)

			// and in case this move has any emt/comments add them
			if game.moves[idx].emt > 0.0 {
				output += fmt.Sprintf("{[%%emt %v]} ", game.moves[idx].emt)
			}
			if game.moves[idx].comments != "" {
				output += fmt.Sprintf("{ %v } ", game.moves[idx].comments)
			}
			idx += 1
		}
	}

	// Next, show the result which is used as a token of end of game
	output += fmt.Sprintf("%v", game.Outcome())

	// and add a blank line
	output += "\n\n"

	// and return the game in PGN format
	return
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

// Return the value of a specific tag. In case the tag is not found in this
// game, an error is return along with any data.
//
// It is intended to be used in LaTeX templates
func (game *PgnGame) GetTagValue(name string) (value string, err error) {

	if value, ok := game.tags[name]; ok {
		return fmt.Sprintf("%v", value), nil
	}

	// when getting here, the required tag has not been found
	return "", errors.New(fmt.Sprintf("tag '%s' not found!", name))
}

// A field is either a tag of the receiver game, or a value that can be
// extracted from it (such as "Moves" or "Result")
//
// This function is to be used in LaTeX templates
func (game *PgnGame) GetField(field string) string {

	// -- Moves
	if field == "Moves" {

		// In case the number of moves were requested, then return the number of
		// moves stored in this game. Some PGN files might contain a tag named
		// "PlyCount", for example, but it is unnecessary to rely on its
		// existence.
		return fmt.Sprintf("%d", len(game.moves))
	}

	// -- Moves
	if field == "Result" {

		if game.outcome.scoreWhite == 0.5 {
			return "½-½"
		} else if game.outcome.scoreWhite == 1 {
			return "1-0"
		} else if game.outcome.scoreBlack == 1 {
			return "0-1"
		} else if game.outcome.scoreWhite == -1 {
			return "*"
		} else {
			log.Fatalln(" Unknown result found!")
		}
	}

	// -- tags

	// after trying special fields, then tags defined in this game are
	// tried. In case they do not exist, an error is automatically raisedx
	value, err := game.GetTagValue(field)
	if err != nil {
		log.Fatalf(" Uknown field '%v'\n", field)
	}
	return fmt.Sprintf("%v", value)
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
		result = append(result, game.GetField(field_str))
	}

	// return the slice of strings computed so far
	return
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
