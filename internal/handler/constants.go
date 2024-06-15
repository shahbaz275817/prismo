package handler

const (
	API                         = "api"
	BadRequest                  = "BAD_REQUEST"
	Unauthorized                = "UNAUTHORIZED"
	InvalidShipmentData         = "INVALID_SHIPMENT_DATA"
	DuplicateBatchFile          = "DUPLICATE_BATCH_FILE"
	MissingDriverNameField      = "MISSING_DRIVER_NAME_FIELD"
	MissingDriverPhoneField     = "MISSING_DRIVER_PHONE_FIELD"
	MissingDriverFleetTypeField = "MISSING_DRIVER_FLEET_TYPE_FIELD"
	MissingDriverVendorField    = "MISSING_DRIVER_VENDOR_FIELD"
	MissingDriverPeopleIDField  = "MISSING_DRIVER_PEOPLE_ID_FIELD"
	MissingDriverEmailField     = "MISSING_DRIVER_EMAIL_FIELD"
	InvalidDriverPhoneFormat    = "INVALID_DRIVER_PHONE_FORMAT"
	InvalidDriverNameFormat     = "INVALID_DRIVER_NAME_FORMAT"
	InvalidDriverFleetType      = "INVALID_DRIVER_FLEET_TYPE"
	InvalidDriverVendor         = "INVALID_DRIVER_VENDOR"
	InvalidDriverEmailFormat    = "INVALID_DRIVER_EMAIL"
	DriverNameLengthExceeded    = "DRIVER_NAME_LENGTH_EXCEEDED"
	DuplicateDriverPhone        = "DUPLICATE_DRIVER_PHONE"
	InternalServerError         = "INTERNAL_SERVER_ERROR"
	MissingAWBNumberField       = "MISSING_AWB_NUMBER_FIELD"
	NotFoundShipment            = "SHIPMENT_NOT_FOUND"
	StatusTransitionError       = "STATUS_TRANSITION_ERROR"
	InvalidCreateBatchRequest   = "INVALID_CREATE_BATCH_REQUEST"
	HubNotFound                 = "HUB_NOT_FOUND"
	InvalidSellerGroupData      = "INVALID_SELLER_GROUP_DATA"
	InvalidShipmentCount        = "INVALID_SHIPMENT_COUNT"
)

// params

const (
	batchCodeParam = "batch_code"
)
