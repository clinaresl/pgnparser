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

// global variables (to the package)
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

// the following regexp captures all the information given from the textual
// description of a move in different groups as follows:
//
// Group #1: Piece
// Group #2: Qualifier
// Group #3: Capture ('x' only if this is a capture)
// Group #4: Target square
// Group #5: Promotion (in the form =<piece>)
// Group #6: Castling (either 'O-O' or 'O-O-O')
var reTextualMove = regexp.MustCompile(`([PNBRQK]?)([a-h]?[1-8]?)(x?)([a-h][1-8]|[NBRQK])(\=[PNBRQK])?|(O(?:-?O){1,2})[\+#]?(\s*[\!\?]+)?`)

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

// The following simple regular expression is used to distinguish criteria given
// for the creation of histograms
var reCriteria = regexp.MustCompile(`\s*;\s*`)

// The following regular expression is used to distinguish the name of a
// var/bool expression from the var/bool expression
var reHistogramName = regexp.MustCompile(`\s*:\s*`)

// The following regular expression is used to distinguish criteria from sorting
// operands
var reSorting = `\s*(?P<direction>[<>])\s*(?P<criteria>.+)\s*`

// the following map stores the translation of literal coordinates to integers
// used to access a PgnBoard
var coords map[string]int

// the following map stores the translation of integer coordinates to literal
// coordinates used to access a PgnBoard
var literal map[int]string

// the following structure contains information about the coordinates (in
// literal form) from which each piece (represented as a white piece, but pawns
// which preserve their color) can access a specific location of the board. For
// example:
//
//	threats ["e4"][WPAWN] = [19][20, 12][21]
//
// which means that a white pawn can access location "e4" from squares 12 (e2),
// 19 (d3, by capturing a piece in e4), 20 (e3) and 21 (f3, again by capturing).
//
// Note that all the locations from which e4 can be accessed are stored in
// separate lists. Each list represents a specific direction.
var threats map[string]map[content][][]int

// The following map relates each content with its utf-8 representation
var utf8 map[content]rune

// The following counter is used to generate LaTeX references
var counter int = 0

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
