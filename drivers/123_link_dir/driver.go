package _123LinkDir

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
)

const DIRVER_API = "https://open-api.123pan.com"

type Pan123LinkDir struct {
	model.Storage
	Addition
}

func (d *Pan123LinkDir) Config() driver.Config {
	return config
}

func (d *Pan123LinkDir) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Pan123LinkDir) Init(ctx context.Context) error {
	// TODO 登录 / 刷新令牌
	req := base.RestyClient.R()
	req.SetHeader(
		"Platform", "open_platform",
	)
	req.SetFormData(map[string]string{
		"client_id":     d.ClientID,
		"client_secret": d.ClientSecret,
	})

	res, err := req.Execute(http.MethodPost, OpenAPIBaseURL+"/api/v1/access_token")
	if err != nil {
		return err
	}

	body := res.Body()

	resStruct := struct {
		Data struct {
			AccessToken string `json:"accessToken"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &resStruct)
	if err != nil {
		return err
	}

	d.access_token = resStruct.Data.AccessToken

	return nil
}

func (d *Pan123LinkDir) Drop(ctx context.Context) error {
	return nil
}

func (d *Pan123LinkDir) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	url := OpenAPIBaseURL + "/api/v2/file/list"

	req := base.RestyClient.R()
	parentID := dir.GetID()
	if parentID == "" && d.RootFolderID != 0 {
		parentID = fmt.Sprintf("%d", d.RootFolderID)
	} else if parentID == "" {
		parentID = "0"
	}

	req.SetQueryParam("parentFileId", parentID)
	req.SetQueryParam("limit", "100")
	req.SetHeader("Authorization", "Bearer "+d.access_token)
	req.SetHeader("Platform", "open_platform")
	res, err := req.Execute(http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	body := res.Body()
	bodyStruct := struct {
		Data struct {
			FileList []File `json:"fileList"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &bodyStruct)
	if err != nil {
		return nil, err
	}

	objs := make([]model.Obj, 0)
	for _, file := range bodyStruct.Data.FileList {
		objs = append(objs, &file)
	}

	return objs, nil
}

func (d *Pan123LinkDir) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	protocol := "http"
	if d.EnableHTTPS {
		protocol = "https"
	}
	var url string
	if d.UUID != "" {
		url = fmt.Sprintf("%s://%s/%s/%s", protocol, d.Domain, d.UUID, file.GetID())
	} else {
		url = fmt.Sprintf("%s://%s/%s", protocol, d.Domain, file.GetID())
	}

	return &model.Link{
		URL: url,
	}, nil
}

func (d *Pan123LinkDir) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) (model.Obj, error) {
	// TODO 创建文件夹，选填
	return nil, errs.NotImplement
}

func (d *Pan123LinkDir) Move(ctx context.Context, srcObj, dstDir model.Obj) (model.Obj, error) {
	// TODO 移动对象，选填
	return nil, errs.NotImplement
}

func (d *Pan123LinkDir) Rename(ctx context.Context, srcObj model.Obj, newName string) (model.Obj, error) {
	// TODO 重命名对象，选填
	return nil, errs.NotImplement
}

func (d *Pan123LinkDir) Copy(ctx context.Context, srcObj, dstDir model.Obj) (model.Obj, error) {
	// TODO 复制对象，选填
	return nil, errs.NotImplement
}

func (d *Pan123LinkDir) Remove(ctx context.Context, obj model.Obj) error {
	// TODO 删除对象，选填
	return errs.NotImplement
}

func (d *Pan123LinkDir) Put(ctx context.Context, dstDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) (model.Obj, error) {
	// TODO 上传文件，选填
	return nil, errs.NotImplement
}

//func (d *Pan123LinkDir) Other(ctx context.Context, args model.OtherArgs) (interface{}, error) {
//	return nil, errs.NotSupport
//}

var _ driver.Driver = (*Pan123LinkDir)(nil)
