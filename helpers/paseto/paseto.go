package paseto

import (
	"context"

	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/o1egl/paseto"
)

func GeneratePaseto(ctx context.Context, privateKeyValue string, data interface{}, footers ...interface{}) (string, error) {
	// initiate paseto
	pg := paseto.NewV2()

	// get privateKey from env
	privateKey, err := getPrivateKey(ctx, privateKeyValue)
	if err != nil {
		tracer.SetError(ctx, err)

		return "", err
	}

	var footer interface{}
	if len(footers) > 0 {
		footer = footers[0]
	}

	payload, err := convert.InterfaceToBytes(data)
	if err != nil {
		tracer.SetError(ctx, err)

		return "", err
	}

	sign, err := pg.Sign(privateKey, payload, footer)
	if err != nil {
		tracer.SetError(ctx, err)

		return "", err
	}

	return sign, nil
}

func VerifyPaseto(ctx context.Context, publicKeyValue, data string) (interface{}, interface{}, error) {
	// initiate paseto
	pg := paseto.NewV2()

	// get publicKey from env
	publicKey, err := getPublicKey(ctx, publicKeyValue)
	if err != nil {
		tracer.SetError(ctx, err)

		return nil, nil, err
	}

	var footer interface{}

	var payload interface{}
	err = pg.Verify(data, publicKey, &payload, &footer)
	if err != nil {
		tracer.SetError(ctx, err)

		return nil, nil, err
	}

	return payload, footer, nil
}
