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
			tags := strings.Split(repoTag, ":")
			im := ImageSummary{
				ID:      im.ID,
				Name:    tags[0],
				Tag:     tags[1],
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
	inputFile, err := os.Open(sourceFile)
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
	outputFile, err := os.Create(targetFile)
	if err != nil {
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
				tagsList := strings.Split(t, ":")
				return &ImageSummary{
					ID:      v.ID,
					Name:    tagsList[0],
					Tag:     tagsList[1],
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

	tagsList := strings.Split(imageResponse.RepoTags[0], ":")

	return &ImageSummary{
		ID:      imageResponse.ID,
		Name:    tagsList[0],
		Tag:     tagsList[1],
		Created: imageResponse.Created,
		Size:    strconv.FormatFloat(float64(imageResponse.Size)/(1000*1000), 'f', 2, 64) + "MB",
	}, nil
}
