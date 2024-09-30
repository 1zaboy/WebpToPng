package main

import (
	"image/png"
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

	// Открываем файл
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
	outFile, err := os.Create("output.png")
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось создать PNG файл")
		return
	}
	defer outFile.Close()

	// Кодируем изображение в PNG формат
	err = png.Encode(outFile, img)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось закодировать PNG изображение")
		return
	}

	// Отправляем результат
	c.File("output.png")
}

func main() {
	r := gin.Default()
	r.POST("/convert", convertWebPToPNG)
	r.Run(":8080")
}
