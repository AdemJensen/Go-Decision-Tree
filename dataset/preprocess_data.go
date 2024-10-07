package dataset

import "DecisionTree/data"

func PreProcessData(valueTable *data.ValueTable) {
	// 75% is <=50K, 25% is >50K
	// Resample the data to make it balanced
	dataInstances := make(map[string][]*data.Instance)
	for _, attr := range valueTable.Instances {
		classValue := attr.ClassValue.Value().(string)
		dataInstances[classValue] = append(dataInstances[classValue], attr)
	}
	valueTable.Instances = append(valueTable.Instances, dataInstances[">50K"]...)
	valueTable.Instances = append(valueTable.Instances, dataInstances[">50K"]...)

	// education and education-num are the same
	// remove education-num
	valueTable.RemoveAttribute("education-num")
}
