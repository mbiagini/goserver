package mockclient

import (
	"context"
	"fmt"
	"goserver/apierrors"
	"goserver/config"
	"goserver/utils/gsclient"
	"goserver/utils/gsvalidation"
	"net/http"
	"strconv"
)

// errorCDO is the error type returned by this client.
type errorCDO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func init() {
}

func GetUsers(ctx context.Context) ([]UserCDO, *apierrors.Error) {

	client, ok := gsclient.GetClient(string(config.MOCK_CLIENT))
	if !ok {
		return nil, apierrors.New(apierrors.CLIENT_NOT_DEFINED)
	}

	url := client.Basepath + "/users"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.Basepath + "/users", nil)
	if err != nil {
		return nil, apierrors.NewWithMsg(apierrors.INTERNAL_SERVER_ERROR, fmt.Sprintf("Error creating request to call %s: %s", url, err.Error()))
	}
	
	response, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, apierrors.NewWithMsg(apierrors.INTERNAL_SERVER_ERROR, fmt.Sprintf("Error found while executing http call: %s", err.Error()))
	}

	var users []UserCDO
	var clientErr errorCDO

	respType, err := gsvalidation.DecodeJSONResponseBody(response, &users, &clientErr)
	switch respType {
	case gsvalidation.OK_RESPONSE:
		return users, nil
	case gsvalidation.ERR_RESPONSE:
		return nil, apierrors.NewWithMsg(
			apierrors.EXTERNAL_API_ERROR,
			fmt.Sprintf("Error received by mock client: code: %s, message: %s", clientErr.Code, clientErr.Message),
		)
	default:
		return nil, apierrors.NewWithMsg(apierrors.RESPONSE_UNMARSHAL_ERROR, fmt.Sprintf("Error unmarshaling response body: %s", err.Error()))
	}
}

func PostUser(ctx context.Context, u CreateUserCDO) *apierrors.Error {
	
	client, ok := gsclient.GetClient("MockClient")
	if !ok {
		return apierrors.New(apierrors.CLIENT_NOT_DEFINED)
	}

	path := "/users"
	req, err := client.NewJSONRequest(ctx, "POST", path, u)
	//req, err := http.NewRequestWithContext(ctx, "POST", url, reader)
	if err != nil {
		return apierrors.NewWithMsg(apierrors.INTERNAL_SERVER_ERROR, fmt.Sprintf("Error creating request to call %s: %s", path, err.Error()))
	}

	_, err = client.HttpClient.Do(req)
	if err != nil {
		return apierrors.NewWithMsg(apierrors.INTERNAL_SERVER_ERROR, fmt.Sprintf("Error found while executing http call: %s", err.Error()))
	}

	return nil
}

func GetUserById(ctx context.Context, id int) (*UserCDO, *apierrors.Error) {

	client, ok := gsclient.GetClient("MockClient")
	if !ok {
		return nil, apierrors.New(apierrors.CLIENT_NOT_DEFINED)
	}

	url := client.Basepath + "/users/" + strconv.Itoa(id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, apierrors.NewWithMsg(apierrors.INTERNAL_SERVER_ERROR, fmt.Sprintf("Error creating request to call %s: %s", url, err.Error()))
	}

	response, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, apierrors.NewWithMsg(apierrors.INTERNAL_SERVER_ERROR, fmt.Sprintf("Error found while executing http call: %s", err.Error()))
	}

	var user UserCDO
	var clientErr errorCDO

	respType, err := gsvalidation.DecodeJSONResponseBody(response, &user, &clientErr)
	switch respType {
	case gsvalidation.OK_RESPONSE:
		return &user, nil
	case gsvalidation.ERR_RESPONSE:
		return nil, apierrors.NewWithMsg(
			apierrors.EXTERNAL_API_ERROR,
			fmt.Sprintf("Error received by mock client: code: %s, message: %s", clientErr.Code, clientErr.Message),
		)
	default:
		return nil, apierrors.NewWithMsg(apierrors.RESPONSE_UNMARSHAL_ERROR, fmt.Sprintf("Error unmarshaling response body: %s", err.Error()))
	}
}