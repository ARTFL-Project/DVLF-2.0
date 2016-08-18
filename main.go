package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

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
	Headword     string         `json:"headword"`
	Dictionaries DictionaryData `json:"dictionaries"`
	Synonyms     []string       `json:"synonyms"`
	Antonyms     []string       `json:"antonyms"`
	Examples     []Example      `json:"examples"`
}

// Dictionary to export
type Dictionary struct {
	Name       string              `json:"name"`
	Label      string              `json:"label"`
	ShortLabel string              `json:"shortLabel"`
	Content    []map[string]string `json:"contentObj"`
	Show       bool                `json:"show"`
}

// DictionaryData to export
type DictionaryData struct {
	Data         []Dictionary `json:"data"`
	TotalDicos   int          `json:"totalDicos"`
	TotalEntries int          `json:"totalEntries"`
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

// ExamplesByID for sorting
type ExamplesByID []Example

func (a ExamplesByID) Len() int {
	return len(a)
}
func (a ExamplesByID) Less(i, j int) bool {
	return a[i].ID < a[j].ID
}
func (a ExamplesByID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
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

// RecaptchaResponse from Google
type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []int     `json:"error-codes"`
}

var defaultConnConfig pgx.ConnConfig
var pool = createConnPool()

var headwordList, headwordMap = getAllHeadwords()

var tokenRegex = regexp.MustCompile(`(?i)([\p{L}]+)|([\.?,;:'!\-]+)|([\s]+)`)

var outputLog = createLog()

var dicoLabels = map[string]map[string]string{
	"feraud": {
		"label":      "Féraud: Dictionaire critique de la langue française (1787-1788)",
		"shortLabel": "Féraud (1787-1788)",
	},
	"nicot": {
		"label":      "Jean Nicot: Thresor de la langue française (1606)",
		"shortLabel": "Jean Nicot (1606)",
	},
	"acad1694": {
		"label":      "Dictionnaire de L'Académie française 1re édition (1694)",
		"shortLabel": "Académie française (1694)",
	},
	"acad1762": {
		"label":      "Dictionnaire de L'Académie française 4e édition (1762)",
		"shortLabel": "Académie française (1762)",
	},
	"acad1798": {
		"label":      "Dictionnaire de L'Académie française 5e édition (1798)",
		"shortLabel": "Académie française (1798)",
	},
	"acad1835": {
		"label":      "Dictionnaire de L'Académie française 6e édition (1835)",
		"shortLabel": "Académie française (1835)",
	},
	"littre": {
		"label":      "Émile Littré: Dictionnaire de la langue française (1872-1877)",
		"shortLabel": "Littré (1872-1877)",
	},
	"acad1932": {
		"label":      "Dictionnaire de L'Académie française 8e édition (1932-1935)",
		"shortLabel": "Académie française (1932-1935)",
	},
	"tlfi": {
		"label":      "Le Trésor de la Langue Française Informatisé",
		"shortLabel": "Trésor Langue Française",
	},
	"bob": {
		"label":      "BOB: Dictionaire d'argot",
		"shortLabel": "BOB: Dictionaire d'argot",
	},
}

var dicoOrder = []string{"tlfi", "acad1932", "littre", "acad1835", "acad1798", "feraud", "acad1762", "acad1694", "nicot", "bob"}

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

func orderDictionaries(dictionaries map[string][]string, userSubmissions []UserSubmit) DictionaryData {
	displayed := 0
	var show bool
	totalEntries := 0
	totalDicos := 0
	var newDicos []Dictionary
	for _, dico := range dicoOrder {
		if len(dictionaries[dico]) == 0 {
			continue
		}
		totalDicos++
		displayed++
		if displayed < 3 {
			show = true
		} else {
			show = false
		}
		totalEntries += len(dictionaries[dico])
		var content []map[string]string
		for _, entry := range dictionaries[dico] {
			content = append(content, map[string]string{"content": entry})
		}
		newDicos = append(newDicos, Dictionary{dico, dicoLabels[dico]["label"], dicoLabels[dico]["shortLabel"], content, show})
	}
	if len(userSubmissions) > 0 {
		totalEntries += len(userSubmissions)
		displayed++
		if displayed < 3 {
			show = true
		} else {
			show = false
		}
		var content []map[string]string
		for _, entry := range userSubmissions {
			content = append(content, map[string]string{"content": entry.Content, "source": entry.Source, "link": entry.Link})
		}
		newDicos = append(newDicos, Dictionary{"userSubmit", "Définition(s) d'utilisateurs", "Définition(s) d'utilisateurs", content, show})
	} else {
		newDicos = append(newDicos, Dictionary{"userSubmit", "Définition(s) d'utilisateurs", "Définition(s) d'utilisateurs", make([]map[string]string, 0), true})
	}
	allDictionaries := DictionaryData{newDicos, totalDicos, totalEntries}
	return allDictionaries
}

func highlightExamples(examples []Example, queryTerm string) []Example {
	query := "SELECT headword FROM word2lemma WHERE lemma=$1"
	forms := []string{queryTerm}
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

func sortExamples(examples []Example) []Example {
	var orderedExamples []Example
	var negativeExamples []Example
	var otherExamples []Example
	for _, example := range examples {
		if example.Score > 0 {
			orderedExamples = append(orderedExamples, example)
		} else if example.Score < 0 {
			negativeExamples = append(negativeExamples, example)
		} else {
			otherExamples = append(otherExamples, example)
		}
	}
	sort.Sort(ExamplesByID(orderedExamples))
	sort.Sort(ExamplesByID(negativeExamples))

	orderedExamples = append(orderedExamples, otherExamples...)
	orderedExamples = append(orderedExamples, negativeExamples...)
	return orderedExamples
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
	allDictionaries := orderDictionaries(dictionaries, userSubmission)
	results = Results{headword, allDictionaries, synonyms, antonyms, highlightedExamples}
	return c.JSON(http.StatusOK, results)
}

func recaptchaValidate(recaptchaResponse string) (r RecaptchaResponse) {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {"6LfhfycTAAAAANpoGOMqHlrhPBQlQoAZwy_O-J5-"}, "response": {recaptchaResponse}})
	if err != nil {
		fmt.Printf("Post error: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read error: could not read body: %s", err)
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Printf("Read error: got invalid JSON: %s", err)
	}
	return
}

func submitDefinition(c echo.Context) error {
	recaptchaCheck := recaptchaValidate(c.FormValue("recaptchaResponse"))
	if recaptchaCheck.Success == false {
		message := map[string]string{"message": "Recaptcha error"}
		return c.JSON(http.StatusOK, message)
	}
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

func submitExample(c echo.Context) error {
	recaptchaCheck := recaptchaValidate(c.FormValue("recaptchaResponse"))
	if recaptchaCheck.Success == false {
		message := map[string]string{"message": "Recaptcha error"}
		return c.JSON(http.StatusOK, message)
	}
	headword := sanitize.HTML(c.FormValue("term"))
	source := sanitize.HTML(c.FormValue("source"))
	link := sanitize.HTML(c.FormValue("link"))
	allowedTags := []string{"i", "b"}
	example, htmlErr := sanitize.HTMLAllowing(c.FormValue("example"), allowedTags)
	if htmlErr != nil {
		fmt.Println(htmlErr)
		message := map[string]string{"message": "error"}
		return c.JSON(http.StatusOK, message)
	}

	if _, ok := headwordMap[headword]; ok {
		query := "SELECT examples FROM headwords WHERE headword=$1"
		var userExamples []Example
		err := pool.QueryRow(query, headword).Scan(&userExamples)
		if err != nil {
			fmt.Println(err)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}

		var storedIDs = make(map[int]bool)
		for _, storedExample := range userExamples {
			storedIDs[storedExample.ID] = true
		}
		id := int(rand.Int31())
		for {
			if _, ok := storedIDs[id]; !ok {
				break
			}
		}
		score := 0
		userExamples = append(userExamples, Example{example, link, score, id, 0.0, source})
		serializedSubmission, jsonError := json.Marshal(userExamples)
		if jsonError != nil {
			fmt.Println(jsonError)
		}

		update := "UPDATE headwords SET examples=$1 WHERE headword=$2"
		_, updateErr := pool.Exec(update, serializedSubmission, headword)
		if updateErr != nil {
			fmt.Println(updateErr)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}
	} else {
		fmt.Println("Headword does not exist for this example")
		message := map[string]string{"message": "error"}
		return c.JSON(http.StatusOK, message)
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

func submitNym(c echo.Context) error {
	recaptchaCheck := recaptchaValidate(c.FormValue("recaptchaResponse"))
	if recaptchaCheck.Success == false {
		message := map[string]string{"message": "Recaptcha error"}
		return c.JSON(http.StatusOK, message)
	}
	headword := sanitize.HTML(c.FormValue("term"))
	nym := sanitize.HTML(c.FormValue("nym"))
	typeOfNym := sanitize.HTML(c.FormValue("type"))
	if _, ok := headwordMap[headword]; ok {
		query := fmt.Sprintf("SELECT %s FROM headwords WHERE headword=$1", typeOfNym)
		var nyms []string
		err := pool.QueryRow(query, headword).Scan(&nyms)
		if err != nil {
			fmt.Println(err)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}

		var storedNyms = make(map[string]bool)
		for _, storedNym := range nyms {
			storedNyms[storedNym] = true
		}
		if _, ok := storedNyms[nym]; ok {
			fmt.Printf("%s is already in the DB\n", nym)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}
		nyms = append(nyms, nym)
		serializedSubmission, jsonError := json.Marshal(nyms)
		if jsonError != nil {
			fmt.Println(jsonError)
		}

		update := fmt.Sprintf("UPDATE headwords SET %s=$1 WHERE headword=$2", typeOfNym)
		_, updateErr := pool.Exec(update, serializedSubmission, headword)
		if updateErr != nil {
			fmt.Println(updateErr)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}
	} else {
		fmt.Println("Headword does not exist for this example")
		message := map[string]string{"message": "error"}
		return c.JSON(http.StatusOK, message)
	}
	message := map[string]string{"message": "success"}
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
	e.GET("/exemple", index)
	e.GET("/synonyme", index)
	e.GET("/antonyme", index)

	// API
	e.GET("/api/mot/:headword", query)
	e.GET("/api/wordwheel", func(c echo.Context) error {
		return c.JSON(http.StatusOK, headwordList)
	})
	e.GET("/api/autocomplete/:prefix", autoComplete)
	e.POST("/api/submit", submitDefinition)
	e.GET("/api/vote/:headword/:id/:vote", voteForExample)
	e.POST("/api/submitExample", submitExample)
	e.POST("/api/submitNym", submitNym)

	// Start server
	e.Run(standard.New(":8080"))
}
