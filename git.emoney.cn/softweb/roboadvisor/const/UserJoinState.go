package _const

const (
	//用户未参与活动
	UserJoinState_NotJoin = 1
	//用户已参与活动，但未开奖
	UserJoinState_Join = 2
	//用户参与活动但未中奖
	UserJoinState_Join_Not_Guess = 3
	//用户参与活动并中奖
	UserJoinState_Join_Guessed = 4
)
