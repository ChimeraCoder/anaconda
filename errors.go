package anaconda

const (
	//Error code defintions match the Twitter documentation
	//https://dev.twitter.com/docs/error-codes-responses
	TwitterErrorCouldNotAuthenticate    = 32
	TwitterErrorDoesNotExist            = 34
	TwitterErrorAccountSuspended        = 64
	TwitterErrorApi1Deprecation         = 68 //This should never be needed
	TwitterErrorRateLimitExceeded       = 88
	TwitterErrorInvalidToken            = 89
	TwitterErrorOverCapacity            = 130
	TwitterErrorInternalError           = 131
	TwitterErrorCouldNotAuthenticateYou = 135
	TwitterErrorStatusIsADuplicate      = 187
	TwitterErrorBadAuthenticationData   = 215
	TwitterErrorUserMustVerifyLogin     = 231
)

//TwitterErrorResponse has an array of Twitter error messages
//It satisfies the "error" interface
//For the most part, Twitter seems to return only a single error message
//Currently, we assume that this always contains exactly one error message
type TwitterErrorResponse struct {
	errors []TwitterError `json:"errors"`
}

func (tr TwitterErrorResponse) First() error {
	return tr.errors[0]
}

func (tr TwitterErrorResponse) Error() string {
	return tr.errors[0].Message
}

//TwitterError represents a single Twitter error messages/code pair
type TwitterError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (te TwitterError) Error() string {
	return te.Message
}
