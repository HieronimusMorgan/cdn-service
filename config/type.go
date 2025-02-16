package config

import (
	"cdn-service/internal/controller"
	"cdn-service/internal/middleware"
	"cdn-service/internal/services"
	"cdn-service/internal/utils"
	"github.com/gin-gonic/gin"
)

// ServerConfig holds all initialized components
type ServerConfig struct {
	Gin        *gin.Engine
	Config     *Config
	Redis      utils.RedisService
	JWTService utils.JWTService
	Controller Controller
	Services   Services
	Repository Repository
	Middleware Middleware
}

// Services holds all service dependencies
type Services struct {
	ImageService services.ImageService
	//AuthService        services.AuthService
	//UserSessionService services.UsersSessionService
	//ResourceService    services.ResourceService
	//RoleService        services.RoleService
}

// Repository contains repository (database access objects)
type Repository struct {
	//AuthRepo         repository.AuthRepository
	//UserRepo         repository.UserRepository
	//ResourceRepo     repository.ResourceRepository
	//RoleResourceRepo repository.RoleResourceRepository
	//RoleRepo         repository.RoleRepository
	//UserRoleRepo     repository.UserRoleRepository
	//UserSessionRepo  repository.UserSessionRepository
}

type Controller struct {
	ImageController controller.ImageController
	//AuthHandler     handler.AuthHandler
	//ResourceHandler handler.ResourceHandler
	//RoleHandler     handler.RoleHandler
}

type Middleware struct {
	AuthMiddleware middleware.AuthMiddleware
}
