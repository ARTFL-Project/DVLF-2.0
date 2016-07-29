package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	"github.com/jackc/pgx"
)

// Results to export
type Results struct {
	Headword     string              `json:"headword"`
	Dictionaries map[string][]string `json:"dictionaries"`
	Synonyms     []string            `json:"synonyms"`
	Antonyms     []string            `json:"antonyms"`
	// Examples     map[string][]Examples
	// UserSubmit   []string
}

// Examples for headwords
type Examples struct {
	Headword string
	Content  string
	Score    int64
	ID       string
	DefSim   float32
}

var defaultConnConfig pgx.ConnConfig
var pool = createConnPool()

var headwordList = getAllHeadwords()

func createConnPool() *pgx.ConnPool {
	defaultConnConfig.Host = "localhost"
	defaultConnConfig.Database = "dvlf"
	defaultConnConfig.User = "artfl"
	defaultConnConfig.Password = "martini"
	defaultConnConfig.RuntimeParams = make(map[string]string)
	defaultConnConfig.RuntimeParams["statement_timeout"] = "40000"
	config := pgx.ConnPoolConfig{ConnConfig: defaultConnConfig, MaxConnections: 50}
	poolConn, err := pgx.NewConnPool(config)
	if err != nil {
		fmt.Printf("Unable to create connection pool: %v", err)
	}
	return poolConn
}

func getAllHeadwords() []string {
	query := "SELECT headword FROM headwords"
	rows, err := pool.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var headwords []string
	for rows.Next() {
		var headword string
		err := rows.Scan(&headword)
		if err != nil {
			fmt.Println(err)
		}
		headwords = append(headwords, headword)
	}

	cl := collate.New(language.French, collate.Loose, collate.IgnoreCase)

	cl.SortStrings(headwords)
	return headwords
}

func query(c echo.Context) error {
	headword, _ := url.QueryUnescape(c.Param("headword"))
	query := "SELECT dictionaries, synonyms, antonyms FROM headwords WHERE headword=$1"
	var results Results
	var dictionaries map[string][]string
	var synonyms []string
	var antonyms []string
	err := pool.QueryRow(query, headword).Scan(&dictionaries, &synonyms, &antonyms)
	if err != nil {
		fmt.Println(err)
		var empty []string
		return c.JSON(http.StatusOK, empty)
	}
	results = Results{headword, dictionaries, synonyms, antonyms}
	return c.JSON(http.StatusOK, results)
}

func index(c echo.Context) error {
	indexByte, _ := ioutil.ReadFile("public/index.html")
	index := string(indexByte)
	return c.HTML(http.StatusOK, index)
}

func main() {
	// Echo instance
	e := echo.New()

	fmt.Println(len(headwordList))
	e.SetDebug(true)

	e.Static("/static", "public/static")
	e.Static("/app", "public/app")

	// e.File("/favicon.ico", "images/favicon.ico")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// Route => handler
	e.GET("/", index)
	e.GET("/mot/*", index)

	// API
	e.GET("/api/mot/:headword", query)
	e.GET("/api/wordwheel", func(c echo.Context) error {
		return c.JSON(http.StatusOK, headwordList)
	})

	// Start server
	e.Run(standard.New(":8080"))
}
