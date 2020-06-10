package main

import (
	"github.com/gin-gonic/gin"
	"go_wyy_micro/mysql_es/router"
)

func main() {
	r := gin.Default()

	r.POST("/api/tag", router.OnNewTag)
	r.GET("/api/tag/search", router.OnSearchTag)
	r.POST("/api/tag/link_entity", router.OnLinkEntity)
	r.GET("/api/tag/entity_tags", router.OnEntityTags)

	r.Run(":9800")

}
