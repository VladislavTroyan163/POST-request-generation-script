package utils

func GetWorkersCount(configWorkersCount int, lenRegistry int) (workersCount int) {
	workersCount = configWorkersCount
	if lenRegistry < workersCount {
		workersCount = lenRegistry
	}
	return
}
