package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hellwind2019/ecomm/cmd/ecomm-api/server"
	"github.com/hellwind2019/ecomm/cmd/ecomm-api/storer"
	"github.com/hellwind2019/ecomm/token"
	"github.com/hellwind2019/ecomm/util"
)

type Handler struct {
	ctx        context.Context
	server     *server.Server
	TokenMaker *token.JWTMaker
}

func NewHandler(srv *server.Server, secretKey string) *Handler {
	return &Handler{
		ctx:        context.Background(),
		server:     srv,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var p ProductRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	product, err := h.server.CreateProduct(h.ctx, toStoreProduct(p))
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}
	res := toResponseProduct(*product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// /product/{id}
func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Error parsing ID", http.StatusBadRequest)
		return
	}

	product, err := h.server.GetProduct(h.ctx, i)
	if err != nil {
		http.Error(w, "Failed to get product", http.StatusInternalServerError)
		return
	}
	res := toResponseProduct(*product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}
func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.server.ListProducts(h.ctx)
	if err != nil {
		http.Error(w, "Failed to list products", http.StatusInternalServerError)
		return
	}
	var res []ProductResponse
	for _, p := range products {
		res = append(res, *toResponseProduct(p))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Error parsing ID", http.StatusBadRequest)
		return
	}
	var p ProductRequest
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	product, err := h.server.GetProduct(h.ctx, i)
	if err != nil {
		http.Error(w, "Failed to get product", http.StatusInternalServerError)
		return
	}
	//patch the product with new values
	pathcProductReq(product, p)
	updated, err := h.server.UpdateProduct(h.ctx, product)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}
	res := toResponseProduct(*updated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}
func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Error parsing ID", http.StatusBadRequest)
		return
	}
	err = h.server.DeleteProduct(h.ctx, i)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func pathcProductReq(product *storer.Product, p ProductRequest) {
	if p.Name != "" {
		product.Name = p.Name
	}
	if p.Image != "" {
		product.Image = p.Image
	}
	if p.Category != "" {
		product.Category = p.Category
	}
	if p.Description != "" {
		product.Description = p.Description
	}
	if p.Rating != 0 {
		product.Rating = p.Rating
	}
	if p.NumReviews != 0 {
		product.NumReviews = p.NumReviews
	}
	if p.Price != 0 {
		product.Price = p.Price
	}
	if p.CountInStock != 0 {
		product.CountInStock = p.CountInStock
	}
	product.UpdatedAt = toTimePtr(time.Now())
}
func toTimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
func toStoreProduct(p ProductRequest) *storer.Product {
	return &storer.Product{
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}
func toResponseProduct(p storer.Product) *ProductResponse {
	return &ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var o OrderReq
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(authKey{}).(*token.UserClaims)
	so := toStorerOrder(o)
	so.UserID = claims.ID
	order, err := h.server.CreateOrder(h.ctx, so)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := toOrderResponse(order)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)	

	order, err := h.server.GetOrder(h.ctx, claims.ID)
	if err != nil {
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	res := toOrderResponse(order)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) listOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.server.ListOrders(h.ctx)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var res []OrderResponse
	for _, o := range orders {
		res = append(res, toOrderResponse(&o))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
func toStorerOrder(o OrderReq) *storer.Order {
	return &storer.Order{
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Items:         toStorerOrderItems(o.Items),
	}
}
func toStorerOrderItems(items []OrderItem) []storer.OrderItem {
	var storeItems []storer.OrderItem
	for _, item := range items {
		storeItems = append(storeItems, storer.OrderItem{
			Name:      item.Name,
			Quantity:  item.Quantity,
			Image:     item.Image,
			Price:     item.Price,
			ProductID: item.ProductID,
		})
	}
	return storeItems
}
func toOrderResponse(o *storer.Order) OrderResponse {
	return OrderResponse{
		ID:            o.ID,
		Items:         toOrderRespItems(o.Items),
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}
}
func toOrderRespItems(items []storer.OrderItem) []OrderItem {
	var respItems []OrderItem
	for _, item := range items {
		respItems = append(respItems, OrderItem{
			Name:      item.Name,
			Quantity:  item.Quantity,
			Image:     item.Image,
			Price:     item.Price,
			ProductID: item.ProductID,
		})
	}
	return respItems
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var u UserRequest
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// hash password
	hashed, err := util.HashPassword(u.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	u.Password = hashed

	user, err := h.server.CreateUser(h.ctx, toStorerUser(u))
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	res := toUserResponse(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
func toStorerUser(u UserRequest) *storer.User {
	return &storer.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}
func toUserResponse(u *storer.User) UserResponse {
	return UserResponse{
		Name:    u.Name,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}
func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.server.ListUsers(h.ctx)
	if err != nil {
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}
	var res ListUserResponse
	for _, u := range users {
		res.Users = append(res.Users, toUserResponse(&u))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	var u UserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	claims := r.Context().Value(authKey{}).(*token.UserClaims)
	user, err := h.server.GetUser(h.ctx, claims.Email)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	pathcUserReq(user, u)
	if user.Email == ""{
		user.Email = claims.Email
	}
	updated, err := h.server.UpdateUser(h.ctx, user)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	res := toUserResponse(updated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}
func pathcUserReq(user *storer.User, u UserRequest) {
	if u.Name != "" {
		user.Name = u.Name
	}
	if u.Email != "" {
		user.Email = u.Email
	}
	if u.Password != "" {
		hashed, err := util.HashPassword(u.Password)
		if err != nil {
			panic("Failed to hash password: " + err.Error())
		}
		user.Password = hashed
	}
	if u.IsAdmin {
		user.IsAdmin = u.IsAdmin
	}
	user.UpdatedAt = toTimePtr(time.Now())
}
func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Error parsing ID", http.StatusBadRequest)
		return
	}
	err = h.server.DeleteUser(h.ctx, i)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var u LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	gu, err := h.server.GetUser(h.ctx, u.Email)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	err = util.CheckPasswordHash(u.Password, gu.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	// create a json web token (JWT)
	accessToken, accessTokenClaims, err := h.TokenMaker.CreateToken(gu.ID, gu.Email, gu.IsAdmin, time.Minute*15)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	refreshToken, refreshClaims, err := h.TokenMaker.CreateToken(gu.ID, gu.Email, gu.IsAdmin, time.Hour*24)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}
	session, _ := h.server.CreateSession(h.ctx, &storer.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    gu.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})

	res := LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User:                  toUserResponse(gu),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
func (h *Handler) logoutUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)	
	err := h.server.DeleteSession(h.ctx, claims.RegisteredClaims.ID)
	if err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req RenewAccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	refreshClaims, err := h.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Error verifying token", http.StatusUnauthorized)
		return
	}
	session, err := h.server.GetSession(h.ctx, refreshClaims.RegisteredClaims.ID)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	if session.IsRevoked {
		http.Error(w, "Session is revoked", http.StatusUnauthorized)
		return
	}
	if session.UserEmail != refreshClaims.Email {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}
	accessToken, accessClaims, err := h.TokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.IsAdmin, time.Minute*15)
	if err != nil {
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}
	res := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
func (h *Handler) revokeSession(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)
	err := h.server.RevokeSession(h.ctx, claims.RegisteredClaims.ID)
	if err != nil {
		http.Error(w, "Failed to revoke session", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
