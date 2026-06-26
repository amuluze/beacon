package service

import (
	"context"
	"testing"
	"time"

	"amprobe/service/container/repository"
	"amprobe/service/schema"
	testutil "amprobe/service/testutil"
	rpcSchema "common/rpc/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newFakeContainerRepoWithDefaults returns a FakeContainerRepo where every Fn
// returns a zero-value result and nil error, so tests only need to override
// the specific Fn they care about.
func newFakeContainerRepoWithDefaults() *testutil.FakeContainerRepo {
	return &testutil.FakeContainerRepo{
		VersionFn: func(_ context.Context, _ rpcSchema.DockerArgs) (rpcSchema.DockerReply, error) {
			return rpcSchema.DockerReply{}, nil
		},
		ContainerListFn: func(_ context.Context, _ rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error) {
			return rpcSchema.ContainerQueryReply{}, nil
		},
		ContainerCountFn: func(_ context.Context, _ rpcSchema.ContainerCountArgs) (rpcSchema.ContainerCountReply, error) {
			return rpcSchema.ContainerCountReply{}, nil
		},
		UsageFn: func(_ context.Context, _ rpcSchema.ContainerUsageArgs) (rpcSchema.ContainerUsageReply, error) {
			return rpcSchema.ContainerUsageReply{}, nil
		},
		ContainersByImageFn: func(_ context.Context, _ string) (int, error) {
			return 0, nil
		},
		ContainerCreateFn: func(_ context.Context, _ rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error) {
			return rpcSchema.ContainerCreateReply{}, nil
		},
		ContainerUpdateFn: func(_ context.Context, _ rpcSchema.ContainerUpdateArgs) (rpcSchema.ContainerUpdateReply, error) {
			return rpcSchema.ContainerUpdateReply{}, nil
		},
		ContainerDeleteFn: func(_ context.Context, _ rpcSchema.ContainerDeleteArgs) error {
			return nil
		},
		ContainerStartFn: func(_ context.Context, _ rpcSchema.ContainerStartArgs) error {
			return nil
		},
		ContainerStopFn: func(_ context.Context, _ rpcSchema.ContainerStopArgs) error {
			return nil
		},
		ContainerRestartFn: func(_ context.Context, _ rpcSchema.ContainerRestartArgs) error {
			return nil
		},
		ContainerLogsFn: func(_ context.Context, _ rpcSchema.ContainerLogsArgs) (rpcSchema.ContainerLogsReply, error) {
			return rpcSchema.ContainerLogsReply{}, nil
		},
		ImageListFn: func(_ context.Context, _ rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error) {
			return rpcSchema.ImageQueryReply{}, nil
		},
		ImageCountFn: func(_ context.Context, _ rpcSchema.ImageCountArgs) (rpcSchema.ImageCountReply, error) {
			return rpcSchema.ImageCountReply{}, nil
		},
		ImagePullFn: func(_ context.Context, _ rpcSchema.ImagePullArgs) error {
			return nil
		},
		ImageTagFn: func(_ context.Context, _ rpcSchema.ImageTagArgs) error {
			return nil
		},
		ImageImportFn: func(_ context.Context, _ rpcSchema.ImageImportArgs) error {
			return nil
		},
		ImageExportFn: func(_ context.Context, _ rpcSchema.ImageExportArgs) (rpcSchema.ImageExportReply, error) {
			return rpcSchema.ImageExportReply{}, nil
		},
		ImageDeleteFn: func(_ context.Context, _ rpcSchema.ImageDeleteArgs) error {
			return nil
		},
		ImagesPruneFn: func(_ context.Context) error {
			return nil
		},
		NetworkListFn: func(_ context.Context, _ rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error) {
			return rpcSchema.NetworkQueryReply{}, nil
		},
		NetworkCountFn: func(_ context.Context, _ rpcSchema.NetworkCountArgs) (rpcSchema.NetworkCountReply, error) {
			return rpcSchema.NetworkCountReply{}, nil
		},
		NetworkCreateFn: func(_ context.Context, _ rpcSchema.NetworkCreateArgs) (rpcSchema.NetworkCreateReply, error) {
			return rpcSchema.NetworkCreateReply{}, nil
		},
		NetworkDeleteFn: func(_ context.Context, _ rpcSchema.NetworkDeleteArgs) error {
			return nil
		},
		GetDockerRegistryMirrorsFn: func(_ context.Context, _ rpcSchema.GetDockerRegistryMirrorsArgs) (rpcSchema.GetDockerRegistryMirrorsReply, error) {
			return rpcSchema.GetDockerRegistryMirrorsReply{}, nil
		},
		SetDockerRegistryMirrorsFn: func(_ context.Context, _ rpcSchema.SetDockerRegistryMirrorsArgs) error {
			return nil
		},
	}
}

// Verify interface compliance at compile time.
var _ repository.IContainerRepo = (*testutil.FakeContainerRepo)(nil)

func TestContainerService_Version(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeContainerRepo)
		want    schema.Docker
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.VersionFn = func(_ context.Context, _ rpcSchema.DockerArgs) (rpcSchema.DockerReply, error) {
					return rpcSchema.DockerReply{
						Data: rpcSchema.Docker{
							Timestamp:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							DockerVersion: "24.0.7",
							APIVersion:    "1.43",
							MinAPIVersion: "1.12",
							GitCommit:     "abc123",
							GoVersion:     "go1.21",
							Os:            "linux",
							Arch:          "amd64",
						},
					}, nil
				}
			},
			want: schema.Docker{
				Timestamp:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				DockerVersion: "24.0.7",
				APIVersion:    "1.43",
				MinAPIVersion: "1.12",
				GitCommit:     "abc123",
				GoVersion:     "go1.21",
				Os:            "linux",
				Arch:          "amd64",
			},
			wantErr: false,
		},
		{
			name: "repo_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.VersionFn = func(_ context.Context, _ rpcSchema.DockerArgs) (rpcSchema.DockerReply, error) {
					return rpcSchema.DockerReply{}, testutil.ErrTest
				}
			},
			want:    schema.Docker{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeContainerRepoWithDefaults()
			tt.setup(repo)
			svc := NewContainerService(repo)

			got, err := svc.Version(ctx)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContainerService_ContainerList(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeContainerRepo)
		args    schema.ContainerQueryArgs
		want    schema.ContainerQueryRely
		wantErr bool
	}{
		{
			name: "success_with_pagination",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.ContainerListFn = func(_ context.Context, _ rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error) {
					return rpcSchema.ContainerQueryReply{
						Data: []rpcSchema.Container{
							{
								ContainerID: "cid-1",
								Name:        "web",
								Image:       "nginx:latest",
								State:       "running",
								CPUPercent:  12.345,
								MemPercent:  45.678,
								MemUsage:    536870912,
								MemLimit:    1073741824,
							},
						},
					}, nil
				}
				repo.ContainerCountFn = func(_ context.Context, _ rpcSchema.ContainerCountArgs) (rpcSchema.ContainerCountReply, error) {
					return rpcSchema.ContainerCountReply{Count: 10}, nil
				}
			},
			args: schema.ContainerQueryArgs{Page: 1, Size: 10},
			want: schema.ContainerQueryRely{
				Data: []schema.Container{
					{
						ID:            "cid-1",
						Name:          "web",
						Image:         "nginx:latest",
						State:         "running",
						CPUPercent:    "12.35 %",
						MemoryPercent: "45.68 %",
						MemoryUsage:   "512.00 MB",
						MemoryLimit:   "1.00 GB",
					},
				},
				Total: 10,
				Page:  1,
				Size:  10,
			},
			wantErr: false,
		},
		{
			name: "repo_list_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.ContainerListFn = func(_ context.Context, _ rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error) {
					return rpcSchema.ContainerQueryReply{}, testutil.ErrTest
				}
			},
			args:    schema.ContainerQueryArgs{Page: 1, Size: 10},
			wantErr: true,
		},
		{
			name: "repo_count_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.ContainerListFn = func(_ context.Context, _ rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error) {
					return rpcSchema.ContainerQueryReply{Data: []rpcSchema.Container{}}, nil
				}
				repo.ContainerCountFn = func(_ context.Context, _ rpcSchema.ContainerCountArgs) (rpcSchema.ContainerCountReply, error) {
					return rpcSchema.ContainerCountReply{}, testutil.ErrTest
				}
			},
			args:    schema.ContainerQueryArgs{Page: 1, Size: 10},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeContainerRepoWithDefaults()
			tt.setup(repo)
			svc := NewContainerService(repo)

			got, err := svc.ContainerList(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.Total, got.Total)
			assert.Equal(t, tt.want.Page, got.Page)
			assert.Equal(t, tt.want.Size, got.Size)
			require.Len(t, got.Data, len(tt.want.Data))
			for i := range tt.want.Data {
				assert.Equal(t, tt.want.Data[i].ID, got.Data[i].ID)
				assert.Equal(t, tt.want.Data[i].Name, got.Data[i].Name)
				assert.Equal(t, tt.want.Data[i].Image, got.Data[i].Image)
				assert.Equal(t, tt.want.Data[i].State, got.Data[i].State)
				assert.Equal(t, tt.want.Data[i].CPUPercent, got.Data[i].CPUPercent)
				assert.Equal(t, tt.want.Data[i].MemoryPercent, got.Data[i].MemoryPercent)
				assert.Equal(t, tt.want.Data[i].MemoryUsage, got.Data[i].MemoryUsage)
				assert.Equal(t, tt.want.Data[i].MemoryLimit, got.Data[i].MemoryLimit)
			}
		})
	}
}

func TestContainerService_Usage(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeContainerRepo)
		args    schema.ContainerUsageArgs
		want    schema.ContainerUsageReply
		wantErr bool
	}{
		{
			name: "success_with_cpu_mem_map",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.UsageFn = func(_ context.Context, _ rpcSchema.ContainerUsageArgs) (rpcSchema.ContainerUsageReply, error) {
					return rpcSchema.ContainerUsageReply{
						Names: []string{"web", "db"},
						CPUUsage: map[string][]rpcSchema.Usage{
							"web": {{Timestamp: 1000, Value: 10.5}, {Timestamp: 2000, Value: 20.5}},
							"db":  {{Timestamp: 1000, Value: 30.5}},
						},
						MemUsage: map[string][]rpcSchema.Usage{
							"web": {{Timestamp: 1000, Value: 50.0}},
							"db":  {{Timestamp: 1000, Value: 60.0}, {Timestamp: 2000, Value: 70.0}},
						},
					}, nil
				}
			},
			args: schema.ContainerUsageArgs{StartTime: 0, EndTime: 3000},
			want: schema.ContainerUsageReply{
				Names: []string{"web", "db"},
				CPUUsage: map[string][]schema.Usage{
					"web": {{Timestamp: 1000, Value: 10.5}, {Timestamp: 2000, Value: 20.5}},
					"db":  {{Timestamp: 1000, Value: 30.5}},
				},
				MemUsage: map[string][]schema.Usage{
					"web": {{Timestamp: 1000, Value: 50.0}},
					"db":  {{Timestamp: 1000, Value: 60.0}, {Timestamp: 2000, Value: 70.0}},
				},
			},
			wantErr: false,
		},
		{
			name: "repo_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.UsageFn = func(_ context.Context, _ rpcSchema.ContainerUsageArgs) (rpcSchema.ContainerUsageReply, error) {
					return rpcSchema.ContainerUsageReply{}, testutil.ErrTest
				}
			},
			args:    schema.ContainerUsageArgs{StartTime: 0, EndTime: 3000},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeContainerRepoWithDefaults()
			tt.setup(repo)
			svc := NewContainerService(repo)

			got, err := svc.Usage(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.Names, got.Names)
			assert.Equal(t, tt.want.CPUUsage, got.CPUUsage)
			assert.Equal(t, tt.want.MemUsage, got.MemUsage)
		})
	}
}

func TestContainerService_ContainerCreate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeContainerRepo)
		args    schema.ContainerCreateArgs
		want    schema.ContainerCreateReply
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.ContainerCreateFn = func(_ context.Context, _ rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error) {
					return rpcSchema.ContainerCreateReply{ContainerID: "new-cid-123"}, nil
				}
			},
			args: schema.ContainerCreateArgs{
				ContainerName: "test-container",
				ImageName:     "nginx:latest",
				NetworkName:   "bridge",
				Ports:         []string{"80:80"},
			},
			want:    schema.ContainerCreateReply{ContainerID: "new-cid-123"},
			wantErr: false,
		},
		{
			name: "repo_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.ContainerCreateFn = func(_ context.Context, _ rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error) {
					return rpcSchema.ContainerCreateReply{}, testutil.ErrTest
				}
			},
			args:    schema.ContainerCreateArgs{ContainerName: "test", ImageName: "nginx"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeContainerRepoWithDefaults()
			tt.setup(repo)
			svc := NewContainerService(repo)

			got, err := svc.ContainerCreate(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContainerService_ImageList(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeContainerRepo)
		args    schema.ImageQueryArgs
		want    schema.ImageQueryReply
		wantErr bool
	}{
		{
			name: "success_with_containers_by_image",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.ImageListFn = func(_ context.Context, _ rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error) {
					return rpcSchema.ImageQueryReply{
						Data: []rpcSchema.Image{
							{ImageID: "img-1", Name: "nginx", Tag: "latest", Created: "2024-01-01", Size: "50MB"},
							{ImageID: "img-2", Name: "redis", Tag: "7", Created: "2024-01-02", Size: "30MB"},
						},
					}, nil
				}
				repo.ContainersByImageFn = func(_ context.Context, image string) (int, error) {
					if image == "nginx:latest" {
						return 3, nil
					}
					return 1, nil
				}
				repo.ImageCountFn = func(_ context.Context, _ rpcSchema.ImageCountArgs) (rpcSchema.ImageCountReply, error) {
					return rpcSchema.ImageCountReply{Count: 2}, nil
				}
			},
			args: schema.ImageQueryArgs{Page: 1, Size: 10},
			want: schema.ImageQueryReply{
				Data: []schema.Image{
					{ID: "img-1", Name: "nginx", Tag: "latest", Created: "2024-01-01", Size: "50MB", Number: 3},
					{ID: "img-2", Name: "redis", Tag: "7", Created: "2024-01-02", Size: "30MB", Number: 1},
				},
				Total: 2,
				Page:  1,
				Size:  10,
			},
			wantErr: false,
		},
		{
			name: "repo_image_list_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.ImageListFn = func(_ context.Context, _ rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error) {
					return rpcSchema.ImageQueryReply{}, testutil.ErrTest
				}
			},
			args:    schema.ImageQueryArgs{Page: 1, Size: 10},
			wantErr: true,
		},
		{
			name: "containers_by_image_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.ImageListFn = func(_ context.Context, _ rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error) {
					return rpcSchema.ImageQueryReply{
						Data: []rpcSchema.Image{
							{ImageID: "img-1", Name: "nginx", Tag: "latest"},
						},
					}, nil
				}
				repo.ContainersByImageFn = func(_ context.Context, _ string) (int, error) {
					return 0, testutil.ErrTest
				}
			},
			args:    schema.ImageQueryArgs{Page: 1, Size: 10},
			wantErr: true,
		},
		{
			name: "repo_image_count_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.ImageListFn = func(_ context.Context, _ rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error) {
					return rpcSchema.ImageQueryReply{Data: []rpcSchema.Image{}}, nil
				}
				repo.ImageCountFn = func(_ context.Context, _ rpcSchema.ImageCountArgs) (rpcSchema.ImageCountReply, error) {
					return rpcSchema.ImageCountReply{}, testutil.ErrTest
				}
			},
			args:    schema.ImageQueryArgs{Page: 1, Size: 10},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeContainerRepoWithDefaults()
			tt.setup(repo)
			svc := NewContainerService(repo)

			got, err := svc.ImageList(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.Total, got.Total)
			assert.Equal(t, tt.want.Page, got.Page)
			assert.Equal(t, tt.want.Size, got.Size)
			require.Len(t, got.Data, len(tt.want.Data))
			for i := range tt.want.Data {
				assert.Equal(t, tt.want.Data[i].ID, got.Data[i].ID)
				assert.Equal(t, tt.want.Data[i].Name, got.Data[i].Name)
				assert.Equal(t, tt.want.Data[i].Tag, got.Data[i].Tag)
				assert.Equal(t, tt.want.Data[i].Number, got.Data[i].Number)
			}
		})
	}
}

func TestContainerService_NetworkList(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(repo *testutil.FakeContainerRepo)
		args    schema.NetworkQueryArgs
		want    schema.NetworkQueryReply
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.NetworkListFn = func(_ context.Context, _ rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error) {
					return rpcSchema.NetworkQueryReply{
						Data: []rpcSchema.Network{
							{
								NetworkID: "net-1",
								Name:      "bridge",
								Driver:    "bridge",
								Created:   "2024-01-01",
								Subnet:    "172.17.0.0/16",
								Gateway:   "172.17.0.1",
								Labels:    `{"com.docker.label":"test"}`,
							},
						},
					}, nil
				}
				repo.NetworkCountFn = func(_ context.Context, _ rpcSchema.NetworkCountArgs) (rpcSchema.NetworkCountReply, error) {
					return rpcSchema.NetworkCountReply{Count: 5}, nil
				}
			},
			args: schema.NetworkQueryArgs{Page: 1, Size: 10},
			want: schema.NetworkQueryReply{
				Data: []schema.Network{
					{
						ID:      "net-1",
						Name:    "bridge",
						Driver:  "bridge",
						Created: "2024-01-01",
						Subnet:  "172.17.0.0/16",
						Gateway: "172.17.0.1",
						Labels:  map[string]string{"com.docker.label": "test"},
					},
				},
				Total: 5,
				Page:  1,
				Size:  10,
			},
			wantErr: false,
		},
		{
			name: "repo_list_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.NetworkListFn = func(_ context.Context, _ rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error) {
					return rpcSchema.NetworkQueryReply{}, testutil.ErrTest
				}
			},
			args:    schema.NetworkQueryArgs{Page: 1, Size: 10},
			wantErr: true,
		},
		{
			name: "repo_count_error",
			setup: func(repo *testutil.FakeContainerRepo) {
				repo.NetworkListFn = func(_ context.Context, _ rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error) {
					return rpcSchema.NetworkQueryReply{Data: []rpcSchema.Network{}}, nil
				}
				repo.NetworkCountFn = func(_ context.Context, _ rpcSchema.NetworkCountArgs) (rpcSchema.NetworkCountReply, error) {
					return rpcSchema.NetworkCountReply{}, testutil.ErrTest
				}
			},
			args:    schema.NetworkQueryArgs{Page: 1, Size: 10},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeContainerRepoWithDefaults()
			tt.setup(repo)
			svc := NewContainerService(repo)

			got, err := svc.NetworkList(ctx, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, testutil.ErrTest)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.Total, got.Total)
			assert.Equal(t, tt.want.Page, got.Page)
			assert.Equal(t, tt.want.Size, got.Size)
			require.Len(t, got.Data, len(tt.want.Data))
			for i := range tt.want.Data {
				assert.Equal(t, tt.want.Data[i].ID, got.Data[i].ID)
				assert.Equal(t, tt.want.Data[i].Name, got.Data[i].Name)
				assert.Equal(t, tt.want.Data[i].Driver, got.Data[i].Driver)
				assert.Equal(t, tt.want.Data[i].Subnet, got.Data[i].Subnet)
				assert.Equal(t, tt.want.Data[i].Gateway, got.Data[i].Gateway)
				assert.Equal(t, tt.want.Data[i].Labels, got.Data[i].Labels)
			}
		})
	}
}
