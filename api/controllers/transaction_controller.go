package controllers

import (
	"dreampay/api/models"
	"fmt"
	"io/ioutil"
	"net/http"

	jsontime "github.com/liamylian/jsontime/v2/v2"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
)

var json = jsontime.ConfigWithCustomTimeFormat

// GetTransactionByIDController ... Transaction Per Account Buyer Or Seller
// @Summary Transaction Per Account Buyer Or Seller
// @Description API URL For Transaction Account Buyer Or Seller
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param receiver query string false "mobile"
// @Param depositor query string false "mobile"
// @Success 202 {array} models.Transaction
// @Router /api/transaction [get]
func (s *Server) GetTransactionByIDController(c *gin.Context) {
	var err error
	receiver := c.Query("receiver")
	depositor := c.Query("depositor")
	trans := []models.Transaction{}
	if receiver == "" && depositor == "" {
		err = s.DB.Preload(clause.Associations).Find(&trans).Error
	} else if receiver != "" && depositor == "" {
		err = s.DB.Where("transaction_receiver = ?", receiver).Preload(clause.Associations).Find(&trans).Error
	} else if receiver == "" && depositor != "" {
		err = s.DB.Where("transaction_depositor = ?", depositor).Preload(clause.Associations).Find(&trans).Error
	} else {
		err = s.DB.Where("transaction_receiver = ? AND transaction_depositor = ?", receiver, depositor).Preload(clause.Associations).Find(&trans).Error
	}

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
	withdraw := []models.Withdraw{}
	var total_withdraw int64 = 0
	err_withdraw := s.DB.Where("seller_id = ?", receiver).Find(&withdraw)
	if err_withdraw != nil {
		total_withdraw += 0
	}

	for _, item_with := range withdraw {
		total_withdraw += item_with.Amount
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"response": map[string]interface{}{
			"record":         trans,
			"debit":          debit,
			"credit":         credit,
			"sum":            summary,
			"withdraw":       withdraw,
			"total_withdraw": total_withdraw,
		},
	})
}

// GetAllTransactionController ... All Transaction
// @Summary All Transaction
// @Description API URL For All Transaction
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Success 202 {object} models.Transaction
// @Router /api/transactions [get]
func (s *Server) GetAllTransactionController(c *gin.Context) {
	tran := []models.Transaction{}
	err := s.DB.Preload(clause.Associations).Find(&tran).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":   http.StatusNotFound,
			"response": "Transaction not found !",
		})
		return
	}
	var total_debit int64 = 0
	var total_credit int64 = 0
	for _, item := range tran {
		if item.IsDebit {
			total_debit += item.TransactionAmount
		} else {
			total_credit += item.TransactionAmount
		}
	}
	uang_sisa := total_debit - total_credit
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"response": map[string]interface{}{
			"record":      &tran,
			"uang_buyer":  total_debit,
			"uang_seller": total_credit,
			"uang_sisa":   uang_sisa,
		},
	})
}

// CreateTransactionController ... Create Transaction
// @Summary Create Transaction
// @Description API URL For Create Transaction
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param Account body models.Transaction true "Transaction Data"
// @Success 201 {object} models.Transaction
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

// GetMoneyStatusController ... Money Status
// @Summary Money Status
// @Description API URL For Money Status
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Success 200 {object} models.MoneyStatus
// @Router /api/money-status [get]
func (s *Server) GetMoneyStatusController(c *gin.Context) {
	trans := []models.Transaction{}
	err := s.DB.Preload(clause.Associations).Find(&trans).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":   http.StatusNotFound,
			"response": "Transaction not found !",
		})
		return
	}
	var total_debit int64 = 0
	var total_credit int64 = 0
	for _, item := range trans {
		if item.IsDebit {
			total_debit += item.TransactionAmount
		} else {
			total_credit += item.TransactionAmount
		}
	}
	uang_sisa := total_debit - total_credit
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"response": map[string]interface{}{
			"uang_buyer":  total_debit,
			"uang_seller": total_credit,
			"uang_sisa":   uang_sisa,
		},
	})
}

// CreateTransactionController ... Create Withdraw
// @Summary Create Withdraw
// @Description API URL For Create Withdraw
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param Account body models.Withdraw true "Transaction Data"
// @Success 201 {object} models.Withdraw
// @Router /api/withdraw [post]
func (s *Server) CreateWithdrawController(c *gin.Context) {
	errList := map[string]string{}
	w := models.Withdraw{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	err = json.Unmarshal(body, &w)
	if err != nil {
		errList["Unmarshal_error"] = err.Error()
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	/*
		trans := []models.Transaction{}
		var total int64
		err_total := s.DB.Where("transaction_receiver = ?", w.SellerID).Preload(clause.Associations).Find(&trans).Error
		if err_total != nil {
			total = 0
		}
		for _, item := range trans {
			total += item.TransactionAmount
		}
		if total > w.Amount {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":   http.StatusInternalServerError,
				"response": "Total withdraw melebihi kapasitas saldo",
			})
			return
		}
	*/
	res, err := w.CreateWithdraw(s.DB)
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

// GetWithdrawController ... Withdraw History
// @Summary Withdraw History
// @Description API URL For Withdraw History
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param seller query string false "mobile"
// @Success 200 {array} models.Withdraw
// @Router /api/withdraw [get]
func (s *Server) GetWithdrawController(c *gin.Context) {
	var err error
	seller := c.Query("seller")
	with := []models.Withdraw{}
	if seller == "" {
		err = s.DB.Preload(clause.Associations).Find(&with).Error
	} else {
		err = s.DB.Where("seller_id = ?", seller).Preload(clause.Associations).Find(&with).Error
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":   http.StatusNotFound,
			"response": "Withdraw not found !",
		})
		return
	}
	var total_withdraw int64 = 0
	for _, item := range with {
		total_withdraw += item.Amount
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"response": map[string]interface{}{
			"record":         &with,
			"total_withdraw": total_withdraw,
		},
	})
}

// DeleteMultipleTransaction ... Delete Multiple Transaction
// @Summary Delete Multiple Transaction
// @Description API URL For Delete Multiple Transaction
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param Account body models.TransactionID true "Transaction Data"
// @Success 202
// @Router /api/transaction/delete/multiple [post]
func (s *Server) DeleteMultipleTransaction(c *gin.Context) {
	errList := map[string]string{}
	ids := models.TransactionID{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	err = json.Unmarshal(body, &ids)
	if err != nil {
		errList["Unmarshal_error"] = err.Error()
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	t := models.Transaction{}
	for _, item := range ids.ID {
		err = s.DB.Where("id = ?", item).Delete(&t).Error
		if err != nil {
			errList["Internal_Error"] = "Process stoped cause failed to delete account " + item
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":   http.StatusAccepted,
		"response": "Delete Complete",
	})
}

// DeleteMultipleWithdraw ... Delete Multiple Withdraw
// @Summary Delete Multiple Withdraw
// @Description API URL For Delete Multiple Withdraw
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param Account body models.TransactionID true "Transaction Data"
// @Success 202
// @Router /api/withdraw/delete/multiple [post]
func (s *Server) DeleteMultipleWithdraw(c *gin.Context) {
	errList := map[string]string{}
	ids := models.TransactionID{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	err = json.Unmarshal(body, &ids)
	if err != nil {
		errList["Unmarshal_error"] = err.Error()
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	t := models.Withdraw{}
	for _, item := range ids.ID {
		err = s.DB.Where("id = ?", item).Delete(&t).Error
		if err != nil {
			errList["Internal_Error"] = "Process stoped cause failed to delete account " + item
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":   http.StatusAccepted,
		"response": "Delete Complete",
	})
}
