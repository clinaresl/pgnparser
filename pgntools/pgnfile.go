// -*- coding: utf-8 -*-
// pgnfile.go
// -----------------------------------------------------------------------------
//
// Started on <jue 02-05-2024 20:25:11.023347448 (1714674311)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package pgntools

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/clinaresl/table"
)

// typedefs
// ----------------------------------------------------------------------------

// A PgnFile contains a collection of chess games in PGN format. It stores no
// information related to the chess games contained in it and it should be used
// solely for creating a PgnCollection
type PgnFile struct {
	name    string    // filename
	size    int64     // size of the file
	modtime time.Time // Last modification time
}

// functions
// ----------------------------------------------------------------------------

// it returns an absolute path of the path given in dirin. It deals with strings
// starting with the symbol '~' and cleans the result (see path.Clean)
func processDirectory(dirin string) (dirout string) {

	// initially, make the dirout to be equal to the dirin
	dirout = dirin

	// first, in case the input directory starts with the symbol
	// '~'
	if dirin[0] == '~' {

		// substitute '~' with the value of the $HOME variable
		dirout = path.Join(os.Getenv("HOME"), dirin[1:])
	}

	// finally, clean the given directory specification
	return path.Clean(dirout)
}

// returns true if the given string names a regular file (ie., that no mode bits
// are set) and false otherwise (thus, it is much like os.IsRegular but it works
// from strings directly). It also returns the fileinfo in case the file exists
func isRegular(path string) (isregular bool, fileinfo os.FileInfo) {

	var err error

	// stat the specified path
	if fileinfo, err = os.Lstat(path); err != nil {
		return false, nil
	}

	// return now whether this is a regular file or not
	return fileinfo.Mode().IsRegular(), fileinfo
}

// Return true if the given filename exists and false otherwise
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Return a slice with all tags in the given string. No error can be returned
// because the string given to this function has already matched the regular
// expression for tags
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
// acknowledged and their information is added to the slice of PgnMove.
//
// Even if the string given in pgn has already matched a regular expression
// other errors might be found and thus an error is returned which can be empty
// if all moves could be extracted. In case of an error, the slice in moves
// returns all moves processed so far
func getMoves(pgn string) (moves []PgnMove, err error) {

	moveNumber := -1     // initialize the move counter to unknown
	color := 0           // initialize the color to unknown
	var moveValue string // move actually parsed in PGN format
	var emt float64      // elapsed move time
	var comments string  // comments of each move

	// process plies in sequence until the whole string is exhausted
	for len(pgn) > 0 {

		// get the next move
		tag := reGroupMoves.FindStringSubmatchIndex(pgn)

		// reGroupMoves contains three groups and therefore legal matches
		// contain 8 characters
		if len(tag) >= 8 {

			// if a move number and color (. or ...) specifier has been found,
			// then process all groups in this matching
			if tag[2] >= 0 && tag[4] >= 0 {

				// update the move counter
				moveNumber, err = strconv.Atoi(pgn[tag[2]:tag[3]])
				if err != nil {
					return moves, errors.New(" Error while extracting the move number")
				}

				// and the color, in case only one character ('.') is found,
				// this is white's move, otherwise, it is black's move
				if tag[5]-tag[4] == 1 {
					color = 1
				} else {
					color = -1
				}
			} else {

				// otherwise, assume that this is the opponent's move
				color *= -1
			}

			// and in any case extract the move value
			moveValue = pgn[tag[6]:tag[7]]
		}

		// and move forward
		pgn = pgn[tag[1]:]

		// are there any comments immediately after? The following loop aims at
		// processing an arbitrary number of comments
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
					return moves, errors.New(" Error while converting emt")
				}
			} else {
				// if not, then just add these comments. In case some comments
				// were already written, make sure to add this in a new line
				if len(comments) > 0 {
					comments += "\r\n"
				}
				comments += pgn[1+tag[2] : tag[3]-1]
			}
			pgn = pgn[tag[1]:]
		}

		// and add this move to the list of moves to return unless there are
		// unknown fields
		if moveNumber == -1 || color == 0 {
			return moves, errors.New(" Either the move number or the color were incorrect")
		}
		moves = append(moves, PgnMove{moveNumber, color, moveValue, float32(emt), comments})
	}

	return
}

// Return an instance of PgnOutcome with the score of each player as specified
// in the given string.
//
// Even if the string given in pgn has already matched a regular expression
// other errors might be found and thus an error is returned which can be empty
// if the outcome could be processed correctly
func getOutcome(pgn string) (outcome *PgnOutcome, err error) {

	// get information about the outcome as given in pgn
	tag := reGroupOutcome.FindStringSubmatchIndex(pgn)

	// process this tag in case it contains 6 indices: <begin/end>-string,
	// <begin/end>-scorewhite, <begin/end>-scoreblack
	if len(tag) >= 6 {

		// if the first tag is three characters long, then this is a
		// draw
		if tag[3]-tag[2] == 3 {
			outcome = &PgnOutcome{0.5, 0.5}
		} else {

			// otherwise, one side won the match
			scoreWhite, err := strconv.Atoi(pgn[tag[2]:tag[3]])
			if err != nil {
				return nil, fmt.Errorf(" Illegal outcome found in string '%s'", pgn)
			}
			outcome = &PgnOutcome{float32(scoreWhite), 1.0 - float32(scoreWhite)}
		}
	} else {

		// In case the grouped regex did not match the given string then the
		// outcome is most likely equal to '*' because 'pgn' already matched the
		// (ungrouped) regexp for the outcome and '*' is not considered in the
		// grouped regexp
		if pgn != "*" {
			return nil, fmt.Errorf(" Unknown outcome found '%v'", pgn)
		} else {

			// In that case the outcome is registered as -1, -1
			outcome = &PgnOutcome{-1, -1}
		}
	}
	return
}

// Return the contents of a chess game from the full transcription of a chess
// game given as a string in PGN format. In case it was not possible to process
// the string, or the information in the game is incorrect (i.e., it could not
// be executed on a chess board) an error is returned
func getGameFromString(pgn string) (*PgnGame, error) {

	// create variables to store different sections of a single PGN game
	var strTags, strMoves, strOutcome string

	// find the tags of the first game in pgn
	endpoints := reTags.FindStringIndex(pgn)
	if endpoints == nil {
		return nil, fmt.Errorf(" No tags were found in the chunk: %v", pgn)
	} else {

		// copy the section of the tags and move forward in the pgn string
		strTags = pgn[endpoints[0]:endpoints[1]]
		pgn = pgn[endpoints[1]:]

		// now, check that this is followed by a legal transcription of chess
		// moves in PGN format
		endpoints = reMoves.FindStringIndex(pgn)
		if endpoints == nil {
			return nil, fmt.Errorf(" No transcription of legal moves were found in the chunk: %v", pgn)
		} else {

			// copy the section with the chess moves and move forward in the pgn
			// string
			strMoves = pgn[endpoints[0]:endpoints[1]]
			pgn = pgn[endpoints[1]:]

			// now, check that the final result is properly written
			endpoints = reOutcome.FindStringIndex(pgn)
			if endpoints == nil {
				return nil, fmt.Errorf(" No lega transcription of the final result was found in the chunk: %v", pgn)
			} else {

				// again, copy the section with the final
				// outcome and move forward in the pgn file
				strOutcome = pgn[endpoints[0]:endpoints[1]]
				pgn = pgn[endpoints[1]:]
			}
		}
	}

	// now, just process the different chunks extracted previously and store
	// them in the game to return. In case processing any of the different parts
	// produces an error, return it immediately
	moves, errMoves := getMoves(strMoves)
	if errMoves != nil {
		return nil, errMoves
	}
	outcome, errOutcome := getOutcome(strOutcome)
	if errOutcome != nil {
		return nil, errOutcome
	}
	return &PgnGame{
		tags:    getTags(strTags),
		moves:   moves,
		outcome: *outcome,
	}, nil
}

// methods
// ----------------------------------------------------------------------------

// A new instance of PgnFile can be created just by providing the file path
// (which is allowed also to contain the character '~'). In case the file does
// not exist, or it is not a regular file then an error is returned
func NewPgnFile(filepath string) (*PgnFile, error) {

	// Substitute the use of the env var $HOME in case it has been given and
	// determine whether the files exists or not
	fullname := processDirectory(filepath)
	if !fileExists(fullname) {
		return nil, fmt.Errorf(" The file '%v' does not exist", filepath)
	} else {
		// verify this is an ordinary regular file
		regularfile, _ := isRegular(fullname)
		if !regularfile {
			return nil, fmt.Errorf(" The file '%v' is not a regular file", fullname)
		}
	}

	// At this point, the file is known both to exist and to be a regular file.
	// Get information about it
	fileinfo, err := os.Stat(fullname)
	if err != nil {
		return nil, fmt.Errorf(" It was not possible to 'stat' the file '%v'", fullname)
	}

	// and return an instance of PgnFile
	return &PgnFile{
		name:    fullname,
		size:    fileinfo.Size(),
		modtime: fileinfo.ModTime(),
	}, nil
}

// Return the filepath of a PgnFile
func (f PgnFile) Name() string {
	return f.name
}

// Return the size in bytes of the given PgnFile
func (f PgnFile) Size() int64 {
	return f.size
}

// Return the last modification time of the given PgnFile
func (f PgnFile) ModTime() time.Time {
	return f.modtime
}

// Return all games stored in the PgnFile f as a collection of PgnGames
func (f PgnFile) Games() (*PgnCollection, error) {

	// Open the PgnFile
	stream, err := os.OpenFile(f.name, os.O_RDONLY, 0644)
	if err != nil {

		// in case of error, return a nil collection of pgn games and the error
		return nil, err
	}

	// Initialize an empty slice of PgGames to return within a PgnCollection
	games := make([]PgnGame, 0)

	// Next, scan the lines of the input file using a buffered input stream
	var text string
	scanner := bufio.NewScanner(stream)

	// Scanning goes line by line
	for scanner.Scan() {

		// text is accumulated until a whole game is found
		text = text + scanner.Text()
		if reGame.MatchString(text) {

			// In case a match has been found, extract the next game
			tag := reGame.FindStringSubmatchIndex(text)

			// Parse this game and get an instance of PgnGame with the
			// information in it
			game, err := getGameFromString(text[tag[0]:tag[1]])
			if err != nil {
				return nil, err
			}

			// parse all moves and ensure the transcription is correct so that
			// the execution is not played ---and this is achieved by providing
			// a huge number of plies to ParseMoves
			game.ParseMoves(-1)

			// add this game to the collection of games to return
			games = append(games, *game)

			// reset the text containing the game just found
			text = ""
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Once done return an instance of PgCollection with all these games
	return &PgnCollection{
		slice:   games,
		nbGames: len(games),
	}, nil
}

// PgnFile are stringers. They just show the information of a PgnFile using a
// table
func (f PgnFile) String() string {

	// Create a table to show the information nicely
	table, err := table.NewTable(" l: l")
	if err != nil {
		log.Fatal(" Fatal error while constructing the table in PgnFile.String")
	}

	table.AddRow("▶ Name", f.Name())
	table.AddRow("▶ Size", fmt.Sprintf("%d", f.Size())+" bytes")
	table.AddRow("▶ Mod Time", f.ModTime())
	table.AddDoubleRule()

	// print the table and return it as a string
	return fmt.Sprintf("%v", table)
}

// Local Variables:
// mode:go
// fill-column:80
// End:
