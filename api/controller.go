package api

import (
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/route"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type Controller struct {
	Service *Service
}

func (x *Controller) Auto(r *gin.Engine) {
	r.POST("/:name", route.Use(x.Create))
	r.GET("/:name", route.Use(x.Find))
	r.GET("/:name/:id", route.Use(x.FindOneById))
	r.PATCH("/:name", route.Use(x.Update))
	r.PATCH("/:name/:id", route.Use(x.UpdateOneById))
	r.PUT("/:name/:id", route.Use(x.ReplaceOneById))
	r.DELETE("/:name/:id", route.Use(x.DeleteOneById))
}

type CreateDto struct {
	Doc bson.M `json:"doc" binding:"required"`
}

// Create 创建文档
func (x *Controller) Create(c *gin.Context) interface{} {
	name := c.Param("name")
	var body CreateDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	result, err := x.Service.Create(c.Request.Context(), name, body.Doc)
	if err != nil {
		return err
	}
	c.Set("status_code", http.StatusCreated)
	return result
}

type PaginationDto struct {
	Index int64 `header:"page" binding:"omitempty,gt=0,number"`
	Size  int64 `header:"pagesize" binding:"omitempty,oneof=10 20 50 100"`
}

type FindDto struct {
	Id     []string `form:"id" binding:"omitempty"`
	Where  bson.M   `form:"where" binding:"omitempty"`
	Sort   []string `form:"sort" binding:"omitempty"`
	Single bool     `form:"single"`
}

// Find 通过获取多个文档
func (x *Controller) Find(c *gin.Context) interface{} {
	name := c.Param("name")
	var page PaginationDto
	if err := c.ShouldBindHeader(&page); err != nil {
		return err
	}
	var query FindDto
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	if query.Single {
		result, err := x.Service.FindOne(ctx, name, query.Where)
		if err != nil {
			return err
		}
		return result
	}
	if len(query.Id) != 0 {
		result, err := x.Service.
			FindById(ctx, name, query.Id, query.Sort)
		if err != nil {
			return err
		}
		return result
	}
	if page.Index != 0 && page.Size != 0 {
		result, err := x.Service.
			FindByPage(ctx, name, page, query.Where, query.Sort)
		if err != nil {
			return err
		}
		return result
	}
	result, err := x.Service.
		Find(ctx, name, query.Where, query.Sort)
	if err != nil {
		return err
	}
	return result
}

// FindOneById 通过 ID 获取单个文档
func (x *Controller) FindOneById(c *gin.Context) interface{} {
	name := c.Param("name")
	id := c.Param("name")
	err, result := x.Service.FindOneById(c.Request.Context(), name, id)
	if err != nil {
		return err
	}
	return result
}

type UpdateQuery struct {
	Id       []string `form:"id" binding:"omitempty"`
	Where    bson.M   `form:"where" binding:"omitempty"`
	Multiple bool     `form:"multiple" binding:"omitempty"`
}

type UpdateDto struct {
	Update bson.M `json:"update" binding:"required"`
}

// Update 更新文档
func (x *Controller) Update(c *gin.Context) interface{} {
	name := c.Param("name")
	var query UpdateQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	var body UpdateDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	if len(query.Id) != 0 {
		result, err := x.Service.
			UpdateManyById(ctx, name, query.Id, body.Update)
		if err != nil {
			return err
		}
		return result
	}
	if query.Multiple {
		result, err := x.Service.
			UpdateMany(ctx, name, query.Where, body.Update)
		if err != nil {
			return err
		}
		return result
	}
	result, err := x.Service.
		UpdateOne(ctx, name, query.Where, body.Update)
	if err != nil {
		return err
	}
	return result
}

func (x *Controller) UpdateOneById(c *gin.Context) interface{} {
	name := c.Param("name")
	id := c.Param("id")
	var body UpdateDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.
		UpdateOneById(ctx, name, id, body.Update)
	if err != nil {
		return err
	}
	return result
}

type ReplaceOneDto struct {
	Doc bson.M `json:"doc" binding:"required"`
}

func (x *Controller) ReplaceOneById(c *gin.Context) interface{} {
	name := c.Param("name")
	id := c.Param("id")
	var body ReplaceOneDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.ReplaceOneById(ctx, name, id, body.Doc)
	if err != nil {
		return err
	}
	return result
}

func (x *Controller) DeleteOneById(c *gin.Context) interface{} {
	name := c.Param("name")
	id := c.Param("id")
	ctx := c.Request.Context()
	result, err := x.Service.DeleteOneById(ctx, name, id)
	if err != nil {
		return err
	}
	return result
}