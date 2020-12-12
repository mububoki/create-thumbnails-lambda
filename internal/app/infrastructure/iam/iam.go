package iam

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"golang.org/x/xerrors"
)

type Handler struct {
	iam iamiface.IAMAPI
}

func NewHandler() (*Handler, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, xerrors.Errorf("failed to NewSession: %w", err)
	}

	return &Handler{
		iam: iam.New(sess),
	}, nil
}

func (h *Handler) createRole(ctx context.Context, roleName string, assumeRolePolicyDocument string) error {
	if _, err := h.iam.CreateRoleWithContext(ctx, &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assumeRolePolicyDocument),
		RoleName:                 aws.String(roleName),
	}); err != nil {
		return xerrors.Errorf("failed to CreateRole: %w", err)
	}

	return nil
}
