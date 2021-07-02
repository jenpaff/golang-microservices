# Validation

## Motivation
- [ ] validate format to avoid errors early rather than failing on a db level
- [ ] validate any input we receive via POST requests to avoid SQL injection
- [ ] having tested a couple of libraries we switched to go-validate when we were forced of adding custom regex pattern

## Usage

I don't even want to attempt to write a better documentation than the go validate package doc. 
- [ ] link to the documentation  

### Using validation

Looking at our `POST /users` request, we want to validate the request body before we process any received request. 
To define our input valdiation we write our input struct like this:

```go
type UserCreationRequest struct {
	UserName    string `json:"name" validate:"notBlank,validRegexInput"`
	Email       string `json:"email" validate:"notBlank"`
	PhoneNumber string `json:"phone_number" validate:"notBlank"`
}
```

### Implementing a validation


## Code

### How the validation library is implemented


## Further Resources
