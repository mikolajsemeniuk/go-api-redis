package order

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type server struct {
	router fiber.Router
	redis  *redis.Client
}

func (s *server) Route() {
	s.router.Get("/:id", s.get)
	s.router.Post("/", s.set)
}

// @Summary Get
// @Schemes
// @Description Get order
// @Tags order
// @Param id path string true "id"
// @Success 200 {object} Order
// @Success 400 {string} bad request
// @Success 404 {string} not found
// @Failure 503 {string} not found
// @Router /order/{id} [get]
func (s *server) get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	value, err := s.redis.Get(context.Background(), id.String()).Result()
	switch {
	case err == redis.Nil:
		return c.Status(http.StatusNotFound).SendString(err.Error())
	case err != nil:
		return c.Status(http.StatusServiceUnavailable).SendString(err.Error())
	default:
		return c.SendString(value)
	}
}

// @Summary Set
// @Schemes
// @Description Set order
// @Tags order
// @Param body body Order true "Body"
// @Success 200 {object} Order
// @Success 400 {string} string "error"
// @Failure 503 {string} string "error"
// @Router /order [post]
func (s *server) set(c *fiber.Ctx) error {
	var request Order
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	order, err := json.Marshal(request)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if err := s.redis.Set(context.Background(), request.Key.String(), order, 0).Err(); err != nil {
		return c.Status(http.StatusServiceUnavailable).SendString(err.Error())
	}

	return c.Send(order)
}

func NewServer(r fiber.Router, c *redis.Client) *server {
	return &server{
		router: r,
		redis:  c,
	}
}
