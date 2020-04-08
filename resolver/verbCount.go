package resolver

import "insultService/model"

type verbCountResolver struct {
	VerbCount model.VerbCount
}

func (v *verbCountResolver) Word() string {
	return v.VerbCount.Verb
}

func (v *verbCountResolver) Count() int32 {
	return int32(v.VerbCount.Count)
}
