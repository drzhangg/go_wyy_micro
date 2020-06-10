package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_wyy_micro/mysql_es/db"
	"go_wyy_micro/mysql_es/es"
	"go_wyy_micro/mysql_es/model"
	"net/http"
	"strings"
)

func OnNewTag(c *gin.Context) {
	var (
		tagBody model.TagBody
		tag     model.Tag
	)
	if err := c.ShouldBindJSON(&tagBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "获取参数错误",
		})
		return
	}
	fmt.Println(tagBody)

	//判断传入的名称是否为空
	tagName := strings.TrimSpace(tagBody.Name)
	if tagName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid name",
		})
		return
	}

	result := db.DbEngine.Table("tag_tbl").Where("name = ?", tagName).Find(&tag).RowsAffected
	if result > 0 {
		// tag 已经存在
		c.JSON(http.StatusOK, gin.H{
			"tag_id": tag.Id,
		})
		return
	}

	var newTag = model.Tag{
		Name: tagName,
	}

	db.DbEngine.Table("tag_tbl").Save(&newTag)
	fmt.Println("name:", newTag)

	db.DbEngine.Table("tag_tbl").Where("name = ?", newTag.Name).Select("id").Find(&tag)
	fmt.Println("id:", tag.Id)

	//添加到es索引
	newTags := &model.Tag{
		Id:   tag.Id,
		Name: tagName,
	}
	go es.ReportTagToES(newTags)

	c.JSON(http.StatusOK, gin.H{
		"tag_id": tag.Id,
	})
}

func OnSearchTag(c *gin.Context) {
	var (
		reqBody model.SearchTagReqBody
	)
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	searchKeyword := strings.TrimSpace(reqBody.Keyword)
	if searchKeyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid keyword",
		})
		return
	}

	tags, err := es.SearchTagsFromEs(reqBody.Keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": tags,
	})
}
func OnEntityTags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"tags": []struct{}{},
	})
	return
}

func OnLinkEntity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"link_id": 0,
	})
}
