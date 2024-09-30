package main

import (
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/webp"
)

func convertWebPToPNG(c *gin.Context) {
	// Получаем файл из запроса
	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, "Не удалось получить файл")
		return
	}

	// Открываем файл из формы
	srcFile, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось открыть файл")
		return
	}
	defer srcFile.Close()

	// Декодируем WebP изображение
	img, err := webp.Decode(srcFile)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось декодировать WebP изображение")
		return
	}

	// Создаем временный файл для PNG изображения
	tempFile, err := os.CreateTemp("./", "image-*.png")
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось создать временный файл")
		return
	}
	defer os.Remove(tempFile.Name()) // Удаляем файл после завершения
	defer tempFile.Close()

	// Кодируем изображение в PNG формат и записываем в временный файл
	err = png.Encode(tempFile, img)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось закодировать PNG изображение")
		return
	}

	// Перемещаем указатель на начало файла для его отправки клиенту
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось прочитать временный файл")
		return
	}

	// Отправляем изображение обратно клиенту
	c.Header("Content-Type", "image/png")
	c.File(tempFile.Name())
}

func main() {
	r := gin.Default()
	r.POST("/convert", convertWebPToPNG)
	r.Run(":8080")
}
