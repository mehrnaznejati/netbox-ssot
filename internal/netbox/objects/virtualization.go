package objects

import "fmt"

type ClusterGroup struct {
	NetboxObject
	// Name is the name of the cluster group. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the cluster group name. This field is read-only.
	Slug string `json:"slug,omitempty"`
	// Description is a description of the cluster group.
}

type ClusterType struct {
	NetboxObject
	// Name is the name of the cluster type. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the cluster type name. This field is read-only.
	Slug string `json:"slug,omitempty"`
}

func (ct ClusterType) String() string {
	return fmt.Sprintf("ClusterType{Name: %s}", ct.Name)
}

type ClusterStatus struct {
	Choice
}

var (
	ClusterStatusActive  = ClusterStatus{Choice{Value: "active", Label: "Active"}}
	ClusterStatusOffline = ClusterStatus{Choice{Value: "offline", Label: "Offline"}}
)

type Cluster struct {
	NetboxObject
	// Name is the name of the cluster. This field is required.
	Name string `json:"name,omitempty"`
	// Type is the type of the cluster. This field is required.
	// e.g. oVirt,VMware...
	Type *ClusterType `json:"type,omitempty"`
	// ClusterGroup is the cluster group to which this cluster belongs.
	Group *ClusterGroup `json:"group,omitempty"`
	// Site is the site to which this cluster belongs.
	Site *Site `json:"site,omitempty"`
	// Status is the operational status of the cluster. This field is required.
	Status ClusterStatus `json:"status,omitempty"`
	// TenantGroup is the tenant group to which this cluster belongs.
	TenantGroup *TenantGroup `json:"tenant_group,omitempty"`
	// Tenant is the tenant to which this cluster belongs.
	Tenant *Tenant `json:"tenant,omitempty"`
}

func (c Cluster) String() string {
	return fmt.Sprintf("Cluster{ID: %d, Name: %s, Type: %s}", c.ID, c.Name, c.Type.Name)
}

type VMStatus struct {
	Choice
}

var (
	VMStatusActive  = VMStatus{Choice{Value: "active", Label: "Active"}}
	VMStatusOffline = VMStatus{Choice{Value: "offline", Label: "Offline"}}
)

// VM represents a netbox's virtual machine.
type VM struct {
	NetboxObject
	// Name is the name of the virtual machine. This field is required.
	Name string `json:"name,omitempty"`
	// Status is the status of the virtual machine. This field is required.
	Status *VMStatus `json:"status,omitempty"`
	// Site is the site to which this virtual machine belongs.
	Site *Site `json:"site,omitempty"`
	// Cluster is the cluster to which this virtual machine belongs.
	Cluster *Cluster `json:"cluster,omitempty"`
	// Host is a specific host that this virtual machine is hosted on.
	Host *Device `json:"device,omitempty"`

	// TenantGroup is the datacenter that this virtual machine belongs to.
	TenantGroup *TenantGroup `json:"tenant_group,omitempty"`
	// Tenant is the tenant to which this virtual machine belongs.
	Tenant *Tenant `json:"tenant,omitempty"`

	// Platform is the platform of the virtual machine.
	Platform *Platform `json:"platform,omitempty"`
	// PrimaryIPv4 is the primary IPv4 address assigned to the virtual machine.
	PrimaryIPv4 *IPAddress `json:"primary_ip4,omitempty"`
	// PrimaryIPv6 is the primary IPv6 address assigned to the virtual machine.
	PrimaryIPv6 *IPAddress `json:"primary_ip6,omitempty"`

	// VCPUs is the number of virtual CPUs allocated to the virtual machine.
	VCPUs float32 `json:"vcpus,omitempty"`
	// Memory is the amount of memory allocated to the virtual machine in MB.
	Memory int `json:"memory,omitempty"`
	// Disk is the amount of disk space allocated to the virtual machine in GB.
	Disk int `json:"disk,omitempty"`

	// Additional Comments
	Comments string `json:"comments,omitempty"`
}

func (vm VM) String() string {
	return fmt.Sprintf("VM{Id: %d, Name: %s}", vm.ID, vm.Name)
}

// 802.1Q VLAN Tagging Mode (Access, Tagged, Tagged All).
type VMInterfaceMode struct {
	Choice
}

var (
	VMInterfaceModeAccess    = VMInterfaceMode{Choice{Value: "access", Label: "Access"}}
	VMInterfaceModeTagged    = VMInterfaceMode{Choice{Value: "tagged", Label: "Tagged"}}
	VMInterfaceModeTaggedAll = VMInterfaceMode{Choice{Value: "tagged-all", Label: "Tagged All"}}
)

type VMInterface struct {
	NetboxObject
	// VM that this interface belongs to. This field is required.
	VM *VM `json:"virtual_machine,omitempty"`
	// Name is the name of the interface. This field is required.
	Name string `json:"name,omitempty"`
	// MAC address of the interface.
	MACAddress string `json:"mac_address,omitempty"`
	// MTU of the interface.
	MTU int `json:"mtu,omitempty"`
	// Enabled is true if interface is enabled, false otherwise.
	Enabled bool `json:"enabled,omitempty"`
	// Related parent interface of this interface.
	ParentInterface *VMInterface `json:"parent,omitempty"`
	// Related bridged interface
	BridgedInterface *VMInterface `json:"bridge,omitempty"`
	// 802.1Q VLAN Tagging Mode
	Mode *VMInterfaceMode `json:"mode,omitempty"`
	// When Mode=VMInterfaceModeTagged: TaggedVlans is a list of all the VLANs that are tagged on the interface.
	TaggedVlans []*Vlan `json:"tagged_vlans,omitempty"`
	// When mode=VMInterfaceModeAccess: UntaggedVlan is the VLAN that is untagged on the interface.
	UntaggedVlan *Vlan `json:"untagged_vlan,omitempty"`
}

func (vmi VMInterface) String() string {
	return fmt.Sprintf("VMInterface{Name: %s, VM: %s}", vmi.Name, vmi.VM.Name)
}
