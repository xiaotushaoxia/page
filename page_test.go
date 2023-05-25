package page

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func TestQueryPage(t *testing.T) {
	var db *gorm.DB // init

	db = db.Where("id in (?)", []int{1, 2, 3, 4})
	db = db.Preload("locations")

	page, err := QueryPage(db, 1, 10, "id desc")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(page)
}

type area struct {
	locations []*location
}

type location struct {
	Id     int
	Name   string
	areaID int
}
