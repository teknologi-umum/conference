package features

type FeatureFlag struct {
	EnableRegistration              bool `yaml:"enable_registration" envconfig:"FEATURE_ENABLE_REGISTRATION" default:"false"`
	EnablePaymentProofUpload        bool `yaml:"enable_payment_proof_upload" envconfig:"FEATURE_ENABLE_PAYMENT_PROOF_UPLOAD" default:"false"`
	EnableCallForProposalSubmission bool `yaml:"enable_call_for_proposal_submission" envconfig:"FEATURE_ENABLE_CALL_FOR_PROPOSAL_SUBMISSION" default:"false"`
	EnableAdministratorMode         bool `yaml:"enable_administrator_mode" envconfig:"FEATURE_ENABLE_ADMINISTRATOR_MODE" default:"false"`
}
