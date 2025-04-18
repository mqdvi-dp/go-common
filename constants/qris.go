package constants

// Example Value for QR Content
// 00020101021226680016ID.CO.TELKOM.WWW011893600898025599662702150001952559966270303UMI51440014ID.CO.QRIS.WWW0215ID10200211817450303UMI520457325303360540825578.005502015802ID5916InterActive Corp6013KOTA SURABAYA61056013662130509413255111630439B7

// QRIS tags
const (
	// PayloadFormatIndicator (00): 01 - Indicates version 1 of the QRIS format
	QRTagPayloadFormatIndicator = "00"
	// PointOfInitiation (01): 12 - Indicates a dynamic QR code (likely generated per transaction).
	QRTagPointOfInitiation      = "01"
	// Merchant Information (26): Contains sub-fields:
	QRTagMerchantInfo           = "26"
	// Additional Data Object (51):
	QRTagAdditionalData         = "51"
	// Merchant Category Code (52): 5732 - This represents the type of merchant, which can be looked up in the MCC code list.
	QRTagMerchantCategoryCode   = "52"
	// Transaction Currency (53): 360 - Represents the currency, where 360 stands for Indonesian Rupiah (IDR).
	QRTagTransactionCurrency    = "53"
	// Transaction Amount (54): 25578.00 - The amount of the transaction.
	QRTagTransactionAmount      = "54"
	// Tip Indicator (55): 01 - Indicates whether a tip is included or not.
	QRTagTipIndicator           = "55"
	// Country Code (58): ID - Represents Indonesia.
	QRTagCountryCode            = "58"
	// Merchant Name (59): InterActive Corp - Name of the merchant.
	QRTagMerchantName           = "59"
	// Merchant City (60): KOTA SURABAYA - The city where the merchant is located.
	QRTagMerchantCity           = "60"
	// Postal Code (61): 60136 - Postal code of the merchant's location.
	QRTagPostalCode             = "61"
	// Additional Data (62):
	QRTagMoreMerchantInfo       = "62"
	// Checksum (63): 39B7 - The checksum to validate the QRIS code integrity
	QRTagChecksum               = "63"
)

// MerchantInfo tags
const (
	// 00: ID.CO.TELKOM.WWW - Reverse domain indicating the merchant.
	MerchantTagReverseDomain = "00"
	// 01: 936008980255996627 - Merchant's account number.
	MerchantTagAccountNumber = "01"
	// 02: 000195255996627 - Merchant ID as registered with the acquirer.
	MerchantTagMerchantID    = "02"
	// 03: UMI - Indicates the type of merchant, likely a small business.
	MerchantTagMerchantType  = "03"
)

// AdditionalDataObject tags
const (
	// 00: ID.CO.QRIS.WWW - Additional merchant account information.
	AdditionalTagMerchantAccount  = "00"
	// 02: ID10200211817450303UMI - Likely more detailed merchant or transaction information.
	AdditionalTagMoreMerchantInfo = "02"
	// 03: UMI - Indicates the type of merchant, likely a small business.
	AdditionalTagMerchantType     = "03"
)
