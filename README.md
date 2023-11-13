# basilisk

Basilisk is a (very rudimentary) machine learning toolkit written in Go.  I started it to explore ML from the standpoint not just of using it but also from implementing it - getting into the guts of how the models and algorithms work (together and separately).

## Components

Basilisk includes three main components:

- the [main library](#models), which at present consists of the underlying representations of a dataset and a K-Nearest Neighbors classifier;
- the [dataset generation](#dataset-generation) library and executable, which provide for the configurable generation of synthetic datasets; and
- the [basilisk server](#basilisk-server), which is a simple REST-based HTTP server you can use to try out the models, and which uses both the main and dataset generation libraries.

### Models

### Dataset Generation

### Basilisk server

#### Supported endpoints and operations

- `/datasets`
  - `POST` - creates a new synthetic dataset.  The configuration is passed as the JSON body, the response is the dataset in JSON.  To generate a dataset as CSV, use the `dsgenerate` command.
- `/models`
  - `POST` - creates a new model, which exists for the lifetime of the server.  There is no persistence layer in `basilisk`, so if you wanted to store a trained model over restarts, you'd need to add your own persistence.
- `/models/:id`
  - `PUT` - trains the specified model.  The dataset to be used is passed as the JSON body, the server will (currently) do a randomized split, with 75% of the records being allocated as training data, and the other 25% being used as testing data.
