package main

import (
  "context"
  "os"

  "github.com/joho/godotenv"
  "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
  if err := godotenv.Load(); err != nil {
    panic(err)
  }
  
  ctx := context.Background()

  pool, err := pgxpool.New(ctx, fmt.Springf(
    "user=%s password=%s host=%s port=%s dbname=%s",
    os.Getenv("WBSCK_DATABASE_USER"),
    os.Getenv("WBSCK_DATABASE_PASSWORD"),
    os.Getenv("WBSCK_DATABASE_HOST"),
    os.Getenv("WBSCK_DATABASE_PORT"),
    os.Getenv("WBSCK_DATABASE_NAME"),
  ))
  
  if err != nil {
    panic(err)
  }
  
  defer pool.Close()
  
  if err := pool.Ping(ctx); err != nil {
    panic(err)
  }
  
  handler := api.NewHandler(pgstore.New(pool))
  
  go func() {
    if err := http.L

  }()

}
