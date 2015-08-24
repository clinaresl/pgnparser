/* 
  tbl_test.go
  Description: Unit tests for the automated generation of tables
  ----------------------------------------------------------------------------- 

  Started on  <Mon Aug 17 18:06:33 2015 Carlos Linares Lopez>
  Last update <lunes, 24 agosto 2015 13:11:00 Carlos Linares Lopez (clinares)>
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
	_"fmt"
	"log"
	"testing"
)

func TestNewTable1 (t *testing.T) {

	var spec = "|||l|cccl|l@{(}r@{)}|||"
	
	table, err := NewTable (spec); if err != nil {
		t.Fatal (" Fatal error while constructing the table")
	}

	table.TopRule ()
	
	if table.AddRow ([]string{"Hola", "me", "llamo", "Carlos", "Linares", "López"})!= nil {
		t.Fatal ("Error adding a new row")
	}

	table.MidRule ()
	
	if table.AddRow ([]string{"", "Y", "tengo", "tres", "hijos"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	table.MidRule ()
	
	if table.AddRow ([]string{"", "", "Roberto", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"", "", "Dario", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}
	
	if table.AddRow ([]string{"", "", "Adriana", "Linares", "Rollán"}) != nil {
		t.Fatal ("Error adding a new row")
	}

	table.BottomRule ()
	
	log.Println (table)
	log.Println (&table)
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
