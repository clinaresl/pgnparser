/*
  pgntools.go
  Description: Simple tools for handling pgn files
  -----------------------------------------------------------------------------

  Started on  <Wed May  6 15:38:56 2015 Carlos Linares Lopez>
  Last update <martes, 29 marzo 2016 21:07:32 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

// This package provides a number of simple services for accessing and handling
// chess games in PGN format
package pgntools

import (
	"regexp"
)

// global variables
// ----------------------------------------------------------------------------

// ungrouped regexps -- they are used just to recognize different chunks of a
// string
// ----------------------------------------------------------------------------
// the following regexp matches a string with an arbitrary number of
// comments
var reTags = regexp.MustCompile(`(\[\s*\w+\s*"[^"]*"\s*\]\s*)+`)

// the following regexp matches an arbitrary sequence of moves which are
// identified by a number, a color (symbolized by either one dot for white or
// three dots for black) and the move in algebraic format. Moves can be followed
// by an arbitrary number of comments
var reMoves = regexp.MustCompile(`(?:(\d+)(\.|\.{3})\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*)+`)

// the outcome is one of the following strings "1-0", "0-1" or "1/2-1/2"
var reOutcome = regexp.MustCompile(`(1\-0|0\-1|1/2\-1/2|\*)`)

// the following regexp is used to parse the description of an entire game,
// including the tags, list of moves and final outcome. It consists of a
// concatenation of the previous expressions where an arbitrary number of spaces
// is allowed between them
var reGame = regexp.MustCompile(`\s*(\[\s*(?P<tagname>\w+)\s*"(?P<tagvalue>[^"]*)"\s*\]\s*)+\s*(?:(\d+)(\.|\.{3})\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*)+\s*(1\-0|0\-1|1/2\-1/2|\*)\s*`)

// grouped regexps -- they are used to extract relevant information from a
// string
// ----------------------------------------------------------------------------

// the following regexp matches a string with an arbitrary number of
// comments. Groups are used to extract the tag name and value
var reGroupTags = regexp.MustCompile(`\[\s*(?P<tagname>\w+)\s*"(?P<tagvalue>[^"]*)"\s*\]\s*`)

// this regexp is used just to extract the textual description of a single move
// which might be preceded by a move number and color identification
var reGroupMoves = regexp.MustCompile(`(?:(?P<moveNumber>\d+)?(?P<color>\.|\.{3})?\s*(?P<moveValue>(?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*)`)

// comments following any move are matched with the following regexp. Note that
// comments are expected to be matched at the beginning of the string (^) and
// its occurrence is required to happen precisely once. This makes sense since
// the whole string is parsed in chunks
var reGroupComment = regexp.MustCompile(`^(?P<comment>{[^{}]*})\s*`)

// A specific type of comments provided by ficsgames.org is the time elapsed to
// make the current move. This is parsed in the following expression. Again,
// note that this expression matches the beginning of the string
var reGroupEMT = regexp.MustCompile(`^{\[%emt (?P<emt>\d+\.\d*)\]}`)

// Groups are used in the following regexp to extract the score of every player
var reGroupOutcome = regexp.MustCompile(`(?P<score1>1/2|0|1)\-(?P<score2>1/2|0|1)`)

// the following regexp is used to match the different sorting criteria

// sorting criteria consists of a sorting direction and a particular variable
// used as a key for sorting games. The direction is specified with either < or
// > meaning increasing and decreasing order respectively; the variable to use
// is preceded by '%' (there is no need actually to use that prefix and this is
// done only for the sake of consistency across different commands of pgnparser)
var reSortingCriteria = regexp.MustCompile(`^\s*(<|>)\s*%([A-Za-z]+)\s*`)

// the following regexps are used to process histogram command lines

// A histogram command line might consist of a title and a variable name
var reHistogramCmdVar = regexp.MustCompile(`^\s*([A-Za-z0-9]+)\s*:\s*%([A-Za-z]+)\s*`)

// Also, a histogram command line might consist of the definition of a case
// which consists of a number of different regular expressions
var reHistogramCmdCase = regexp.MustCompile(`^\s*(?P<title>[A-Za-z0-9]+)\s*:\s*\[(?P<cases>[^\]]+)\]`)
var reHistogramCmdSubcase = regexp.MustCompile(`^\s*(?P<subtitle>[A-Za-z0-9]+)\s*:\s*{(?P<expression>[^}]+)}\s*`)

// functions
// ----------------------------------------------------------------------------

// Initializes various structures necessary for the proper execution of this
// package: 1. the map of coordinates to specific cells in the chess board; 2.
// the utf-8 representation of each content
func init() {

	// Coordinates

	// first, initialize the transformation from literal coordinates to
	// indexes used to access a PgnBoard
	coords = make(map[string]int)
	for row := 0; row < 8; row++ {
		for column := 0; column < 8; column++ {

			// and store the transformation from literal coordinates
			// to integers
			coords[string('a'+byte(column))+string('0'+byte(1+row))] = row*8 + column
		}
	}

	// second, makes the opposite and compute the translation from integer
	// coordinates to literal coordinates
	literal = make(map[int]string)
	for index := 0; index < 64; index++ {
		literal[index] = string('a'+byte(index%8)) + string('0'+byte(1+index/8))
	}

	// now, compute all threats
	threats = make(map[string]map[content][][]int)

	// for all squares of the board represented as a pair (row,
	// column)
	for row := 0; row < 8; row++ {
		for column := 0; column < 8; column++ {

			threat := make(map[content][][]int) // create an empty map

			// and all pieces where color is ignored but for the
			// pawns (because they are the only chess pieces which
			// have direction) are computed
			for piece := BKING; piece <= WKING; piece++ {
				if piece == BLANK {
					continue
				}
				threat[piece] = getThreat(row*8+column, piece)
			}
			threats[string('a'+byte(column))+string('0'+byte(1+row))] = threat
		}
	}

	// utf-8 representation of contents
	utf8 = make(map[content]rune)
	utf8[BKING] = '♚'
	utf8[BQUEEN] = '♛'
	utf8[BROOK] = '♜'
	utf8[BBISHOP] = '♝'
	utf8[BKNIGHT] = '♞'
	utf8[BPAWN] = '♟'
	utf8[BLANK] = ' '
	utf8[WKING] = '♔'
	utf8[WQUEEN] = '♕'
	utf8[WROOK] = '♖'
	utf8[WBISHOP] = '♗'
	utf8[WKNIGHT] = '♘'
	utf8[WPAWN] = '♙'
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
