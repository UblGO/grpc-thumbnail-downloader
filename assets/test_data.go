package testData

var (
	ValidURLs = []string{
		"https://www.youtube.com/watch?v=y0sF5xhGreA",
		"https://www.youtube.com/watch?v=MlDtL2hIj-Q",
		"https://www.youtube.com/watch?v=DixdzXAFS18",
		"https://www.youtube.com/watch?v=KvE92fCMbmc"}
	InvalidURLs = []string{
		"https://www.youtube.com/watch?v=y0shGreA",
		"https://www.google.com",
		"https://dzen.ru",
		"127.0.0.1"}
	BrokenLinks = []string{
		"https://www.youtub0shGreA",
		"https://wwwfeoogle.com",
		"syoutube.com/watch?v=KvE92fCMbmc",
		"d,.dwd.;sadw",
		""}
)
