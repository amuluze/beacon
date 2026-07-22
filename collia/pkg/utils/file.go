// Package utils
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SanitizePath validates that the path is absolute and does not contain
// directory traversal sequences. It returns the cleaned absolute path or an error.
func SanitizePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	cleaned := filepath.Clean(path)
	if !filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("path must be absolute: %s", path)
	}
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path contains traversal sequence: %s", path)
	}
	return cleaned, nil
}

// resolveForSandbox 解析路径的符号链接，返回可用于根目录前缀比较的绝对路径。
// 对尚不存在的路径（如待创建的文件），EvalSymlinks 会失败，此时逐级向上
// 找到最近存在的祖先目录求值，再拼回尾部——保证新建文件也被约束在沙箱内，
// 且与根目录使用一致的符号链接解析（例如 macOS /var -> /private/var）。
func resolveForSandbox(path string) (string, error) {
	cleaned := filepath.Clean(path)
	if !filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("path must be absolute: %s", path)
	}
	resolved, err := filepath.EvalSymlinks(cleaned)
	if err == nil {
		return resolved, nil
	}
	// 路径不存在（创建场景）：逐级向上找到最近存在的祖先目录求值，
	// 再拼回尚未存在的尾部，保证符号链接被一致解析。
	dir := filepath.Dir(cleaned)
	base := filepath.Base(cleaned)
	tail := []string{base}
	for dir != "/" && dir != "." {
		resolvedDir, err := filepath.EvalSymlinks(dir)
		if err == nil {
			parts := append([]string{resolvedDir}, reverse(tail)...)
			return filepath.Join(parts...), nil
		}
		tail = append(tail, filepath.Base(dir))
		dir = filepath.Dir(dir)
	}
	// 整条路径都不存在，直接用 cleaned 做前缀比较（至少拦绝对越界）。
	return cleaned, nil
}

func reverse(s []string) []string {
	out := make([]string, len(s))
	for i, v := range s {
		out[len(s)-1-i] = v
	}
	return out
}

// withinRoot 报告 path 解析后是否严格位于 root 之内（含等于 root）。
// 通过解析符号链接防止软链逃逸：即使 /data/amprobe/evil -> /etc，
// 求值后 /etc 也不会匹配 root 前缀。
func withinRoot(path, root string) bool {
	resolvedPath, err := resolveForSandbox(path)
	if err != nil {
		return false
	}
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		resolvedRoot = filepath.Clean(root)
	}
	if resolvedPath == resolvedRoot {
		return true
	}
	return strings.HasPrefix(resolvedPath, resolvedRoot+string(filepath.Separator))
}

// SanitizePathWithin 校验 path 解析后严格位于 root 沙箱之内，
// 返回 cleaned 后的绝对路径或错误。相比 SanitizePath：
//   - 限制根目录，防止访问 /etc/passwd 等系统文件；
//   - 解析符号链接，防止软链逃逸；
//   - 对尚不存在的路径回退到父目录求值，保证创建场景仍受约束。
//
// root 必须是绝对路径；为空时报错（调用方应注入配置的根目录）。
func SanitizePathWithin(path, root string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	if root == "" {
		return "", fmt.Errorf("sandbox root must not be empty")
	}
	cleanedRoot := filepath.Clean(root)
	if !filepath.IsAbs(cleanedRoot) {
		return "", fmt.Errorf("sandbox root must be absolute: %s", root)
	}
	cleaned := filepath.Clean(path)
	if !filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("path must be absolute: %s", path)
	}
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path contains traversal sequence: %s", path)
	}
	if !withinRoot(cleaned, cleanedRoot) {
		return "", fmt.Errorf("path %s is outside sandbox root %s", path, cleanedRoot)
	}
	return cleaned, nil
}

func CopyFile(src, dst string) (int64, error) {
	src, err := SanitizePath(src)
	if err != nil {
		return 0, err
	}
	dst, err = SanitizePath(dst)
	if err != nil {
		return 0, err
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src) //#nosec G304 -- path sanitized by SanitizePath
	if err != nil {
		return 0, err
	}
	defer func() { _ = source.Close() }()

	destination, err := os.Create(dst) //#nosec G304 -- path sanitized by SanitizePath
	if err != nil {
		return 0, err
	}
	defer func() { _ = destination.Close() }()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
