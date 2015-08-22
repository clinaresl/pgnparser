/* 
  tbl_test.go
  Description: Unit tests for the automated generation of tables
  ----------------------------------------------------------------------------- 

  Started on  <Mon Aug 17 18:06:33 2015 Carlos Linares Lopez>
  Last update <sábado, 22 agosto 2015 18:08:11 Carlos Linares Lopez (clinares)>
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
	"log"
	"testing"
)

func TestNewTable (t *testing.T) {

	var spec1 = "|||l|cccl|l|||"
	var spec2 = "l@{ (}r@{)}"
	var spec3 = "||l|clrrc||"
	
	table1, err1 := NewTable (spec1); if err1 != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table1.TopRule ()
	
	if table1.AddRow ([]string{"Hola", "me", "llamo", "Carlos", "Linares", "López"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	table1.MidRule ()
	
	if table1.AddRow ([]string{"", "Y", "tengo", "tres", "hijos"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	table1.MidRule ()
	
	if table1.AddRow ([]string{"", "", "Roberto", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table1.AddRow ([]string{"", "", "Dario", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table1.AddRow ([]string{"", "", "Adriana", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	table1.BottomRule ()

	fmt.Println (table1)

	table2, err2 := NewTable (spec2); if err2 != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table2.TopRule ()
	
	if table2.AddRow ([]string{"Player", "ELO"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	table2.MidRule ()

	if table2.AddRow ([]string{"clinares", "1582"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	if table2.AddRow ([]string{"nemesis", "1631"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	if table2.AddRow ([]string{"jemma", "1811"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table2.AddRow ([]string{"zco", "1880"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table2.AddRow ([]string{"miercoles", "1893"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	table2.BottomRule ()

	fmt.Println (table2)

	table3, err3 := NewTable (spec3); if err3 != nil {
		t.Fatal (" Fatal error while constructing the table")
	}
	log.Println (table3)
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
