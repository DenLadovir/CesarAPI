package routes

import (
	"CesarAPI/database"
	"CesarAPI/models"
	"CesarAPI/utils"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func AuthRoutes(r *chi.Mux) {
	r.Post("/register", Register)
	r.Post("/login", Login)
}

func Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Register endpoint hit")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error("Error decoding JSON: ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Registering user: %s\n", user.Username)

	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		http.Error(w, "User already exists!", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		http.Error(w, "Error hashing password!", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	database.DB.Create(&user)
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var foundUser models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("JWT Validation error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := database.DB.Where("username = ?", user.Username).First(&foundUser).Error; err != nil {
		log.Printf("JWT Validation error: %v", err)
		http.Error(w, "Invalid credentials!", http.StatusUnauthorized)
		return
	}

	if foundUser.Username == "" {
		http.Error(w, "Invalid credentials!", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		log.Printf("JWT Validation error: %v", err)
		http.Error(w, "Invalid credentials!", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(foundUser.Username)
	if err != nil {
		log.Printf("JWT Validation error: %v", err)
		http.Error(w, "Could not generate token!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
