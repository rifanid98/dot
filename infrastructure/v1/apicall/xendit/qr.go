package xendit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"dot/core"
	"dot/pkg/api"
	"dot/pkg/helper"
)

func (x *xenditApiCallImpl) QRCreate(ic *core.InternalContext, data map[string]any) (map[string]any, *core.CustomError) {
	var apiReq api.HttpParam

	payload := new(QrCodesCreate).Bind(data)

	header := make(map[string]string)
	header["api-version"] = x.cfg.ApiVersion20220731
	header["Authorization"] = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(x.cfg.ApiKey+":")))
	header["Content-Type"] = "application/json"

	apiReq.Header = header
	apiReq.Method = "post"
	apiReq.Timeout = x.cfg.Timeout
	apiReq.Url = fmt.Sprintf("%s%s", x.cfg.BaseUrl, x.cfg.PathQrCodes)
	apiReq.Body = helper.DataToString(payload)

	log.Info(ic.ToContext(), "xendit create qr code : %v", apiReq)

	res, err := x.client.HttpDo(apiReq)
	if err != nil {
		log.Error(ic.ToContext(), "failed xendit create qr code : %v", err)

		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVICE_ERROR,
		}
	}
	defer res.Body.Close()

	log.Info(ic.ToContext(), "response xendit create qr code : %v", api.DumpResponse(res))

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(ic.ToContext(), "failed to read response body http request from xendit create qr code", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	response := make(map[string]any)
	if err = json.Unmarshal(resp, &response); err != nil {
		log.Error(ic.ToContext(), "failed to bind response data http request from xendit service", err)
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	errCode := response["error_code"]
	if errCode == nil {
		return response, nil
	}
	errMsg := response["message"]
	if errMsg == nil {
		return response, nil
	}

	errorCode := helper.DataToString(errCode)
	errorMessage := helper.DataToString(errMsg)

	if errorCode == "API_VALIDATION_ERROR" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_BAD_REQUEST,
			Message: errorMessage,
		}
	}
	if errorCode == "UNSUPPORTED_CURRENCY" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_BAD_REQUEST,
			Message: errorMessage,
		}
	}
	if errorCode == "INVALID_API_KEY" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_UNAUTHORIZED,
			Message: errorMessage,
		}
	}
	if errorCode == "INVALID_MERCHANT_CREDENTIALS" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_UNAUTHORIZED,
			Message: errorMessage,
		}
	}
	if errorCode == "REQUEST_FORBIDDEN_ERROR" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_FORBIDDEN,
			Message: errorMessage,
		}
	}
	if errorCode == "CHANNEL_NOT_ACTIVATED" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_FORBIDDEN,
			Message: errorMessage,
		}
	}
	if errorCode == "DUPLICATE_ERROR" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_CONFLICT,
			Message: errorMessage,
		}
	}
	if errorCode == "CALLBACK_URL_NOT_FOUND" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_NOT_FOUND,
			Message: errorMessage,
		}
	}
	if errorCode == "SERVER_ERROR" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVICE_ERROR,
			Message: errorMessage,
		}
	}
	if errorCode == "CHANNEL_UNAVAILABLE" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVICE_UNAVAILABLE,
			Message: errorMessage,
		}
	}

	return response, nil
}

func (x *xenditApiCallImpl) QRCheck(ic *core.InternalContext, data map[string]any) (map[string]any, *core.CustomError) {
	var apiReq api.HttpParam

	header := make(map[string]string)
	header["api-version"] = x.cfg.ApiVersion20220731
	header["Authorization"] = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(x.cfg.ApiKey+":")))
	header["Content-Type"] = "application/json"

	apiReq.Header = header
	apiReq.Method = "get"
	apiReq.Timeout = x.cfg.Timeout
	apiReq.Url = fmt.Sprintf("%s%s/%v", x.cfg.BaseUrl, x.cfg.PathQrCodes, helper.DataToString(data["id_qr"]))

	log.Info(ic.ToContext(), "xendit get qr code : %v", apiReq)

	res, err := x.client.HttpDo(apiReq)
	if err != nil {
		log.Error(ic.ToContext(), "failed xendit get qr code : %v", err)

		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVICE_ERROR,
		}
	}
	defer res.Body.Close()

	log.Info(ic.ToContext(), "response xendit get qr code : %v", api.DumpResponse(res))

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(ic.ToContext(), "failed to read response body http request from xendit get qr code", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	response := make(map[string]any)
	if err = json.Unmarshal(resp, &response); err != nil {
		log.Error(ic.ToContext(), "failed to bind response data http request from xendit service", err)
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	errCode := response["error_code"]
	if errCode == nil {
		return response, nil
	}
	errMsg := response["message"]
	if errMsg == nil {
		return response, nil
	}

	errorCode := helper.DataToString(errCode)
	errorMessage := helper.DataToString(errMsg)

	if errorCode == "API_VALIDATION_ERROR" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_BAD_REQUEST,
			Message: errorMessage,
		}
	}
	if errorCode == "INVALID_API_KEY" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_UNAUTHORIZED,
			Message: errorMessage,
		}
	}
	if errorCode == "DATA_NOT_FOUND" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_NOT_FOUND,
			Message: errorMessage,
		}
	}
	if errorCode == "SERVER_ERROR" {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVICE_ERROR,
			Message: errorMessage,
		}
	}

	return response, nil
}
