package models

func GetNotificationsForUser(user *User) ([]Notification, error) {
	// Return notifications user is authorized to view:
	// - All
	// - Owners if user is owner
	// - Board if user has a board role
	// - Specific roles
	// Implement Datastore query based on that logic
}

func SaveNotification(n *Notification) error {
	return saveToDatastore(n, "notifications", n.ID.String())
}

func ResolveRecipients(n Notification) struct{ Emails, Phones []string } {
	// Lookup users based on audience/roles
	// Return their verified emails and phone numbers
}
