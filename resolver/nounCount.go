package resolver

import "insultService/model"

type nounCountResolver struct {
	NounCount model.NounCount
}

func (n *nounCountResolver) Word() string {
	return n.NounCount.Noun
}

func (n *nounCountResolver) Count() int32 {
	return int32(n.NounCount.Count)
}
