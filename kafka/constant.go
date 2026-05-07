package kafka

const (
	TopicOrderCreated     = "order.created"
	TopicStockReserved    = "stock.reserved"
	TopicStockRejected    = "stock.rejected"
	TopicStockUpdated     = "stock.updated"
	TopicStockRollback    = "stock.rollback"
	TopicDLQOrderCreated  = "order.created.dlq"
	TopicDLQStockUpdated  = "stock.updated.dlq"
	TopicDLQStockRollback = "stock.rollback.dlq"

	SchemaVersionStockEvent = 1
)
