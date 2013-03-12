package anaconda

const (
	//As defined on https://dev.twitter.com/docs/error-codes-responses
	TwitterErrorDoesNotExist            = 34
	TwitterErrorRateLimitExceeded       = 88
	TwitterErrorInvalidToken            = 89
	TwitterErrorOverCapacity            = 130
	TwitterErrorInternalError           = 131
	TwitterErrorCouldNotAuthenticateYou = 135
	TwitterErrorBadAuthenticationData   = 215
)

//This is probably unnecessary, but technically Twitter error codes are separate from the HTTP statuses, so it's separate for now.
type ApiError struct {
	errorString   string
	httpStatus    int
	TwitterErrors error //If non-nil, this will be a TwitterError struct.
	//Using 'error' as the type should solve the nil error:
	//http://golang.org/doc/faq#nil_error
	requestUrl string //If this was in response to a request, which endpoint?
}

func (e ApiError) Error() string {
	return e.errorString
}

//HttpCode provides the HTTP status code returned, if the error originated with a HTTP request
func (e ApiError) HTTPStatus() int {
	return e.httpStatus
}

//TwitterError corresponds to the JSON errors that Twitter may return in API queries
type TwitterError struct {
	Message   string
	Code      int
	NextError error //Will be non-nil if Twitter returned more than one error
}

//OrMap returns true if the function evalutes to true on any TwitterError later in the list
func (c TwitterError) OrMap(f func(TwitterError) bool) bool {
	if f(c) {
		return true
	}
	if c.NextError == nil {
		return false
	}
	return c.NextError.(TwitterError).OrMap(f)
}

//ContainsError returns true if the current error or any later error in the list matches the error code specified.
func (err TwitterError) ContainsError(code int) bool {
	return err.OrMap(func(e TwitterError) bool {
		return e.Code == code
	})
}

//Error is in included just to sastisfy the error interface.
func (e TwitterError) Error() string {
	return e.Message
}

//Internal struct used only to unmarshal the error response that Twitter provides
type twitterErrorResponse struct {
	Errors []TwitterError
}
