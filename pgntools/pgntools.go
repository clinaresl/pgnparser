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

import "regexp"

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

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
