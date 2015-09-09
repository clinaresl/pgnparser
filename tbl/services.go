/* 
  services.go
  Description: General services used by various functions
  ----------------------------------------------------------------------------- 

  Started on  <Wed Sep  9 08:06:09 2015 Carlos Linares Lopez>
  Last update <miércoles, 09 septiembre 2015 17:56:09 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package tbl

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


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
