package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	// Capture connection properties
	cfg := mysql.Config{
		User:                 "fkrfdllh",
		Passwd:               "18januari",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "recordings",
		AllowNativePasswords: true,
	}

	// Get a database handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Albums found: %v\n", albums)

	album, err := albumByID(2)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Album found: %v\n", album)

	albumID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID of added album: %v\n", albumID)
}

// albumsByArtist queries for albums that have the specified artist name
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)

	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields
	for rows.Next() {
		var album Album

		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
		}

		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

// albumByID queries for the album with the specified ID
func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row
	var album Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumByID %d: no such album", id)
		}

		return album, fmt.Errorf("albumByID %d: %v", id, err)
	}

	return album, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(album Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", album.Title, album.Artist, album.Price)

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}
