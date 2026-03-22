package jira

type Issue struct {
	ID     string
	Key    string
	Fields IssueFields
}

type IssueFields struct {
	Summary     string
	Description string
	Status      IssueStatus
	Assignee    *User
	Reporter    *User
	Parent      *IssueParent
}

type IssueParent struct {
	Key string
}

type IssueStatus struct {
	Name string
}

type User struct {
	DisplayName  string
	EmailAddress string
}

type Transition struct {
	ID   string
	Name string
	To   TransitionTo
}

type TransitionTo struct {
	Name string
}
