package controller

import (
	"goserver/apierrors"
	"goserver/model"
	"goserver/presentation/dto"
	"goserver/service"
	"goserver/utils/gslog"
	"goserver/utils/gsmiddleware"
	"goserver/utils/gsrender"
	"goserver/utils/gsvalidation"
	"net/http"
	"strconv"

	dtomapper "github.com/dranikpg/dto-mapper"
	"github.com/go-chi/chi/v5"
)

// GetUsers godoc
// @Summary 	Busca todos los usuarios
// @Description Permite la búsqueda de todos los usuarios (no utiliza paginación)
// @ID 			get-users
// @Tags 		Users
// @Success 	200 {array}  dto.UserDTO
// @Failure 	500 {object} gserrors.Error
// @Router 		/users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, gserror := service.FindUsers(r.Context())
	if gserror != nil {
		gsrender.WriteJSON(w, http.StatusInternalServerError, gserror)
		return
	}
	var userDTOs []dto.UserDTO
	err := dtomapper.Map(userDTOs, users)
	if err != nil {
		gslog.ErrorFrom(err, gsmiddleware.GetTraceID(r.Context()))
		gsrender.WriteJSON(w, http.StatusInternalServerError, apierrors.New(apierrors.INTERNAL_SERVER_ERROR))
		return
	}
	gsrender.WriteJSON(w, http.StatusOK, users)
}

// PostUser godoc
// @Summary 	Crea un nuevo usuario
// @Description Permite crear un nuevo usuario
// @ID 			post-user
// @Tags 		Users
// @Success 	201
// @Failure 	500 {object} gserrors.Error
// @Router 		/users [post]
func PostUser(w http.ResponseWriter, r *http.Request) {

	var reqDTO dto.CreateUserRequestDTO
	if httpSuggestion := gsvalidation.DecodeJSONRequestBody(r, &reqDTO); httpSuggestion != nil {
		code := apierrors.INTERNAL_SERVER_ERROR
		if httpSuggestion.Status != 500 {
			code = apierrors.INVALID_ARGUMENT
		}
		gsrender.WriteJSON(w, httpSuggestion.Status, apierrors.NewWithMsg(code, httpSuggestion.Message))
		return
	}

	var reqModel model.CreateUserRequest
	err := dtomapper.Map(&reqModel, &reqDTO)
	if err != nil {
		gslog.ErrorFrom(err, gsmiddleware.GetTraceID(r.Context()))
		gsrender.WriteJSON(w, http.StatusInternalServerError, apierrors.New(apierrors.INTERNAL_SERVER_ERROR))
		return
	}
	gserror := service.CreateUser(r.Context(), reqModel)
	if gserror != nil {
		gsrender.WriteJSON(w, http.StatusInternalServerError, gserror)
		return
	}
	gsrender.Status(w, http.StatusCreated)
}

// GetUserById godoc
// @Summary 	Busca un usuario por su ID
// @Description Permite la búsqueda de un usuario a través de su ID
// @ID 			get-user-by-id
// @Tags 		Users
// @Param 		id 	path 	 int true "The ID of a user"
// @Success 	200 {object} dto.UserDTO
// @Router 		/users/{id} [get]
func GetUserById(w http.ResponseWriter, r *http.Request) {

	strID := chi.URLParam(r, "id")

	userID, e := strconv.Atoi(strID)
	if (e != nil) {
		gsrender.WriteJSON(w, http.StatusBadRequest, apierrors.New(apierrors.INVALID_ARGUMENT))
		return
	}
	
	user, gserror := service.FindUserById(r.Context(),userID)
	if gserror != nil {
		gsrender.WriteJSON(w, http.StatusInternalServerError, gserror)
		return
	}

	if (user == nil) {
		gslog.ErrorFrom(gserror, gsmiddleware.GetTraceID(r.Context()))
		gsrender.WriteJSON(w, http.StatusNotFound, apierrors.New(apierrors.USER_NOT_FOUND))
		return
	}

	var userDTO dto.UserDTO
	err := dtomapper.Map(&userDTO, &user)
	if err != nil {
		gslog.ErrorFrom(err, gsmiddleware.GetTraceID(r.Context()))
		gsrender.WriteJSON(w, http.StatusInternalServerError, apierrors.New(apierrors.INTERNAL_SERVER_ERROR))
		return
	}
	gsrender.WriteJSON(w, http.StatusOK, userDTO)
}