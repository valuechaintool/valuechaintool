package main

import "github.com/spf13/viper"

func setViperDefaults() {
	viper.SetDefault("title", "Value Chain Tracker")
	viper.SetDefault("minimumPasswordStrength", 3)
	viper.SetDefault("defaultRoles", []string{"28a98657-2c5b-435e-bfe5-18081066af8d"})
	viper.SetDefault("tiers", []map[string]interface{}{
		{
			"id":    0,
			"name":  "Undefined",
			"color": "FFFFFF",
		},
		{
			"id":    1,
			"name":  "None",
			"color": "000000",
		},
		{
			"id":    2,
			"name":  "Alliance",
			"color": "CD7F32",
		},
		{
			"id":    3,
			"name":  "Premium",
			"color": "B4B4B4",
		},
		{
			"id":    4,
			"name":  "Preferred",
			"color": "AF9500",
		},
	})
}
