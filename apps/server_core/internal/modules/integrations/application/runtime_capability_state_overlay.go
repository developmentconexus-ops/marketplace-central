package application

import "marketplace-central/apps/server_core/internal/modules/integrations/domain"

func overlayRuntimeCapabilityStates(runtime []domain.RuntimeCapability, resolved []domain.CapabilityState) []domain.RuntimeCapability {
	if len(runtime) == 0 || len(resolved) == 0 {
		return runtime
	}

	byCode := make(map[string]domain.CapabilityState, len(resolved))
	for _, state := range resolved {
		if state.CapabilityCode == "" {
			continue
		}
		byCode[state.CapabilityCode] = state
	}

	items := make([]domain.RuntimeCapability, len(runtime))
	copy(items, runtime)

	for i := range items {
		state, ok := byCode[string(items[i].Code)]
		if !ok {
			items[i].LastValidatedAt = nil
			continue
		}
		applyCapabilityState(&items[i], state)
	}

	return items
}

func applyCapabilityState(runtime *domain.RuntimeCapability, state domain.CapabilityState) {
	if runtime == nil {
		return
	}

	runtime.LastValidatedAt = state.LastEvaluatedAt
	runtime.UnavailableReason = firstNonEmpty(state.ReasonCode, runtime.UnavailableReason)

	switch state.Status {
	case domain.CapabilityStatusEnabled:
		runtime.LiveValidated = state.LastEvaluatedAt != nil
	case domain.CapabilityStatusDegraded:
		runtime.State = domain.RuntimeCapabilityStateDegraded
		runtime.LiveValidated = false
	case domain.CapabilityStatusRequiresReauth:
		runtime.State = domain.RuntimeCapabilityStateNeedsAuth
		runtime.Executable = false
		runtime.LiveValidated = false
	case domain.CapabilityStatusUnsupported:
		runtime.State = domain.RuntimeCapabilityStateUnavailable
		runtime.Executable = false
		runtime.LiveValidated = false
	}
}
