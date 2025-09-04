package constants

// Error messages
const (
	ErrTaskNotFound          = "task not found"
	ErrUserNotFound          = "user not found"
	ErrInvalidTaskTitle      = "task title cannot be empty"
	ErrInvalidTaskStatus     = "invalid task status given in req"
	ErrInvalidTaskPriority   = "invalid task priority given in req"
	ErrInternalServer        = "internal server error"
	ErrInvalidRequestBody    = "Invalid request body"
	ErrFailedToCreateTask    = "Failed to create task"
	ErrFailedToGetTask       = "Failed to get task"
	ErrFailedToUpdateTask    = "Failed to update task"
	ErrFailedToDeleteTask    = "Failed to delete task"
	ErrFailedToListTasks     = "Failed to list tasks"
	ErrFailedToConnectDB     = "Failed to connect to database"
	ErrFailedToGetSqlDB      = "Failed to get sql.DB"
	ErrFailedToMigrateDB     = "Failed to migrate database"
	ErrorStartingApplication = "Error starting application"
	ErrorClosingDb           = "Error closing postgres db"
)

// Default values
const (
	DefaultPage     = 1
	DefaultPageSize = 10
)

// Query parameter names
const (
	QueryParamStatus   = "status"
	QueryParamUserID   = "user_id"
	QueryParamPriority = "priority"
	QueryParamPage     = "page"
	QueryParamPageSize = "pageSize"
)

// URL parameter names
const (
	URLParamUUID = "uuid"
)

// Default string values
const (
	DefaultPageStr     = "1"
	DefaultPageSizeStr = "10"
)

// Response messages
const (
	HealthCheckOKMessage = "ok"
)

const (
	KafkaHosts     = "KAFKA_HOSTS"
	KafkaUsername  = "KAFKA_JAAS_CONFIG_USERNAME"
	KafkaPassword  = "KAFKA_JAAS_CONFIG_PASSWORD"
	KafkaGroupId   = "KAFKA_GROUP_ID"
	KafkaAuthAlgo  = "KAFKA_AUTH_ALGO"
	Err            = "err"
	AppName        = "APP_NAME"
	AppVersion     = "APP_VERSION"
	UserAppBaseUri = "TESSERACT_BASE_URI"

	KafkaRetryTopic  = "KAFKA_RETRY_TOPIC"
	PostgresAddress  = "POSTGRES_ADDRESS"
	PostgresUsername = "POSTGRES_USERNAME"
	PostgresPassword = "POSTGRES_PASSWORD"
	PostgresDbName   = "POSTGRES_DB_NAME"
	PostgresHost     = "POSTGRES_HOST"
	PostgresPort     = "POSTGRES_PORT"
	MaxDbConnections = "MAX_DB_CONNECTIONS"
	HOST             = "HOST"
	PORT             = "PORT"
	TCP              = "tcp"
	DriverName       = "postgres"
)
