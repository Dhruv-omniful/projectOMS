# GoCommons Onboarding Assignment: OMS & IMS

This repository contains the source code for the GoCommons Onboarding Assignment. It implements two microservices, an **Order Management Service (OMS)** and an **Inventory Management Service (IMS)**, built using Go and leveraging various GoCommons libraries and patterns.

## : Microservices Overview

### : Order Management Service (OMS)

The OMS is responsible for handling all aspects of customer orders.

- **Order Creation via CSV**: Orders can be created in bulk by uploading a CSV file. A REST endpoint accepts a file path (S3 upload).
- **Asynchronous Processing**: The service pushes a message to an SQS queue (`CreateBulkOrder`) to trigger asynchronous processing of the CSV file.
- **CSV Processor**: A dedicated worker consumes from SQS, downloads the file from S3, and parses it.
- **Data Validation**: It validates SKUs and Hubs by calling the IMS APIs.
- **Order Persistence**: Valid orders are saved to a MongoDB collection with an `on_hold` status. Invalid rows are logged and made available for download via a public endpoint.
- **Event-Driven Finalization**: Upon successful creation, an `order.created` event is published to a Kafka topic.
- **Inventory Check**: A Kafka consumer listens for `order.created` events and checks for inventory availability via an IMS call.
- **Atomic Updates**: If inventory is available, the order status is updated to `new_order`, and inventory is reduced in an atomic transaction. Otherwise, the order remains `on_hold`.
- **Public APIs**:
  - Filtered list of orders (by tenant, seller, status, date).
- **Webhooks**: Provides a mechanism for other services to register webhooks and receive push notifications for order updates, with per-tenant filtering.
- **Internationalization**: User-facing error messages support i18n.

### : Inventory Management Service (IMS)

The IMS is the source of truth for product and stock information.

- **Hub Management**: Provides CRUD APIs for managing hubs (warehouses/locations).
- **SKU Management**: Provides CRUD APIs for managing SKUs (products). Includes filtering by tenant, seller, and SKU codes.
- **Inventory Management**:
  - **Atomic Upserts**: An endpoint for atomically upserting inventory levels.
  - **Inventory View**: An endpoint to view current inventory for a given hub and a list of SKUs. Missing entries default to `0.
- **Caching**: Uses Redis to cache SKU and hub validation responses to improve performance.
- # OMS & IMS API Overview

This document provides a detailed overview of the APIs for the Order Management Service (OMS) and Inventory Management Service (IMS).

---

### 1. Order Management Service (OMS)

The OMS handles order intake, processing, and fulfillment orchestration.

**Order Creation via CSV**
- **Endpoint**: `POST /orders/csv`
- **Alternative for Local Testing**: `POST /orders/upload-local` with a JSON body `{"path":"path/to/local.csv"}`.
- **Process**:
  1. Accepts a CSV file containing bulk order data.
  2. Pushes a message to the `CreateBulkOrder` SQS queue to trigger asynchronous processing.

**CSV Processor (SQS Consumer)**
- **Trigger**: New message in the `CreateBulkOrder` SQS queue.
- **Process**:
  1. Downloads the corresponding CSV file from S3.
  2. Parses the CSV row by row.
  3. **Validation**:
     - Calls IMS endpoints (`GET /inventory/query`) to validate that the Hub and SKU for each order row exist and are valid. These calls are cached in Redis to improve performance.
  4. **Outcome**:
     - **Valid Rows**: Saved as individual orders to the `orders` collection in MongoDB with an `on_hold` status. An `order.created` event is then published to Kafka for each valid order.
     - **Invalid Rows**: Logged to a separate error file, which is made available for download.

**Order Finalizer (Kafka Consumer)**
- **Trigger**: `order.created` event on the Kafka topic.
- **Process**:
  1. Checks current inventory levels by calling the IMS `GET /inventory/query` endpoint.
  2. **If sufficient inventory exists**:
     - Atomically decrements the stock by calling the IMS `POST /inventory/consume` endpoint.
     - Updates the order status in MongoDB from `on_hold` to `new_order`.
     - Publishes an `order.updated` event to Kafka.
  3. **If inventory is insufficient**:
     - The order remains in the `on_hold` status for a future retry or manual intervention.

**Webhook Dispatcher**
- **Trigger**: Listens for `order.created` and `order.updated` Kafka events.
- **Process**:
  1. Finds all registered webhooks for the specific tenant associated with the order.
  2. Delivers the full order payload to the registered callback URLs.
  3. Implements a retry mechanism for transient delivery failures.

**Public REST APIs**
  - `GET /orders`: Retrieves a paginated and filtered list of orders. Supports filtering by `tenant_id`, `seller_id`, `status`, and a date range.
  - `POST /orders`: Creates a single order. Performs the same validations as the bulk process and emits an `order.created` event.
  - `POST /orders/csv`: Kicks off the bulk order creation process.
  - `POST /orders/upload-local`: For local testing of the bulk order process.
  - `GET /orders/errors/:file_id`: Downloads the CSV file containing invalid rows from a specific bulk upload.
  - `POST /webhooks`: Registers a new webhook URL for a tenant to receive order event notifications.
  - `GET /webhooks`: Lists all webhooks for a tenant.
  - `PUT /webhooks/:id`: Updates an existing webhook.

---

### 2. Inventory Management Service (IMS)

IMS is the source of truth for all data related to tenants, sellers, hubs, SKUs, and inventory levels.

**Entity CRUD APIs**
- **Tenants**: Full CRUD at `/tenants`
- **Sellers**: Full CRUD at `/sellers`
- **Hubs (Warehouses)**: Full CRUD at `/hubs`. Includes lookup by code at `/hubs/code/:hub_code`. Hub data is cached in Redis.
- **SKUs (Products)**: Full CRUD at `/skus`. Includes lookup by code at `/skus/code/:sku_code`. SKU data is cached in Redis.

**Inventory APIs**
- `POST /inventory`: Atomically creates or updates (upserts) the inventory quantity for a given SKU at a specific hub.
- `PUT /inventory/:id`: Updates a specific inventory record.
- `POST /inventory/consume`: Atomically decrements stock for a given SKU and hub. Used by OMS during order finalization.
- `GET /inventory/query`: Returns the current inventory levels for a list of SKUs at a specific hub. If an entry for a Hub/SKU combination doesn't exist, it defaults to a quantity of `0`.
- `GET /inventory`: Lists all inventory records (paginated).

---


## : Tech Stack

- **Backend**: Go
- **Framework/Libraries**: GoCommons for server setup, Kafka, Redis, logging, and error handling.
- **Databases**:
  - **PostgreSQL**: Primary database for IMS (hubs, SKUs, inventory).
  - **MongoDB**: Database for OMS (orders).
  - **Redis**: Caching layer for IMS.
- **Messaging**:
  - **Apache Kafka**: For event-driven communication between services (`order.created`).
  - **AWS SQS**: For queueing background jobs (CSV processing).
- **APIs**: REST
- **Infrastructure**:
  - **Docker & Docker Compose**: For containerizing and orchestrating services.
  - **AWS S3**: For storing uploaded CSV files (emulated with LocalStack).
  **Logging**: *Request/response logging via GoCommons http.RequestLogMiddleware*

## : Folder Structure

The repository is structured as a monorepo with each microservice in its own directory:

```
omni_project/
├── ims/                  # Inventory Management Service
│   ├── configs/
│   ├── controllers/
│   ├── model/
│   ├── postgres/
│   ├── go.mod
│   └── main.go
├── oms/                  # Order Management Service
│   ├── api/
│   ├── client/
│   ├── configs/
│   ├── model/
│   ├── services/
│   ├── worker/
│   ├── go.mod
│   └── main.go
└── docker-compose.yaml   # Docker orchestration for dependencies
```
## : Screenshots
--Starting IMS on PORT 8081
![Image](https://github.com/user-attachments/assets/6b8dc143-dec9-4e2a-9268-344cf6c03751)
--OMS
![Screenshot 2025-06-24 163909](https://github.com/user-attachments/assets/1a19a732-3be2-4857-b3b2-770053f57279)
--Uploading CSV to s3
![Image](https://github.com/user-attachments/assets/1dc08037-cf2f-4b0b-a3ba-5e963f58465e)
--valid row(call ims api and then check hubcode and skucode, also check quantity in csv>0)
![Image](https://github.com/user-attachments/assets/a6dd4d47-b327-4a80-9314-16b0541b80b7)
--invalid row(if either hub/skucode mismatched or invalid inventory then move to s3://oms-bucket/errors)
![Image](https://github.com/user-attachments/assets/93703091-9703-41ed-9395-9c190c8309dc)
![Screenshot 2025-06-24 164804](https://github.com/user-attachments/assets/71fba16e-b1f4-45a5-8ba7-0a4541598514)
--inventory update
![Screenshot 2025-06-24 164557](https://github.com/user-attachments/assets/ab50bbd7-2db9-447b-a45f-e1b83134b430)
--mongo DB

![Screenshot 2025-06-24 164225](https://github.com/user-attachments/assets/fff31cbc-7349-4092-a56f-3f0bfefc1669)
![Screenshot 2025-06-24 164349](https://github.com/user-attachments/assets/b82c23e6-ee42-47ce-9b36-a721216d6059)
![Screenshot 2025-06-24 164928](https://github.com/user-attachments/assets/a3f2a9e4-e7f3-4b33-93d8-2cd32fc46758)

--webhook

![Image](https://github.com/user-attachments/assets/1cc4025d-a2e9-43d7-a6b1-571643707188)

## : Getting Started

Follow these steps to get the project up and running locally.

### Prerequisites

- [Go](https://golang.org/doc/install) (latest version recommended)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- PostgreSQL — A local PostgreSQL instance is used as the primary database in this project.
- go_commons — Common utilities and libraries required by this project. Ensure it is set up and properly integrated.

### 1. Clone the Repository and Navigate to the Project Directory

First, clone this repository to your local machine.

```sh
git clone https://github.com/Dhruv-omniful/projectOMS.git
powershell
cd omni_project
```

### 2. Configure IMS Database

The IMS service requires a PostgreSQL database. The provided `docker-compose.yaml` does not include a Postgres instance. You can add one by appending the following service definition to your `docker-compose.yaml` file:

```yaml
# In docker-compose.yaml, add this service alongside the others:

  postgres:
    image: postgres:14
    container_name: postgres
    environment:
      POSTGRES_USER: "postgres"
      POSTGRES_PASSWORD: "root"
      POSTGRES_DB: "ims_db"
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

# At the end of docker-compose.yaml, add this:
volumes:
  postgres_data:
```

### 3. Start All Docker Containers

```powershell
docker-compose up -d
```

### 4. Verify Containers Are Running

```powershell
docker ps
```
You should see `kafka`, `zookeeper`, `redis`, `mongo`, `localstack`, and `postgres` running.

### 5. Set AWS Environment Variables for LocalStack

This step configures your current PowerShell session to use LocalStack.
```powershell
$env:AWS_ACCESS_KEY_ID = "test"
$env:AWS_SECRET_ACCESS_KEY = "test"
$env:AWS_REGION = "us-east-1"
$env:LOCAL_SQS_ENDPOINT = "http://localhost:4566"
```

### 6. Create S3 Bucket in LocalStack

```powershell
aws --endpoint-url=http://localhost:4566 s3 mb s3://oms-bucket --region us-east-1
```

### 7. Create SQS Queue in LocalStack

```powershell
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name CreateBulkOrder --region us-east-1
```

### 8. Start IMS (Inventory Management Service)

 Run this in a **new terminal** and keep it running.

```powershell
cd ./ims
$env:CONFIG_SOURCE = "local"
go run main.go
```


### 9. Start OMS (Order Management Service)

 Run this in another **new terminal** and keep it running.

```powershell
# Navigate to the project directory first
cd omni_project
# Set environment variables and run the service
$env:CONFIG_SOURCE = "local"
$env:AWS_ACCESS_KEY_ID = "test"
$env:AWS_SECRET_ACCESS_KEY = "test"
$env:AWS_REGION = "us-east-1"
$env:LOCAL_SQS_ENDPOINT = "http://localhost:4566"
cd ./oms
go run main.go
```

### 10. Upload Local CSV File to Trigger Order Processing

```powershell
Invoke-WebRequest `
  -Uri http://localhost:8080/orders/upload-local `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"path":"csv/sample.csv"}'
```

### 11. View Kafka Topic (`order.created`)

```powershell
docker exec -it kafka /bin/bash
```
Then, inside the Kafka container's shell, run:
```bash
kafka-console-consumer --bootstrap-server localhost:9092 --topic order.created --from-beginning
```

### 12. Check Orders in MongoDB

```powershell
docker exec -it mongo mongosh
```
Then, inside the Mongo shell, run:
```javascript
use oms_db
db.orders.find().pretty()
```

### 13. (Optional) Post Inventory to IMS

```powershell
Invoke-WebRequest -Uri http://localhost:8081/inventory `
  -Method POST `
  -ContentType "application/json" `
  -Body '{
    "tenant_id": "t1",
    "seller_id": "s1",
    "hub_code": "H2",
    "sku_code": "SKUU-001",
    "quantity": 5
  }'
```

### 14.(Optional) Register Webhook (`order.updated`)

```powershell
$body = @{
    tenant_id = "t1"
    callback_url = "https://webhook.site/53931814-d7a8-497c-b222-f1df4e2ca484"
    events = @("order.updated")
    headers = @{}
    secret = "mysecret"
    is_active = $true
} | ConvertTo-Json -Depth 3

Invoke-WebRequest -Uri http://localhost:8080/webhooks `
  -Method POST `
  -ContentType "application/json" `
  -Body $body
``` 
