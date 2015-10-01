/* 
  services.go
  Description: General services used by various functions
  ----------------------------------------------------------------------------- 

  Started on  <Wed Sep  9 08:06:09 2015 Carlos Linares Lopez>
  Last update <jueves, 01 octubre 2015 09:49:06 Carlos Linares Lopez (clinares)>
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
	"sort"			// used for sorting rules
	"strconv"		// Atoi
)

// global variables
// ----------------------------------------------------------------------------

// a cline can be specified as a single line or an arbitrary collection of them
// separated by either commas or blank characters. The following regexp is used
// just to process each cline separately.
var reCLines = regexp.MustCompile (`^(\s*\d+\s*-\s*\d+)[,]?`)

// single clines are recognized with the following regular expression
// var reCLine = regexp.MustCompile (`(?P<from>\d+)-(?P<to>\d+)`)
var reCLine = regexp.MustCompile (`^\s*(?P<from>\d+)\s*-\s*(?P<to>\d+)[,]?`)


// Methods
// ----------------------------------------------------------------------------

// Return the number of items in a specific collection
func (rules tblRuleCollection) Len () int {
	return len (rules)
}

// Swap two rules in the same collection
func (rules tblRuleCollection) Swap (i, j int) {
	rules[i], rules[j] = rules[j], rules[i]
}

// Return whether the first rule is less than the second one
func (rules tblRuleCollection) Less (i, j int) bool {
	return rules[i].from < rules[j].from
}

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
func (table *Tbl) hrule (content, light_sw, light_se, light_s, double_sw, double_se, double_s, thick_sw, thick_se, thick_s contentType) {

	// simply draw a line (ie., a single rule) that goes over all columns of
	// the table, ie., from 0 to the last column ---and this is specified
	// with a slice of rules which consist of a single rule whose bounds are
	// literally specified
	table.line ([]tblRule{tblRule{content, 0, len (table.column) - 1}}, content, light_sw, light_se, light_s, double_sw, double_se, double_s, thick_sw, thick_se, thick_s)
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
	// function. Obviously, the rule goes from the first column until the
	// last one
	var newRow tblLine
	newRow = tblLine{content,
		tblRule{content, 0, len (table.column)-1},
		[]cellType{}}
	
	for idx := range table.column {
		newRow.cell = append (newRow.cell, cellType {thickness,
			table.width[idx], ""})
	}
	table.row = append (table.row, newRow)
}

// Add a partial line (or more, see below) to the bottom of the current table as
// in the LaTeX command \cline. These lines intersect with the vertical
// separators provided that any have been specified. The bounds of the line are
// specified in a command string which follows the LaTeX syntax: "begin-end"
// where begin and end are user columns or columns that contain text ---whereas
// *effective* columns are those defined in the specification string given when
// creating the table.
//
// This specification can be extended to indicate an arbitrary number of clines
// in the same line with the syntax: "begin-end[, begin-end]*"
//
// The type of lines to draw is specified with content which should take one of
// the values HORIZONTAL_SINGLE, HORIZONTAL_DOUBLE or HORIZONTAL_THICK
func (table *Tbl) cline (cmd string, content, light_sw, light_se, light_s, double_sw, double_se, double_s, thick_sw, thick_se, thick_s contentType) {
		
	var err error
	var from, to int
	
	// the following slice shall hold the rules specified in the
	// specification string
	var rules tblRuleCollection
	
	// While a specification of a cline is found at the beginning of the strinng
	for ;reCLines.MatchString (cmd); {

		// extract the specification of the next cline
		tag := reCLines.FindStringSubmatchIndex (cmd)
		interval := cmd[tag[0]:tag[1]]

		// and now process its components extracting the user columns
		// given as 'from' and 'to'
		if reCLine.MatchString (interval) {
			itag := reCLine.FindStringSubmatchIndex (interval)

			// extract the limits of this cline
			from, err = strconv.Atoi (interval[itag[2]:itag[3]]); if err != nil {
				log.Fatalf (" Error while extracting the first bound from '%v'",
					interval[itag[2]:itag[3]])
			}
			to, err = strconv.Atoi (interval[itag[4]:itag[5]]); if err != nil {
				log.Fatalf (" Error while extracting the second bound from '%v'",
					interval[itag[4]:itag[5]])
			}

			// 'from' and 'to' are given as user column indexes. Translate them into
			// effective column indexes
			from = table.getEffectiveColumn (from)
			to = table.getEffectiveColumn (to)

			// there is however two exceptions:
			// 1. if the user specified a user column as 'from' which is preceded of
			// a vertical separator, then start the cline in the previous column
			if from >= 1 && table.column[from-1].content != LEFT &&
				table.column[from-1].content != CENTER &&
				table.column[from-1].content != RIGHT &&
				table.column[from-1].content != VERTICAL_VERBATIM &&
				table.column[from-1].content != VERTICAL_FIXED_WIDTH {
				from -= 1
			}
			
			// 2. if the user specified a user column as 'to' which is continued by
			// a vertical separator, then end the cline in the next column
			if (to < len (table.column) - 1) && table.column[to+1].content != LEFT &&
				table.column[to+1].content != CENTER &&
				table.column[to+1].content != RIGHT &&
				table.column[to+1].content != VERTICAL_VERBATIM &&
				table.column[to+1].content != VERTICAL_FIXED_WIDTH {
				to += 1
			}

			// add this two bounds to the current slice
			rules = append (rules, tblRule{content, from, to})
		}

		// move forward in the specification string
		cmd = cmd[tag[1]:]
	}

	// verify now that the whole specification string was exhausted. If not,
	// a regexp could not be applied, thus, there were a syntax error
	if len (cmd) > 0 {
		log.Fatalf (" Syntax error in the cline specification string '%v'\n", cmd)
	}

	// since the user could have provided the rules in any order but they
	// have to be processed in ascending order of the 'from' field, sort
	// them now
	sort.Sort (rules)
	
	// and now, simply draw the rules in the current line
	table.line (rules, content, light_sw, light_se, light_s, double_sw, double_se, double_s, thick_sw, thick_se, thick_s)
}

// Add a horizontal line (ie., a partial rule from two specified effective
// column indexes 'from' and 'to') that intersects with the vertical separators
// provided that any have been specified. The type of rule is defined by the
// parameter:
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
func (table *Tbl) line (rules []tblRule, content, light_sw, light_se, light_s, double_sw, double_se, double_s, thick_sw, thick_se, thick_s contentType) {

	// INVARIANT: this code assumes that rule consists of a disjoint
	// sequence of rules which are sorted in increasing order of 'from'
	
	// Since it is possible to concatenate horizontal rules, redo the last
	// one if necessary
	table.redoLastLine ()

	from, to := rules[0].from, rules[0].to

	// create a new row with a single line between the specified bounds
	// whose contents will be computed in this function.
	newRow := tblLine{content,
		tblRule{content, from, to},
		[]cellType{}}

	// traverse the slice of disjoint rules in ascending order of
	// 'from'. jdx holds the index of the first rule (which is initially -1)
	// and idx will contain the index of the column under consideration
	jdx := -1
	
	// consider now all columns from the general specification of the table
	// and draw the intersections accordingly
	for idx, column := range table.column {

		// update jdx in case a new rule is been visited. This check
		// includes verifying that incrementing jdx makes sense, ie.,
		// that there are more rules to consider
		if jdx < len (rules) - 1 &&  idx == rules[1+jdx].from {
			jdx += 1
		}

		// in case we are in a column not covered by any of the rules,
		// preserve the type of the column without any contents
		if jdx == -1 || idx > rules[jdx].to {
			newRow.cell = append (newRow.cell,
				cellType {column.content, table.width[idx], ""})
		} else {

			// otherwise, choose the right character to show
			// according to the type of this vertical separator
			switch column.content {
			case VERTICAL_SINGLE:
				table.lineColumn (idx, rules[jdx], &newRow, light_sw, light_se, light_s)

			case VERTICAL_DOUBLE:
				table.lineColumn (idx, rules[jdx], &newRow, double_sw, double_se, double_s)
			
			case VERTICAL_THICK:
				table.lineColumn (idx, rules[jdx], &newRow, thick_sw, thick_se, thick_s)

			default:
				newRow.cell = append (newRow.cell,
					cellType {content, 
						table.width[idx], ""})
			}
		}
	}
	
	// and add this row to the bottom of the table
	table.row = append (table.row, newRow)
}

// Add a single character to 'row' wrt to the effective column index 'idx'. This
// function already takes into account that the row to draw is delimited by
// [from, to] as stored in row.rule. The type of rule generated by consecutive
// invocations to this service is defined by the following parameters:
//
//    sw, se, s - south/west, south/east and south separators used for different
//    types of vertical separators
//
// INVARIANT - This function is invoked solely to draw characters that fall
// within the interval of the rule in row. Character falling outside the rule
// are drawn easily by preserving the content of other vertical separators (ie.,
// columns)
func (table *Tbl) lineColumn (idx int, rule tblRule, row *tblLine, sw, se, s contentType) {

	// if a line starts at this particular location, draw the sw character
	if idx==rule.from {
		row.cell = append (row.cell,
			cellType {sw,
				table.width[idx], ""})
	} else if idx == rule.to {

		// otherwise, in case a line is ended at this specific column,
		// then draw the se character
		row.cell = append (row.cell,
			cellType {se,
				table.width[idx], ""})
	} else {

		// and, by default, just draw the south item
		row.cell = append (row.cell,
			cellType{s,
				table.width[idx], ""})
	}
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
