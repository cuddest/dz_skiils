package controllers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/cuddest/dz-skills/models"
	"github.com/gin-gonic/gin"
)

// SQL queries as constants
const (
	createVideoQuery = `
		INSERT INTO videos (title, link, course_id) 
		VALUES ($1, $2, $3) RETURNING id`

	getVideoQuery = `
		SELECT id, title, link, course_id 
		FROM videos WHERE id = $1`

	getAllVideosQuery = `
		SELECT id, title, link, course_id 
		FROM videos`

	getVideosByCourseQuery = `
		SELECT id, title, link, course_id 
		FROM videos WHERE course_id = $1`

	updateVideoQuery = `
		UPDATE videos 
		SET title = $1, link = $2, course_id = $3 
		WHERE id = $4`

	deleteVideoQuery = `
		DELETE FROM videos WHERE id = $1`
)

type VideoController struct {
	db *sql.DB
}

func NewVideoController(db *sql.DB) *VideoController {
	return &VideoController{db: db}
}

func (h *VideoController) validateVideo(video *models.Video) error {
	if video.Title == "" {
		return errors.New("title is required")
	}
	if video.Link == "" {
		return errors.New("link is required")
	}
	if video.CourseID <= 0 {
		return errors.New("valid course ID is required")
	}
	return nil
}

// @Summary Create a new video
// @Description Create a new video with the provided information
// @Tags videos
// @Accept json
// @Produce json
// @Param video body models.Video true "Video information"
// @Success 201 {object} models.Video
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /videos/createVideo [post]
func (h *VideoController) CreateVideo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var video models.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateVideo(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify course exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", video.CourseID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify course"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	err = h.db.QueryRowContext(ctx, createVideoQuery,
		video.Title, video.Link, video.CourseID).Scan(&video.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video"})
		return
	}

	c.JSON(http.StatusCreated, video)
}

// @Summary Get a specific video
// @Description Get a video by its ID
// @Tags videos
// @Accept json
// @Produce json
// @Param id path int true "Video ID"
// @Success 200 {object} models.Video
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /videos/get [post]
func (h *VideoController) GetVideo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var video models.Video
	err = h.db.QueryRowContext(ctx, getVideoQuery, id).Scan(
		&video.ID, &video.Title, &video.Link, &video.CourseID,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve video"})
		return
	}

	c.JSON(http.StatusOK, video)
}

// @Summary Get all videos
// @Description Retrieve all videos
// @Tags videos
// @Accept json
// @Produce json
// @Success 200 {array} models.Video
// @Failure 500 {object} map[string]interface{}
// @Router /videos/all [get]
func (h *VideoController) GetAllVideos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, getAllVideosQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve videos"})
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var video models.Video
		if err := rows.Scan(
			&video.ID, &video.Title, &video.Link, &video.CourseID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process videos"})
			return
		}
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing videos"})
		return
	}

	c.JSON(http.StatusOK, videos)
}

// @Summary Get videos by course
// @Description Get all videos belonging to a specific course
// @Tags videos
// @Accept json
// @Produce json
// @Param courseId path int true "Course ID"
// @Success 200 {array} models.Video
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /videos/GetVideosByCourse [post]
func (h *VideoController) GetVideosByCourse(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	courseID, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
		return
	}

	// Verify course exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", courseID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify course"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	rows, err := h.db.QueryContext(ctx, getVideosByCourseQuery, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve videos"})
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var video models.Video
		if err := rows.Scan(
			&video.ID, &video.Title, &video.Link, &video.CourseID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process videos"})
			return
		}
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing videos"})
		return
	}

	c.JSON(http.StatusOK, videos)
}

// @Summary Update a video
// @Description Update an existing video's information
// @Tags videos
// @Accept json
// @Produce json
// @Param id path int true "Video ID"
// @Param video body models.Video true "Updated video information"
// @Success 200 {object} models.Video
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /videos/updateVideo [put]
func (h *VideoController) UpdateVideo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var video models.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateVideo(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify course exists
	var exists bool
	err = h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", video.CourseID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify course"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	result, err := h.db.ExecContext(ctx, updateVideoQuery,
		video.Title, video.Link, video.CourseID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm update"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	video.ID = uint(id)
	c.JSON(http.StatusOK, video)
}

// @Summary Delete a video
// @Description Delete a video by its ID
// @Tags videos
// @Accept json
// @Produce json
// @Param id path int true "Video ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /videos/DeleteVideo [delete]
func (h *VideoController) DeleteVideo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.db.ExecContext(ctx, deleteVideoQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm deletion"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}
