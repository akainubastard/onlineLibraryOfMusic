// Package database предоставляет функции для инициализации и работы с базой данных.
package database

import (
	"fmt"
	"log"
	"work/details" // Импортируем пакет types

	"gorm.io/driver/postgres" // Импортируем драйвер PostgreSQL для GORM
	"gorm.io/gorm"            // Импортируем GORM для работы с базой данных
)

// Db - глобальная переменная для хранения соединения с базой данных.
var Db *gorm.DB

// InitDb инициализирует соединение с базой данных и выполняет миграции для структур из пакета types.
func InitDb() {
    // Строка подключения к базе данных PostgreSQL

    var err error
    c := details.Config
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",c.POSTGRES_HOST,c.POSTGRES_USER,c.POSTGRES_PASSWORD,c.POSTGRES_DBNAME,c.LOCAL_PORT)
    // Открываем соединение с базой данных
    Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        // Если не удалось подключиться, выводим сообщение об ошибке и завершаем программу
        log.Fatalf("Неудалось подключиться к базе данных: %v", err)
    }

    // Выполняем миграцию для структуры SongDetail
    if err := Db.AutoMigrate(&details.SongDetail{}); err != nil {
        // Если не удалось создать таблицу SongDetail, выводим сообщение об ошибке
        log.Fatalf("Неудалось создать таблицу SongDetail: %v", err)
    }

    // Выполняем миграцию для структуры Song
    if err := Db.AutoMigrate(&details.Song{}); err != nil {
        // Если не удалось создать таблицу Song, выводим сообщение об ошибке
        log.Fatalf("Неудалось создать таблицу Song: %v", err)
    }

    // Если все прошло успешно, выводим сообщение об успешной инициализации базы данных
    log.Println("База данных успешно инициализирована")
}