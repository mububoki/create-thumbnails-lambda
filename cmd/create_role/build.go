package main

import (
	"golang.org/x/xerrors"

	"github.com/mububoki/create-thumbnails-lambda/internal/app/infrastructure/env"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/infrastructure/iam"
	"github.com/mububoki/create-thumbnails-lambda/internal/app/infrastructure/lambda"
	ifauthnz "github.com/mububoki/create-thumbnails-lambda/internal/app/interface/authnz"
	ucauthnz "github.com/mububoki/create-thumbnails-lambda/internal/app/usecase/authnz"
)

func build() (*ucauthnz.Interactor, error) {
	iamHandler, err := iam.NewHandler()
	if err != nil {
		return nil, xerrors.Errorf("failed to NewHandler: %w", err)
	}
	lambdaInformant := new(lambda.Informant)

	repo := ifauthnz.NewRepository(iamHandler, lambdaInformant)

	return ucauthnz.NewInteractor(repo, env.AuthNZ.RoleName, env.AuthNZ.Actions), nil
}
