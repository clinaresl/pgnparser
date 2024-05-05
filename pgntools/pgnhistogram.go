/*
  pgnhistogram.go
  Description: Definition of histograms of any order
  -----------------------------------------------------------------------------

  Started on  <Thu Jul  2 08:06:07 2015 Carlos Linares Lopez>
  Last update <lunes, 04 abril 2016 17:46:52 Carlos Linares Lopez (clinares)>
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

	"github.com/clinaresl/table"
)

// typedefs
// ----------------------------------------------------------------------------

// Given a sequence of criteria (either variables or boolean expressions), a
// histogram is implemented as a decision tree where internal nodes store the
// result of the criteria and lead to other histograms or decision trees until
// the leaves are reached which simply store the number of occurrences of all
// variables/boolean expressions from the root to it.
//
// In addition, a histogram contains the total number of observations stored in
// it so that percentages can be computed for every inner/leaf node
type PgnHistogram struct {
	criteria []string
	data     map[string]any
	nbhits   uint64
}

// Functions
// ----------------------------------------------------------------------------

// Return the result of executing the given criteria as a string with
// information in the specified game and nil if no error happened.
func getResult(criteria string, game PgnGame) (string, error) {

	// execute the ith-criteria of this histogram
	env := game.getEnv()
	output, err := evaluateExpr(criteria, env)
	if err != nil {
		return "", err
	}

	// return the result casted as a string with success
	return fmt.Sprintf("%v", output), nil
}

// Methods
// ----------------------------------------------------------------------------

// Return a brand new PgnHistogram defined with a string, which consists of a
// semicolon list of variables/boolean expressions in the form: "<var/expr>+".
// At least one should be given, and an arbitrary number of them can be
// specified.
func NewPgnHistogram(spec string) PgnHistogram {

	// Compute the sequence of criteria from the specification string
	criteria := reHistogramCriteria.Split(spec, -1)

	// finally, return a new histogram with the decision tree built above and no
	// hits
	return PgnHistogram{
		criteria: criteria,
		data:     make(map[string]any),
		nbhits:   0,
	}
}

// Return the nbhits that are reached by using all values in the given sequence.
// This function assumes that such value can be effectively achieved by using
// the given sequence
func (histogram PgnHistogram) getHits(sequence []any) uint64 {

	// The implementation is performed iteratively
	data := histogram.data

	// Traverse all keys but the last one
	for idx := 0; idx < len(sequence)-1; idx++ {
		data = data[sequence[idx].(string)].(map[string]any)
	}

	// Once the last value has been found, just return it
	return data[sequence[len(sequence)-1].(string)].(uint64)
}

// Updates this histogram with information in the given game, and nil if no
// error was found
func (histogram *PgnHistogram) Add(game PgnGame) error {

	// get the map of this histogram
	data := histogram.data

	// process all criteria in this histogram
	idx := 0
	for idx < len(histogram.criteria)-1 {

		// execute the ith-criteria of this histogram
		result, err := getResult(histogram.criteria[idx], game)
		if err != nil {
			return err
		}

		// Next verify whether this result is already stored in the current map
		if value, ok := data[result]; !ok {

			// in case it did not exist, then create a nexted map[string]any and
			// update the current data
			data[result] = make(map[string]any)
			data = data[result].(map[string]any)
		} else {

			// if it exists, just update the current data
			data = value.(map[string]any)
		}

		// and move forward
		idx += 1
	}

	// Once the leaf has been found, then add a new observation. Do as before,
	// evaluate the last criteria and add data to the histogram adding a new
	// keyword if necessary
	result, err := getResult(histogram.criteria[idx], game)
	if err != nil {
		return err
	}

	// Next verify whether this result is already stored in the current map
	if _, ok := data[result]; !ok {

		// in case it did not exist, then add the first observation
		data[result] = uint64(1)
	} else {

		// otherwise, increment it
		value := data[result].(uint64)
		data[result] = value + 1
	}

	// Update the number of observations of this histogram and return with
	// success
	histogram.nbhits += 1
	return nil
}

// Histograms are stringers, so that they can be shown on any writer
func (histogram PgnHistogram) String() string {

	// create a table to show the data in this histogram where all columns but
	// first are criteria, and the last is the number of observations
	nocols := 0
	spec := " c "
	for ; nocols < len(histogram.criteria); nocols++ {
		spec += "| c "
	}
	tab, _ := table.NewTable(spec)

	// The headers of the table are just the criteria
	line := make([]any, 0)
	for _, icriteria := range histogram.criteria {
		line = append(line, icriteria)
	}

	// add the header for the last column and add this line to the table
	// followed by a horizontal rule
	line = append(line, "# Obs.")
	tab.AddRow(line...)
	tab.AddThickRule()

	// Next, add data to a slice of slices of strings. Some preprocessing is
	// necessary because the table shows horizontal rules to distinguish each
	// line from the *next* one. Unless all rows are computed, it is not
	// possible to know where to place the horizontal rules
	contents := make([][]any, 0)

	// Data is generated by accessing the number of hits for all combinations of
	// the criteria. These combinations are obtained by "flatting" the data map
	lines := flatMap(histogram.data)
	for idx, ikey := range lines {

		// if this is the first line, then use the values of all keys.
		// Otherwise,, get the diff slice
		var irow []any
		if idx == 0 {
			irow = ikey
		} else {
			irow = diffSlice(lines[idx-1], ikey)
		}

		// And add the value of all criteria and, at the end, the number of hits
		// for this specific combination
		irow = append(irow, fmt.Sprintf("%v", histogram.getHits(ikey)))
		contents = append(contents, irow)
	}

	// Once the contents of the entire table have been computed, add the rows to
	// the table
	for idx, iline := range contents {

		// Add this line to the table
		tab.AddRow(iline...)

		// Add a horizontal rule distinguishing the columns with values which
		// are different than the values of the next one, so that for all lines
		// but the last one
		if idx < len(contents)-1 {

			eqcols := 0
			for contents[idx+1][eqcols] == "" {
				eqcols += 1
			}

			// And show the horizontal rule
			tab.AddSingleRule(eqcols, nocols+1)
		}
	}

	// Add a bottom row and return the table
	tab.AddThickRule()

	return fmt.Sprintf("%v", tab)
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
