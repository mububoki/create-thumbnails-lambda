package iam

import (
	"context"
	"fmt"

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

func (h *Handler) CreateRole(ctx context.Context, roleName string, serviceName string, actions []string) error {
	return h.createRole(ctx, roleName, assumeRolePolicyDocument(serviceName, actions))
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

const formatAssumeRolePolicyDocument = `{
    	"AssumeRolePolicyDocument": {
    	"Version" : "2012-10-17",
        	"Statement": [ {
            	"Effect": "Allow",
                "Principal": {
                	"Service": [ "%s" ]
				},
                "Action": %s
			} ]
		}
	}`

func assumeRolePolicyDocument(serviceName string, actions []string) string {
	a := "["

	for _, action := range actions {
		a += `"` + action + `",`
	}
	a = a[:len(a)-1] + "]"

	return fmt.Sprintf(formatAssumeRolePolicyDocument, serviceName, a)
}
