/* 
  tbl.go
  Description: Automated generation of text and LaTeX tables
  ----------------------------------------------------------------------------- 

  Started on  <Mon Aug 17 17:48:55 2015 Carlos Linares Lopez>
  Last update <jueves, 27 agosto 2015 02:33:33 Carlos Linares Lopez (clinares)>
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

// the following map relates integer constants to characters to be printed. It
// is initialized in the init function of this module and it is used to print
// cells
var characterSet map[contentType]string

// typedefs
// ----------------------------------------------------------------------------

// Any specific cell of a table can be one among different types: either
// separators or text cells. The legal values are represented as integer
// constants. 
type contentType int

// Cells are the bricks used to build up the tables. They can be either
// separators or text cells and they have their own width. In case the cell
// contains text, its contents are specified separately in a string. Cells
// contain the minimal necessary information to draw themselves. This means that
// they have no idea about their current location or their surrounding
// environment. Thus, they have to be carefully set to the right type of cell.
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
// every line. Finally, there are three flags used to *remember* whether the
// last line was a horizontal rule or not and the type of horizontal rule. This
// is necessary to redraw the connectors in case that more lines are added.
type Tbl struct {
	column []tblColumn
	row []tblLine
	width []int
	horizontalSingleRule bool
	horizontalDoubleRule bool
	horizontalThickRule bool
}


// constants
// ----------------------------------------------------------------------------

// Any specific cell of a table can be one among different types: either
// separators or text cells.
const (

	// generic separators
	VOID contentType = iota			// nothing
	BLANK					// blank character

	// vertical separators
	VERTICAL_SINGLE				// 2502: │
	VERTICAL_DOUBLE				// 2551: ║
	VERTICAL_THICK				// 2503: ┃
	VERTICAL_VERBATIM			// text separator

	// horizontal separators
	HORIZONTAL_SINGLE			// 2500: ─
	HORIZONTAL_DOUBLE			// 2550: ═
	HORIZONTAL_THICK			// 2501: ━

	// horizontal separators with vertical bars
	DOWN_LIGHT_AND_RIGHT_HEAVY		// 250d: ┍
	HEAVY_DOWN_AND_RIGHT			// 250f: ┏
	DOWN_LIGHT_AND_LEFT_HEAVY		// 2511: ┑
	HEAVY_DOWN_AND_LEFT			// 2513: ┓
	UP_LIGHT_AND_RIGHT_HEAVY		// 2515: ┕
	HEAVY_UP_AND_RIGHT			// 2517: ┗
	UP_LIGHT_AND_LEFT_HEAVY			// 2519: ┙
	HEAVY_UP_AND_LEFT			// 251b: ┛
	VERTICAL_LIGHT_AND_RIGHT_HEAVY		// 251d: ┝
	HEAVY_VERTICAL_AND_RIGHT		// 2523: ┣
	VERTICAL_LIGHT_AND_LEFT_HEAVY		// 2525: ┥
	HEAVY_VERTICAL_AND_LEFT			// 252b: ┫
	DOWN_LIGHT_AND_HORIZONTAL_HEAVY		// 252f: ┯
	HEAVY_DOWN_AND_HORIZONTAL		// 2533: ┳
	UP_LIGHT_AND_HORIZONTAL_HEAVY		// 2537: ┷
	HEAVY_UP_AND_HORIZONTAL			// 253b: ┻
	VERTICAL_LIGHT_AND_HORIZONTAL_HEAVY	// 253f: ┿
	HEAVY_VERTICAL_AND_HORIZONTAL		// 254b: ╋
	
	// text cells
	LEFT					// left justified
	CENTER					// centered
	RIGHT					// right justified
)

// Functions
// ----------------------------------------------------------------------------

// Return the information of a column according to the given string. Note that
// columns are specified as cells and thus, their width and text (if applicable)
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

			// if this has been recognized as a legal verbatim
			// column, extract its contents and return them. Note
			// that in this case no extra spaces are added either
			// before or after the separator
			tag := reVerbatimSeparator.FindStringSubmatchIndex (cmd)
			return tblColumn{VERTICAL_VERBATIM,
				utf8.RuneCountInString (cmd[tag[2]:tag[3]]),
				cmd[tag[2]:tag[3]]}
		}

		// otherwise, raise an error
		log.Fatalf (" Unknown column type '%v'\n", cmd)
	}

	return
}

// initializes this module by setting the right values in the characterSet map
func init () {

	// initialize the map of utf-8 characters and set its contents
	characterSet = make (map[contentType]string)

	// -- generic separators
	characterSet[VOID]  = ""
	characterSet[BLANK] = " "
	
	// -- vertical separators
	characterSet[VERTICAL_SINGLE] = "\u2502"
	characterSet[VERTICAL_DOUBLE] = "\u2551"
	characterSet[VERTICAL_THICK]  = "\u2503"

	// -- horizontal separators
	characterSet[HORIZONTAL_SINGLE] = "\u2500"
	characterSet[HORIZONTAL_DOUBLE] = "\u2550"
	characterSet[HORIZONTAL_THICK]  = "\u2501"

	// -- horizontal separators with vertical bars
	characterSet[DOWN_LIGHT_AND_RIGHT_HEAVY]            = "\u250d"
	characterSet[HEAVY_DOWN_AND_RIGHT]                  = "\u250f"
	characterSet[DOWN_LIGHT_AND_LEFT_HEAVY]             = "\u2511"
	characterSet[HEAVY_DOWN_AND_LEFT]                   = "\u2513"
	characterSet[UP_LIGHT_AND_RIGHT_HEAVY]              = "\u2515"
	characterSet[HEAVY_UP_AND_RIGHT]                    = "\u2517"
	characterSet[UP_LIGHT_AND_LEFT_HEAVY]               = "\u2519"
	characterSet[HEAVY_UP_AND_LEFT]                     = "\u251b"
	characterSet[VERTICAL_LIGHT_AND_RIGHT_HEAVY]        = "\u251d"
	characterSet[HEAVY_VERTICAL_AND_RIGHT]              = "\u2523"
	characterSet[VERTICAL_LIGHT_AND_LEFT_HEAVY]         = "\u2525"
	characterSet[HEAVY_VERTICAL_AND_LEFT]               = "\u252b"
	characterSet[DOWN_LIGHT_AND_HORIZONTAL_HEAVY]       = "\u252f"
	characterSet[HEAVY_DOWN_AND_HORIZONTAL]             = "\u2533"
	characterSet[UP_LIGHT_AND_HORIZONTAL_HEAVY]         = "\u2537"
	characterSet[HEAVY_UP_AND_HORIZONTAL]               = "\u253b"
	characterSet[VERTICAL_LIGHT_AND_HORIZONTAL_HEAVY]   = "\u253f"
	characterSet[HEAVY_VERTICAL_AND_HORIZONTAL]         = "\u254b"
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

// Redraw the last line in case it is a horizontal thick rule. This is necessary
// in case more lines are added after a horizontal rule so that the connectors
// are now set properly
func (table *Tbl) redoThickRule () {
	last := len (table.row) - 1
	for idx, column := range table.column {
		switch column.content {
		case VERTICAL_SINGLE:
			if idx==0 {
				if last == 0 {
					table.row[last][idx]=cellType{DOWN_LIGHT_AND_RIGHT_HEAVY,
						column.width, ""}
				} else {
					table.row[last][idx]=cellType{VERTICAL_LIGHT_AND_RIGHT_HEAVY,
						column.width, ""}
				}
			} else if idx == len (table.column) - 1 {
				if last == 0 {
					table.row[last][idx]=cellType{DOWN_LIGHT_AND_LEFT_HEAVY,
						column.width, ""}
				} else {
					table.row[last][idx]=cellType{VERTICAL_LIGHT_AND_LEFT_HEAVY,
						column.width, ""}
				}
			} else {
				if last == 0 {
					table.row[last][idx]=cellType{DOWN_LIGHT_AND_HORIZONTAL_HEAVY,
						column.width, ""}
				} else {
					table.row[last][idx]=cellType{VERTICAL_LIGHT_AND_HORIZONTAL_HEAVY,
						column.width, ""}
				}
			}
		case VERTICAL_DOUBLE, VERTICAL_THICK:

			// note that both cases are dealt with in the same way
			// since there are no UTF-8 characters which combine
			// them
			if idx==0 {
				if last == 0 {
					table.row[last][idx]=cellType{HEAVY_DOWN_AND_RIGHT,
						column.width, ""}
				} else {
					table.row[last][idx] = cellType{HEAVY_VERTICAL_AND_RIGHT,
						column.width, ""}
				}
			} else if idx == len (table.column) - 1 {
				if last == 0 {
					table.row[last][idx] = cellType{HEAVY_DOWN_AND_LEFT,
						column.width, ""}
				} else {
					table.row[last][idx] = cellType{HEAVY_VERTICAL_AND_LEFT,
						column.width, ""}
				}
			} else {
				if last == 0 {
					table.row[last][idx] = cellType{HEAVY_DOWN_AND_HORIZONTAL,
						column.width, ""}
				} else {
					table.row[last][idx] = cellType{HEAVY_VERTICAL_AND_HORIZONTAL,
						column.width, ""}
				}
			}
		}
	}
}

// Add a single line of text to the bottom of the receiver table. The contents
// are specified as a slice of strings. In case the number of items is less than
// the number of columns, the row is paddled with empty strings. If the number
// of items in the given slice exceeds the number of columns in this table, an
// error is raised
func (table *Tbl) AddRow (row []string) (err error) {

	// First of all, in case the last line was a thick rule, redo it since
	// we are about to generate a new line
	if table.horizontalThickRule {
		table.redoThickRule ()
	}
	
	// insert all cells of this line: those provided by the user and others
	// provided in the specification string
	var newRow tblLine
	idx := 0
	for jdx, value := range (table.column) {

		// depending upon the type of this column cell
		switch (value.content) {

		case VOID, BLANK, VERTICAL_SINGLE, VERTICAL_DOUBLE, VERTICAL_THICK, VERTICAL_VERBATIM:

			// in case this cell does not contain text specified by
			// the user, just copy its specification
			newRow = append (newRow, cellType (value))

		case LEFT, CENTER, RIGHT:

			// if it contains text provided by the user, then create
			// a new cell with those contents (along with two
			// surrounding spaces) and move forward in the slice
			// provided by the user

			// if there are no more contents provided by the user,
			// paddle the remainin entries with blank spaces
			content := " "
			if idx < len (row) {
				content = " " + row[idx] + " "
			}

			// compute the width of this cell as the maximum between
			// the current width of the column and the width of the
			// text to insert here (since only one line is used!)
			// and add it to the table with a blank space after as
			// well
			table.width [jdx] = int (math.Max (float64 (table.width[jdx]),
				float64 (utf8.RuneCountInString (content))))
			newRow = append (newRow,
				cellType{value.content,
					table.width [jdx],
					content})
			idx += 1
		default:
			return errors.New (fmt.Sprintf("Unknown column type '%v'\n", value))
		}
	}
	
	// add this row to the table and exit with no error
	table.row = append (table.row, newRow)

	// and make sure to update all flags of the horizontal rules properly
	table.horizontalSingleRule = false
	table.horizontalDoubleRule = false
	table.horizontalThickRule = false
	return nil
}

// Add a thick horizontal rule that intersects with the vertical separators
// provided that any have been specified.
func (table *Tbl) HThickRule () {

	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	if table.horizontalThickRule {
		table.redoThickRule ()
	}
	
	var newRow tblLine
	for idx, column := range table.column {
		switch column.content {
		case VERTICAL_SINGLE:
			if idx==0 {
				newRow = append (newRow,
					cellType {UP_LIGHT_AND_RIGHT_HEAVY,
						table.width[idx], ""})
			} else if idx == len (table.column) - 1 {
				newRow = append (newRow,
					cellType {UP_LIGHT_AND_LEFT_HEAVY,
						table.width[idx], ""})
			} else {
				newRow = append (newRow,
					cellType{UP_LIGHT_AND_HORIZONTAL_HEAVY,
						table.width[idx], ""})
			}

		case VERTICAL_DOUBLE, VERTICAL_THICK:

			// note that both cases are dealt with in the same way
			// since there are no UTF-8 characters which combine
			// them
			if idx==0 {
				newRow = append (newRow,
					cellType {HEAVY_UP_AND_RIGHT,
						table.width[idx], ""})
			} else if idx == len (table.column) - 1 {
				newRow = append (newRow,
					cellType {HEAVY_UP_AND_LEFT,
						table.width[idx], ""})
			} else {
				newRow = append (newRow, cellType{HEAVY_UP_AND_HORIZONTAL,
					table.width[idx], ""})
			}
		default:
			newRow = append (newRow,
				cellType {HORIZONTAL_THICK, 
					table.width[idx], ""})
		}
	}
	table.row = append (table.row, newRow)

	// Before leaving, set the flag of a thick horizontal rule
	table.horizontalSingleRule = false
	table.horizontalDoubleRule = false
	table.horizontalThickRule = true
}

// Add a thick horizontal rule to the current table. Top rules do not draw
// intersections with column separators (they break them instead).
//
// This function is implemented in imitation to the LaTeX package booktabs
func (table *Tbl) TopRule () {

	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	if table.horizontalThickRule {
		table.redoThickRule ()
	}
	
	// Top rules consist of thick lines. Just add a thick line with no text
	// at all in every column of this line
	var newRow tblLine	
	for idx := range table.column {
		newRow = append (newRow, cellType {HORIZONTAL_THICK,
			table.width[idx], ""})
	}
	table.row = append (table.row, newRow)

	// and make sure to update all flags of the horizontal rules properly
	table.horizontalSingleRule = false
	table.horizontalDoubleRule = false
	table.horizontalThickRule = false
	
}

// Add a thin horizontal rule to the current table. Mid rules do not draw
// intersections with column separators (they break them instead)
//
// This function is implemented in imitation to the LaTeX package booktabs
func (table *Tbl) MidRule () {

	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	if table.horizontalThickRule {
		table.redoThickRule ()
	}
	
	// Mid rules consist of thin lines. Just add a thin line with no text at
	// all in every column of this line
	var newRow tblLine	
	for idx := range table.column {
		newRow = append (newRow, cellType {HORIZONTAL_SINGLE,
			table.width[idx], ""})
	}
	table.row = append (table.row, newRow)

	// and make sure to update all flags of the horizontal rules properly
	table.horizontalSingleRule = false
	table.horizontalDoubleRule = false
	table.horizontalThickRule = false	
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
	case VERTICAL_VERBATIM:
		output = cell.text
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
	default:
		output = strings.Repeat (characterSet[cell.content], cell.width)
	}
	return output
}

// A table is drawn just by drawing its cells one after the other
func (table *Tbl) String () string {

	var output string

	// for every single line
	for _, line := range table.row {

		// and for every column
		for jdx, cell := range line {

			// draw this cell after updating its width
			cell.width = table.width [jdx]
			output += fmt.Sprintf ("%v", cell)
		}

		// and start a newline
		output += "\n"
	}
	return output
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
