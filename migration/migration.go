package migration

// import (
// 	"embed"

// 	"database/sql"

// 	"github.com/pressly/goose/v3"
// )

// //go:embed schema/*.sql
// var embedMigrations embed.FS

// func New(db *sql.DB) *migration {
// 	return &migration{ db }
// }

// type migration struct {
// 	db *sql.DB
// }

// func (m *migration) Up() error {
// 	goose.SetBaseFS(embedMigrations)

//     if err := goose.SetDialect("postgres"); err != nil {
//         return err
//     }

//     if err := goose.Up(m.db, "schema"); err != nil {
//         return err
//     }

// 	return nil
// }

// func (m *migration) Down() error {
// 	goose.SetBaseFS(embedMigrations)

//     if err := goose.SetDialect("postgres"); err != nil {
//         return err
//     }

//     if err := goose.Down(m.db, "schema"); err != nil {
//         return err
//     }

// 	return nil
// }
