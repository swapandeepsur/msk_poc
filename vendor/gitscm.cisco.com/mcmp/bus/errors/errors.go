/*
Package errors defines common event bus errors.
*/
package errors

// ConfigurationError is the type of error returned from a constructor (e.g. NewProducer, or NewConsumer)
// when the specified configuration is invalid.
type ConfigurationError string

func (err ConfigurationError) Error() string {
	return "invalid configuration (" + string(err) + ")"
}
