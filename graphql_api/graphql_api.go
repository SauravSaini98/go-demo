package graphql_api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "my_project/database"
    "github.com/graphql-go/graphql"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var userInputType = graphql.NewInputObject(graphql.InputObjectConfig{
    Name: "UserInput",
    Fields: graphql.InputObjectConfigFieldMap{
        "name": &graphql.InputObjectFieldConfig{
            Type: graphql.NewNonNull(graphql.String),
        },
        "email": &graphql.InputObjectFieldConfig{
            Type: graphql.NewNonNull(graphql.String),
        },
    },
})

var userType = graphql.NewObject(graphql.ObjectConfig{
    Name: "User",
    Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.Int,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
            "email": &graphql.Field{
                Type: graphql.String,
            },
    },
})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Mutation",
    Fields: graphql.Fields{
        "createUser": &graphql.Field{
            Type: userType,
            Args: graphql.FieldConfigArgument{
                "input": &graphql.ArgumentConfig{
                    Type: userInputType,
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                // Simulate database insert here
                input := p.Args["input"].(map[string]interface{})
                name, _ := input["name"].(string)
                email, _ := input["email"].(string)
                newUser := User{Name: name, Email: email}

                var db = database.GetDB() // Make sure you have a valid DB connection
                var insertedID int
                err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", newUser.Name, newUser.Email).Scan(&insertedID)
                if err != nil {
                    // Handle the error
                    return nil, err
                }
                newUser.ID =  insertedID
                return newUser, nil
            },
        },
    },
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Query",
    Fields: graphql.Fields{
        "getUser": &graphql.Field{
            Type: userType,
            Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                    Type: graphql.Int,
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                id, _ := p.Args["id"].(int)
                query := fmt.Sprintf("SELECT id, name, email FROM users where id = %d limit 1", id)
                var db = database.GetDB()
                row := db.QueryRow(query)
                var user User
                err := row.Scan(&user.ID, &user.Name, &user.Email)
                if err != nil {
                    // Handle the error appropriately
                    fmt.Println("Error scanning row:", err)
                    return nil, err
                }

                return user, nil
            },
        },
        "getUsers": &graphql.Field{
            Type: graphql.NewList(userType), // userType is the type for a single user
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                query := "SELECT id, name, email FROM users"
                var db = database.GetDB()
                rows, err := db.Query(query)
                if err != nil {
                    // Handle the error appropriately
                    fmt.Println("Error querying users:", err)
                    return nil, err
                }
                defer rows.Close()

                var users []User
                for rows.Next() {
                    var user User
                    err := rows.Scan(&user.ID, &user.Name, &user.Email)
                    if err != nil {
                        // Handle the error appropriately
                        fmt.Println("Error scanning user row:", err)
                        return nil, err
                    }
                    users = append(users, user)
                }

                return users, nil
            },
        },
    },
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: queryType,
    Mutation: mutationType,
})

func HandleGraphQL(w http.ResponseWriter, r *http.Request) {

    var requestBody struct {
        Query string `json:"query"`
    }

    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        fmt.Println(err, "error message Invalid request body")
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    result := graphql.Do(graphql.Params{
            Schema:        schema,
        RequestString: requestBody.Query,
    })


    if len(result.Errors) > 0 {
        for _, err := range result.Errors {
            fmt.Println(err.Message, "error message")
        }
    }

    jsonResponse, _ := json.Marshal(result)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}