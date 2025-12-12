package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"properties-api/utils"

	"github.com/gin-gonic/gin"
)

// UploadController maneja las operaciones de upload de archivos
type UploadController struct{}

// NewUploadController crea una nueva instancia de UploadController
func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadImage maneja el upload de una imagen
// POST /api/upload
// Body: multipart/form-data con campo "image"
// Respuesta: { "url": "/images/filename.jpg" }
func (c *UploadController) UploadImage(ctx *gin.Context) {
	// Obtener el archivo del formulario
	file, err := ctx.FormFile("image")
	if err != nil {
		utils.LogError("POST /api/upload", err, ctx)
		utils.BadRequest(ctx, "No se pudo obtener el archivo. Asegúrate de enviar un campo 'image' con el archivo.", nil)
		return
	}

	// Validar que sea una imagen
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		utils.LogError("POST /api/upload", fmt.Errorf("extensión no permitida: %s", ext), ctx)
		utils.BadRequest(ctx, "Formato de imagen no permitido. Solo se aceptan: JPG, JPEG, PNG, GIF, WEBP", nil)
		return
	}

	// Validar tamaño (máximo 5MB)
	if file.Size > 5*1024*1024 {
		utils.LogError("POST /api/upload", fmt.Errorf("archivo muy grande: %d bytes", file.Size), ctx)
		utils.BadRequest(ctx, "El archivo es demasiado grande. Tamaño máximo: 5MB", nil)
		return
	}

	// Generar nombre único para el archivo
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	
	// Ruta donde se guardará el archivo (dentro del contenedor)
	uploadPath := "/app/uploads/images"
	filePath := filepath.Join(uploadPath, filename)

	// Guardar el archivo
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		utils.LogError("POST /api/upload", err, ctx)
		utils.InternalServerError(ctx, "Error al guardar el archivo", nil)
		return
	}

	// Retornar la URL relativa de la imagen
	imageURL := fmt.Sprintf("/images/%s", filename)
	
	utils.LogInfo("POST /api/upload", fmt.Sprintf("Imagen subida exitosamente: %s", filename), ctx)
	utils.SendSuccess(ctx, http.StatusOK, gin.H{
		"url": imageURL,
		"filename": filename,
	})
}


