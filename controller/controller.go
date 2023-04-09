package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sonus21/db-read-write/service"
	"log"
	"net/http"
)

func invalidBodyError() map[string]interface{} {
	return map[string]interface{}{
		"errors": map[string]interface{}{
			"code":    400,
			"message": "Invalid body",
		},
	}
}

func internalServerError() map[string]interface{} {
	return map[string]interface{}{
		"errors": map[string]interface{}{
			"code":    500,
			"message": "Internal server error",
		},
	}
}

func writeBody(statusCode int, w http.ResponseWriter, payload interface{}) {
	data, err := json.Marshal(payload)
	if err == nil {
		w.WriteHeader(statusCode)
		_, _ = w.Write(data)
	} else {
		log.Println("Error in JSON Marshal", err.Error())
		writeBody(http.StatusInternalServerError, w, internalServerError())
	}
}

func HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	// validate data
	req := service.OrderCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeBody(http.StatusBadRequest, w, invalidBodyError())
		return
	}
	resp, err := service.CreateOrder(r.Context(), &req)
	if err != nil {
		writeBody(http.StatusInternalServerError, w, internalServerError())
	} else {
		writeBody(http.StatusOK, w, resp)
	}
}

func OrderDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "orderId")
	resp, err := service.OrderDetail(r.Context(), id)
	// handle 404 here if you want to return 404 instead of 200 if order id is not found
	if err != nil {
		log.Println(err)
		writeBody(http.StatusInternalServerError, w, internalServerError())
	} else {
		writeBody(http.StatusOK, w, resp)
	}
}
