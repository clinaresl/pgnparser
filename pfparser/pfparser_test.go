/* 
  pfparser_test.go
  Description: Unit tests for the propositional formula parser package
  ----------------------------------------------------------------------------- 

  Started on  <Sun May 24 23:26:09 2015 Carlos Linares Lopez>
  Last update <domingo, 07 junio 2015 17:00:31 Carlos Linares Lopez (clinares)>
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

// the following function parses the given pformula and evaluates it using the
// specified symbol table. In case the result differs from the expected one, a
// Fatal error is raised using the testing framework specified in t
func assert (t *testing.T, pformula string,
	symtable map[string]RelationalInterface, expected TypeBool) {

	var err error
	var logicalEvaluator LogicalEvaluator

	log.Println (pformula, expected)

	// Parse the given formula
	logicalEvaluator, err = Parse (&pformula, 0); if err != nil {
		log.Fatalf ("%v\n", err)
	}

	// and now evaluate it, in case the result is not the expected
	// one then raise a fatal error
	if logicalEvaluator.Evaluate (symtable) != expected {
		t.Fatalf (" Error in pformula %v", pformula)
	}	
}

func TestConstIntegerUnparenthesized (t *testing.T) {

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

	// create an empty symbol table
	symtable := make (map[string]RelationalInterface)
	
	// -- simple relational expressions
	for expression, value := range expected {
		assert (t, expression, symtable, TypeBool (value))
	}
	
	// -- compound relational expressions

	// ---- two relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {

			// OR
			assert (t, expression1 + " or " + expression2,
				symtable, TypeBool (value1 || value2))

			// AND
			assert (t, expression1 + " and " + expression2,
				symtable, TypeBool (value1 && value2))
		}
	}
	
	// ---- three relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {

				// OR
				assert (t, expression1 + " or " + expression2 + " or " + expression3,
					symtable, TypeBool (value1 || value2 || value3))

				// AND
				assert (t, expression1 + " and " + expression2 + " and " + expression3,
					symtable, TypeBool (value1 && value2 && value3))

				// OR/AND
				assert (t, expression1 + " or " + expression2 + " and " + expression3,
					symtable, TypeBool (value1 || value2 && value3))
				assert (t, expression1 + " and " + expression2 + " or " + expression3,
					symtable, TypeBool (value1 && value2 || value3))
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
						symtable, TypeBool (value1 || value2 || value3 || value4))

					// AND
					assert (t, expression1 + " and " + expression2 + " and " + expression3 + " and " + expression4,
						symtable, TypeBool (value1 && value2 && value3 && value4))

					// OR/AND
					assert (t, expression1 + " or " + expression2 + " or " + expression3 + " and " + expression4,
						symtable, TypeBool (value1 || value2 || value3 && value4))
					assert (t, expression1 + " or " + expression2 + " and " + expression3 + " or " + expression4,
						symtable, TypeBool (value1 || value2 && value3 || value4))
					assert (t, expression1 + " and " + expression2 + " or " + expression3 + " or " + expression4,
						symtable, TypeBool (value1 && value2 || value3 || value4))
					assert (t, expression1 + " or " + expression2 + " and " + expression3 + " and " + expression4,
						symtable, TypeBool (value1 || value2 && value3 && value4))
					assert (t, expression1 + " and " + expression2 + " or " + expression3 + " and " + expression4,
						symtable, TypeBool (value1 && value2 || value3 && value4))
					assert (t, expression1 + " and " + expression2 + " and " + expression3 + " or " + expression4,
						symtable, TypeBool (value1 && value2 && value3 || value4))
				}
			}
		}
	}
}

func TestConstStringUnparenthesized (t *testing.T) {

	// create a map that associates simple relational expressions to their
	// expected result using integer constants
	expected := map[string]bool {
		"'dario' <  'adriana'" : false,
		"'dario' <= 'adriana'" : false,
		"'dario' =  'adriana'" : false,
		"'dario' != 'adriana'" : true ,
		"'dario' >= 'adriana'" : true ,
		"'dario' >  'adriana'" : true ,
		"'dario' in 'adriana'" : false,
		"'dario' not_in 'adriana'" : true,
		
		"'dario' <  'dario'" : false,
		"'dario' <= 'dario'" : true ,
		"'dario' =  'dario'" : true ,
		"'dario' != 'dario'" : false,
		"'dario' >= 'dario'" : true ,
		"'dario' >  'dario'" : false,
		"'dario' in 'dario'" : true,
		"'dario' not_in 'dario'" : false,
	}
	
	// create an empty symbol table
	symtable := make (map[string]RelationalInterface)
	
	// -- simple relational expressions
	for expression, value := range expected {
		assert (t, expression, symtable, TypeBool (value))
	}
	
	// -- compound relational expressions

	// ---- two relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {

			// OR
			assert (t, expression1 + " or " + expression2,
				symtable, TypeBool (value1 || value2))

			// AND
			assert (t, expression1 + " and " + expression2,
				symtable, TypeBool (value1 && value2))
		}
	}
	
	// ---- three relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {

				// OR
				assert (t, expression1 + " or " + expression2 + " or " + expression3,
					symtable, TypeBool (value1 || value2 || value3))

				// AND
				assert (t, expression1 + " and " + expression2 + " and " + expression3,
					symtable, TypeBool (value1 && value2 && value3))

				// OR/AND
				assert (t, expression1 + " or " + expression2 + " and " + expression3,
					symtable, TypeBool (value1 || value2 && value3))
				assert (t, expression1 + " and " + expression2 + " or " + expression3,
					symtable, TypeBool (value1 && value2 || value3))
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
						symtable, TypeBool (value1 || value2 || value3 || value4))
					
					// AND
					assert (t, expression1 + " and " + expression2 + " and " + expression3 + " and " + expression4,
						symtable, TypeBool (value1 && value2 && value3 && value4))

					// OR/AND
					assert (t, expression1 + " or " + expression2 + " or " + expression3 + " and " + expression4,
						symtable, TypeBool (value1 || value2 || value3 && value4))
					assert (t, expression1 + " or " + expression2 + " and " + expression3 + " or " + expression4,
						symtable, TypeBool (value1 || value2 && value3 || value4))
					assert (t, expression1 + " and " + expression2 + " or " + expression3 + " or " + expression4,
						symtable, TypeBool (value1 && value2 || value3 || value4))
					assert (t, expression1 + " or " + expression2 + " and " + expression3 + " and " + expression4,
						symtable, TypeBool (value1 || value2 && value3 && value4))
					assert (t, expression1 + " and " + expression2 + " or " + expression3 + " and " + expression4,
						symtable, TypeBool (value1 && value2 || value3 && value4))
					assert (t, expression1 + " and " + expression2 + " and " + expression3 + " or " + expression4,
						symtable, TypeBool (value1 && value2 && value3 || value4))
				}
			}
		}
	}	
}

func TestConstIntegerParenthesized (t *testing.T) {

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
	
	// create an empty symbol table
	symtable := make (map[string]RelationalInterface)
	
	// -- simple relational expressions
	for expression, value := range expected {
		assert (t, "(" + expression + ")", symtable, TypeBool (value))
	}
	
	// -- compound relational expressions

	// ---- two relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {

			// OR
			assert (t, "(" + expression1 + " or " + expression2 + ")",
				symtable, TypeBool (value1 || value2))

			// AND
			assert (t, "(" + expression1 + " and " + expression2 + ")",
				symtable, TypeBool (value1 && value2))
		}
	}
	
	// ---- three relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {

				// OR
				assert (t, "(" + expression1 + " or " + expression2 + ")" + " or " + expression3,
					symtable, TypeBool ((value1 || value2) || value3))
				assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + ")",
					symtable, TypeBool (value1 || (value2 || value3)))

				// AND
				assert (t, "(" + expression1 + " and " + expression2 + ")"+ " and " + expression3,
					symtable, TypeBool ((value1 && value2) && value3))
				assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + ")",
					symtable, TypeBool (value1 && (value2 && value3)))

				// OR/AND
				assert (t, "(" + expression1 + " or " + expression2 + ")"+ " and " + expression3,
					symtable, TypeBool ((value1 || value2) && value3))
				assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")",
					symtable, TypeBool (value1 || (value2 && value3)))
				assert (t, "(" + expression1 + " and " + expression2 + ")"+ " or " + expression3,
					symtable, TypeBool ((value1 && value2) || value3))
				assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")",
					symtable, TypeBool (value1 && (value2 || value3)))
			}
		}
	}

	// --- four relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {
				for expression4, value4 := range expected {

					// expressions with only ORs or only
					// ANDs are avoided here.

					// OR/AND

					// (p || q) || r && s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " or " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 || value2) || value3 && value4))

					// p || (q || r) && s
					assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 || (value2 || value3) && value4))
					
					// p || q || (r && s)
					assert (t, expression1 + " or " + expression2 + " or " + "(" + expression3 + " and " + expression4 + ")" ,
						symtable, TypeBool (value1 || value2 || (value3 && value4)))

					// (p || q || r) && s
					assert (t, "(" + expression1 + " or " + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 || value2 || value3) && value4))

					// p || (q || r && s)
					assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + " and " + expression4 + ")" ,
						symtable, TypeBool (value1 || (value2 || value3 && value4)))


					// (p || q) && r || s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " and " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 || value2) && value3 || value4))
					
					// p || (q && r) || s
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 || (value2 && value3) || value4))
					
					// p || q && (r || s)
					assert (t, expression1 + " or " + expression2 + " and " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 || value2 && (value3 || value4)))
					
					// (p || q && r) || s
					assert (t, "(" + expression1 + " or " + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 || value2 && value3) || value4))
					
					// p || (q && r || s)
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 || (value2 && value3 || value4)))


					// (p && q) || r || s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " or " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 && value2) || value3 || value4))
					
					// p && (q || r) || s
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 && (value2 || value3) || value4))
					
					// p && q || (r || s)
					assert (t, expression1 + " and " + expression2 + " or " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && value2 || (value3 || value4)))
					
					// (p && q || r) || s
					assert (t, "(" + expression1 + " and " + expression2 + " or " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 && value2 || value3) || value4))
					
					// p && (q || r || s)
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + " or " + expression4 + ")",
						symtable, TypeBool (value1 && (value2 || value3 || value4)))


					// (p || q) && r && s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " and " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 || value2) && value3 && value4))
					
					// p || (q && r) && s
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 || (value2 && value3) && value4))
					
					// p || q && (r && s)
					assert (t, expression1 + " or " + expression2 + " and " + "(" + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 || value2 && (value3 && value4)))
					
					// (p || q && r) && s
					assert (t, "(" + expression1 + " or " + expression2 + " and " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 || value2 && value3) && value4))
					
					// p || (q && r && s)
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 || (value2 && value3 && value4)))



					// (p && q) || r && s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " or " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 && value2) || value3 && value4))
					
					// p && (q || r) && s
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 && (value2 || value3) && value4))
					
					// p && q || (r && s)
					assert (t, expression1 + " and " + expression2 + " or " + "(" + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 && value2 || (value3 && value4)))
					
					// (p && q || r) && s
					assert (t, "(" + expression1 + " and " + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 && value2 || value3) && value4))
					
					// p && (q || r && s)
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 && (value2 || value3 && value4)))



					// (p && q) && r || s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " and " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 && value2) && value3 || value4))
					
					// p && (q && r) || s
					assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 && (value2 && value3) || value4))
					
					// p && q && (r || s)
					assert (t, expression1 + " and " + expression2 + " and " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && value2 && (value3 || value4)))
					
					// (p && q && r) || s
					assert (t, "(" + expression1 + " and " + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 && value2 && value3) || value4))
					
					// p && (q && r || s)
					assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && (value2 && value3 || value4)))

				}
			}
		}
	}
}

func TestConstStringParenthesized (t *testing.T) {

	// create a map that associates simple relational expressions to their
	// expected result using integer constants
	expected := map[string]bool {
		"'dario' <  'adriana'" : false,
		"'dario' <= 'adriana'" : false,
		"'dario' =  'adriana'" : false,
		"'dario' != 'adriana'" : true ,
		"'dario' >= 'adriana'" : true ,
		"'dario' >  'adriana'" : true ,
		"'dario' in 'adriana'" : false,
		"'dario' not_in 'adriana'" : true,
		
		"'dario' <  'dario'" : false,
		"'dario' <= 'dario'" : true ,
		"'dario' =  'dario'" : true ,
		"'dario' != 'dario'" : false,
		"'dario' >= 'dario'" : true ,
		"'dario' >  'dario'" : false,
		"'dario' in 'dario'" : true,
		"'dario' not_in 'dario'" : false,
	}
	
	// create an empty symbol table
	symtable := make (map[string]RelationalInterface)
	
	// -- simple relational expressions
	for expression, value := range expected {
		assert (t, "(" + expression + ")", symtable, TypeBool (value))
	}
	
	// -- compound relational expressions

	// ---- two relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {

			// OR
			assert (t, "(" + expression1 + " or " + expression2 + ")",
				symtable, TypeBool (value1 || value2))

			// AND
			assert (t, "(" + expression1 + " and " + expression2 + ")",
				symtable, TypeBool (value1 && value2))
		}
	}
	
	// ---- three relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {

				// OR
				assert (t, "(" + expression1 + " or " + expression2 + ")" + " or " + expression3,
					symtable, TypeBool ((value1 || value2) || value3))
				assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + ")",
					symtable, TypeBool (value1 || (value2 || value3)))

				// AND
				assert (t, "(" + expression1 + " and " + expression2 + ")"+ " and " + expression3,
					symtable, TypeBool ((value1 && value2) && value3))
				assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + ")",
					symtable, TypeBool (value1 && (value2 && value3)))

				// OR/AND
				assert (t, "(" + expression1 + " or " + expression2 + ")"+ " and " + expression3,
					symtable, TypeBool ((value1 || value2) && value3))
				assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")",
					symtable, TypeBool (value1 || (value2 && value3)))
				assert (t, "(" + expression1 + " and " + expression2 + ")"+ " or " + expression3,
					symtable, TypeBool ((value1 && value2) || value3))
				assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")",
					symtable, TypeBool (value1 && (value2 || value3)))
			}
		}
	}

	// --- four relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {
				for expression4, value4 := range expected {

					// expressions with only ORs or only
					// ANDs are avoided here.

					// OR/AND

					// (p || q) || r && s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " or " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 || value2) || value3 && value4))

					// p || (q || r) && s
					assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 || (value2 || value3) && value4))
					
					// p || q || (r && s)
					assert (t, expression1 + " or " + expression2 + " or " + "(" + expression3 + " and " + expression4 + ")" ,
						symtable, TypeBool (value1 || value2 || (value3 && value4)))

					// (p || q || r) && s
					assert (t, "(" + expression1 + " or " + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 || value2 || value3) && value4))

					// p || (q || r && s)
					assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + " and " + expression4 + ")" ,
						symtable, TypeBool (value1 || (value2 || value3 && value4)))


					// (p || q) && r || s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " and " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 || value2) && value3 || value4))
					
					// p || (q && r) || s
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 || (value2 && value3) || value4))
					
					// p || q && (r || s)
					assert (t, expression1 + " or " + expression2 + " and " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 || value2 && (value3 || value4)))
					
					// (p || q && r) || s
					assert (t, "(" + expression1 + " or " + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 || value2 && value3) || value4))
					
					// p || (q && r || s)
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 || (value2 && value3 || value4)))


					// (p && q) || r || s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " or " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 && value2) || value3 || value4))
					
					// p && (q || r) || s
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 && (value2 || value3) || value4))
					
					// p && q || (r || s)
					assert (t, expression1 + " and " + expression2 + " or " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && value2 || (value3 || value4)))
					
					// (p && q || r) || s
					assert (t, "(" + expression1 + " and " + expression2 + " or " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 && value2 || value3) || value4))
					
					// p && (q || r || s)
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + " or " + expression4 + ")",
						symtable, TypeBool (value1 && (value2 || value3 || value4)))


					// (p || q) && r && s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " and " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 || value2) && value3 && value4))
					
					// p || (q && r) && s
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 || (value2 && value3) && value4))
					
					// p || q && (r && s)
					assert (t, expression1 + " or " + expression2 + " and " + "(" + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 || value2 && (value3 && value4)))
					
					// (p || q && r) && s
					assert (t, "(" + expression1 + " or " + expression2 + " and " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 || value2 && value3) && value4))
					
					// p || (q && r && s)
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 || (value2 && value3 && value4)))



					// (p && q) || r && s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " or " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 && value2) || value3 && value4))
					
					// p && (q || r) && s
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 && (value2 || value3) && value4))
					
					// p && q || (r && s)
					assert (t, expression1 + " and " + expression2 + " or " + "(" + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 && value2 || (value3 && value4)))
					
					// (p && q || r) && s
					assert (t, "(" + expression1 + " and " + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 && value2 || value3) && value4))
					
					// p && (q || r && s)
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 && (value2 || value3 && value4)))



					// (p && q) && r || s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " and " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 && value2) && value3 || value4))
					
					// p && (q && r) || s
					assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 && (value2 && value3) || value4))
					
					// p && q && (r || s)
					assert (t, expression1 + " and " + expression2 + " and " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && value2 && (value3 || value4)))
					
					// (p && q && r) || s
					assert (t, "(" + expression1 + " and " + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 && value2 && value3) || value4))
					
					// p && (q && r || s)
					assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && (value2 && value3 || value4)))

				}
			}
		}
	}
}

func TestVarIntegerParenthesized (t *testing.T) {

	// create a map that associates simple relational expressions to their
	// expected result using integer constants
	expected := map[string]bool {
		"%var1 <  %var2" : true ,
		"%var1 <= %var2" : true ,
		"%var1 =  %var2" : false,
		"%var1 != %var2" : true ,
		"%var1 >= %var2" : false,
		"%var1 >  %var2" : false,
		
		"%var1 <  %var1" : false,
		"%var1 <= %var1" : true ,
		"%var1 =  %var1" : true ,
		"%var1 != %var1" : false,
		"%var1 >= %var1" : true ,
		"%var1 >  %var1" : false,
	}
	
	// create a symbol table
	symtable := make (map[string]RelationalInterface)
	symtable["var1"] = ConstInteger(3)
	symtable["var2"] = ConstInteger(7)
	
	// -- simple relational expressions
	for expression, value := range expected {
		assert (t, "(" + expression + ")", symtable, TypeBool (value))
	}
	
	// -- compound relational expressions

	// ---- two relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {

			// OR
			assert (t, "(" + expression1 + " or " + expression2 + ")",
				symtable, TypeBool (value1 || value2))

			// AND
			assert (t, "(" + expression1 + " and " + expression2 + ")",
				symtable, TypeBool (value1 && value2))
		}
	}
	
	// ---- three relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {

				// OR
				assert (t, "(" + expression1 + " or " + expression2 + ")" + " or " + expression3,
					symtable, TypeBool ((value1 || value2) || value3))
				assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + ")",
					symtable, TypeBool (value1 || (value2 || value3)))

				// AND
				assert (t, "(" + expression1 + " and " + expression2 + ")"+ " and " + expression3,
					symtable, TypeBool ((value1 && value2) && value3))
				assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + ")",
					symtable, TypeBool (value1 && (value2 && value3)))

				// OR/AND
				assert (t, "(" + expression1 + " or " + expression2 + ")"+ " and " + expression3,
					symtable, TypeBool ((value1 || value2) && value3))
				assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")",
					symtable, TypeBool (value1 || (value2 && value3)))
				assert (t, "(" + expression1 + " and " + expression2 + ")"+ " or " + expression3,
					symtable, TypeBool ((value1 && value2) || value3))
				assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")",
					symtable, TypeBool (value1 && (value2 || value3)))
			}
		}
	}

	// --- four relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {
				for expression4, value4 := range expected {

					// expressions with only ORs or only
					// ANDs are avoided here.

					// OR/AND

					// (p || q) || r && s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " or " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 || value2) || value3 && value4))

					// p || (q || r) && s
					assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 || (value2 || value3) && value4))
					
					// p || q || (r && s)
					assert (t, expression1 + " or " + expression2 + " or " + "(" + expression3 + " and " + expression4 + ")" ,
						symtable, TypeBool (value1 || value2 || (value3 && value4)))

					// (p || q || r) && s
					assert (t, "(" + expression1 + " or " + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 || value2 || value3) && value4))

					// p || (q || r && s)
					assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + " and " + expression4 + ")" ,
						symtable, TypeBool (value1 || (value2 || value3 && value4)))


					// (p || q) && r || s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " and " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 || value2) && value3 || value4))
					
					// p || (q && r) || s
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 || (value2 && value3) || value4))
					
					// p || q && (r || s)
					assert (t, expression1 + " or " + expression2 + " and " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 || value2 && (value3 || value4)))
					
					// (p || q && r) || s
					assert (t, "(" + expression1 + " or " + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 || value2 && value3) || value4))
					
					// p || (q && r || s)
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 || (value2 && value3 || value4)))


					// (p && q) || r || s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " or " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 && value2) || value3 || value4))
					
					// p && (q || r) || s
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 && (value2 || value3) || value4))
					
					// p && q || (r || s)
					assert (t, expression1 + " and " + expression2 + " or " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && value2 || (value3 || value4)))
					
					// (p && q || r) || s
					assert (t, "(" + expression1 + " and " + expression2 + " or " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 && value2 || value3) || value4))
					
					// p && (q || r || s)
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + " or " + expression4 + ")",
						symtable, TypeBool (value1 && (value2 || value3 || value4)))


					// (p || q) && r && s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " and " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 || value2) && value3 && value4))
					
					// p || (q && r) && s
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 || (value2 && value3) && value4))
					
					// p || q && (r && s)
					assert (t, expression1 + " or " + expression2 + " and " + "(" + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 || value2 && (value3 && value4)))
					
					// (p || q && r) && s
					assert (t, "(" + expression1 + " or " + expression2 + " and " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 || value2 && value3) && value4))
					
					// p || (q && r && s)
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 || (value2 && value3 && value4)))



					// (p && q) || r && s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " or " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 && value2) || value3 && value4))
					
					// p && (q || r) && s
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 && (value2 || value3) && value4))
					
					// p && q || (r && s)
					assert (t, expression1 + " and " + expression2 + " or " + "(" + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 && value2 || (value3 && value4)))
					
					// (p && q || r) && s
					assert (t, "(" + expression1 + " and " + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 && value2 || value3) && value4))
					
					// p && (q || r && s)
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 && (value2 || value3 && value4)))



					// (p && q) && r || s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " and " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 && value2) && value3 || value4))
					
					// p && (q && r) || s
					assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 && (value2 && value3) || value4))
					
					// p && q && (r || s)
					assert (t, expression1 + " and " + expression2 + " and " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && value2 && (value3 || value4)))
					
					// (p && q && r) || s
					assert (t, "(" + expression1 + " and " + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 && value2 && value3) || value4))
					
					// p && (q && r || s)
					assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && (value2 && value3 || value4)))

				}
			}
		}
	}
}

func TestVarStringParenthesized (t *testing.T) {

	// create a map that associates simple relational expressions to their
	// expected result using integer constants
	expected := map[string]bool {
		"%var1 <  %var2" : true ,
		"%var1 <= %var2" : true ,
		"%var1 =  %var2" : false,
		"%var1 != %var2" : true ,
		"%var1 >= %var2" : false,
		"%var1 >  %var2" : false,
		"%var1 in %var2" : false,
		"%var1 not_in %var2": true,
		
		"%var1 <  %var1" : false,
		"%var1 <= %var1" : true ,
		"%var1 =  %var1" : true ,
		"%var1 != %var1" : false,
		"%var1 >= %var1" : true ,
		"%var1 >  %var1" : false,
		"%var1 in %var1" : true,
		"%var1 not_in %var1": false,
	}
	
	// create a symbol table
	symtable := make (map[string]RelationalInterface)
	symtable["var1"] = ConstString("Monica")
	symtable["var2"] = ConstString("Roberto")
	
	// -- simple relational expressions
	for expression, value := range expected {
		assert (t, "(" + expression + ")", symtable, TypeBool (value))
	}
	
	// -- compound relational expressions

	// ---- two relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {

			// OR
			assert (t, "(" + expression1 + " or " + expression2 + ")",
				symtable, TypeBool (value1 || value2))

			// AND
			assert (t, "(" + expression1 + " and " + expression2 + ")",
				symtable, TypeBool (value1 && value2))
		}
	}
	
	// ---- three relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {

				// OR
				assert (t, "(" + expression1 + " or " + expression2 + ")" + " or " + expression3,
					symtable, TypeBool ((value1 || value2) || value3))
				assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + ")",
					symtable, TypeBool (value1 || (value2 || value3)))

				// AND
				assert (t, "(" + expression1 + " and " + expression2 + ")"+ " and " + expression3,
					symtable, TypeBool ((value1 && value2) && value3))
				assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + ")",
					symtable, TypeBool (value1 && (value2 && value3)))

				// OR/AND
				assert (t, "(" + expression1 + " or " + expression2 + ")"+ " and " + expression3,
					symtable, TypeBool ((value1 || value2) && value3))
				assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")",
					symtable, TypeBool (value1 || (value2 && value3)))
				assert (t, "(" + expression1 + " and " + expression2 + ")"+ " or " + expression3,
					symtable, TypeBool ((value1 && value2) || value3))
				assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")",
					symtable, TypeBool (value1 && (value2 || value3)))
			}
		}
	}

	// --- four relational expressions
	for expression1, value1 := range expected {
		for expression2, value2 := range expected {
			for expression3, value3 := range expected {
				for expression4, value4 := range expected {

					// expressions with only ORs or only
					// ANDs are avoided here.

					// OR/AND

					// (p || q) || r && s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " or " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 || value2) || value3 && value4))

					// p || (q || r) && s
					assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 || (value2 || value3) && value4))
					
					// p || q || (r && s)
					assert (t, expression1 + " or " + expression2 + " or " + "(" + expression3 + " and " + expression4 + ")" ,
						symtable, TypeBool (value1 || value2 || (value3 && value4)))

					// (p || q || r) && s
					assert (t, "(" + expression1 + " or " + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 || value2 || value3) && value4))

					// p || (q || r && s)
					assert (t, expression1 + " or " + "(" + expression2 + " or " + expression3 + " and " + expression4 + ")" ,
						symtable, TypeBool (value1 || (value2 || value3 && value4)))


					// (p || q) && r || s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " and " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 || value2) && value3 || value4))
					
					// p || (q && r) || s
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 || (value2 && value3) || value4))
					
					// p || q && (r || s)
					assert (t, expression1 + " or " + expression2 + " and " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 || value2 && (value3 || value4)))
					
					// (p || q && r) || s
					assert (t, "(" + expression1 + " or " + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 || value2 && value3) || value4))
					
					// p || (q && r || s)
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 || (value2 && value3 || value4)))


					// (p && q) || r || s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " or " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 && value2) || value3 || value4))
					
					// p && (q || r) || s
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 && (value2 || value3) || value4))
					
					// p && q || (r || s)
					assert (t, expression1 + " and " + expression2 + " or " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && value2 || (value3 || value4)))
					
					// (p && q || r) || s
					assert (t, "(" + expression1 + " and " + expression2 + " or " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 && value2 || value3) || value4))
					
					// p && (q || r || s)
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + " or " + expression4 + ")",
						symtable, TypeBool (value1 && (value2 || value3 || value4)))


					// (p || q) && r && s
					assert (t, "(" + expression1 + " or " + expression2 + ")" + " and " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 || value2) && value3 && value4))
					
					// p || (q && r) && s
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 || (value2 && value3) && value4))
					
					// p || q && (r && s)
					assert (t, expression1 + " or " + expression2 + " and " + "(" + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 || value2 && (value3 && value4)))
					
					// (p || q && r) && s
					assert (t, "(" + expression1 + " or " + expression2 + " and " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 || value2 && value3) && value4))
					
					// p || (q && r && s)
					assert (t, expression1 + " or " + "(" + expression2 + " and " + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 || (value2 && value3 && value4)))



					// (p && q) || r && s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " or " + expression3 + " and " + expression4,
						symtable, TypeBool ((value1 && value2) || value3 && value4))
					
					// p && (q || r) && s
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool (value1 && (value2 || value3) && value4))
					
					// p && q || (r && s)
					assert (t, expression1 + " and " + expression2 + " or " + "(" + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 && value2 || (value3 && value4)))
					
					// (p && q || r) && s
					assert (t, "(" + expression1 + " and " + expression2 + " or " + expression3 + ")" + " and " + expression4,
						symtable, TypeBool ((value1 && value2 || value3) && value4))
					
					// p && (q || r && s)
					assert (t, expression1 + " and " + "(" + expression2 + " or " + expression3 + " and " + expression4 + ")",
						symtable, TypeBool (value1 && (value2 || value3 && value4)))



					// (p && q) && r || s
					assert (t, "(" + expression1 + " and " + expression2 + ")" + " and " + expression3 + " or " + expression4,
						symtable, TypeBool ((value1 && value2) && value3 || value4))
					
					// p && (q && r) || s
					assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool (value1 && (value2 && value3) || value4))
					
					// p && q && (r || s)
					assert (t, expression1 + " and " + expression2 + " and " + "(" + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && value2 && (value3 || value4)))
					
					// (p && q && r) || s
					assert (t, "(" + expression1 + " and " + expression2 + " and " + expression3 + ")" + " or " + expression4,
						symtable, TypeBool ((value1 && value2 && value3) || value4))
					
					// p && (q && r || s)
					assert (t, expression1 + " and " + "(" + expression2 + " and " + expression3 + " or " + expression4 + ")" ,
						symtable, TypeBool (value1 && (value2 && value3 || value4)))

				}
			}
		}
	}
}




/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
