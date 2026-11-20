// Package docker
// Date: 2024/07/09 14:14:10
// Author: Amu
// Description:
package docker

import (
	"bytes"
	"context"
	"errors"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type ImageSummary struct {
	ID      string
	Name    string
	Tag     string
	Created string
	Size    string
}

// splitRepoTag 按镜像引用的最后一个「标签冒号」拆分名称与标签：
// registry:5000/img:v1 中的端口冒号不是标签分隔符；无冒号或冒号后
// 仍含 "/" 的引用整体视为名称并返回空标签，避免固定下标越界。
func splitRepoTag(repoTag string) (name, tag string) {
	idx := strings.LastIndex(repoTag, ":")
	if idx >= 0 && !strings.Contains(repoTag[idx+1:], "/") {
		return repoTag[:idx], repoTag[idx+1:]
	}
	return repoTag, ""
}

// ShortImageID 返回用于展示的短镜像 ID：剥离 sha256: 前缀后截取
// 前 12 个字符；非 sha256 前缀或更短的 ID 原样保留，不做固定下标截取。
func ShortImageID(id string) string {
	const digestPrefix = "sha256:"
	if strings.HasPrefix(id, digestPrefix) {
		id = id[len(digestPrefix):]
	}
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func (m *Manager) ListImage(ctx context.Context) ([]ImageSummary, error) {
	images, err := m.client.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var imageList []ImageSummary
	for _, im := range images {
		if len(im.RepoTags) == 0 {
			continue
		}
		for _, repoTag := range im.RepoTags {
			name, tag := splitRepoTag(repoTag)
			im := ImageSummary{
				ID:      im.ID,
				Name:    name,
				Tag:     tag,
				Created: time.Unix(im.Created, 0).Format("2006-01-02 15:04:05"),
				Size:    strconv.FormatFloat(float64(im.Size)/(1000*1000), 'f', 2, 64) + "MB",
			}
			imageList = append(imageList, im)
		}
	}
	return imageList, nil
}

func (m *Manager) DeleteImage(ctx context.Context, imageID string) error {
	_, err := m.client.ImageRemove(ctx, imageID, image.RemoveOptions{Force: true})
	return err
}

func (m *Manager) PruneImages(ctx context.Context) error {
	_, err := m.client.ImagesPrune(ctx, filters.NewArgs(filters.Arg("dangling", "true")))
	return err
}

func (m *Manager) SearchImage(ctx context.Context, imageName string) ([]registry.SearchResult, error) {
	return m.client.ImageSearch(ctx, imageName, registry.SearchOptions{
		Limit: 10,
	})
}

func (m *Manager) PullImage(ctx context.Context, imageName string) error {
	pullReader, err := m.client.ImagePull(ctx, imageName, image.PullOptions{All: false, PrivilegeFunc: nil, RegistryAuth: ""})
	if err != nil {
		return err
	}
	defer func(pullReader io.ReadCloser) {
		err := pullReader.Close()
		if err != nil {
			return
		}
	}(pullReader)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(pullReader)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) TagImage(ctx context.Context, oldTag, newTag string) error {
	return m.client.ImageTag(ctx, oldTag, newTag)
}

func (m *Manager) ImportImage(ctx context.Context, sourceFile string) error {
	root, name, err := openFileRoot(sourceFile)
	if err != nil {
		return err
	}
	defer root.Close()
	inputFile, err := root.Open(name)
	if err != nil {
		return err
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			return
		}
	}(inputFile)

	resp, err := m.client.ImageLoad(ctx, inputFile, true)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// 读取并输出导入过程
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) ExportImage(ctx context.Context, imageIDs []string, targetFile string) error {
	resp, err := m.client.ImageSave(ctx, imageIDs)
	if err != nil {
		return err
	}
	defer func(resp io.ReadCloser) {
		err := resp.Close()
		if err != nil {
			return
		}
	}(resp)
	root, name, err := openFileRoot(targetFile)
	if err != nil {
		return err
	}
	defer root.Close()
	outputFile, err := root.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if err := root.Chmod(name, 0o600); err != nil {
		_ = outputFile.Close()
		return err
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			return
		}
	}(outputFile)

	_, err = io.Copy(outputFile, resp)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetImageByName(ctx context.Context, imageName string) (*ImageSummary, error) {
	images, err := m.client.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, v := range images {
		for _, t := range v.RepoTags {
			if t == imageName {
				name, tag := splitRepoTag(t)
				return &ImageSummary{
					ID:      v.ID,
					Name:    name,
					Tag:     tag,
					Created: time.Unix(v.Created, 0).Format("2006-01-02 15:04:05"),
					Size:    strconv.FormatFloat(float64(v.Size)/(1000*1000), 'f', 2, 64) + "MB",
				}, nil
			}
		}
	}
	return nil, errors.New("not found image")
}

func (m *Manager) GetImageByID(ctx context.Context, imageID string) (*ImageSummary, error) {
	imageResponse, _, err := m.client.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return nil, err
	}

	if len(imageResponse.RepoTags) == 0 {
		return nil, errors.New("image has no repo tags")
	}
	name, tag := splitRepoTag(imageResponse.RepoTags[0])

	return &ImageSummary{
		ID:      imageResponse.ID,
		Name:    name,
		Tag:     tag,
		Created: imageResponse.Created,
		Size:    strconv.FormatFloat(float64(imageResponse.Size)/(1000*1000), 'f', 2, 64) + "MB",
	}, nil
}
