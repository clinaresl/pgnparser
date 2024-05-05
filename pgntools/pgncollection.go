/*
  pgncollection.go
  Description: Simple tools for handling collections of games in PGN
               format
  -----------------------------------------------------------------------------

  Started on  <Sat May  9 16:50:49 2015 Carlos Linares Lopez>
  Last update <sábado, 07 mayo 2016 17:26:19 Carlos Linares Lopez (clinares)>
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
	"fmt"
	"io"  // io streams
	"log" // logging services

	"text/template" // go facility for processing templates

	// import my favourite package to automatically create tables
	"github.com/clinaresl/table"
)

// typedefs
// ----------------------------------------------------------------------------

// // PGN games can be sorted either in ascending or descending order. The
// // direction is then defined as an integer
// type sortingDirection int

// // A pgnSorting consists of two items: a constant value for distinguishing
// // ascending from descending order and a variable name used as a key for sorting
// // pgn games
// type pgnSorting struct {
// 	direction sortingDirection
// 	variable  string
// }

// A PgnCollection consists of an arbitrary number of PgnGames along with a
// count of the number of games stored in it ---this is given to check for
// consistency so that the difference between nbGames and len (slice) shall be
// always null.
//
// In addition, a PGN collection contains a sort descriptor which consists of a
// slice of pairs that contain for each variable whether PGN games should be
// sorted in increasing or decreasing order
type PgnCollection struct {
	slice []PgnGame
	// sortDescriptor []pgnSorting
	nbGames int
}

// consts
// ----------------------------------------------------------------------------

// // PGN games can be sorted either in ascending or descending order
// const (
// 	increasing sortingDirection = 1 << iota // increasing order
// 	decreasing                              // decreasing order
// )

// Methods
// ----------------------------------------------------------------------------

// -- Accessors

// the following are getters over the attributes of a PgnCollection

// Return all games as instances of PgnGame that are stored in this particular
// collection
func (games *PgnCollection) GetGames() []PgnGame {
	return games.slice
}

// Return the index-th game stored in this particular collection. The result is
// undefined in case an index out of bounds is passed
func (games *PgnCollection) GetGame(index int) PgnGame {
	return games.slice[index]
}

// Return the number of items in the collection
func (games PgnCollection) Len() int {
	return games.nbGames
}

// Methods
// ----------------------------------------------------------------------------

// Return an empty collection of PGN games
func NewPgnCollection() PgnCollection {
	return PgnCollection{}
}

// Add the given PgnGame to this collection
func (c *PgnCollection) Add(game PgnGame) {

	// Add this game to the slice of games and increment the counter
	c.slice = append(c.slice, game)
	c.nbGames += 1
}

// Play this collection of games on the given writer showing the board
// repeteadly after the given number of plies on the specified writer, in case
// it is strictly positive.
func (c PgnCollection) Play(plies int, writer io.Writer) {

	// the table has to be shown if an only if plies is greater than zero
	showBoard := (plies > 0)

	// if the board is not going to be shown, make the number of plies to be
	// equal to one so that the game is played anyway
	if !showBoard {
		plies = 1
	}

	// use tables to show the execution of chess games
	tab, _ := table.NewTable(" l c", "cc")
	if showBoard {
		tab.AddThickRule()
	}

	// For each game
	for _, igame := range c.slice {

		// Only in case the board is to be shown, create a table, otherwise,
		// skip the process
		if showBoard {

			// Create a nested table to show the tags of this game
			tab_tags, _ := table.NewTable(" l : l")
			for name, value := range igame.Tags() {
				tab_tags.AddRow(name, value)
			}

			// The tags are shown in a single column containing the table of tags
			// centered
			tab.AddRow(table.Multicolumn(2, "c", tab_tags))
			tab.AddSingleRule()
		}

		// Create a new board and access the list of moves to show
		board := NewPgnBoard()
		imoves := igame.Moves()

		// and now show the requested number of plies along with the resulting
		// chess board
		idx := 0
		for idx < len(imoves)/plies {

			// compute the resulting board
			for jdx := idx * plies; jdx < (idx+1)*plies; jdx += 1 {
				board.UpdateBoard(imoves[jdx])
			}

			if showBoard {

				// add a new row with the list of moves in vertical mode and the
				// updated board
				tab.AddRow(igame.prettyMoves((idx*plies), (idx+1)*plies), board)
				tab.AddRow()
			}

			// and move forward
			idx += 1
		}

		// in case there are still moves to show
		if idx*plies < len(imoves) {

			// update the board with those additional moves
			for jdx := idx * plies; jdx < len(imoves); jdx += 1 {
				board.UpdateBoard(imoves[jdx])
			}

			if showBoard {
				// and add the last row
				tab.AddRow(igame.prettyMoves(idx*plies, len(imoves)), board)
			}
		}

		if showBoard {
			// and add a separator with the next game
			tab.AddThickRule()
		}
	}

	// and write the result of the execution in the given writer only in case it
	// has been requested
	if showBoard {
		io.WriteString(writer, fmt.Sprintf("%v\n", tab))
	}
}

// Create a brand new PgnCollection with games found in this collection which
// satisfy the given expression
func (c PgnCollection) Filter(expression string) (*PgnCollection, error) {

	// Create an empty collection of chess games
	collection := NewPgnCollection()

	// Process each game in this collection
	for _, igame := range c.slice {

		// In case this game satisfies the given query, then add it to the
		// filtered collection
		if result, err := igame.Filter(expression); err != nil {
			return nil, err
		} else {
			if result {
				collection.Add(igame)
			}
		}
	}

	// and return the collection processed so far without errors
	return &collection, nil
}

// Write all games in this collection in the specified io.Writer in PGN format.
// In case it was not possible it returns an error and nil otherwise
func (c PgnCollection) GetPGN(writer io.Writer) error {

	// get the contents of each game in PGN format
	for _, igame := range c.slice {
		if _, err := io.WriteString(writer, igame.GetPGN()); err != nil {
			return err
		}
	}

	// at this point return no error
	return nil
}

// Return a histogram defined with the given specification criteria computed
// over all games in this collection. It returns any error found or nil in case
// the histogram was successfully computed
func (c PgnCollection) GetHistogram(spec string) (*PgnHistogram, error) {

	// Create a new GetHistogram
	histogram, err := NewPgnHistogram(spec)
	if err != nil {
		return nil, err
	}

	// and update the histogram with the information of all games in this
	// collection
	for _, igame := range c.slice {
		if err := histogram.Add(igame); err != nil {
			return nil, err
		}
	}

	// and return the histogram computed so far
	return histogram, nil
}

// Templates
//
// All the following methods are used to handle templates both for generating
// ascii and LaTeX output
// ----------------------------------------------------------------------------

// This function is used in text/templates and it is the equivalent to the
// homonym function defined for PgnGame.
//
// It returns the empty string if the given name is not defined as a tag in the
// first game of this collection and the value of the tag (casted to a string)
// otherwise. This function assumes that all games within the same collection
// share the same tags.
func (games *PgnCollection) GetTagValue(name string) string {

	// first, attempt at reading the specified tag from the first game of
	// the collection
	val, err := games.slice[0].GetTagValue(name)

	// if such tag does not exist, then return the empty string
	if err != nil {
		return ""
	}

	// otherwise, return its value
	return val
}

// This is an auxiliary function used in text/templates to generate slices of
// strings to be given as argument to other methods
func (games *PgnCollection) GetSlice(fields ...any) []any {
	return fields
}

// Returns a table according to the specification given in first place. Columns
// are populated with the tags given in fields. It is intended to be used in
// ascii table templates
func (games *PgnCollection) GetTable(specline string, fields []any) table.Table {

	// Create a table according to the given specification
	table, err := table.NewTable(specline)
	if err != nil {
		log.Fatal(" Fatal error while constructing the table in PgnCollection.GetTable")
	}

	// Add the header
	table.AddThickRule()
	table.AddRow(fields...)
	table.AddDoubleRule()

	// Now, add a row per game
	for idx, game := range games.slice {

		// show a separator every ten lines to make the table easier to
		// read
		if idx > 0 && idx%10 == 0 {
			table.AddSingleRule()
		}

		// and show here the information from the specified fields for
		// this game
		table.AddRow(game.getFields(fields)...)
	}

	// End the table and return the table as a string
	table.AddThickRule()
	return *table
}

// Writes into the specified writer the result of instantiating the given
// template file with information of all games in this collection. The template
// acknowledges all tags of a pgngame plus others. For a full description, see
// the manual.
func (games *PgnCollection) GamesToWriterFromTemplate(dst io.Writer, templateFile string) {

	// access a template and parse its contents
	template, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatal(err)
	}

	// and now execute the template
	err = template.Execute(dst, games)
	if err != nil {
		log.Fatal(err)
	}
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
