/* 
  services.go
  Description: General services used by various functions
  ----------------------------------------------------------------------------- 

  Started on  <Wed Sep  9 08:06:09 2015 Carlos Linares Lopez>
  Last update <miércoles, 09 septiembre 2015 22:55:01 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package tbl

import (
	"log"			// Fatal messages
	"regexp"		// for processing specification strings
	"strconv"		// Atoi
)

// global variables
// ----------------------------------------------------------------------------

// While the full width is passed to the LaTeX code, only the integer part is
// used to set the width of a column. Thus, an additional regexp is used just to
// extract it
var reIntegerFixedWidth = regexp.MustCompile (`^(?P<value>[\d]+).*`)

// clines are recognized with the following regular expression
var reCLine = regexp.MustCompile (`(?P<from>\d+)-(?P<to>\d+)`)


// Methods
// ----------------------------------------------------------------------------

// Add a horizontal rule that intersects with the vertical separators
// provided that any have been specified. The type of rule is defined
// by the parameter:
//
//    content - specifies whether this is a single/double/thick horizonntal
//    rule. Legal values are: HORIZONTAL_SINGLE, HORIZONTAL_DOUBLE and
//    HORIZONTAL_THICK
//
// When adding a rule, intersections with vertical separators specified in the
// creation of the table are taken into account as well. What characters should
// be used is specified in the following parameters:
//
//    *_sw, *_se, *_s - south/west, south/east and south separators used for
//    different types of vertical separators as specified in '*' that can take
//    the following values: light, double and thick
func (table *Tbl) hrule (content,
	light_sw, light_se, light_s,
	double_sw, double_se, double_s,
	thick_sw, thick_se, thick_s contentType) {

	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	table.redoLastLine ()

	// create a new row whose contents will be computed in this
	// function. Importantly, the beginning of the rule depends on whether
	// there is an initial column at location 0 or not: if there is a column
	// at location 0, the rule starts at location 1 so that when redrawing
	// this horizontal rule the first character is set properly
	var newRow tblLine
	if table.column[0].content >= VERTICAL_SINGLE &&
		table.column[0].content <= VERTICAL_THICK {
		newRow = tblLine{content,
			tblRule{content, 1, len (table.column)-1},
			[]cellType{}}
	} else {
		newRow = tblLine{content,
			tblRule{content, 0, len (table.column)-1},
			[]cellType{}}
	}

	// consider now all columns from the general specification of the table
	// and draw the intersections accordingly
	for idx, column := range table.column {
		switch column.content {
		case VERTICAL_SINGLE:
			if idx==0 {
				newRow.cell = append (newRow.cell,
					cellType {light_sw,
						table.width[idx], ""})
			} else if idx == len (table.column) - 1 {
				newRow.cell = append (newRow.cell,
					cellType {light_se,
						table.width[idx], ""})
			} else {
				newRow.cell = append (newRow.cell,
					cellType{light_s,
						table.width[idx], ""})
			}

		case VERTICAL_DOUBLE:

			if idx==0 {
				newRow.cell = append (newRow.cell,
					cellType {double_sw,
						table.width[idx], ""})
			} else if idx == len (table.column) - 1 {
				newRow.cell = append (newRow.cell,
					cellType {double_se,
						table.width[idx], ""})
			} else {
				newRow.cell = append (newRow.cell,
					cellType{double_s,
						table.width[idx], ""})
			}
			
		case VERTICAL_THICK:

			if idx==0 {
				newRow.cell = append (newRow.cell,
					cellType {thick_sw,
						table.width[idx], ""})
			} else if idx == len (table.column) - 1 {
				newRow.cell = append (newRow.cell,
					cellType {thick_se,
						table.width[idx], ""})
			} else {
				newRow.cell = append (newRow.cell,
					cellType{thick_s,
						table.width[idx], ""})
			}
		default:
			newRow.cell = append (newRow.cell,
				cellType {content, 
					table.width[idx], ""})
		}
	}
	table.row = append (table.row, newRow)	
}

// Add a horizontal rule to the bottom of the current table as in the LaTeX
// package booktabs. This type of rules do not draw intersections with column
// separators (they break them instead).
//
// The type of rule is identified with the content which might be either
// HORIZONTAL_TOP_RULE or HORIZONTAL_MID_RULE. Bottom rules are exactly equal to
// top rules
//
// Additionally, the thickness of each type of rule is described with an
// additional attribute, thickness, which should usually be either
// HORIZONTAL_SINGLE or HORIZONTAL_THICK
func (table *Tbl) rule (content, thickness contentType) {

	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	table.redoLastLine ()
	
	// create a new row whose contents will be computed in this
	// function. Importantly, the beginning of the rule depends on whether
	// there is an initial column at location 0 or not: if there is a column
	// at location 0, the rule starts at location 1 so that when redrawing
	// this horizontal rule the first character is set properly
	var newRow tblLine
	if table.column[0].content >= VERTICAL_SINGLE &&
		table.column[0].content <= VERTICAL_THICK {
		newRow = tblLine{content,
			tblRule{content, 1, len (table.column)-1},
			[]cellType{}}
	} else {
		newRow = tblLine{content,
			tblRule{content, 0, len (table.column)-1},
			[]cellType{}}
	}
	
	for idx := range table.column {
		newRow.cell = append (newRow.cell, cellType {thickness,
			table.width[idx], ""})
	}
	table.row = append (table.row, newRow)
}

// Add a partial line to the bottom of the current table as in the LaTeX command
// \cline. These lines intersect with the vertical separators provided that any
// have been specified. The bounds of the line are specified in a command string
// which follows the LaTeX syntax: "begin-end" where begin and end are user
// columns or columns that contain text ---whereas *effective* columns are those
// defined in the specification string given when creating the table.
//
// The type of line is specified with content which should take one of the
// values HORIZONTAL_SINGLE, HORIZONTAL_DOUBLE or HORIZONTAL_THICK
func (table *Tbl) cline (cmd string, content contentType) {

	var err error
	var from, to int
	
	// parse the given command
	if reCLine.MatchString (cmd) {
		tag := reCLine.FindStringSubmatchIndex (cmd)

		// extract the limits of this cline
		from, err = strconv.Atoi (cmd[tag[2]:tag[3]]); if err != nil {
			log.Fatalf (" It was not feasible to extract the first bound from '%v'",
				cmd[tag[2]:tag[3]])
		}
		to, err = strconv.Atoi (cmd[tag[4]:tag[5]]); if err != nil {
			log.Fatalf (" It was not feasible to extract the second bound from '%v'",
				cmd[tag[4]:tag[5]])
		}
	} else {
		log.Fatalf ("Incorrect cline specification: '%v'",
			cmd)
	}

	// 'from' and 'to' are given as user column indexes. Translate them into
	// effective column indexes
	from, to = table.getEffectiveColumn (from), table.getEffectiveColumn (to)
	
	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	table.redoLastLine ()
	
	// A cline consists of lines in those areas specified by the user and
	// blank characters otherwise
	newRow := tblLine{content,
		tblRule{content, from, to},
		[]cellType{}}
	for idx, cell := range table.column {

		// first update the column number of this one wrt columns with
		// content ---ie., ignoring separators
		if cell.content == LEFT ||
			cell.content == CENTER ||
			cell.content == RIGHT ||
			cell.content == VERTICAL_VERBATIM ||
			cell.content == VERTICAL_FIXED_WIDTH {

			// and now check whether this one shall be drawn
			if idx < from || idx > to {
				newRow.cell = append (newRow.cell,
					cellType {BLANK,
						table.width[idx], ""})
			} else {
				newRow.cell = append (newRow.cell,
					cellType {content,
						table.width[idx], ""})
			}
		} else {

			// in case we are out of bounds, just simply preserve
			// the type of column at this position
			newRow.cell = append (newRow.cell,
				cellType {cell.content, table.width[idx], ""})
		}
	}

	// and add this row to the bottom of the table
	table.row = append (table.row, newRow)
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
