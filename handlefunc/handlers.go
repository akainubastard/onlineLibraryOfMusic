package handlefunc

import (
	"encoding/json" // Импортируем пакет для работы с JSON
	"log"           // Импортируем пакет для логирования
	"net/http"      // Импортируем пакет для работы с HTTP
	"strconv"       // Импортируем пакет для преобразования строк в числа
	"work/database" // Импортируем пакет database для работы с базой данных
	"work/details"  // Импортируем пакет details
	"github.com/go-playground/validator" // Импортируем пакет для валидации структур
	"github.com/labstack/echo/v4"        // Импортируем Echo - веб-фреймворк для Go
	"gorm.io/gorm"                       // Импортируем GORM для работы с базой данных
)

var validate *validator.Validate

// init - функция инициализации, которая вызывается автоматически при загрузке пакета.
func init() {
	// Создаем новый экземпляр валидатора.
	validate = validator.New()
}

// Получение списка песен с фильтрацией и пагинацией
func GetSongs(c echo.Context) error {
	// Получаем параметры фильтрации из запроса
	groupFilter := c.QueryParam("group")
	songFilter := c.QueryParam("song")
	pageStr := c.QueryParam("page")
	sizeStr := c.QueryParam("size")

	// Преобразуем параметры страницы и размера в целые числа с обработкой ошибок
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Устанавливаем значение по умолчанию для страницы
		log.Printf("Некорректный параметр 'page': %s. Устанавливаем значение по умолчанию: %d\n", pageStr, page)
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10 // Устанавливаем значение по умолчанию для размера
		log.Printf("Некорректный параметр 'size': %s. Устанавливаем значение по умолчанию: %d\n", sizeStr, size)
	}

	// Логируем параметры фильтрации и пагинации
	log.Printf("Фильтрация - Группа: %s, Песня: %s, Страница: %d, Размер: %d\n", groupFilter, songFilter, page, size)

	// Создаем срез для хранения песен
	songs := make([]details.Song, 0, size)

	// Вычисляем смещение для пагинации
	offset := (page - 1) * size

	// Начинаем запрос к базе данных с учетом пагинации
	db := database.Db.Offset(offset).Limit(size)

	// Применяем фильтры, если они указаны
	if groupFilter != "" {
		db = db.Where("\"group\" = ?", groupFilter)
		log.Printf("Применен фильтр для группы: %s\n", groupFilter)
	}
	if songFilter != "" {
		db = db.Where("song = ?", songFilter)
		log.Printf("Применен фильтр для песни: %s\n", songFilter)
	}

	// Выполняем запрос и обрабатываем ошибки
	result := db.Find(&songs)
	if  result.Error != nil {
		log.Printf("Ошибка при выполнении запроса: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	// Логируем количество найденных песен
	log.Printf("Найдено песен: %d\n", len(songs))

	// Возвращаем успешный ответ с найденными песнями
	return c.JSON(http.StatusOK, &songs)
}

// Получение текста песни
func GetSongText(c echo.Context) error {
    // Извлекаем 'id' песни из параметров URL
    stringID := c.Param("id")
    log.Printf("Получен запрос на получение текста песни с ID: %s", stringID)

    // Определяем структуру для хранения текста песни
    var textSong struct {
        Text string `json:"text"`
    }

    // Получаем доступ к базе данных
    db := database.Db
    
    // Выполняем запрос для получения текста песни по ID
    result := db.Select("text").Table("song_details").Where("id = ?", stringID).Scan(&textSong)

    // Проверяем, произошла ли ошибка во время выполнения запроса
    if result.Error != nil {
        log.Printf("Ошибка базы данных при получении текста песни: %v", result.Error)
        // Возвращаем ошибку 500 с описанием
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
    }

    // Проверяем, найден ли запросом какой-либо результат
    if result.RowsAffected == 0 {
        log.Printf("Песня с ID: %s не найдена", stringID)
        // Возвращаем ошибку 404, если песня не найдена
        return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
    }

    // Успешно найдено - возвращаем текст песни
    log.Printf("Успешно получен текст песни с ID: %s", stringID)
    return c.JSON(http.StatusOK, textSong)
}
// Добавление песни
func PutSongs(c echo.Context) error {
	// Создаем переменную для хранения информации о песне
	var song details.Song

	// Логируем информацию о начале обработки запроса
	log.Println("Начинаем обработку запроса на добавление/обновление песни")

	// Привязываем входящие данные к структуре песни
	if err := c.Bind(&song); err != nil {
		log.Printf("Ошибка при привязке входящих данных: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Input"})
	}

	// Логируем полученные данные о песне
	log.Printf("Получены данные о песне: %+v\n", song)

	// Запрос к внешнему API для проверки информации о песне
	apiURL := "http://api-url/info" // ЗАМЕНИТЬ НА РЕАЛЬНЫЙ API
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Printf("Ошибка создания запроса к API: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	// Устанавливаем параметры запроса
	q := req.URL.Query()
	q.Add("group", song.Group)
	q.Add("song", song.Song)
	req.URL.RawQuery = q.Encode()

	// Отправка запроса к API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка выполнения запроса к API: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	defer resp.Body.Close()

	// Обрабатываем ответ от API
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка API: %s\n", resp.Status)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}

	// Декодируем ответ от API
	var songDetail details.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		log.Printf("Ошибка декодирования ответа от API: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := database.Db.Create(&songDetail).Error; err != nil {
		log.Printf("Ошибка при создании записи о деталях песне в базе данных: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Логируем информацию о состоянии перед добавлением в базу данных
	log.Println("Проверяем данные перед добавлением в базу данных")

	// Проверяем, прошла ли песня валидацию
	if err := validate.Struct(song); err != nil {
		log.Printf("Сообщение не прошло проверку валидации: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Пытаемся создать новую запись песни в базе данных
	if err := database.Db.Create(&song).Error; err != nil {
		log.Printf("Ошибка при создании записи о песне в базе данных: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Логируем успешное добавление песни на уровне info
	log.Printf("Песня успешно добавлена: %+v\n", song)

	// Если создание прошло успешно, возвращаем статус 200 и сообщение об успешном добавлении
	return c.JSON(http.StatusOK, "The song was successfully added")
}

// Удаление песни
func DeleteSong(c echo.Context) error {
	// Получаем параметр "id" из URL
	stringID := c.Param("id")

	// Логируем информацию о начале процесса удаления
	log.Printf("Попытка удалить песню с ID: %s\n", stringID)

	// Пытаемся удалить запись песни из базы данных по указанному ID
	if err := database.Db.Delete(&details.Song{}, stringID).Error; err != nil {
		// Если произошла ошибка, она будет записана в лог
		log.Printf("Ошибка при попытке удалить песню с ID %s: %v\n", stringID, err)

		// Если запись не найдена, возвращаем статус 404 и сообщение о том, что песня не найдена
		if err == gorm.ErrRecordNotFound {
			log.Printf("Песня с ID %s не найдена\n", stringID)
			return c.JSON(http.StatusNotFound, "The song not found")
		}

		// Если произошла другая ошибка, возвращаем статус 500 и сообщение об ошибке
		log.Println("Не удалось удалить песню:", err)
		return c.JSON(http.StatusInternalServerError, "Could not delete the song")
	}

	// Логируем успешное удаление
	log.Printf("Песня с ID %s успешно удалена\n", stringID)

	// Если удаление прошло успешно, возвращаем статус 200 и сообщение об успешном удалении
	return c.JSON(http.StatusOK, "The song was successfully deleted")
}

// Обновление песни
func UpdateSong(c echo.Context) error {
	// Получаем параметр "id" из URL
	stringID := c.Param("id")
	var song details.Song

	// Логируем информацию о начале процесса обновления
	log.Printf("Начинаем обновление песни с ID: %s\n", stringID)

	// Проверка на ошибки при связывании данных
	if err := c.Bind(&song); err != nil {
		log.Printf("Ошибка при привязке входящих данных: %v\n", err)
		return c.JSON(http.StatusBadRequest, "Invalid Input")
	}

	// Логируем полученные данные о песне на уровне debug
	log.Printf("Получены данные для обновления песни: %+v\n", song)

	// Обновление записи в базе данных
	if err := database.Db.Model(&song).Where("id = ?", stringID).Updates(&song).Error; err != nil {
		log.Printf("Ошибка при обновлении песни с ID %s: %v\n", stringID, err)
		return c.JSON(http.StatusInternalServerError, "Could not update the song")
	}

	// Логируем успешное обновление
	log.Printf("Песня с ID %s успешно обновлена\n", stringID)

	// Возврат успешного ответа
	return c.JSON(http.StatusOK, "The song was successfully updated")
}
