package repository

import (
	"context"
	"database/sql"
	"strings"

	dbent "github.com/x-frame/x-frame/ent"
	dbuser "github.com/x-frame/x-frame/ent/user"
	"github.com/x-frame/x-frame/internal/pkg/pagination"
	"github.com/x-frame/x-frame/internal/service"
)

type userRepository struct {
	client *dbent.Client
}

func NewUserRepository(client *dbent.Client, _ *sql.DB) service.UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) Create(ctx context.Context, userIn *service.User) error {
	if userIn == nil {
		return nil
	}
	created, err := r.client.User.Create().
		SetEmail(userIn.Email).
		SetUsername(userIn.Username).
		SetNotes(userIn.Notes).
		SetPasswordHash(userIn.PasswordHash).
		SetRole(userIn.Role).
		SetBalance(userIn.Balance).
		SetConcurrency(userIn.Concurrency).
		SetStatus(userIn.Status).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, service.ErrEmailExists)
	}
	applyUserEntityToService(userIn, created)
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*service.User, error) {
	m, err := r.client.User.Query().Where(dbuser.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrUserNotFound, nil)
	}
	return userEntityToService(m), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*service.User, error) {
	m, err := r.client.User.Query().Where(dbuser.EmailEQ(email)).Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrUserNotFound, nil)
	}
	return userEntityToService(m), nil
}

func (r *userRepository) GetFirstAdmin(ctx context.Context) (*service.User, error) {
	m, err := r.client.User.Query().
		Where(dbuser.RoleEQ(service.RoleAdmin)).
		Order(dbent.Asc(dbuser.FieldID)).
		First(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrUserNotFound, nil)
	}
	return userEntityToService(m), nil
}

func (r *userRepository) Update(ctx context.Context, userIn *service.User) error {
	if userIn == nil {
		return nil
	}
	upd := r.client.User.UpdateOneID(userIn.ID).
		SetEmail(userIn.Email).
		SetUsername(userIn.Username).
		SetNotes(userIn.Notes).
		SetPasswordHash(userIn.PasswordHash).
		SetRole(userIn.Role).
		SetBalance(userIn.Balance).
		SetConcurrency(userIn.Concurrency).
		SetStatus(userIn.Status)
	if userIn.TotpSecretEncrypted != nil {
		upd.SetTotpSecretEncrypted(*userIn.TotpSecretEncrypted)
	} else {
		upd.ClearTotpSecretEncrypted()
	}
	upd.SetTotpEnabled(userIn.TotpEnabled)
	if userIn.TotpEnabledAt != nil {
		upd.SetTotpEnabledAt(*userIn.TotpEnabledAt)
	} else {
		upd.ClearTotpEnabledAt()
	}
	updated, err := upd.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, service.ErrEmailExists)
	}
	applyUserEntityToService(userIn, updated)
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	return translatePersistenceError(err, service.ErrUserNotFound, nil)
}

func (r *userRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.User, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, service.UserListFilters{})
}

func (r *userRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters service.UserListFilters) ([]service.User, *pagination.PaginationResult, error) {
	q := r.client.User.Query()
	if filters.Status != "" {
		q = q.Where(dbuser.StatusEQ(filters.Status))
	}
	if filters.Role != "" {
		q = q.Where(dbuser.RoleEQ(filters.Role))
	}
	if search := strings.TrimSpace(filters.Search); search != "" {
		q = q.Where(dbuser.Or(
			dbuser.EmailContainsFold(search),
			dbuser.UsernameContainsFold(search),
		))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	page, pageSize := params.Page, params.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	rows, err := q.Order(dbent.Desc(dbuser.FieldID)).Offset(offset).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, nil, err
	}
	out := make([]service.User, 0, len(rows))
	for i := range rows {
		out = append(out, *userEntityToService(rows[i]))
	}
	return out, &pagination.PaginationResult{Total: int64(total), Page: page, PageSize: pageSize}, nil
}

func (r *userRepository) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	return r.client.User.UpdateOneID(id).SetBalance(amount).Exec(ctx)
}

func (r *userRepository) DeductBalance(ctx context.Context, id int64, amount float64) error {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	newBalance := user.Balance - amount
	if newBalance < 0 {
		newBalance = 0
	}
	return r.UpdateBalance(ctx, id, newBalance)
}

func (r *userRepository) UpdateConcurrency(ctx context.Context, id int64, amount int) error {
	return r.client.User.UpdateOneID(id).SetConcurrency(amount).Exec(ctx)
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.client.User.Query().Where(dbuser.EmailEQ(email)).Exist(ctx)
}

func (r *userRepository) UpdateTotpSecret(ctx context.Context, userID int64, encryptedSecret *string) error {
	upd := r.client.User.UpdateOneID(userID)
	if encryptedSecret == nil {
		upd.ClearTotpSecretEncrypted()
	} else {
		upd.SetTotpSecretEncrypted(*encryptedSecret)
	}
	return upd.Exec(ctx)
}

func (r *userRepository) EnableTotp(ctx context.Context, userID int64) error {
	return r.client.User.UpdateOneID(userID).SetTotpEnabled(true).Exec(ctx)
}

func (r *userRepository) DisableTotp(ctx context.Context, userID int64) error {
	return r.client.User.UpdateOneID(userID).
		SetTotpEnabled(false).
		ClearTotpSecretEncrypted().
		ClearTotpEnabledAt().
		Exec(ctx)
}

func userEntityToService(m *dbent.User) *service.User {
	if m == nil {
		return nil
	}
	u := &service.User{}
	applyUserEntityToService(u, m)
	return u
}

func applyUserEntityToService(dst *service.User, m *dbent.User) {
	if dst == nil || m == nil {
		return
	}
	dst.ID = m.ID
	dst.Email = m.Email
	dst.Username = m.Username
	dst.Notes = m.Notes
	dst.PasswordHash = m.PasswordHash
	dst.Role = m.Role
	dst.Balance = m.Balance
	dst.Concurrency = m.Concurrency
	dst.Status = m.Status
	dst.CreatedAt = m.CreatedAt
	dst.UpdatedAt = m.UpdatedAt
	dst.TotpSecretEncrypted = m.TotpSecretEncrypted
	dst.TotpEnabled = m.TotpEnabled
	dst.TotpEnabledAt = m.TotpEnabledAt
}
