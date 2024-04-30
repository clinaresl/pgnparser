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
	"log"     // logging services
	"regexp"  // pgn files are parsed with a regexp
	"sort"    // for sorting games
	"strconv" // to convert from strings to other types

	// import a user package to manage paths
	"github.com/clinaresl/pgnparser/fstools"

	// import the parser of propositional formulae
	"github.com/clinaresl/pgnparser/pfparser"
)

// global variables
// ----------------------------------------------------------------------------

// ungrouped regexps -- they are used just to recongize different chunks of a
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

// functions
// ----------------------------------------------------------------------------

// Return a slice with all tags in the given string.
func getTags(pgn string) (tags map[string]dataInterface) {

	// create the map
	tags = make(map[string]dataInterface)

	// get information about all pgn tags in the given string
	for _, tag := range reGroupTags.FindAllStringSubmatchIndex(pgn, -1) {

		// process this tag in case it contains relevant info. Every
		// matching should consist of a slice with 6 indices:
		// <begin/end>-string, <begin/end>-tagname, <begin/end>-tagvalue
		if len(tag) >= 6 {

			// add this tag to the map to return. In case this
			// string can be interpreted as an integer number
			value, err := strconv.Atoi(pgn[tag[4]:tag[5]])
			if err == nil {

				// then store it as an integer constant
				tags[pgn[tag[2]:tag[3]]] = constInteger(value)
			} else {

				// otherwise, store it as a string constant
				tags[pgn[tag[2]:tag[3]]] = constString(pgn[tag[4]:tag[5]])
			}
		}
	}

	return
}

// Return a slice of PgnMove with the information in the string 'pgn' which
// shall consist of a legal transcription of legal PGN moves that might be
// annotated (an arbitrary number of times) or not. 'emt' annotations are also
// acknowledged and their information is added to the slice of PgnMove
func getMoves(pgn string) (moves []PgnMove) {

	moveNumber := -1     // initialize the move counter to unknown
	color := 0           // initialize the color to unknown
	var moveValue string // move actually parsed in PGN format
	var emt float64      // elapsed move time
	var comments string  // comments of each move
	var err error

	// process plies in sequence until the whole string is exhausted
	for len(pgn) > 0 {

		// get the next move
		tag := reGroupMoves.FindStringSubmatchIndex(pgn)

		// reGroupMoves contains three groups and therefore legal
		// matches contain 8 characters
		if len(tag) >= 8 {

			// if a move number and color (. or ...) specifier has
			// been found, then process all groups in this matching
			if tag[2] >= 0 && tag[4] >= 0 {

				// update the move counter
				moveNumber, err = strconv.Atoi(pgn[tag[2]:tag[3]])
				if err != nil {
					log.Fatalf(" Fatal error while converting the move number")
				}

				// and the color, in case only one character
				// ('.') is found, this is white's move,
				// otherwise, it is black's move
				if tag[5]-tag[4] == 1 {
					color = 1
				} else {
					color = -1
				}
			} else {

				// otherwise, assume that this is the opponent's
				// move
				color *= -1
			}

			// and in any case extract the move value
			moveValue = pgn[tag[6]:tag[7]]
		}

		// and move forward
		pgn = pgn[tag[1]:]

		// are there any comments immediately after? The following loop
		// aims at processing an arbitrary number of comments
		emt = -1.0    // initialize the elapsed move time to unknown
		comments = "" // initialize the comments to the empty string
		for reGroupComment.MatchString(pgn) {

			// Yeah, a comment has been found! extract it
			tag = reGroupComment.FindStringSubmatchIndex(pgn)

			// is this an emt field?
			if reGroupEMT.MatchString(pgn) {
				tagEMT := reGroupEMT.FindStringSubmatchIndex(pgn)
				emt, err = strconv.ParseFloat(pgn[tagEMT[2]:tagEMT[3]], 32)
				if err != nil {
					log.Fatalf(" Fatal error while converting emt")
				}
			} else {
				// if not, then just add these comments. In case
				// some comments were already written, make sure
				// to add this in a new line
				if len(comments) > 0 {
					comments += "\r\n"
				}
				comments += pgn[1+tag[2] : tag[3]-1]
			}
			pgn = pgn[tag[1]:]
		}

		// and add this move to the list of moves to return unless there
		// are unknown fields
		if moveNumber == -1 || color == 0 {
			log.Fatalf(" Either the move number or the color were incorrect")
		}
		moves = append(moves, PgnMove{moveNumber, color, moveValue, float32(emt), comments})
	}

	return
}

// Return an instance of PgnOutcome with the score of each player as specified
// in the given string
func getOutcome(pgn string) (outcome PgnOutcome) {

	// get information about the outcome as given in pgn
	tag := reGroupOutcome.FindStringSubmatchIndex(pgn)

	// process this tag in case it contains 6 indices: <begin/end>-string,
	// <begin/end>-scorewhite, <begin/end>-scoreblack
	if len(tag) >= 6 {

		// if the first tag is three characters long, then this is a
		// draw
		if tag[3]-tag[2] == 3 {
			outcome = PgnOutcome{0.5, 0.5}
		} else {

			// otherwise, one side won the match
			scoreWhite, err := strconv.Atoi(pgn[tag[2]:tag[3]])
			if err != nil {
				log.Fatalf(" Ilegal outcome found in string '%s'", pgn)
			}
			outcome = PgnOutcome{float32(scoreWhite), 1.0 - float32(scoreWhite)}
		}
	}

	return
}

// Return the contents of a chess game from the full transcription of a chess
// game given in a string in PGN format. In case verbose is given, it shows
// additional information
func getGameFromString(pgn string, verbose bool) PgnGame {

	// create variables to store different sections of a single PGN game
	var strTags, strMoves, strOutcome string

	// find the tags of the first game in pgn
	endpoints := reTags.FindStringIndex(pgn)
	if endpoints == nil {
		log.Fatalf(" the PGN tags have not been found")
	} else {

		// copy the section of the tags and move forward in the pgn string
		strTags = pgn[endpoints[0]:endpoints[1]]
		pgn = pgn[endpoints[1]:]

		if verbose {
			log.Printf(" Legal tags of a PGN game have been found:\n%v", strTags)
		}

		// now, check that this is followed by a legal transcription of
		// chess moves in PGN format
		endpoints = reMoves.FindStringIndex(pgn)
		if endpoints == nil {
			log.Fatalf(" no transcription of legal moves has been found")
		} else {

			// copy the section with the chess moves and move
			// forward in the pgn string
			strMoves = pgn[endpoints[0]:endpoints[1]]
			pgn = pgn[endpoints[1]:]

			if verbose {
				log.Printf(" A legal transcription of chess moves in PGN format has been found:\n%v\n\n", strMoves)
			}

			// now, check that the final result is properly written
			endpoints = reOutcome.FindStringIndex(pgn)
			if endpoints == nil {
				log.Fatalf(" there is no legal transcription of the final result")
			} else {

				// again, copy the section with the final
				// outcome and move forward in the pgn file
				strOutcome = pgn[endpoints[0]:endpoints[1]]
				pgn = pgn[endpoints[1]:]

				if verbose {
					log.Printf(" The outcome has been properly identified:\n%v\n\n", strOutcome)
				}
			}
		}
	}

	// now, just process the different chunks extracted previously and store
	// them in the game to return
	tags := getTags(strTags)          // -- PGN tags
	moves := getMoves(strMoves)       // -- PGN moves
	outcome := getOutcome(strOutcome) // -- PGN outcome
	return PgnGame{tags, moves, outcome}
}

// Return the contents of all chess games that satisfiy the given query from the
// specified string which shall be formattted in PGN format. Games are sorted
// according to the criteria given in sort if any is given; if not, they are
// listed in the same order they were found in the file. For each game, the
// board is shown every showboard plies
//
// In case verbose is given, it shows additional information
func GetGamesFromString(pgn string, showboard int, query string, sortString string, verbose bool) (games PgnCollection) {

	var err error
	var logEvaluator pfparser.LogicalEvaluator

	// since parsing queries affect its contents, make a backup copy
	queryString := query

	// in case a string has been given ...
	if query != "" {

		// ... parse it (with null depth, of course!)
		logEvaluator, err = pfparser.Parse(&queryString, 0)
		if err != nil {
			log.Fatal(err)
		}
	}

	// just iterate over the string extracting the information of every game
	for reGame.MatchString(pgn) {

		// In case a match has been found, extract the next game
		tag := reGame.FindStringSubmatchIndex(pgn)

		// Parse this game and add it to the slice of games to return
		game := getGameFromString(pgn[tag[0]:tag[1]], verbose)

		symtable := make(map[string]pfparser.RelationalInterface)

		// in case a query has been given, then process it
		if query != "" {

			// first, start by creating a symbol table with all the
			// information appearing in the headers of this game
			for key, content := range game.tags {

				// first, verify whether this is an integer
				value, ok := content.(constInteger)
				if ok {

					symtable[key] = pfparser.ConstInteger(value)
				} else {

					// if not, check if it is a string
					value, ok := content.(constString)
					if ok {
						symtable[key] = pfparser.ConstString(value)
					} else {
						log.Fatal(" Unknown type")
					}
				}
			}
		}

		// if no query was given, or if one was given and this game
		// satisfies it
		if query == "" ||
			(query != "" &&
				logEvaluator.Evaluate(symtable) == pfparser.TypeBool(true)) {

			// parse all moves of this game
			game.ParseMoves(showboard)

			games.slice = append(games.slice, game)
			games.nbGames += 1
		}

		// and move forward
		pgn = pgn[tag[1]:]
	}

	// at this point the resulting string should be empty, otherwise, an
	// error should be raised
	if pgn != "" {
		log.Fatalf(" Some games were not processed: %q\n\n", pgn)
	}

	// and finally sort the games in case a sorting string was given
	if sortString != "" {
		games.GetSortDescriptor(sortString)
		sort.Sort(games)
	}

	// and return the slice computed so far
	return
}

// Return the contents of all chess games that satisfiy the given query from the
// specified file which shall be formattted in PGN format. Games are sorted
// according to the criteria given in sort if any is given; if not, they are
// listed in the same order they were found in the file. For each game, the
// board is shown every showboard plies
//
// In case verbose is given, it shows additional information
func GetGamesFromFile(pgnfile string, showboard int, query string, sortString string, verbose bool) (games PgnCollection) {

	// Open and read the given file and retrieve its contents
	contents := fstools.Read(pgnfile, -1)
	strContents := string(contents[:len(contents)])

	// and now, just return the results of parsing these contents
	return GetGamesFromString(strContents, showboard, query, sortString, verbose)
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
