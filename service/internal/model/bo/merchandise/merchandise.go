package boMerchandise

type GetRecommendationArgs struct{}
type GetRecommendationReply struct {
	Recommendation string `json:"recommendation"`
}
