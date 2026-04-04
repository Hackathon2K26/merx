package uniswap

import (
	"context"
	"fmt"
	"net/http"
)

// CheckApproval verifies whether token spending approval is needed.
// If Approval is nil in the response, the token is already approved.
func (c *Client) CheckApproval(ctx context.Context, req *ApprovalRequest) (*ApprovalResponse, error) {
	var resp ApprovalResponse
	if err := c.do(ctx, http.MethodPost, "/check_approval", req, &resp); err != nil {
		return nil, fmt.Errorf("check approval: %w", err)
	}
	return &resp, nil
}
