/*
  tbl.go
  Description: Automated generation of text and LaTeX tables
  -----------------------------------------------------------------------------

  Started on  <Mon Aug 17 17:48:55 2015 Carlos Linares Lopez>
  Last update <martes, 29 marzo 2016 20:13:51 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

// This package provides means to automatically generating text and LaTeX tables
// from simple specification strings as those used in LaTeX.
//
// A specification string consists of an indication how the text is justified in
// all cells in the same column and also how separators shall be
// formatted. Thus, they can refer to either text cells or column separators.
//
// A separator can be present or not and it can be one among the following
// types:
//    |      - a single bar
//    ||     - a double bar
//    |||    - a thick bar
//    @{...} - a verbatim separator where '...' stands for anything but '}'
//    p{width} - creates a column with a fixed width which is specified in LaTeX
//    format. When formatting the table in LaTeX the units specified are
//    honoured; when rendering the table in textual mode, the units are always
//    assumed to be chars
//
// The column is one of the types:
//    c - centered
//    l - left
//    r - right
//
// The tbl package fully supports utf-8 characters
package tbl

import (
	"errors"       // for raising errors
	"fmt"          // printing services
	"log"          // Fatal messages
	"math"         // Max
	"regexp"       // for processing specification strings
	"strconv"      // Atoi
	"strings"      // for repeating characters
	"unicode/utf8" // provides support for UTF-8 characters
)

// global variables
// ----------------------------------------------------------------------------

// A specification string consists of an indication how the text is justified in
// all cells in the same column and also how separators shall be
// formatted. Thus, they can refer to either text cells or column separators.
//
// A separator can be present or not and it can be one among the following
// types:
//    |      - a single bar
//    ||     - a double bar
//    |||    - a thick bar
//    @{...} - a verbatim separator where '...' stands for anything but '}'
//    p{width} - creates a column with a fixed width which is specified in LaTeX
//    format. When formatting the table in LaTeX the units specified are
//    honoured; when rendering the table in textual mode, the units are always
//    assumed to be chars
//
// The column is one of the types:
//    c - centered
//    l - left
//    r - right
var reSpecification = regexp.MustCompile(`^(p{[^}]*\}|\@\{[^}]*\}|\|\|\||\|\||\||c|l|r)`)

// Verbatim separators are also processed with their own regular expression for
// extracting the text
var reVerbatimSeparator = regexp.MustCompile(`^@\{(?P<text>[^}]*)\}`)

// Likewise, fixed widths are processed separately with an additional regular
// expression to extract the width
var reFixedWidth = regexp.MustCompile(`^p\{(?P<width>[^}]*)\}`)

// While the full width is passed to the LaTeX code, only the integer part is
// used to set the width of a column in textual mode. Thus, an additional regexp
// is used just to extract it
var reIntegerFixedWidth = regexp.MustCompile(`^(?P<value>[\d]+).*`)

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
// they have no information about their current location or their surrounding
// environment.
type cellType struct {
	content contentType
	width   int
	text    string
}

// A column type is specified with a cell type so that it has its own type (and
// all items in the same column are expected to be of the same type), it has its
// own width and, in case it is a verbatim column, it holds a copy of the text
// to write
type tblColumn cellType

// A horizontal rule (of a specific type) is characterized by its beginning and
// end.
type tblRule struct {
	content       contentType
	from, to      int
	specification string
}

// Rules are grouped in collections of rules so that the same line can hold an
// arbitrary number of rules if required
type tblRuleCollection []tblRule

// Lines can either contain text (content=TEXT) or a horizontal rule (of any
// type, ie, booktabs rules, ordinary rules or clines). In case it is a
// horizontal rule, the slice rules stores information about them. In any case,
// lines are made up of cells of different types
type tblLine struct {
	content contentType
	rules   tblRuleCollection
	cell    []cellType
}

// A table consists mainly of two components: information about the columns and
// information about the rows. Additionally, a table contains a slice of widths
// with the overall width of each cell in every line. It also stores the
// specification string used to create it.
type Tbl struct {
	column        []tblColumn
	row           []tblLine
	width         []int
	specification string
}

// Functions
// ----------------------------------------------------------------------------

// Return the information of a column according to the given string. Note that
// columns are specified as cells and thus, their width and text (if applicable)
// shall be computed as well.
func getColumnType(cmd string) (column tblColumn) {

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
		if reVerbatimSeparator.MatchString(cmd) {

			// if this has been recognized as a legal verbatim
			// column, extract its contents and return them. Note
			// that in this case no extra spaces are added either
			// before or after the separator
			tag := reVerbatimSeparator.FindStringSubmatchIndex(cmd)
			return tblColumn{VERTICAL_VERBATIM,
				utf8.RuneCountInString(cmd[tag[2]:tag[3]]),
				cmd[tag[2]:tag[3]]}
		} else if reFixedWidth.MatchString(cmd) {

			// ... or a fixed widht column
			tag0 := reFixedWidth.FindStringSubmatchIndex(cmd)
			arg := cmd[tag0[2]:tag0[3]]

			// If this is the case, create a column of the proper
			// type indicating the width provided by the user
			if reIntegerFixedWidth.MatchString(arg) {

				tag1 := reIntegerFixedWidth.FindStringSubmatchIndex(arg)
				width, err := strconv.Atoi(arg[tag1[2]:tag1[3]])
				if err != nil {
					log.Fatalf(" Impossible to extract the width from '%v'",
						arg[tag1[2]:tag1[3]])
				} else {
					return tblColumn{VERTICAL_FIXED_WIDTH,
						width, ""}
				}
			} else {
				log.Fatalf(" Impossible to extract an integer width from '%v'",
					arg)
			}
		}

		// otherwise, raise an error
		log.Fatalf(" Unknown column type '%v'\n", cmd)
	}

	return
}

// Return a new instance of Tbl from a specification string
func NewTable(cmd string) (table Tbl, err error) {

	// before starting to consume the specification string, copy it
	table.specification = cmd

	// just simply process the specification string
	for reSpecification.MatchString(cmd) {

		// get the next item in the specification string and add it to
		// the collection of columns of this table
		tag := reSpecification.FindStringSubmatchIndex(cmd)
		column := getColumnType(cmd[tag[2]:tag[3]])
		table.column = append(table.column, column)

		// Initialize also the width of each cell
		table.width = append(table.width, column.width)

		// move forward
		cmd = cmd[tag[1]:]
	}

	// In case the specification string has not been fully processed, a
	// syntax error has been found
	if cmd != "" {

		return Tbl{},
			errors.New(fmt.Sprintf("Syntax error in a specification string at '%v'\n", cmd))
	}

	// otherwise, copy the specification string and return the table with no
	// error
	return
}

// Methods
// ----------------------------------------------------------------------------

// translate a *user* column index into an *effective* column index:
//
//    * User columns: those with user contents.
//    * Effective columns: those defined in the specification string, either
//                         text or column separators.
//
// Importantly, while the user column index is 1-based, the effective column
// index is 0-based
func (table *Tbl) getEffectiveColumn(user int) int {

	// iterate over all columns until the specified user column has been
	// found
	for idx, jdx := 0, 0; idx < len(table.column); idx += 1 {
		if table.column[idx].content == LEFT ||
			table.column[idx].content == CENTER ||
			table.column[idx].content == RIGHT ||
			table.column[idx].content == VERTICAL_VERBATIM ||
			table.column[idx].content == VERTICAL_FIXED_WIDTH {

			// increment the number of user column indexes found so
			// far
			jdx += 1

			// if this is the user column index requested then
			// return its effective counter
			if jdx == user {
				return idx
			}
		}
	}

	// if all columns have been processed and the user column index
	// requested has not been found, then show an error
	log.Fatalf(" User column index '%v' out of bounds!", user)
	return -1
}

// Add a single line of text to the bottom of the receiver table. The contents
// are specified as a slice of strings. In case the number of items is less than
// the number of columns, the row is paddled with empty strings. If the number
// of items in the given slice exceeds the number of columns in this table, an
// error is raised
func (table *Tbl) AddRow(row []string) (err error) {

	// First of all, in case the last line was a horizontal rule, redo it
	// since we are about to generate a new line
	table.redoLastLine()

	// insert all cells of this line: those provided by the user and others
	// provided in the specification string. Since this line does not
	// contain horizontal rules, an empty rule is used
	newRow := tblLine{TEXT,
		tblRuleCollection{},
		[]cellType{}}
	idx := 0
	for jdx, value := range table.column {

		// depending upon the type of this column cell
		switch value.content {

		case VOID, BLANK, VERTICAL_SINGLE, VERTICAL_DOUBLE, VERTICAL_THICK, VERTICAL_VERBATIM:

			// in case this cell does not contain text specified by
			// the user, just copy its specification
			newRow.cell = append(newRow.cell, cellType(value))

		case VERTICAL_FIXED_WIDTH:

			// compute the contents of this cell: in case a string
			// has been given, and it fits the cell, then just add
			// it after addig as many spaces as necessary to make it
			// fit the width of the cell
			var text string
			if idx < len(row) {
				content := row[idx]
				if utf8.RuneCountInString(content) <= value.width {
					text = content + strings.Repeat(" ",
						value.width-
							utf8.RuneCountInString(content))
				} else {

					// otherwise trim the string
					text = content[0:value.width-1] + "►"
				}
			} else {

				// and, if no string was given, then
				// artificially add empty strings
				text = strings.Repeat(" ", value.width)
			}

			// finally, make sure user text is surrounded by blank
			// spaces unless the previous/next column are verbatim
			if jdx == 0 ||
				table.column[jdx-1].content != VERTICAL_VERBATIM {
				text = " " + text
			}
			if jdx == len(table.column)-1 ||
				table.column[jdx+1].content != VERTICAL_VERBATIM {
				text = text + " "
			}

			// update the width of this column after taking into
			// account the surrounding blank spaces
			table.width[jdx] = int(math.Max(float64(table.width[jdx]),
				float64(utf8.RuneCountInString(text))))

			// create the cell and add it to this row
			newRow.cell = append(newRow.cell, cellType{value.content,
				utf8.RuneCountInString(text),
				text})

			// and move to the next entry provided by the user
			idx += 1

		case LEFT, CENTER, RIGHT:

			// compute the contents of this cell which also consists
			// of surrounding blank characters if needed
			content := " "
			if idx < len(row) {

				// make sure user text is surrounded by blank
				// spaces unless the previous/next column are
				// verbatim column
				if jdx == 0 ||
					table.column[jdx-1].content != VERTICAL_VERBATIM {
					content = " " + row[idx]
				} else {
					content = row[idx]
				}
				if jdx == len(table.column)-1 {
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
			// and add it to the table with a blank space following
			// immediate after as well
			table.width[jdx] = int(math.Max(float64(table.width[jdx]),
				float64(utf8.RuneCountInString(content))))
			newRow.cell = append(newRow.cell,
				cellType{value.content,
					table.width[jdx],
					content})

			// and move to the next entry provided by the user
			idx += 1
		default:
			return errors.New(fmt.Sprintf("Unknown column type '%v'\n", value))
		}
	}

	// Check there are no left values in the given row
	if idx < len(row) {
		return errors.New(fmt.Sprintf("%v items were given but there are only %v columns", len(row), idx))
	}

	// add this row to the table and exit with no error
	table.row = append(table.row, newRow)
	return nil
}

// Add a single horizontal rule that intersects with the vertical separators
// provided that any have been specified.
func (table *Tbl) HSingleRule() {

	table.hrule(HORIZONTAL_SINGLE,
		LIGHT_UP_AND_RIGHT, LIGHT_UP_AND_LEFT, LIGHT_UP_AND_HORIZONTAL,
		UP_DOUBLE_AND_RIGHT_SINGLE, UP_DOUBLE_AND_LEFT_SINGLE, UP_DOUBLE_AND_HORIZONTAL_SINGLE,
		UP_HEAVY_AND_RIGHT_LIGHT, UP_HEAVY_AND_LEFT_LIGHT, UP_HEAVY_AND_HORIZONTAL_LIGHT)
}

// Add a double horizontal rule that intersects with the vertical separators
// provided that any have been specified.
func (table *Tbl) HDoubleRule() {

	// notice that the intersection of double rules with either double or
	// thick vertical separators is computed with the same UTF-8 characters
	// since other combinations are not allowed by UTF-8
	table.hrule(HORIZONTAL_DOUBLE,
		UP_SINGLE_AND_RIGHT_DOUBLE, UP_SINGLE_AND_LEFT_DOUBLE, UP_SINGLE_AND_HORIZONTAL_DOUBLE,
		DOUBLE_UP_AND_RIGHT, DOUBLE_UP_AND_LEFT, DOUBLE_UP_AND_HORIZONTAL,
		DOUBLE_UP_AND_RIGHT, DOUBLE_UP_AND_LEFT, DOUBLE_UP_AND_HORIZONTAL)
}

// Add a thick horizontal rule that intersects with the vertical separators
// provided that any have been specified.
func (table *Tbl) HThickRule() {

	// notice that the intersection of thick rules with either double or
	// thick vertical separators is computed with the same UTF-8 characters
	// since other combinations are not allowed by UTF-8
	table.hrule(HORIZONTAL_THICK,
		UP_LIGHT_AND_RIGHT_HEAVY, UP_LIGHT_AND_LEFT_HEAVY, UP_LIGHT_AND_HORIZONTAL_HEAVY,
		HEAVY_UP_AND_RIGHT, HEAVY_UP_AND_LEFT, HEAVY_UP_AND_HORIZONTAL,
		HEAVY_UP_AND_RIGHT, HEAVY_UP_AND_LEFT, HEAVY_UP_AND_HORIZONTAL)
}

// Add a thick horizontal rule to the current table. Top rules do not draw
// intersections with column separators (they break them instead).
//
// This function is implemented in imitation to the LaTeX package booktabs
func (table *Tbl) TopRule() {

	table.rule(HORIZONTAL_TOP_RULE, HORIZONTAL_THICK)
}

// Add a single horizontal rule to the current table. Mid rules do not draw
// intersections with column separators (they break them instead)
//
// This function is implemented in imitation to the LaTeX package booktabs
func (table *Tbl) MidRule() {

	table.rule(HORIZONTAL_MID_RULE, HORIZONTAL_SINGLE)
}

// Add a thick horizontal rule to the current table. Bottom rules do not draw
// intersections with column separators (they break them instead)
//
// This function is implemented in imitation to the LaTeX package booktabs
func (table *Tbl) BottomRule() {

	table.rule(HORIZONTAL_TOP_RULE, HORIZONTAL_THICK)
}

// Draw an arbitrary number of single horizontal lines from a specific user
// column to another as specified in cmd which consists of a comma-separated list
// of pairs [begin]-[end] with the user column of the first and last column of
// each line.
//
// In case that more than a single cline is given, they can be provided in any
// order
func (table *Tbl) CSingleLine(cmd string) {

	rules := table.parseCLine(cmd)
	table.cline(rules, HORIZONTAL_SINGLE,
		LIGHT_UP_AND_RIGHT, LIGHT_UP_AND_LEFT, LIGHT_UP_AND_HORIZONTAL,
		UP_DOUBLE_AND_RIGHT_SINGLE, UP_DOUBLE_AND_LEFT_SINGLE, UP_DOUBLE_AND_HORIZONTAL_SINGLE,
		UP_HEAVY_AND_RIGHT_LIGHT, UP_HEAVY_AND_LEFT_LIGHT, UP_HEAVY_AND_HORIZONTAL_LIGHT)
}

// Draw an arbitrary number of double horizontal lines from a specific user
// column to another as specified in cmd which consists of a comma-separated list
// of pairs [begin]-[end] with the user column of the first and last column of
// each line.
//
// In case that more than a single cline is given, they can be provided in any
// order
func (table *Tbl) CDoubleLine(cmd string) {

	rules := table.parseCLine(cmd)
	table.cline(rules, HORIZONTAL_DOUBLE,
		UP_SINGLE_AND_RIGHT_DOUBLE, UP_SINGLE_AND_LEFT_DOUBLE, UP_SINGLE_AND_HORIZONTAL_DOUBLE,
		DOUBLE_UP_AND_RIGHT, DOUBLE_UP_AND_LEFT, DOUBLE_UP_AND_HORIZONTAL,
		DOUBLE_UP_AND_RIGHT, DOUBLE_UP_AND_LEFT, DOUBLE_UP_AND_HORIZONTAL)
}

// Draw an arbitrary number of thick horizontal lines from a specific user
// column to another as specified in cmd which consists of a comma-separated list
// of pairs [begin]-[end] with the user column of the first and last column of
// each line.
//
// In case that more than a single cline is given, they can be provided in any
// order
func (table *Tbl) CThickLine(cmd string) {

	rules := table.parseCLine(cmd)
	table.cline(rules, HORIZONTAL_THICK,
		UP_LIGHT_AND_RIGHT_HEAVY, UP_LIGHT_AND_LEFT_HEAVY, UP_LIGHT_AND_HORIZONTAL_HEAVY,
		HEAVY_UP_AND_RIGHT, HEAVY_UP_AND_LEFT, HEAVY_UP_AND_HORIZONTAL,
		HEAVY_UP_AND_RIGHT, HEAVY_UP_AND_LEFT, HEAVY_UP_AND_HORIZONTAL)
}

// Cells draw themselves producing a string which takes into account the width
// of the cell.
func (cell cellType) String() (output string) {

	// depending upon the type of cell
	switch cell.content {
	case VERTICAL_VERBATIM:
		output = cell.text
	case VERTICAL_FIXED_WIDTH:
		output = cell.text
	case LEFT:
		output = cell.text + strings.Repeat(" ",
			cell.width-utf8.RuneCountInString(cell.text))
	case CENTER:
		output = strings.Repeat(" ", (cell.width-utf8.RuneCountInString(cell.text))/2) +
			cell.text +
			strings.Repeat(" ", (cell.width-utf8.RuneCountInString(cell.text))/2+
				(cell.width-utf8.RuneCountInString(cell.text))%2)
	case RIGHT:
		output = strings.Repeat(" ", cell.width-utf8.RuneCountInString(cell.text)) +
			cell.text
	default:
		output = strings.Repeat(characterSet[cell.content], cell.width)
	}

	// and return the string computed so far
	return
}

// Return a string with a representation of the contents of the table
func (table Tbl) String() (output string) {

	// for every single line
	for _, line := range table.row {

		// and for every column
		for jdx, cell := range line.cell {

			// draw this cell after updating its width
			cell.width = table.width[jdx]
			output += fmt.Sprintf("%v", cell)
		}

		// and start a newline
		output += "\n"
	}

	// and return the string computed so far
	return
}

// Return a LaTeX implementation of the contents of this table in a tabular
// environment
func (table Tbl) ToLaTeX() (output string) {

	// start the string with the center and tabular environments
	output = `\begin{center}
  \begin{tabular}{`

	// next, add the specification string
	output += table.specification + "}\n"

	// for every single line in the table
	for _, line := range table.row {

		// now, according to the type of line
		switch line.content {

		case TEXT:

			// if this is a line of text, then extract the contents
			// of all the user columns
			var contents []string
			for _, column := range line.cell[:len(line.cell)] {

				// if this cell contains text, then add it
				if column.content == VERTICAL_FIXED_WIDTH ||
					column.content == LEFT ||
					column.content == CENTER ||
					column.content == RIGHT {
					contents = append(contents, column.text)
				}
			}

			// if and only if there are items to show
			if len(contents) > 0 {

				// now, interleave the contents extracted from this line
				// with the symbol "&"
				for _, column := range contents[:len(contents)-1] {
					output += fmt.Sprintf(" %v &", column)
				}

				// show now the contents of the last user column
				output += fmt.Sprintf(" %v ", contents[len(contents)-1])
			}

			// and end the current line
			output += "\\\\ \n"

		default:

			// if this is not a line of text, then it is expected to
			// be either a horizontal rule or a cline. Traverse then
			// all the rules and print out all the specification
			// strings
			for _, rule := range line.rules {

				output += fmt.Sprintf("%v ", rule.specification)
			}

			// if this is not a line of text, then it is expected to
			// be either a horizontal rule or a cline
			// output += "\\hline\n"
		}

	}

	// end the string closing the tabular and center environments
	output += `
  \end{tabular}
\end{center}`

	// finally exit returning the string computed so far
	return
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
