package plugins

// OpenShiftAzureConfig is a new, Red Hat owned type defined in an external project.  It
// represents a superset of the configuration required for all steps in the plugin process *after*
// the manifest is validated.
//
// OpenShiftAzureConfig is a versioned type.  It is expected that the configuration will change
// between releases and may require migration steps in the future.  The storage of this type
// is opaque to the plugins, they only require that it be passed in.  However, it must be retained
// to allow it to be lifecycled from version to version.
type OpenShiftAzureConfig struct{}

// ManagedOpenShiftCluster is the external schema that is passed to the RP.  It is defined in
// acs-engine.  ManagedOpenShiftCluster will reuse many of the lower level objects shared by
// ManagedCluster but is separated in order to keep the relationship with AKS clear.
type ManagedOpenShiftCluster struct{}

// Validator interface is the first plugin in the chain.  It is responsible for validation of
// both ManagedOpenShiftCluster.  ManagedOpenShiftCluster enters this process in a versioned state.
// It is validated both according to the version and against the previous state of the world.
// Versioned validation should occur first.  The manifest, if acceptable, may then be converted
// into the internal representation and validated against the previous state of the world.
//
// Since this plugin validates with both a versioned manifest and an unversioned manifest it will
// be necessary to cycle the plugin along side the api versions.
//
// TODO: is the lifecycle here acceptable?  It is how api.validate.go works.  Does this make it harder
// on the RP calling the method?  Might have to provide an api version aware factory to return the correct
// implementation.
type Validator interface {
	Validate(new *ManagedOpenShiftCluster, old *ManagedOpenShiftCluster) (ManagedOpenShiftCluster, []error)
}

// ConfigManager is responsible for the generate and/or upgrade of the OpenShiftAzureConfig.  The config
// generated here is a superset of all configurations used by other plugins. If an existing config
// exists for the ManagedOpenShiftCluster it is passed in so it may be used to transfer values to the new config.
// A new configuration is always returned.
type ConfigManager interface {
	Generate(apimodel *ManagedOpenShiftCluster, existingConfig *OpenShiftAzureConfig) (OpenShiftAzureConfig, error)
}

// HCPManager generates the values required by helm charts that are created for the customer control plane.
type HCPManager interface {
	Generate(apimodel *ManagedOpenShiftCluster, existingConfig *OpenShiftAzureConfig) ([]byte, error)
}

// NodeManager
// TODO: is anything other than the call to acs-engine needed here? This one is still the most unclear!
type NodeManager interface {
	Generate(apimodel *ManagedOpenShiftCluster, existingConfig *OpenShiftAzureConfig) ([]byte, error)
}

// AddOnManager generates any configuration necessary for add ons that will be created by helm.  It is
// separate from HCPManager for clarity.
type AddOnManager interface {
	Generate(apimodel *ManagedOpenShiftCluster, existingConfig *OpenShiftAzureConfig) ([]byte, error)
}

// HealthChecker ensures that a cluster may be handed over to the consumer.  This should check that
// all created components are in a running state and may wait for some time to confirm stability.
type HealthChecker interface {
	Check(apimodel *ManagedOpenShiftCluster, config *OpenShiftAzureConfig) []error
}
