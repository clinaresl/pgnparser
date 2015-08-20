/* 
  tbl_test.go
  Description: Unit tests for the automated generation of tables
  ----------------------------------------------------------------------------- 

  Started on  <Mon Aug 17 18:06:33 2015 Carlos Linares Lopez>
  Last update <jueves, 20 agosto 2015 20:19:08 Carlos Linares Lopez (clinares)>
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

	var spec1 = "||l|cccll||"
	var spec2 = "|l|clrrc"
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

	fmt.Println ()
	fmt.Println (table1)
	fmt.Println ()

	table2, err2 := NewTable (spec2); if err2 != nil {
		t.Fatal (" Fatal error while constructing the table")
	}
	log.Println (table2)

	table3, err3 := NewTable (spec3); if err3 != nil {
		t.Fatal (" Fatal error while constructing the table")
	}
	log.Println (table3)
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
