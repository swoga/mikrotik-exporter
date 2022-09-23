package config

type Metrics = *MetricArray
type MetricArray []Metric

func (x *MetricArray) GetByName(name string) (*Metric, int) {
	for i, metric := range *x {
		if metric.MetricName == name {
			return &metric, i
		}
	}
	return nil, -1
}

func (x *MetricArray) Add(item Metric) {
	tmp := *x
	*x = append(tmp, item)
}

func (x *MetricArray) RemoveByIndex(i int) {
	tmp := *x
	*x = append(tmp[:i], tmp[i+1:]...)
}
