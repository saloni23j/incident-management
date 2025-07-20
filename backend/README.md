# Incident Management API

A simple incident management system with AI-powered severity and category classification.

## Features

- **POST /api/v1/incidents** - Create a new incident with AI analysis
- **GET /api/v1/incidents** - Get all incidents
- **AI Integration** - Automatically determines severity (low/medium/high) and category (network/software/hardware/security)
- **Comprehensive Validation** - Input validation with detailed error messages
- **Simple & Clean** - Single model approach with JSON, GORM, and validation tags

## Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Set OpenAI API key (optional):**
   ```bash
   export OPENAI_API_KEY="your-openai-api-key-here"
   ```
   If no API key is set, the system will use default values (medium severity, software category).

3. **Run the application:**
   ```bash
   go run main.go
   ```

## 🛠️ Setup Instructions

### 🔧 Backend (Go)

1. **Clone and install dependencies**
   ```bash
   git clone https://github.com/your-username/incident-management.git
   cd incident-management
   go mod tidy

## API Usage

### Create Incident (POST /api/v1/incidents)

**Request Body:**
```json
{
  "title": "Server Down",
  "description": "Production server is not responding to requests",
  "status": "open",
  "priority": "high"
}
```

**Response:**
```json
{
  "id": "uuid-here",
  "title": "Server Down",
  "description": "Production server is not responding to requests",
  "status": "open",
  "priority": "high",
  "ai_severity": "high",
  "ai_category": "hardware",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### Get All Incidents (GET /api/v1/incidents)

**Response:**
```json
[
  {
    "id": "uuid-here",
    "title": "Server Down",
    "description": "Production server is not responding to requests",
    "status": "open",
    "priority": "high",
    "ai_severity": "high",
    "ai_category": "hardware",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
]
```

### Health Check (GET /health)

**Response:**
```json
{
  "status": "ok",
  "message": "Incident Management API is running",
  "version": "1.0.0"
}
```

## Single Model Design with Validation

The `Incident` model serves both as input and output, with comprehensive tags for:
- **JSON serialization**: `json:"field_name"`
- **Database mapping**: `gorm:"constraints"`
- **Input validation**: `validate:"rules"`

```go
type Incident struct {
    ID          string    `json:"id" gorm:"primaryKey;type:varchar(36)" validate:"omitempty,uuid4"`
    Title       string    `json:"title" gorm:"not null" validate:"required,min=1,max=200"`
    Description string    `json:"description" gorm:"type:text" validate:"required,min=1,max=1000"`
    Status      string    `json:"status" gorm:"default:'open'" validate:"omitempty,oneof=open in_progress resolved closed"`
    Priority    string    `json:"priority" gorm:"default:'medium'" validate:"omitempty,oneof=low medium high critical"`
    AISeverity  string    `json:"ai_severity" gorm:"default:'medium'" validate:"omitempty,oneof=low medium high"`
    AICategory  string    `json:"ai_category" gorm:"default:'software'" validate:"omitempty,oneof=network software hardware security"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

## Validation Rules

The system validates all inputs with the following rules:

### Required Fields
- **Title**: Required, 1-200 characters
- **Description**: Required, 1-1000 characters

### Optional Fields with Constraints
- **Status**: Must be one of: `open`, `in_progress`, `resolved`, `closed`
- **Priority**: Must be one of: `low`, `medium`, `high`, `critical`
- **AISeverity**: Must be one of: `low`, `medium`, `high`
- **AICategory**: Must be one of: `network`, `software`, `hardware`, `security`

### Validation Error Response
```json
{
  "error": "Validation failed",
  "details": {
    "title": "title is required",
    "status": "status must be one of: open in_progress resolved closed"
  }
}
```

## Testing

The project includes comprehensive unit and integration tests.

### Running Tests

**Run all tests:**
```bash
go test ./...
```

**Run tests with verbose output:**
```bash
go test -v ./...
```

**Run tests with coverage:**
```bash
go test -cover ./...
```

**Run specific test packages:**
```bash
# Unit tests for services
go test ./services

# Unit tests for repository
go test ./repository

# Unit tests for handlers
go test ./handlers

# Unit tests for validation
go test ./utils

# Integration tests
go test -run TestIntegration
```

### Test Structure

- **Unit Tests:**
  - `services/ai_service_test.go` - Tests AI analysis functionality
  - `services/incident_services_test.go` - Tests incident service logic
  - `repository/incident_repository_test.go` - Tests database operations
  - `handlers/incident_handler_test.go` - Tests HTTP request handling
  - `utils/validator_test.go` - Tests validation functionality

- **Integration Tests:**
  - `integration_test.go` - Tests complete API flow from HTTP to database

- **Test Helpers:**
  - `test_helpers.go` - Common test utilities and validation functions

### Test Features

- **Database Testing:** Uses SQLite with automatic cleanup between tests
- **HTTP Testing:** Uses Gin's test utilities for API endpoint testing
- **AI Testing:** Tests both with and without OpenAI API key
- **Validation Testing:** Comprehensive validation rule testing
- **Coverage:** Tests cover all major functionality including edge cases

### Test Scenarios

1. **AI Service Tests:**
   - Service initialization
   - Incident analysis with and without API key
   - Text extraction and validation
   - Severity and category validation

2. **Incident Service Tests:**
   - Service initialization
   - Incident creation with defaults
   - Incident creation with provided values
   - Retrieving all incidents

3. **Repository Tests:**
   - Database connection
   - Creating and retrieving incidents
   - Multiple incident operations
   - Data persistence validation

4. **Handler Tests:**
   - HTTP request handling
   - JSON parsing and validation
   - Response status codes
   - Error handling
   - Validation error responses

5. **Validation Tests:**
   - Valid struct validation
   - Invalid struct validation
   - OneOf validation rules
   - Error message formatting
   - String sanitization

6. **Integration Tests:**
   - Complete API flow
   - Multiple incident creation and retrieval
   - AI analysis integration
   - End-to-end functionality

## Database

The application uses SQLite with automatic schema migration. The database file (`incidents.db`) will be created automatically when you first run the application.

## AI Classification

The AI service analyzes incident titles and descriptions to determine:

- **Severity Levels:** low, medium, high
- **Categories:** network, software, hardware, security

The AI uses GPT-3.5-turbo with a low temperature setting for consistent results.

## Architecture Benefits

The system combines the best of both worlds:

1. **Simplified Model**: Single struct for input/output reduces complexity
2. **Comprehensive Validation**: Ensures data integrity and provides clear error messages
3. **Type Safety**: Single source of truth for the data structure
4. **Clean Architecture**: Clear separation of concerns with validation layer
5. **Developer Experience**: Easy to understand and maintain

## Error Handling

The API provides detailed error responses:

- **400 Bad Request**: Invalid JSON format or validation failures
- **500 Internal Server Error**: Database or service errors

All error responses include:
- `error`: Human-readable error message
- `details`: Detailed error information (validation errors or technical details) 




LOGGING and PROMPT ENGINEERING

Asked AI to integrate with OPEN AI API and accepted the changes
Asked AI to setup the DB layer and it gave the suggestions for various calls, removed the extra code and used only POST AND GET
Asked AI to help with go-playground validations, used it as and asked AI to not create different structs for the same thing
Utilized AI to cover the test cases in both backend and frontend
Utilized AI for writing frontend code and for CSS