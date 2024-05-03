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
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/clinaresl/table"
)

// constants
// ----------------------------------------------------------------------------
// MAXLEN is the size of the largest block read at once when reading the
// contents of PgnFiles. By default, 1Kbyte
var MAXLEN int32 = 1024

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

	//
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
			game := getGameFromString(text[tag[0]:tag[1]], true)

			// parse all moves and ensure the transcription is correct so that
			// the execution is not played ---and this is achieved by providing
			// a huge number of plies to ParseMoves
			game.ParseMoves(-1)

			// add this game to the collection of games to return
			games = append(games, game)

			// reset the text containing the game just found
			text = ""
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Once done return an instance of PgCollection with all these games
	return &PgnCollection{
		slice:          games,
		sortDescriptor: nil,
		nbGames:        len(games),
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
