package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// connStr := "host=localhost port=5432 user=postgres dbname=datahara password=postgres sslmode=disable"
	connStr := "host=postgres-map port=5432 user=postgres dbname=datahara password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dir := "./maps"

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".geojson") {
			filename := d.Name()

			// 1. Check if filename already exists in DB
			var exists bool
			err := db.QueryRowContext(context.Background(),
				"SELECT EXISTS(SELECT 1 FROM geodata.kondisi WHERE name = $1)", filename).
				Scan(&exists)
			if err != nil {
				return fmt.Errorf("DB check failed for %s: %v", filename, err)
			}

			if exists {
				fmt.Printf("Skipping %s (already in database).\n", filename)
				return nil
			}

			// 2. Read and validate JSON
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read %s: %v", filename, err)
			}

			var tmp interface{}
			if err := json.Unmarshal(content, &tmp); err != nil {
				return fmt.Errorf("invalid JSON in %s: %v", filename, err)
			}

			// 3. Insert into DB
			_, err = db.ExecContext(context.Background(),
				"INSERT INTO geodata.kondisi (name, map) VALUES ($1, $2)",
				filename, string(content))
			if err != nil {
				return fmt.Errorf("failed to insert %s: %v", filename, err)
			}

			fmt.Printf("Inserted %s into the database.\n", filename)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}
}
