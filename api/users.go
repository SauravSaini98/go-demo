package api

import (
	"net/http"
	"encoding/json"
	"fmt"
	"strconv"
	"my_project/models"
	"my_project/serializers"
	"my_project/database"
)

// MyHandler is an HTTP handler function that responds with a custom message.
type Message struct {
    Text string `json:"text"`
}

type User struct {
    ID   int    `json:"id"`
    Email string `json:"email"`
}

const (
    defaultPage  = 1
    defaultLimit = 10
)

type Pagination struct {
    Count       int `json:"count"`
    CurrentPage int `json:"current_page"`
    Items       int `json:"items"`
}

type Response struct {
    Data []map[string]interface{} `json:"data"`
    Meta struct {
        Pagination Pagination `json:"pagination"`
    } `json:"meta"`
}

func HandleUsersApi(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        GetAllUsers(w, r)
    } else if r.Method == http.MethodPost {
        CreateUser(w, r)
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        fmt.Println("Request Error:", err, "Request Body:", r.Body) // Print the request body
        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
        return
    }

    var db = database.GetDB()
    // Insert the user record into the database
    insertQuery := "INSERT INTO users (name, email) VALUES ($1, $2)"
    _, err = db.Exec(insertQuery, user.Name, user.Email)
    if err != nil {
        http.Error(w, "Failed to insert data into the database", http.StatusInternalServerError)
        return
    }

    // Respond with a success message
    response := map[string]string{"message": "User created successfully"}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(jsonResponse)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
    var db = database.GetDB()
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")

    // Parse page and limit parameters with default values
    page := defaultPage
    limit := defaultLimit

    if pageStr != "" {
        page, _ = strconv.Atoi(pageStr)
    }

    if limitStr != "" {
        limit, _ = strconv.Atoi(limitStr)
    }

    // Calculate offset based on page and limit
    offset := (page - 1) * limit

    query := fmt.Sprintf("SELECT id, name, email FROM users order by id asc LIMIT %d OFFSET %d", limit, offset)
    rows, err := db.Query(query)
    if err != nil {
        http.Error(w, "Failed to execute query", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Fetch and serialize users
    var serializedUsers []map[string]interface{}
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.Name, &user.Email)
        if err != nil {
            http.Error(w, "Failed to scan row", http.StatusInternalServerError)
            return
        }

        serializer := serializers.NewUserSerializer(user)
        serializedUser := serializer.Serialize()
        serializedUsers = append(serializedUsers, serializedUser)
    }

    countQuery := "SELECT COUNT(*) FROM users"
    var totalCount int
    err = db.QueryRow(countQuery).Scan(&totalCount)
    if err != nil {
        http.Error(w, "Failed to retrieve total count", http.StatusInternalServerError)
        return
    }

    // Create the API response
    response := Response{
        Data: serializedUsers,
        Meta: struct {
            Pagination Pagination `json:"pagination"`
        }{
            Pagination: Pagination{
                Count:       totalCount,
                CurrentPage: page,
                Items:       limit,
            },
        },
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
