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

//All errors, whether originating from Twitter or not, can be safely cast as CodedError
type CodedError interface {
	Error() string
	ErrorCode() int //Renamed from Code() so as not to collide with TwitterError.Code
	//which is based on the field names returned by Twitter
}

type ApiError struct {
	errorString   string
	httpStatus    int
	TwitterErrors error //If non-nil, this will be a TwitterError struct.
	//Using 'error' as the type should solve the nil error:
	//http://golang.org/doc/faq#nil_error
	requestUrl string //If this was in response to a request, which endpoint?
}

//Error returns the Twitter error message
//If the error did not originate from the Twitter API query, it returns the local error message
func (e ApiError) Error() string {
	if e.TwitterErrors == nil {
		return e.errorString
	}
	return e.TwitterErrors.(TwitterError).Error()
}

//Code returns the Twitter error code, or 0 if the error did not originate from the Twitter API
//These can be compared directly with the predefined constants - eg, TwitterErrorDoesNotExist
func (e ApiError) ErrorCode() int {
	if e.TwitterErrors == nil {
		return 0
	}
	return e.TwitterErrors.(TwitterError).Code
}

//HttpCode provides the HTTP status code returned, if the error originated with a HTTP request
func (e ApiError) HTTPStatus() int {
	return e.httpStatus
}

//TwitterError corresponds to the JSON errors that Twitter may return in API queries
//It satisifes the CodedError interface
type TwitterError struct {
	Message   string //The message returned by Error(). Field exported for now to allow JSON unmarshaling.
	Code      int    //Error code according to Twitter's error code scheme
	NextError error  //Will be non-nil if Twitter returned more than one error
	//Equivalent to Messages []CodedError
}

//Error is in included just to satisfy the error and CodedError interfaces.
func (e TwitterError) Error() string {
	return e.Message
}

//Code is included to satisfy the CodedError interface
func (e TwitterError) ErrorCode() int {
	return e.Code
}

//Leaving these two functions in only since there's no reason to remove them (yet)

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
		return e.ErrorCode() == code
	})
}
