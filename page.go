package page

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Result[T any] struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	PageSize    int `json:"page_size"`
	Total       int `json:"total"`
	Data        []T `json:"data"`
}

func GinQuery[T any](db *gorm.DB, c *gin.Context) (*Result[T], error) {
	return Query[T](db, c.Query("page"), c.Query("size"), c.Query("sort"))
}

func Query[T any](db *gorm.DB, pageI, sizeI any, order string) (*Result[T], error) {
	page, err := toInt(pageI)
	if err != nil {
		return nil, err
	}
	if page == 0 {
		page = 1
	}
	size, err := toInt(sizeI)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		size = 20
	}

	result := Result[T]{
		CurrentPage: page,
		PageSize:    size,
	}
	if db.Statement.Model == nil {
		var t T
		db = db.Model(t)
	}
	if order != "" {
		db = db.Order(order)
	}
	var total int64
	err = db.Count(&total).Limit(result.PageSize).Offset(result.PageSize * (result.CurrentPage - 1)).Find(&result.Data).Error
	if err != nil {
		return nil, err
	}
	result.Total = int(total)
	totalPage := result.Total / result.PageSize
	if result.Total%result.PageSize > 0 {
		totalPage++
	}
	result.TotalPage = totalPage
	return &result, nil
}

func toInt(i any) (int, error) {
	switch i.(type) {
	case string:
		s := i.(string)
		if len(s) == 0 {
			return 0, nil
		} else {
			i2, err := strconv.Atoi(s)
			if err != nil {
				return 0, fmt.Errorf("num format error: %v", i)
			}
			return i2, nil
		}
	case int:
		return i.(int), nil
	case int32:
		return int(i.(int32)), nil
	default:
		return 0, fmt.Errorf("input type error not int or string: %T", i)
	}
}
