package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type filmTS struct {
	ID       int64  `json:"id"`
	Judul    string `json:"judul"`
	Kategori string `json:"kategori"`
}

var MyConn *sql.DB

func getAllFilm(c *gin.Context) {
	var films []filmTS
	rows, _ := MyConn.Query("SELECT id, judul, kategori FROM tb_films")
	for rows.Next() {
		var film filmTS
		if err := rows.Scan(&film.ID, &film.Judul, &film.Kategori); err == nil {
			films = append(films, film)
		}
	}
	c.JSON(http.StatusOK, films)
}

func getFilm(c *gin.Context) {
	id := c.Param("id")
	row := MyConn.QueryRow("SELECT id, judul, kategori FROM tb_films WHERE id=?", id)
	var film filmTS
	err := row.Scan(&film.ID, &film.Judul, &film.Kategori)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, film)
}

func putFilm(c *gin.Context) {
	id := c.Param("id")
	var newfilm filmTS
	if err := c.ShouldBindJSON(&newfilm); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	_, err := MyConn.Exec("UPDATE tb_films SET judul=?, kategori=? WHERE id=?", newfilm.Judul, newfilm.Kategori, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, newfilm)
}

func delFilm(c *gin.Context) {
	id := c.Param("id")
	_, err := MyConn.Exec("DELETE FROM tb_films WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hapus Data Berhasil"})
}

func addFilm(c *gin.Context) {
	var newfilm filmTS
	if err := c.ShouldBindJSON(&newfilm); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	_, err := MyConn.Exec("INSERT INTO tb_films (judul, kategori) VALUES (?, ?)", newfilm.Judul, newfilm.Kategori)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, newfilm)
}

func main() {
	var err error
	MyConn, err = sql.Open("mysql", "user_bioskop:IniRahasia@tcp(192.168.87.197:13306)/new_bioskop")
	if err != nil {
		panic(err)
	}
	MyConn.SetConnMaxIdleTime(time.Minute)
	MyConn.SetMaxOpenConns(50)
	MyConn.SetConnMaxLifetime(time.Hour)
	defer MyConn.Close()
	router := gin.Default()
	router.GET("/films", getAllFilm)
	router.GET("/film/:id", getFilm)
	router.PUT("/film/:id", putFilm)
	router.DELETE("/film/:id", delFilm)
	router.POST("/film", addFilm)
	router.Run(":8085")
}
