package inventory

import (
	"slices"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// AddTag adds the newTag from source sourceName to the local inventory.
func (nbi *NetboxInventory) AddTag(newTag *objects.Tag) (*objects.Tag, error) {
	existingTagIndex := slices.IndexFunc(nbi.Tags, func(t *objects.Tag) bool {
		return t.Name == newTag.Name
	})
	if existingTagIndex == -1 {
		nbi.Logger.Debug("Tag ", newTag.Name, " does not exist in Netbox. Creating it...")
		createdTag, err := service.Create[objects.Tag](nbi.NetboxAPI, newTag)
		if err != nil {
			return nil, err
		}
		nbi.Tags = append(nbi.Tags, createdTag)
		return createdTag, nil
	}
	nbi.Logger.Debug("Tag ", newTag.Name, " already exists in Netbox...")
	oldTag := nbi.Tags[existingTagIndex]
	diffMap, err := utils.JSONDiffMapExceptID(newTag, oldTag, false, nbi.SourcePriority)
	if err != nil {
		return nil, err
	}
	if len(diffMap) > 0 {
		patchedTag, err := service.Patch[objects.Tag](nbi.NetboxAPI, oldTag.ID, diffMap)
		if err != nil {
			return nil, err
		}
		nbi.Tags[existingTagIndex] = patchedTag
		return patchedTag, nil
	}
	return oldTag, nil
}

// AddContact adds a contact to the local netbox inventory.
func (nbi *NetboxInventory) AddSite(newSite *objects.Site) (*objects.Site, error) {
	newSite.Tags = append(newSite.Tags, nbi.SsotTag)
	if _, ok := nbi.SitesIndexByName[newSite.Name]; ok {
		oldSite := nbi.SitesIndexByName[newSite.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newSite, oldSite, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Site ", newSite.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedSite, err := service.Patch[objects.Site](nbi.NetboxAPI, oldSite.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.SitesIndexByName[newSite.Name] = patchedSite
		} else {
			nbi.Logger.Debug("Site ", newSite.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Site ", newSite.Name, " does not exist in Netbox. Creating it...")
		createdContact, err := service.Create[objects.Site](nbi.NetboxAPI, newSite)
		if err != nil {
			return nil, err
		}
		nbi.SitesIndexByName[newSite.Name] = createdContact
	}
	return nbi.SitesIndexByName[newSite.Name], nil
}

// AddContactRole adds the newContactRole to the local netbox inventory.
func (nbi *NetboxInventory) AddContactRole(newContactRole *objects.ContactRole) (*objects.ContactRole, error) {
	newContactRole.NetboxObject.Tags = []*objects.Tag{nbi.SsotTag}
	if _, ok := nbi.ContactRolesIndexByName[newContactRole.Name]; ok {
		oldContactRole := nbi.ContactRolesIndexByName[newContactRole.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newContactRole, oldContactRole, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Contact role ", newContactRole.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedContactRole, err := service.Patch[objects.ContactRole](nbi.NetboxAPI, oldContactRole.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ContactRolesIndexByName[newContactRole.Name] = patchedContactRole
		} else {
			nbi.Logger.Debug("Contact role ", newContactRole.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Contact role ", newContactRole.Name, " does not exist in Netbox. Creating it...")
		newContactRole, err := service.Create[objects.ContactRole](nbi.NetboxAPI, newContactRole)
		if err != nil {
			return nil, err
		}
		nbi.ContactRolesIndexByName[newContactRole.Name] = newContactRole
	}
	return nbi.ContactRolesIndexByName[newContactRole.Name], nil
}

// AddContactGroup adds contact group to the local netbox inventory.
func (nbi *NetboxInventory) AddContactGroup(newContactGroup *objects.ContactGroup) (*objects.ContactGroup, error) {
	if _, ok := nbi.ContactGroupsIndexByName[newContactGroup.Name]; ok {
		oldContactGroup := nbi.ContactGroupsIndexByName[newContactGroup.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newContactGroup, oldContactGroup, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Contact group ", newContactGroup.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedContactGroup, err := service.Patch[objects.ContactGroup](nbi.NetboxAPI, oldContactGroup.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ContactGroupsIndexByName[newContactGroup.Name] = patchedContactGroup
		} else {
			nbi.Logger.Debug("Contact group ", newContactGroup.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Contact group ", newContactGroup.Name, " does not exist in Netbox. Creating it...")
		newContactGroup, err := service.Create[objects.ContactGroup](nbi.NetboxAPI, newContactGroup)
		if err != nil {
			return nil, err
		}
		nbi.ContactGroupsIndexByName[newContactGroup.Name] = newContactGroup
	}
	return nbi.ContactGroupsIndexByName[newContactGroup.Name], nil
}

// AddContact adds a contact to the local netbox inventory.
func (nbi *NetboxInventory) AddContact(newContact *objects.Contact) (*objects.Contact, error) {
	newContact.Tags = append(newContact.Tags, nbi.SsotTag)
	if _, ok := nbi.ContactsIndexByName[newContact.Name]; ok {
		oldContact := nbi.ContactsIndexByName[newContact.Name]
		delete(nbi.OrphanManager[service.ContactsAPIPath], oldContact.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newContact, oldContact, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Contact ", newContact.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedContact, err := service.Patch[objects.Contact](nbi.NetboxAPI, oldContact.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ContactsIndexByName[newContact.Name] = patchedContact
		} else {
			nbi.Logger.Debug("Contact ", newContact.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Contact ", newContact.Name, " does not exist in Netbox. Creating it...")
		createdContact, err := service.Create[objects.Contact](nbi.NetboxAPI, newContact)
		if err != nil {
			return nil, err
		}
		nbi.ContactsIndexByName[newContact.Name] = createdContact
	}
	return nbi.ContactsIndexByName[newContact.Name], nil
}

// AddContact assignment adds a contact assignment to the local netbox inventory.
// TODO: Make index check less code and more universal, checking each level is ugly.
func (nbi *NetboxInventory) AddContactAssignment(newCA *objects.ContactAssignment) (*objects.ContactAssignment, error) {
	if nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType] == nil {
		nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType] = make(map[int]map[int]map[int]*objects.ContactAssignment)
	}
	if nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType][newCA.ObjectID] == nil {
		nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType][newCA.ObjectID] = make(map[int]map[int]*objects.ContactAssignment)
	}
	if nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType][newCA.ObjectID][newCA.Contact.ID] == nil {
		nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType][newCA.ObjectID][newCA.Contact.ID] = make(map[int]*objects.ContactAssignment)
	}
	newCA.Tags = append(newCA.Tags, nbi.SsotTag)
	if _, ok := nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID]; ok {
		oldCA := nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID]
		delete(nbi.OrphanManager[service.ContactAssignmentsAPIPath], oldCA.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newCA, oldCA, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("ContactAssignment ", newCA.ID, " already exists in Netbox but is out of date. Patching it... ")
			patchedCA, err := service.Patch[objects.ContactAssignment](nbi.NetboxAPI, oldCA.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID] = patchedCA
		} else {
			nbi.Logger.Debug("ContactAssignment ", newCA.ID, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debugf("ContactAssignment %s does not exist in Netbox. Creating it...", newCA)
		newCA, err := service.Create[objects.ContactAssignment](nbi.NetboxAPI, newCA)
		if err != nil {
			return nil, err
		}
		nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID] = newCA
	}
	return nbi.ContactAssignmentsIndexByContentTypeAndObjectIDAndContactIDAndRoleID[newCA.ContentType][newCA.ObjectID][newCA.Contact.ID][newCA.Role.ID], nil
}

func (nbi *NetboxInventory) AddCustomField(newCf *objects.CustomField) error {
	if _, ok := nbi.CustomFieldsIndexByName[newCf.Name]; ok {
		oldCustomField := nbi.CustomFieldsIndexByName[newCf.Name]
		diffMap, err := utils.JSONDiffMapExceptID(newCf, oldCustomField, false, nbi.SourcePriority)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Custom field ", newCf.Name, " already exists in Netbox but is out of date. Patching it... ")
			patchedCf, err := service.Patch[objects.CustomField](nbi.NetboxAPI, oldCustomField.ID, diffMap)
			if err != nil {
				return err
			}
			nbi.CustomFieldsIndexByName[newCf.Name] = patchedCf
		} else {
			nbi.Logger.Debug("Custom field ", newCf.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Custom field ", newCf.Name, " does not exist in Netbox. Creating it...")
		newCf, err := service.Create[objects.CustomField](nbi.NetboxAPI, newCf)
		if err != nil {
			return err
		}
		nbi.CustomFieldsIndexByName[newCf.Name] = newCf
	}
	return nil
}

func (nbi *NetboxInventory) AddClusterGroup(newCg *objects.ClusterGroup) (*objects.ClusterGroup, error) {
	newCg.Tags = append(newCg.Tags, nbi.SsotTag)
	if _, ok := nbi.ClusterGroupsIndexByName[newCg.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldCg := nbi.ClusterGroupsIndexByName[newCg.Name]
		delete(nbi.OrphanManager[service.ClusterGroupsAPIPath], oldCg.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newCg, oldCg, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Cluster group ", newCg.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedCg, err := service.Patch[objects.ClusterGroup](nbi.NetboxAPI, oldCg.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ClusterGroupsIndexByName[newCg.Name] = patchedCg
		} else {
			nbi.Logger.Debug("Cluster group ", newCg.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Cluster group ", newCg.Name, " does not exist in Netbox. Creating it...")
		newCg, err := service.Create[objects.ClusterGroup](nbi.NetboxAPI, newCg)
		if err != nil {
			return nil, err
		}
		nbi.ClusterGroupsIndexByName[newCg.Name] = newCg
	}
	// Delete id from orphan manager
	return nbi.ClusterGroupsIndexByName[newCg.Name], nil
}

func (nbi *NetboxInventory) AddClusterType(newClusterType *objects.ClusterType) (*objects.ClusterType, error) {
	newClusterType.Tags = append(newClusterType.Tags, nbi.SsotTag)
	if _, ok := nbi.ClusterTypesIndexByName[newClusterType.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldClusterType := nbi.ClusterTypesIndexByName[newClusterType.Name]
		delete(nbi.OrphanManager[service.ClusterTypesAPIPath], oldClusterType.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newClusterType, oldClusterType, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedClusterType, err := service.Patch[objects.ClusterType](nbi.NetboxAPI, oldClusterType.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ClusterTypesIndexByName[newClusterType.Name] = patchedClusterType
			return patchedClusterType, nil
		}
		nbi.Logger.Debug("Cluster type ", newClusterType.Name, " already exists in Netbox and is up to date...")
		existingClusterType := nbi.ClusterTypesIndexByName[newClusterType.Name]
		return existingClusterType, nil
	}
	nbi.Logger.Debug("Cluster type ", newClusterType.Name, " does not exist in Netbox. Creating it...")
	newClusterType, err := service.Create[objects.ClusterType](nbi.NetboxAPI, newClusterType)
	if err != nil {
		return nil, err
	}
	nbi.ClusterTypesIndexByName[newClusterType.Name] = newClusterType
	return newClusterType, nil
}

func (nbi *NetboxInventory) AddCluster(newCluster *objects.Cluster) error {
	newCluster.Tags = append(newCluster.Tags, nbi.SsotTag)
	if _, ok := nbi.ClustersIndexByName[newCluster.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldCluster := nbi.ClustersIndexByName[newCluster.Name]
		delete(nbi.OrphanManager[service.ClustersAPIPath], oldCluster.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newCluster, oldCluster, false, nbi.SourcePriority)
		if err != nil {
			return err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Cluster ", newCluster.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedCluster, err := service.Patch[objects.Cluster](nbi.NetboxAPI, oldCluster.ID, diffMap)
			if err != nil {
				return err
			}
			nbi.ClustersIndexByName[newCluster.Name] = patchedCluster
		} else {
			nbi.Logger.Debug("Cluster ", newCluster.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Cluster ", newCluster.Name, " does not exist in Netbox. Creating it...")
		newCluster, err := service.Create[objects.Cluster](nbi.NetboxAPI, newCluster)
		if err != nil {
			return err
		}
		nbi.ClustersIndexByName[newCluster.Name] = newCluster
	}
	return nil
}

func (nbi *NetboxInventory) AddDeviceRole(newDeviceRole *objects.DeviceRole) (*objects.DeviceRole, error) {
	newDeviceRole.Tags = append(newDeviceRole.Tags, nbi.SsotTag)
	if _, ok := nbi.DeviceRolesIndexByName[newDeviceRole.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldDeviceRole := nbi.DeviceRolesIndexByName[newDeviceRole.Name]
		delete(nbi.OrphanManager[service.DeviceRolesAPIPath], nbi.DeviceRolesIndexByName[newDeviceRole.Name].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newDeviceRole, oldDeviceRole, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedDeviceRole, err := service.Patch[objects.DeviceRole](nbi.NetboxAPI, oldDeviceRole.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.DeviceRolesIndexByName[newDeviceRole.Name] = patchedDeviceRole
		} else {
			nbi.Logger.Debug("Device role ", newDeviceRole.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Device role ", newDeviceRole.Name, " does not exist in Netbox. Creating it...")
		newDeviceRole, err := service.Create[objects.DeviceRole](nbi.NetboxAPI, newDeviceRole)
		if err != nil {
			return nil, err
		}
		nbi.DeviceRolesIndexByName[newDeviceRole.Name] = newDeviceRole
	}
	return nbi.DeviceRolesIndexByName[newDeviceRole.Name], nil
}

func (nbi *NetboxInventory) AddManufacturer(newManufacturer *objects.Manufacturer) (*objects.Manufacturer, error) {
	newManufacturer.Tags = append(newManufacturer.Tags, nbi.SsotTag)
	if _, ok := nbi.ManufacturersIndexByName[newManufacturer.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldManufacturer := nbi.ManufacturersIndexByName[newManufacturer.Name]
		delete(nbi.OrphanManager[service.ManufacturersAPIPath], oldManufacturer.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newManufacturer, oldManufacturer, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Manufacturer ", newManufacturer.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedManufacturer, err := service.Patch[objects.Manufacturer](nbi.NetboxAPI, oldManufacturer.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.ManufacturersIndexByName[newManufacturer.Name] = patchedManufacturer
		} else {
			nbi.Logger.Debug("Manufacturer ", newManufacturer.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Manufacturer ", newManufacturer.Name, " does not exist in Netbox. Creating it...")
		newManufacturer, err := service.Create[objects.Manufacturer](nbi.NetboxAPI, newManufacturer)
		if err != nil {
			return nil, err
		}
		nbi.ManufacturersIndexByName[newManufacturer.Name] = newManufacturer
	}
	return nbi.ManufacturersIndexByName[newManufacturer.Name], nil
}

func (nbi *NetboxInventory) AddDeviceType(newDeviceType *objects.DeviceType) (*objects.DeviceType, error) {
	newDeviceType.Tags = append(newDeviceType.Tags, nbi.SsotTag)
	if _, ok := nbi.DeviceTypesIndexByModel[newDeviceType.Model]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldDeviceType := nbi.DeviceTypesIndexByModel[newDeviceType.Model]
		delete(nbi.OrphanManager[service.DeviceTypesAPIPath], oldDeviceType.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newDeviceType, oldDeviceType, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Device type ", newDeviceType.Model, " already exists in Netbox but is out of date. Patching it...")
			patchedDeviceType, err := service.Patch[objects.DeviceType](nbi.NetboxAPI, oldDeviceType.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.DeviceTypesIndexByModel[newDeviceType.Model] = patchedDeviceType
		} else {
			nbi.Logger.Debug("Device type ", newDeviceType.Model, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Device type ", newDeviceType.Model, " does not exist in Netbox. Creating it...")
		newDeviceType, err := service.Create[objects.DeviceType](nbi.NetboxAPI, newDeviceType)
		if err != nil {
			return nil, err
		}
		nbi.DeviceTypesIndexByModel[newDeviceType.Model] = newDeviceType
	}
	return nbi.DeviceTypesIndexByModel[newDeviceType.Model], nil
}

func (nbi *NetboxInventory) AddPlatform(newPlatform *objects.Platform) (*objects.Platform, error) {
	newPlatform.Tags = append(newPlatform.Tags, nbi.SsotTag)
	if _, ok := nbi.PlatformsIndexByName[newPlatform.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldPlatform := nbi.PlatformsIndexByName[newPlatform.Name]
		delete(nbi.OrphanManager[service.PlatformsAPIPath], oldPlatform.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newPlatform, oldPlatform, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Platform ", newPlatform.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedPlatform, err := service.Patch[objects.Platform](nbi.NetboxAPI, oldPlatform.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.PlatformsIndexByName[newPlatform.Name] = patchedPlatform
		} else {
			nbi.Logger.Debug("Platform ", newPlatform.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Platform ", newPlatform.Name, " does not exist in Netbox. Creating it...")
		newPlatform, err := service.Create[objects.Platform](nbi.NetboxAPI, newPlatform)
		if err != nil {
			return nil, err
		}
		nbi.PlatformsIndexByName[newPlatform.Name] = newPlatform
	}
	return nbi.PlatformsIndexByName[newPlatform.Name], nil
}

func (nbi *NetboxInventory) AddDevice(newDevice *objects.Device) (*objects.Device, error) {
	newDevice.Tags = append(newDevice.Tags, nbi.SsotTag)
	if _, ok := nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID]; ok {
		oldDevice := nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID]
		delete(nbi.OrphanManager[service.DevicesAPIPath], oldDevice.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newDevice, oldDevice, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Device ", newDevice.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedDevice, err := service.Patch[objects.Device](nbi.NetboxAPI, oldDevice.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID] = patchedDevice
		} else {
			nbi.Logger.Debug("Device ", newDevice.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Device ", newDevice.Name, " does not exist in Netbox. Creating it...")
		newDevice, err := service.Create[objects.Device](nbi.NetboxAPI, newDevice)
		if err != nil {
			return nil, err
		}
		if nbi.DevicesIndexByNameAndSiteID[newDevice.Name] == nil {
			nbi.DevicesIndexByNameAndSiteID[newDevice.Name] = make(map[int]*objects.Device)
		}
		nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID] = newDevice
	}
	return nbi.DevicesIndexByNameAndSiteID[newDevice.Name][newDevice.Site.ID], nil
}

func (nbi *NetboxInventory) AddVlanGroup(newVlanGroup *objects.VlanGroup) (*objects.VlanGroup, error) {
	newVlanGroup.Tags = append(newVlanGroup.Tags, nbi.SsotTag)
	if _, ok := nbi.VlanGroupsIndexByName[newVlanGroup.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldVlanGroup := nbi.VlanGroupsIndexByName[newVlanGroup.Name]
		delete(nbi.OrphanManager[service.VlanGroupsAPIPath], oldVlanGroup.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newVlanGroup, oldVlanGroup, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("VlanGroup ", newVlanGroup.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVlanGroup, err := service.Patch[objects.VlanGroup](nbi.NetboxAPI, oldVlanGroup.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VlanGroupsIndexByName[newVlanGroup.Name] = patchedVlanGroup
		} else {
			nbi.Logger.Debug("Vlan ", newVlanGroup.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Vlan ", newVlanGroup.Name, " does not exist in Netbox. Creating it...")
		newVlan, err := service.Create[objects.VlanGroup](nbi.NetboxAPI, newVlanGroup)
		if err != nil {
			return nil, err
		}
		nbi.VlanGroupsIndexByName[newVlan.Name] = newVlan
	}
	return nbi.VlanGroupsIndexByName[newVlanGroup.Name], nil
}

func (nbi *NetboxInventory) AddVlan(newVlan *objects.Vlan) (*objects.Vlan, error) {
	newVlan.Tags = append(newVlan.Tags, nbi.SsotTag)
	if _, ok := nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		oldVlan := nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid]
		delete(nbi.OrphanManager[service.VlansAPIPath], oldVlan.ID)
		diffMap, err := utils.JSONDiffMapExceptID(newVlan, oldVlan, false, nbi.SourcePriority)
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Vlan ", newVlan.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVlan, err := service.Patch[objects.Vlan](nbi.NetboxAPI, oldVlan.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid] = patchedVlan
		} else {
			nbi.Logger.Debug("Vlan ", newVlan.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Vlan ", newVlan.Name, " does not exist in Netbox. Creating it...")
		newVlan, err := service.Create[objects.Vlan](nbi.NetboxAPI, newVlan)
		if err != nil {
			return nil, err
		}
		if nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID] == nil {
			nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID] = make(map[int]*objects.Vlan)
		}
		nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid] = newVlan
	}
	return nbi.VlansIndexByVlanGroupIDAndVID[newVlan.Group.ID][newVlan.Vid], nil
}

func (nbi *NetboxInventory) AddInterface(newInterface *objects.Interface) (*objects.Interface, error) {
	newInterface.Tags = append(newInterface.Tags, nbi.SsotTag)
	if _, ok := nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[service.InterfacesAPIPath], nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newInterface, nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name], false, nbi.SourcePriority)
		oldIntf := nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Interface ", newInterface.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedInterface, err := service.Patch[objects.Interface](nbi.NetboxAPI, oldIntf.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name] = patchedInterface
		} else {
			nbi.Logger.Debug("Interface ", newInterface.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("Interface ", newInterface.Name, " does not exist in Netbox. Creating it...")
		newInterface, err := service.Create[objects.Interface](nbi.NetboxAPI, newInterface)
		if err != nil {
			return nil, err
		}
		if nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID] == nil {
			nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID] = make(map[string]*objects.Interface)
		}
		nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name] = newInterface
	}
	return nbi.InterfacesIndexByDeviceIDAndName[newInterface.Device.ID][newInterface.Name], nil
}

func (nbi *NetboxInventory) AddVM(newVM *objects.VM) (*objects.VM, error) {
	newVM.Tags = append(newVM.Tags, nbi.SsotTag)
	if _, ok := nbi.VMsIndexByName[newVM.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[service.VirtualMachinesAPIPath], nbi.VMsIndexByName[newVM.Name].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newVM, nbi.VMsIndexByName[newVM.Name], false, nbi.SourcePriority)
		oldVM := nbi.VMsIndexByName[newVM.Name]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("VM ", newVM.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVM, err := service.Patch[objects.VM](nbi.NetboxAPI, oldVM.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VMsIndexByName[newVM.Name] = patchedVM
		} else {
			nbi.Logger.Debug("VM ", newVM.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("VM ", newVM.Name, " does not exist in Netbox. Creating it...")
		newVM, err := service.Create[objects.VM](nbi.NetboxAPI, newVM)
		if err != nil {
			return nil, err
		}
		nbi.VMsIndexByName[newVM.Name] = newVM
		return newVM, nil
	}
	return nbi.VMsIndexByName[newVM.Name], nil
}

func (nbi *NetboxInventory) AddVMInterface(newVMInterface *objects.VMInterface) (*objects.VMInterface, error) {
	newVMInterface.Tags = append(newVMInterface.Tags, nbi.SsotTag)
	if _, ok := nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name]; ok {
		// Remove id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[service.VMInterfacesAPIPath], nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newVMInterface, nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name], false, nbi.SourcePriority)
		oldVMIface := nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("VM interface ", newVMInterface.Name, " already exists in Netbox but is out of date. Patching it...")
			patchedVMInterface, err := service.Patch[objects.VMInterface](nbi.NetboxAPI, oldVMIface.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name] = patchedVMInterface
		} else {
			nbi.Logger.Debug("VM interface ", newVMInterface.Name, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("VM interface ", newVMInterface.Name, " does not exist in Netbox. Creating it...")
		newVMInterface, err := service.Create[objects.VMInterface](nbi.NetboxAPI, newVMInterface)
		if err != nil {
			return nil, err
		}
		if nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID] == nil {
			nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID] = make(map[string]*objects.VMInterface)
		}
		nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name] = newVMInterface
	}
	return nbi.VMInterfacesIndexByVMIdAndName[newVMInterface.VM.ID][newVMInterface.Name], nil
}

func (nbi *NetboxInventory) AddIPAddress(newIPAddress *objects.IPAddress) (*objects.IPAddress, error) {
	newIPAddress.Tags = append(newIPAddress.Tags, nbi.SsotTag)
	if _, ok := nbi.IPAdressesIndexByAddress[newIPAddress.Address]; ok {
		// Delete id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[service.IPAddressesAPIPath], nbi.IPAdressesIndexByAddress[newIPAddress.Address].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newIPAddress, nbi.IPAdressesIndexByAddress[newIPAddress.Address], false, nbi.SourcePriority)
		oldIPAddress := nbi.IPAdressesIndexByAddress[newIPAddress.Address]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("IP address ", newIPAddress.Address, " already exists in Netbox but is out of date. Patching it...")
			patchedIPAddress, err := service.Patch[objects.IPAddress](nbi.NetboxAPI, oldIPAddress.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.IPAdressesIndexByAddress[newIPAddress.Address] = patchedIPAddress
		} else {
			nbi.Logger.Debug("IP address ", newIPAddress.Address, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("IP address ", newIPAddress.Address, " does not exist in Netbox. Creating it...")
		newIPAddress, err := service.Create[objects.IPAddress](nbi.NetboxAPI, newIPAddress)
		if err != nil {
			return nil, err
		}
		nbi.IPAdressesIndexByAddress[newIPAddress.Address] = newIPAddress
		return newIPAddress, nil
	}
	return nbi.IPAdressesIndexByAddress[newIPAddress.Address], nil
}

func (nbi *NetboxInventory) AddPrefix(newPrefix *objects.Prefix) (*objects.Prefix, error) {
	newPrefix.Tags = append(newPrefix.Tags, nbi.SsotTag)
	if _, ok := nbi.PrefixesIndexByPrefix[newPrefix.Prefix]; ok {
		// Delete id from orphan manager, because it still exists in the sources
		delete(nbi.OrphanManager[service.PrefixesAPIPath], nbi.PrefixesIndexByPrefix[newPrefix.Prefix].ID)
		diffMap, err := utils.JSONDiffMapExceptID(newPrefix, nbi.PrefixesIndexByPrefix[newPrefix.Prefix], false, nbi.SourcePriority)
		oldPrefix := nbi.PrefixesIndexByPrefix[newPrefix.Prefix]
		if err != nil {
			return nil, err
		}
		if len(diffMap) > 0 {
			nbi.Logger.Debug("Prefix ", newPrefix.Prefix, " already exists in Netbox but is out of date. Patching it...")
			patchedPrefix, err := service.Patch[objects.Prefix](nbi.NetboxAPI, oldPrefix.ID, diffMap)
			if err != nil {
				return nil, err
			}
			nbi.PrefixesIndexByPrefix[newPrefix.Prefix] = patchedPrefix
		} else {
			nbi.Logger.Debug("IP address ", newPrefix.Prefix, " already exists in Netbox and is up to date...")
		}
	} else {
		nbi.Logger.Debug("IP address ", newPrefix.Prefix, " does not exist in Netbox. Creating it...")
		newPrefix, err := service.Create[objects.Prefix](nbi.NetboxAPI, newPrefix)
		if err != nil {
			return nil, err
		}
		nbi.PrefixesIndexByPrefix[newPrefix.Prefix] = newPrefix
		return newPrefix, nil
	}
	return nbi.PrefixesIndexByPrefix[newPrefix.Prefix], nil
}
