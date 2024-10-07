# Go-Decision-Tree

A decision tree implementation in Go.

Note that this project is only a Course Assignment project, and it is not recommended to use this in production.

## How to build from source

To run this code, you need to install golang SDK version >= 1.23.2.

After you get your golang SDK installed, just run the following command to build the project:

```bash
make all
```

To build the decision tree using preset dataset, run the following command:

```bash
go run main.go
```

And run the test using the following command:

```bash
go test predict_test.go
```

Accuracy on this dataset (dataset has been resampled to balance the class data):

```text
=========================== TRAIN DATASET ===========================
Accuracy: 80.94%
Pessimistic error: 20.23%
Class Data [<=50K] frequency: 51.24%
Within class [<=50K] predict accuracy: 73.73%
Class Data [>50K] frequency: 48.76%
Within class [>50K] predict accuracy: 88.52%
=========================== TEST DATASET ===========================
Accuracy: 78.08%
Pessimistic error: 24.26%
Class Data [<=50K] frequency: 51.87%
Within class [<=50K] predict accuracy: 72.80%
Class Data [>50K] frequency: 48.13%
Within class [>50K] predict accuracy: 83.78%
```

# Basic Usages

## Loading Dataset

A dataset should at least consists of 2 parts: Names and Data.

A names file should be like this:

```text
| This is a comment

| Class definition must be the first attribute to be defined
| Class must be a nominal attribute.
Class Name: Class A, Class B.
| You can also make the class anonamous:
| Class A, Class B.
| By doing this, the class will be automatically named as "Class".

| For attribute definition, we have 2 types: continuous and nominal.
| An example of continuous attribute definition:
Attr1: continuous.

| An example of nominal attribute definition:
Attr2: Value1, Value2, Value3.

| Note that the ordinal attribute is not supported in this implementation.
| If you really need an ordinal attribute, you can convert it to a continuous attribute.
```

A data file should be like this:

```text
| This is a comment

| According to the definition above, the data line (or we call it an "instance") should be like this:
| Attr1, Attr2, Class.
1.5, Value1, Class A.
1.8, Value3, Class B.

| For missing value, just replace it with a question mark "?".
4.5, ?, Class B.
```

To load a dataset from file:
```go
attrTable, err := data.ReadAttributes(attributesFile)
if err != nil {
    log.Fatalf("failed to read attributes: %v", err)
    return
}

trainData, err := data.ReadValues(config.Conf, attrTable, trainDataFile)
if err != nil {
    log.Fatalf("failed to read training data: %v", err)
    return
}
```

## Building Decision Tree

To build a decision tree, you can use the following code:

```go
t, err := tree.BuildTree(config.Conf, trainData)
if err != nil {
    log.Fatalf("failed to build tree: %v", err)
    return
}
```

The tree building process consists of following steps:
1. Data washing: Remove instances with missing class values.
2. Node building: Build nodes by splitting nodes based on Entropy:
   1. For continuous attribute, we support binary split.
   2. For nominal attribute, we support multi-way split and binary split.
3. Post-Pruning: Prune the tree to avoid overfitting.

After these processes, the returned object `t` is a decision tree. You can either save the tree into json format, or use it to predict.

## Predicting

To predict a value, you can use the following code:

```go
predicted, err := t.Predict(dataInstance)
if err != nil {
    log.Fatalf("failed to predict: %v", err)
    return
}
```

Return value is of type `string`, indicating the value of class prediction.

## Serialize / Deserialize

You can read your tree from a json file, or save your tree to a json file.

To save a tree to a json file:

```go
err = tree.WriteTreeToFile(t, "tree.json")
if err != nil {
    log.Fatalf("failed to save tree: %v", err)
    return
}
```

To read a tree from a json file:

```go
tr, err := tree.ReadTreeFromFile("tree.json")
if err != nil {
    log.Fatalf("failed to read tree: %v", err)
    return
}
```

## Testing

To test the tree, you can use the following code:

```go
res, err = tree.TestRun(tr, testData)
if err != nil {
    t.Fatalf("failed to do test run: %v", err)
    return
}
t.Logf("Accuracy: %.2f%%", res.Accuracy*100)
t.Logf("Pessimistic error: %.2f%%", res.PessimisticError*100)
for class, count := range res.ClassDataCount {
    t.Logf("Class Data [%s] frequency: %.2f%%", class, float64(count)/float64(len(testData.Instances))*100)
    t.Logf("Within class [%s] predict accuracy: %.2f%%", class, float64(res.ClassCorrectCount[class])/float64(count)*100)
}
```

The return value is of type `tree.TestResult`, which contains the following fields:
1. Correct count / error count / accuracy: As the name shows.
2. PessimisticError: The pessimistic error of the prediction ($PessimisticError = (N_{TrainPredictErr} + 0.5 * N_{leafNodes}) / N_{trainInstances}$).
3. Class value data count: The count of each class value in the test data.
4. Class correct count / error count / accuracy: The metrics of each class value in the test data.

# License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
