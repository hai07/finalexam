package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Cust struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func createCustomer(c *gin.Context) {

	url := os.Getenv("DATABASE_URL")
	fmt.Println("url", url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()
	var cust2 Cust
	err = c.ShouldBindJSON(&cust2)
	fmt.Println("name:", cust2.Name, "email:", cust2.Email)
	row := db.QueryRow("INSERT INTO customer (name, email,status) values ($1, $2, $3) RETURNING id", cust2.Name, cust2.Email, cust2.Status)
	//var id int
	err = row.Scan(&cust2.ID)
	if err != nil {
		fmt.Println("can't scan id", err)
		return
	}

	c.JSON(http.StatusCreated, cust2)
}

func getOneCustomer(c *gin.Context) {
	url := os.Getenv("DATABASE_URL")
	fmt.Println("url", url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()
	var cust2 Cust
	id := c.Param("id")

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customer where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	//rowId := cust2.ID
	row := stmt.QueryRow(id)
	err = row.Scan(&cust2.ID, &cust2.Name, &cust2.Email, &cust2.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, cust2)

}

func getAllCustomer(c *gin.Context) {

	url := os.Getenv("DATABASE_URL")
	fmt.Println("url", url)
	db, err := sql.Open("postgres", url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	custs := []Cust{}
	var cust2 Cust

	smt, err := db.Prepare("SELECT id, name, email, status FROM customer")
	rows, err := smt.Query()
	for rows.Next() {
		err = rows.Scan(&cust2.ID, &cust2.Name, &cust2.Email, &cust2.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		custs = append(custs, cust2)
	}

	c.JSON(http.StatusOK, custs)

}

func updateOneCustomer(c *gin.Context) {

	url := os.Getenv("DATABASE_URL")
	fmt.Println("url", url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	var cust2 Cust
	err = c.ShouldBindJSON(&cust2)
	id := c.Param("id")

	stmt, err := db.Prepare("UPDATE customer SET name = $2 , email = $3, status = $4 where id=$1")

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := stmt.Exec(id, cust2.Name, cust2.Email, cust2.Status); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	getOneCustomer(c)

}

func deleteOneCustomer(c *gin.Context) {

	url := os.Getenv("DATABASE_URL")
	fmt.Println("url", url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	id := c.Param("id")

	stmt, err := db.Prepare("DELETE FROM customer where id=$1")

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := stmt.Exec(id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})

}

func authMiddleware(c *gin.Context) {
	fmt.Println("This is a middlewear")
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized."})
		c.Abort()
		return
	}

	c.Next()

}

func main() {
	url := os.Getenv("DATABASE_URL")
	fmt.Println("url", url)
	db, err := sql.Open("postgres", url)

	createTb := `
	CREATE TABLE IF NOT EXISTS customer (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	r := gin.Default()
	r.Use(authMiddleware)
	r.POST("/customers", createCustomer)
	r.GET("/customers/:id", getOneCustomer)
	r.GET("/customers", getAllCustomer)
	r.PUT("/customers/:id", updateOneCustomer)
	r.DELETE("/customers/:id", deleteOneCustomer)
	r.Run(":2019")
}
