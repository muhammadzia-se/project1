package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

type Customer struct {
	gorm.Model  `json:"-"`
	CompanyId   int    `json:"company_id"`
	UserId      string `json:"user_id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	CreditCards string `json:"credit_cards"`
}

type Order struct {
	gorm.Model `json:"-"`
	Id         int    `json:"-"`
	OrderName  string `json:"order_name"`
	CustomerId string `json:"customer_id"`
	CreatedAt  string `json:"created_at"`
}

type OrderItem struct {
	gorm.Model   `json:"-"`
	OrderId      int     `json:"order_id"`
	Quantity     int     `json:"quantity"`
	Product      string  `json:"product"`
	PricePerUnit float32 `json:"price_per_unit"`
}

type CustomerCompany struct {
	gorm.Model  `json:"-"`
	CompanyId   int    `json:"company_id"`
	CompanyName string `json:"company_name"`
}

type Result struct {
	Product         string
	Quantity        int
	DeliveredAmount float32
	TotalAmount     float32
	CompanyName     string
	CustomerName    string
	OrderDate       string
	OrderName       string
}

// creating virtual object for the api response
type Response struct {
	Status bool     `json:"status"`
	Code   int      `json:"code"`
	Data   []Result `json:"data"`
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var results []Result
	DB.Table("order_items oi").Select("DISTINCT c.name as customer_name, cc.company_name,oi.product, oi.quantity, round((oi.quantity * oi.price_per_unit)::NUMERIC, 2) as delivered_amount, round((oi.quantity * oi.price_per_unit)::NUMERIC, 2) as total_amount, o.created_at as order_date, o.order_name").Joins("inner join orders o on oi.order_id=o.id").Joins("inner join customers c on o.customer_id=c.user_id").Joins("inner join customer_companies cc on c.company_id=cc.company_id").Find(&results)

	resp := Response{Status: true, Code: 200, Data: results}

	json.NewEncoder(w).Encode(resp)
}
func searchOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orderItems []OrderItem

	params := mux.Vars(r)
	// var orderItems OrderItem
	var p = params["pn"]

	DB.Where("product = ?", p).Find(&orderItems)
	json.NewEncoder(w).Encode(orderItems)
}

func initialMigration() {
	DNS := os.Getenv("DB_CONNECTION")
	DB, err = gorm.Open(postgres.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	} else {
		fmt.Println("DB is connected")
	}
	DB.AutoMigrate(&Customer{}, &Order{}, &OrderItem{}, &CustomerCompany{})
}
