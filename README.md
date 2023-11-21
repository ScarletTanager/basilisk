# basilisk

Basilisk is a (very rudimentary) machine learning toolkit written in Go.  I started it to explore ML from the standpoint not just of using it but also from implementing it - getting into the guts of how the models and algorithms work (together and separately).

## Components

Basilisk includes three main components:

- the [main library](#models), which at present consists of the underlying representations of a dataset and a K-Nearest Neighbors classifier;
- the [dataset generation](#dataset-generation) library and executable, which provide for the configurable generation of synthetic datasets; and
- the [basilisk server](#basilisk-server), which is a simple REST-based HTTP server you can use to try out the models, and which uses both the main and dataset generation libraries.

### Models

### Dataset Generation

You can install the `dsgenerate` utility for creating synthetic datasets with:

```sh
$ go install github.com/ScarletTanager/basilisk/dsgen/dsgenerate
```

To define a new dataset, you:
1. Define a dataset configuration in JSON.  Let's assume you save this file as `./dsconf.json`.
2. Run `dsgenerate` to create the dataset and save it (in JSON) as `./dataset.json`: `dsgenerate -config ./dsconf.json -output ./dataset.json -format json`

#### Dataset configuration

To define a synthetic dataset to be used for testing a classification model, you create a JSON configuration file.  The general structure of the file is as follows:

```json
{
  "recordCount": "number",
  "classes": {
    "<classname1": [
      {
        "name": "string",
        "lower": "number",
        "upper": "number",
        "allocationsByQuintile": [
          "number",
          "number",
          "number",
          "number",
          "number",
        ]
      },
      {...}
    ],
    "classnameN": [...] 
  }
}
```

The fields should be populated as follows:
- **recordCount** - A positive integer indicating the number of record to create in the dataset
- **classes** - A map of **class** objects - the keys are the names of the classes
  - Each **class** is an array of **attribute** objects:
    - **name** - The name of the attribute
    - **lower** - The lower bound on values for the attribute.  `dsgenerate` _currently_ treats all attribute values as `float64` values internally.
    - **upper** - The upper bound on allowable values for this attribute
    - **allocationsByQuintile** - An array of values which define the attribute's value distribution (by percentage) within its allowable range.  _These values must add to  100._  Although the generator will accept values which total to less than 100, this can produce weird behavior, and I don't recommend it.  Fixing that is not an incredibly high priority, so just make sure these add to 100.  _There must be five (5) values in the array, one for each quintile._  Be careful that you treat these as percentages, not decimal fractions - so one-fifth should be `20` or `20.0`, _NOT_ `.2`.

Some rules about the configuration:

1. All classes must have the exact same number of attributes.
1. Attribute names must match exactly across classes.
1. Class names must be unique within your dataset.
1. The lower bound must be less than the upper bound.

There is a sample configuration in this repository at `dsgen/dsgenerate/sampleconf.json`.

#### Dataset format

The generation code (which is the same whether you use `dsgenerate` or the REST API in the `basilisk` server) can output either JSON or CSV (`dsgenerate`) or JSON (`basilisk` server).  Examples of the two formats can be found at `datasets/shorebirds.json` and `datasets/shorebirds.csv`.  The model training API in the REST server only accepts datasets in JSON.

The CSV should work with standard data toolkits (e.g. `scikit-learn`), but if not, please open an issue.

### Basilisk server

To run the `basilisk` REST API server, simply enter:

```sh
$ go run basilisk/basilisk.go
```

#### Supported endpoints and operations

At present the REST server provides the following endpoints:

- `/datasets`
  - `POST` - creates a new synthetic dataset.  The configuration (see [above](#dataset-configuration) for format) is passed as the JSON body, the response is the dataset in JSON.  To generate a dataset as CSV, use the `dsgenerate` command.
- `/models`
  - `GET` - lists the currently running models and their configurations
  - `POST` - creates a new model, which exists for the lifetime of the server.  There is no persistence layer in `basilisk`, so if you wanted to store a trained model over restarts, you'd need to add your own persistence layer.  The payload is pretty simple at present - `{"K": <int>, "distance_method": <string>}`.  `K` must be a positive integer (the only supported model right now is `KNearestNeighbors`), and `distance_method` must be one of `euclidean` or `manhattan`.  Euclidean distance is the magnitude of two difference of the two vectors, Manhattan (or "city block") distance is the sum of the differences of each vector component.  In two dimensions, this can be visualized as the distance along the rectilinear grid lines, thus the "city block" moniker.
- `/models/:id/data`
  - `PUT` - trains the specified model.  The dataset to be used is passed as the JSON body, the server will (currently) do a randomized split, with 75% of the records being allocated as training data, and the other 25% being used as testing data.
- `models/:id/results`
  - `GET` - tests the specified model and returns the test results analysis
