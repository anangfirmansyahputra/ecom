package types

import (
	"context"
	"time"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user User) error
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type ProductStore interface {
	GetProducts() ([]*Product, error)
	CreateProduct(product Product) error
	GetProductByID(id int) (*Product, error)
	UpdateProduct(id int, product Product) error
	DeleteProduct(id int) error
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1,max=1000"`
}

type OrderStore interface {
	CreateOrder(ctx context.Context, payload CreateOrderPayload) error
}

type OrderItem struct {
	ProductID int `json:"productId" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
}

type CreateOrderPayload struct {
	Items []OrderItem `json:"items" validate:"required"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
