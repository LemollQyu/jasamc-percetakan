package handler

import (
	"errors"
	"fmt"
	"jasamc/infrastructure/log"
	"jasamc/models"
	"jasamc/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ini masing bingung dibuat route api search by id tidak (masih nanti unutk handler)
func (h *JasaHandler) GetCategoryJasaByID(c *gin.Context) {
	// ambil param ID dan konvers ke int64
	idParam := c.Param("id")
	jasaID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"id":    idParam,
			"error": err.Error(),
		}).Error("Gagal ambil ID, strconv.ParseInt")

		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid jasa id",
		})
		return
	}

	// panggil usecase
	category, err := h.JasaUsecase.GetCategoryJasaByID(c.Request.Context(), jasaID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"id":    jasaID,
			"error": err.Error(),
		}).Error("Handler gagal ambil jasa by ID, h.JasaUsecase.GetJasaByID")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "Kesalahan dari system",
		})
		return
	}

	// cek category kosong atau nggak
	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_message": "Category tidak ditemukan",
		})

		return
	}

	// response sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    category,
	})
}

// get all acategory
func (h *JasaHandler) GetAllCategoryJasa(c *gin.Context) {
	categories, err := h.JasaUsecase.GetAllCategoryJasa(c.Request.Context())
	if err != nil {
		log.Logger.Error("JasaHandler: h.JasaUsecase.GetAllCategoryJasa", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "Kesalahan dari system",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    categories,
	})

}

// create category jasa
func (h *JasaHandler) CreateCategoryJasa(c *gin.Context) {
	var param models.ParamCreateCategoryJasa

	// cek inputan valid atau nggak
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Input",
		})

		return
	}

	// cek param inputan ID tidak boleh != 0
	if param.ID != 0 {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Error("invalid request - jasa category id is not empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Request",
		})

		return
	}

	// create data category jasa
	categoryJasaID, err := h.JasaUsecase.CreateCategoryJasa(c.Request.Context(), &param)
	if err != nil {
		if errors.Is(err, utils.ErrNameExists) || errors.Is(err, utils.ErrSlugExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err.Error(),
			})
			return
		}

		// error lain adalah internal
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "Terjadi kesalahan dari system",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Successfully create new category product %d", categoryJasaID),
	})

}

// handle delete
func (h *JasaHandler) DeleteCategoryJasa(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid category id",
		})
		return
	}

	err = h.JasaUsecase.DeleteCategoryJasa(c.Request.Context(), id)
	if err != nil {
		// kategori tidak ditemukan
		if err.Error() == "category jasa tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})
			return
		}

		// kategori masih dipakai oleh services
		if err.Error() == "category tidak bisa dihapus karena masih digunakan oleh service lain" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err.Error(),
			})
			return
		}

		// error lain internal
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "terjadi kesalahan dari server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "category jasa berhasil dihapus",
	})
}

// handle update icon di meta
func (h *JasaHandler) UpdateCategoryIcon(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	file, err := c.FormFile("icon")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "icon file required"})
		return
	}

	if err := h.JasaUsecase.UpdateCategoryIcon(c, id, file); err != nil {
		if err.Error() == "category jasa tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})
			return
		}

		if err.Error() == "icon hanya png dan svg" || err.Error() == "ukuran icon maksimal 2MB" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Icon berhasil terupload di category id %d", id),
	})
}

// handle status category jasa
func (h *JasaHandler) SetStatusJasaCategory(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	category, err := h.JasaUsecase.UpdateStatusCategoryJasa(c.Request.Context(), id)
	if err != nil {

		log.Logger.WithFields(logrus.Fields{
			"id":    id,
			"error": err.Error(),
		}).Error("Handler gagal, h.JasaUsecase.UpdateStatusCategoryJasa")

		if err.Error() == "category jasa tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "Kesalahan dari system",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully set status category %d", id),
		"data":    category,
	})
}
