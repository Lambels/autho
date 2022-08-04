package facebook

// User represents fields accessible on public facebook accounts.
//
// https://developers.facebook.com/docs/graph-api/reference/user/#default-public-profile-fields
type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}
