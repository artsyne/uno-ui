// Shim to access global Uno UI objects in ES modules
// Since Uno is loaded via script tag (UMD), it exposes 'Uno' to window

const getAlpineStore = (name) => {
  // Check if Alpine is available globally or on window
  const alpine = window.Alpine || (window.top && window.top.Alpine);
  if (alpine && alpine.store) {
    return alpine.store(name);
  }
  return null;
};

export const Box = {
  confirm: (options) => {
    const store = getAlpineStore('box');
    // 'confirm' is the type, options is the second arg
    return store
      ? store.open('confirm', options)
      : console.log('Uno.box.open (Store not found)', 'confirm', options);
  },
  dialog: (options) => {
    const store = getAlpineStore('box');
    // 'dialog' is the default type
    return store
      ? store.open('dialog', options)
      : console.log('Uno.box.open (Store not found)', 'dialog', options);
  },
  close: (id) => {
    const store = getAlpineStore('box');
    return store
      ? store.close(id)
      : console.log('Uno.box.close (Store not found)', id);
  },
};

export const Notify = {
  success: (msg) => {
    const store = getAlpineStore('notify');
    return store
      ? store.success(msg, { position: 'bottom-right' })
      : console.log('Success:', msg);
  },
  error: (msg) => {
    const store = getAlpineStore('notify');
    return store
      ? store.error(msg, { position: 'bottom-right' })
      : console.error('Error:', msg);
  },
};

export const Valid = {
  validateForm: async (form) => {
    const store = getAlpineStore('valid');
    if (store && typeof store.validateForm === 'function') {
      return await store.validateForm(form);
    }
    // If store missing or method missing, log warning and allow submit (or block?)
    // Given user wants validation, let's warn.
    console.warn('Uno.valid.validateForm not available');
    return true;
  },
};

export default {
  Box,
  Notify,
  Valid,
};
