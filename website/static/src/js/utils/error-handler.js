/**
 * Standardized Error Handler
 * Adapts backend error format to frontend UI
 */

export function handleApiError(errorResponse) {
  // Default fallback
  const defaultError = {
    message: 'common.op_fail',
    localizedMessage: 'Operation Failed',
  };

  if (!errorResponse) return defaultError;

  // Destructure with default values
  const { message, data } = errorResponse;

  // 1. Get translation key (message from backend)
  const key = message || defaultError.message;

  // 2. Localize and interpolate
  let localizedMessage = key;

  // Access global Alpine store safely
  // Note: In tests, window.Alpine might not be available unless mocked
  const i18n = window.Alpine?.store('i18n');

  if (i18n) {
    // i18n.t supports params replacement
    localizedMessage = i18n.t(key, data || {});
  } else {
    // Basic interpolation fallback if i18n is missing (e.g. early init error)
    if (data && typeof data === 'object') {
      Object.keys(data).forEach((param) => {
        localizedMessage = localizedMessage.replace(
          new RegExp(`{${param}}`, 'g'),
          data[param],
        );
      });
    }
  }

  return {
    ...errorResponse,
    localizedMessage,
  };
}
