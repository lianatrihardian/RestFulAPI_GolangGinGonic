package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/animal_shelter")
	err = db.Ping()
	if err != nil {
		panic("Gagal Menghubungkan ke Database")
	}
	defer db.Close()

	router := gin.Default()

	type Hewan struct {
		Id_hewan    int    `json: "id_hewan"`
		Nama_hewan  string `json: "nama_hewan"`
		Jenis_hewan string `json: "jenis_hewan"`
		Jekel       string `json: "jekel"`
		Warna       string `json: "warna"`
	}

	// Menampilkan Detail Data Berdasarkan ID
	router.GET("/:id_hewan", func(c *gin.Context) {
		var (
			hewan  Hewan
			result gin.H
		)
		id_hewan := c.Param("id_hewan")
		row := db.QueryRow("select id_hewan, nama_hewan, jenis_hewan, jekel, warna from hewan where id_hewan = ?;", id_hewan)
		err = row.Scan(&hewan.Id_hewan, &hewan.Nama_hewan, &hewan.Jenis_hewan, &hewan.Jekel, &hewan.Warna)
		if err != nil {
			// If no results send null
			result = gin.H{
				"Hasilnya": "Tidak ada data hewan yang ditemukan",
			}
		} else {
			result = gin.H{
				"Hasilnya":  hewan,
				"Jumlahnya": 1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET all persons
	router.GET("/", func(c *gin.Context) {
		var (
			hewan  Hewan
			hewans []Hewan
		)
		rows, err := db.Query("select id_hewan, nama_hewan, jenis_hewan, jekel, warna from hewan;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&hewan.Id_hewan, &hewan.Nama_hewan, &hewan.Jenis_hewan, &hewan.Jekel, &hewan.Warna)
			hewans = append(hewans, hewan)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"Hasilnya":  hewans,
			"Jumlahnya": len(hewans),
		})
	})

	// POST new person details
	router.POST("/", func(c *gin.Context) {
		var buffer bytes.Buffer
		id_hewan := c.PostForm("id_hewan")
		nama_hewan := c.PostForm("nama_hewan")
		jenis_hewan := c.PostForm("jenis_hewan")
		jekel := c.PostForm("jekel")
		warna := c.PostForm("warna")
		stmt, err := db.Prepare("insert into hewan (id_hewan, nama_hewan, jenis_hewan, jekel, warna) values(?,?,?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id_hewan, nama_hewan, jenis_hewan, jekel, warna)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(nama_hewan)
		buffer.WriteString(" ")
		buffer.WriteString(jenis_hewan)
		buffer.WriteString(" ")
		buffer.WriteString(jekel)
		buffer.WriteString(" ")
		buffer.WriteString(warna)
		defer stmt.Close()
		datanya := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"Pesane": fmt.Sprintf(" Berhasil menambahkan Hewan %s ", datanya),
		})
	})

	// PUT - update a person details
	router.PUT("/", func(c *gin.Context) {
		var buffer bytes.Buffer
		id_hewan := c.PostForm("id_hewan")
		nama_hewan := c.PostForm("nama_hewan")
		jenis_hewan := c.PostForm("jenis_hewan")
		jekel := c.PostForm("jekel")
		warna := c.PostForm("warna")
		stmt, err := db.Prepare("update hewan set nama_hewan= ?, jenis_hewan = ?, jekel= ?, warna= ? where id_hewan= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(nama_hewan, jenis_hewan, jekel, warna, id_hewan)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(nama_hewan)
		buffer.WriteString(" ")
		buffer.WriteString(jenis_hewan)
		buffer.WriteString(" ")
		buffer.WriteString(jekel)
		buffer.WriteString(" ")
		buffer.WriteString(warna)
		defer stmt.Close()
		datanya := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"Pesannya": fmt.Sprintf("Berhasil Merubah Id %s Menjadi %s", id_hewan, datanya),
		})
	})

	// Delete resources
	router.DELETE("/", func(c *gin.Context) {
		id_hewan := c.PostForm("id_hewan")
		stmt, err := db.Prepare("delete from hewan where id_hewan= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id_hewan)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Pesannya": fmt.Sprintf("Berhasil Menghapus %s", id_hewan),
		})
	})
	router.Run(":8080")
}
