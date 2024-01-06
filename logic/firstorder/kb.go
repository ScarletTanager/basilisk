package firstorder

//
// kb.go - KnowledgeBase-related code
//

type KnowledgeBase interface {
	Tell()
}

type knowledgeBase struct {
}
