// Interface definition

package emailer

type Mailer interface {
	// OTP token or String of Pass code : SendOTP
	// If pass code is empty , it will auto-generate it
	SendOTP(identifier string , token string ) (string, error)

	// Welcome message can either be sent to the user 
	// through and email or phone number
	// TODO:: Whats Welcome message
	SendWelcomeMessage(identifier string) (string, error)
}