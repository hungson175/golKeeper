package data

import (
	"database/sql"
	"time"

	"github.com/hungson175/golKeeper/utils"
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
	RefID      int       `json: "ref_id"`
	RefName    string    `json: "ref_name"`
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
	return goalSource.createGoal(goal)
}

/*
Private area: implement CRUD
*/

//Create
func (src *GoalSrc) createGoal(goal *Goal) (*Goal, error) {
	tx, err := src.db.Begin()

	strDate := utils.Date2String(&goal.StartDate) //TODO: in production code, this will be tested / assert merciless
	// query := fmt.Sprintf("insert into goals (desc,start_date,weeks,user_id,ref_id,lose_commit) Values (%v,%v,%v,%v,%v,%v)",
	// 	goal.Desc,
	// 	strDate,
	// 	goal.Weeks,
	// 	goal.UserID,
	// 	goal.RefID,
	// 	goal.LoseCommit)
	// fmt.Println("Query: ", query)
	// x := 1
	// y := 1
	// if x == y {
	// 	panic(errors.New("Just test"))
	// }
	result, err := tx.Exec("insert into goals (description,start_date,weeks,user_id,ref_id,lose_commit) Values (?,?,?,?,?,?)",
		goal.Desc,
		strDate,
		goal.Weeks,
		goal.UserID,
		goal.RefID,
		goal.LoseCommit)
	var newGoal Goal
	if err != nil {
		tx.Rollback()
		return nil, err
	} else {
		id64, _ := result.LastInsertId()
		newGoal = *goal
		newGoal.ID = int(id64)
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	} else {
		tx.Commit()
		return &newGoal, nil
	}
}

//insert even the id

//Read
//Update
//Delete
