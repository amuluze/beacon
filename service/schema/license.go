// Package schema
// Date: 2023/6/7 13:03
// Author: Amu
// Description:
package schema

import "mime/multipart"

type LicenseInfo struct {
}

type GetLicenseArgs struct {
}

type GetLicenseReply struct {
	Info    LicenseInfo `json:"info" validate:"omitempty"`
	Expired bool        `json:"expired" description:"是否过期"`
	Exists  bool        `json:"exists" description:"是否存在"`
}

type UploadLicenseArgs struct {
	File *multipart.FileHeader `json:"file" validate:"required" description:"证书文件"`
}
