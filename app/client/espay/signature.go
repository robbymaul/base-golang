package espay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/enums"
	"strings"
	"time"
)

type ISignature interface {
	HashBashSignatureGenerate(payload PaymentRequest, action enums.EnumSignatureEspay) string
	AsymmetricSignatureGenerate(body any, key string) (string, error)
	StringToSign(httpMethod string, path string, signature string, timestamp time.Time) string
	GenerateXSignature(sign string, key string) (string, error)
}

type Signature struct {
	//Req          PaymentRequest
	signatureKey string
	CommCode     string
	//privateKey   string
}

func (s *Signature) StringToSign(httpMethod string, path string, signature string, timestamp time.Time) string {
	stringToSign := httpMethod + ":" + path + ":" + signature + ":" + timestamp.Format(time.RFC3339)
	//
	//block, _ := pem.Decode([]byte(""))
	//if block == nil {
	//	return "", errors.New("failed to decode private key")
	//}
	//
	//var privKey any
	//var err error
	//if privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
	//	return "", fmt.Errorf("failed to parse private key: %w", err)
	//}
	//
	//rsaKey, ok := privKey.(*rsa.PrivateKey)
	//if !ok {
	//	return "", errors.New("failed to get private key")
	//}
	//
	//hash := sha256.Sum256([]byte(stringToSign))
	//
	//sign, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hash[:])
	//if err != nil {
	//	return "", fmt.Errorf("failed to sign: %w", err)
	//}

	//return base64.StdEncoding.EncodeToString(sign), nil
	return stringToSign
}

func NewSignature(signatureKey string, commCode string) ISignature {
	return &Signature{
		//Req:          req,
		signatureKey: signatureKey,
		CommCode:     commCode,
		//privateKey:   privateKey,
	}
}

/*
ACTION = SENDINVOICE | INQUIRY | | CHECKSTATUS | EXPIRETRANSACTION
*/
func (s *Signature) HashBashSignatureGenerate(payload PaymentRequest, action enums.EnumSignatureEspay) string {
	log.Debug().Interface("action", action).Interface("context", "signature espay").Msg("generate signature espay")

	//var signature string
	//stringAction := fmt.Sprint(action)

	//data := s.signatureKey + s.Req.RQUUID + s.Req.OrderID + strings.ToUpper(stringAction)

	//sum := md5.Sum([]byte(data))
	//signature = hex.EncodeToString(sum[:])

	rawString := fmt.Sprintf(
		"##%v##%v##%v##%v##%v##%v##%v##%v##",
		s.signatureKey,
		payload.RQUUID,
		payload.RQDateTime.Format(time.DateTime),
		payload.OrderID,
		payload.Amount,
		payload.CCY,
		s.CommCode,
		action,
	)
	upperString := strings.ToUpper(rawString)
	log.Debug().Str("rawString", rawString).Str("upper string", upperString).Msg("raw string to hash sha 256 generate signature espay")

	hash := sha256.Sum256([]byte(upperString))
	log.Debug().Interface("hash", hash).Msg("hash to sha256 generate signature espay")

	return hex.EncodeToString(hash[:])
}

func (s *Signature) AsymmetricSignatureGenerate(body any, key string) (string, error) {
	marshal, err := json.Marshal(body)
	if err != nil {
		log.Error().Err(err).Interface("body", body).Msg("failed to marshal body")
		return "", fmt.Errorf("invalid body minify signature")
	}
	log.Debug().Interface("byte marshal body", marshal).Str("string marshal body", string(marshal)).Msg("minify signature byte body ")

	sum256 := sha256.Sum256(marshal)
	log.Debug().Interface("sum256", sum256).Msg("sum256 signature")

	encodeString := hex.EncodeToString(sum256[:])
	log.Debug().Str("encode string", encodeString).Msg("encode string signature")

	lowerString := strings.ToLower(encodeString)
	log.Debug().Str("lower string", lowerString).Msg("lower string signature")

	sign := s.StringToSign(http.MethodPost, PathPaymentHostToHost, lowerString, time.Now())
	log.Debug().Str("sign", sign).Msg("sign signature")

	generateXSignature, err := s.GenerateXSignature(sign, key)
	if err != nil {
		return "", err
	}
	log.Debug().Str("generateXSignature", generateXSignature).Msg("generate x signature")

	return generateXSignature, nil
}

func (s *Signature) GenerateXSignature(sign string, key string) (string, error) {

	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return "", errors.New("failed to decode private key")
	}

	var privKey any
	var err error
	if privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("failed to get private key")
	}

	hash := sha256.Sum256([]byte(sign))

	signPkc, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signPkc), nil
}
