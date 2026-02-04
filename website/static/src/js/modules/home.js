export function HeroAnimation() {
  return {
    show: false,
    init() {
      setTimeout(() => {
        this.show = true;
      }, 100);
    },
  };
}
