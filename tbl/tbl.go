/* 
  tbl.go
  Description: Automated generation of text and LaTeX tables
  ----------------------------------------------------------------------------- 

  Started on  <Mon Aug 17 17:48:55 2015 Carlos Linares Lopez>
  Last update <lunes, 24 agosto 2015 13:43:59 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

// This package provides means to automatically generating text and LaTeX tables
// from simple string specifications as those used in LaTeX
//
// A string specification consists of an indication how the text is justified in
// all cells in the same column and also how separators shall be
// formatted. Thus, they can refer to either text cells or column separators.
// 
// A separator can be present or not and it can be one among the following
// types:
//    void   - no separator
//    |      - a single bar
//    ||     - a double bar
//    |||    - a thick bar
//    @{...} - a verbatim separator where '...' stands for anything but '}'
// 
// The column is one of the types:
//    c - centered
//    l - left
//    r - right
//
// Text cells and separators can be provided in any order. The tbl package fully
// supports utf-8 characters
package tbl

import (
	"errors"		// for raising errors
	"fmt"			// printing services
	"log"			// Fatal messages
	"math"			// Max
	"regexp"		// for processing specification strings
	"strings"		// for repeating characters
	"unicode/utf8"		// provides support for UTF-8 characters
)

// global variables
// ----------------------------------------------------------------------------

// A string specification consists of an indication how the text is justified in
// a cell and also about the separators. It is made of pairs separator/column
// specification.
// 
// The separator can be present or not and it can be one among the following
// types:
//    void   - no separator
//    |      - a single bar
//    ||     - a double bar
//    |||    - a thick bar
//    @{...} - a verbatim separator where '...' stands for anything but '}'
// 
// The column is one of the types:
//    c - centered
//    l - left
//    r - right
var reSpecification = regexp.MustCompile (`^(@\{[^}]*\}|\|\|\||\|\||\||c|l|r)`)

// Verbatim separators are also processed with their own regular expression for
// extracting the text
var reVerbatimSeparator = regexp.MustCompile (`^@\{(?P<text>[^}]*)\}`)

// typedefs
// ----------------------------------------------------------------------------

// Any specific cell of a table can be one among different types: either
// separators or text cells. The legal values are represented as integer
// constants
type contentType int

// Cells are the bricks used to build up the tables. They can be either
// separators or text cells and they have their own width. In case the cell
// contains text, its contents are specified separately in a string
type cellType struct {
	content contentType
	width int
	text string
}

// A column is specified with a cell type
type tblColumn cellType

// A line is just made up of cells
type tblLine []cellType

// A table consists mainly of two components: information about the columns and
// information about the rows. The former is stored as a slice of columns. Rows
// are specified as a slice of lines, each one with its own cells. Additionally,
// a table contains a slice of widths with the overall width of each cell in
// every line
type Tbl struct {
	column []tblColumn
	row []tblLine
	width []int
}


// constants
// ----------------------------------------------------------------------------

// Any specific cell of a table can be one among different types: either
// separators or text cells.
const (

	// generic separators
	VOID contentType = 1 << iota	// nothing
	BLANK				// blank character

	// vertical separators
	VERTICAL_SINGLE			// single bar
	VERTICAL_DOUBLE			// double bar
	VERTICAL_THICK			// thick bar
	VERTICAL_VERBATIM		// text separator

	// horizontal separators
	HORIZONTAL_SINGLE		// single bar
	HORIZONTAL_DOUBLE		// double bar
	HORIZONTAL_THICK		// thick bar
	
	// text cells
	LEFT				// left justified
	CENTER				// centered
	RIGHT				// right justified
)

// Functions
// ----------------------------------------------------------------------------

// Return the information of a column according to the given string. Note that
// columns are specified as cell and thus, the width and text (if applicable)
// shall be computed as well.
func getColumnType (cmd string) (column tblColumn) {

	switch cmd {
	case "":
		column = tblColumn{VOID, 0, ""}
	case " ":
		column = tblColumn{BLANK, 1, ""}
	case "|":
		column = tblColumn{VERTICAL_SINGLE, 1, ""}
	case "||":
		column = tblColumn{VERTICAL_DOUBLE, 1, ""}
	case "|||":
		column = tblColumn{VERTICAL_THICK, 1, ""}
	case "l":
		column = tblColumn{LEFT, 0, ""}
	case "c":
		column = tblColumn{CENTER, 0, ""}
	case "r":
		column = tblColumn{RIGHT, 0, ""}
	default:

		// Still, this might be a legal column only if it is a verbatim
		if reVerbatimSeparator.MatchString (cmd) {

			// if this has been recognized as a legal verbatim column,
			// extract its contents and return them
			tag := reVerbatimSeparator.FindStringSubmatchIndex (cmd)
			content := cmd[tag[2]:tag[3]]
			width := utf8.RuneCountInString (content)
			return tblColumn{VERTICAL_VERBATIM, width, content}
		}

		// otherwise, raise an error
		log.Fatalf (" Unknown column type '%v'\n", cmd)
	}

	return
}

// Return a new instance of Tbl from a string specification
func NewTable (cmd string) (table Tbl, err error) {

	// just simply process the string specification
	for ; reSpecification.MatchString (cmd) ; {

		// get the next item in the specification string and add it to
		// the collection of columns of this table
		tag := reSpecification.FindStringSubmatchIndex (cmd)
		column := getColumnType (cmd[tag[2]:tag[3]])
		table.column = append (table.column, column)

		// Initialize also the width of each cell
		table.width = append (table.width, column.width)
		
		// move forward
		cmd = cmd[tag[1]:]
	}

	// In case the specification string has not been fully processed, a
	// syntax error has been found
	if (cmd != "") {

		return Tbl{},
		errors.New (fmt.Sprintf ("Syntax error in a specification string at '%v'\n", cmd))
	}

	// otherwise, return with the table and no error
	return
}


// Methods
// ----------------------------------------------------------------------------

// Add a single line of text to the bottom of the receiver table. The contents
// are specified as a slice of strings. In case the number of items is less than
// the number of columns, the row is paddled with empty strings. If the number
// of items in the given slice exceeds the number of columns in this table, an
// error is raised
func (table *Tbl) AddRow (row []string) (err error) {

	// insert all cells of this line: those provided by the user and others
	// provided in the specification string
	var newRow tblLine
	idx := 0
	for jdx, value := range (table.column) {

		// depending upon the type of this column cell
		switch (value.content) {

		case VOID, BLANK, VERTICAL_SINGLE, VERTICAL_DOUBLE, VERTICAL_THICK, VERTICAL_VERBATIM:

			// in case this cell does not contain text specified by
			// the user, just copy it its specification
			newRow = append (newRow, cellType (value))
			
		case LEFT, CENTER, RIGHT:

			// if it contains text provided by the user, then create
			// a new cell with those contents and move forward in
			// the slice provided by the user

			// if there are no more contents provided by the user,
			// paddle the remainin entries with blank spaces
			content := " "
			if idx < len (row) {
				content = row[idx]
			}

			// compute the width of this cell as the maximum between
			// the current width of the column and the width of the
			// text to insert here (since only one line is used!)
			table.width [jdx] = int (math.Max (float64 (table.width[jdx]),
				float64 (utf8.RuneCountInString (content))))
			newRow = append (newRow, cellType{value.content,
				table.width [jdx],
				content})
			idx += 1
		default:
			return errors.New (fmt.Sprintf("Unknown column type '%v'\n", value))
		}
	}
	
	// add this row to the table and exit with no error
	table.row = append (table.row, newRow)
	return nil
}

// Add a thick horizontal rule to the current table. Top rules do not draw
// intersections with column separators (they break them instead).
//
// This function is implemented in imitation to the LaTeX package booktabs
func (table *Tbl) TopRule () {

	// Top rules consist of thick lines. Just add a thick line with no text
	// at all in every column of this line
	var newRow tblLine	
	for idx := range table.column {
		newRow = append (newRow, cellType {HORIZONTAL_THICK,
			table.width[idx], ""})
	}
	table.row = append (table.row, newRow)
}

// Add a thin horizontal rule to the current table. Mid rules do not draw
// intersections with column separators (they break them instead)
//
// This function is implemented in imitation to the LaTeX package booktabs
func (table *Tbl) MidRule () {

	// Top rules consist of thick lines. Just add a thick line with no text
	// at all in every column of this line
	var newRow tblLine
	for idx := range table.column {
		newRow = append (newRow, cellType {HORIZONTAL_SINGLE,
			table.width[idx], ""})
	}
	table.row = append (table.row, newRow)
}

// Add a thick horizontal rule to the current table. Bottom rules do not draw
// intersections with column separators (they break them instead)
//
// This function is implemented in imitation to the LaTeX package booktabs
func (table *Tbl) BottomRule () {

	table.TopRule ()
}

// Cells draw themselves producing a string which takes into account the width
// of the cell.
func (cell cellType) String () string {

	var output string
	
	// depending upon the type of cell
	switch cell.content {
	case VOID:
		output = ""
	case BLANK:
		output = strings.Repeat(" ", cell.width)
	case VERTICAL_SINGLE:
		output = strings.Repeat("\u2502", cell.width)
	case VERTICAL_DOUBLE:
		output = strings.Repeat("\u2551", cell.width)
	case VERTICAL_THICK:
		output = strings.Repeat ("\u2503", cell.width)
	case VERTICAL_VERBATIM:
		output = cell.text
	case HORIZONTAL_SINGLE:
		output = strings.Repeat("\u2500", cell.width)
	case HORIZONTAL_DOUBLE:
		output = strings.Repeat("\u2550", cell.width)
	case HORIZONTAL_THICK:
		output = strings.Repeat("\u2501", cell.width)
	case LEFT:
		output = cell.text + strings.Repeat (" ",
			cell.width - utf8.RuneCountInString (cell.text))
	case CENTER:
		output = strings.Repeat (" ", (cell.width - utf8.RuneCountInString (cell.text))/2) +
			cell.text +
			strings.Repeat (" ", (cell.width - utf8.RuneCountInString (cell.text))/2 +
			(cell.width - utf8.RuneCountInString (cell.text)) % 2)
	case RIGHT:
		output = strings.Repeat (" ", cell.width - utf8.RuneCountInString (cell.text)) +
			cell.text
	}
	return output
}

// A table is drawn just by drawing its cells one after the other
func (table Tbl) String () string {

	var output string
	
	// for every single line
	for _, line := range table.row {

		// and for every column
		for jdx, cell := range line {

			// draw this cell after updating its width
			cell.width = table.width [jdx]
			output += fmt.Sprintf ("%v", cell)

			// add a blank space between adjacent cells unless this
			// column or the next one are verbatim cells or this is
			// the last cell in the row
			if jdx < len (table.column) -1 &&
				table.column[jdx].content != VERTICAL_VERBATIM &&
				(jdx==len (table.column) - 1 ||
				table.column[1+jdx].content != VERTICAL_VERBATIM) {

				// in this case, add a horizontal rule if this
				// cell is of that kind or a blank space
				// otherwise
				if cell.content == HORIZONTAL_SINGLE {
					output += "\u2500"
				} else if cell.content == HORIZONTAL_DOUBLE {
					output += "\u2550"
				} else if cell.content == HORIZONTAL_THICK {
					output += "\u2501"
				} else {
					output += " "
				}
			}
		}
		output += "\n"
	}
	return output
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
