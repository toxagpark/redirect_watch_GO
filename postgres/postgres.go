package postgres

import (
	"log"
	"math/rand"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB() {
	PostgresLog, err := os.ReadFile("postgres/PostgresLog.txt")
	if err != nil {
		log.Fatalf("Ошибка чтения информации о БД: %w", err)
	}
	PostgresLogStr := string(PostgresLog)
	DB, err = gorm.Open(postgres.Open(PostgresLogStr), nil)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %w", err)
	}

	log.Println("Успешное подключение к БД")
}

func randomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func AddURLPostgres(original_url string) string {
	short_url := randomString(6)
	query := `INSERT INTO short_urls (original_url, short_code) VALUES ($1, $2)`
	res := DB.Exec(query, original_url, short_url)
	if res.Error != nil {
		log.Fatalf("не удалось сохранить ориг урл: %w", res.Error)
	}
	return short_url
}

func GetURLPostgres(short_url string) string {
	var result struct {
		Original_url string `gorm:"original_url"`
	}

	err := DB.Table("short_urls").Where("short_code = ?", short_url).First(&result).Error
	if err != nil {
		log.Fatalf("ошибка Постгрес %v", err)
	}
	return result.Original_url
}

func AddVisitPostgres(ip, userAgent, original_url string) {
	query := "INSERT INTO visits_go (original_link, ip, user_agent) VALUES ($1, $2, $3)"
	res := DB.Exec(query, original_url, ip, userAgent)
	if res.Error != nil {
		log.Fatalf("Не записываются гости: %w", res.Error)
	}
}
