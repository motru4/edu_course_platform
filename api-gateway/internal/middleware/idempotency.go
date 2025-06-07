package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IdempotencyMiddleware обеспечивает идемпотентность запросов
type IdempotencyMiddleware struct {
	// Храним информацию о запросах в памяти
	// В реальном приложении лучше использовать Redis или другое хранилище
	processedRequests sync.Map
	// Время жизни записей в кэше
	cacheTTL time.Duration
}

// RequestKey структура для хранения информации о запросе
type RequestKey struct {
	Method string
	Path   string
	Body   string
}

// CachedResponse структура для хранения ответа
type CachedResponse struct {
	Status       int
	Body         []byte
	Headers      map[string][]string
	LastAccessed time.Time
}

// NewIdempotencyMiddleware создает новый экземпляр IdempotencyMiddleware
func NewIdempotencyMiddleware(ttl time.Duration) *IdempotencyMiddleware {
	middleware := &IdempotencyMiddleware{
		cacheTTL: ttl,
	}

	// Запускаем горутину для очистки кэша
	go middleware.cleanCache()

	fmt.Println("Идемпотентность запросов инициализирована с TTL:", ttl)
	return middleware
}

// Middleware возвращает Gin middleware для обеспечения идемпотентности запросов
func (m *IdempotencyMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Игнорируем GET запросы, так как они идемпотентны по определению
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		// Получаем идентификатор идемпотентности из заголовка, если он есть
		idempotencyKey := c.GetHeader("X-Idempotency-Key")
		if idempotencyKey == "" {
			// Если заголовок не предоставлен, создаем ключ на основе запроса
			// Читаем тело запроса
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				fmt.Printf("Ошибка при чтении тела запроса: %v\n", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "невозможно прочитать тело запроса"})
				c.Abort()
				return
			}

			// Восстанавливаем тело запроса для дальнейшего использования
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Создаем уникальный ключ для запроса
			idempotencyKey = m.createRequestKey(c.Request.Method, c.Request.URL.Path, string(bodyBytes))
		}

		fmt.Printf("Обработка запроса с ключом идемпотентности: %s\n", idempotencyKey)

		// Проверяем, был ли такой запрос уже обработан
		if cachedResp, found := m.processedRequests.Load(idempotencyKey); found {
			resp := cachedResp.(*CachedResponse)
			// Обновляем время последнего доступа
			resp.LastAccessed = time.Now()
			m.processedRequests.Store(idempotencyKey, resp)

			fmt.Printf("Найден кэшированный ответ для ключа: %s (статус: %d)\n", idempotencyKey, resp.Status)

			// Устанавливаем заголовки
			for name, values := range resp.Headers {
				for _, value := range values {
					c.Header(name, value)
				}
			}

			// Добавляем заголовок идемпотентности
			c.Header("X-Idempotency-Key", idempotencyKey)
			c.Header("X-Idempotency-From-Cache", "true")

			// Возвращаем кэшированный ответ
			c.Data(resp.Status, "application/json", resp.Body)
			c.Abort()
			return
		}

		// Создаем ResponseWriter для перехвата ответа
		w := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		// Добавляем заголовок идемпотентности в оригинальный запрос
		c.Header("X-Idempotency-Key", idempotencyKey)

		// Выполняем запрос
		c.Next()

		// Кэшируем ответ
		cachedResp := &CachedResponse{
			Status:       w.status,
			Body:         w.body.Bytes(),
			Headers:      make(map[string][]string),
			LastAccessed: time.Now(),
		}

		// Копируем заголовки
		for k, v := range w.Header() {
			cachedResp.Headers[k] = v
		}

		// Сохраняем в кэше
		m.processedRequests.Store(idempotencyKey, cachedResp)
		fmt.Printf("Кэширован новый ответ для ключа: %s (статус: %d)\n", idempotencyKey, w.status)
	}
}

// createRequestKey создает уникальный ключ для запроса
func (m *IdempotencyMiddleware) createRequestKey(method, path, body string) string {
	h := sha256.New()
	h.Write([]byte(method + path + body))
	return hex.EncodeToString(h.Sum(nil))
}

// cleanCache периодически удаляет устаревшие записи из кэша
func (m *IdempotencyMiddleware) cleanCache() {
	ticker := time.NewTicker(m.cacheTTL / 2)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		count := 0
		m.processedRequests.Range(func(key, value interface{}) bool {
			resp := value.(*CachedResponse)
			if now.Sub(resp.LastAccessed) > m.cacheTTL {
				m.processedRequests.Delete(key)
				count++
			}
			return true
		})
		fmt.Printf("Очистка кэша: удалено %d устаревших записей\n", count)
	}
}

// responseWriter перехватывает ответ от обработчика
type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

// Write записывает данные в буфер и в исходный ResponseWriter
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteHeader записывает статус ответа
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
