package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"strconv"
	"strings"
	"time"

	"cafe/bootstrap"
	"cafe/database/seed"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	dialect     = "postgres"
	migrationsDir = "database/migrations"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go [up|down|status|create|redo|reset]")
	}

	command := os.Args[1]

	env := bootstrap.NewEnv()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Ho_Chi_Minh",
		env.DBHost, env.DBUser, env.DBPass, env.DBName, env.DBPort)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect(dialect); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if command == "create" {
		if len(os.Args) < 3 {
			log.Fatal("Migration name required: go run cmd/migrate/main.go create [name]")
		}
		name := os.Args[2]

		// Custom naming: YYYYMMDDXXXX
		now := time.Now()
		datePart := now.Format("20060102")

		// Find highest XXXX for today
		files, _ := os.ReadDir(migrationsDir)
		maxSeq := 0
		for _, f := range files {
			if strings.HasPrefix(f.Name(), datePart) {
				if len(f.Name()) >= 12 {
					seqStr := f.Name()[8:12]
					seq, _ := strconv.Atoi(seqStr)
					if seq > maxSeq {
						maxSeq = seq
					}
				}
			}
		}

		version := fmt.Sprintf("%s%04d", datePart, maxSeq+1)
		filename := fmt.Sprintf("%s/%s_%s.sql", migrationsDir, version, name)

		content := "-- +goose Up\n-- +goose StatementBegin\nSELECT 'up';\n-- +goose StatementEnd\n\n-- +goose Down\n-- +goose StatementBegin\nSELECT 'down';\n-- +goose StatementEnd\n"

		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			log.Fatalf("failed to create migration file: %v", err)
		}

		fmt.Printf("Created migration: %s\n", filename)
		return
	}

	// Check version before Up
	initialVersion, _ := goose.GetDBVersion(db)

	if err := goose.Run(command, db, migrationsDir, os.Args[2:]...); err != nil {
		log.Fatalf("goose run failed: %v", err)
	}

	// If it was up command and version was 0, run seed
	if command == "up" && initialVersion == 0 {
		gormDB := bootstrap.NewPostgresDatabase(env)
		defer bootstrap.ClosePostgresDBConnection(gormDB)
		seed.Seed(gormDB, env)
	}
}
