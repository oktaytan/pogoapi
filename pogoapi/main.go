package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"github.com/oktaytan/pogoapi/models"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)



// Post User Struct
type Author struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

// Post Comment Struct
type Comment struct {
	ID            string `json:"id"`
	Comment       string `json:"comment"`
	CommentAuthor Author `json:"author"`
}

type Comments []Comment

// POST Struct
type Post struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Likes     string   `json:"likes"`
	Author    Author   `json:"author"`
	Comments  Comments `json:"comments"`
}

type Posts []Post

// POST Struct
type OwnPost struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Likes     string   `json:"likes"`
	Comments  Comments `json:"comments"`
}

type OwnPosts []OwnPost

// POST Struct
type UpdatedPost struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserID    string `json:"user_id"`
	Likes     string `json:"likes"`
}

type Error struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

var mainDB *sql.DB

// Login
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
}

// Tüm kullanıcılar
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	rows, err := mainDB.Query("SELECT * FROM users")
	checkErr(err)
	var users Users
	for rows.Next() {
		var user User
		var password string
		err = rows.Scan(&user.ID, &user.UserName, &user.Email, &password)
		checkErr(err)
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

// Tüm postlar
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")

	// Query
	rows, err := mainDB.Query("SELECT * FROM posts AS p INNER JOIN (SELECT id, username, email FROM users) AS u ON p.user_id = u.id ORDER BY p.created_at DESC")

	checkErr(err)

	var posts Posts

	for rows.Next() {
		var post Post
		var author Author
		var user_id string
		var post_id string
		var comment_post_id string
		var comments Comments

		err = rows.Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt, &user_id, &post.Likes, &author.ID, &author.UserName, &author.Email)

		checkErr(err)

		post.Author = author

		stmt, err := mainDB.Prepare("SELECT * FROM comments AS c INNER JOIN ( SELECT id, username, email FROM users) u ON c.user_id = u.id WHERE c.post_id = ?")

		var commentsRow *sql.Rows
		commentsRow, err = stmt.Query(post.ID)

		for commentsRow.Next() {
			var comment Comment
			var commentAuthor Author

			err = commentsRow.Scan(&comment.ID, &comment.Comment, &post_id, &comment_post_id, &commentAuthor.ID, &commentAuthor.UserName, &commentAuthor.Email)

			checkErr(err)
			comment.CommentAuthor = commentAuthor
			comments = append(comments, comment)
		}
		post.Comments = comments
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}

// Kullanıcının kendi postları
func getOwnPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	params := mux.Vars(r)
	stmt, err := mainDB.Prepare("SELECT posts.id, posts.title, posts.body, posts.created_at, posts.updated_at, posts.user_id, posts.likes, users.id FROM posts INNER JOIN users ON users.id = posts.user_id WHERE users.username = ? ORDER BY posts.created_at DESC")
	checkErr(err)
	rows, errQuery := stmt.Query(params["username"])
	checkErr(errQuery)
	var ownPosts OwnPosts
	for rows.Next() {
		var post OwnPost
		var post_user_id string
		var comment_post_id string
		var post_id string
		var user_id string
		var comments Comments

		err = rows.Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt, &post_user_id, &post.Likes, &user_id)
		checkErr(err)

		stmt, err := mainDB.Prepare("SELECT * FROM comments AS c INNER JOIN ( SELECT id, username, email FROM users) u ON c.user_id = u.id WHERE c.post_id = ?")

		var commentsRow *sql.Rows
		commentsRow, err = stmt.Query(post.ID)

		for commentsRow.Next() {
			var comment Comment
			var commentAuthor Author

			err = commentsRow.Scan(&comment.ID, &comment.Comment, &post_id, &comment_post_id, &commentAuthor.ID, &commentAuthor.UserName, &commentAuthor.Email)

			checkErr(err)
			comment.CommentAuthor = commentAuthor
			comments = append(comments, comment)
		}
		post.Comments = comments
		ownPosts = append(ownPosts, post)
	}
	json.NewEncoder(w).Encode(ownPosts)
}

// Tek post
func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	params := mux.Vars(r)
	stmt, err := mainDB.Prepare("SELECT * FROM posts AS p INNER JOIN (SELECT id, username, email from users) AS u ON p.user_id = u.id WHERE p.id = ?")
	checkErr(err)
	rows, errQuery := stmt.Query(params["id"])
	checkErr(errQuery)
	var post Post
	var author Author

	for rows.Next() {
		var user_id string
		var post_id string
		var comments Comments
		var comment_post_id string

		err = rows.Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt, &user_id, &post.Likes, &author.ID, &author.UserName, &author.Email)
		checkErr(err)

		post.Author = author

		var stmtNew *sql.Stmt
		stmtNew, err := mainDB.Prepare("SELECT * FROM comments AS c INNER JOIN ( SELECT id, username, email FROM users) u ON c.user_id = u.id WHERE c.post_id = ?")

		var commentsRow *sql.Rows
		commentsRow, err = stmtNew.Query(post.ID)

		for commentsRow.Next() {
			var comment Comment
			var commentAuthor Author

			err = commentsRow.Scan(&comment.ID, &comment.Comment, &post_id, &comment_post_id, &commentAuthor.ID, &commentAuthor.UserName, &commentAuthor.Email)

			checkErr(err)
			comment.CommentAuthor = commentAuthor
			comments = append(comments, comment)
			post.Comments = comments
		}
	}
	isID, _ := strconv.ParseInt(post.ID, 10, 0)
	if isID == 0 {
		var error Error
		error.Error = true
		error.Message = "Post not found"
		json.NewEncoder(w).Encode(error)
		return
	} else {
		json.NewEncoder(w).Encode(post)
	}
}

// Send Post Struct
type SendPost struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

// Post Ekle
func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	var post SendPost
	_ = json.NewDecoder(r.Body).Decode(&post)
	post.ID = strconv.Itoa(rand.Intn(100000000000))
	stmt, err := mainDB.Prepare("INSERT INTO posts(id, title, body, user_id) VALUES (?, ?, ?, ?)")
	checkErr(err)
	_, errExec := stmt.Exec(post.ID, post.Title, post.Body, post.UserID)
	checkErr(errExec)
	json.NewEncoder(w).Encode(post)
}

// Post güncelle
func updatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	params := mux.Vars(r)

	var updatedPost UpdatedPost
	_ = json.NewDecoder(r.Body).Decode(&updatedPost)

	stmt, err := mainDB.Prepare("UPDATE posts SET id = ?, title = ?, body = ?, created_at = ?, updated_at = ?, user_id = ?, likes = ? WHERE id = ?")
	checkErr(err)
	result, errExec := stmt.Exec(updatedPost.ID, updatedPost.Title, updatedPost.Body, updatedPost.CreatedAt, updatedPost.UpdatedAt, updatedPost.UserID, updatedPost.Likes, params["id"])
	checkErr(errExec)
	rowAffected, errLast := result.RowsAffected()
	checkErr(errLast)
	if rowAffected == 0 {
		var error Error
		error.Error = true
		error.Message = "Post not found"
		json.NewEncoder(w).Encode(error)
		return
	} else {
		json.NewEncoder(w).Encode(updatedPost)
	}
}

// Post sil
func deletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	params := mux.Vars(r)
	stmt, err := mainDB.Prepare("DELETE FROM posts WHERE id = ?")
	checkErr(err)
	result, errExec := stmt.Exec(params["id"])
	checkErr(errExec)
	rows, errRow := result.RowsAffected()
	log.Println(rows)
	checkErr(errRow)
	if rows == 0 {
		var error Error
		error.Error = true
		error.Message = "Post not found"
		json.NewEncoder(w).Encode(error)
		return
	} else {
		json.NewEncoder(w).Encode("Post Deleted")
	}
}

func main() {
	// Sqlite db bağlantısı yapılıyor
	db, errOpenDB := sql.Open("sqlite3", "./db/pogo.db")
	checkErr(errOpenDB)
	mainDB = db

	// Router örneği oluşturuluyor
	r := mux.NewRouter()

	// Kullanılacak endpointler
	r.HandleFunc("/api/login", login).Methods("POST")
	r.HandleFunc("/api/users", getUsers).Methods("GET")
	r.HandleFunc("/api/posts", getPosts).Methods("GET")
	r.HandleFunc("/api/{username}/posts", getOwnPosts).Methods("GET")
	r.HandleFunc("/api/posts/{id}", getPost).Methods("GET")
	r.HandleFunc("/api/posts", createPost).Methods("POST")
	r.HandleFunc("/api/posts/{id}", updatePost).Methods("PUT")
	r.HandleFunc("/api/posts/{id}", deletePost).Methods("DELETE")

	log.Print("Server running on port 5000")

	log.Fatal(http.ListenAndServe(":5000", r))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
