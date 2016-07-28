package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	"github.com/jackc/pgx"
)

// Results to export
type Results struct {
	Headword     string              `json:"headword"`
	Dictionaries map[string][]string `json:"dictionaries"`
	// Examples     map[string][]Examples
	// Synonyms     []string
	// Antonyms     []string
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

type caseInsensitive []string

func (s caseInsensitive) Len() int {
	return len(s)
}
func (s caseInsensitive) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s caseInsensitive) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
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
	sort.Sort(caseInsensitive(headwords))
	return headwords
}

func query(c echo.Context) error {
	headword, _ := url.QueryUnescape(c.Param("headword"))
	query := "SELECT dictionaries FROM headwords WHERE headword=$1"
	var results Results
	var dictionaries map[string][]string
	err := pool.QueryRow(query, headword).Scan(&dictionaries)
	if err != nil {
		fmt.Println(err)
		var empty []string
		return c.JSON(http.StatusOK, empty)
	}
	results = Results{headword, dictionaries}
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
	e.GET("/api/mot/:headword", query)
	e.GET("/api/wordwheel", func(c echo.Context) error {
		return c.JSON(http.StatusOK, headwordList)
	})

	// Start server
	e.Run(standard.New(":8080"))
}
