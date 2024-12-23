package linodego

import (
	"context"
)

// IPAddressUpdateOptions fields are those accepted by UpdateToken
type IPAddressUpdateOptions struct {
	// The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to a default value provided by Linode if set to nil.
	Reserved *bool   `json:"reserved,omitempty"`
	RDNS     *string `json:"rdns,omitempty"`
}

// LinodeIPAssignment stores an assignment between an IP address and a Linode instance.
type LinodeIPAssignment struct {
	Address  string `json:"address"`
	LinodeID int    `json:"linode_id"`
}

type AllocateReserveIPOptions struct {
	Type     string `json:"type"`
	Public   bool   `json:"public"`
	Reserved bool   `json:"reserved,omitempty"`
	Region   string `json:"region,omitempty"`
	LinodeID int    `json:"linode_id,omitempty"`
}

// LinodesAssignIPsOptions fields are those accepted by InstancesAssignIPs.
type LinodesAssignIPsOptions struct {
	Region string `json:"region"`

	Assignments []LinodeIPAssignment `json:"assignments"`
}

// IPAddressesShareOptions fields are those accepted by ShareIPAddresses.
type IPAddressesShareOptions struct {
	IPs      []string `json:"ips"`
	LinodeID int      `json:"linode_id"`
}

// ListIPAddressesQuery fields are those accepted as query params for the
// ListIPAddresses function.
type ListIPAddressesQuery struct {
	SkipIPv6RDNS bool `query:"skip_ipv6_rdns"`
}

// GetUpdateOptions converts a IPAddress to IPAddressUpdateOptions for use in UpdateIPAddress
func (i InstanceIP) GetUpdateOptions() (o IPAddressUpdateOptions) {
	o.RDNS = copyString(&i.RDNS)
	return
}

// ListIPAddresses lists IPAddresses
func (c *Client) ListIPAddresses(ctx context.Context, opts *ListOptions) ([]InstanceIP, error) {
	response, err := getPaginatedResults[InstanceIP](ctx, c, "networking/ips", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetIPAddress gets the IPAddress with the provided IP
func (c *Client) GetIPAddress(ctx context.Context, id string) (*InstanceIP, error) {
	e := formatAPIPath("networking/ips/%s", id)
	response, err := doGETRequest[InstanceIP](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateIPAddress updates the IPAddress with the specified id
func (c *Client) UpdateIPAddress(ctx context.Context, id string, opts IPAddressUpdateOptions) (*InstanceIP, error) {
	e := formatAPIPath("networking/ips/%s", id)
	response, err := doPUTRequest[InstanceIP](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// InstancesAssignIPs assigns multiple IPv4 addresses and/or IPv6 ranges to multiple Linodes in one Region.
// This allows swapping, shuffling, or otherwise reorganizing IPs to your Linodes.
func (c *Client) InstancesAssignIPs(ctx context.Context, opts LinodesAssignIPsOptions) error {
	e := "networking/ips/assign"
	_, err := doPOSTRequest[InstanceIP](ctx, c, e, opts)
	return err
}

// ShareIPAddresses allows IP address reassignment (also referred to as IP failover)
// from one Linode to another if the primary Linode becomes unresponsive.
func (c *Client) ShareIPAddresses(ctx context.Context, opts IPAddressesShareOptions) error {
	e := "networking/ips/share"
	_, err := doPOSTRequest[InstanceIP](ctx, c, e, opts)
	return err
}

// AllocateReserveIP allocates a new IPv4 address to the Account, with the option to reserve it
// and optionally assign it to a Linode.
func (c *Client) AllocateReserveIP(ctx context.Context, opts AllocateReserveIPOptions) (*InstanceIP, error) {
	e := "networking/ips"
	result, err := doPOSTRequest[InstanceIP](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}
