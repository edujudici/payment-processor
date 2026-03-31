package domain

type Bank struct {
	Id                      string `dynamodbav:"id"`
	Company                 string `dynamodbav:"company"`
	BankAccountAbbreviation string `dynamodbav:"bankAccountAbbreviation"`
	BankCode                int    `dynamodbav:"bankCode"`
	BankAgreement           string `dynamodbav:"bankAgreement"`
	PixKey                  string `dynamodbav:"pixKey"`
	ModelOfGood             string `dynamodbav:"modelOfGood"`
	OriginCode              string `dynamodbav:"originCode"`
	NameBank                string `dynamodbav:"nameBank"`
	Agency                  int    `dynamodbav:"agency"`
	Account                 int    `dynamodbav:"account"`
	DigitAccount            int    `dynamodbav:"digitAccount"`
	PaymentMethod           string `dynamodbav:"paymentMethod"`
}
