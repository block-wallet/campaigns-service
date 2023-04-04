package signatures

import "fmt"

const VERIFICATION_MESSAGE_PREFIX = "Sign this message to enroll in"

func GenerateDefaultCampaignVerificationMessage(campaignName string) string {
	return fmt.Sprintf("%v %v", VERIFICATION_MESSAGE_PREFIX, campaignName)
}
