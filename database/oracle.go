package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/sijms/go-ora/v2"
)

func OpenOracle() *sql.DB {
	// Pega valores das variáveis de ambiente
	user := os.Getenv("ORACLE_USER")
	pass := os.Getenv("ORACLE_PASS")
	host := os.Getenv("ORACLE_HOST")
	port := os.Getenv("ORACLE_PORT")
	service := os.Getenv("ORACLE_SERVICE")

	// Se alguma variável estiver vazia, usa valores padrão
	if user == "" {
		user = "DEV"
	}
	if pass == "" {
		pass = "DEV"
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "1521"
	}
	if service == "" {
		service = "XEPDB1"
	}

	dsn := fmt.Sprintf("oracle://%s:%s@%s:%s/%s", user, pass, host, port, service)

	db, err := sql.Open("oracle", dsn)
	if err != nil {
		log.Fatal("Erro ao abrir conexão:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Erro ao conectar no Oracle:", err)
	}

	log.Println("✅ Conectado ao Oracle com database/sql")
	return db
}
