package handler

import (
	"fmt"
	"jasamc/infrastructure/log"
	"jasamc/models"
	"jasamc/utils"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ambil semua data jasa
func (h *JasaHandler) GetAllJasa(c *gin.Context) {
	categories, err := h.JasaUsecase.GetAllJasa(c.Request.Context())
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

// buat data jasanya dan juga upload icon, gallery, thumbnail
func (h *JasaHandler) CreateService(c *gin.Context) {
	// 1️⃣ Parse form fields
	var param models.RequestService
	if err := c.ShouldBind(&param); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invlaid param",
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "multipart form tidak valid",
		})
		return
	}

	files := map[string][]*multipart.FileHeader{
		"icon":      form.File["icon"],
		"thumbnail": form.File["thumbnail"],
		"gallery":   form.File["gallery"],
	}

	// 3️⃣ Panggil usecase (PAKAI param)
	service, err := h.JasaUsecase.CreateAndUploadFileService(c, param, files)
	if err != nil {
		switch err.Error() {
		case "ukuran file terlalu besar", "ekstensi file tidak diperbolehkan", "thumbnail wajib diisi", "maksimal thumbnail 4", "gallery wajib diisi", "name sudah dipakai", "slug sudah dipakai":
			c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
		case "kategori tidak ditemukan", "category tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{"error_message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
		}
		return
	}

	// 4️⃣ Response sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "Service berhasil dibuat",
		"data":    service,
	})
}

// set status service
func (h *JasaHandler) SetStatusService(c *gin.Context) {
	// ambil id nya
	idService := c.Param("serviceID")
	id, err := strconv.ParseInt(idService, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid service id",
		})
		return
	}
	service, err := h.JasaUsecase.SetStatusService(c.Request.Context(), id)
	if err != nil {

		log.Logger.WithFields(logrus.Fields{
			"id":    id,
			"error": err.Error(),
		}).Error("Handler gagal, h.JasaUsecase.SetStatusService")

		if err.Error() == "service tidak ditemukan" {
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
		"message": fmt.Sprintf("Successfully set status service %d", id),
		"data":    service.IsActive,
	})
}

// unutk delete service tapi belum dulu karena service_spesification sama yang value belum ada / nunggu itu
// func (h *JasaHandler) DeleteService(c *gin.Context) {

// }

// hapus media di jasa itu
func (h *JasaHandler) DeleteMediaInService(c *gin.Context) {

	// ambil id nya
	id := c.Param("serviceID")
	serviceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid service id",
		})
		return
	}

	id = c.Param("mediaID")
	mediaID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid media id",
		})
		return
	}

	err = h.JasaUsecase.DeleteMediaInServiceFromDBAndStorage(c.Request.Context(), serviceID, mediaID)

	if err != nil {
		switch err.Error() {
		case "invalid file path", "empty media path", "invalid media url":
			c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
		case "service tidak ditemukan", "media tidak ditemukan", "media tidak terkait dengan service ini":
			c.JSON(http.StatusNotFound, gin.H{"error_message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil dihapus",
	})
}

// create service spesification
func (h *JasaHandler) CreateServiceSpesification(c *gin.Context) {
	var param models.RequestServiceSpesification

	if err := c.ShouldBind(&param); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invlaid param",
		})
		return
	}

	// validasi kondisi jika typenya select option wajib diisi
	if param.InputType == "select" {
		if utils.IsEmptyJSON(param.Options) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Options wajib diisi dan tidak boleh kosong",
			})
			return
		}
	} else {
		param.Options = nil
	}

	_, err := h.JasaUsecase.CreateServiceSpesification(c.Request.Context(), &param)
	if err != nil {
		switch err.Error() {
		case "name sudah digunakan":
			c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
		case "service tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{"error_message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Spesifikasi service  berhasil dibuat",
	})

}

func (h *JasaHandler) AddServiceMedia(c *gin.Context) {

	id := c.Param("serviceID")
	serviceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid service id",
		})
		return
	}

	var param models.RequestAddServiceMedia
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid param",
		})
		return
	}

	// ambil file sesuai type
	file, err := c.FormFile(param.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": fmt.Sprintf("file %s wajib diupload", param.Type),
		})
		return
	}

	err = h.JasaUsecase.AddServiceMediaByType(
		c,
		serviceID,
		param,
		file,
	)

	if err != nil {
		switch err.Error() {

		case
			"file tidak sesuai dengan type",
			"icon sudah ada",
			"thumbnail maksimal 3",
			"file icon wajib diupload",
			"file thumbnail wajib diupload",
			"file gallery wajib diupload",
			"ekstensi file tidak diperbolehkan",
			"ukuran file terlalu besar":

			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err.Error(),
			})

		case "service tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "media berhasil ditambahkan",
	})
}

// handler service spesification value
func (h *JasaHandler) CreateServiceSpesificationValue(c *gin.Context) {
	var param models.RequestServiceSpesificationValue

	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid param",
		})
		return
	}

	err := h.JasaUsecase.CreateServiceSpesificationValue(c.Request.Context(), param)
	if err != nil {
		switch err.Error() {

		case "value spesification sudah digunakan / dibuat":
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})

		case "service tidak ditemukan", "service spesification tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})

		default:
			log.Logger.Error("jasaHandler: h.JasaUsecase.CreateServiceSpesificationValue")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "service spesification value berhasil ditambahkan",
	})

}

func (h *JasaHandler) UpdateServiceSpesificationValue(c *gin.Context) {

	var param models.RequestUpdateServiceSpesificationValue

	id := c.Param("serviceID")
	serviceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid service id",
		})
		return
	}

	id = c.Param("specID")
	specID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid spec id",
		})
		return
	}

	id = c.Param("valueID")
	valueID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid value id",
		})
		return
	}

	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid param",
		})
		return
	}

	err = h.JasaUsecase.UpdateServiceSpesificationValue(c.Request.Context(), serviceID, specID, valueID, param)

	if err != nil {
		switch err.Error() {
		case "service spesification value tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})

		default:
			log.Logger.Error("jasaHandler: h.JasaUsecase.CreateServiceSpesificationValue")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "service spesification value berhasil diupdate",
	})

}

func (h *JasaHandler) DeleteServiceSpesification(c *gin.Context) {
	id := c.Param("serviceID")
	serviceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid service id",
		})
		return
	}

	id = c.Param("specID")
	specID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid spec id",
		})
		return
	}

	err = h.JasaUsecase.DeleteServiceSpesification(c.Request.Context(), serviceID, specID)
	if err != nil {
		switch err.Error() {

		case "service tidak ditemukan", "service spesification tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})

		default:
			log.Logger.Error("jasaHandler: h.JasaUsecase.CreateServiceSpesificationValue")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "service spesification berhasil dihapus",
	})
}

func (h *JasaHandler) DeleteService(c *gin.Context) {
	id := c.Param("serviceID")
	serviceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid service id",
		})
		return
	}

	err = h.JasaUsecase.DeleteServiceByID(c.Request.Context(), serviceID)
	if err != nil {
		switch err.Error() {

		case "service tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})
		case "invalid media url", "empty media path", "invalid file path":
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err.Error(),
			})

		default:
			log.Logger.Error("jasaHandler: h.JasaUsecase.DeleteServiceByID")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": "Kesalahan dari system",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service berhasil dihapus",
	})
}

func (h *JasaHandler) ToggleServiceSpesificationStatus(c *gin.Context) {
	id := c.Param("serviceID")
	serviceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid service id",
		})
		return
	}

	id = c.Param("specID")
	specID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid spec id",
		})
		return
	}

	serviceSpec, err := h.JasaUsecase.ToggleServiceSpesificationStatus(c.Request.Context(), serviceID, specID)
	if err != nil {
		switch err.Error() {

		case "service tidak ditemukan", "service spesification tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})

		default:
			log.Logger.Error("jasaHandler: h.JasaUsecase.ToggleServiceSpesificationStatus")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Spesification id %d Status berhasil diubah", specID),
		"data":    serviceSpec.IsActive,
	})

}

func (h *JasaHandler) ToggleServiceSpesificationRequired(c *gin.Context) {
	id := c.Param("serviceID")
	serviceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid service id",
		})
		return
	}

	id = c.Param("specID")
	specID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid spec id",
		})
		return
	}

	serviceSpec, err := h.JasaUsecase.ToggleServiceSpesificationRequired(c.Request.Context(), serviceID, specID)
	if err != nil {
		switch err.Error() {

		case "spesification tidak boleh required jika tidak aktif":
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err.Error(),
			})

		case "service tidak ditemukan", "service spesification tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{
				"error_message": err.Error(),
			})

		default:
			log.Logger.Error("jasaHandler: h.JasaUsecase.ToggleServiceSpesificationRequired")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Spesification id %d required berhasil diubah", specID),
		"data":    serviceSpec.IsRequired,
	})

}

func (h *JasaHandler) GetService(c *gin.Context) {
	id := c.Param("serviceID")
	serviceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "invalid service id",
		})
		return
	}
	service, err := h.JasaUsecase.GetServiceByIDFromRead(c.Request.Context(), serviceID)
	if err != nil {
		log.Logger.Error("JasaHandler: h.JasaUsecase.GetServiceByIDfromRead", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "Kesalahan dari system",
		})
		return
	}

	if service == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Service tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    service,
	})
}
