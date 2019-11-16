package model

import (
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// PaginationData ...
type PaginationData struct {
	Total    uint64 `json:"total"`
	PageSize uint64 `json:"page_size"`
	Current  uint64 `json:"current"`
}

// PaginationQuery ..
type PaginationQuery struct {
	Sorter      string `json:"sorter,omitempty" form:"sorter"`
	PageSize    uint64 `json:"page_size,omitempty" form:"pageSize"`
	CurrentPage uint64 `json:"current_page,omitempty" form:"currentPage"`
}

// ListData ..
type ListData struct {
	Pagination PaginationData `json:"pagination,omitempty"`
	List       interface{}    `json:"list,omitempty"`
}

// AfterPagination ...
func AfterPagination(db *gorm.DB, pq PaginationQuery, data *ListData) *gorm.DB {

	if pq.PageSize <= 0 {
		pq.PageSize = 20
	}
	data.Pagination.PageSize = pq.PageSize

	if pq.CurrentPage <= 0 {
		pq.CurrentPage = 1
	}
	data.Pagination.Current = pq.CurrentPage

	if pq.Sorter != "" {
		index := strings.LastIndex(pq.Sorter, "_")
		if index != -1 {
			var sortDirection string
			if pq.Sorter[index:] == "_ascend" {
				sortDirection = "ASC"
			} else if pq.Sorter[index:] == "_descend" {
				sortDirection = "DESC"
			}
			if sortDirection != "" {
				db = db.Order(pq.Sorter[:index]+" "+sortDirection, true)
			}
		}
	}

	db.Count(&data.Pagination.Total)
	return db.Limit(pq.PageSize).Offset((pq.CurrentPage - 1) * pq.PageSize)
}

// BeforePagenation ...
func BeforePagenation(c *gin.Context) PaginationQuery {
	var pq PaginationQuery
	c.BindQuery(&pq)
	return pq
}

// WhereQuery ...
func WhereQuery(db *gorm.DB, data interface{}) *gorm.DB {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	for i := 0; i < v.NumField(); i++ {
		col := gorm.ToColumnName(t.Field(i).Name)
		fType := t.Field(i).Type.Name()
		switch fType {
		case "string":
			if v.Field(i).String() != "" {
				if strings.Contains(col, "status") {
					// 为包含 status 的字段设置数组索引
					db = db.Where(col+" IN (?)", strings.Split(v.Field(i).String(), ","))
				} else {
					db = db.Where(col+" LIKE ?", "%"+v.Field(i).String()+"%")
				}
			}
			continue
		case "uint8", "uint16", "uint32", "uint64":
			if v.Field(i).Uint() > 0 {
				db = db.Where(col+" = ?", v.Field(i).Uint())
			}
		case "Time":
			date := v.Field(i).Interface().(time.Time)
			if !date.IsZero() {
				db = db.Where(col+" BETWEEN ? AND ?", date, date.Add(time.Second*(23*60*60+59*60+59)))
			}
		default:
			log.Println("WhereQuery", "Unknown Type", t.Field(i).Name, v.Field(i).Interface())
		}
	}
	return db
}
