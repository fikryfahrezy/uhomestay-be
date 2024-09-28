package history

type HistoryDeps struct {
	HistoryRepository *HistoryRepository
}

func NewDeps(
	historyRepository *HistoryRepository,
) *HistoryDeps {
	return &HistoryDeps{
		HistoryRepository: historyRepository,
	}
}
