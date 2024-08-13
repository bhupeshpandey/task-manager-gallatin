package models

import (
	"time"
)

type Event struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Task struct {
	Id          string    `json:"id"`
	ParentID    string    `json:"parent_id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	ParentId    string `json:"parent_id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskResponse struct {
	Id string `json:"id"`
}

type UpdateTaskRequest struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskResponse struct {
	Success bool `json:"success"`
}

type ListTasksResponse struct {
	Tasks []*Task `json:"tasks"`
}

type ListTasksRequest struct {
	StartTime string
	EndTime   string
	Limit     int
	Offset    int
}

type DeleteTaskRequest struct {
	Id string `json:"id"`
}

type DeleteTaskResponse struct {
	Success bool `json:"success"`
}

// Config represents the application configuration
type Config struct {
	Server       ServerConfig       `yaml:"server"`
	Database     DatabaseConfig     `yaml:"database"`
	Logging      LoggingConfig      `yaml:"logging"`
	MessageQueue MessageQueueConfig `yaml:"message_queue"`
	Cache        CacheConfig        `yaml:"cache"`
}

type CacheConfig struct {
	Type  string      `yaml:"type"`
	Redis RedisConfig `yaml:"redis"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// ServerConfig represents the server-related configuration
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Type     string         `yaml:"type"`
	Postgres PostgresConfig `yaml:"postgres"`
	SQLite   SQLiteConfig   `yaml:"sqlite"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type SQLiteConfig struct {
	Filepath string `yaml:"filepath"`
}

// LoggingConfig represents the logging-related configuration
type LoggingConfig struct {
	Type        string `yaml:"type"`
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"logLevel"`
}

// MessageQueueConfig represents the message queue-related configuration
type MessageQueueConfig struct {
	Type     string          `yaml:"type"`
	RabbitMQ *RabbitMQConfig `yaml:"rabbitmq,omitempty"`
	Kafka    *KafkaConfig    `yaml:"kafka,omitempty"`
}

// RabbitMQConfig represents the RabbitMQ-specific configuration
type RabbitMQConfig struct {
	URL        string `yaml:"url"`
	Exchange   string `yaml:"exchange"`
	Queue      string `yaml:"queue"`
	RoutingKey string `yaml:"routing_key"`
}

// KafkaConfig represents the Kafka-specific configuration
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"group_id"`
}
