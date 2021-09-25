package config

type Labels = *LabelArray
type LabelArray []Label

func (x *LabelArray) GetByName(name string) (*Label, int) {
	for i, label := range *x {
		if label.GetName() == name {
			return &label, i
		}
	}
	return nil, -1
}

func (x *LabelArray) Add(item Label) {
	tmp := *x
	*x = append(tmp, item)
}

func (x *LabelArray) RemoveByIndex(i int) {
	tmp := *x
	*x = append(tmp[:i], tmp[i+1:]...)
}
