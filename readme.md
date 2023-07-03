# 基于gorm的简单分页工具

用法示例
```go
package page

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB // init

type area struct {
	Id        int
	Name      string
	Size      int
	locations []*location
}

type location struct {
	Id     int
	Name   string
	areaID int
}

func TestQueryPage(t *testing.T) {
    // 根据自己需要给db填入条件
	tx := db.Where("id in (?)", []int{1, 2, 3, 4})
	tx = db.Preload("locations")

	page, err := Query[area](tx, 1, 10, "id desc")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(page)
}

func TestGinQuery(t *testing.T) {
	// GinQuery是对Query的简单封装。（因为我个人gin用的多，所以有这个东西）
	// RequestURI: /area/query?page=2&size=10&sort=id desc,name,size
	// sort里的排序条件和sql语法一致
}

func QueryArea(c *gin.Context) {
	query, err := GinQuery[area](db, c)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, query)
}
```