import {
  Navbar,
  Footer,
} from '../../../website/static/src/js/modules/layout.js';

describe('Layout Modules', () => {
  test('Navbar should return object with init function', () => {
    const navbar = Navbar();
    expect(typeof navbar).toBe('object');
    expect(typeof navbar.init).toBe('function');
  });

  test('Footer should exist', () => {
    const footer = Footer();
    expect(typeof footer.init).toBe('function');
  });
});
