package routes

import (
	"errors"
	"net/http"

	"github.com/PacoDw/maze_challenge/repository"
	"github.com/gin-gonic/gin"
)

// CreateQuadrant creates a quadrant.
var CreateQuadrant = func(c *gin.Context) {
	repo, ok := c.MustGet("mongoRepoConn").(*repository.MongoDBService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("no connection with database").Error(),
		})

		return
	}

	quadrant := &repository.Quadrant{}

	if err := c.ShouldBindJSON(quadrant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	id, err := repo.Quadrant.Create(c, quadrant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	quadrant.ID = id

	c.JSON(http.StatusOK, quadrant)
}

// GetQuadrant creates a quadrant.
var GetQuadrant = func(c *gin.Context) {
	repo, ok := c.MustGet("mongoRepoConn").(*repository.MongoDBService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("no connection with database").Error(),
		})

		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("param :id must not be empty")})

		return
	}

	quadrant, err := repo.Quadrant.Get(c, &repository.QuadrantFilter{ID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, quadrant)
}

// UpdateQuadrant updates a quadrant.
var UpdateQuadrant = func(c *gin.Context) {
	repo, ok := c.MustGet("mongoRepoConn").(*repository.MongoDBService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("no connection with database").Error(),
		})

		return
	}

	quadrant := &repository.Quadrant{}

	if err := c.ShouldBindJSON(quadrant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	quadrant, err := repo.Quadrant.Update(c, quadrant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, quadrant)
}

// DeleteQuadrant deletes a quadrant.
var DeleteQuadrant = func(c *gin.Context) {
	repo, ok := c.MustGet("mongoRepoConn").(*repository.MongoDBService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("no connection with database").Error(),
		})

		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("param :id must not be empty")})

		return
	}

	_, err := repo.Quadrant.Delete(c, &repository.QuadrantFilter{ID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, "{}")
}
