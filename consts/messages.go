package consts

const (
	ERR_UPDATE_FAIL              = "Fail to update"
	ERR_LIST_PATHFINDER          = "Error listing pathfinders"
	ERR_GET_PATHFINDER_REGION    = "Error getting pathfinder region"
	ERR_REGION_UNSPECIFIED       = "Region Unspecified"
	ERR_SERVICE_NAME_UNSPECIFIED = "Service Name Unspecified"
	ERR_WEBHOOK_INIT_FAIL        = "Unable to initialize webhook"

	INFO_UPDATINGPATHFINDER = "Updating PathFinder"
	INFO_START_CLEANUP      = "Starting cleanup"

	WARN_REGION_UNSPECIFIED   = "Region unspecified. Using default"
	WARN_NO_SERVICE_IN_REGION = "No service found in region"
	WARN_REGION_NOT_FOUND     = "Region not found"
	WARN_REGION_INCONSISTENT  = "In consistent region"
)

const (
	F_ERR_REGION_NOT_FOUND  = "Region not found %s %s"
	F_ERR_DUPLICATED_REGION = "Duplicated region found %s %s"
)

type ErrCode int

const (
	CODE_DUP_PF               ErrCode = 10000
	CODE_REGION_NOT_FOUND     ErrCode = 10001
	CODE_SVC_NAME_UNSPECIFIED ErrCode = 10002
)
