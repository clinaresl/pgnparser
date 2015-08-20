/* 
  tbl.go
  Description: Automated generation of text and LaTeX tables
  ----------------------------------------------------------------------------- 

  Started on  <Mon Aug 17 17:48:55 2015 Carlos Linares Lopez>
  Last update <miércoles, 19 agosto 2015 17:53:18 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

// This package provides means to automatically generating text and LaTeX tables
// from very simple specifications
package tbl

import (
	"errors"		// for raising errors
	"fmt"			// printing services
	"log"			// Fatal messages
	"regexp"		// for processing specification strings
	"strings"		// for repeating characters
)

// global variables
// ----------------------------------------------------------------------------

// A string specification consists of an indication how the text is justified in
// a cell and also about the separators. It is made of pairs separator/column
// specification.
// The separator can be present or not and it can be one among the following
// types:
//    void - no separator
//    |    - a thin bar
//    ||   - a thick bar
// The column is one of the types:
//    c - centered
//    l - left
//    r - right
var reSpecification = regexp.MustCompile (`^(\|\||\|)?(c|l|r)`)

// There is also a specific regexp to recognize separators on their own when
// processing the last one before the whole string specification is exhausted
var reLastSeparator = regexp.MustCompile (`^(\|\||\|)$`)


// typedefs
// ----------------------------------------------------------------------------

// A style specifies the way text is formatted within cells in the same
// column. The legal values are represented as integer constants
type stylet int

// A separator specifies the type of separator (if any) between adjacent
// columns. The legal values are represented as integer constants
type separatort int

// A table consists just of a slice of slices of strings. Each slice is a single
// line to be generated which might consist of either user text or
// separators. Additionally, tables are built from specification strings
// (similar to those used in LaTeX) which result in a slice of styles for every
// column whose width is re-computed every time a new line is added. In
// addition, the specification string can be also used to select among various
// separators
type Tbl struct {
	style []stylet			// style of each column
	separator []separatort		// separator between columns
	width []int			// width of each column
	content [][]string		// contents of each cell indexed by row
}


// constants
// ----------------------------------------------------------------------------

// A style specifies the location of text within a single cell
const (
	LEFT stylet = 1	<< iota		// left justified
	CENTER				// centered
	RIGHT				// right justified
)

// A separator specifies one among various characters to be used for separating
// columns
const (
	VOID separatort = 1 << iota	// no separator
	THIN				// single bar
	THICK				// double bar
)


// Functions
// ----------------------------------------------------------------------------

// The following is a private function that returns the style represented by the
// given string
func getStyle (cmd string) (style stylet) {

	switch cmd {
	case "l":
		style = LEFT
	case "c":
		style = CENTER
	case "r":
		style = RIGHT
	default:
		log.Fatalf (" Unknown style string '%v'\n", cmd)
	}

	return
}

// The following is a private function that returns the separator represented by the
// given string
func getSeparator (cmd string) (separator separatort) {

	switch cmd {
	case "":
		separator = VOID
	case "|":
		separator = THIN
	case "||":
		separator = THICK
	default:
		log.Fatalf (" Unknown separator string '%v'\n", cmd)
	}

	return
}

// Return the character used as separator according to the parameter given
func getSeparatorChr (separator separatort) string {

	var output string
	
	switch separator {
	case VOID:
		output = " "
	case THIN:
		output = "\u2502"
	case THICK:
		output = "\u2503"
	default:
		log.Fatalf (" Unknown separator '%v'\n", separator)
	}

	return output
}

// return a string made of blank characters which automatically adjusts the
// specified contents within a cell with the given width according to the given
// style if they are inserted *before* the contents
//
// This function assumes that the cell consists of a single line
func preBlank (contents string, width int, style stylet) string {

	// first, verify that the length of the contents is less or equal than
	// the width
	if len (contents) > width {
		log.Fatalf (" It is not possible to insert '%v' within a cell with %v positions",
			contents, width)
	}

	var nbspaces int		// number of spaces to insert
	
	// Now, acccording to the given style, compute the number of spaces to
	// insert
	switch style {
	case LEFT:
		nbspaces = 1
	case CENTER:
		nbspaces = 1 + (width - len (contents))/2
	case RIGHT:
		nbspaces = 1 + width - len (contents)
	}

	// and return a string with as many blank characters as computed above
	return strings.Repeat (" ", nbspaces)
}

// return a string made of blank characters which automatically adjusts the
// specified contents within a cell with the given width according to the given
// style if they are inserted *after* the contents
//
// This function assumes that the cell consists of a single line
func postBlank (contents string, width int, style stylet) string {

	// first, verify that the length of the contents is less or equal than
	// the width
	if len (contents) > width {
		log.Fatalf (" It is not possible to insert '%v' within a cell with %v positions",
			contents, width)
	}

	var nbspaces int		// number of spaces to insert
	
	// Now, acccording to the given style, compute the number of spaces to
	// insert
	switch style {
	case LEFT:
		nbspaces = 1 + width - len (contents)
	case CENTER:

		// if extra spaces are required, they are inserted after the
		// text (ie., in this function)
		nbspaces = 1 + (width - len (contents))/2 + (width - len (contents)) % 2
	case RIGHT:
		nbspaces = 1
	}

	// and return a string with as many blank characters as computed above
	return strings.Repeat (" ", nbspaces)
}

// Return a new instance of Tbl from a string specification
func NewTable (cmd string) (table Tbl, err error) {

	// INVARIANT: The number of separators shall be equal to the number of
	// columns plus one
	
	// just simply process the string specification
	for ; reSpecification.MatchString (cmd) ; {

		// get the next item in the specification string
		tag := reSpecification.FindStringSubmatchIndex (cmd)

		// get the information on both the separator and the column
		// specification
		sep, column := "", cmd[tag[4]:tag[5]]
		if tag[2] >= 0 {
			sep = cmd[tag[2]:tag[3]]
		}

		// update the information on the separator and the style
		table.separator = append (table.separator, getSeparator (sep))
		table.style = append (table.style, getStyle (column))

		// and now move forward in the specification string 
		cmd = cmd[tag[1]:]
	}

	// At this point, the specification string might be non-empty. In this
	// case, however, the only allowed content of the specification string
	// is just a last separator. Otherwise, an error is returned along with
	// an empty table
	if (cmd != "") {

		// If this is a legal separator ...
		if reLastSeparator.MatchString (cmd) {

			// ... then process it
			tag := reLastSeparator.FindStringSubmatchIndex (cmd)
			table.separator = append (table.separator, getSeparator (cmd[tag[2]:tag[3]]))

			// and return the table along with no error
			return table, nil
			
		} else {

			// otherwise, the specification string is not empty and
			// it is not recognized as the last separator, so that
			// signal an error
			return Tbl{},
			errors.New (fmt.Sprintf ("Syntax error in a specification string at point '%v'\n", cmd))
		}
	}

	// In case of success, return an instance of a new table and no error
	// after ensuring that no separator is inserted at the end (so that the
	// invariant that the number of separators equals the number of columns
	// plus one is preserved)
	table.separator = append (table.separator, VOID)
	return table, nil
}


// Methods
// ----------------------------------------------------------------------------

// Add a single line to the contents of the current table at the bottom. The
// contents are specified as a slice of strings. In case the number of items is
// less than the number of columns, the row is paddled with empty strings. If
// the number of items in the given slice exceeds the number of columns in this
// table, an error is raised
func (table *Tbl) AddRow (row []string) (err error) {

	// First, verify that the number of items in this row is less or equal
	// than the number of columns in this table
	if len (row) > len (table.style) {
		return errors.New (fmt.Sprintf (" The row '%v' exceeds the number of columns of this table (%v != %v)\n", row, len (row), len (table.style)))
	}

	// Second, verify that this table has a legal specification string with
	// non-empty style and separators
	if len (table.style) == 0 {
		return errors.New (" This table can not accept any contents!")
	}

	// Create a slice with the contents of the next row to be inserted at
	// the bottom of the table
	newRow := row

	// And add empty cells if necessary
	for ; len (newRow) < len (table.style) ; {
		newRow = append (newRow, "")
	}

	// now, update the maximum width of each column
	for idx, value := range (newRow) {

		// First, if no content was ever processed, init the maximum
		// width of this column to the length of this cell
		if idx == len (table.width) {
			table.width = append (table.width, len (value))
		} else {

			// Otherwise, compare the length of this item with the
			// maximum width computed so far
			if len (value) > table.width [idx] {
				table.width [idx] = len (value)
			}
		}
	}
	
	// and insert this row to the table
	table.content = append (table.content, newRow)
	return nil
}

// Return the contents of the current table as a string.
func (table Tbl) String () string {

	var output string
	
	// For every single line
	for _, line := range table.content {

		// and for every column
		for idx, content := range line {
		
			// Show first the separator
			output += getSeparatorChr (table.separator [idx])

			// show the contents of this cell according to the style
			// of this column. This is done in three steps: first, a
			// string with blank characters is inserted before, next
			// the contents of this cell are printed out and
			// finally, a last string made of blanks is inserted
			// again. The first and last strings are used to justify
			// the contents of the text in this cell according to
			// its style and they already take into account the
			// extra space between the contents and the two
			// separators surrounding it
			output += preBlank (content, table.width[idx], table.style[idx])
			output += content
			output += postBlank (content, table.width[idx], table.style[idx])
		}

		// show the last separator and end the current line
		output += getSeparatorChr (table.separator [len (table.separator) - 1])
		output += "\n"
	}

	// and return the string
	return output
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
