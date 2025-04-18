package helpers

import (
	"fmt"
	"strconv"

	"github.com/mqdvi-dp/go-common/constants"
)

type QRIS struct {
	PayloadFormatIndicator string
	PointOfInitiation      string
	MerchantInfo           *MerchantInfo
	TransactionCurrency    string
	TransactionAmount      string
	MerchantName           string
	MerchantCity           string
	MerchantCategoryCode   string
	Checksum               string
	CountryCode            string
	AdditionalData         *AdditionalDataObject
	TipIndicator           string
	PostalCode             string
}

type MerchantInfo struct {
	ReverseDomain string
	AccountNumber string
	MerchantID    string
	MerchantType  string
}

type AdditionalDataObject struct {
	AdditionalMerchantAccount string
	MoreMerchantInfo          string
	MerchantType              string
}

func ParseQRIS(qrContent string) (*QRIS, error) {
	qris := new(QRIS)

	for len(qrContent) > 0 {
		if len(qrContent) < 4 {
			return nil, fmt.Errorf("QR string too short to parse")
		}

		tag := qrContent[:2]
		length, err := strconv.Atoi(qrContent[2:4])
		if err != nil {
			return nil, fmt.Errorf("invalid length format: %v", err)
		}

		if len(qrContent) < 4+length {
			return nil, fmt.Errorf("length mismatch for tag %s: expected %d, got %d", tag, length, len(qrContent[4:]))
		}

		value := qrContent[4 : 4+length]
		qrContent = qrContent[4+length:]

		switch tag {
		case constants.QRTagPayloadFormatIndicator:
			qris.PayloadFormatIndicator = value
		case constants.QRTagPointOfInitiation:
			qris.PointOfInitiation = value
		case constants.QRTagMerchantInfo:
			merchantInfo, err := parseMerchantInfo(value)
			if err != nil {
				return nil, err
			}
			qris.MerchantInfo = merchantInfo
		case constants.QRTagAdditionalData:
			additionalData, err := parseAdditionalDataObject(value)
			if err != nil {
				return nil, err
			}
			qris.AdditionalData = additionalData
		case constants.QRTagMerchantCategoryCode:
			qris.MerchantCategoryCode = value
		case constants.QRTagTransactionCurrency:
			qris.TransactionCurrency = value
		case constants.QRTagTransactionAmount:
			qris.TransactionAmount = value
		case constants.QRTagTipIndicator:
			qris.TipIndicator = value
		case constants.QRTagCountryCode:
			qris.CountryCode = value
		case constants.QRTagMerchantName:
			qris.MerchantName = value
		case constants.QRTagMerchantCity:
			qris.MerchantCity = value
		case constants.QRTagPostalCode:
			qris.PostalCode = value
		case constants.QRTagMoreMerchantInfo:
			qris.AdditionalData.MoreMerchantInfo = value
		case constants.QRTagChecksum:
			qris.Checksum = value
		default:
			fmt.Printf("Unhandled tag: %s with value %s\n", tag, value)
		}
	}

	return qris, nil
}

func parseMerchantInfo(merchantString string) (*MerchantInfo, error) {
	var merchantInfo MerchantInfo
	var err error

	for len(merchantString) > 0 {
		tag := merchantString[:2]

		var length int
		length, err = strconv.Atoi(merchantString[2:4])
		if err != nil {
			err = fmt.Errorf("Invalid length format in merchant info: %v", length)
			break
		}

		if len(merchantString) < 4+length {
			err = fmt.Errorf("Length mismatch for merchant subtag %s: expected %d, got %d\n", tag, length, len(merchantString[4:]))
			break
		}

		value := merchantString[4 : 4+length]
		merchantString = merchantString[4+length:]

		switch tag {
		case constants.MerchantTagReverseDomain:
			merchantInfo.ReverseDomain = value
		case constants.MerchantTagAccountNumber:
			merchantInfo.AccountNumber = value
		case constants.MerchantTagMerchantID:
			merchantInfo.MerchantID = value
		case constants.MerchantTagMerchantType:
			merchantInfo.MerchantType = value
		default:
			fmt.Printf("Unhandled merchant subtag: %s with value %s\n", tag, value)
		}
	}

	return &merchantInfo, err
}

func parseAdditionalDataObject(additionalString string) (*AdditionalDataObject, error) {
	var additionalData AdditionalDataObject
	var err error

	for len(additionalString) > 0 {
		tag := additionalString[:2]

		var length int
		length, err = strconv.Atoi(additionalString[2:4])
		if err != nil {
			err = fmt.Errorf("Invalid length format in additional data: %v", length)
			break
		}

		if len(additionalString) < 4+length {
			err = fmt.Errorf("Length mismatch for additional subtag %s: expected %d, got %d\n", tag, length, len(additionalString[4:]))
			break
		}

		value := additionalString[4 : 4+length]
		additionalString = additionalString[4+length:]

		switch tag {
		case constants.AdditionalTagMerchantAccount:
			additionalData.AdditionalMerchantAccount = value
		case constants.AdditionalTagMoreMerchantInfo:
			additionalData.MoreMerchantInfo = value
		case constants.AdditionalTagMerchantType:
			additionalData.MerchantType = value
		default:
			fmt.Printf("Unhandled additional subtag: %s with value %s\n", tag, value)
		}
	}

	return &additionalData, err
}
