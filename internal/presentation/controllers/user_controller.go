package controllers

import (
	"net/http" // For http status codes

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus" // Assuming logrus from go.mod
	"github.com/yourusername/electricity-shop-go/internal/application/commands"
	"github.com/yourusername/electricity-shop-go/internal/application/dtos"
	"github.com/yourusername/electricity-shop-go/internal/presentation/responses" // For standard responses
	"github.com/yourusername/electricity-shop-go/pkg/mediator"
)

type UserController struct {
	mediator mediator.Mediator
	logger   *logrus.Logger
}

func NewUserController(m mediator.Mediator, l *logrus.Logger) *UserController {
	return &UserController{mediator: m, logger: l}
}

func (uc *UserController) Register(c *gin.Context) {
	var req dtos.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Errorf("Failed to bind request for user registration: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", "INVALID_INPUT"))
		return
	}

	// TODO: Add validation using go-playground/validator for req DTO

	cmd := commands.RegisterUserCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	// Assuming RegisterUserCommandHandler's Handle method returns `error`
	err := uc.mediator.Send(c.Request.Context(), &cmd)
	if err != nil {
		uc.logger.Errorf("User registration failed: %v", err)
		// TODO: Map domain errors to HTTP status codes more granularly
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(err.Error(), "REGISTRATION_FAILED"))
		return
	}

	c.JSON(http.StatusCreated, responses.NewSuccessResponse(nil, "User registered successfully"))
}

func (uc *UserController) Login(c *gin.Context) {
	var req dtos.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.logger.Errorf("Failed to bind request for user login: %v", err)
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Invalid request body", "INVALID_INPUT"))
		return
	}

	// TODO: Add validation for req DTO

	cmd := commands.LoginUserCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	// LoginUserCommandHandler's Handle method returns (interface{}, error)
	// The mediator's Send method is for commands that return error only.
	// The mediator's Query method is for queries that return (interface{}, error).
	// Login is a command that returns data. This is a common CQRS nuance.
	// LoginUserCommand already has GetName(), so it can satisfy mediator.Query interface.
	result, err := uc.mediator.Query(c.Request.Context(), &cmd) // cmd must implement Query
	if err != nil {
		uc.logger.Errorf("User login failed: %v", err)
		// TODO: Map domain errors (e.g., invalid credentials) to specific HTTP status codes
		c.JSON(http.StatusUnauthorized, responses.NewErrorResponse(err.Error(), "LOGIN_FAILED"))
		return
	}

	loginResponse, ok := result.(*dtos.LoginUserResponse)
	if !ok {
		uc.logger.Errorf("Login handler returned unexpected type: %T", result)
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse("Internal server error", "LOGIN_UNEXPECTED_TYPE"))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccessResponse(loginResponse, "Login successful"))
}
