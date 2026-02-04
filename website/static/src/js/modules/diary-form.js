export function DiaryForm(initialTheme) {
  return {
    theme: initialTheme || 16,
    setTheme(t) {
      this.theme = t;
    },
  };
}
