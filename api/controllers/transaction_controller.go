package controllers

import (
	"dreampay/api/models"
	"fmt"
	"io/ioutil"
	"net/http"

	jsontime "github.com/liamylian/jsontime/v2/v2"

	"github.com/gin-gonic/gin"
)

var json = jsontime.ConfigWithCustomTimeFormat

// GetTransactionByIDController ... Transaction Per Account
// @Summary Transaction Per Account
// @Description API URL For Transaction Account
// @Tags Transaction Per Account
// @Accept  json
// @Produce  json
// @Param mobile path string true "mobile"
// @Success 202 {object} models.Transaction
// @Router /api/transaction/{id} [get]
func (s *Server) GetTransactionByIDController(c *gin.Context) {
	receiver := c.Param("id")
	trans := []models.Transaction{}
	err := s.DB.Where("transaction_depositor = ?", receiver).Find(&trans).Error
	var debit int64 = 0
	var credit int64 = 0

	for _, item := range trans {

		if item.IsDebit {
			debit += item.TransactionAmount
		} else {
			credit += item.TransactionAmount
		}
	}
	summary := debit - credit
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":   http.StatusNotFound,
			"response": "Transaction not found !",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": trans,
		"debit":    debit,
		"credit":   credit,
		"sum":      summary,
	})
}

// GetAllTransactionController ... All Transaction
// @Summary All Transaction
// @Description API URL For All Transaction
// @Tags All Transaction
// @Accept  json
// @Produce  json
// @Success 202 {object} models.Transaction
// @Router /api/transaction [get]
func (s *Server) GetAllTransactionController(c *gin.Context) {
	tran := []models.Transaction{}
	err := s.DB.Find(&tran).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":   http.StatusNotFound,
			"response": "Transaction not found !",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": &tran,
	})
}

// CreateTransactionController ... Create Transaction
// @Summary Create Transaction
// @Description API URL For Create Transaction
// @Tags Create Transaction
// @Accept  json
// @Produce  json
// @Param Account body models.Transaction true "Transaction Data"
// @Success 202 {object} models.Transaction
// @Router /api/transaction [post]
func (s *Server) CreateTransactionController(c *gin.Context) {
	errList := map[string]string{}
	i := models.Transaction{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	err = json.Unmarshal(body, &i)
	if err != nil {
		errList["Unmarshal_error"] = err.Error()
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	fmt.Println("Transaction section")
	fmt.Println(&i)
	res, err := i.CreateTransaction(s.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   http.StatusInternalServerError,
			"response": res,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": res,
	})
}
