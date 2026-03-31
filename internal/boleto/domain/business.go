package domain

type Business struct {
	Id                string `dynamodbav:"id"`
	Company           string `dynamodbav:"company"`
	Name              string `dynamodbav:"name"`
	CNPJ              string `dynamodbav:"cnpj"`
	StateRegistration string `dynamodbav:"stateRegistration"`
}
