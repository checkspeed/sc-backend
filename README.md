# sc-backend

## Speed Test API Documentation

### Overview

### Repository
The repositories package encapsulates the logic required to interact with the database, allowing other parts of the application to perform CRUD operations without directly dealing with SQL queries or database connections.

#### Key Components

- **`Datastore Interface`**: Defines the contract for database operations, including creating speed test results, fetching them based on filters, and closing database connections.
- **`store Struct`**: Implements the Datastore interface, providing concrete methods to interact with the database using GORM.
- **`GetSpeedTestResultsFilter`**: A struct used to filter speed test results by country code.

#### Functions

**NewStore**

Creates a new instance of the store struct, initializing it with a connection to the specified database URL.

```go
func NewStore(dbUrl string) (store, error)
```

**RunAutoMigrate**
Automatically migrates the database schema to match the current model definitions.

```go
func RunAutoMigrate() (store, error)
```

**RunDropTable**
Drops existing tables defined in the models, effectively resetting the database schema.

```go
func (s *store) RunDropTable() error
```

**CloseConn**
Closes the underlying database connection.

```go
func (s store) CloseConn(ctx context.Context) error
```

**CreateSpeedtestResult**

Creates a new speed test result entry in the database.

```Go
func (s store) CreateSpeedtestResult(ctx context.Context, speedTestResult *models.SpeedTestResult) error

```

**GetSpeedTestResults**
Retrieves speed test results from the database, optionally filtered by a country code

```go
func (s store) GetSpeedTestResults(ctx context.Context, filters GetSpeedTestResultsFilter) ([]models.SpeedTestResult, error)
```

### Controllers
The package contains the core logic for handling HTTP requests and responses within the application.

#### Network Information Controller

**Purpose**  
Fetches and returns network-related information for a specific IP address, including geographical location data.

**Input Parameters**  
- **IP Address**: The IP address for which to fetch network information. Provided as a query parameter in the request URL.
- **Geo API key**: The API key for accessing geographical data.

#### CreateSpeedtestResults

**Purpose**  
Creates a new speed test result entry in the database.

**Input Parameters**  
- A JSON payload containing the speed test result data.

**Output**  
- A JSON response indicating success or failure in creating the speed test result.

#### GetSpeedtestResults

**Purpose**  
Retrieves speed test results from the database, optionally filtered by specific criteria.

**Input Parameters**  
- A JSON payload containing filter criteria for the speed test results.

**Output**  
- A JSON response containing the retrieved speed test results or an error message.

### Routes

**POST /speed_test_result/list/list**  
This endpoint retrieves speed test results based on the provided filter criteria in the request body.

```go
	r.POST("/speed_test_result/list", ctrl.GetSpeedtestResults)

```

**Request Body**
The request body should be a JSON object containing the filter criteria.

Example Request Body

```json
{
  "country_code": ""
}
```

**POST /speed_test_result/list**
This endpoint to create a speed test result

```Go
	r.POST("/speed_test_result", ctrl.CreateSpeedtestResults)
```

**Get /network**
This endpoint is to get network information based on the IP address.

````Go
	r.GET("/network", ctrl.GetNetworkInfo)
  ```
````
