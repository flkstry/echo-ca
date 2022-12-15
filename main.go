package main

import (
	"database/sql"
	"fmt"
	"log"
	"login-system/server/app/account/delivery/httphandler"
	"login-system/server/app/account/repository"
	"login-system/server/app/account/usecase"
	"login-system/server/internals"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	p := &internals.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDB := os.Getenv("DB_DB")

	conn := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbPort, dbUsername, dbPassword, dbDB)
	dbCon, err := sql.Open(`postgres`, conn)
	if err != nil {
		log.Fatalln(err)
	}

	err = dbCon.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err := dbCon.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	e := echo.New()
	accountRepo := repository.NewPsqlAccountRepository(dbCon)
	timeoutCtx := time.Duration(5) * time.Second

	accountUsecase := usecase.NewAccountUsecase(accountRepo, timeoutCtx, *p)
	httphandler.NewAccountHandler(e, accountUsecase)

	log.Fatal(e.Start("127.0.0.1:3000"))
}
