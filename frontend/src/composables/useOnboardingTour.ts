/**
 * Onboarding tour stub — example Guide steps were removed in Plan A scaffold.
 * Keeps AppHeader / store call sites compiling without driver.js tour content.
 */

export interface OnboardingOptions {
  storageKey?: string
  autoStart?: boolean
}

export function useOnboardingTour(_options: OnboardingOptions) {
  const replayTour = () => {
    // no-op: example tour steps deleted with Guide/
  }

  return {
    replayTour,
    startTour: async () => {},
    stopTour: () => {}
  }
}
