package httpserver

import (
	utilis "go_chi_pgx/utils"
	"net/http"
)

var ValidDataNotFound = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "The provided information is invalid. Please recheck and try again.",
}

var InvalidEmailPassword = utilis.ResponseState{
	StatusCode: http.StatusUnauthorized,
	Message:    "Invalid email or password",
}
var UserNotActive = utilis.ResponseState{
	StatusCode: http.StatusUnauthorized,
	Message:    "User not active",
}
var InvalidToken = utilis.ResponseState{
	StatusCode: http.StatusUnauthorized,
	Message:    "Invalid token",
}
var InvalidId = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "Invalid contact ID",
}
var InvalidUserId = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "Invalid user ID",
}
var InternalError = utilis.ResponseState{
	StatusCode: http.StatusInternalServerError,
	Message:    "Internal server error",
}
var Unauthorized = utilis.ResponseState{
	StatusCode: http.StatusUnauthorized,
	Message:    "Unauthorized user",
}
var UserAlreadyExist = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "User Already Exist With this Email",
}

var UserActivated = utilis.ResponseState{
	StatusCode: http.StatusOK,
	Message:    "User activated successfully",
}

var UserCreated = utilis.ResponseState{
	StatusCode: http.StatusCreated,
	Message:    "User created successfully",
}

var loginSuccess = utilis.ResponseState{
	StatusCode: http.StatusOK,
	Message:    "Login Successful",
}

var ContactRetrieved = utilis.ResponseState{
	StatusCode: http.StatusOK,
	Message:    "Contacts Retrieved successfully",
}

var ContactCreated = utilis.ResponseState{
	StatusCode: http.StatusCreated,
	Message:    "Contacts Created successfully",
}
var UserNotFound = utilis.ResponseState{
	StatusCode: http.StatusNotFound,
	Message:    "User not found",
}
var NotFound = utilis.ResponseState{
	StatusCode: http.StatusNotFound,
	Message:    "Contact Not found",
}
var ContactUpdated = utilis.ResponseState{
	StatusCode: http.StatusCreated,
	Message:    "Contacts Updated successfully",
}
var RateLimitExceeded = utilis.ResponseState{
	StatusCode: http.StatusTooManyRequests,
	Message:    "Rate Limit Exceeded",
}

var BadRequestError = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "Bad Request",
}
