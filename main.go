package main

import (
	"encoding/json"
	"fmt"
	"html"
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

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/agext/levenshtein"
	"github.com/jackc/pgx"
	"github.com/kennygrant/sanitize"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

// App config
type config struct {
	DatabaseName     string `json:"databaseName"`
	DatabaseUser     string `json:"user"`
	DatabasePassword string `json:"password"`
	Debug            bool   `json:"debug"`
}

// Words of the day
type wordOfTheDay struct {
	Headword string `json:"headword"`
	Date     string `json:"date"`
}

// Results to export
type Results struct {
	Headword         string         `json:"headword"`
	Dictionaries     DictionaryData `json:"dictionaries"`
	Synonyms         []Nym          `json:"synonyms"`
	Antonyms         []Nym          `json:"antonyms"`
	Examples         []Example      `json:"examples"`
	TimeSeries       [][]float64    `json:"timeSeries"`
	Collocates       []Collocates   `json:"collocates"`
	NearestNeighbors []string       `json:"nearestNeighbors"`
	FuzzyResults     []FuzzyResult  `json:"fuzzyResults"`
}

// Dictionary to export
type Dictionary struct {
	Name       string              `json:"name"`
	Label      string              `json:"label"`
	ShortLabel string              `json:"shortLabel"`
	Content    []map[string]string `json:"contentObj"`
	Show       bool                `json:"show"`
}

// Collocates Type
type Collocates struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

// DictionaryData to export
type DictionaryData struct {
	Data         []Dictionary `json:"data"`
	TotalDicos   int          `json:"totalDicos"`
	TotalEntries int          `json:"totalEntries"`
}

// Example for headwords
type Example struct {
	Content    string `json:"content"`
	Link       string `json:"link"`
	Score      int    `json:"score"`
	ID         int    `json:"id"`
	Source     string `json:"source"`
	UserSubmit bool   `json:"userSubmit"`
	Date       string `json:"date"`
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
	Date    string `json:"date"`
}

// Nym type
type Nym struct {
	Label      string `json:"label"`
	UserSubmit bool   `json:"UserSubmit"`
	Date       string `json:"date"`
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

// FuzzyResult is the result of fuzzy searching
type FuzzyResult struct {
	Word  string  `json:"word"`
	Score float64 `json:"score"`
}

var defaultConnConfig pgx.ConnConfig
var pool = createConnPool()

var headwordList, headwordMap = getAllHeadwords()

var tokenRegex = regexp.MustCompile(`(?i)([\p{L}]+)|([\.?,;:'’!\-]+)|([\s]+)|([\d]+)`)

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

var webConfig = loadConfig()

var appJs = []string{
	"static/js/jquery.min.js",
	"static/js/bootstrap.min.js",
	"static/js/angular.min.js",
	"static/js/angular-route.min.js",
	"static/js/angular-resource.min.js",
	"static/js/angular-route.min.js",
	"static/js/angular-animate.min.js",
	"static/js/angular-touch.min.js",
	"static/js/angular-sanitize.min.js",
	"static/js/angular-cookies.min.js",
	"static/js/Chart.min.js",
	"static/js/angucomplete-alt.min.js",
	"static/js/sticky.min.js",
	"static/js/angular-recaptcha.min.js",
	"app/app.js",
	"app/config.js",
	"app/filters.js",
	"app/values.js",
	"app/components/apropos/aproposDirective.js",
	"app/components/results/resultsController.js",
	"app/components/dictionaryEntries/dictionaryEntriesDirective.js",
	"app/components/wordwheel/wordwheelDirective.js",
	"app/components/synAntoNyms/synAntoNymsDirective.js",
	"app/components/examples/examplesDirective.js",
	"app/components/newDefinition/newDefinitionController.js",
	"app/components/newExample/newExampleController.js",
	"app/components/newSynAnto/newSynAntoController.js",
	"app/components/timeSeries/timeSeriesDirective.js",
	"app/components/total/totalDirective.js",
	"app/components/collocations/collocationDirective.js",
	"app/components/nearestNeighbors/nearestNeighborsDirective.js",
}

var appCSS = []string{
	"static/css/bootstrap.min.css",
	"static/css/style.css",
}

var indexHTML, mainJS, mainCSS = getIndexHTML()

var fuzzySearchParams = levenshtein.NewParams()

var wordsOfTheDay = loadWordsOfTheDay()

func loadWordsOfTheDay() map[string]string {
	wordOfTheDayFile, err := os.Open("words_of_the_day.json")
	if err != nil {
		fmt.Println("opening words of the day file", err.Error())
	}
	var wordList []wordOfTheDay
	jsonParser := json.NewDecoder(wordOfTheDayFile)
	if err = jsonParser.Decode(&wordList); err != nil {
		fmt.Println("parsing word of the day file", err.Error())
	}
	dateToWords := make(map[string]string)
	for _, wordElement := range wordList {
		dateToWords[wordElement.Date] = wordElement.Headword
	}
	return dateToWords
}

func getIndexHTML() (string, string, string) {
	t := time.Now()
	secs := t.Unix()
	suffix := strconv.Itoa(int(secs))
	dvlfCSSPath := fmt.Sprintf("public/static/css/dvlf-%s.css", suffix)
	dvlfJsPath := fmt.Sprintf("public/static/js/dvlf-%s.js", suffix)
	indexByte, _ := ioutil.ReadFile("public/index.html")
	index := string(indexByte)
	cssPaths := ""
	javascript := ""
	if webConfig.Debug == true {
		for _, jsFile := range appJs {
			javascript += fmt.Sprintf("<script src='%s'></script>", jsFile)
		}
		for _, cssFile := range appCSS {
			cssPaths += fmt.Sprintf("<link href='%s' rel='stylesheet'>", cssFile)
		}
	} else {
		m := minify.New()
		m.AddFunc("text/css", css.Minify)
		m.AddFunc("text/javascript", js.Minify)
		jsCode := ""
		for _, jsFile := range appJs {
			jsByte, _ := ioutil.ReadFile(fmt.Sprintf("public/%s", jsFile))
			jsString := string(jsByte)
			minifiedJs, err := m.String("text/javascript", jsString)
			if err != nil {
				fmt.Println(err)
			}
			jsCode += minifiedJs
		}
		if _, notExists := os.Stat(dvlfJsPath); notExists == nil {
			os.Remove(dvlfJsPath)
		}
		f, err := os.Create(dvlfJsPath)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		_, writeErr := f.WriteString(jsCode)
		if writeErr != nil {
			fmt.Println(writeErr)
		}
		f.Sync()
		cssCode := ""
		for _, cssFile := range appCSS {
			cssByte, _ := ioutil.ReadFile(fmt.Sprintf("public/%s", cssFile))
			cssString := string(cssByte)
			minifiedCSS, err := m.String("text/css", cssString)
			if err != nil {
				fmt.Println(err)
			}
			cssCode += minifiedCSS
		}
		if _, notExists := os.Stat(dvlfCSSPath); notExists == nil {
			os.Remove(dvlfCSSPath)
		}
		c, err := os.Create(dvlfCSSPath)
		if err != nil {
			fmt.Println(err)
		}
		defer c.Close()
		_, cssWriteErr := c.WriteString(cssCode)
		if cssWriteErr != nil {
			fmt.Println(cssWriteErr)
		}
		c.Sync()
		javascript = fmt.Sprintf("<script async src='static/js/dvlf-%s.js'></script>", suffix)
		cssPaths = fmt.Sprintf("<link href='static/css/dvlf-%s.css' rel='stylesheet'>", suffix)
	}
	index = strings.Replace(index, "$APP_CSS$", cssPaths, 1)
	index = strings.Replace(index, "$APP_SCRIPTS$", javascript, 1)
	return index, dvlfJsPath, dvlfCSSPath
}

func loadConfig() config {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}
	var settings config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&settings); err != nil {
		fmt.Println("parsing config file", err.Error())
	}
	return settings
}

func createConnPool() *pgx.ConnPool {
	defaultConnConfig.Host = "localhost"
	defaultConnConfig.Database = webConfig.DatabaseName
	defaultConnConfig.User = webConfig.DatabaseUser
	defaultConnConfig.Password = webConfig.DatabasePassword
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

func getAllHeadwords() ([]string, map[string]int) {
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
	var headwordHash = make(map[string]int)
	for pos, value := range headwords {
		headwordHash[value] = pos
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
			content = append(content, map[string]string{"content": entry.Content, "source": entry.Source, "link": entry.Link, "date": entry.Date})
		}
		newDicos = append(newDicos, Dictionary{"userSubmit", "Définition(s) d'utilisateurs/trices", "Définition(s) d'utilisateurs/trices", content, show})
	} else {
		newDicos = append(newDicos, Dictionary{"userSubmit", "Définition(s) d'utilisateurs/trices", "Définition(s) d'utilisateurs/trices", make([]map[string]string, 0), true})
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

	formRegex := regexp.MustCompile("(?i)^(" + strings.Join(forms, "|") + ")$")
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
	var otherExamples []Example
	var userExamplesWithNoScore []Example
	for _, example := range examples {
		if example.Score > 0 {
			orderedExamples = append(orderedExamples, example)
		} else if example.UserSubmit && example.Score == 0 {
			userExamplesWithNoScore = append(userExamplesWithNoScore, example)
		} else if example.Score == 0 {
			otherExamples = append(otherExamples, example)
		}
	}

	sort.Sort(ExamplesByID(orderedExamples))

	orderedExamples = append(orderedExamples, userExamplesWithNoScore...)
	orderedExamples = append(orderedExamples, otherExamples...)
	if len(orderedExamples) > 30 {
		return orderedExamples[:30]
	}
	return orderedExamples
}

func getWordwheel(c echo.Context) error {
	headword, _ := url.QueryUnescape(c.FormValue("headword"))
	if _, ok := headwordMap[headword]; ok {
		index := headwordMap[headword]
		startIndex := index - 100
		if startIndex < 0 {
			startIndex = 0
		}
		endIndex := index + 100
		if endIndex > len(headwordList) {
			endIndex = len(headwordList) - 1
		}
		return c.JSON(http.StatusOK, headwordList[startIndex:endIndex])
	}
	return c.JSON(http.StatusOK, headwordList[0:200])
}

func query(c echo.Context) error {
	headword, _ := url.QueryUnescape(c.Param("headword"))
	query := "SELECT user_submit, dictionaries, synonyms, antonyms, examples, time_series, collocations, nearest_neighbors FROM headwords WHERE headword=$1"
	var results Results
	var dictionaries map[string][]string
	var synonyms []Nym
	var antonyms []Nym
	var userSubmission []UserSubmit
	var examples []Example
	var timeSeries [][]float64
	var collocations []Collocates
	var nearestNeighbors []string
	fuzzyResults := []FuzzyResult{}
	err := pool.QueryRow(query, headword).Scan(&userSubmission, &dictionaries, &synonyms, &antonyms, &examples, &timeSeries, &collocations, &nearestNeighbors)
	if err != nil {
		fmt.Println(err)
		fuzzyResults = getSimilarHeadWords(headword)
		empty := Results{headword, DictionaryData{[]Dictionary{}, 0, 0}, []Nym{}, []Nym{}, []Example{}, [][]float64{}, []Collocates{}, []string{}, fuzzyResults}
		return c.JSON(http.StatusOK, empty)
	}
	highlightedExamples := highlightExamples(examples, headword)
	sortedExamples := sortExamples(highlightedExamples)
	allDictionaries := orderDictionaries(dictionaries, userSubmission)
	if allDictionaries.TotalEntries < 2 {
		fuzzyResults = getSimilarHeadWords(headword)
	}
	results = Results{headword, allDictionaries, synonyms, antonyms, sortedExamples, timeSeries, collocations, nearestNeighbors, fuzzyResults}
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

func convertTimeStampToString(timeStamp string) string {
	months := []string{"Janvier", "Février", "Mars", "Avril", "Mai", "Juin", "Juillet", "Août", "Septembre", "Octobre", "Novembre", "Décembre"}
	year := strings.Split(timeStamp, "-")[0]
	month := strings.Split(timeStamp, "-")[1]
	if strings.HasPrefix(month, "0") {
		month = month[1:]
	}
	monthIndex, indexErr := strconv.Atoi(month)
	if indexErr != nil {
		fmt.Printf("Index error: %s", indexErr)
		return timeStamp
	}
	month = months[monthIndex-1]
	day := strings.Split(timeStamp, "-")[2]
	if strings.HasPrefix(day, "0") {
		day = day[1:]
	}
	if day == "1" {
		day = "1er"
	}
	fullDate := fmt.Sprintf("%s %s %s", day, month, year)
	return fullDate
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
	definition = html.UnescapeString(definition)
	if htmlErr != nil {
		fmt.Println(htmlErr)
		message := map[string]string{"message": "error"}
		return c.JSON(http.StatusOK, message)
	}
	timeStamp := convertTimeStampToString(strings.Split(time.Now().String(), " ")[0])
	var newSubmission = UserSubmit{definition, source, link, timeStamp}
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
		update := "UPDATE headwords SET user_submit=$1 WHERE headword=$2"
		_, updateErr := pool.Exec(update, userSubmission, headword)
		if updateErr != nil {
			fmt.Println(updateErr)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}
	} else {
		var userSubmission = []UserSubmit{newSubmission}
		dictionaries := make(map[string][]string)
		synonyms := []string{}
		antonyms := []string{}
		examples := []Example{}
		insert := "INSERT INTO headwords (headword, dictionaries, synonyms, antonyms, user_submit, examples) VALUES ($1, $2, $3, $4, $5, $6)"
		_, insertErr := pool.Exec(insert, headword, dictionaries, synonyms, antonyms, userSubmission, examples)
		if insertErr != nil {
			fmt.Println(insertErr)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}
		headwordList = append(headwordList, headword)
		cl := collate.New(language.French, collate.Loose, collate.IgnoreCase)
		cl.SortStrings(headwordList)
		for index, word := range headwordList {
			if word == headword {
				headwordMap[headword] = index
				break
			}
		}
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
	example = html.UnescapeString(example)
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
		timeStamp := convertTimeStampToString(strings.Split(time.Now().String(), " ")[0])
		userExamples = append(userExamples, Example{example, link, score, id, source, true, timeStamp})
		update := "UPDATE headwords SET examples=$1 WHERE headword=$2"
		_, updateErr := pool.Exec(update, userExamples, headword)
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
	timeStamp := convertTimeStampToString(strings.Split(time.Now().String(), " ")[0])
	nym := Nym{sanitize.HTML(c.FormValue("nym")), true, timeStamp}
	typeOfNym := sanitize.HTML(c.FormValue("type"))
	if _, ok := headwordMap[headword]; ok {
		query := fmt.Sprintf("SELECT %s FROM headwords WHERE headword=$1", typeOfNym)
		var nyms []Nym
		err := pool.QueryRow(query, headword).Scan(&nyms)
		if err != nil {
			fmt.Println(err)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}

		var storedNyms = make(map[string]bool)
		for _, storedNym := range nyms {
			storedNyms[storedNym.Label] = true
		}
		if _, ok := storedNyms[nym.Label]; ok {
			fmt.Printf("%s is already in the DB\n", nym.Label)
			message := map[string]string{"message": "error"}
			return c.JSON(http.StatusOK, message)
		}
		nyms = append(nyms, nym)
		update := fmt.Sprintf("UPDATE headwords SET %s=$1 WHERE headword=$2", typeOfNym)
		_, updateErr := pool.Exec(update, nyms, headword)
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

func getSimilarHeadWords(queryTerm string) []FuzzyResult {
	results := []FuzzyResult{}
	normQueryTerm := norm.NFC.String(queryTerm)
	maxScore := float64(1)
	for _, word := range headwordList {
		normalizedWord := norm.NFC.String(word)
		score := levenshtein.Match(normQueryTerm, normalizedWord, fuzzySearchParams)
		if score >= 0.7 && score < maxScore {
			results = append(results, FuzzyResult{word, score})
		}
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	return results
}

func index(c echo.Context) error {
	return c.HTML(http.StatusOK, indexHTML)
}

func main() {
	// Echo instance
	e := echo.New()

	e.Static("/static", "public/static")
	e.Static("/app", "public/app")

	e.File("/dvlf.ico", "public/static/images/dvlf.ico")

	e.Debug = webConfig.Debug

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.Secure())

	// Enable TLS
	e.AutoTLSManager.Cache = autocert.DirCache(".cache")
	// Redirect http traffic to https
	e.Pre(middleware.NonWWWRedirect())
	e.Pre(middleware.HTTPSWWWRedirect())
	e.Pre(middleware.HTTPSRedirect())

	e.GET("/", index)
	e.GET("/mot/*", index)
	e.GET("/apropos", index)
	e.GET("/definition", index)
	e.GET("/exemple", index)
	e.GET("/synonyme", index)
	e.GET("/antonyme", index)

	e.GET("/static/js/dvlf*", func(c echo.Context) error {
		c.Response().Header().Add("Cache-Control", "max-age=946080000")
		return c.File(mainJS)
	})
	e.GET("/static/css/dvlf*", func(c echo.Context) error {
		c.Response().Header().Add("Cache-Control", "max-age=946080000")
		return c.File(mainCSS)
	})
	e.GET("static/images/dvlf_logo_medium_no_beta_transparent.png", func(c echo.Context) error {
		c.Response().Header().Add("Cache-Control", "max-age=2592000")
		return c.File("public/static/images/dvlf_logo_medium_no_beta_transparent.png")
	})
	e.GET("/static/fonts/glyphicons-halflings-regular.woff2", func(c echo.Context) error {
		c.Response().Header().Add("Cache-Control", "max-age=2592000")
		return c.File("public/static/fonts/glyphicons-halflings-regular.woff2")
	})
	e.GET("/static/fonts/glyphicons-halflings-regular.woff", func(c echo.Context) error {
		c.Response().Header().Add("Cache-Control", "max-age=2592000")
		return c.File("public/static/fonts/glyphicons-halflings-regular.woff")
	})

	// API
	e.GET("/api/mot/:headword", query)
	e.GET("/api/wordwheel", getWordwheel)
	e.GET("/api/autocomplete/:prefix", autoComplete)
	e.POST("/api/submit", submitDefinition)
	e.GET("/api/vote/:headword/:id/:vote", voteForExample)
	e.POST("/api/submitExample", submitExample)
	e.POST("/api/submitNym", submitNym)
	e.GET("/api/wordoftheday", func(c echo.Context) error {
		date := strings.Split(time.Now().String(), " ")[0]
		println(date)
		word := wordsOfTheDay[date]
		return c.JSON(http.StatusOK, word)
	})

	// Start server
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}
