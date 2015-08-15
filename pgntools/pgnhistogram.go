/* 
  pgnhistogram.go
  Description: Definition of histograms of any order
  ----------------------------------------------------------------------------- 

  Started on  <Thu Jul  2 08:06:07 2015 Carlos Linares Lopez>
  Last update <domingo, 16 agosto 2015 01:17:05 Carlos Linares Lopez (clinares)>
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
	"fmt"			// printing services
	"log"			// logging services
	"strconv"		// string conversion from integers
)

// typedefs
// ----------------------------------------------------------------------------

// A histogram contains either an integer value or a reference to a new
// histogram (so that histograms of any order can be easily created). The
// following typedefs allow the definition of different values

// histograms can store integer values. To allow the definition of values
// arbitrarily large double precision is used
type dataHistValue int64;

// The values stored in a histogram should support increment operations and
// lookups with keys of any length. Keys are specified as slices of
// strings. Additionally, they should support convertion to strings that can be
// printed on a terminal
type histogramCounter interface {
	Increment (index []string, increment dataHistValue) dataHistValue
	Lookup (index []string) dataHistValue
	String () string
}

// A histogram consists simply of a map of strings to counters which can be
// either integers or nested histograms. Importantly, every level of the
// histogram stores the number of items below it
type Histogram struct {
	nbitems int64;
	key map[string]histogramCounter
}

// Functions
// ----------------------------------------------------------------------------

// Return a new instance of Histogram
func NewHistogram () (hist Histogram) {
	return Histogram {0, make (map [string]histogramCounter)}
}

// Methods
// ----------------------------------------------------------------------------

// The following methods allow incrementing operations over different types so
// that they can be used as values in histograms

// A scalar is incremented only in case null indexes are given. Otherwise, an
// error is raised.
//
// The value added is returned in case the operation was successful.
func (value dataHistValue) Increment (index []string, increment dataHistValue) dataHistValue {

	// first, verify that the given index is null. If not, raise an error
	if len (index) > 0 {
		log.Fatal (" A non-null index was given to a terminal entry of a histogram")
	}

	// otherwise, just return the increment
	return increment
}

// A histogram (of an arbitrary order) can be incremented if a non-null index is
// given.  The number of items in the index shall be equal to the order of the
// histogram so that from the root it is possible to reach a terminal entry. In
// case the length of the index and the histogram differ, an error is raised
//
// The value added is returned in case the operation was successful.
func (hist *Histogram) Increment (index []string, increment dataHistValue) dataHistValue {

	// first, in case this is a null index, raise an error
	if len (index) == 0 {
		log.Fatal (" A null index was given to a non-terminal entry of a histogram")
	}

	// in other case, just select the right entry. In case it does not
	// exist, then create it
	_, ok := hist.key [index[0]]; if !ok {

		// Case #1 - This is the last index, so that just create an
		// entry for this key, initialize it to zero and increment its
		// content
		if len (index) == 1 {
			hist.key [index[0]] = dataHistValue (0)
		} else {

			// Case #2 - Otherwise, the histogram should point to a
			// nested histogram
			hist.key [index[0]] = &Histogram {0, make (map[string]histogramCounter)}
		}
	}

	// Before incrementing the number of samples below this histogram,
	// update the private count of items
	hist.nbitems += int64 (increment)
	
	// in case the entry exists then, in case this is the last key, then
	// increment its content
	if len (index) == 1 {

		// make sure these numbers can be added, ie., assert the type of
		// this entry
		value, ok := hist.key [index[0]].(dataHistValue); if !ok {
			log.Fatal (" It was not possible to add an increment to a non-terminal location")
		}
		hist.key [index[0]] = value + increment
		return increment
	}
	
	// otherwise, proceed recursively
	return hist.key[index[0]].Increment (index[1:], increment)
}

// The following methods allow lookups with keys of any length which are
// specified as slices of strings so that they can be used as values in histograms

// Return the value of this particular integer. If the given index is not empty
// an error is raised.
func (value dataHistValue) Lookup (index []string) dataHistValue {

	// first, verify that the given index is null. If not, raise an error
	if len (index) > 0 {
		log.Fatal (" A non-null index was given to a terminal entry of a histogram")
	}

	// otherwise, just return this value
	return value
}

// This method acknowledges either full or partial keys.
//
// If a full key is given, it returns the value attached to it.
//
// In case a partial key is given, it returns the number of items stored below
// the given key.
//
// If the given index is not found, an error is raised
func (hist *Histogram) Lookup (index []string) dataHistValue {

	// first, in case this is a null index, then return the private count of
	// items below it
	if len (index) == 0 {
		return dataHistValue (hist.nbitems)
	}

	// in other case, just select the right entry. In case it does not
	// exist, then raise an error
	entry, ok := hist.key [index[0]]; if !ok {
		log.Fatal (" A null index was given to a non-terminal entry of a histogram")
	}

	// finally, look for this specific value recursively
	return entry.Lookup (index[1:])
}

// The following service just returns a string representation of this value
// which is known to be a frequency expressed as a double-precision integer
func (value dataHistValue) String () string {

	// note that Itoa is used instead of Sprintf to avoid an infinite
	// recursion
	return strconv.Itoa (int (value))
}

// The following method routinely converts the information in a histogram into a
// string that can be printed to a terminal
func (hist *Histogram) String () string {

	var output string
	for index, value := range hist.key {

		// Check the type of the value of this key. In case it is
		// another histogram ...
		_, ok := value.(dataHistValue); if !ok {
			
			// ... compute the string that corresponds to every
			// entry of this nested string
			output += fmt.Sprintf ("%10v:\n%v", index, value.String ())
		} else {

			// otherwise, just add this value
			output += fmt.Sprintf (" %10v: %10v\n", index, value)
		}
	}

	return output
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
