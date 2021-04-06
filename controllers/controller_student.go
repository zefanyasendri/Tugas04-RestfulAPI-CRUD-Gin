package controllers

import (
	"fmt"
	"net/http"

	"github.com/Tugas04-RestfulAPI-CRUD-Gin/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Sebagai Input JSON
type StudentInput struct {
	NIM     string `json:"nim"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

// Read Data Students / Get All Students
func ReadDataStudent(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var mhs []models.Student
	db.Find(&mhs)
	c.JSON(http.StatusOK, gin.H{"data": mhs})
}

// Read Data One Student / Get One User
func ReadDataOneStudent(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	name := c.Params.ByName("name")
	var mhs models.Student
	if err := db.Where("name = ?", name).First(&mhs).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	db.Find(&mhs)
	c.JSON(http.StatusOK, gin.H{"data": mhs})
}

// Create/Input Data Student
func CreateDataStudent(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validasi Input
	var dataInput StudentInput

	// Cek yang di input data JSON atau bukan
	if err := c.ShouldBindJSON(&dataInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Proses Input Data
	mhs := models.Student{
		//Sesuai Model
		NIM:     dataInput.NIM,
		Name:    dataInput.Name,
		Age:     dataInput.Age,
		Address: dataInput.Address,
	}

	// Proses Create DB
	db.Create(&mhs)

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Input Success",
		"data":    mhs,
	})
}

// Update Data Student
func UpdateDataStudent(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Cek Data
	var mhs models.Student
	if err := db.Where("nim = ?", c.Param("nim")).First(&mhs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data Not Found"})
		return
	}

	// Validasi Input
	var dataInput StudentInput

	// Cek yang di input data JSON atau bukan
	if err := c.ShouldBindJSON(&dataInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Proses Update Data
	db.Model(&mhs).Update(dataInput)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Update Success",
		"data":    mhs,
	})
}

// Delete Data Student
func DeleteDataStudent(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Cek Data
	var mhs models.Student
	if err := db.Where("nim = ?", c.Param("nim")).First(&mhs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data Not Found"})
		return
	}

	// Proses Delete Data
	db.Delete(&mhs)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Delete Success",
		"data":    true,
	})
}

// Login by Name
func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	name := c.Params.ByName("name")
	var mhs models.Student
	if err := db.Where("name = ?", name).First(&mhs).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	db.Find(&mhs)
	c.JSON(http.StatusOK, gin.H{
		"Login Success, Welcome": mhs.Name,
	})
}

// Logout
// func Logout(w http.ResponseWriter, r *http.Request){
// 	// reset token yang dikirim, token yang lama ditimpa
// 	resetUserToken(w)

// 	var response UserResponse
// 	response.Status = 200
// 	response.Message = "Success"

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
