package service

import (
	"context"
	"goserver/apierrors"
	"goserver/client/mockclient"
	"goserver/model"
	"goserver/utils/gslog"
	"goserver/utils/gsmiddleware"

	dtomapper "github.com/dranikpg/dto-mapper"
)

func FindUsers(ctx context.Context) ([]model.User, *apierrors.Error) {
	
	users, err := mockclient.GetUsers(ctx)
	if err != nil {
		gslog.ErrorFrom(err, gsmiddleware.GetTraceID(ctx))
		return nil, apierrors.New(apierrors.EXTERNAL_API_ERROR)
	}

	var resp []model.User

	for _, user := range users {
		u, err := user.ToModel()
		if err != nil {
			gslog.Error(err.Error(), gsmiddleware.GetTraceID(ctx))
			return []model.User{}, apierrors.NewWithMsg(
				apierrors.INTERNAL_SERVER_ERROR,
				"Error found while mapping external API schema to model",
			)
		}
		resp = append(resp, *u)
	}

	return resp, nil
}

func FindUserById(ctx context.Context, id int) (*model.User, *apierrors.Error) {
	
	user, apierror := mockclient.GetUserById(ctx, id)
	if apierror != nil {
		gslog.ErrorFrom(apierror, gsmiddleware.GetTraceID(ctx))
		return nil, apierrors.New(apierrors.EXTERNAL_API_ERROR)
	}

	u, err := user.ToModel()
	if err != nil {
		gslog.Error(err.Error(), gsmiddleware.GetTraceID(ctx))
		return nil, apierrors.NewWithMsg(
			apierrors.INTERNAL_SERVER_ERROR,
			"Error found while mapping external API schema to model",
		)
	}
	
	return u, nil
}

func CreateUser(ctx context.Context, createUserRequest model.CreateUserRequest) *apierrors.Error {
	
	var createUserCDO mockclient.CreateUserCDO
	err := dtomapper.Map(&createUserCDO, &createUserRequest)
	if err != nil {
		gslog.ErrorFrom(err, gsmiddleware.GetTraceID(ctx))
		return apierrors.New(apierrors.INTERNAL_SERVER_ERROR)
	}

	gserror := mockclient.PostUser(ctx, createUserCDO)
	if gserror != nil {
		gslog.ErrorFrom(err, gsmiddleware.GetTraceID(ctx))
		return apierrors.New(apierrors.EXTERNAL_API_ERROR)
	}

	return nil
}
