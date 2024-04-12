package db

import "minitwit-api/model"

type Idb interface {
	Connect_db()
	QueryMessage(message *model.Message)
	QueryFollow(args []int)
	QueryUnfollow(args []int)
	QueryDelete(args []int)
	QueryRegister(args []string)
	GetMessages(args []int) []map[string]any
	GetMessagesForUser(args []int) []map[string]any
	GetFollowees(args []int) []string
	Get_user_id(username string) (int, error)
	IsNil(i interface{}) bool
	IsZero(i int) bool

	GetAllUsers() []model.User
	GetAllFollowers() []model.Follower
	GetAllMessages() []model.Message
}

/*
docker run \
    --rm \
    -e SONAR_HOST_URL="http://localhost:9000" \
    -e SONAR_SCANNER_OPTS="-Dsonar.projectKey=minitwit" \
    -e SONAR_TOKEN="sqp_023e532decb800cad5951d606c213eef41b33ab7" \
    sonarsource/sonar-scanner-cli
*/
