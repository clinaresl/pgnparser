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
	"path"
	"regexp"
	"sort"
	"text/template"

	// go facility for processing templates

	// import my favourite package to automatically create tables

	"github.com/clinaresl/pgnparser/metatemplate"
	"github.com/clinaresl/table"
	// and also the replacement of text/templates known as multi-template
)

// typedefs
// ----------------------------------------------------------------------------

// PGN games can be sorted either in ascending or descending order. The
// direction is then defined as an integer
type sortingDirection int

// A pgnSorting consists of two items: a constant value for distinguishing
// ascending from descending order and a criteria (either a variable or a bool
// expression) which is used for sorting elements
type pgnSorting struct {
	direction sortingDirection
	criteria  string
}

// So that a sorting criteria consists of a sequence of pgnSorting pairs
type criteriaSorting []pgnSorting

// A PgnCollection consists of an arbitrary number of PgnGames
type PgnCollection struct {
	slice   []PgnGame
	nbGames int
}

// consts
// ----------------------------------------------------------------------------

// PGN games can be sorted either in ascending or descending order
const (
	increasing sortingDirection = 1 << iota // increasing order
	decreasing                              // decreasing order
)

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

// Sort the games in this collection according to the specific criteria which
// consists of a semicolon separated list of pairs (direction var/bool expr).
// The direction can be either '<' (ascending order) or '>' (descending order),
// next either a variable or a bool expression can be used so that games are
// sorted according to the value of the variable or the result of the evaluation
// of the bool expr
//
// The result is returned in a brand new collection of Pgn games
func (c *PgnCollection) Sort(spec string) (*PgnCollection, error) {

	// parse the given specification string. First, distinguish the different
	// parts and get the sorting direction and criteria (either a variable or a
	// bool expression) of each one
	cmds := reCriteria.Split(spec, -1)
	if len(cmds) == 0 {
		return nil, fmt.Errorf(" Empty sorting string '%v'\n", spec)
	}

	// Process all chunks to get a sorting criteria to be used for sorting games
	criteria := make(criteriaSorting, 0)
	for _, icmd := range cmds {

		// Next, process this specific chunk
		if match, err := regexp.MatchString(reSorting, icmd); err != nil {
			return nil, err
		} else {

			// In case no match is detected then return an error
			if !match {
				return nil, fmt.Errorf(" Syntax eerror in sorting command '%v'\n", icmd)
			} else {

				// Extract the groups
				indices := regexp.MustCompile(reSorting).FindSubmatchIndex([]byte(icmd))

				// Get the direction and the variable/bool expression
				var sortingDirection = increasing
				if icmd[indices[2]:indices[3]] == ">" {
					sortingDirection = decreasing
				}

				// Create a sorting criteria and add it to the slice of sorting
				// criteria to be used for sorting games
				criteria = append(criteria,
					pgnSorting{
						direction: sortingDirection,
						criteria:  icmd[indices[4]:indices[5]],
					})
			}
		}
	}

	// Now, sort the slice of games in this collection
	sort.SliceStable(c.slice, func(i, j int) bool {
		result, err := c.GetGame(i).lessGame(c.GetGame(j), criteria)
		if err != nil {
			log.Fatalf(" Error while sorting games: '%v'\n", err)
		}
		return result
	})

	return c, nil
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
//
// It is intended to be used in LaTeX templates
func (games *PgnCollection) GetField(name string) string {

	// first, attempt at reading the specified tag from the first game of the
	// collection. Note that because GetField is intended to be used in
	// templates, no error is captured (and a Fatalf is issued in such case)
	val := games.slice[0].GetField(name)

	// otherwise, return its value
	return val
}

// Returns a table according to the specification given in first place. Columns
// are populated with the tags given in fields. It is intended to be used in
// ascii table templates
//
// It is intended to be used in LaTeX templates
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

	// create a dictionary of meta-variables
	variables := make(map[string]string)

	// access a template and parse its contents
	tpl, err := metatemplate.New(path.Base(templateFile)).Funcs(template.FuncMap{
		"getSlice": func(fields ...interface{}) []interface{} {
			return fields
		},
	}).ParseFiles(variables, templateFile)

	if err != nil {
		log.Fatal(err)
	}

	// and now execute the template
	err = tpl.ExecuteTemplate(dst, tpl.Name(), games)
	if err != nil {
		log.Fatal(err)
	}
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
