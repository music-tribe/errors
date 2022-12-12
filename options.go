package errors

func SetCorrelationIDOption(correlationID string) CloudErrorOption {
	return func(se *CloudError) {
		se.CorrelationID = correlationID
	}
}
