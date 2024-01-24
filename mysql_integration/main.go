package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	Id     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "recordings",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!!!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Albums found %v\n", albums)

	album, err := albumById(2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Album found: %v\n", album)

	albId, err := addAlbum(Album{
		Title:  "Nuevo disco",
		Artist: "Jordi",
		Price:  640.99,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Id of added album %v\n", albId)
}

func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("Select * from album where artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.Id, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func albumById(id int64) (Album, error) {
	var alb Album

	row := db.QueryRow("Select * from album where id = ?", id)
	if err := row.Scan(&alb.Id, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}

	return alb, nil
}

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("Insert into album (title, artist, price) VALUES (?,?,?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum %v", err)
	}

	return id, nil
}
