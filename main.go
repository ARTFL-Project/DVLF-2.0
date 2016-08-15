package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	"github.com/jackc/pgx"
	"github.com/kennygrant/sanitize"
)

// Results to export
type Results struct {
	Headword     string              `json:"headword"`
	Dictionaries map[string][]string `json:"dictionaries"`
	Synonyms     []string            `json:"synonyms"`
	Antonyms     []string            `json:"antonyms"`
	UserSubmit   []UserSubmit        `json:"userSubmit"`
	Examples     []Example           `json:"examples"`
}

// Example for headwords
type Example struct {
	Content string  `json:"content"`
	Link    string  `json:"link"`
	Score   int     `json:"score"`
	ID      int     `json:"id"`
	DefSim  float32 `json:"defSim"`
	Source  string  `json:"source"`
}

// UserSubmit fields
type UserSubmit struct {
	Content string `json:"content"`
	Source  string `json:"source"`
	Link    string `json:"link"`
}

// AutoCompleteList is the top 10 words
type AutoCompleteList []AutoCompleteHeadword

// AutoCompleteHeadword is just the object in the AutoCompleteList
type AutoCompleteHeadword struct {
	Headword string `json:"headword"`
}

var defaultConnConfig pgx.ConnConfig
var pool = createConnPool()

var headwordList, headwordMap = getAllHeadwords()

var tokenRegex = regexp.MustCompile(`(?i)([\p{L}]+)|([\.?,;:'!\-]+)|([\s]+)`)

var outputLog = createLog()

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

func createLog() *os.File {
	f, _ := os.Create("output.log")
	defer f.Close()
	return f
}

func getAllHeadwords() ([]string, map[string]bool) {
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
	var headwordHash = make(map[string]bool)
	for _, value := range headwords {
		headwordHash[value] = true
	}
	return headwords, headwordHash
}

func autoComplete(c echo.Context) error {
	prefix, _ := url.QueryUnescape(c.Param("prefix"))
	prefix = strings.TrimSpace(prefix)
	prefix = strings.ToLower(prefix)
	prefix += "%"
	query := "SELECT headword FROM headwords WHERE headword LIKE $1 ORDER BY headword LIMIT 10"
	rows, err := pool.Query(query, prefix)
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var headwords AutoCompleteList
	for rows.Next() {
		var headword string
		err := rows.Scan(&headword)
		if err != nil {
			fmt.Println(err)
		}
		headwords = append(headwords, AutoCompleteHeadword{headword})
	}
	results := make(map[string]AutoCompleteList)
	results["headwords"] = headwords
	return c.JSON(http.StatusOK, results)
}

func highlightExamples(examples []Example, queryTerm string) []Example {
	query := "SELECT headword FROM word2lemma WHERE lemma=$1"
	var forms []string
	rows, err := pool.Query(query, queryTerm)

	if err != nil {
		fmt.Println(err)
		return examples
	}

	defer rows.Close()

	for rows.Next() {
		var form string
		err := rows.Scan(&form)
		if err != nil {
			fmt.Println(err)
			return examples
		}
		forms = append(forms, form)
	}

	formRegex := regexp.MustCompile("^(" + strings.Join(forms, "|") + ")$")
	var newExamples []Example
	for example := range examples {
		var newContent []string
		matches := tokenRegex.FindAllString(examples[example].Content, -1)
		for match := range matches {
			if formRegex.MatchString(matches[match]) {
				newContent = append(newContent, fmt.Sprintf(`<span class="highlight">%s</span>`, matches[match]))
			} else {
				newContent = append(newContent, matches[match])
			}
		}
		examples[example].Content = strings.Join(newContent, "")
		newExamples = append(newExamples, examples[example])
	}
	return newExamples
}

func query(c echo.Context) error {
	headword, _ := url.QueryUnescape(c.Param("headword"))
	query := "SELECT user_submit, dictionaries, synonyms, antonyms, examples FROM headwords WHERE headword=$1"
	var results Results
	var dictionaries map[string][]string
	var synonyms []string
	var antonyms []string
	var userSubmission []UserSubmit
	var examples []Example
	err := pool.QueryRow(query, headword).Scan(&userSubmission, &dictionaries, &synonyms, &antonyms, &examples)
	if err != nil {
		fmt.Println(err)
		var empty []string
		return c.JSON(http.StatusOK, empty)
	}
	highlightedExamples := highlightExamples(examples, headword)
	results = Results{headword, dictionaries, synonyms, antonyms, userSubmission, highlightedExamples}
	return c.JSON(http.StatusOK, results)
}

func submitDefinition(c echo.Context) error {
	headword := sanitize.HTML(c.FormValue("term"))
	source := sanitize.HTML(c.FormValue("source"))
	link := sanitize.HTML(c.FormValue("link"))
	allowedTags := []string{"i", "b"}
	definition, htmlErr := sanitize.HTMLAllowing(c.FormValue("definition"), allowedTags)
	if htmlErr != nil {
		fmt.Println(htmlErr)
		message := map[string]string{"message": "error"}
		return c.JSON(http.StatusOK, message)
	}
	var newSubmission = UserSubmit{definition, source, link}
	if _, ok := headwordMap[headword]; ok {
		query := "SELECT user_submit FROM headwords WHERE headword=$1"
		var userSubmission []UserSubmit
		err := pool.QueryRow(query, headword).Scan(&userSubmission)
		if err != nil {
			fmt.Println(err)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}

		userSubmission = append(userSubmission, newSubmission)
		serializedSubmission, jsonError := json.Marshal(userSubmission)
		if jsonError != nil {
			fmt.Println(jsonError)
		}

		update := "UPDATE headwords SET user_submit=$1 WHERE headword=$2"
		_, updateErr := pool.Exec(update, serializedSubmission, headword)
		if updateErr != nil {
			fmt.Println(updateErr)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}
	} else {
		var userSubmission = []UserSubmit{newSubmission}
		serializedSubmission, jsonError := json.Marshal(userSubmission)
		if jsonError != nil {
			fmt.Println(jsonError)
		}
		dictionaries, _ := json.Marshal(map[string][]string{})
		synonyms, _ := json.Marshal([]string{})
		antonyms, _ := json.Marshal([]string{})
		examples, _ := json.Marshal([]Example{})
		insert := "INSERT INTO headwords (headword, dictionaries, synonyms, antonyms, user_submit, examples) VALUES ($1, $2, $3, $4, $5, $6)"
		_, insertErr := pool.Exec(insert, headword, dictionaries, synonyms, antonyms, serializedSubmission, examples)
		if insertErr != nil {
			fmt.Println(insertErr)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}
		headwordMap[headword] = true
		headwordList = append(headwordList, headword)
		cl := collate.New(language.French, collate.Loose, collate.IgnoreCase)
		cl.SortStrings(headwordList)
	}
	message := map[string]string{"message": "success"}
	return c.JSON(http.StatusOK, message)
}

func voteForExample(c echo.Context) error {
	headword, _ := url.QueryUnescape(c.Param("headword"))
	exampleID, _ := strconv.Atoi(c.Param("id"))
	var examples []Example
	query := "SELECT examples FROM headwords WHERE headword=$1"
	err := pool.QueryRow(query, headword).Scan(&examples)
	if err != nil {
		fmt.Println(err)
		message := map[string]string{"message": "error"}
		return c.JSON(http.StatusOK, message)
	}
	var newExamples []Example
	var newScore int
	for _, example := range examples {
		if example.ID == exampleID {
			if c.Param("vote") == "up" {
				example.Score++
			} else {
				example.Score--
			}
			newScore = example.Score
		}
		newExamples = append(newExamples, example)
	}
	update := fmt.Sprintf("UPDATE headwords SET examples=$1 WHERE headword=$2")
	_, updateErr := pool.Exec(update, newExamples, headword)
	if updateErr != nil {
		fmt.Println(updateErr)
		message := map[string]string{"message": "error"}
		return c.JSON(http.StatusOK, message)
	}
	message := map[string]interface{}{"message": "success", "score": newScore}
	return c.JSON(http.StatusOK, message)
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
	e.Use(middleware.Secure())

	e.GET("/", index)
	e.GET("/mot/*", index)
	e.GET("/apropos", index)
	e.GET("/definition", index)

	// API
	e.GET("/api/mot/:headword", query)
	e.GET("/api/wordwheel", func(c echo.Context) error {
		return c.JSON(http.StatusOK, headwordList)
	})
	e.GET("/api/autocomplete/:prefix", autoComplete)
	e.POST("/api/submit", submitDefinition)
	e.GET("/api/vote/:headword/:id/:vote", voteForExample)

	// Start server
	e.Run(standard.New(":8080"))
}
