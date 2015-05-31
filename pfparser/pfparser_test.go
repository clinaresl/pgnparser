/* 
  pfparser_test.go
  Description: Unit tests for the propositional formula parser package
  ----------------------------------------------------------------------------- 

  Started on  <Sun May 24 23:26:09 2015 Carlos Linares Lopez>
  Last update <domingo, 31 mayo 2015 17:51:52 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package pfparser

import (
	"log"
	"testing"
)

// the following function parses the given pformula and evaluates
// it. In case the result differs from the expected one, a Fatal error
// is raised using the testing framework specified in t
func assert (t *testing.T, pformula string, expected TypeBool) {

	var err error
	var logicalEvaluator LogicalEvaluator

	log.Println (pformula, expected)

	// Parse the given formula
	logicalEvaluator, err = Parse (&pformula, 0); if err != nil {
		log.Fatalf ("%v\n", err)
	}

	// and now evaluate it, in case the result is not the expected
	// one then raise a fatal error
	if logicalEvaluator.Evaluate () != expected {
		t.Fatalf (" Error in pformula %v", pformula)
	}	
}

func TestIntegerUnparenthesized (t *testing.T) {

	// create a map that associates simple relational expressions to their
	// expected result using integer constants
	expected := map[string]bool {
		"10 <  1" : false,
		"10 <= 1" : false,
		"10 =  1" : false,
		"10 != 1" : true ,
		"10 >= 1" : true ,
		"10 >  1" : true ,
		
		"10 <  10" : false,
		"10 <= 10" : true ,
		"10 =  10" : true ,
		"10 != 10" : false,
		"10 >= 10" : true ,
		"10 >  10" : false,
	}
	
	// -- simple relational expressions
	for expression, value := range expected {
		assert (t, expression, TypeBool (value))
	}
	
	// -- compound relational expressions

	// ---- two relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {

			// OR
			assert (t, expression1 + " or " + expression2, TypeBool (value1 || value2))

			// AND
			assert (t, expression1 + " and " + expression2, TypeBool (value1 && value2))
		}
	}
	
	// ---- three relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {

				// OR
				assert (t, expression1 + " or " + expression2 + " or " + expression3,
					TypeBool (value1 || value2 || value3))

				// AND
				assert (t, expression1 + " and " + expression2 + " and " + expression3,
					TypeBool (value1 && value2 && value3))

				// OR/AND
				assert (t, expression1 + " or " + expression2 + " and " + expression3,
					TypeBool (value1 || value2 && value3))
				assert (t, expression1 + " and " + expression2 + " or " + expression3,
					TypeBool (value1 && value2 || value3))
			}
		}
	}

	// --- four relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {
				for expression4, value4 := range expected {

					// OR
					assert (t, expression1 + " or " + expression2 + " or " + expression3 + " or " + expression4,
						TypeBool (value1 || value2 || value3 || value4))

					// AND
					assert (t, expression1 + " and " + expression2 + " and " + expression3 + " and " + expression4,
						TypeBool (value1 && value2 && value3 && value4))

					// OR/AND
					assert (t, expression1 + " or " + expression2 + " or " + expression3 + " and " + expression4,
						TypeBool (value1 || value2 || value3 && value4))
					assert (t, expression1 + " or " + expression2 + " and " + expression3 + " or " + expression4,
						TypeBool (value1 || value2 && value3 || value4))
					assert (t, expression1 + " and " + expression2 + " or " + expression3 + " or " + expression4,
						TypeBool (value1 && value2 || value3 || value4))
					assert (t, expression1 + " or " + expression2 + " and " + expression3 + " and " + expression4,
						TypeBool (value1 || value2 && value3 && value4))
					assert (t, expression1 + " and " + expression2 + " or " + expression3 + " and " + expression4,
						TypeBool (value1 && value2 || value3 && value4))
					assert (t, expression1 + " and " + expression2 + " and " + expression3 + " or " + expression4,
						TypeBool (value1 && value2 && value3 || value4))
				}
			}
		}
	}
}

func TestStringUnparenthesized (t *testing.T) {

	// create a map that associates simple relational expressions to their
	// expected result using integer constants
	expected := map[string]bool {
		"'dario' <  'adriana'" : false,
		"'dario' <= 'adriana'" : false,
		"'dario' =  'adriana'" : false,
		"'dario' != 'adriana'" : true ,
		"'dario' >= 'adriana'" : true ,
		"'dario' >  'adriana'" : true ,
		
		"'dario' <  'dario'" : false,
		"'dario' <= 'dario'" : true ,
		"'dario' =  'dario'" : true ,
		"'dario' != 'dario'" : false,
		"'dario' >= 'dario'" : true ,
		"'dario' >  'dario'" : false,
	}
	
	// -- simple relational expressions
	for expression, value := range expected {
		assert (t, expression, TypeBool (value))
	}
	
	// -- compound relational expressions

	// ---- two relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {

			// OR
			assert (t, expression1 + " or " + expression2, TypeBool (value1 || value2))

			// AND
			assert (t, expression1 + " and " + expression2, TypeBool (value1 && value2))
		}
	}
	
	// ---- three relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {

				// OR
				assert (t, expression1 + " or " + expression2 + " or " + expression3,
					TypeBool (value1 || value2 || value3))

				// AND
				assert (t, expression1 + " and " + expression2 + " and " + expression3,
					TypeBool (value1 && value2 && value3))

				// OR/AND
				assert (t, expression1 + " or " + expression2 + " and " + expression3,
					TypeBool (value1 || value2 && value3))
				assert (t, expression1 + " and " + expression2 + " or " + expression3,
					TypeBool (value1 && value2 || value3))
			}
		}
	}

	// --- four relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {
				for expression4, value4 := range expected {

					// OR
					assert (t, expression1 + " or " + expression2 + " or " + expression3 + " or " + expression4,
						TypeBool (value1 || value2 || value3 || value4))

					// AND
					assert (t, expression1 + " and " + expression2 + " and " + expression3 + " and " + expression4,
						TypeBool (value1 && value2 && value3 && value4))

					// OR/AND
					assert (t, expression1 + " or " + expression2 + " or " + expression3 + " and " + expression4,
						TypeBool (value1 || value2 || value3 && value4))
					assert (t, expression1 + " or " + expression2 + " and " + expression3 + " or " + expression4,
						TypeBool (value1 || value2 && value3 || value4))
					assert (t, expression1 + " and " + expression2 + " or " + expression3 + " or " + expression4,
						TypeBool (value1 && value2 || value3 || value4))
					assert (t, expression1 + " or " + expression2 + " and " + expression3 + " and " + expression4,
						TypeBool (value1 || value2 && value3 && value4))
					assert (t, expression1 + " and " + expression2 + " or " + expression3 + " and " + expression4,
						TypeBool (value1 && value2 || value3 && value4))
					assert (t, expression1 + " and " + expression2 + " and " + expression3 + " or " + expression4,
						TypeBool (value1 && value2 && value3 || value4))
				}
			}
		}
	}	
}

func TestParenthesized (t *testing.T) {

	// tests using constant integers
	assert (t, "(10>=1)", true)
	assert (t, " ( 3 = 4 and 5<2 )or 3>=2", true)
	assert (t, " 3 = 4 and (5<2 or 3>=2)", false)
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
