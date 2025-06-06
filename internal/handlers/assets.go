package handlers

import (
	"backend-service/internal/utils"
	"bytes"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"image"
	_ "image/gif" // Импортируйте нужные форматы изображений
	_ "image/jpeg"
	_ "image/png"
	"io"
)

func (h *Handler) UploadFile(c *fiber.Ctx) error {
	// получаем файл из формы
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "missing file")
	}

	// открываем файл
	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to open file")
	}
	defer file.Close()

	// читаем в память
	data, err := io.ReadAll(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to read file")
	}

	// генерим уникальное имя
	token := uuid.New().String()

	// загружаем в S3
	err = h.S3.Upload(context.Background(), token, data, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "upload failed")
	}

	// Инициализируем переменные для ширины и высоты
	var width, height int

	// Проверяем, является ли файл изображением
	if utils.IsImage(fileHeader.Header.Get("Content-Type")) {
		// Декодируем изображение
		img, _, err := image.Decode(bytes.NewReader(data))
		if err == nil {
			bounds := img.Bounds()
			width = bounds.Dx()
			height = bounds.Dy()
		}
	}

	// возвращаем имя файла
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"details": fiber.Map{
			"token":    token,
			"url":      "/api/v1/assets/" + token,
			"size":     fileHeader.Size,
			"name":     fileHeader.Filename,
			"mimeType": fileHeader.Header.Get("Content-Type"),
			"width":    width,
			"height":   height,
		},
	})
}

func (h *Handler) GetFile(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing token")
	}

	object, err := h.S3.Minio.GetObject(context.Background(), h.S3.Bucket, token, minio.GetObjectOptions{})
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "object not found")
	}

	// Получаем метаданные (например, размер, контент-тайп)
	info, err := object.Stat()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "object stat failed")
	}

	// Можно указать inline или attachment
	c.Set("Content-Type", info.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", info.Size))
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", token))

	return c.SendStream(object)
}

func (h *Handler) DeleteFile(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing token")
	}

	err := h.S3.Delete(context.Background(), token)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
