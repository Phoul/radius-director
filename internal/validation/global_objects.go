package validation

import (
	"fmt"
	"net"

	"github.com/gobcn/radius-director/internal/model"
)

func validateGlobalObjects(globalObjects model.GlobalObjects) []error {
	validationErrors := validateCredentialProfiles(globalObjects.CredentialProfiles)
	validationErrors = append(validationErrors, validateAuthenticationProfiles(globalObjects.AuthenticationProfiles)...)
	validationErrors = append(validationErrors, validateAccountingProfiles(globalObjects.AccountingProfiles)...)
	validationErrors = append(validationErrors, validateMonitoringProfiles(globalObjects.MonitoringProfiles)...)
	validationErrors = append(validationErrors, validateNASDevices(globalObjects.NASDevices)...)

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
