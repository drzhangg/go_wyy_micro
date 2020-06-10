package orm

import (
	"go_wyy_micro/mysql_es/db"
	"go_wyy_micro/mysql_es/model"
)

func GetTag(name string) int64 {
	var (
		tag model.Tag
	)
	result := db.DbEngine.Table("tag_tbl").Where("name = ?", name).Find(&tag).RowsAffected
	return result
}
