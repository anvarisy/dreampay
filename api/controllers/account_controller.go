package controllers

import (
	"dreampay/api/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateAccountController ... Create Account
// @Summary Create New Account
// @Description API URL For Create New Account
// @Tags Create Account
// @Accept  json
// @Produce  json
// @Param Account body models.AccountRegister true "Account Data"
// @Success 201 {object} models.AccountRegister
// @Router /api/register [post]
func (s *Server) CreateAccountController(c *gin.Context) {
	account := models.Account{}
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	res, err := account.CreateAccount(s.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	// RunWa(res.Mobile, res.Token, res.Name)
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": res,
	})
}

// VerificationAccountController ... Verification Account
// @Summary Verification Account
// @Description API URL For Verification Account
// @Tags Verification Account
// @Accept  json
// @Produce  json
// @Param Account body models.AccountVerification true "Account Data"
// @Success 202 {object} models.AccountVerification
// @Router /api/verification [post]
func (s *Server) VerificationAccountController(c *gin.Context) {
	account := models.Account{}
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err := s.DB.Where("account_mobile = ? and is_active = ?", account.AccountMobile, false).Take(&account).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":   http.StatusNotFound,
			"response": "Mobile already confirmed!",
		})
		return
	} else {
		s.DB.Where("account_mobile = ?", account.AccountMobile).Updates(models.Account{IsActive: true})
		c.JSON(http.StatusAccepted, gin.H{
			"status":   http.StatusAccepted,
			"response": "Verification Completed !",
		})
		return
	}

}

// LoginAccountController ... Login Account
// @Summary Login Account
// @Description API URL For Login Account
// @Tags Login Account
// @Accept  json
// @Produce  json
// @Param Account body models.AccountLogin true "Account Data"
// @Success 202 {object} models.AccountLogin
// @Router /api/login [post]
func (s *Server) LoginAccountController(c *gin.Context) {
	account := models.Account{}
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err := s.DB.Where("account_mobile = ? AND is_active = ?", account.AccountMobile, true).Take(&account).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":   http.StatusNotFound,
			"response": "Mobile not found !",
		})
		return
	}
	fmt.Println(&account)
	c.JSON(http.StatusAccepted, gin.H{
		"status":   http.StatusAccepted,
		"response": "Login Completed!",
		"data":     account,
	})

}

// DeleteAccountController ... Delete Account
// @Summary Delete Account
// @Description API URL For Delete Account
// @Tags Delete Account
// @Accept  json
// @Produce  json
// @Param mobile path string true "mobile"
// @Success 202 {object} models.Account
// @Router /api/delete/{mobile} [post]
func (s *Server) DeleteAccountController(c *gin.Context) {
	account := models.Account{}
	mobile := c.Param("mobile")
	err := s.DB.Where("account_mobile = ?", mobile).Delete(&account).Error
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"status": http.StatusNoContent,
			"error":  err,
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":   http.StatusAccepted,
		"response": "Delete Completed !",
	})
}

// GetAllAccountController ... List Account
// @Summary -
// @Description -
// @Tags Get Account
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Account
// @Router /api/account [get]
func (s *Server) GetAllAccountController(c *gin.Context) {
	accounts := models.Account{}
	res, err := accounts.GetAllAccount(s.DB)
	if err != nil {
		c.JSON(http.StatusNoContent, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": res,
	})
}

// LoginAccountController ... Update Account
// @Summary Update Account
// @Description API URL For Update Account
// @Tags Update Account
// @Accept  json
// @Produce  json
// @Param Account body models.Account true "Account Data"
// @Success 202 {object} models.Account
// @Router /api/update/{account_mobile} [post]
func (s *Server) UpdateAccountController(c *gin.Context) {
	id := c.Param("id")
	account := models.Account{}
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err := s.DB.Where("account_mobile = ?", id).Updates(&account).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   http.StatusInternalServerError,
			"response": err,
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":   http.StatusAccepted,
		"response": "Update Completed!",
		"data":     account,
	})

}
