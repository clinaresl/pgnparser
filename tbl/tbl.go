/* 
  tbl.go
  Description: Automated generation of text and LaTeX tables
  ----------------------------------------------------------------------------- 

  Started on  <Mon Aug 17 17:48:55 2015 Carlos Linares Lopez>
  Last update <lunes, 31 agosto 2015 02:33:29 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

// This package provides means to automatically generating text and LaTeX tables
// from simple string specifications as those used in LaTeX.
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
//    p{width} - creates a column with a fixed width which can be specified in
//    LaTeX format. This package however, takes only the first digits and
//    interpretes them as a number of characters
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
	"strconv"		// Atoi
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
//    |      - a single bar
//    ||     - a double bar
//    |||    - a thick bar
//    @{...} - a verbatim separator where '...' stands for anything but '}'
//    p{width} - creates a column with a fixed width which can be specified in
//    LaTeX format. This package however, takes only the first digits and
//    interpretes them as a number of characters
// 
// The column is one of the types:
//    c - centered
//    l - left
//    r - right
var reSpecification = regexp.MustCompile (`^(p{[^}]*\}|\@\{[^}]*\}|\|\|\||\|\||\||c|l|r)`)

// Verbatim separators are also processed with their own regular expression for
// extracting the text
var reVerbatimSeparator = regexp.MustCompile (`^@\{(?P<text>[^}]*)\}`)

// Likewise, fixed widths are processed separately with an additional regular
// expression to extract the width
var reFixedWidth = regexp.MustCompile (`^p\{(?P<width>[^}]*)\}`)

// While the full width is passed to the LaTeX code, only the integer part is
// used to set the width of a column. Thus, an additional regexp is used just to
// extract it
var reIntegerFixedWidth = regexp.MustCompile (`^(?P<value>[\d]+).*`)


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

		// Still, this might be a legal column only if it is a verbatim,
		// or ...
		if reVerbatimSeparator.MatchString (cmd) {

			// if this has been recognized as a legal verbatim
			// column, extract its contents and return them. Note
			// that in this case no extra spaces are added either
			// before or after the separator
			tag := reVerbatimSeparator.FindStringSubmatchIndex (cmd)
			return tblColumn{VERTICAL_VERBATIM,
				utf8.RuneCountInString (cmd[tag[2]:tag[3]]),
				cmd[tag[2]:tag[3]]}
		} else if reFixedWidth.MatchString (cmd) {

			// ... or a fixed widht column
			tag0 := reFixedWidth.FindStringSubmatchIndex (cmd)
			arg  := cmd[tag0[2]:tag0[3]]
			
			// If this is the case, create a column of the proper
			// type indicating the width provided by the user
			if reIntegerFixedWidth.MatchString (arg) {

				tag1 := reIntegerFixedWidth.FindStringSubmatchIndex (arg)
				width, err := strconv.Atoi (arg[tag1[2]:tag1[3]]); if err != nil {
					log.Fatalf (" Impossible to extract the width from '%v'",
						arg[tag1[2]:tag1[3]])
				} else {
					return tblColumn {VERTICAL_FIXED_WIDTH,
						width, ""}
				}
			} else {
				log.Fatalf (" Impossible to extract an integer width from '%v'",
					arg)
			}
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

	// First of all, in case the last line was a horizontal rule, redo it
	// since we are about to generate a new line
	if table.horizontalSingleRule {
		table.redoSingleRule ()
	} else if table.horizontalDoubleRule {
		table.redoDoubleRule ()
	} else if table.horizontalThickRule {
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

		case VERTICAL_FIXED_WIDTH:

			// if there are no more contents provided by the user,
			// paddle the remainin entries with blank
			// spaces. Otherwise, add the user text
			var text string
			if idx < len (row) {
				content := row[idx]
				if utf8.RuneCountInString (content) <= value.width {
					text = content + strings.Repeat (" ",
						value.width -
							utf8.RuneCountInString (content))
				} else {
					text = content[0:value.width-1] + "►"
				}				
			} else {
				text = strings.Repeat (" ", value.width)
			}

			// finally, make sure user text is surrounded by blank
			// spaces unless the previous/next column are verbatim
			// column
			if jdx == 0 ||
				table.column[jdx-1].content != VERTICAL_VERBATIM {
				text = " " + text
			}
			if jdx == len (table.column) - 1 ||
				table.column[jdx+1].content != VERTICAL_VERBATIM {
				text = text + " "
			}

			// update the width of this column after taking into
			// account the surrounding blank spaces
			table.width [jdx] = int (math.Max (float64 (table.width[jdx]),
				float64 (utf8.RuneCountInString (text))))
			
			// create the cell and add it to this row
			newRow = append (newRow, cellType{value.content,
				utf8.RuneCountInString (text),
				text})
			
			// and move to the next entry provided by the user
			idx += 1
			
		case LEFT, CENTER, RIGHT:

			// if there are no more contents provided by the user,
			// paddle the remainin entries with blank
			// spaces. Otherwise, add the user text
			content := " "
			if idx < len (row) {

				// make sure user text is surrounded by blank
				// spaces unless the previous/next column are
				// verbatim column
				if jdx == 0 ||
					table.column[jdx-1].content != VERTICAL_VERBATIM {
					content = " " + row[idx]
				} else {
					content = row[idx]
				}
				if jdx == len (table.column) - 1 {
					if table.column[jdx].content != VERTICAL_VERBATIM {
						content += " "
					}
				} else if table.column[jdx+1].content != VERTICAL_VERBATIM {
						content += " "
				}
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
			
			// and move to the next entry provided by the user
			idx += 1
		default:
			return errors.New (fmt.Sprintf("Unknown column type '%v'\n", value))
		}
	}

	// Check there are no left values in the given row
	if idx < len (row) {
		return errors.New (fmt.Sprintf ("%v items were given but there are only %v columns", len (row), idx))
	}
	
	// add this row to the table and exit with no error
	table.row = append (table.row, newRow)

	// and make sure to update all flags of the horizontal rules properly
	table.horizontalSingleRule = false
	table.horizontalDoubleRule = false
	table.horizontalThickRule = false
	return nil
}

// Add a single horizontal rule that intersects with the vertical separators
// provided that any have been specified.
func (table *Tbl) HSingleRule () {

	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	if table.horizontalSingleRule {
		table.redoSingleRule ()
	} else if table.horizontalDoubleRule {
		table.redoDoubleRule ()
	} else if table.horizontalThickRule {
		table.redoThickRule ()
	}
	
	var newRow tblLine
	for idx, column := range table.column {
		switch column.content {
		case VERTICAL_SINGLE:
			if idx==0 {
				newRow = append (newRow,
					cellType {LIGHT_UP_AND_RIGHT,
						table.width[idx], ""})
			} else if idx == len (table.column) - 1 {
				newRow = append (newRow,
					cellType {LIGHT_UP_AND_LEFT,
						table.width[idx], ""})
			} else {
				newRow = append (newRow,
					cellType{LIGHT_UP_AND_HORIZONTAL,
						table.width[idx], ""})
			}

		case VERTICAL_DOUBLE, VERTICAL_THICK:

			// note that both cases are dealt with in the same way
			// since there are no UTF-8 characters which combine
			// them
			if idx==0 {
				newRow = append (newRow,
					cellType {UP_HEAVY_AND_RIGHT_LIGHT,
						table.width[idx], ""})
			} else if idx == len (table.column) - 1 {
				newRow = append (newRow,
					cellType {UP_HEAVY_AND_LEFT_LIGHT,
						table.width[idx], ""})
			} else {
				newRow = append (newRow, cellType{UP_HEAVY_AND_HORIZONTAL_LIGHT,
					table.width[idx], ""})
			}
		default:
			newRow = append (newRow,
				cellType {HORIZONTAL_SINGLE, 
					table.width[idx], ""})
		}
	}
	table.row = append (table.row, newRow)

	// Before leaving, set the flag of a thick horizontal rule
	table.horizontalSingleRule = true
	table.horizontalDoubleRule = false
	table.horizontalThickRule = false
}

// Add a double horizontal rule that intersects with the vertical separators
// provided that any have been specified.
func (table *Tbl) HDoubleRule () {

	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	if table.horizontalSingleRule {
		table.redoSingleRule ()
	} else if table.horizontalDoubleRule {
		table.redoDoubleRule ()
	} else if table.horizontalThickRule {
		table.redoThickRule ()
	}
	
	var newRow tblLine
	for idx, column := range table.column {
		switch column.content {
		case VERTICAL_SINGLE:
			if idx==0 {
				newRow = append (newRow,
					cellType {UP_SINGLE_AND_RIGHT_DOUBLE,
						table.width[idx], ""})
			} else if idx == len (table.column) - 1 {
				newRow = append (newRow,
					cellType {UP_SINGLE_AND_LEFT_DOUBLE,
						table.width[idx], ""})
			} else {
				newRow = append (newRow,
					cellType{UP_SINGLE_AND_HORIZONTAL_DOUBLE,
						table.width[idx], ""})
			}

		case VERTICAL_DOUBLE, VERTICAL_THICK:

			// note that both cases are dealt with in the same way
			// since there are no UTF-8 characters which combine
			// them
			if idx==0 {
				newRow = append (newRow,
					cellType {DOUBLE_UP_AND_RIGHT,
						table.width[idx], ""})
			} else if idx == len (table.column) - 1 {
				newRow = append (newRow,
					cellType {DOUBLE_UP_AND_LEFT,
						table.width[idx], ""})
			} else {
				newRow = append (newRow, cellType{DOUBLE_UP_AND_HORIZONTAL,
					table.width[idx], ""})
			}
		default:
			newRow = append (newRow,
				cellType {HORIZONTAL_DOUBLE, 
					table.width[idx], ""})
		}
	}
	table.row = append (table.row, newRow)

	// Before leaving, set the flag of a thick horizontal rule
	table.horizontalSingleRule = false
	table.horizontalDoubleRule = true
	table.horizontalThickRule = false
}

// Add a thick horizontal rule that intersects with the vertical separators
// provided that any have been specified.
func (table *Tbl) HThickRule () {

	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	if table.horizontalSingleRule {
		table.redoSingleRule ()
	} else if table.horizontalDoubleRule {
		table.redoDoubleRule ()
	} else if table.horizontalThickRule {
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
	if table.horizontalSingleRule {
		table.redoSingleRule ()
	} else if table.horizontalDoubleRule {
		table.redoDoubleRule ()
	} else if table.horizontalThickRule {
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
	if table.horizontalSingleRule {
		table.redoSingleRule ()
	} else if table.horizontalDoubleRule {
		table.redoDoubleRule ()
	} else if table.horizontalThickRule {
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
	case VERTICAL_FIXED_WIDTH:
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

// Return a string with a representation of the contents of the table
func (table Tbl) String () string {

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
