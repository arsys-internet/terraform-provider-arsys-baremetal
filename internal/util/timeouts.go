package util

import (
	"time"
)

type TimeoutConfig struct {
	Default       time.Duration
	RetryInterval time.Duration
	MinTimeout    time.Duration
}

func (tc TimeoutConfig) Get() (time.Duration, time.Duration, time.Duration) {
	return tc.Default, tc.RetryInterval, tc.MinTimeout
}

func GetResourceTimeouts(resourcePrefix string) TimeoutConfig {
	defaults := getDefaultsForResource(resourcePrefix)

	timeout, err := GetEnvTimeValues(resourcePrefix+"_DEFAULT_TIMEOUT", time.Minute)
	if err != nil {
		timeout = defaults.timeout
	}

	retryInterval, err := GetEnvTimeValues(resourcePrefix+"_DEFAULT_RETRY_INTERVAL", time.Second)
	if err != nil {
		retryInterval = defaults.retryInterval
	}

	minTimeout, err := GetEnvTimeValues(resourcePrefix+"_DEFAULT_MIN_TIMEOUT", time.Second)
	if err != nil {
		minTimeout = defaults.minTimeout
	}

	return TimeoutConfig{
		Default:       timeout,
		RetryInterval: retryInterval,
		MinTimeout:    minTimeout,
	}
}

type resourceDefaults struct {
	timeout       time.Duration
	retryInterval time.Duration
	minTimeout    time.Duration
}

func getDefaultsForResource(resourcePrefix string) resourceDefaults {
	switch resourcePrefix {
	case "SERVER":
		return resourceDefaults{40 * time.Minute, 30 * time.Second, 20 * time.Second}
	case "PUBLIC_NETWORK":
		return resourceDefaults{8 * time.Minute, 15 * time.Second, 15 * time.Second}
	case "FIREWALL_POLICY":
		return resourceDefaults{5 * time.Minute, 10 * time.Second, 5 * time.Second}
	case "FIREWALL_POLICY_SERVER_IPS":
		return resourceDefaults{20 * time.Minute, 20 * time.Second, 5 * time.Second}
	case "PUBLIC_IP":
		return resourceDefaults{5 * time.Minute, 10 * time.Second, 5 * time.Second}
	default:
		// Generic default fallback
		return resourceDefaults{20 * time.Minute, 20 * time.Second, 5 * time.Second}
	}
}
