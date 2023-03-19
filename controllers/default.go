package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	i "github.com/pergamenum/go-consensus-standards/interfaces"
	t "github.com/pergamenum/go-consensus-standards/types"
	r "github.com/pergamenum/go-utils-gin/responses"
)

type Controller[M, D any] struct {
	service i.Service[M]
	mapper  i.ControllerMapper[M, D]
}

type ControllerConfig[M, D any] struct {
	Service i.Service[M]
	Mapper  i.ControllerMapper[M, D]
}

func NewController[M, D any](conf ControllerConfig[M, D]) *Controller[M, D] {

	c := &Controller[M, D]{
		service: conf.Service,
		mapper:  conf.Mapper,
	}

	return c
}

func (c *Controller[M, D]) Create(ctx *gin.Context) {

	var dto D
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		r.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	model := c.mapper.FromDTO(dto)

	err = c.service.Create(ctx, model)
	if err != nil {
		r.StandardResponses(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller[M, D]) Read(ctx *gin.Context) {

	id := ctx.Param("id")

	model, err := c.service.Read(ctx, id)
	if err != nil {
		r.StandardResponses(ctx, err)
		return
	}

	dto := c.mapper.ToDTO(model)
	ctx.JSON(http.StatusOK, dto)
}

func (c *Controller[M, D]) Update(ctx *gin.Context) {

	var dto D
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		r.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	update := c.mapper.ToUpdate(dto)
	err = c.service.Update(ctx, update)
	if err != nil {
		r.StandardResponses(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller[M, D]) Delete(ctx *gin.Context) {

	id := ctx.Param("id")

	err := c.service.Delete(ctx, id)
	if err != nil {
		r.StandardResponses(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller[M, D]) Search(ctx *gin.Context) {

	q := t.Query{}
	qs, err := q.FromURL(ctx.Request.URL.Query())
	if err != nil {
		r.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	models, err := c.service.Search(ctx, qs)
	if err != nil {
		r.StandardResponses(ctx, err)
		return
	}
	// This prevents null response bodies.
	dtos := make([]D, 0)
	for _, model := range models {
		dto := c.mapper.ToDTO(model)
		dtos = append(dtos, dto)
	}
	ctx.JSON(http.StatusOK, dtos)
}
