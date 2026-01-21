package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kadraman/InsecureGoMonorepo/pkg/database"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Product struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

type Order struct {
	UserID     int     `json:"user_id"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

func main() {
	// Start simple in-process HTTP servers that mimic the microservices.
	usersDB, _ := database.NewDatabase("", "", "", 0)
	productsDB, _ := database.NewDatabase("", "", "", 0)
	ordersDB, _ := database.NewDatabase("", "", "", 0)

	// Users server
	go func() {
		http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "method", http.StatusMethodNotAllowed)
				return
			}
			var u User
			_ = json.NewDecoder(r.Body).Decode(&u)
			hashed := database.HashPassword(u.Password)
			_ = usersDB.CreateUser(u.Username, u.Email, hashed)
			w.WriteHeader(http.StatusCreated)
		})

		http.HandleFunc("/users/id/", func(w http.ResponseWriter, r *http.Request) {
			// path: /users/id/{id}
			id := r.URL.Path[len("/users/id/"):]
			query := "SELECT id, username, email FROM users WHERE id = " + id
			results, _ := usersDB.ExecuteQuery(query)
			if results == nil || len(results) == 0 {
				http.NotFound(w, r)
				return
			}
			json.NewEncoder(w).Encode(results[0])
		})

		log.Println("Users demo server listening on :8081")
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	// Products server
	go func() {
		http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "method", http.StatusMethodNotAllowed)
				return
			}
			var p Product
			_ = json.NewDecoder(r.Body).Decode(&p)
			query := fmt.Sprintf("INSERT INTO products (name, description, price, category) VALUES ('%s','%s',%f,'%s')",
				p.Name, p.Description, p.Price, p.Category)
			_, _ = productsDB.ExecuteQuery(query)
			w.WriteHeader(http.StatusCreated)
		})

		http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
			// path: /products/{id}
			id := r.URL.Path[len("/products/"):]
			query := "SELECT * FROM products WHERE id = " + id
			results, _ := productsDB.ExecuteQuery(query)
			if results == nil || len(results) == 0 {
				http.NotFound(w, r)
				return
			}
			json.NewEncoder(w).Encode(results[0])
		})

		log.Println("Products demo server listening on :8082")
		log.Fatal(http.ListenAndServe(":8082", nil))
	}()

	// Orders server
	go func() {
		http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "method", http.StatusMethodNotAllowed)
				return
			}
			var o Order
			_ = json.NewDecoder(r.Body).Decode(&o)

			// Use timeouted client to snapshot
			client := &http.Client{Timeout: 2 * time.Second}

			userSnap := map[string]interface{}{}
			userResp, err := client.Get(fmt.Sprintf("http://localhost:8081/users/id/%d", o.UserID))
			if err == nil && userResp.StatusCode == http.StatusOK {
				defer userResp.Body.Close()
				_ = json.NewDecoder(userResp.Body).Decode(&userSnap)
			}

			prodSnap := map[string]interface{}{}
			prodResp, err := client.Get(fmt.Sprintf("http://localhost:8082/products/%d", o.ProductID))
			if err == nil && prodResp.StatusCode == http.StatusOK {
				defer prodResp.Body.Close()
				_ = json.NewDecoder(prodResp.Body).Decode(&prodSnap)
			}

			userJSON, _ := json.Marshal(userSnap)
			prodJSON, _ := json.Marshal(prodSnap)

			query := fmt.Sprintf("INSERT INTO orders (user_id, product_id, quantity, total_price, status, created_at, user_snapshot, product_snapshot) VALUES (%d, %d, %d, %f, 'pending', '%s', '%s', '%s')",
				o.UserID, o.ProductID, o.Quantity, o.TotalPrice, time.Now().Format(time.RFC3339), string(userJSON), string(prodJSON))
			_, _ = ordersDB.ExecuteQuery(query)
			w.WriteHeader(http.StatusCreated)
		})

		http.HandleFunc("/orders/list", func(w http.ResponseWriter, r *http.Request) {
			results, _ := ordersDB.ExecuteQuery("SELECT * FROM orders")
			json.NewEncoder(w).Encode(results)
		})

		log.Println("Orders demo server listening on :8083")
		log.Fatal(http.ListenAndServe(":8083", nil))
	}()

	// Allow servers to start
	time.Sleep(200 * time.Millisecond)

	// Demo sequence: create user, product, then order
	createUser := User{Username: "alice", Email: "alice@example.com", Password: "password123"}
	b, _ := json.Marshal(createUser)
	http.Post("http://localhost:8081/users", "application/json", bytes.NewReader(b))

	createProduct := Product{Name: "Demo Product", Description: "Demo", Price: 9.99, Category: "Demo"}
	pb, _ := json.Marshal(createProduct)
	http.Post("http://localhost:8082/products", "application/json", bytes.NewReader(pb))

	// create order snapshotting the above (IDs: user=1, product=1)
	order := Order{UserID: 1, ProductID: 1, Quantity: 2, TotalPrice: 19.98}
	ob, _ := json.Marshal(order)
	http.Post("http://localhost:8083/orders", "application/json", bytes.NewReader(ob))

	// Fetch orders and print
	time.Sleep(100 * time.Millisecond)
	resp, err := http.Get("http://localhost:8083/orders/list")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var orders []map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&orders)
	fmt.Printf("Stored orders: %+v\n", orders)

	// Keep demo servers running briefly
	time.Sleep(2 * time.Second)
}
