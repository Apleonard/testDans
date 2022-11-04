package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"testDans/jwt"
	"testDans/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var whiteListTokens = make([]string, 5)

type Handlers interface {
	Register(c *gin.Context)
	Check(c *gin.Context)
	Login(c *gin.Context)

	// job
	GetJobList(c *gin.Context)
	GetJobDetail(c *gin.Context)
}

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Handlers {
	return &handler{
		db: db,
	}
}
func (h *handler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hae"})
}

func (h *handler) Register(c *gin.Context) {
	var input models.UserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	passwordString := string(hashedPassword)

	u := models.User{
		Password: passwordString,
		Username: html.EscapeString(strings.TrimSpace(input.Username)),
	}

	err = h.db.Create(&u).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

func (h *handler) Login(c *gin.Context) {
	var input models.UserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	err := h.db.Model(models.User{}).Where("username = ?", input.Username).Take(&u).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(input.Password), []byte(u.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	token, err := jwt.GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token."})
		return
	}

	whiteListTokens = append(whiteListTokens, string(token))

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *handler) GetJobList(c *gin.Context) {
	var params *models.JobListParams

	if err := c.BindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := "http://dev3.dansmultipro.co.id/api/recruitment/positions.json?"

	if params.Description != "" {
		url = url + models.Description + params.Description
	}

	if params.Location != "" {
		url = url + models.Location + params.Location
	}

	if params.IsFullTime == true {
		url = url + models.IsFullTime
	}

	if params.Page > 0 {
		url = url + models.Page + strconv.Itoa(params.Page)
	}

	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	res := new(models.JobList)
	err = json.Unmarshal(resBody, &res)

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *handler) GetJobDetail(c *gin.Context) {
	var params models.JobDetailParam

	if err := c.BindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := "http://dev3.dansmultipro.co.id/api/recruitment/positions/" + params.ID
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	res := new(models.JobDetail)

	err = json.Unmarshal(resBody, &res)

	c.JSON(http.StatusOK, gin.H{"data": res})
}
