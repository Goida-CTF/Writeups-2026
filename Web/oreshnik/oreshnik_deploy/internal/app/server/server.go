package server

import (
	"crypto/rsa"
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"oreshnik/internal/app/auth"
	"oreshnik/internal/app/db"
	"oreshnik/internal/app/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Server struct {
	DB        *gorm.DB
	PrivKey   *rsa.PrivateKey
	PubKey    *rsa.PublicKey
	Templates *template.Template
	Flag      string
	rand      *rand.Rand
}

func New() (*Server, error) {
	s := &Server{Flag: os.Getenv("FLAG")}
	if s.Flag == "" {
		s.Flag = "goidactf{fake_flag}"
	}

	s.Templates = template.Must(template.ParseGlob("templates/*.html"))

	var err error
	s.DB, err = db.ConnectAndMigrate(&models.User{}, &models.RevokedToken{}, &models.Product{}, &models.Purchase{})
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat("public"); os.IsNotExist(err) {
		_ = os.Mkdir("public", 0755)
	}

	s.PrivKey, s.PubKey, err = auth.LoadKeys()
	if err != nil {
		return nil, err
	}

	s.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	if err := s.seed(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", s.homeHandler)
	mux.HandleFunc("/register", s.registerHandler)
	mux.HandleFunc("/login", s.loginHandler)
	mux.HandleFunc("/admin", s.adminHandler)
	mux.HandleFunc("/revoked", s.revokedHandler)
	mux.HandleFunc("/product/", s.productHandler)
	mux.HandleFunc("/buy/", s.buyHandler)
	mux.HandleFunc("/logout", s.logoutHandler)
	mux.HandleFunc("/my-purchases", s.myPurchasesHandler)
	mux.HandleFunc("/order/", s.orderHandler)
	return mux
}

func (s *Server) randomPassword(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[s.rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *Server) seed() error {
	s.DB.Exec("DELETE FROM revoked_tokens")
	s.DB.Exec("DELETE FROM products")
	s.DB.Exec("DELETE FROM users")

	products := []models.Product{
		{Name: "Грецкий орех", Description: "Классический орех, богат омега-3. Отлично подходит для выпечки и салатов.", Price: 250.0},
		{Name: "Миндаль", Description: "Полезен для сердца и кожи. Идеален в качестве перекуса или добавки в мюсли.", Price: 400.0},
		{Name: "Кешью", Description: "Сладкий и маслянистый, идеален для закусок и азиатских блюд.", Price: 550.0},
		{Name: "Фисташки", Description: "Соленые и хрустящие, прекрасная закуска к напиткам.", Price: 600.0},
		{Name: "Фундук", Description: "Лесной орех с насыщенным вкусом, хорош в шоколаде и десертах.", Price: 450.0},
	}
	if err := s.DB.Create(&products).Error; err != nil {
		return err
	}

	adminPass := s.randomPassword(16)
	hashAdmin, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
	log.Printf("Создан администратор с паролем: %s", adminPass)
	admin := models.User{Username: "админ", Password: string(hashAdmin), IsAdmin: true}
	if err := s.DB.Create(&admin).Error; err != nil {
		return err
	}

	users := []models.User{{Username: "пользователь1"}, {Username: "тест"}}
	for i := range users {
		p := s.randomPassword(16)
		h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		users[i].Password = string(h)
		if err := s.DB.Create(&users[i]).Error; err != nil {
			return err
		}
		log.Printf("Создан пользователь '%s' с паролем: %s", users[i].Username, p)
	}

	adminTok, _ := auth.GenerateJWT(admin.ID, admin.IsAdmin, s.PrivKey)
	if err := s.DB.Create(&models.RevokedToken{Token: adminTok}).Error; err != nil {
		return err
	}
	return nil
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	data := s.getTemplateData(r)
	if data["LoggedIn"].(bool) {
		var products []models.Product
		s.DB.Find(&products)
		data["Products"] = products
	}
	_ = s.Templates.ExecuteTemplate(w, "index.html", data)
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = s.Templates.ExecuteTemplate(w, "register.html", s.getTemplateData(r))
		return
	}
	_ = r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := models.User{Username: username, Password: string(h)}
	if err := s.DB.Create(&u).Error; err != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = s.Templates.ExecuteTemplate(w, "login.html", s.getTemplateData(r))
		return
	}
	_ = r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	var user models.User
	s.DB.First(&user, "username = ?", username)
	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	tok, _ := auth.GenerateJWT(user.ID, user.IsAdmin, s.PrivKey)
	http.SetCookie(w, &http.Cookie{Name: "token", Value: tok, Path: "/"})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) adminHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	tokenString := cookie.Value
	var revoked models.RevokedToken
	if s.DB.First(&revoked, "token = ?", tokenString).Error == nil {
		http.Error(w, "Token has been revoked", http.StatusForbidden)
		return
	}
	_, claims, err := auth.Parse(tokenString, s.PubKey)
	if err != nil || !claims["is_admin"].(bool) {
		http.Error(w, "Invalid token or not admin", http.StatusForbidden)
		return
	}
	data := s.getTemplateData(r)
	data["Flag"] = s.Flag
	_ = s.Templates.ExecuteTemplate(w, "admin.html", data)
}

func (s *Server) revokedHandler(w http.ResponseWriter, r *http.Request) {
	var tokens []models.RevokedToken
	s.DB.Find(&tokens)
	var tokenStrings []string
	for _, t := range tokens {
		tokenStrings = append(tokenStrings, t.Token)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tokenStrings)
}

func (s *Server) productHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/product/")
	var product models.Product
	s.DB.First(&product, id)
	if product.ID == 0 {
		http.NotFound(w, r)
		return
	}
	data := s.getTemplateData(r)
	data["Product"] = product
	_ = s.Templates.ExecuteTemplate(w, "product.html", data)
}

func (s *Server) buyHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	_, claims, err := auth.Parse(cookie.Value, s.PubKey)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userID := uint(claims["user_id"].(float64))
	productID := strings.TrimPrefix(r.URL.Path, "/buy/")

	purchase := models.Purchase{UserID: userID, ProductID: 0}
	s.DB.First(&models.Product{}, productID).Scan(&purchase.Product)
	if purchase.Product.ID == 0 {
		http.NotFound(w, r)
		return
	}
	purchase.ProductID = purchase.Product.ID

	s.DB.Create(&purchase)

	_ = s.Templates.ExecuteTemplate(w, "purchase_success.html", s.getTemplateData(r))
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/", Expires: time.Unix(0, 0)})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) myPurchasesHandler(w http.ResponseWriter, r *http.Request) {
	data := s.getTemplateData(r)
	if !data["LoggedIn"].(bool) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	var purchases []models.Purchase
	s.DB.Preload("Product").Where("user_id = ?", data["UserID"]).Find(&purchases)
	data["Purchases"] = purchases

	_ = s.Templates.ExecuteTemplate(w, "my_purchases.html", data)
}

func (s *Server) orderHandler(w http.ResponseWriter, r *http.Request) {
	data := s.getTemplateData(r)
	if !data["LoggedIn"].(bool) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	orderID := strings.TrimPrefix(r.URL.Path, "/order/")
	var purchase models.Purchase
	s.DB.Preload("Product").First(&purchase, orderID)

	data["Purchase"] = purchase
	_ = s.Templates.ExecuteTemplate(w, "order_detail.html", data)
}

func (s *Server) getTemplateData(r *http.Request) map[string]interface{} {
	data := map[string]interface{}{
		"LoggedIn": false,
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		return data
	}
	_, claims, err := auth.Parse(cookie.Value, s.PubKey)
	if err != nil {
		return data
	}
	var user models.User
	s.DB.First(&user, claims["user_id"])
	if user.ID != 0 {
		data["LoggedIn"] = true
		data["Username"] = user.Username
		data["UserID"] = user.ID
	}
	return data
}
