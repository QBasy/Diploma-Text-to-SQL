package controllers

import (
	"database-service/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type AuthController struct {
	JwtSecret []byte
	DB        *gorm.DB
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	if err := ac.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	databasePath := fmt.Sprintf("static/user_databases/%s.sqlite", user.ID)
	sqlDB, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		http.Error(w, "Failed to create user database", http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close()

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS example_table (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	if _, err := sqlDB.Exec(createTableQuery); err != nil {
		http.Error(w, "Failed to initialize user database", http.StatusInternalServerError)
		return
	}

	userDatabase := models.Database{
		UserID: user.ID,
		Path:   databasePath,
	}

	if err := ac.DB.Create(&userDatabase).Error; err != nil {
		http.Error(w, "Failed to save database info", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login user
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var reqUser models.User
	var dbUser models.User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if err := ac.DB.Where("email = ?", reqUser.Email).First(&dbUser).Error; err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(reqUser.Password)); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := ac.generateJWTToken(dbUser.ID)
	fmt.Println("hi")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (ac *AuthController) generateJWTToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	fmt.Println("3")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(ac.JwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
