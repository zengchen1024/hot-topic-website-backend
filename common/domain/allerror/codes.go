package allerror

const (
	errorCodeNotFound     = "not_found"
	errorCodeOverLimited  = "over_limited"
	errorCodeNoPermission = "no_permission"

	ErrorCodeMissingDS = "missing_ds"
	ErrorCodeMissingHT = "missing_ht"

	ErrorCodeReviewDuplicateDS      = "duplicate_ds"
	ErrorCodeReviewDuplicateTopic   = "duplicate_topic"
	ErrorCodeReviewNotConstantOrder = "not_constant_order"

	ErrorCodeInvokeTimeRestricted = "invoke_time_restricted"

	ErrorCodeNoMatchedTopicsToReview = "no_matched_topics_to_review"
)
