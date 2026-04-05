package kafka

const (
	TopicStockUpdated    = "stock.updated"
	TopicStockRollback   = "stock.rollback"
	TopicDLQStockUpdated = "stock.updated.dlq"
	TopicDLQStockRollback = "stock.rollback.dlq"

	SchemaVersionStockEvent = 1
)
