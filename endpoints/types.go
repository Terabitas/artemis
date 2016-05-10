package endpoints

type (
	// NetworkInterface type
	NetworkInterface struct {
		ID string
		IP string
	}

	// Provider type
	Provider struct {
		ID     string
		Region string
		Size   string
		APIKey string
		Image  string
		SSHKey string
	}

	// Node type
	Node struct {
		ID           string
		Provider     Provider
		PublicIFace  NetworkInterface
		PrivateIFace NetworkInterface
	}

	// HealthPolicy type
	HealthPolicy struct {
		ID                string
		Min               int
		Max               int
		Desired           int
		HealthyThreshold  float64
		CheckInterval     int
		Provider          Provider
		ConsecutiveChecks int
	}
)
