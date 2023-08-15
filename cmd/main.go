package main

import (
	"context"
	_ "github.com/lib/pq"
	"google-gen/internal/handlers"
	"google-gen/internal/postgres"
	"google-gen/internal/postgres/CRUD"
	"google-gen/internal/postgres/repo"
	"google-gen/internal/service"
	"google-gen/pkg/config"
	"log"
	"os"
	"os/signal"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewPostgresConnection(postgres.ConnectionInfo{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username_DB,
		DBName:   cfg.DBName,
		SSLMode:  cfg.SSLMode,
		Password: cfg.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	user := CRUD.NewUsers(db)
	userService := repo.UserRepoI(user)

	question := CRUD.NewQuestion(db)
	questRepo := repo.QuestionRepo(question)
	questService := service.NewQuestion(questRepo)

	label := CRUD.NewLabel(db)
	labelRepo := repo.LabelRepo(label)
	labelService := service.NewLabel(labelRepo)

	choice := CRUD.NewChoice(db)
	choiceRepo := repo.ChoicesRepo(choice)
	choiceService := service.NewChoice(choiceRepo)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	handlers.New(ctx, cfg.TelegramToken, userService, questService, labelService, choiceService)
	defer cancel()

}
