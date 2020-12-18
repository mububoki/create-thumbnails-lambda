package lambda

type Informant struct{}

const serviceName = "lambda.amazonaws.com"

func (i *Informant) ServiceName() string {
	return serviceName
}
