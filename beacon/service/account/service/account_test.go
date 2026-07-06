package service

import (
	"context"
	"testing"
	"time"

	"beacon/service/account/repository"
	"beacon/service/model"
	"beacon/service/schema"
	testutil "beacon/service/testutil"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestEnforcer creates an in-memory Casbin enforcer with a basic RBAC model
// suitable for unit tests.
func newTestEnforcer(t *testing.T) *casbin.SyncedEnforcer {
	t.Helper()
	m, err := casbinModel.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
	require.NoError(t, err, "failed to create casbin model")
	e, err := casbin.NewSyncedEnforcer(m)
	require.NoError(t, err, "failed to create synced enforcer")
	return e
}

// Verify interface compliance at compile time.
var _ repository.IAccountRepository = (*testutil.FakeAccountRepo)(nil)

func TestAccountService_UserQuery(t *testing.T) {
	ctx := context.Background()
	enforcer := newTestEnforcer(t)

	roleID1 := uuid.Must(uuid.NewV7())
	userID1 := uuid.Must(uuid.NewV7())
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeAccountRepo)
		args    schema.UserQueryArgs
		want    schema.UserQueryReply
		wantErr bool
	}{
		{
			name: "success_with_user_role_conversion",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.UserCountFn = func(_ context.Context) (int64, error) {
					return 1, nil
				}
				repo.UserQueryFn = func(_ context.Context, _ schema.UserQueryArgs) (model.Users, error) {
					return model.Users{
						&model.User{
							ID:        userID1,
							Username:  "admin",
							Remark:    "super admin",
							IsAdmin:   1,
							Status:    1,
							CreatedAt: now,
							Roles: model.Roles{
								&model.Role{
									ID:     roleID1,
									Name:   "admin",
									Status: 1,
									Remark: "administrator role",
								},
							},
						},
					}, nil
				}
			},
			args: schema.UserQueryArgs{},
			want: schema.UserQueryReply{
				Total: 1,
				Data: []schema.User{
					{
						ID:        userID1.String(),
						Username:  "admin",
						Remark:    "super admin",
						IsAdmin:   1,
						Status:    1,
						CreatedAt: now.Format("2006-01-02 15:04:05"),
						Roles: []schema.Role{
							{
								ID:     roleID1.String(),
								Name:   "admin",
								Status: 1,
								Remark: "administrator role",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "user_count_error",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.UserCountFn = func(_ context.Context) (int64, error) {
					return 0, testutil.ErrTest
				}
			},
			args:    schema.UserQueryArgs{},
			wantErr: true,
		},
		{
			name: "user_query_error",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.UserCountFn = func(_ context.Context) (int64, error) {
					return 0, nil
				}
				repo.UserQueryFn = func(_ context.Context, _ schema.UserQueryArgs) (model.Users, error) {
					return nil, testutil.ErrTest
				}
			},
			args:    schema.UserQueryArgs{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testutil.NewFakeAccountRepo()
			tt.setup(repo)
			svc := NewAccountService(repo, enforcer)

			got, err := svc.UserQuery(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.Total, got.Total)
			require.Len(t, got.Data, len(tt.want.Data))
			for i := range tt.want.Data {
				assert.Equal(t, tt.want.Data[i].ID, got.Data[i].ID)
				assert.Equal(t, tt.want.Data[i].Username, got.Data[i].Username)
				assert.Equal(t, tt.want.Data[i].Remark, got.Data[i].Remark)
				assert.Equal(t, tt.want.Data[i].IsAdmin, got.Data[i].IsAdmin)
				assert.Equal(t, tt.want.Data[i].Status, got.Data[i].Status)
				assert.Equal(t, tt.want.Data[i].CreatedAt, got.Data[i].CreatedAt)
				require.Len(t, got.Data[i].Roles, len(tt.want.Data[i].Roles))
				for j := range tt.want.Data[i].Roles {
					assert.Equal(t, tt.want.Data[i].Roles[j].ID, got.Data[i].Roles[j].ID)
					assert.Equal(t, tt.want.Data[i].Roles[j].Name, got.Data[i].Roles[j].Name)
				}
			}
		})
	}
}

func TestAccountService_UserCreate(t *testing.T) {
	ctx := context.Background()
	enforcer := newTestEnforcer(t)

	userID := uuid.Must(uuid.NewV7())
	roleID1 := uuid.Must(uuid.NewV7())
	roleID2 := uuid.Must(uuid.NewV7())

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeAccountRepo)
		args    schema.UserCreateArgs
		wantErr bool
	}{
		{
			name: "success_with_casbin_policy_added",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.UserCreateFn = func(_ context.Context, _ schema.UserCreateArgs) (model.User, error) {
					return model.User{
						ID:       userID,
						Username: "newuser",
						Roles: model.Roles{
							&model.Role{ID: roleID1},
							&model.Role{ID: roleID2},
						},
					}, nil
				}
			},
			args: schema.UserCreateArgs{
				Username: "newuser",
				Password: "secret",
				RoleIDs:  []string{roleID1.String(), roleID2.String()},
			},
			wantErr: false,
		},
		{
			name: "repo_error",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.UserCreateFn = func(_ context.Context, _ schema.UserCreateArgs) (model.User, error) {
					return model.User{}, testutil.ErrTest
				}
			},
			args: schema.UserCreateArgs{
				Username: "newuser",
				Password: "secret",
				RoleIDs:  []string{roleID1.String()},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testutil.NewFakeAccountRepo()
			tt.setup(repo)
			svc := NewAccountService(repo, enforcer)

			err := svc.UserCreate(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)

			// Verify Casbin grouping policies were added
			for _, roleID := range []uuid.UUID{roleID1, roleID2} {
				has, err := enforcer.HasGroupingPolicy(userID.String(), roleID.String())
				require.NoError(t, err)
				assert.True(t, has, "expected grouping policy for user %s and role %s", userID, roleID)
			}
		})
	}
}

func TestAccountService_UserUpdate(t *testing.T) {
	ctx := context.Background()
	enforcer := newTestEnforcer(t)

	userID := uuid.Must(uuid.NewV7())
	roleID1 := uuid.Must(uuid.NewV7())
	roleID2 := uuid.Must(uuid.NewV7())

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeAccountRepo, e *casbin.SyncedEnforcer)
		args    schema.UserUpdateArgs
		wantErr bool
	}{
		{
			name: "success_with_casbin_policy_updated",
			setup: func(repo *testutil.FakeAccountRepo, e *casbin.SyncedEnforcer) {
				// Pre-add old grouping policies
				_, _ = e.AddGroupingPolicy(userID.String(), roleID1.String())

				repo.UserUpdateFn = func(_ context.Context, _ schema.UserUpdateArgs) (model.User, error) {
					return model.User{
						ID:       userID,
						Username: "updated",
						Roles: model.Roles{
							&model.Role{ID: roleID2},
						},
					}, nil
				}
			},
			args: schema.UserUpdateArgs{
				ID:       userID.String(),
				Username: "updated",
				RoleIDs:  []string{roleID1.String()},
			},
			wantErr: false,
		},
		{
			name: "repo_error",
			setup: func(repo *testutil.FakeAccountRepo, _ *casbin.SyncedEnforcer) {
				repo.UserUpdateFn = func(_ context.Context, _ schema.UserUpdateArgs) (model.User, error) {
					return model.User{}, testutil.ErrTest
				}
			},
			args: schema.UserUpdateArgs{
				ID:       userID.String(),
				Username: "updated",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testutil.NewFakeAccountRepo()
			tt.setup(repo, enforcer)
			svc := NewAccountService(repo, enforcer)

			err := svc.UserUpdate(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)

			// Old grouping policy should be removed
			has, err := enforcer.HasGroupingPolicy(userID.String(), roleID1.String())
			require.NoError(t, err)
			assert.False(t, has, "old grouping policy should be removed")

			// New grouping policy should be added
			has, err = enforcer.HasGroupingPolicy(userID.String(), roleID2.String())
			require.NoError(t, err)
			assert.True(t, has, "new grouping policy should be added")
		})
	}
}

func TestAccountService_UserDelete(t *testing.T) {
	ctx := context.Background()
	enforcer := newTestEnforcer(t)

	userID := uuid.Must(uuid.NewV7())
	roleID1 := uuid.Must(uuid.NewV7())
	roleID2 := uuid.Must(uuid.NewV7())

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeAccountRepo, e *casbin.SyncedEnforcer)
		args    schema.UserDeleteArgs
		wantErr bool
	}{
		{
			name: "success_with_casbin_policy_removed",
			setup: func(repo *testutil.FakeAccountRepo, e *casbin.SyncedEnforcer) {
				// Pre-add grouping policies
				_, _ = e.AddGroupingPolicy(userID.String(), roleID1.String())
				_, _ = e.AddGroupingPolicy(userID.String(), roleID2.String())

				repo.UserQueryFn = func(_ context.Context, _ schema.UserQueryArgs) (model.Users, error) {
					return model.Users{
						&model.User{
							ID: userID,
							Roles: model.Roles{
								&model.Role{ID: roleID1},
								&model.Role{ID: roleID2},
							},
						},
					}, nil
				}
				repo.UserDeleteFn = func(_ context.Context, _ schema.UserDeleteArgs) error {
					return nil
				}
			},
			args: schema.UserDeleteArgs{
				IDs: []string{userID.String()},
			},
			wantErr: false,
		},
		{
			name: "user_query_error",
			setup: func(repo *testutil.FakeAccountRepo, _ *casbin.SyncedEnforcer) {
				repo.UserQueryFn = func(_ context.Context, _ schema.UserQueryArgs) (model.Users, error) {
					return nil, testutil.ErrTest
				}
			},
			args:    schema.UserDeleteArgs{IDs: []string{userID.String()}},
			wantErr: true,
		},
		{
			name: "user_delete_error",
			setup: func(repo *testutil.FakeAccountRepo, _ *casbin.SyncedEnforcer) {
				repo.UserQueryFn = func(_ context.Context, _ schema.UserQueryArgs) (model.Users, error) {
					return model.Users{}, nil
				}
				repo.UserDeleteFn = func(_ context.Context, _ schema.UserDeleteArgs) error {
					return testutil.ErrTest
				}
			},
			args:    schema.UserDeleteArgs{IDs: []string{userID.String()}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testutil.NewFakeAccountRepo()
			tt.setup(repo, enforcer)
			svc := NewAccountService(repo, enforcer)

			err := svc.UserDelete(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)

			// Grouping policies should be removed
			has, err := enforcer.HasGroupingPolicy(userID.String(), roleID1.String())
			require.NoError(t, err)
			assert.False(t, has, "grouping policy for role1 should be removed")

			has, err = enforcer.HasGroupingPolicy(userID.String(), roleID2.String())
			require.NoError(t, err)
			assert.False(t, has, "grouping policy for role2 should be removed")
		})
	}
}

func TestAccountService_RoleQuery(t *testing.T) {
	ctx := context.Background()
	enforcer := newTestEnforcer(t)

	roleID := uuid.Must(uuid.NewV7())
	resourceID := uuid.Must(uuid.NewV7())
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeAccountRepo)
		args    schema.RoleQueryArgs
		want    schema.RoleQueryReply
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.RoleCountFn = func(_ context.Context) (int64, error) {
					return 1, nil
				}
				repo.RoleQueryFn = func(_ context.Context, _ schema.RoleQueryArgs) (model.Roles, error) {
					return model.Roles{
						&model.Role{
							ID:        roleID,
							Name:      "admin",
							Status:    1,
							Remark:    "administrator",
							CreatedAt: now,
							Resources: model.Resources{
								&model.Resource{
									ID:     resourceID,
									Name:   "users",
									Path:   "/api/users",
									Method: "GET",
									Status: 1,
								},
							},
						},
					}, nil
				}
			},
			args: schema.RoleQueryArgs{},
			want: schema.RoleQueryReply{
				Total: 1,
				Data: []schema.Role{
					{
						ID:        roleID.String(),
						Name:      "admin",
						Status:    1,
						Remark:    "administrator",
						CreatedAt: now.Format("2006-01-02 15:04:05"),
						Resources: []schema.Resource{
							{
								ID:     resourceID.String(),
								Name:   "users",
								Path:   "/api/users",
								Method: "GET",
								Status: 1,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "role_count_error",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.RoleCountFn = func(_ context.Context) (int64, error) {
					return 0, testutil.ErrTest
				}
			},
			args:    schema.RoleQueryArgs{},
			wantErr: true,
		},
		{
			name: "role_query_error",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.RoleCountFn = func(_ context.Context) (int64, error) {
					return 0, nil
				}
				repo.RoleQueryFn = func(_ context.Context, _ schema.RoleQueryArgs) (model.Roles, error) {
					return nil, testutil.ErrTest
				}
			},
			args:    schema.RoleQueryArgs{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testutil.NewFakeAccountRepo()
			tt.setup(repo)
			svc := NewAccountService(repo, enforcer)

			got, err := svc.RoleQuery(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.Total, got.Total)
			require.Len(t, got.Data, len(tt.want.Data))
			for i := range tt.want.Data {
				assert.Equal(t, tt.want.Data[i].ID, got.Data[i].ID)
				assert.Equal(t, tt.want.Data[i].Name, got.Data[i].Name)
				assert.Equal(t, tt.want.Data[i].Remark, got.Data[i].Remark)
				assert.Equal(t, tt.want.Data[i].CreatedAt, got.Data[i].CreatedAt)
				require.Len(t, got.Data[i].Resources, len(tt.want.Data[i].Resources))
				for j := range tt.want.Data[i].Resources {
					assert.Equal(t, tt.want.Data[i].Resources[j].ID, got.Data[i].Resources[j].ID)
					assert.Equal(t, tt.want.Data[i].Resources[j].Name, got.Data[i].Resources[j].Name)
					assert.Equal(t, tt.want.Data[i].Resources[j].Path, got.Data[i].Resources[j].Path)
					assert.Equal(t, tt.want.Data[i].Resources[j].Method, got.Data[i].Resources[j].Method)
				}
			}
		})
	}
}

func TestAccountService_RoleCreate(t *testing.T) {
	ctx := context.Background()
	enforcer := newTestEnforcer(t)

	roleID := uuid.Must(uuid.NewV7())
	resourceID := uuid.Must(uuid.NewV7())

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeAccountRepo)
		args    schema.RoleCreateArgs
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.RoleCreateFn = func(_ context.Context, _ schema.RoleCreateArgs) (model.Role, error) {
					return model.Role{
						ID:   roleID,
						Name: "editor",
						Resources: model.Resources{
							&model.Resource{
								ID:     resourceID,
								Path:   "/api/posts",
								Method: "POST",
							},
						},
					}, nil
				}
			},
			args: schema.RoleCreateArgs{
				Name:        "editor",
				Status:      1,
				ResourceIDs: []string{resourceID.String()},
			},
			wantErr: false,
		},
		{
			name: "repo_error",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.RoleCreateFn = func(_ context.Context, _ schema.RoleCreateArgs) (model.Role, error) {
					return model.Role{}, testutil.ErrTest
				}
			},
			args: schema.RoleCreateArgs{
				Name:   "editor",
				Status: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testutil.NewFakeAccountRepo()
			tt.setup(repo)
			svc := NewAccountService(repo, enforcer)

			err := svc.RoleCreate(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)

			// Verify Casbin policy was added
			has, err := enforcer.HasPolicy(roleID.String(), "/api/posts", "POST")
			require.NoError(t, err)
			assert.True(t, has, "expected policy for role %s", roleID)
		})
	}
}

func TestAccountService_ResourceQuery(t *testing.T) {
	ctx := context.Background()
	enforcer := newTestEnforcer(t)

	resourceID := uuid.Must(uuid.NewV7())

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeAccountRepo)
		args    schema.ResourceQueryArgs
		want    schema.ResourceQueryReply
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.ResourceCountFn = func(_ context.Context) (int64, error) {
					return 1, nil
				}
				repo.ResourceQueryFn = func(_ context.Context, _ schema.ResourceQueryArgs) (model.Resources, error) {
					return model.Resources{
						&model.Resource{
							ID:     resourceID,
							Name:   "users",
							Path:   "/api/users",
							Method: "GET",
							Status: 1,
						},
					}, nil
				}
			},
			args: schema.ResourceQueryArgs{},
			want: schema.ResourceQueryReply{
				Total: 1,
				Data: []schema.Resource{
					{
						ID:     resourceID.String(),
						Name:   "users",
						Path:   "/api/users",
						Method: "GET",
						Status: 1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "resource_count_error",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.ResourceCountFn = func(_ context.Context) (int64, error) {
					return 0, testutil.ErrTest
				}
			},
			args:    schema.ResourceQueryArgs{},
			wantErr: true,
		},
		{
			name: "resource_query_error",
			setup: func(repo *testutil.FakeAccountRepo) {
				repo.ResourceCountFn = func(_ context.Context) (int64, error) {
					return 0, nil
				}
				repo.ResourceQueryFn = func(_ context.Context, _ schema.ResourceQueryArgs) (model.Resources, error) {
					return nil, testutil.ErrTest
				}
			},
			args:    schema.ResourceQueryArgs{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testutil.NewFakeAccountRepo()
			tt.setup(repo)
			svc := NewAccountService(repo, enforcer)

			got, err := svc.ResourceQuery(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.Total, got.Total)
			require.Len(t, got.Data, len(tt.want.Data))
			for i := range tt.want.Data {
				assert.Equal(t, tt.want.Data[i].ID, got.Data[i].ID)
				assert.Equal(t, tt.want.Data[i].Name, got.Data[i].Name)
				assert.Equal(t, tt.want.Data[i].Path, got.Data[i].Path)
				assert.Equal(t, tt.want.Data[i].Method, got.Data[i].Method)
				assert.Equal(t, tt.want.Data[i].Status, got.Data[i].Status)
			}
		})
	}
}
