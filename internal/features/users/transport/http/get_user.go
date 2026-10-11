package users_transport_http

import (
	"net/http"

	core_logger "gitverse.ru/IvanVerse/go-todoApp/internal/core/logger"
	core_http_request "gitverse.ru/IvanVerse/go-todoApp/internal/core/transport/http/request"
	core_http_response "gitverse.ru/IvanVerse/go-todoApp/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromCOntext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}
	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}
	response := GetUserResponse(userDtoFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)

}
