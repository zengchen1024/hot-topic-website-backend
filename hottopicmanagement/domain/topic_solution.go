package domain

type DiscussionSourceSolution struct {
	ResolvedOne int   // id of resolved discussion source
	RelatedOnes []int // id sets of related open discussion source
}

type TopicSolution struct {
	TopicId   string
	Solutions []DiscussionSourceSolution
}
