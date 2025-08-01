package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/rest"

	"github.com/canonical/gomaasclient/client"
	"github.com/canonical/gomaasclient/entity"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/kogito-ops/cert-manager-webhook-maas/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

var GroupName = os.Getenv("GROUP_NAME")

func main() {
	if GroupName == "" {
		GroupName = "acme.maas.io"
		klog.Infof("GROUP_NAME environment variable not set, using default: %s", GroupName)
	}

	cmd.RunWebhookServer(GroupName,
		&maasDNSProviderSolver{},
	)
}

type maasDNSProviderSolver struct {
	client *kubernetes.Clientset
}

type maasDNSProviderConfig struct {
	SecretRef  string `json:"secretName"`
	ZoneName   string `json:"zoneName"`
	ApiUrl     string `json:"apiUrl"`
	ApiVersion string `json:"apiVersion,omitempty"`
}

func (c *maasDNSProviderSolver) Name() string {
	return "maas"
}

func (c *maasDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	klog.V(6).Infof("call function Present: namespace=%s, zone=%s, fqdn=%s",
		ch.ResourceNamespace, ch.ResolvedZone, ch.ResolvedFQDN)

	config, err := clientConfig(c, ch)
	if err != nil {
		return fmt.Errorf("unable to get secret `%s`; %v", ch.ResourceNamespace, err)
	}

	err = addTxtRecord(config, ch)
	if err != nil {
		return fmt.Errorf("unable to create TXT record: %v", err)
	}

	klog.Infof("Presented txt record %v", ch.ResolvedFQDN)
	return nil
}

func (c *maasDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	config, err := clientConfig(c, ch)
	if err != nil {
		return fmt.Errorf("unable to get secret `%s`; %v", ch.ResourceNamespace, err)
	}

	err = deleteTxtRecord(config, ch)
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("unable to delete TXT record: %v", err)
	}

	klog.Infof("Cleaned up TXT record for %v", ch.ResolvedFQDN)
	return nil
}

func (c *maasDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	k8sClient, err := kubernetes.NewForConfig(kubeClientConfig)
	klog.V(6).Infof("Input variable stopCh is %d length", len(stopCh))
	if err != nil {
		return err
	}

	c.client = k8sClient
	return nil
}

func loadConfig(cfgJSON *extapi.JSON) (maasDNSProviderConfig, error) {
	cfg := maasDNSProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	return cfg, nil
}

func stringFromSecretData(secretData map[string][]byte, key string) (string, error) {
	data, ok := secretData[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret data", key)
	}
	return string(data), nil
}

func addTxtRecord(config internal.Config, ch *v1alpha1.ChallengeRequest) error {
	maasClient, err := createMaasClient(config)
	if err != nil {
		return fmt.Errorf("failed to create MAAS client: %v", err)
	}

	// Create DNS resource record
	// For ACME challenges, we need to create TXT records
	// Extract the record name from the FQDN
	fqdn := strings.TrimSuffix(ch.ResolvedFQDN, ".")
	
	// Split the FQDN to get name and domain parts
	// The first part is the record name, the rest is the domain
	var name string
	var domain string
	
	parts := strings.SplitN(fqdn, ".", 2)
	if len(parts) == 2 {
		name = parts[0]
		domain = parts[1]
	} else {
		// Fallback to using FQDN if split fails
		return fmt.Errorf("unable to split FQDN %s into name and domain", fqdn)
	}
	
	klog.Infof("Creating TXT record - Name: %s, Domain: %s, Key: %s", name, domain, ch.Key)
	
	params := &entity.DNSResourceRecordParams{
		Name:   name,
		Domain: domain,
		RRData: ch.Key,
		RRType: "TXT",
		TTL:    120,
	}

	_, err = maasClient.DNSResourceRecords.Create(params)
	if err != nil {
		return fmt.Errorf("failed to create TXT record: %v", err)
	}

	klog.Infof("Added TXT record for %s", ch.ResolvedFQDN)
	return nil
}

func clientConfig(c *maasDNSProviderSolver, ch *v1alpha1.ChallengeRequest) (internal.Config, error) {
	var config internal.Config

	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return config, err
	}
	config.ZoneName = cfg.ZoneName
	config.ApiUrl = cfg.ApiUrl
	config.ApiVersion = cfg.ApiVersion
	if config.ApiVersion == "" {
		config.ApiVersion = "2.0" // Default API version
	}

	secretName := cfg.SecretRef
	sec, err := c.client.CoreV1().Secrets(ch.ResourceNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return config, fmt.Errorf("unable to get secret `%s/%s`; %v", secretName, ch.ResourceNamespace, err)
	}

	apiKey, err := stringFromSecretData(sec.Data, "api-key")
	if err != nil {
		return config, fmt.Errorf("unable to get api-key from secret `%s/%s`; %v", secretName, ch.ResourceNamespace, err)
	}
	config.ApiKey = apiKey

	// Get ZoneName by domain search if not provided by config
	if config.ZoneName == "" {
		foundZone, err := searchZoneName(config, ch.ResolvedZone)
		if err != nil {
			return config, err
		}
		config.ZoneName = foundZone
	}

	return config, nil
}

func createMaasClient(config internal.Config) (*client.Client, error) {
	return client.GetClient(config.ApiUrl, config.ApiKey, config.ApiVersion)
}

func deleteTxtRecord(config internal.Config, ch *v1alpha1.ChallengeRequest) error {
	maasClient, err := createMaasClient(config)
	if err != nil {
		return fmt.Errorf("failed to create MAAS client: %v", err)
	}

	// Remove trailing dot if present (MAAS doesn't use it)
	fqdn := strings.TrimSuffix(ch.ResolvedFQDN, ".")
	
	// Find the DNS resource record to delete
	records, err := maasClient.DNSResourceRecords.Get(&entity.DNSResourceRecordsParams{})
	if err != nil {
		return fmt.Errorf("failed to get DNS resource records: %v", err)
	}

	for _, record := range records {
		if record.FQDN == fqdn && record.RRType == "TXT" && record.RRData == ch.Key {
			err = maasClient.DNSResourceRecord.Delete(record.ID)
			if err != nil {
				return fmt.Errorf("failed to delete TXT record: %v", err)
			}
			klog.Infof("Deleted TXT record for %s", ch.ResolvedFQDN)
			return nil
		}
	}

	klog.Infof("TXT record not found for deletion: %s", ch.ResolvedFQDN)
	return nil
}

func searchZoneName(config internal.Config, searchZone string) (string, error) {
	maasClient, err := createMaasClient(config)
	if err != nil {
		return "", fmt.Errorf("failed to create MAAS client: %v", err)
	}

	// Get all domains and find matching one
	domains, err := maasClient.Domains.Get()
	if err != nil {
		return "", fmt.Errorf("failed to get domains: %v", err)
	}

	// Try to find exact match first
	for _, domain := range domains {
		if domain.Name == searchZone {
			return domain.Name, nil
		}
	}

	// Try to find parent domains
	parts := strings.Split(searchZone, ".")
	parts = parts[:len(parts)-1]
	for i := 0; i <= len(parts)-2; i++ {
		testZone := strings.Join(parts[i:], ".")
		for _, domain := range domains {
			if domain.Name == testZone {
				klog.Infof("Found domain: %s", domain.Name)
				return domain.Name, nil
			}
		}
	}

	return "", fmt.Errorf("unable to find MAAS domain for: %s", searchZone)
}
