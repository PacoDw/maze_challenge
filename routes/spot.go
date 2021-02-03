package routes

import (
	"errors"
	"net/http"

	"github.com/PacoDw/maze_challenge/repository"
	"github.com/gin-gonic/gin"
)

// CreateSpot creates a spot.
var CreateSpot = func(c *gin.Context) {
	repo, ok := c.MustGet("mongoRepoConn").(*repository.MongoDBService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("no connection with database").Error(),
		})

		return
	}

	spot := &repository.Spot{}

	if err := c.ShouldBindJSON(spot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	id, err := repo.Spot.Create(c, spot)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	spot.ID = id

	c.JSON(http.StatusOK, spot)
}

// GetSpot creates a spot.
var GetSpot = func(c *gin.Context) {
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

	spot, err := repo.Spot.Get(c, &repository.SpotFilter{ID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, spot)
}

// UpdateSpot updates a spot.
var UpdateSpot = func(c *gin.Context) {
	repo, ok := c.MustGet("mongoRepoConn").(*repository.MongoDBService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("no connection with database").Error(),
		})

		return
	}

	spot := &repository.Spot{}

	if err := c.ShouldBindJSON(spot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	spot, err := repo.Spot.Update(c, spot)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, spot)
}

// DeleteSpot deletes a spot.
var DeleteSpot = func(c *gin.Context) {
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

	_, err := repo.Spot.Delete(c, &repository.SpotFilter{ID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, "")
}
