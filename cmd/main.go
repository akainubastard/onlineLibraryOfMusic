// Package main содержит основную функцию для запуска веб-сервера.
package main

import (
    "log"           // Импортируем пакет для логирования
    // "work/database"  // Импортируем пакет database для работы с базой данных
    "work/details"   // Импортируем пакет details для конфигурации
    handlers "work/handlefunc" // Импортируем пакет handlefunc для работы с обработчиками
    "github.com/labstack/echo/v4" // Импортируем Echo - веб-фреймворк для Go
    echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Music API
// @version 1.0
// @description API для управления музыкальной библиотекой
// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
    // Инициализация базы данных
    // database.InitDb()

    // Создаем новый экземпляр Echo
    e := echo.New()

    // Определяем маршруты для работы с песнями
    e.GET("/swagger/*", echoSwagger.WrapHandler)
    // @Summary Получить список песен
    // @Description Получает список песен с возможностью фильтрации по группе и названию песни, а также с пагинацией.
    // @Produce json
    // @Param group query string false "Фильтр по группе"
    // @Param song query string false "Фильтр по названию песни"
    // @Param page query int false "Номер страницы" default(1)
    // @Param size query int false "Количество песен на странице" default(10)
    // @Success 200 {array} details.Song
    // @Failure 400 {object} map[string]string {"error": "Ошибка при выполнении запроса"}
    // @Failure 404 {object} map[string]string {"error": "Song not found"}
    // @Router /songs [get]
    e.GET("/songs", handlers.GetSongs) // Обработчик для получения списка песен
    // @Summary Получить текст песни по ID
    // @Description Получает текст песни по указанному ID
    // @Produce json
    // @Param id path string true "ID песни"
    // @Success 200 {object} struct{ Text string "text" } "Текст песни"
    // @Failure 400 {object} map[string]string {"error": "Ошибка при выполнении запроса"}
    // @Failure 404 {object} map[string]string {"error": "Song not found"}
    // @Failure 500 {object} map[string]string {"error": "Ошибка выполнения запроса к базе данных"}
    // @Router /songs/{id}/text [get]
    e.GET("/songs/:id", handlers.GetSongText) // Обработчик для получения текста песни по ID
    // @Summary Добавить или обновить информацию о песне
    // @Description Добавляет новую песню или обновляет существующую песню в базе данных. Сначала проверяет информацию о песне через внешний API.
    // @Accept json
    // @Produce json
    // @Param song body details.Song true "Детали песни"
    // @Success 200 {string} string "The song was successfully added"
    // @Failure 400 {object} map[string]string {"error": "Invalid Input"}
    // @Failure 500 {object} map[string]string {"error": "Internal server error"}
    // @Router /songs [put]
    e.POST("/songs", handlers.PutSongs) // Обработчик для добавления новой песни
    // @Summary Обновить информацию о песне по ID
    // @Description Обновляет данные песни в базе данных по указанному ID
    // @Accept json
    // @Produce json
    // @Param id path string true "ID песни"
    // @Param song body details.Song true "Детали песни"
    // @Success 200 {string} string "The song was successfully updated"
    // @Failure 400 {string} string "Invalid Input"
    // @Failure 404 {string} string "The song not found"
    // @Failure 500 {string} string "Could not update the song"
    // @Router /songs/{id} [put]
    e.PATCH("/songs/:id", handlers.UpdateSong) // Обработчик для обновления существующей песни по ID
    // @Summary Удалить песню по ID
    // @Description Удаляет песню из базы данных по указанному ID
    // @Produce json
    // @Param id path string true "ID песни"
    // @Success 200 {string} string "The song was successfully deleted"
    // @Failure 404 {string} string "The song not found"
    // @Failure 500 {string} string "Could not delete the song"
    // @Router /songs/{id} [delete]
    e.DELETE("/songs/:id", handlers.DeleteSong) // Обработчик для удаления песни по ID

    // Запускаем сервер на порту 8080
    if err := e.Start(details.Config.LOCAL_PORT); err != nil {
        log.Fatalf("Не удалось создать сервер: %v", err) // Логируем ошибку, если сервер не удалось запустить
    }
    
    log.Printf("Сервер запущен на порту %s\n", details.Config.LOCAL_PORT) // Логируем информацию о запуске сервера
}