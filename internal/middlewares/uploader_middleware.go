package middlewares

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UploaderMiddleware struct{}

func NewUploaderMiddleware() *UploaderMiddleware {
	return &UploaderMiddleware{}
}

func (u *UploaderMiddleware) Upload(maxSizeMB int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "No file uploaded",
			})
		}

		maxSize := int64(maxSizeMB * 1024 * 1024)
		if file.Size > maxSize {
			return c.Status(400).JSON(fiber.Map{
				"error": "File too large (max " + string(rune(maxSizeMB)) + "MB)",
			})
		}

		c.Locals("uploadedFile", file)
		return c.Next()
	}
}

func (u *UploaderMiddleware) FileFilter(allowedTypes []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "No file uploaded",
			})
		}

		src, err := file.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to open file",
			})
		}
		defer src.Close()

		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to read file",
			})
		}

		mimeType := http.DetectContentType(buffer)

		isAllowed := false
		for _, allowedType := range allowedTypes {
			if mimeType == allowedType {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(400).JSON(fiber.Map{
				"error": "File type " + mimeType + " is not allowed",
			})
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := u.getExtensionsFromMimeTypes(allowedTypes)
		isValidExt := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				isValidExt = true
				break
			}
		}

		if !isValidExt {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid file extension",
			})
		}

		c.Locals("uploadedFile", file)
		c.Locals("mimeType", mimeType)
		return c.Next()
	}
}

func (u *UploaderMiddleware) getExtensionsFromMimeTypes(mimeTypes []string) []string {
	extMap := map[string][]string{
		"image/jpeg": {".jpg", ".jpeg"},
		"image/png":  {".png"},
		"image/gif":  {".gif"},
		"image/webp": {".webp"},
		"image/bmp":  {".bmp"},
		"image/tiff": {".tiff", ".tif"},
	}

	var extensions []string
	for _, mimeType := range mimeTypes {
		if exts, exists := extMap[mimeType]; exists {
			extensions = append(extensions, exts...)
		}
	}
	return extensions
}

func (u *UploaderMiddleware) ImageUpload(maxSizeMB int, allowedTypes []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "No file uploaded",
			})
		}

		maxSize := int64(maxSizeMB * 1024 * 1024)
		if file.Size > maxSize {
			return c.Status(400).JSON(fiber.Map{
				"error": "File too large (max " + string(rune(maxSizeMB)) + "MB)",
			})
		}

		src, err := file.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to open file",
			})
		}
		defer src.Close()

		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to read file",
			})
		}

		mimeType := http.DetectContentType(buffer)

		isAllowed := false
		for _, allowedType := range allowedTypes {
			if mimeType == allowedType {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(400).JSON(fiber.Map{
				"error": "File type " + mimeType + " is not allowed",
			})
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := u.getExtensionsFromMimeTypes(allowedTypes)
		isValidExt := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				isValidExt = true
				break
			}
		}

		if !isValidExt {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid file extension",
			})
		}

		c.Locals("uploadedFile", file)
		c.Locals("mimeType", mimeType)
		return c.Next()
	}
}
