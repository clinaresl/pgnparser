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

import "fmt"

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

// Updates this histogram with information in the given game, and nil if no
// error was found
func (histogram *PgnHistogram) Add(game PgnGame) error {

	// get the map of this histogram
	data := histogram.data

	// process all criteria in this histogram
	idx := 0
	for idx < len(histogram.criteria)-1 {

		// execute the ith-criteria of this histogram
		env := game.getEnv()
		output, err := evaluateExpr(histogram.criteria[idx], env)
		if err != nil {
			return err
		}
		result := fmt.Sprintf("%v", output)

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
	env := game.getEnv()
	output, err := evaluateExpr(histogram.criteria[idx], env)
	if err != nil {
		return err
	}
	result := fmt.Sprintf("%v", output)

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

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
