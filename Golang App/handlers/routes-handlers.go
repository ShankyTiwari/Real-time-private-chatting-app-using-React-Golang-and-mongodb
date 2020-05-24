package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"

	constants "private-chat/constants"
)

// RenderHome Rendering the Home Page
func RenderHome(responseWriter http.ResponseWriter, request *http.Request) {
	response := APIResponseStruct{
		Code:     http.StatusOK,
		Status:   http.StatusText(http.StatusOK),
		Message:  constants.APIWelcomeMessage,
		Response: nil,
	}
	ReturnResponse(responseWriter, request, response)
}

//IsUsernameAvailable function will handle the availability of username
func IsUsernameAvailable(responseWriter http.ResponseWriter, request *http.Request) {
	type usernameAvailableResposeStruct struct {
		IsUsernameAvailable bool `json:"isUsernameAvailable"`
	}
	var response APIResponseStruct
	var IsAlphaNumeric = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$`).MatchString
	username := mux.Vars(request)["username"]

	// Checking if username is not empty & has only AlphaNumeric charecters
	if !IsAlphaNumeric(username) {
		response := APIResponseStruct{
			Code:     http.StatusBadRequest,
			Status:   http.StatusText(http.StatusBadRequest),
			Message:  constants.UsernameCantBeEmpty,
			Response: nil,
		}
		ReturnResponse(responseWriter, request, response)
	} else {
		isUsernameAvailable := IsUsernameAvailableQueryHandler(username)
		if isUsernameAvailable {
			response = APIResponseStruct{
				Code:    http.StatusOK,
				Status:  http.StatusText(http.StatusOK),
				Message: constants.UsernameIsAvailable,
				Response: usernameAvailableResposeStruct{
					IsUsernameAvailable: isUsernameAvailable,
				},
			}
		} else {
			response = APIResponseStruct{
				Code:    http.StatusOK,
				Status:  http.StatusText(http.StatusOK),
				Message: constants.UsernameIsNotAvailable,
				Response: usernameAvailableResposeStruct{
					IsUsernameAvailable: isUsernameAvailable,
				},
			}
		}
		ReturnResponse(responseWriter, request, response)
	}
}

//Login function will login the users
func Login(responseWriter http.ResponseWriter, request *http.Request) {
	var userDetails UserDetailsRequestPayloadStruct

	decoder := json.NewDecoder(request.Body)
	requestDecoderError := decoder.Decode(&userDetails)
	defer request.Body.Close()

	if requestDecoderError != nil {
		response := APIResponseStruct{
			Code:     http.StatusBadRequest,
			Status:   http.StatusText(http.StatusBadRequest),
			Message:  constants.UsernameAndPasswordCantBeEmpty,
			Response: nil,
		}
		ReturnResponse(responseWriter, request, response)
	} else {
		if userDetails.Username == "" {
			response := APIResponseStruct{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.UsernameCantBeEmpty,
				Response: nil,
			}
			ReturnResponse(responseWriter, request, response)
		} else if userDetails.Password == "" {
			response := APIResponseStruct{
				Code:     http.StatusInternalServerError,
				Status:   http.StatusText(http.StatusInternalServerError),
				Message:  constants.PasswordCantBeEmpty,
				Response: nil,
			}
			ReturnResponse(responseWriter, request, response)
		} else {

			userDetails, loginErrorMessage := LoginQueryHandler(userDetails)

			if loginErrorMessage != nil {
				response := APIResponseStruct{
					Code:     http.StatusNotFound,
					Status:   http.StatusText(http.StatusNotFound),
					Message:  loginErrorMessage.Error(),
					Response: nil,
				}
				ReturnResponse(responseWriter, request, response)
			} else {
				response := APIResponseStruct{
					Code:     http.StatusOK,
					Status:   http.StatusText(http.StatusOK),
					Message:  constants.UserLoginCompleted,
					Response: userDetails,
				}
				ReturnResponse(responseWriter, request, response)
			}
		}
	}
}

//Registertation function will login the users
func Registertation(responseWriter http.ResponseWriter, request *http.Request) {
	var userDetailsRequestPayload UserDetailsRequestPayloadStruct

	decoder := json.NewDecoder(request.Body)
	requestDecoderError := decoder.Decode(&userDetailsRequestPayload)
	defer request.Body.Close()

	if requestDecoderError != nil {
		response := APIResponseStruct{
			Code:     http.StatusBadRequest,
			Status:   http.StatusText(http.StatusBadRequest),
			Message:  constants.ServerFailedResponse,
			Response: nil,
		}
		ReturnResponse(responseWriter, request, response)
	} else {
		if userDetailsRequestPayload.Username == "" {
			response := APIResponseStruct{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.UsernameCantBeEmpty,
				Response: nil,
			}
			ReturnResponse(responseWriter, request, response)
		} else if userDetailsRequestPayload.Password == "" {
			response := APIResponseStruct{
				Code:     http.StatusInternalServerError,
				Status:   http.StatusText(http.StatusInternalServerError),
				Message:  constants.PasswordCantBeEmpty,
				Response: nil,
			}
			ReturnResponse(responseWriter, request, response)
		} else {
			userObjectID, registrationError := RegisterQueryHandler(userDetailsRequestPayload)
			if registrationError != nil {
				response := APIResponseStruct{
					Code:     http.StatusInternalServerError,
					Status:   http.StatusText(http.StatusInternalServerError),
					Message:  constants.ServerFailedResponse,
					Response: nil,
				}
				ReturnResponse(responseWriter, request, response)
			} else {
				response := APIResponseStruct{
					Code:    http.StatusOK,
					Status:  http.StatusText(http.StatusOK),
					Message: constants.UserRegistrationCompleted,
					Response: UserDetailsResponsePayloadStruct{
						Username: userDetailsRequestPayload.Username,
						UserID:   userObjectID,
					},
				}
				ReturnResponse(responseWriter, request, response)
			}
		}
	}
}

//UserSessionCheck function will check login status of the user
func UserSessionCheck(responseWriter http.ResponseWriter, request *http.Request) {
	var IsAlphaNumeric = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$`).MatchString
	userID := mux.Vars(request)["userID"]

	if !IsAlphaNumeric(userID) {
		response := APIResponseStruct{
			Code:     http.StatusBadRequest,
			Status:   http.StatusText(http.StatusBadRequest),
			Message:  constants.UsernameCantBeEmpty,
			Response: nil,
		}
		ReturnResponse(responseWriter, request, response)
	} else {
		uerDetails := GetUserByUserID(userID)
		if uerDetails == (UserDetailsStruct{}) {
			response := APIResponseStruct{
				Code:     http.StatusOK,
				Status:   http.StatusText(http.StatusOK),
				Message:  constants.YouAreNotLoggedIN,
				Response: false,
			}
			ReturnResponse(responseWriter, request, response)
		} else {
			response := APIResponseStruct{
				Code:     http.StatusOK,
				Status:   http.StatusText(http.StatusOK),
				Message:  constants.YouAreLoggedIN,
				Response: uerDetails.Online == "Y",
			}
			ReturnResponse(responseWriter, request, response)
		}
	}
}

//GetMessagesHandler function will fetch the messages between two users
func GetMessagesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var IsAlphaNumeric = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$`).MatchString
	toUserID := mux.Vars(request)["toUserID"]
	fromUserID := mux.Vars(request)["fromUserID"]

	if !IsAlphaNumeric(toUserID) {
		response := APIResponseStruct{
			Code:     http.StatusBadRequest,
			Status:   http.StatusText(http.StatusBadRequest),
			Message:  constants.UsernameCantBeEmpty,
			Response: nil,
		}
		ReturnResponse(responseWriter, request, response)
	} else if !IsAlphaNumeric(fromUserID) {
		response := APIResponseStruct{
			Code:     http.StatusBadRequest,
			Status:   http.StatusText(http.StatusBadRequest),
			Message:  constants.UsernameCantBeEmpty,
			Response: nil,
		}
		ReturnResponse(responseWriter, request, response)
	} else {
		conversations := GetConversationBetweenTwoUsers(toUserID, fromUserID)
		response := APIResponseStruct{
			Code:     http.StatusOK,
			Status:   http.StatusText(http.StatusOK),
			Message:  constants.UsernameIsAvailable,
			Response: conversations,
		}
		ReturnResponse(responseWriter, request, response)
	}
}
