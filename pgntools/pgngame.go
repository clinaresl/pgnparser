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
	// for signaling errors
	"fmt" // printing msgs
	"io"
	"log" // logging services
	"sort"
	"strings"

	"github.com/expr-lang/expr"
)

// typedefs
// ----------------------------------------------------------------------------

// A PGN move consist of a single ply. For each move the move number, color
// (with -1 representing black and +1 representing white) and actual move value
// (both in short and long algebraic notation) is stored. Additionally, in case
// that the elapsed move time was present in the PGN file, it is also stored
// here.
//
// Finally, any combination of moves after the move are combined into the
// same field (comments). In case various comments were given they are then
// separated by '\n'.
type PgnMove struct {
	number         int
	color          int
	shortAlgebraic string
	longAlgebraic
	emt      float32
	comments string
}

// A move in the long algebraic notation consists of a explicity description of
// the starting and end positions of the move
type longAlgebraic struct {
	from, to string
}

// The outcome of a chess game consists of the score obtained by every player as
// two float32 numbers such that their sum equals 1. Plausible outcomes are (0,
// 1), (1, 0) and (0.5, 0.5). In addition, the pair (-1, -1) is considered for
// those games which are not properly ended
type PgnOutcome struct {
	scoreWhite, scoreBlack float32
}

// A game consists just of a map that stores information of all PGN tags, the
// sequence of moves and successive boards and the outcome. For various purposes
// it contains also an id which is an integer index and is used to uniquely
// refer to each game.
type PgnGame struct {
	tags    map[string]any
	moves   []PgnMove
	boards  []PgnBoard
	outcome PgnOutcome
	id      int
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

// Return true if and only if the FEN piece placement of the first string
// matches the FEN piece placement of the second, and false otherwise. Both
// strings are supposed to contain only the piece placement of the FEN code and
// not the entire FEN code
func matchFENPiecePlacement(expr, code string) bool {

	// The piece placement is the same if and only if they are exactly equal and
	// false otherwise
	return expr == code
}

// Return true if and only if the FEN active color of the first string matches
// the FEN active color of the second, and false otherwise. Both strings are
// supposed to contain only the active color of the FEN code and not the
// entire FEN code
func matchFENActiveColor(expr, code string) bool {

	// The active color is the same if and only if they are exactly equal and
	// false otherwise
	return expr == code
}

// Return true if and only if the FEN castling rights of the first string
// matches the FEN castling rights of the second, and false otherwise. Both
// strings are supposed to contain only the castling rights of the FEN code and
// not the entire FEN code
func matchFENCastlingRights(expr, code string) bool {

	// The castling rights are the same if and only if they contain exactly the
	// same characters, even if they are in different order. Thus, sort each
	// string
	sexpr, scode := strings.Split(expr, ""), strings.Split(code, "")
	sort.Strings(sexpr)
	sort.Strings(scode)

	// and check now if they are the same
	return strings.Join(sexpr, "") == strings.Join(scode, "")
}

// Return true if and only if the FEN en passant targets of the first string
// matches the FEN en passant targets of the second, and false otherwise. Both
// strings are supposed to contain only the en passant targets of the FEN code
// and not the entire FEN code
func matchFENEnPassantTargets(expr, code string) bool {

	// The en passant target consists of a square in short algebraic notation
	// and thus, being equal they should be given exactly in the same way
	return expr == code
}

// Return true if and only if the FEN halfmove clock of the first string matches
// the FEN halfmove clock of the second, and false otherwise. Both strings are
// supposed to contain only the halfmove clock of the FEN code and not the
// entire FEN code
func matchFENHalfMoveClock(expr, code string) bool {

	// The halfmove clock consists of an integer, so they both should be the
	// same
	return expr == code
}

// Return true if and only if the FEN fullmove number of the first string
// matches the FEN fullmove number of the second, and false otherwise. Both
// strings are supposed to contain only the fullmove number of the FEN code and
// not the entire FEN code
func matchFENFullMoveNumber(expr, code string) bool {

	// The fullmove number consists of an integer, so they both should be the
	// same
	return expr == code
}

// Return true if and only if the first fen code matches the second. Matching
// means that they are actually the same even if they are written in different
// ways
func matchFEN(expr, code string) bool {

	// split both fen codes into their fields. Since they are assumed to be
	// correct, it just suffices splitting with the blank
	exprFields := strings.Split(expr, " ")
	codeFields := strings.Split(code, " ")

	// Piece placement
	if !matchFENPiecePlacement(exprFields[0], codeFields[0]) {
		return false
	}

	// Active Color
	if !matchFENActiveColor(exprFields[1], codeFields[1]) {
		return false
	}

	// Castling rights
	if !matchFENCastlingRights(exprFields[2], codeFields[2]) {
		return false
	}

	// En passant targets
	if !matchFENEnPassantTargets(exprFields[3], codeFields[3]) {
		return false
	}

	// Half move clock
	if !matchFENHalfMoveClock(exprFields[4], codeFields[4]) {
		return false
	}

	// Fullmove number
	if !matchFENFullMoveNumber(exprFields[4], codeFields[4]) {
		return false
	}

	// at this point, they are proven to be equal
	return true
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

// Return the actual move in short algebraic notation
func (move PgnMove) Move() string {
	return move.shortAlgebraic
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

	output += fmt.Sprintf("%v ", move.shortAlgebraic)
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

// Return true if and only if a board in this game contains a position with the
// given fen code
func (game *PgnGame) checkFEN(fencode string) bool {

	// Examine all positions in this game
	for _, iboard := range game.boards {

		// if this board has the given fen code immediately return true
		if matchFEN(fencode, iboard.fen) {
			return true
		}
	}

	// At this point, no position in this game has the given fen fencode
	return false
}

// return a string showing all moves in the specified interval in vertical mode,
// i.e. from move number 'from' until move number 'to' not included.
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

	// And also, add all the available functions
	env["FEN"] = func(fen string) bool {
		return game.checkFEN(fen)
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

// Return a list of the boards of this game as a slice of PgnBoards
func (game *PgnGame) Boards() []PgnBoard {
	return game.boards
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
		output += fmt.Sprintf("%v. %v ", game.moves[idx].number, game.moves[idx].shortAlgebraic)

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
			output += fmt.Sprintf("%v ", game.moves[idx].shortAlgebraic)

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

// The following function returns the input string after substituting all the
// special LaTeX characters so that they can be correctly processed
func substituteLaTeX(input string) (output string) {

	// Just substitute the special LaTeX characters one after the other. Note
	// that because most replacements start with a backslash and this is,
	// indeed, a special LaTeX character, backslash substitution takes
	// precedence
	output = strings.Replace(input, `\`, `\textbackslash `, -1)
	output = strings.Replace(output, "#", `\#`, -1)
	output = strings.Replace(output, "$", `\$`, -1)
	output = strings.Replace(output, "%", `\%`, -1)
	output = strings.Replace(output, "&", `\&`, -1)
	output = strings.Replace(output, "~", `\~`, -1)
	output = strings.Replace(output, "_", `\_`, -1)
	output = strings.Replace(output, "^", `\^`, -1)
	output = strings.Replace(output, "{", `\{`, -1)
	output = strings.Replace(output, "}", `\}`, -1)
	return
}

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

// Returns a closure that generates a \mainline{...} LaTeX command with the next
// "nbplies" noves and the resulting chessboard, starting from the beginning. It
// also shows other information for every single move. In case the game has been
// exhausted it returns the empty string and io.EOF.
//
// This function specifically takes care of special LaTeX character appearing in
// any comment
func (game *PgnGame) getMainLineWithComments(nbplies int) func() (string, error) {

	// Initially, all moves are generated from the first one
	start := 0

	// return a closure which produces the LaTeX command for the next nbplies
	// moves
	return func() (string, error) {

		// Ensure the game has not been fully reported yet
		if start >= len(game.moves) {

			// If so, return the empty string and io.EOF
			return "", io.EOF
		}

		output := ""

		// the variable newMainLine is used to determine whether the next move
		// should start with a \mainline or not. Obviously, this is true at the
		// beginning
		newMainLine := true

		// Iterate from the given position
		last := min(start+nbplies, len(game.moves))
		for idx, move := range game.moves[start:last] {

			// if we are starting a new mainline (either because we are about to
			// generate the first move or because a comment or other information
			// was printed in the last iteration)
			if newMainLine {
				output += `\mainline{`
			}

			// now in case we are either starting a new mainline or it is
			// white's move, then show all the details of the move including
			// counter and color prefix
			if newMainLine || move.color == 1 {

				// now, show the actual move with all details
				output += fmt.Sprintf("%v%v %v ", move.number, move.getColorPrefix(), move.shortAlgebraic)
			} else {

				// otherwise, just show the actual move
				output += fmt.Sprintf("%v ", move.shortAlgebraic)
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
					output += fmt.Sprintf("\\textcolor{CadetBlue}{%v}", substituteLaTeX(move.comments))
				}
			} else if idx == last-start-1 {

				// if this is the last move to show in this mainline, and no
				// emt/comments were produced, then make sure to close the mainline
				// anyway
				output += "} "
			}

			// and check whether a new mainline has to be started in the
			// next iteration
			newMainLine = (move.emt != -1 || move.comments != "")
		}

		// update the position of the next location to examine
		start = last

		// and return the string produced so far
		return output, nil
	}
}

// Produces a LaTeX string with the list of moves of this game along with the
// different annotations.
//
// This method successively processes the moves in this PgnGame until a comment
// is found.
//
// It is intended to be used in LaTeX templates
func (game *PgnGame) GetLaTeXMovesWithComments() string {

	// capture the closure that generates the moves
	result, _ := game.getMainLineWithComments(len(game.moves))()

	// and return all moves of this game
	return result
}

// Produces a LaTeX string with a long table showing the moves every nbplies and
// the chess board
//
// This method successively processes the moves in this PgnGame until a comment
// is found.
//
// It is intended to be used in LaTeX templates
func (game *PgnGame) GetLaTeXMovesWithCommentsTabular(width1, width2 string, nbplies int) (output string) {

	// Declare a long table which can span over several pages to show the entire
	// game
	output += fmt.Sprintf(`\begin{longtable}{>{\centering\arraybackslash}m{%v} | >{\centering\arraybackslash}m{%v}}`, width1, width2)
	output += "\n"

	// Get the generator of the mainlines that shows the chess board after
	// nbplies plies
	generator := game.getMainLineWithComments(nbplies)

	// Now, produce the lines of the table. Each line shows a mainline (along
	// with comments and other information) in the left cell, and the resulting
	// chess board to the right
	for {

		// get the next mainline to show and in case the game was exhausted,
		// exit from the main loop
		if mainline, err := generator(); err == io.EOF {
			break
		} else {

			// Otherwise, add a new line to the table
			output += fmt.Sprintf("%v & \\chessboard[smallboard,print,showmover=true] \\\\ \n", mainline)
		}
	}

	// Before leaving ensure the longtable environment is closed
	output += "\n"
	output += `\end{longtable}`
	output += "\n"

	// and return the string computed so far
	return
}

// A field is either a tag of the receiver game, or a value that can be
// extracted from it (such as "Moves" or "Result")
//
// This function specifically takes care of special LaTeX character appearing in
// any comment
//
// This function does not perform any error checking. In case a field does not
// exist it returns the empty string and it should be the author of the template
// who should handle such cases
//
// This function is to be used in LaTeX templates
func (game *PgnGame) GetField(field string) string {

	// -- Id
	if field == "Id" {

		// In case the id of this game is requested just return it as a string
		return fmt.Sprintf("%d", game.id)
	}

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
	if value, ok := game.tags[field]; ok {
		return substituteLaTeX(fmt.Sprintf("%v", value))
	}

	// at this point, the field is known to not exist
	return ""
}

// Return an index entry of a specific game for any slice of fields. The first
// argument serves to determine where to add a horizontal single rule so that
// every block consists of sep entries.
//
// It assumes that each game is properly indexed with labels (with the usage of
// game.SetLabel ())
//
// It is intended to be used in LaTeX templates
func (game *PgnGame) GetIndexEntry(sep int, fields []any) (output string) {

	// for all requested fields
	for idx, field := range fields {

		// cast this field into a string
		if value, ok := field.(string); !ok {
			log.Fatalf(" It was not possible to cast field '%v' into a string\n", field)
		} else {

			// Ids are slightly different because they have to be generated with
			// a hyperref
			if value == "Id" {
				output += fmt.Sprintf("\\hyperref[game:%v]{\\#%v}", game.id, game.id)
			} else {

				// Otherwise just reteurn the value of the given field
				output += game.GetField(value)
			}
		}

		// in case this is not the last entry add a column separator
		if idx < len(fields)-1 {
			output += ` & `
		}

	}

	// And end this entry
	output += `\\`

	// in case a block has been ended with this entry then add a single
	// horizontal rule
	if game.id%sep == 0 {
		output += `\midrule`
	}

	// and return the output
	return
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
