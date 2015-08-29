/* 
  tbl_test.go
  Description: Unit tests for the automated generation of tables
  ----------------------------------------------------------------------------- 

  Started on  <Mon Aug 17 18:06:33 2015 Carlos Linares Lopez>
  Last update <sábado, 29 agosto 2015 02:16:36 Carlos Linares Lopez (clinares)>
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
	"fmt"
	"testing"
)

func TestNewTable1 (t *testing.T) {

	var spec = "|l|c|ccl|l|"
	
	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table.HSingleRule ()
	
	if table.AddRow ([]string{"Hola", "me", "llamo", "Carlos", "Linares", "López"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	table.HSingleRule ()
	
	if table.AddRow ([]string{"", "Y", "tengo", "tres", "hijos"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	table.HSingleRule ()
	
	if table.AddRow ([]string{"", "", "Roberto", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"", "", "Dario", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"", "", "Adriana", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	table.HSingleRule ()
	
	fmt.Println (table)
}

func TestNewTable2 (t *testing.T) {
	var spec = "l@{ (}r@{)}@{ }c"

	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table.TopRule ()
	
	if table.AddRow ([]string{"Player", "ELO", "Country"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	table.MidRule ()

	if table.AddRow ([]string{"clinares", "1588", "Spain"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	if table.AddRow ([]string{"nemesis", "1631", "Egypt"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	if table.AddRow ([]string{"jemma", "1811", "Germany"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"zco", "1880", "United Kingdom"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"miercoles", "1893", "Spain"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	table.BottomRule ()

	fmt.Println (table)	
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
