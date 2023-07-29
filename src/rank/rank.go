package rank

import "github.com/gocolly/colly/v2"

func RankPage(response *colly.Response) int {
	return 17
	//return int(response.Body[12]) //something random (hopefully not OOB)
}
