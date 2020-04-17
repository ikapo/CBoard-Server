package pages

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

//------Constants--------
const baseAPIURL string = "http://api:8000/"
const templatesDir string = "./templates/"

//-----------------------

//------Globals--------
var templates map[string]*template.Template

//-----------------------

//-------Structs---------
type post struct {
	Id        uint32 `json:"post_id"`
	CreatedAt string `json:"created_at"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Board     uint8  `json:"board"`
	Ext       string `json:"ext"`
}

type comment struct {
	Id        uint32 `json:"com_id"`
	CreatedAt string `json:"created_at"`
	Content   string `json:"content"`
	Board     uint8  `json:"board"`
	Ext       string `json:"ext"`
}

type postsList struct {
	Posts    []post
	BoardNum uint8
}

type postAndComments struct {
	Post     post      `json:"post"`
	Comments []comment `json:"comments"`
}

//-----------------------

// logOnError logs err if err != nil
func logOnError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// makeContectRequest requests the API
// for content (posts/comments)
func makeContentRequest(url string) []byte {
	// Calling the API
	resp, err := http.Get(url)
	logOnError(err)

	// Reading the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	logOnError(err)

	return body

}

// getBoard gets and parses posts from a board
func getBoard(boardNumStr string) postsList {

	// Getting the posts JSON
	url := baseAPIURL + "board/" + boardNumStr
	body := makeContentRequest(url)
	var posts postsList
	num, err := strconv.Atoi(boardNumStr)
	if err != nil {
		log.Fatal(err)
	} else {
		posts.BoardNum = uint8(num)
		//Parsing the posts
		err := json.Unmarshal(body, &posts.Posts)
		logOnError(err)
	}

	return posts

}

// getPostPage gets requests and parses
// a post and its comments
func getPostPage(urlAddon string) postAndComments {

	// Getting the posts JSON
	url := baseAPIURL + urlAddon
	body := makeContentRequest(url)
	var postcom postAndComments
	err := json.Unmarshal(body, &postcom)
	logOnError(err)
	return postcom

}

// Init intializes the templates
func Init() {

	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	//Getting all the templates
	layouts, err := filepath.Glob(templatesDir + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates with base.tmpl
	for _, layout := range layouts {
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(layout, "base.tmpl"))
	}

}

// renderTemplate is a wrapper around template.ExecuteTemplate
// that allows template inheritance
func renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Rendering the template with base.tmpl and the data
	return tmpl.ExecuteTemplate(w, "base", data)
}

//-----Page Functions-------
// Index renders the index page (home page)
func Index(w http.ResponseWriter, r *http.Request) {
	// Executing the template
	renderTemplate(w, "index.tmpl", nil)
}

// Board renders the page for a board
func Board(w http.ResponseWriter, r *http.Request) {
	// Getting the board number
	vars := mux.Vars(r)
	board := vars["board"]
	postslist := getBoard(board)

	// Executing the template
	renderTemplate(w, "board.tmpl", postslist)
}

// Post renders the page for a post
func Post(w http.ResponseWriter, r *http.Request) {
	// Getting the parent post id
	vars := mux.Vars(r)
	id := vars["id"]

	// Getting parent post, if it is none set
	// The parent as an empty post
	postComs := getPostPage("post/" + id)

	if len(postComs.Comments) == 1 {
		if postComs.Comments[0].Id == 0 {
			var emptyComs []comment
			postComs.Comments = emptyComs

		}

	}
	renderTemplate(w, "post.tmpl", postComs)
}

//-----------------------------
