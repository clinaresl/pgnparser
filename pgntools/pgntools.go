/* 
  pgntools.go
  Description: Simple tools for handling pgn files
  ----------------------------------------------------------------------------- 

  Started on  <Wed May  6 15:38:56 2015 Carlos Linares Lopez>
  Last update <sábado, 09 mayo 2015 03:56:35 Carlos Linares Lopez (clinares)>
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
// the following regexp matches a string with an arbitrary number of
// comments
var reTags = regexp.MustCompile (`(\[\s*\w+\s*"[^"]*"\s*\]\s*)+`)

// the following regexp matches an arbitrary sequence of moves which are
// identified by a number, a color (symbolized by either one dot for white or
// three dots for black) and the move in algebraic format. Moves can be followed
// by an arbitrary number of comments
var reMoves = regexp.MustCompile (`(?:(\d+)(\.|\.{3})\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*)+`)

// the outcome is one of the following strings "1-0", "0-1" or "1/2-1/2"
var reOutcome = regexp.MustCompile (`(1\-0|0\-1|1/2\-1/2|\*)`)

// the following regexp is used to parse the description of an entire game,
// including the tags, list of moves and final outcome. It consists of a
// concatenation of the previous expressions where an arbitrary number of spaces
// is allowed between them
var reGame = regexp.MustCompile (`\s*(\[\s*(?P<tagname>\w+)\s*"(?P<tagvalue>[^"]*)"\s*\]\s*)+\s*(?:(\d+)(\.|\.{3})\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*((?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*({[^{}]*}\s*)*\s*)+\s*(1\-0|0\-1|1/2\-1/2|\*)\s*`)

// grouped regexps -- they are used to extract relevant information from a
// string
// ----------------------------------------------------------------------------

// the following regexp matches a string with an arbitrary number of
// comments. Groups are used to extract the tag name and value
var reGroupTags = regexp.MustCompile (`\[\s*(?P<tagname>\w+)\s*"(?P<tagvalue>[^"]*)"\s*\]\s*`)

// this regexp is used just to extract the textual description of a single move
// which might be preceded by a move number and color identification
var reGroupMoves = regexp.MustCompile (`(?:(?P<moveNumber>\d+)?(?P<color>\.|\.{3})?\s*(?P<moveValue>(?:[PNBRQK]?[a-h]?[1-8]?x?(?:[a-h][1-8]|[NBRQK])(?:\=[PNBRQK])?|O(?:-?O){1,2})[\+#]?(?:\s*[\!\?]+)?)\s*)`)

// comments following any move are matched with the following regexp. Note that
// comments are expected to be matched at the beginning of the string (^) and
// its occurrence is required to happen precisely once. This makes sense since
// the whole string is parsed in chunks
var reGroupComment = regexp.MustCompile (`^(?P<comment>{[^{}]*})\s*`)

// A specific type of comments provided by ficsgames.org is the time elapsed to
// make the current move. This is parsed in the following expression. Again,
// note that this expression matches the beginning of the string
var reGroupEMT = regexp.MustCompile (`^{\[%emt (?P<emt>\d+\.\d*)\]}`)

// Groups are used in the following regexp to extract the score of every player
var reGroupOutcome = regexp.MustCompile (`(?P<score1>1/2|0|1)\-(?P<score2>1/2|0|1)`)


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

	tags []PgnTag;               // A game consists of a collection of tags
	moves []PgnMove;                      // sequence of moves of this game
	outcome PgnOutcome;                                    // final outcome
}

type PgnCollection struct {

	slice []PgnGame                                  // collection of games
	nbGames int;                                  // number of games stored
}

// Methods
// ----------------------------------------------------------------------------

// the following methods overwrite the string output method
func (tag PgnTag) String () string {
	return fmt.Sprintf ("%v: %v", tag.name, tag.value)
}

func (move PgnMove) String () string {
	if move.color == 1 {
		return fmt.Sprintf ("%v. %v", move.moveNumber, move.moveValue)
	}
	return fmt.Sprintf (" %v ", move.moveValue)
}

func (outcome PgnOutcome) String () string {
	return fmt.Sprintf ("%v - %v", outcome.scoreWhite, outcome.scoreBlack)
}

// the following service just prints all the sequence of moves in the given game
func (game PgnGame) String () string {
	output := ""
	for _, move := range game.moves {
		output += fmt.Sprintf ("%v", move)
	}
	return output
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

// the following are getters over the attributes of a PgnCollection
func (games *PgnCollection) GetGames () []PgnGame {
	return games.slice
}

func (games *PgnCollection) GetGame (index int) PgnGame {
	return games.slice [index]
}

func (games *PgnCollection) GetNbGames () int {
	return games.nbGames
}

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

// ShowHeaders summarizes the main information stored in the tags of all games
// in the given collection
func (games *PgnCollection) ShowHeaders () string {

	// show the header
	output := " |  DBGameNo  | Date                | White                     | Black                     | ECO | Time  | Moves | Result |\n +------------+---------------------+---------------------------+---------------------------+-----+-------+-------+--------+\n"

	// and now, add to output information of every single game in the given
	// collection
	for _, game := range games.slice {
		output += game.ShowHeader () + "\n"
	}

	// and add a bottom line
	output += " +------------+---------------------+---------------------------+---------------------------+-----+-------+-------+--------+"

	// and return the string
	return output
}

// getLaTeXbody computes the main part of the LaTeX document that shows
// information of a specific game
func (game PgnGame) getLaTeXbody () string {

	// first, verify that all necessary tags are available
	event, err := game.GetTagValue ("Event")
	if err != nil {
		log.Fatalf ("Event not found!")
	}
	
	date, err := game.GetTagValue ("Date")
	if err != nil {
		log.Fatalf ("Date not found!")
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

	var scoreWhite, scoreBlack string;
	outcome := game.GetOutcome ()
	if outcome.scoreWhite == 0.5 {
		scoreWhite, scoreBlack = `\textonehalf`, `\textonehalf`
	} else if outcome.scoreWhite == 1 {
		scoreWhite, scoreBlack = "1", "0"
	} else {
		scoreWhite, scoreBlack = "0", "1"
	}
	
	// now, initialize the output with the main contents of the LaTeX body
	output := fmt.Sprintf (`\begin{center}
  {\Large %v (%v)}  
\end{center}

\hrule
\noindent
\WhiteKnightOnWhite %v (%v) \hfill %v\\
\BlackKnightOnWhite %v (%v) \hfill %v
\hrule

\vspace{0.5cm}

\newgame

\mainline{%v}\hfill{\textbf{%v}-\textbf{%v}}

\begin{center}
  \showboard
\end{center}`, event, timeControl, white, whiteELO, date, black, blackELO, ECO, game, scoreWhite, scoreBlack)

	// and return the string computed so far
	return output
}

// GameToLaTeX produces LaTeX code that uses package skak to show the given game
func (game PgnGame) GameToLaTeX () string {

	// justsubstitute values over a standard template
	output := fmt.Sprintf (`\documentclass{article}
\usepackage[utf8]{inputenc}
\usepackage[english]{babel}
\usepackage{mathpazo}
\usepackage{nicefrac}
\usepackage{skak}
\def\hrulefill{\leavevmode\leaders\hrule height 10pt\hfill\kern\z@}
\begin{document}

%v

\end{document}`, game.getLaTeXbody ())

	// and return the string
	return output
}

// GameToLaTeX produces LaTeX code that uses package skak to show all games in a
// given collection
func (games PgnCollection) GameToLaTeX () string {

	// start with the preamble of the document
	output := `\documentclass{article}
\usepackage[utf8]{inputenc}
\usepackage[english]{babel}
\usepackage{mathpazo}
\usepackage{nicefrac}
\usepackage{skak}
\def\hrulefill{\leavevmode\leaders\hrule height 10pt\hfill\kern\z@}
\begin{document}`

	// now, process each game in succession
	for _, game := range games.slice {

		output += fmt.Sprintf(`%v
\clearpage
`, game.getLaTeXbody ())
	}

	// and end the document
	output += `
\end{document}`
	
	// and return the string
	return output
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

	moveNumber := -1              // initialize the move counter to unknown
	color := 0;                          // initialize the color to unknown
	var moveValue string;             // move actually parsed in PGN format
	var emt float64;                                   // elapsed move time
	var comments string;                           // comments of each move
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
func GetGamesFromString (pgn string, verbose bool) (games PgnCollection) {

	// just iterate over the string extracting the information of every game
	for ;reGame.MatchString (pgn); {

		// In case a match has been found, extract the next game
		tag := reGame.FindStringSubmatchIndex (pgn)

		// Parse this game and add it to the slice of games to return
		game := getGameFromString (pgn[tag[0]:tag[1]], verbose)
		games.slice = append (games.slice, game)
		games.nbGames += 1

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
func GetGamesFromFile (pgnfile string, verbose bool) (games PgnCollection) {

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
