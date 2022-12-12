package errors

func SetCorrelationIDOption(correlationID string) StorageErrorOption {
	return func(se *StorageError) {
		se.CorrelationID = correlationID
	}
}
