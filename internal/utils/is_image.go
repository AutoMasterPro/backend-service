package utils

// IsImage Вспомогательная функция для проверки, является ли MIME-тип изображением
func IsImage(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/png", "image/gif", "image/bmp", "image/webp", "image/tiff":
		return true
	default:
		return false
	}
}
