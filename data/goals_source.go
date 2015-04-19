package data

import (
	"database/sql"
	"time"
)

type GoalSrc struct {
	db *sql.DB
}
type Goal struct {
	ID         int       `json: "id"`
	Desc       string    `json: "desc"`
	StartDate  time.Time `json: "start_date"`
	Weeks      int       `json: "weeks"`
	UserID     int       `json: "user_id"`
	LoseCommit string    `json: "lose_commit"`
}

type Goals []Goal

var goalSource = GoalSrc{db: GetDBInstance()}

func GetGoalsOfUser(userID int) Goals {
	return nil
}

func GetGoal(goalID int) *Goal {
	return nil
}

func GetReferee(goalID int) *User {
	return nil
}

//UpdateGoal only update the content of goal - not the id
func UpdateGoal(goalID int, goal *Goal) error {
	return NotImplementedError{}
}

//CreateGoal creates a goal in data store, and return the newly created goal
// otherwise: return non-nil error: (TODO: define error in more detail in "real" version)
func CreateGoal(goal *Goal) (*Goal, error) {
	return nil, NotImplementedError{}
}
