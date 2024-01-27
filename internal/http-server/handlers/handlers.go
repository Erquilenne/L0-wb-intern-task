package handlers

import (
	"L0-wb-intern-task/internal/storage/cache"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetOrderHandler(w http.ResponseWriter, r *http.Request, cache *cache.Cache) {
	orderIDStr := r.URL.Path[len("/order/"):]
	log.Println("in GetOrderHandler")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	order, exist := cache.GetOrder(orderID)
	if !exist {
		http.Error(w, "Order with this ID does not exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("in HomeHandler")
	http.ServeFile(w, r, "./frontend/index.html")
}
