package validation

import (
	"fmt"
	"net"
	"time"

	"github.com/gobcn/radius-director/internal/model"
)

func validateGlobalObjects(globalObjects model.GlobalObjects) []error {
	validationErrors := validateCredentialProfiles(globalObjects.CredentialProfiles)
	validationErrors = append(validationErrors, validateAuthenticationProfiles(globalObjects.AuthenticationProfiles)...)
	validationErrors = append(validationErrors, validateAccountingProfiles(globalObjects.AccountingProfiles)...)
	validationErrors = append(validationErrors, validateMonitoringProfiles(globalObjects.MonitoringProfiles)...)
	validationErrors = append(validationErrors, validateNASDevices(globalObjects.NASDevices)...)
	validationErrors = append(validationErrors, validateTrustedRADIUSClients(globalObjects.TrustedRADIUSClients)...)

	return validationErrors
}

func validateCredentialProfiles(profiles map[string]model.CredentialProfile) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(profiles) {
		validationErrors = append(validationErrors, validateCredentialProfile(identifier, profiles[identifier])...)
	}

	return validationErrors
}

func validateCredentialProfile(identifier string, profile model.CredentialProfile) []error {
	if profile.SharedSecret == "" {
		return []error{fmt.Errorf("credential profile %q: shared_secret must be specified", identifier)}
	}

	return nil
}

func validateAuthenticationProfiles(profiles map[string]model.AuthenticationProfile) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(profiles) {
		validationErrors = append(validationErrors, validateAuthenticationProfile(identifier, profiles[identifier])...)
	}

	return validationErrors
}

func validateAuthenticationProfile(identifier string, profile model.AuthenticationProfile) []error {
	if profile.SimultaneousUse != nil && *profile.SimultaneousUse < 1 {
		return []error{
			fmt.Errorf(
				"authentication profile %q: simultaneous_use must be greater than zero",
				identifier,
			),
		}
	}

	return nil
}

func validateAccountingProfiles(profiles map[string]model.AccountingProfile) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(profiles) {
		validationErrors = append(validationErrors, validateAccountingProfile(identifier, profiles[identifier])...)
	}

	return validationErrors
}

func validateAccountingProfile(identifier string, profile model.AccountingProfile) []error {
	if profile.StaleSessionTimeout == "" {
		return nil
	}

	duration, err := time.ParseDuration(profile.StaleSessionTimeout)
	if err != nil {
		return []error{fmt.Errorf("accounting profile %q: stale_session_timeout must be a valid duration", identifier)}
	}
	if duration <= 0 {
		return []error{fmt.Errorf("accounting profile %q: stale_session_timeout must be greater than zero", identifier)}
	}

	return nil
}

func validateMonitoringProfiles(profiles map[string]model.MonitoringProfile) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(profiles) {
		validationErrors = append(validationErrors, validateMonitoringProfile(identifier, profiles[identifier])...)
	}

	return validationErrors
}

func validateMonitoringProfile(identifier string, profile model.MonitoringProfile) []error {
	return nil
}

func validateNASDevices(devices map[string]model.NASDevice) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(devices) {
		validationErrors = append(validationErrors, validateNASDevice(identifier, devices[identifier])...)
	}

	return validationErrors
}

func validateNASDevice(identifier string, device model.NASDevice) []error {
	var validationErrors []error
	if net.ParseIP(device.IPAddress) == nil {
		validationErrors = append(validationErrors, fmt.Errorf("nas device %q: ip_address must be a valid IPv4 or IPv6 address", identifier))
	}
	if device.Vendor == "" {
		validationErrors = append(validationErrors, fmt.Errorf("nas device %q: vendor must be specified", identifier))
	}

	return validationErrors
}

func validateTrustedRADIUSClients(clients map[string]model.TrustedRADIUSClient) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(clients) {
		validationErrors = append(validationErrors, validateTrustedRADIUSClient(identifier, clients[identifier])...)
	}

	return validationErrors
}

func validateTrustedRADIUSClient(identifier string, client model.TrustedRADIUSClient) []error {
	if net.ParseIP(client.IPAddress) == nil {
		return []error{fmt.Errorf("trusted radius client %q: ip_address must be a valid IPv4 or IPv6 address", identifier)}
	}

	return nil
}
