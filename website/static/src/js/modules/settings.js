export function Settings() {
  return {
    init() {
      // Any initialization if needed
    },
  };
}

export function HandleSettings(modifiedTime, cooldownDays, originalHandle) {
  return {
    modifiedTime: modifiedTime,
    cooldownDays: cooldownDays,
    originalHandle: originalHandle,
    currentHandle: originalHandle,
    get remainingDays() {
      const now = new Date().getTime();
      const cooldownMs = this.cooldownDays * 24 * 60 * 60 * 1000;
      const passedMs = now - this.modifiedTime;
      if (passedMs < cooldownMs) {
        return Math.ceil((cooldownMs - passedMs) / (24 * 60 * 60 * 1000));
      }
      return 0;
    },
    get isChanged() {
      return this.currentHandle !== this.originalHandle;
    },
    handleBeforeRequest(event) {
      if (this.remainingDays > 0 || !this.isChanged) {
        event.preventDefault();
        return;
      }
    },
    handleAfterRequest(event) {
      if (event.detail.successful) {
        this.modifiedTime = new Date().getTime();
        this.originalHandle = this.currentHandle;
      }
    },
  };
}
