import { EntryMixin } from './entry.js';

export default function Drafts() {
  return {
    ...EntryMixin,

    init() {
      window.openEntryModal = (arg1, arg2, arg3) => {
        if (typeof arg1 === 'string' && arg1.startsWith('/entry/')) {
          this.openEntryModal(arg1, arg2);
          return;
        }

        if (arg3 === 'entry.edit') {
          this.openEditEntry(arg1, arg2);
        }
      };
    },
  };
}
