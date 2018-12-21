package viewmodel

//我的参与记录与奖品记录-猜涨跌
type GuessChangeJoinInfo struct {
	//猜涨跌参与记录
	GuessChangeResultList []*GuessChangeResult
	//猜涨跌奖品记录
	GuessChangeAwardList []*GuessChangeAward
	//本周参与次数
	JoinCount int
	//猜中次数
	GuessedCount int
}
