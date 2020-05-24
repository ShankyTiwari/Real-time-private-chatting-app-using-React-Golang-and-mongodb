package constants

// Constant will export all application Constants
const (
	// Request Paramyter Validation error/success messages
	UsernameCantBeEmpty            = "Username can't be empty."
	UsernameIsAvailable            = "Username is available."
	UsernameIsNotAvailable         = "Username is not available."
	PasswordCantBeEmpty            = "Password can't be empty."
	UsernameAndPasswordCantBeEmpty = "Username and Password can't be empty."
	LoginPasswordIsInCorrect       = "Your Login Password is incorrect."
	UserRegistrationCompleted      = "User Registration Completed."
	UserLoginCompleted             = "User Login is Completed."
	YouAreNotLoggedIN              = "You are not logged in."
	YouAreLoggedIN                 = "You are logged in."
	UserIsNotRegisteredWithUs      = "This account does not exist in our system."

	// Application response messages
	SuccessfulResponse   = "Request completed successfully"
	ServerFailedResponse = "Request failed to complete, we are working on it"
	APIWelcomeMessage    = "This is an API for Realtime Private chat application build in GoLang"
)
