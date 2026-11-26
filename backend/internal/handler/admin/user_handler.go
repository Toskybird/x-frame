package admin

import (
	"strconv"
	"strings"

	"github.com/x-frame/x-frame/internal/handler/dto"
	"github.com/x-frame/x-frame/internal/pkg/response"
	"github.com/x-frame/x-frame/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler handles admin user management
type UserHandler struct {
	adminService service.AdminService
}

// NewUserHandler creates a new admin user handler
func NewUserHandler(adminService service.AdminService) *UserHandler {
	return &UserHandler{adminService: adminService}
}

type CreateUserRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	Username    string  `json:"username"`
	Notes       string  `json:"notes"`
	Balance     float64 `json:"balance"`
	Concurrency int     `json:"concurrency"`
}

type UpdateUserRequest struct {
	Email       string   `json:"email" binding:"omitempty,email"`
	Password    string   `json:"password" binding:"omitempty,min=6"`
	Username    *string  `json:"username"`
	Notes       *string  `json:"notes"`
	Balance     *float64 `json:"balance"`
	Concurrency *int     `json:"concurrency"`
	Status      string   `json:"status" binding:"omitempty,oneof=active disabled"`
}

func (h *UserHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	search := strings.TrimSpace(c.Query("search"))
	if runes := []rune(search); len(runes) > 100 {
		search = string(runes[:100])
	}
	filters := service.UserListFilters{
		Status:     c.Query("status"),
		Role:       c.Query("role"),
		Search:     search,
		Attributes: parseAttributeFilters(c),
	}
	users, total, err := h.adminService.ListUsers(c.Request.Context(), page, pageSize, filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]*dto.AdminUser, 0, len(users))
	for i := range users {
		out = append(out, dto.UserFromServiceAdmin(&users[i]))
	}
	response.Paginated(c, out, total, page, pageSize)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}
	user, err := h.adminService.GetUser(c.Request.Context(), userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.UserFromServiceAdmin(user))
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, err := h.adminService.CreateUser(c.Request.Context(), &service.CreateUserInput{
		Email:       req.Email,
		Password:    req.Password,
		Username:    req.Username,
		Notes:       req.Notes,
		Balance:     req.Balance,
		Concurrency: req.Concurrency,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.UserFromServiceAdmin(user))
}

func (h *UserHandler) Update(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, err := h.adminService.UpdateUser(c.Request.Context(), userID, &service.UpdateUserInput{
		Email:       req.Email,
		Password:    req.Password,
		Username:    req.Username,
		Notes:       req.Notes,
		Balance:     req.Balance,
		Concurrency: req.Concurrency,
		Status:      req.Status,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.UserFromServiceAdmin(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}
	if err := h.adminService.DeleteUser(c.Request.Context(), userID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"ok": true})
}
