# Security Policy

## Supported versions

Uno UI is pre-1.0. `0.10.x` is published to npm. Security fixes are applied to the current
`0.10.x` line; older versions are not maintained.

| Version | Supported |
| :--- | :--- |
| 0.10.x | Yes |
| < 0.10 | No |

## Reporting a vulnerability

Do **not** open a public GitHub issue for a security vulnerability.

Use GitHub's private vulnerability reporting for this repository:

<https://github.com/artsyne/uno-ui/security/advisories/new>

If private reporting is unavailable, contact the maintainers through
<https://github.com/artsyne> and ask for a private channel before sharing details.

Please include:

- Affected version/commit.
- A minimal reproduction (HTML, configuration, browser/Alpine/HTMX versions).
- Impact and the conditions required to trigger it.
- Any suggested fix or mitigation.

We will acknowledge the report and investigate on a best-effort basis. No response-time or
disclosure deadline is promised. Please give us a reasonable opportunity to fix the issue before
public disclosure.

## Scope

In scope:

- Box, Notify, Tooltip, Tab and Valid behavior.
- HTML sanitization and any path that inserts HTML.
- CSP nonce handling.
- The built `dist/` bundles and the source in `src/`.

Out of scope:

- Vulnerabilities in Alpine.js, Tailwind CSS, HTMX or the browser; report those upstream.
- Applications that pass untrusted data into APIs documented as trusted.
- Missing CSP headers, insecure server configuration or deployment mistakes.
- Social engineering and denial-of-service against hosted services.

## Security model and limitations

Uno is a client-side UI library. It does not provide authentication, authorization, a server-side
template engine or a general-purpose HTML sanitizer.

### HTML sanitization

`Uno.dom.sanitize()` is a regular-expression pass. It removes `<script>`, `<img>`, quoted `on*`
attributes and `javascript:` URLs. It is not an HTML parser and can be bypassed, especially with
unusual encoding or markup.

- Box `content` is passed through this sanitizer.
- Tooltip `content` and `target` innerHTML are passed through it.
- Box `prompt` message and default value are passed through it before interpolation. Do not pass
  untrusted values as the prompt default; use a dedicated sanitizer for untrusted input.
- Valid error messages are rendered as text.
- Notify `message` is rendered with `x-text`.

For untrusted HTML, use a maintained sanitizer such as DOMPurify on the server or in the client
before passing content to Uno.

### HTMX responses

Box, Notify and Tooltip insert HTMX responses as HTML. Uno does not sanitize those responses.
Treat every server response as trusted output and encode untrusted data on the server.

### Trusted HTML helper

`Uno.dom.setHTML(el, html, { trusted: true })` skips sanitization by design. Only use it with
content you control.

### CSP nonce

`csp.nonce` is written to Uno-created containers and to `<style>` elements created by
`Uno.dom.injectStyles()`. It does not make Alpine's standard expression evaluator compatible with
a policy that forbids `'unsafe-eval'`. See [CSP](docs/csp.md). Do not describe a nonce-only
deployment as a complete CSP.

### URLs and attributes

Uno does not validate URLs, `href` values or arbitrary HTML attributes in application content.
Encode and validate them in your application.

## Dependencies

The only runtime dependency is `alpinejs`. `htmx.org` is an optional peer dependency. Dependency
updates are reviewed as part of normal maintenance.

## Related

- [CSP guide](docs/csp.md)
- [Concepts: sanitization boundary](docs/concepts.md#sanitization-boundary)
- [Contributing](CONTRIBUTING.md)
