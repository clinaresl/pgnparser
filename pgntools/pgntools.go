/* 
  pgntools.go
  Description: Simple tools for handling pgn files
  ----------------------------------------------------------------------------- 

  Started on  <Wed May  6 15:38:56 2015 Carlos Linares Lopez>
  Last update <sábado, 09 mayo 2015 00:17:36 Carlos Linares Lopez (clinares)>
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

// ungrouped regexps -- they are used just to recongize different chunks of a
// string
// ----------------------------------------------------------------------------
var reTags = regexp.MustCompile (`(\[\s*(?P<tagname>\w+)\s*"(?P<tagvalue>[^"]*)"\s*\]\s*)+`)

var reMoves = regexp.MustCompile (`(?:(\d+)(\.|\.{3})\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*)+`)

var reOutcome = regexp.MustCompile (`(1\-0|0\-1|1/2\-1/2|\*)`)

// the following regexp is used to parse the description of an entire game,
// including the tags, list of moves and final outcome. It consists of a
// concatenation of the previous expressions where an arbitrary number of spaces
// is allowed
var reGame = regexp.MustCompile (`\s*(\[\s*(?P<tagname>\w+)\s*"(?P<tagvalue>[^"]*)"\s*\]\s*)+\s*(?:(\d+)(\.|\.{3})\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*)+\s*(1\-0|0\-1|1/2\-1/2|\*)\s*`)

// grouped regexps -- they are used to extract relevant information from a
// string
// ----------------------------------------------------------------------------
var reGroupTags = regexp.MustCompile (`\[\s*(?P<tagname>\w+)\s*"(?P<tagvalue>[^"]*)"\s*\]\s*`)

var reGroupFullMoves = regexp.MustCompile (`(?:(?P<moveNumber>\d+)(?P<color>\.|\.{3})\s*(?P<moveValue1>(?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*(?P<comment1>{[^{}]*})?\s*(?P<moveValue2>(?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*(?P<comment2>{[^{}]*})?\s*)`)

var reGroupMoves = regexp.MustCompile (`(?:(?P<moveNumber>\d+)?(?P<color>\.|\.{3})?\s*(?P<moveValue>(?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*)`)

// note that comments are expected to be matched at the beginning of the string
// (^) and its occurrence is required to happen precisely once. This makes sense
// since the whole string is parsed in chunks
var reGroupComment = regexp.MustCompile (`^(?P<comment>{[^{}]*})\s*`)

// A specific type of comments provided by ficsgames.org is the time elapsed to
// make the current move. This is parsed in the following expression. Again,
// note that this expression matches the beginning of the string
var reGroupEMT = regexp.MustCompile (`^{\[%emt (?P<emt>\d+\.\d*)\]}`)

var reGroupOutcome = regexp.MustCompile (`(?P<score1>1/2|0|1)\-(?P<score2>1/2|0|1)`)


// typedefs
// ----------------------------------------------------------------------------
type PgnTag struct {

	// every tag is identified with a name and value
	name, value string;
}

type PgnMove struct {

	// current move number
	moveNumber int;

	// color: 1=white; -1=black
	color int;
	
	// move value in PGN format
	moveValue string;

	// elapsed move time
	emt float32;

	// comments - in case there are various comments, each one is added
	// after \r\n
	comments string;
}

type PgnOutcome struct {

	// the outcome is qualified with two floating-point numbers with the
	// score of each player
	scoreWhite, scoreBlack float32;
}

type PgnGame struct {

	// A game consists of a collection of tags
	tags []PgnTag;

	// a sequence of moves with related information such as comments and emt
	moves []PgnMove;

	// and an outcome
	outcome PgnOutcome;
}

// Methods
// ----------------------------------------------------------------------------

// the following methods overwrite the string output method
func (tag PgnTag) String () string {
	return fmt.Sprintf ("%v: %v", tag.name, tag.value)
}

func (move PgnMove) String () string {
	if move.color == 1 {
		return fmt.Sprintf ("%3v. %v [%v]\n[%v]", move.moveNumber, move.moveValue, move.emt, move.comments)
	}
	return fmt.Sprintf ("%3v... %v [%v]\n[%v]", move.moveNumber, move.moveValue, move.emt, move.comments)	
}

func (outcome PgnOutcome) String () string {
	return fmt.Sprintf ("%v - %v", outcome.scoreWhite, outcome.scoreBlack)
}

func (game PgnGame) String () string {
	return fmt.Sprintf ("%v\n%v\n%v", game.tags, game.moves, game.outcome)
}

// the following are getters over the attributes of a PgnGame
func (game *PgnGame) GetTags () []PgnTag {
	return game.tags
}

func (game *PgnGame) GetMoves () []PgnMove {
	return game.moves
}

func (game *PgnGame) GetOutcome () PgnOutcome {
	return game.outcome
}

// next, there are a few methods that provide additional support over the
// information stored in various attributes

// GetTagValue return the value of a specific tag and nil if it exists or any
// value and err in case it does not exist
func (game *PgnGame) GetTagValue (name string) (value string, err error) {

	// usually there are just a few tags, so that a linear search does not
	// cause any harm
	for _, tag := range game.tags {

		// in case this tag contains the specified name return it
		if tag.name == name {
			return tag.value, nil
		}
	}

	// when getting here, the required tag has not been found
	return "", errors.New ("tag not found!")
}

// ShowHeader summarizes the main information stored in the tags of a specific
// game
func (game *PgnGame) ShowHeader () string {

	// first, verify that all necessary tags are available
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

	outcome := game.GetOutcome ()
	
	return fmt.Sprintf (" | %v %v | %-18v (%4v) | %-18v (%4v) | %v | %v | %2v | %3v-%-3v |", date, time, white, whiteELO, black, blackELO, ECO, timeControl, moves, outcome.scoreWhite, outcome.scoreBlack)
}

// functions
// ----------------------------------------------------------------------------

// getTags
//
// Return a slice with all tags in the given string. 
// ----------------------------------------------------------------------------
func getTags (pgn string) (tags []PgnTag) {

	// get information about all pgn tags in the given string
	for _, tag := range reGroupTags.FindAllStringSubmatchIndex (pgn, -1) {

		// process this tag in case it contains relevant info. Every
		// matching should consist of a slice with 6 indices:
		// <begin/end>-string, <begin/end>-tagname, <begin/end>-tagvalue
		if len (tag) >= 6 {

			// create an instance of a pgntag and add it to the
			// slice to return
			tags = append (tags, PgnTag{pgn[tag[2]:tag[3]], pgn[tag[4]:tag[5]]})
		}
	}

	return
}

// getMoves
//
// Return a slice of PgnMove with the information in the string 'pgn' which
// shall consist of a legal transcription of legal PGN moves that might be
// annotated (an arbitrary number of times) or not. 'emt' annotations are also
// acknowledged and their information is added to the slice of PgnMove
// ----------------------------------------------------------------------------
func getMoves (pgn string) (moves []PgnMove) {

	moveNumber := -1		// initialize the move counter to unknown
	color := 0;			// initialize the color to unknown
	var moveValue string;		// move actually parsed in PGN format
	var emt float64;		// elapsed move time
	var comments string;		// comments of each move
	var err error;
	
	// process plies in sequence until the whole string is exhausted
	for ;len (pgn) > 0; {

		// get the next move 
		tag := reGroupMoves.FindStringSubmatchIndex (pgn)

		// reGroupMoves contains three groups and therefore legal
		// matches contain 8 characters
		if len (tag) >= 8 {

			// if a move number and color (. or ...) specifier has
			// been found, then process all groups in this matching
			if tag[2]>=0 && tag[4]>=0 {
			
				// update the move counter
				moveNumber, err = strconv.Atoi (pgn[tag[2]:tag[3]])
				if err != nil {
					log.Fatalf (" Fatal error while converting the move number")
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
		pgn = pgn [tag[1]:]

		// are there any comments immediately after? The following loop
		// aims at processing an arbitrary number of comments
		emt = -1.0;		// initialize the elapsed move time to unknown
		comments = ""		// initialize the comments to the empty string
		for ;reGroupComment.MatchString (pgn); {

			// Yeah, a comment has been found! extract it
			tag = reGroupComment.FindStringSubmatchIndex (pgn)

			// is this an emt field?
			if reGroupEMT.MatchString (pgn) {
				tagEMT := reGroupEMT.FindStringSubmatchIndex (pgn)

				emt, err = strconv.ParseFloat (pgn[tagEMT[2]:tagEMT[3]], 32)
				if err != nil {
					log.Fatalf (" Fatal error while converting emt")
				}
			} else {
				// if not, then just add these comments. In case
				// some comments were already written, make sure
				// to add this in a new line
				if len (comments) > 0 {
					comments += "\r\n"
				}
				comments += pgn[1+tag[2]:tag[3]-1]
			}
			
			
			pgn = pgn[tag[1]:]
		}

		// and add this move to the list of moves to return unless there
		// are unknown fields
		if moveNumber == -1 || color == 0 {
			log.Fatalf (" Either the move number or the color were incorrect")
		}
		moves = append (moves, PgnMove{moveNumber, color, moveValue, float32 (emt), comments})
	}
	
	return
}

// getOutcome
//
// Return an instance of PgnOutcome with the score of each player as specified
// in the given string
// ----------------------------------------------------------------------------
func getOutcome (pgn string) (outcome PgnOutcome) {

	// get information about the outcome as given in pgn
	tag := reGroupOutcome.FindStringSubmatchIndex (pgn)

	// process this tag in case it contains 6 indices: <begin/end>-string,
	// <begin/end>-scorewhite, <begin/end>-scoreblack
	if len (tag) >= 6 {

		// if the first tag is three characters long, then this is a
		// draw
		if tag [3]-tag[2] == 3 {
			outcome = PgnOutcome {0.5, 0.5}
		} else {

			// otherwise, one side won the match
			scoreWhite, err := strconv.Atoi (pgn[tag[2]:tag[3]])
			if err != nil {
				log.Fatalf (" Ilegal outcome found in string '%s'", pgn)
			}
			outcome = PgnOutcome {float32 (scoreWhite), 1.0-float32 (scoreWhite)}
		}
	}

	return
}

// getGameFromString
//
// Return the contents of a chess game from the full transcription of a chess
// game given in a string in PGN format. In case verbose is given, it shows
// additional information
// ----------------------------------------------------------------------------
func getGameFromString (pgn string, verbose bool) PgnGame {

	// create variables to store different sections of a single PGN game
	var strTags, strMoves, strOutcome string;
	
	// find the tags of the first game in pgn
	endpoints := reTags.FindStringIndex (pgn); if endpoints == nil {
		log.Fatalf (" the PGN tags have not been found")
	} else {

		// copy the section of the tags and move forward in the pgn string
		strTags = pgn [endpoints[0]:endpoints[1]]
		pgn = pgn [endpoints[1]:]

		if verbose {
			log.Printf (" Legal tags of a PGN game have been found:\n%v", strTags)
		}

		// now, check that this is followed by a legal transcription of
		// chess moves in PGN format
		endpoints = reMoves.FindStringIndex (pgn); if endpoints == nil {
			log.Fatalf (" no transcription of legal moves has been found")
		} else {

			// copy the section with the chess moves and move
			// forward in the pgn string
			strMoves = pgn[endpoints[0]:endpoints[1]]
			pgn = pgn [endpoints[1]:]

			if verbose {
				log.Printf (" A legal transcription of chess moves in PGN format has been found:\n%v\n\n", strMoves)
			}

			// now, check that the final result is properly written
			endpoints = reOutcome.FindStringIndex (pgn); if endpoints == nil {
				log.Fatalf (" there is no legal transcription of the final result")
			} else {

				// again, copy the section with the final
				// outcome and move forward in the pgn file
				strOutcome = pgn[endpoints[0]:endpoints[1]]
				pgn = pgn [endpoints[1]:]

				if verbose {
					log.Printf (" The outcome has been properly identified:\n%v\n\n", strOutcome)
				}
			}
		}
	}

	// now, just process the different chunks extracted previously and store
	// them in the game to return	
	tags := getTags (strTags)                                // -- PGN tags
	moves := getMoves (strMoves)                            // -- PGN moves
	outcome := getOutcome (strOutcome)                    // -- PGN outcome
	return PgnGame {tags, moves, outcome}
}

// GetGamesFromString
//
// Return the contents of all chess games included the given string in PGN
// format. In case verbose is given, it shows additional information
// ----------------------------------------------------------------------------
func GetGamesFromString (pgn string, verbose bool) (games []PgnGame) {

	// just iterate over the string extracting the information of every game
	for ;reGame.MatchString (pgn); {

		// In case a match has been found, extract the next game
		tag := reGame.FindStringSubmatchIndex (pgn)

		// Parse this game and add it to the slice of games to return
		game := getGameFromString (pgn[tag[0]:tag[1]], verbose)
		games = append (games, game)

		// and move forward
		pgn = pgn[tag[1]:]
	}

	// at this point the resulting string should be empty, otherwise, an
	// error should be raised
	if pgn != "" {
		log.Fatalf (" Some games were not processed: %q\n\n", pgn)
	}
	
	// and return the slice computed so far
	return
}

// GetGamesFromFile
//
// Return the contents of all chess games in the given file in PGN format. In
// case verbose is given, it shows additional information
// ----------------------------------------------------------------------------
func GetGamesFromFile (pgnfile string, verbose bool) (games []PgnGame) {

	// Open and read the given file and retrieve its contents
	contents := fstools.Read (pgnfile, -1)
	strContents := string (contents[:len (contents)])

	// and now, just return the results of parsing these contents
	return GetGamesFromString (strContents, verbose)
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
