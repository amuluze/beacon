// Package service
// Date       : 2024/9/4 14:56
// Author     : Amu
// Description:
package service

import (
	"amprobe/service/account/repository"
	"amprobe/service/schema"
	"context"
	"log/slog"

	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
)

var AccountServiceSet = wire.NewSet(NewAccountService, wire.Bind(new(IAccountService), new(*AccountService)))

type IAccountService interface {
	UserQuery(context.Context, schema.UserQueryArgs) (schema.UserQueryReply, error)
	UserCreate(context.Context, schema.UserCreateArgs) error
	UserUpdate(context.Context, schema.UserUpdateArgs) error
	UserDelete(context.Context, schema.UserDeleteArgs) error

	RoleQuery(context.Context, schema.RoleQueryArgs) (schema.RoleQueryReply, error)
	RoleCreate(context.Context, schema.RoleCreateArgs) error
	RoleUpdate(context.Context, schema.RoleUpdateArgs) error
	RoleDelete(context.Context, schema.RoleDeleteArgs) error

	ResourceQuery(context.Context, schema.ResourceQueryArgs) (schema.ResourceQueryReply, error)
}

type AccountService struct {
	Enforcer          *casbin.SyncedEnforcer
	AccountRepository repository.IAccountRepository
}

func NewAccountService(accountRepository repository.IAccountRepository, enforcer *casbin.SyncedEnforcer) *AccountService {
	return &AccountService{AccountRepository: accountRepository, Enforcer: enforcer}
}

func (a *AccountService) UserQuery(ctx context.Context, args schema.UserQueryArgs) (schema.UserQueryReply, error) {
	result := schema.UserQueryReply{Total: 0, Data: make([]schema.User, 0)}
	count, err := a.AccountRepository.UserCount(ctx)
	if err != nil {
		return result, err
	}
	result.Total = int(count)
	users, err := a.AccountRepository.UserQuery(ctx, args)
	if err != nil {
		return result, err
	}
	for _, u := range users {
		schemaRoles := make([]schema.Role, 0)
		for _, r := range u.Roles {
			schemaRoles = append(schemaRoles, schema.Role{
				ID:     r.ID.String(),
				Name:   r.Name,
				Status: r.Status,
				Remark: r.Remark,
			})
		}
		result.Data = append(result.Data, schema.User{
			ID:        u.ID.String(),
			Username:  u.Username,
			Remark:    u.Remark,
			Status:    u.Status,
			Roles:     schemaRoles,
			IsAdmin:   u.IsAdmin,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

func (a *AccountService) UserCreate(ctx context.Context, args schema.UserCreateArgs) error {
	user, err := a.AccountRepository.UserCreate(ctx, args)
	if err != nil {
		return err
	}
	for _, role := range user.Roles {
		_, err := a.Enforcer.AddNamedGroupingPolicy("g", user.ID.String(), role.ID.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccountService) UserUpdate(ctx context.Context, args schema.UserUpdateArgs) error {
	user, err := a.AccountRepository.UserUpdate(ctx, args)
	if err != nil {
		slog.Error("Error updating user", "error", err)
		return err
	}
	slog.Info("User updated", "User", user)
	for _, roleID := range args.RoleIDs {
		_, err := a.Enforcer.RemoveNamedGroupingPolicy("g", user.ID.String(), roleID)
		if err != nil {
			return err
		}
	}
	for _, role := range user.Roles {
		_, err := a.Enforcer.AddNamedGroupingPolicy("g", user.ID.String(), role.ID.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccountService) UserDelete(ctx context.Context, args schema.UserDeleteArgs) error {
	queryArgs := schema.UserQueryArgs{
		IDs: args.IDs,
	}
	users, err := a.AccountRepository.UserQuery(ctx, queryArgs)
	if err != nil {
		return err
	}

	err = a.AccountRepository.UserDelete(ctx, args)
	if err != nil {
		return err
	}

	for _, user := range users {
		for _, role := range user.Roles {
			_, err := a.Enforcer.RemoveNamedGroupingPolicy("g", user.ID.String(), role.ID.String())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AccountService) RoleQuery(ctx context.Context, args schema.RoleQueryArgs) (schema.RoleQueryReply, error) {
	result := schema.RoleQueryReply{Total: 0, Data: make([]schema.Role, 0)}
	count, err := a.AccountRepository.RoleCount(ctx)
	if err != nil {
		return result, err
	}
	result.Total = int(count)
	roles, err := a.AccountRepository.RoleQuery(ctx, args)
	if err != nil {
		return result, err
	}
	for _, role := range roles {
		schemaResources := make([]schema.Resource, 0)
		for _, resource := range role.Resources {
			if resource.Name == "" {
				continue
			}
			schemaResources = append(schemaResources, schema.Resource{
				ID:     resource.ID.String(),
				Name:   resource.Name,
				Status: resource.Status,
				Path:   resource.Path,
				Method: resource.Method,
			})
		}
		result.Data = append(result.Data, schema.Role{
			ID:        role.ID.String(),
			Name:      role.Name,
			Status:    role.Status,
			Remark:    role.Remark,
			CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
			Resources: schemaResources,
		})
	}
	return result, nil
}

func (a *AccountService) RoleCreate(ctx context.Context, args schema.RoleCreateArgs) error {
	role, err := a.AccountRepository.RoleCreate(ctx, args)
	if err != nil {
		return err
	}
	for _, resource := range role.Resources {
		_, err := a.Enforcer.AddNamedPolicy("p", role.ID.String(), resource.Path, resource.Method)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccountService) RoleUpdate(ctx context.Context, args schema.RoleUpdateArgs) error {
	role, err := a.AccountRepository.RoleUpdate(ctx, args)
	if err != nil {
		return err
	}
	for _, resourceID := range args.ResourceIDs {
		resources, err := a.AccountRepository.ResourceQuery(ctx, schema.ResourceQueryArgs{ID: resourceID})
		if err != nil {
			return err
		}
		_, err = a.Enforcer.RemoveNamedPolicy("p", role.ID.String(), resources[0].Path, resources[0].Method)
		if err != nil {
			return err
		}
	}
	for _, resource := range role.Resources {
		_, err := a.Enforcer.AddNamedPolicy("p", role.ID.String(), resource.Path, resource.Method)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccountService) RoleDelete(ctx context.Context, args schema.RoleDeleteArgs) error {
	queryArgs := schema.RoleQueryArgs{
		IDs: args.IDs,
	}
	roles, err := a.AccountRepository.RoleQuery(ctx, queryArgs)
	if err != nil {
		return err
	}
	err = a.AccountRepository.RoleDelete(ctx, args)
	if err != nil {
		return err
	}
	for _, role := range roles {
		for _, resource := range role.Resources {
			_, err := a.Enforcer.RemoveNamedPolicy("p", role.ID.String(), resource.Path, resource.Method)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AccountService) ResourceQuery(ctx context.Context, args schema.ResourceQueryArgs) (schema.ResourceQueryReply, error) {
	result := schema.ResourceQueryReply{Data: make([]schema.Resource, 0), Total: 0}
	count, err := a.AccountRepository.ResourceCount(ctx)
	if err != nil {
		return result, err
	}
	result.Total = int(count)
	resources, err := a.AccountRepository.ResourceQuery(ctx, args)
	if err != nil {
		return result, err
	}
	for _, resource := range resources {
		if resource.Name == "" {
			continue
		}

		result.Data = append(result.Data, schema.Resource{
			ID:     resource.ID.String(),
			Name:   resource.Name,
			Status: resource.Status,
			Path:   resource.Path,
			Method: resource.Method,
		})
	}
	return result, nil
}
