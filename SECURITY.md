# Security Policy

## Reporting a Vulnerability

We take the security of GoMeta Redirector seriously. If you discover a security vulnerability, please follow these steps:

### 🔒 Private Disclosure

**Do NOT open a public issue.** Instead, please report security vulnerabilities privately by emailing:

📧 **gyula@pilab.hu**

### What to Include

When reporting a vulnerability, please include:

- **Description** - Clear description of the vulnerability
- **Impact** - What an attacker could achieve
- **Reproduction Steps** - How to reproduce the issue
- **Affected Versions** - Which versions are affected
- **Suggested Fix** - If you have a fix or mitigation (optional)

### Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 1 week
- **Status Updates**: Weekly until resolved
- **Fix Release**: Depends on severity (see below)

### Severity Levels

| Severity | Description | Response Time |
|----------|-------------|---------------|
| Critical | Remote code execution, data breach | 24-48 hours |
| High | Authentication bypass, major vulnerability | 1 week |
| Medium | Limited impact vulnerabilities | 2-4 weeks |
| Low | Minor issues, information disclosure | Next release |

### Disclosure Policy

- We follow **responsible disclosure**
- We'll work with you to understand and resolve the issue
- We'll credit you in the security advisory (unless you prefer to remain anonymous)
- We'll publish a security advisory after the fix is released

### Supported Versions

| Version | Supported |
|---------|-----------|
| Latest main branch | ✅ |
| Latest release | ✅ |
| Older releases | ❌ |

### Security Best Practices

When deploying GoMeta Redirector:

- Keep the software updated to the latest version
- Use HTTPS in production (with valid TLS certificates)
- Restrict network access to the server
- Review and minimize the `repos.yaml` configuration
- Monitor logs for suspicious activity
- Run as a non-privileged user (not root)

### 🏆 Hall of Fame

We thank the following security researchers for responsibly disclosing vulnerabilities:

_No disclosures yet - be the first!_

---

Thank you for helping keep GoMeta Redirector secure! 🔒
