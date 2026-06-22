// Package compose
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package compose

import "testing"

func TestGenerateDockerCompose(t *testing.T) {
	filePath := "/Users/amu/Desktop/docker-compose.yml"
	err := GenerateDockerCompose(filePath)
	if err != nil {
		t.Fatal(err)
	}
}
